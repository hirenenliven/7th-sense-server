package service

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"7th-sense-server/internal/db"
	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/response"
)

// PartService holds the business rules for parts. A part belongs to a customer
// (user_id) and is saved with its processes, their applicable machines and sub
// programs in one transaction. The clientele and every applicable machine must be the customer's
// and every drawing must have been uploaded (UploadDrawing).
//
// The superadmin reads and changes every customer's parts; everyone else reads their
// company's, and managers (customer, admin — managersOnly routes) change them.
//
// Production times are worked out here, never taken from the request:
// process = setup + cycle + loading + its sub programs' cycle times; part = its processes.
type PartService struct {
	dbPool           *pgxpool.Pool
	partRepository   *repository.PartRepository
	userRepository   *repository.UserRepository
	machineService   *MachineService
	clienteleService *ClienteleService
	drawingStorage   *repository.PartDrawingStorage
}

// NewPartService creates a PartService that stores parts in dbPool, drawings in
// drawingStorage, and checks customers, machines and clienteles with the other arguments.
func NewPartService(
	dbPool *pgxpool.Pool,
	userRepository *repository.UserRepository,
	machineService *MachineService,
	clienteleService *ClienteleService,
	drawingStorage *repository.PartDrawingStorage,
) *PartService {
	return &PartService{
		dbPool:           dbPool,
		partRepository:   repository.NewPartRepository(dbPool),
		userRepository:   userRepository,
		machineService:   machineService,
		clienteleService: clienteleService,
		drawingStorage:   drawingStorage,
	}
}

// GetAll returns the parts actor may read (see companyOwner), without their processes.
// A non-nil customerID keeps only that customer's parts.
func (partService *PartService) GetAll(ctx context.Context, actor model.AuthUser, customerID *int64) ([]model.Part, error) {
	ownerID, err := companyOwner(actor, customerID)
	if err != nil {
		return nil, err
	}
	return partService.partRepository.FindAll(ctx, ownerID)
}

// GetByID returns the part with its processes when actor may read it, else
// "Part not found".
func (partService *PartService) GetByID(ctx context.Context, actor model.AuthUser, customerID *int64, id int64) (model.Part, error) {
	ownerID, err := companyOwner(actor, customerID)
	if err != nil {
		return model.Part{}, err
	}
	return partService.details(ctx, partService.partRepository, ownerID, id)
}

// details returns part id of ownerID with its processes and their sub programs.
func (partService *PartService) details(ctx context.Context, partRepository *repository.PartRepository, ownerID *int64, id int64) (model.Part, error) {
	part, err := partRepository.FindByID(ctx, ownerID, id)
	if err != nil {
		return model.Part{}, toPartResponseError(err)
	}
	processes, err := partRepository.FindProcesses(ctx, id)
	if err != nil {
		return model.Part{}, err
	}
	subPrograms, err := partRepository.FindSubPrograms(ctx, id)
	if err != nil {
		return model.Part{}, err
	}

	byProcess := make(map[int64][]model.PartSubProgram, len(processes))
	for _, subProgram := range subPrograms {
		byProcess[subProgram.ProcessID] = append(byProcess[subProgram.ProcessID], subProgram)
	}
	for index := range processes {
		processes[index].SubPrograms = byProcess[processes[index].ID]
		if processes[index].SubPrograms == nil {
			processes[index].SubPrograms = []model.PartSubProgram{}
		}
	}
	part.Processes = processes
	return part, nil
}

// Create adds a part with its processes. A customer or admin (actor) adds it to their
// own company, whatever customer the request names; the superadmin to the customer named.
func (partService *PartService) Create(ctx context.Context, actor model.AuthUser, partRequest dto.CreatePartRequest) (model.Part, error) {
	if err := ownCompanyCustomer(actor, &partRequest.UserID); err != nil {
		return model.Part{}, err
	}
	if err := partService.checkReferences(ctx, partRequest.PartFields); err != nil {
		return model.Part{}, err
	}
	for _, process := range partRequest.Processes {
		if process.ID != nil {
			return model.Part{}, response.NewError(http.StatusBadRequest, "Please enter valid processes")
		}
	}
	part, processes := partFromFields(partRequest.PartFields)
	part.Status = dto.PartStatusActive

	var created model.Part
	err := db.WithTransaction(ctx, partService.dbPool, func(tx pgx.Tx) error {
		partRepository := repository.NewPartRepository(tx)

		id, err := partRepository.Create(ctx, part)
		if err != nil {
			return toPartResponseError(err)
		}
		if err := saveProcesses(ctx, partRepository, id, processes); err != nil {
			return err
		}
		created, err = partService.details(ctx, partRepository, nil, id)
		return err
	})
	return created, err
}

