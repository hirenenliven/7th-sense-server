package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"7th-sense-server/internal/db"
	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/response"
)

// ProductionScheduleService holds the business rules for production schedules. A
// schedule belongs to a customer (user_id), plans one of its machines for one day or a
// range of days and is saved with its items (jobs) in one transaction. Everything an
// item names must be the customer's: a supervisor and an operator of the company, a
// shift, a part and one of its processes that runs on the machine. A job runs inside its
// shift, and a machine has one schedule per day.
//
// Access is the same as parts (companyOwner / ownCompanyCustomer): the superadmin works
// on every customer's schedules, everyone else reads their company's and managers
// (managersOnly routes) change them.
type ProductionScheduleService struct {
	dbPool             *pgxpool.Pool
	scheduleRepository *repository.ProductionScheduleRepository
	shiftRepository    *repository.ShiftRepository
	userRepository     *repository.UserRepository
	machineService     *MachineService
}

// NewProductionScheduleService creates a ProductionScheduleService that stores schedules
// in dbPool and checks customers with userRepository and machines with machineService.
func NewProductionScheduleService(dbPool *pgxpool.Pool, userRepository *repository.UserRepository, machineService *MachineService) *ProductionScheduleService {
	return &ProductionScheduleService{
		dbPool:             dbPool,
		scheduleRepository: repository.NewProductionScheduleRepository(dbPool),
		shiftRepository:    repository.NewShiftRepository(dbPool),
		userRepository:     userRepository,
		machineService:     machineService,
	}
}

// GetAll returns the schedules actor may read (see companyOwner), without their items.
// A non-nil customerID keeps only that customer's schedules.
func (scheduleService *ProductionScheduleService) GetAll(ctx context.Context, actor model.AuthUser, customerID *int64) ([]model.ProductionSchedule, error) {
	ownerID, err := companyOwner(actor, customerID)
	if err != nil {
		return nil, err
	}
	return scheduleService.scheduleRepository.FindAll(ctx, ownerID)
}

// GetByID returns the schedule with its items when actor may read it, else
// "Production schedule not found".
func (scheduleService *ProductionScheduleService) GetByID(ctx context.Context, actor model.AuthUser, customerID *int64, id int64) (model.ProductionSchedule, error) {
	ownerID, err := companyOwner(actor, customerID)
	if err != nil {
		return model.ProductionSchedule{}, err
	}
	return scheduleDetails(ctx, scheduleService.scheduleRepository, ownerID, id)
}

// Create adds a schedule with its items. A customer or admin (actor) adds it to their
// own company, whatever customer the request names; the superadmin to the customer named.
func (scheduleService *ProductionScheduleService) Create(ctx context.Context, actor model.AuthUser, scheduleRequest dto.CreateProductionScheduleRequest) (model.ProductionSchedule, error) {
	if err := ownCompanyCustomer(actor, &scheduleRequest.UserID); err != nil {
		return model.ProductionSchedule{}, err
	}
	for _, item := range scheduleRequest.Items {
		if item.ID != nil {
			return model.ProductionSchedule{}, response.NewError(http.StatusBadRequest, "Please enter valid jobs")
		}
	}
	if err := scheduleService.checkReferences(ctx, scheduleRequest.ProductionScheduleFields, 0); err != nil {
		return model.ProductionSchedule{}, err
	}
	schedule, items := scheduleFromFields(scheduleRequest.ProductionScheduleFields)
	schedule.Status = dto.ScheduleStatusActive

	var created model.ProductionSchedule
	err := db.WithTransaction(ctx, scheduleService.dbPool, func(tx pgx.Tx) error {
		scheduleRepository := repository.NewProductionScheduleRepository(tx)

		id, err := scheduleRepository.Create(ctx, schedule)
		if err != nil {
			return toProductionScheduleResponseError(err)
		}
		if err := saveScheduleItems(ctx, scheduleRepository, id, items); err != nil {
			return err
		}
		created, err = scheduleDetails(ctx, scheduleRepository, nil, id)
		return err
	})
	return created, err
}

