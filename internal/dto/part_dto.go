package dto

import (
	"bytes"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"7th-sense-server/internal/response"
)

// Lengths of the parts, part_processes and part_sub_programs columns (migration 016).
const (
	partNoMaxLength          = 100
	partNameMaxLength        = 200
	partProcessTextMaxLength = 100 // process, program name, drawing no., revision no.
	partSubProgramMaxLength  = 200
	partDrawingNameMax       = 255
	partMaxCount             = 100000
	partMaxProcesses         = 50
	partMaxSubPrograms       = 50
	partMaxMachines          = 100
	// partMaxSeconds caps every time field (999:59:59).
	partMaxSeconds = 999*3600 + 59*60 + 59
)

// PartDrawingMaxBytes is the largest part drawing PDF accepted (10 MB).
const PartDrawingMaxBytes = 10 << 20

// Part statuses accepted by UpdatePartRequest.
const (
	PartStatusActive   = "active"
	PartStatusInactive = "inactive"
)

// partDrawingFilePattern is the stored name of an uploaded drawing: 32 hex characters
// and .pdf (see service.PartService.UploadDrawing). Checking it here keeps a request
// from pointing a process at any other file.
var partDrawingFilePattern = regexp.MustCompile(`^[a-f0-9]{32}\.pdf$`)

// PartSubProgramFields is one sub program of a process.
type PartSubProgramFields struct {
	Name      string `json:"name"`
	CycleTime int    `json:"cycleTime"` // seconds
}

// PartProcessFields is one process of the Add / Edit Part form. ID names a saved
// process of the part when editing; nil adds a new one.
type PartProcessFields struct {
	ID                *int64                 `json:"id"`
	Process           string                 `json:"process"`
	ProgramName       string                 `json:"programName"`
	MachineIDs        []int64                `json:"machineIds"` // the customer's machines
	DrawingNo         string                 `json:"drawingNo"`
	DrawingRevisionNo string                 `json:"drawingRevisionNo"`
	DrawingFile       string                 `json:"drawingFile"` // from POST /parts/drawings
	DrawingName       string                 `json:"drawingName"`
	AccuracyID        int64                  `json:"accuracyId"`
	SetupTime         int                    `json:"setupTime"`   // seconds
	CycleTime         int                    `json:"cycleTime"`   // seconds
	LoadingTime       int                    `json:"loadingTime"` // seconds
	HasSubProgram     bool                   `json:"hasSubProgram"`
	SubPrograms       []PartSubProgramFields `json:"subPrograms"`
}

// PartFields are the Add / Edit Part form fields shared by create and update. The
// production times are not sent: the service adds them up from the processes.
type PartFields struct {
	UserID      int64               `json:"userId"`      // customer that owns the part
	ClienteleID *int64              `json:"clienteleId"` // optional; no clientele table yet
	PartNo      string              `json:"partNo"`
	PartName    string              `json:"partName"`
	PartCount   int                 `json:"partCount"`
	Processes   []PartProcessFields `json:"processes"`
}

// CreatePartRequest is the body for POST /parts.
type CreatePartRequest struct {
	PartFields
}

// Validate trims the fields and checks they are usable.
func (request *CreatePartRequest) Validate() error {
	return request.PartFields.validate()
}

// UpdatePartRequest is the body for PUT /parts/{id}. A blank status keeps the saved one.
type UpdatePartRequest struct {
	PartFields
	Status string `json:"status"`
}

// Validate trims the fields and checks they are usable.
func (request *UpdatePartRequest) Validate() error {
	if err := request.PartFields.validate(); err != nil {
		return err
	}
	request.Status = strings.ToLower(strings.TrimSpace(request.Status))
	if request.Status != "" && request.Status != PartStatusActive && request.Status != PartStatusInactive {
		return response.NewError(http.StatusBadRequest, "Please enter valid status")
	}
	return nil
}

func (fields *PartFields) validate() error {
	fields.PartNo = strings.TrimSpace(fields.PartNo)
	fields.PartName = strings.TrimSpace(fields.PartName)

	switch {
	case fields.UserID <= 0:
		return response.NewError(http.StatusBadRequest, "Please select customer")
	case fields.ClienteleID != nil && *fields.ClienteleID <= 0:
		return response.NewError(http.StatusBadRequest, "Please select valid clientele")
	case !validText(fields.PartNo, partNoMaxLength):
		return response.NewError(http.StatusBadRequest, "Please enter valid part no.")
	case !validText(fields.PartName, partNameMaxLength):
		return response.NewError(http.StatusBadRequest, "Please enter valid part name")
	case fields.PartCount < 1 || fields.PartCount > partMaxCount:
		return response.NewError(http.StatusBadRequest, "Please enter valid part count")
	case len(fields.Processes) == 0:
		return response.NewError(http.StatusBadRequest, "Please add at least one process")
	case len(fields.Processes) > partMaxProcesses:
		return response.NewError(http.StatusBadRequest, fmt.Sprintf("A part can have up to %d processes", partMaxProcesses))
	}

	seenIDs := make(map[int64]bool, len(fields.Processes))
	for index := range fields.Processes {
		process := &fields.Processes[index]
		if err := process.validate(index + 1); err != nil {
			return err
		}
		if process.ID != nil {
			if seenIDs[*process.ID] {
				return response.NewError(http.StatusBadRequest, "Please enter valid processes")
			}
			seenIDs[*process.ID] = true
		}
	}
	return nil
}