// Update replaces the part and its processes. Processes with an id are updated (they
// must be the part's), processes without one are added and saved processes missing
// from the request are deleted. A blank status keeps the saved one. A customer or
// admin (actor) changes only their company's parts ("Part not found" for others) and
// cannot move one to another company.
func (partService *PartService) Update(ctx context.Context, actor model.AuthUser, id int64, partRequest dto.UpdatePartRequest) (model.Part, error) {
	if err := ownCompanyCustomer(actor, &partRequest.UserID); err != nil {
		return model.Part{}, err
	}
	ownerID, err := companyOwner(actor, nil)
	if err != nil {
		return model.Part{}, err
	}
	if _, err := partService.partRepository.FindByID(ctx, ownerID, id); err != nil {
		return model.Part{}, toPartResponseError(err)
	}
	if err := partService.checkReferences(ctx, partRequest.PartFields); err != nil {
		return model.Part{}, err
	}
	part, processes := partFromFields(partRequest.PartFields)
	part.ID = id
	part.Status = partRequest.Status

	var updated model.Part
	err = db.WithTransaction(ctx, partService.dbPool, func(tx pgx.Tx) error {
		partRepository := repository.NewPartRepository(tx)

		if err := partRepository.Update(ctx, part); err != nil {
			return toPartResponseError(err)
		}
		keepIDs := make([]int64, 0, len(processes))
		for _, process := range processes {
			if process.ID != 0 {
				keepIDs = append(keepIDs, process.ID)
			}
		}
		if err := partRepository.DeleteProcessesExcept(ctx, id, keepIDs); err != nil {
			return err
		}
		if err := saveProcesses(ctx, partRepository, id, processes); err != nil {
			return err
		}
		updated, err = partService.details(ctx, partRepository, nil, id)
		return err
	})
	return updated, err
}

// Delete soft-deletes the part. A customer or admin (actor) deletes only their
// company's parts ("Part not found" for others).
func (partService *PartService) Delete(ctx context.Context, actor model.AuthUser, id int64) error {
	ownerID, err := companyOwner(actor, nil)
	if err != nil {
		return err
	}
	return toPartResponseError(partService.partRepository.Delete(ctx, ownerID, id))
}

// UploadDrawing stores a part drawing PDF and returns the name a process refers to.
func (partService *PartService) UploadDrawing(ctx context.Context, drawingRequest dto.UploadPartDrawingRequest) (model.PartDrawing, error) {
	if err := ctx.Err(); err != nil {
		return model.PartDrawing{}, err
	}
	file, err := partService.drawingStorage.Save(drawingRequest.Content)
	if err != nil {
		return model.PartDrawing{}, err
	}
	return model.PartDrawing{File: file, Name: drawingRequest.Name}, nil
}

// OpenDrawing opens the drawing of process processID of part id when actor may read
// the part. The caller closes the file.
func (partService *PartService) OpenDrawing(ctx context.Context, actor model.AuthUser, id, processID int64) (*os.File, string, error) {
	part, err := partService.GetByID(ctx, actor, nil, id)
	if err != nil {
		return nil, "", err
	}
	for _, process := range part.Processes {
		if process.ID != processID {
			continue
		}
		if !dto.ValidPartDrawingFile(process.DrawingFile) {
			return nil, "", partDrawingNotFound
		}
		file, err := partService.drawingStorage.Open(process.DrawingFile)
		if errors.Is(err, repository.ErrNotFound) {
			return nil, "", partDrawingNotFound
		}
		return file, process.DrawingName, err
	}
	return nil, "", response.NewError(http.StatusNotFound, "Process not found")
}

var partDrawingNotFound = response.NewError(http.StatusNotFound, "Part drawing not found")