// Update replaces the schedule and its items. Items with an id are updated (they must
// be the schedule's), items without one are added and saved items missing from the
// request are deleted. A blank status keeps the saved one. A customer or admin (actor)
// changes only their company's schedules and cannot move one to another company.
func (scheduleService *ProductionScheduleService) Update(ctx context.Context, actor model.AuthUser, id int64, scheduleRequest dto.UpdateProductionScheduleRequest) (model.ProductionSchedule, error) {
	if err := ownCompanyCustomer(actor, &scheduleRequest.UserID); err != nil {
		return model.ProductionSchedule{}, err
	}
	ownerID, err := companyOwner(actor, nil)
	if err != nil {
		return model.ProductionSchedule{}, err
	}
	if _, err := scheduleService.scheduleRepository.FindByID(ctx, ownerID, id); err != nil {
		return model.ProductionSchedule{}, toProductionScheduleResponseError(err)
	}
	if err := scheduleService.checkReferences(ctx, scheduleRequest.ProductionScheduleFields, id); err != nil {
		return model.ProductionSchedule{}, err
	}
	schedule, items := scheduleFromFields(scheduleRequest.ProductionScheduleFields)
	schedule.ID = id
	schedule.Status = scheduleRequest.Status

	var updated model.ProductionSchedule
	err = db.WithTransaction(ctx, scheduleService.dbPool, func(tx pgx.Tx) error {
		scheduleRepository := repository.NewProductionScheduleRepository(tx)

		if err := scheduleRepository.Update(ctx, schedule); err != nil {
			return toProductionScheduleResponseError(err)
		}
		keepIDs := make([]int64, 0, len(items))
		for _, item := range items {
			if item.ID != 0 {
				keepIDs = append(keepIDs, item.ID)
			}
		}
		if err := scheduleRepository.DeleteItemsExcept(ctx, id, keepIDs); err != nil {
			return err
		}
		if err := saveScheduleItems(ctx, scheduleRepository, id, items); err != nil {
			return err
		}
		updated, err = scheduleDetails(ctx, scheduleRepository, nil, id)
		return err
	})
	return updated, err
}

// Delete soft-deletes the schedule. A customer or admin (actor) deletes only their
// company's schedules ("Production schedule not found" for others).
func (scheduleService *ProductionScheduleService) Delete(ctx context.Context, actor model.AuthUser, id int64) error {
	ownerID, err := companyOwner(actor, nil)
	if err != nil {
		return err
	}
	return toProductionScheduleResponseError(scheduleService.scheduleRepository.Delete(ctx, ownerID, id))
}