// validate checks process number position (1-based, used in the messages).
func (process *PartProcessFields) validate(position int) error {
	process.Process = strings.TrimSpace(process.Process)
	process.ProgramName = strings.TrimSpace(process.ProgramName)
	process.DrawingNo = strings.TrimSpace(process.DrawingNo)
	process.DrawingRevisionNo = strings.TrimSpace(process.DrawingRevisionNo)
	process.DrawingFile = strings.TrimSpace(process.DrawingFile)
	process.DrawingName = strings.TrimSpace(process.DrawingName)
	invalid := func(message string) error {
		return response.NewError(http.StatusBadRequest, fmt.Sprintf("%s (process %d)", message, position))
	}

	switch {
	case process.ID != nil && *process.ID <= 0:
		return invalid("Please enter valid process")
	case !validText(process.Process, partProcessTextMaxLength):
		return invalid("Please enter valid process")
	case !validText(process.ProgramName, partProcessTextMaxLength):
		return invalid("Please enter valid program name")
	case len(process.MachineIDs) == 0:
		return invalid("Please select applicable machine")
	case len(process.MachineIDs) > partMaxMachines:
		return invalid("Please select valid applicable machines")
	case !validText(process.DrawingNo, partProcessTextMaxLength):
		return invalid("Please enter valid drawing no.")
	case !validText(process.DrawingRevisionNo, partProcessTextMaxLength):
		return invalid("Please enter valid drawing revision no.")
	case !partDrawingFilePattern.MatchString(process.DrawingFile):
		return invalid("Please upload part drawing (PDF)")
	case !validText(process.DrawingName, partDrawingNameMax):
		return invalid("Please upload part drawing (PDF)")
	case process.AccuracyID <= 0:
		return invalid("Please select machine accuracy level")
	case process.SetupTime < 0 || process.SetupTime > partMaxSeconds:
		return invalid("Please enter valid setup time")
	case process.CycleTime <= 0 || process.CycleTime > partMaxSeconds:
		return invalid("Please enter valid cycle time")
	case process.LoadingTime < 0 || process.LoadingTime > partMaxSeconds:
		return invalid("Please enter valid loading / unloading time")
	}

	machineIDs, ok := uniquePositiveIDs(process.MachineIDs)
	if !ok {
		return invalid("Please select valid applicable machines")
	}
	process.MachineIDs = machineIDs

	if !process.HasSubProgram {
		process.SubPrograms = nil
		return nil
	}
	if len(process.SubPrograms) == 0 {
		return invalid("Please add at least one sub program")
	}
	if len(process.SubPrograms) > partMaxSubPrograms {
		return invalid(fmt.Sprintf("A process can have up to %d sub programs", partMaxSubPrograms))
	}
	for index := range process.SubPrograms {
		subProgram := &process.SubPrograms[index]
		subProgram.Name = strings.TrimSpace(subProgram.Name)
		if !validText(subProgram.Name, partSubProgramMaxLength) {
			return invalid("Please enter valid sub program name")
		}
		if subProgram.CycleTime < 0 || subProgram.CycleTime > partMaxSeconds {
			return invalid("Please enter valid sub program cycle time")
		}
	}
	return nil
}

// UploadPartDrawingRequest is the file of POST /parts/drawings (multipart field "file").
type UploadPartDrawingRequest struct {
	Name    string
	Content []byte
}

// Validate keeps the base name of the file and checks it is a PDF of at most
// PartDrawingMaxBytes.
func (request *UploadPartDrawingRequest) Validate() error {
	request.Name = strings.TrimSpace(filepath.Base(strings.ReplaceAll(request.Name, `\`, "/")))
	switch {
	case len(request.Content) == 0:
		return response.NewError(http.StatusBadRequest, "Please upload part drawing (PDF)")
	case len(request.Content) > PartDrawingMaxBytes:
		return response.NewError(http.StatusBadRequest, "Part drawing must be 10 MB or smaller")
	case !strings.EqualFold(filepath.Ext(request.Name), ".pdf") || !bytes.HasPrefix(request.Content, []byte("%PDF-")):
		return response.NewError(http.StatusBadRequest, "Part drawing must be a PDF file")
	case !validText(request.Name, partDrawingNameMax):
		return response.NewError(http.StatusBadRequest, "Please enter valid file name")
	}
	return nil
}

// ValidPartDrawingFile reports whether file is the stored name of an uploaded drawing.
func ValidPartDrawingFile(file string) bool {
	return partDrawingFilePattern.MatchString(file)
}

// validText reports whether value is 1 to maxLength characters without control characters.
func validText(value string, maxLength int) bool {
	if value == "" || utf8.RuneCountInString(value) > maxLength {
		return false
	}
	return !strings.ContainsFunc(value, unicode.IsControl)
}

// uniquePositiveIDs drops duplicate ids (keeping the order); ok is false when an id is not positive.
func uniquePositiveIDs(ids []int64) (unique []int64, ok bool) {
	unique = make([]int64, 0, len(ids))
	seen := make(map[int64]bool, len(ids))
	for _, id := range ids {
		if id <= 0 {
			return nil, false
		}
		if !seen[id] {
			seen[id] = true
			unique = append(unique, id)
		}
	}
	return unique, true
}