// checkReferences returns a "Please select valid …" error unless the customer exists,
// the clientele (when given) and every applicable machine are the customer's and every
// drawing was uploaded.
func (partService *PartService) checkReferences(ctx context.Context, fields dto.PartFields) error {
	isCustomer, err := partService.userRepository.IsCustomer(ctx, fields.UserID)
	if err != nil {
		return err
	}
	if !isCustomer {
		return response.NewError(http.StatusBadRequest, "Please select valid customer")
	}
	if fields.ClienteleID != nil {
		if err := partService.clienteleService.CheckOwnedBy(ctx, fields.UserID, *fields.ClienteleID); err != nil {
			return err
		}
	}

	var machineIDs []int64
	seen := make(map[int64]bool)
	for _, process := range fields.Processes {
		for _, machineID := range process.MachineIDs {
			if !seen[machineID] {
				seen[machineID] = true
				machineIDs = append(machineIDs, machineID)
			}
		}
	}
	if err := partService.machineService.CheckOwnedBy(ctx, fields.UserID, machineIDs); err != nil {
		return asInvalidSelection(err, "Please select valid machines", "Please select applicable machines of the selected customer")
	}

	for _, process := range fields.Processes {
		exists, err := partService.drawingStorage.Exists(process.DrawingFile)
		if err != nil {
			return err
		}
		if !exists {
			return response.NewError(http.StatusBadRequest, "Please upload part drawing (PDF) again")
		}
	}
	return nil
}

// saveProcesses updates the processes with an id and adds the others to partID, each
// with its applicable machines and sub programs. A process id that is not one of the
// part's is a bad request.
func saveProcesses(ctx context.Context, partRepository *repository.PartRepository, partID int64, processes []model.PartProcess) error {
	for _, process := range processes {
		process.PartID = partID
		if process.ID == 0 {
			id, err := partRepository.CreateProcess(ctx, process)
			if err != nil {
				return err
			}
			process.ID = id
		} else if err := partRepository.UpdateProcess(ctx, process); err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return response.NewError(http.StatusBadRequest, "Please enter valid processes")
			}
			return err
		}

		if err := partRepository.ReplaceProcessMachines(ctx, process.ID, process.MachineIDs); err != nil {
			return err
		}
		if err := partRepository.ReplaceSubPrograms(ctx, process.ID, process.SubPrograms); err != nil {
			return err
		}
	}
	return nil
}

// partFromFields copies the form fields into a part and its processes (in form order)
// and works out the production times.
func partFromFields(fields dto.PartFields) (model.Part, []model.PartProcess) {
	part := model.Part{
		UserID:      fields.UserID,
		ClienteleID: fields.ClienteleID,
		PartNo:      fields.PartNo,
		PartName:    fields.PartName,
		PartCount:   fields.PartCount,
	}

	processes := make([]model.PartProcess, len(fields.Processes))
	for index, processFields := range fields.Processes {
		subPrograms := make([]model.PartSubProgram, len(processFields.SubPrograms))
		subProgramTime := 0
		for subIndex, subProgram := range processFields.SubPrograms {
			subPrograms[subIndex] = model.PartSubProgram{Name: subProgram.Name, CycleTime: subProgram.CycleTime}
			subProgramTime += subProgram.CycleTime
		}

		process := model.PartProcess{
			SortOrder:         index + 1,
			Process:           processFields.Process,
			ProgramName:       processFields.ProgramName,
			DrawingNo:         processFields.DrawingNo,
			DrawingRevisionNo: processFields.DrawingRevisionNo,
			DrawingFile:       processFields.DrawingFile,
			DrawingName:       processFields.DrawingName,
			AccuracyID:        processFields.AccuracyID,
			SetupTime:         processFields.SetupTime,
			CycleTime:         processFields.CycleTime,
			LoadingTime:       processFields.LoadingTime,
			ProductionTime:    processFields.SetupTime + processFields.CycleTime + processFields.LoadingTime + subProgramTime,
			HasSubProgram:     processFields.HasSubProgram,
			MachineIDs:        processFields.MachineIDs,
			SubPrograms:       subPrograms,
		}
		if processFields.ID != nil {
			process.ID = *processFields.ID
		}
		processes[index] = process
		part.ProductionTime += process.ProductionTime
	}
	return part, processes
}

// toPartResponseError maps repository errors to part responses. Other errors pass
// through and are sent as ServerError.
func toPartResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.NewError(http.StatusNotFound, "Part not found")
	case errors.Is(err, repository.ErrDuplicateName):
		return response.NewError(http.StatusBadRequest, "This customer already has a part with this part no. and name")
	case errors.Is(err, repository.ErrInvalidUser):
		return response.NewError(http.StatusBadRequest, "Please select valid customer")
	case errors.Is(err, repository.ErrInvalidClientele):
		return response.NewError(http.StatusBadRequest, "Please select valid clientele")
	default:
		return err
	}
}