// checkReferences returns a "Please select valid …" error unless the customer exists,
// the machine is the customer's and free on the schedule's days (another schedule than
// excludeID), and every item names the customer's supervisor, operator, shift, part and
// a process of that part that runs on the machine, inside its shift.
func (scheduleService *ProductionScheduleService) checkReferences(ctx context.Context, fields dto.ProductionScheduleFields, excludeID int64) error {
	isCustomer, err := scheduleService.userRepository.IsCustomer(ctx, fields.UserID)
	if err != nil {
		return err
	}
	if !isCustomer {
		return response.NewError(http.StatusBadRequest, "Please select valid customer")
	}
	if err := scheduleService.machineService.CheckOwnedBy(ctx, fields.UserID, []int64{fields.MachineID}); err != nil {
		return asInvalidSelection(err, "Please select valid machines", "Please select a machine of the selected customer")
	}
	booked, err := scheduleService.scheduleRepository.MachineBooked(ctx, fields.MachineID, fields.StartDate, fields.EndDate, excludeID)
	if err != nil {
		return err
	}
	if booked {
		return response.NewError(http.StatusBadRequest, "This machine is already scheduled on the selected date(s)")
	}

	var supervisorIDs, operatorIDs, partIDs, pairParts, pairProcesses []int64
	seen := make(map[string]bool)
	addOnce := func(list *[]int64, kind string, id int64) bool {
		key := fmt.Sprintf("%s:%d", kind, id)
		if seen[key] {
			return false
		}
		seen[key] = true
		*list = append(*list, id)
		return true
	}
	for _, item := range fields.Items {
		addOnce(&supervisorIDs, "supervisor", item.SupervisorID)
		addOnce(&operatorIDs, "operator", item.OperatorID)
		addOnce(&partIDs, "part", item.PartID)
		if addOnce(&pairProcesses, fmt.Sprintf("process:%d", item.PartID), item.ProcessID) {
			pairParts = append(pairParts, item.PartID)
		}
	}

	checks := []struct {
		count   func() (int, error)
		want    int
		message string
	}{
		{func() (int, error) {
			return scheduleService.scheduleRepository.CountCompanyMembers(ctx, fields.UserID, model.RoleSupervisor, supervisorIDs)
		}, len(supervisorIDs), "Please select a supervisor / programmer of the selected customer"},
		{func() (int, error) {
			return scheduleService.scheduleRepository.CountCompanyMembers(ctx, fields.UserID, model.RoleOperator, operatorIDs)
		}, len(operatorIDs), "Please select an operator of the selected customer"},
		{func() (int, error) {
			return scheduleService.scheduleRepository.CountCustomerParts(ctx, fields.UserID, partIDs)
		}, len(partIDs), "Please select a part of the selected customer"},
		{func() (int, error) {
			return scheduleService.scheduleRepository.CountMachineProcesses(ctx, fields.MachineID, pairParts, pairProcesses)
		}, len(pairProcesses), "Please select a process of the part that runs on the selected machine"},
	}
	for _, check := range checks {
		count, err := check.count()
		if err != nil {
			return err
		}
		if count != check.want {
			return response.NewError(http.StatusBadRequest, check.message)
		}
	}

	shifts, err := scheduleService.shiftRepository.FindAll(ctx, &fields.UserID)
	if err != nil {
		return err
	}
	shiftsByID := make(map[int64]model.Shift, len(shifts))
	for _, shift := range shifts {
		shiftsByID[shift.ID] = shift
	}
	spans := make([]jobSpan, len(fields.Items))
	for index, item := range fields.Items {
		shift, ok := shiftsByID[item.ShiftID]
		if !ok {
			return response.NewError(http.StatusBadRequest, fmt.Sprintf("Please select a shift of the selected customer (job %d)", index+1))
		}
		span, ok := shiftSpan(shift, item, index+1)
		if !ok {
			return response.NewError(http.StatusBadRequest, fmt.Sprintf(
				"Start and end time must be inside the shift %s (%s – %s) (job %d)", shift.Name, shift.StartTime, shift.EndTime, index+1))
		}
		spans[index] = span
	}
	return checkJobsOverlap(spans)
}

// minutesPerDay is the minutes of a day.
const minutesPerDay = 24 * 60

// jobSpan is when job position runs on the machine, in minutes since 1970-01-01.
type jobSpan struct{ start, end, position int }

// shiftSpan returns when item runs: its schedule date is the day its shift starts, so a
// 02:00 job of a 22:00 – 06:00 shift runs the next morning. ok is false when start … end
// ("HH:MM", past midnight when end is not after start) is not inside the shift.
func shiftSpan(shift model.Shift, item dto.ProductionScheduleItemFields, position int) (jobSpan, bool) {
	shiftLength := dto.MinutesBetween(shift.StartTime, shift.EndTime)
	startOffset := 0
	if item.StartTime != shift.StartTime {
		startOffset = dto.MinutesBetween(shift.StartTime, item.StartTime)
	}
	endOffset := startOffset + dto.MinutesBetween(item.StartTime, item.EndTime)
	shiftStart := dto.ScheduleDay(item.ScheduleDate)*minutesPerDay + dto.MinutesBetween("00:00", shift.StartTime)%minutesPerDay
	return jobSpan{shiftStart + startOffset, shiftStart + endOffset, position},
		startOffset < shiftLength && endOffset <= shiftLength
}

// checkJobsOverlap refuses two jobs that run on the machine at the same time.
func checkJobsOverlap(spans []jobSpan) error {
	sorted := slices.Clone(spans)
	slices.SortFunc(sorted, func(a, b jobSpan) int { return a.start - b.start })
	for index := 1; index < len(sorted); index++ {
		if sorted[index].start < sorted[index-1].end {
			first, second := min(sorted[index-1].position, sorted[index].position), max(sorted[index-1].position, sorted[index].position)
			return response.NewError(http.StatusBadRequest, fmt.Sprintf("Job %d and job %d overlap on the machine", first, second))
		}
	}
	return nil
}

// scheduleDetails returns schedule id of ownerID with its items.
func scheduleDetails(ctx context.Context, scheduleRepository *repository.ProductionScheduleRepository, ownerID *int64, id int64) (model.ProductionSchedule, error) {
	schedule, err := scheduleRepository.FindByID(ctx, ownerID, id)
	if err != nil {
		return model.ProductionSchedule{}, toProductionScheduleResponseError(err)
	}
	items, err := scheduleRepository.FindItems(ctx, id)
	if err != nil {
		return model.ProductionSchedule{}, err
	}
	schedule.Items = items
	return schedule, nil
}

// saveScheduleItems updates the items with an id and adds the others to scheduleID. An
// item id that is not one of the schedule's is a bad request.
func saveScheduleItems(ctx context.Context, scheduleRepository *repository.ProductionScheduleRepository, scheduleID int64, items []model.ProductionScheduleItem) error {
	for _, item := range items {
		item.ScheduleID = scheduleID
		if item.ID == 0 {
			if err := scheduleRepository.CreateItem(ctx, item); err != nil {
				return err
			}
			continue
		}
		if err := scheduleRepository.UpdateItem(ctx, item); err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return response.NewError(http.StatusBadRequest, "Please enter valid jobs")
			}
			return err
		}
	}
	return nil
}

// scheduleFromFields copies the form fields into a schedule and its items (in form
// order) and works out each item's minutes.
func scheduleFromFields(fields dto.ProductionScheduleFields) (model.ProductionSchedule, []model.ProductionScheduleItem) {
	schedule := model.ProductionSchedule{
		UserID:       fields.UserID,
		MachineID:    fields.MachineID,
		ScheduleType: fields.ScheduleType,
		StartDate:    fields.StartDate,
		EndDate:      fields.EndDate,
	}
	items := make([]model.ProductionScheduleItem, len(fields.Items))
	for index, itemFields := range fields.Items {
		item := model.ProductionScheduleItem{
			SortOrder:    index + 1,
			ScheduleDate: itemFields.ScheduleDate,
			JobID:        itemFields.JobID,
			SupervisorID: itemFields.SupervisorID,
			OperatorID:   itemFields.OperatorID,
			ShiftID:      itemFields.ShiftID,
			PartID:       itemFields.PartID,
			ProcessID:    itemFields.ProcessID,
			Quantity:     itemFields.Quantity,
			StartTime:    itemFields.StartTime,
			EndTime:      itemFields.EndTime,
			TotalMinutes: dto.MinutesBetween(itemFields.StartTime, itemFields.EndTime),
		}
		if itemFields.ID != nil {
			item.ID = *itemFields.ID
		}
		items[index] = item
	}
	return schedule, items
}

// toProductionScheduleResponseError maps repository errors to schedule responses. Other
// errors pass through and are sent as ServerError.
func toProductionScheduleResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.NewError(http.StatusNotFound, "Production schedule not found")
	case errors.Is(err, repository.ErrInvalidUser):
		return response.NewError(http.StatusBadRequest, "Please select valid customer")
	default:
		return err
	}
}
