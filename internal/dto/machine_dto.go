package dto

import (
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"7th-sense-server/internal/response"
)

var (
	machineNamePattern     = regexp.MustCompile(`^[a-zA-Z0-9 _./()#-]{1,50}$`)
	serialPattern          = regexp.MustCompile(`^[a-zA-Z0-9/_-]{1,50}$`)
	controllerVerPattern   = regexp.MustCompile(`^[a-zA-Z0-9._-]{1,10}$`)
	machineUsernamePattern = regexp.MustCompile(`^[a-zA-Z0-9_.]{1,15}$`)
	decimalPattern         = regexp.MustCompile(`^\d{1,8}(\.\d{1,2})?$`) // hourly_cost NUMERIC(10,2)
	healthIDsPattern       = regexp.MustCompile(`^\d+(,\d+)*$`)
)

const (
	machinePasswordMaxLength = 50
	machineAMCMaxLength      = 50
	machineHealthIDsMax      = 100
	dateLayout               = "2006-01-02"
)

// Machine statuses accepted by UpdateMachineRequest.
const (
	MachineStatusActive   = "active"
	MachineStatusInactive = "inactive"
)

// Machine system types.
const (
	MachineSystemTypeNewPlc = "NewPlc"
	MachineSystemTypeOldPlc = "OldPlc"
)

// Security modes.
const (
	SecurityModeNone           = "None"
	SecurityModeBasic256Sha256 = "Basic256Sha256"
)

// Controller version types (machines.controller_version_type).
const (
	ControllerVersionType828D = 1
	ControllerVersionType840D = 2
)

// MachineHealthAll in healthIds means every health sensor is monitored.
const MachineHealthAll = "All"

// MachineFields are the Add / Edit Machine form fields shared by create and update.
// Runtime columns written by the data collector (online status, mode, online and
// offline dates, …) are not part of the form and are never changed through the API.
type MachineFields struct {
	// Machine details
	UserID              int64  `json:"userId"`  // customer that owns the machine
	ModelID             int64  `json:"modelId"` // a model of the controller (models table)
	Name                string `json:"name"`
	SerialNumber        string `json:"serialNumber"`
	BrandID             int64  `json:"brandId"`
	ControllerID        int64  `json:"controllerId"`
	ControllerVersion   string `json:"controllerVersion"`
	AccuracyID          *int64 `json:"accuracyId"`
	LineID              *int64 `json:"lineId"` // one of the customer's lines (user_line)
	IsInsured           bool   `json:"isInsured"`
	MaintenanceRequired bool   `json:"maintenanceRequired"`
	InstallationDate    string `json:"installationDate"`
	BillDate            string `json:"billDate"`     // optional
	WarrantyDate        string `json:"warrantyDate"` // optional
	HourlyCost          string `json:"hourlyCost"`
	AMC                 string `json:"amc"`
	MachineSystemType   string `json:"machineSystemType"`
	SecurityMode        string `json:"securityMode"`
	// Data access
	LocalIP               string `json:"localIp"`
	LiveIP                string `json:"liveIp"`
	Port                  int    `json:"port"`
	Username              string `json:"username"`
	ControllerVersionType int    `json:"controllerVersionType"`
	Longitude             string `json:"longitude"`
	Latitude              string `json:"latitude"`
	HealthIDs             string `json:"healthIds"`
}

// CreateMachineRequest is the body for POST /machines.
type CreateMachineRequest struct {
	MachineFields
	Password string `json:"password"`
}

// Validate trims the fields and checks they are usable.
func (request *CreateMachineRequest) Validate() error {
	if err := request.MachineFields.validate(); err != nil {
		return err
	}
	if request.Password == "" || len(request.Password) > machinePasswordMaxLength {
		return response.NewError(http.StatusBadRequest, "Please enter valid password")
	}
	return nil
}

// UpdateMachineRequest is the body for PUT /machines/{id}. A blank password keeps the
// saved one; a blank status keeps the saved status.
type UpdateMachineRequest struct {
	MachineFields
	Password string `json:"password"`
	Status   string `json:"status"`
}

// Validate trims the fields and checks they are usable.
func (request *UpdateMachineRequest) Validate() error {
	if err := request.MachineFields.validate(); err != nil {
		return err
	}
	if len(request.Password) > machinePasswordMaxLength {
		return response.NewError(http.StatusBadRequest, "Please enter valid password")
	}
	request.Status = strings.ToLower(strings.TrimSpace(request.Status))
	if request.Status != "" && request.Status != MachineStatusActive && request.Status != MachineStatusInactive {
		return response.NewError(http.StatusBadRequest, "Please enter valid status")
	}
	return nil
}

func (fields *MachineFields) validate() error {
	fields.trim()

	switch {
	case fields.UserID <= 0:
		return response.NewError(http.StatusBadRequest, "Please select customer")
	case fields.ModelID <= 0:
		return response.NewError(http.StatusBadRequest, "Please select model")
	case !machineNamePattern.MatchString(fields.Name):
		return response.NewError(http.StatusBadRequest, "Please enter valid machine name")
	case !serialPattern.MatchString(fields.SerialNumber):
		return response.NewError(http.StatusBadRequest, "Please enter valid serial number")
	case fields.BrandID <= 0:
		return response.NewError(http.StatusBadRequest, "Please select brand")
	case fields.ControllerID <= 0:
		return response.NewError(http.StatusBadRequest, "Please select controller")
	case !controllerVerPattern.MatchString(fields.ControllerVersion):
		return response.NewError(http.StatusBadRequest, "Please enter valid controller software version")
	case fields.AccuracyID != nil && *fields.AccuracyID <= 0:
		return response.NewError(http.StatusBadRequest, "Please select valid machine accuracy")
	case fields.LineID != nil && *fields.LineID <= 0:
		return response.NewError(http.StatusBadRequest, "Please select valid line")
	case !validDate(fields.InstallationDate):
		return response.NewError(http.StatusBadRequest, "Please enter valid installation date")
	case fields.BillDate != "" && !validDate(fields.BillDate):
		return response.NewError(http.StatusBadRequest, "Please enter valid bill date")
	case fields.WarrantyDate != "" && !validDate(fields.WarrantyDate):
		return response.NewError(http.StatusBadRequest, "Please enter valid warranty end date")
	case fields.WarrantyDate != "" && fields.WarrantyDate < fields.InstallationDate: // YYYY-MM-DD sorts as text
		return response.NewError(http.StatusBadRequest, "Warranty end date cannot be before the installation date")
	case !decimalPattern.MatchString(fields.HourlyCost):
		return response.NewError(http.StatusBadRequest, "Please enter valid hourly cost")
	case len(fields.AMC) > machineAMCMaxLength:
		return response.NewError(http.StatusBadRequest, "Please enter valid AMC")
	case fields.MachineSystemType != MachineSystemTypeNewPlc && fields.MachineSystemType != MachineSystemTypeOldPlc:
		return response.NewError(http.StatusBadRequest, "Please select valid machine system type")
	case fields.SecurityMode != SecurityModeNone && fields.SecurityMode != SecurityModeBasic256Sha256:
		return response.NewError(http.StatusBadRequest, "Please select valid security type")
	case !validIPv4(fields.LocalIP):
		return response.NewError(http.StatusBadRequest, "Please enter valid local IP")
	case !validIPv4(fields.LiveIP):
		return response.NewError(http.StatusBadRequest, "Please enter valid live IP")
	case fields.Port < 1 || fields.Port > 65535:
		return response.NewError(http.StatusBadRequest, "Please enter valid port")
	case !machineUsernamePattern.MatchString(fields.Username):
		return response.NewError(http.StatusBadRequest, "Please enter valid user name")
	case fields.ControllerVersionType != ControllerVersionType828D && fields.ControllerVersionType != ControllerVersionType840D:
		return response.NewError(http.StatusBadRequest, "Please select controller version")
	case fields.Longitude != "" && !validCoordinate(fields.Longitude, 180):
		return response.NewError(http.StatusBadRequest, "Please enter valid longitude")
	case fields.Latitude != "" && !validCoordinate(fields.Latitude, 90):
		return response.NewError(http.StatusBadRequest, "Please enter valid latitude")
	case !validHealthIDs(fields.HealthIDs):
		return response.NewError(http.StatusBadRequest, "Please select valid machine health")
	}
	return nil
}

func (fields *MachineFields) trim() {
	fields.Name = strings.TrimSpace(fields.Name)
	fields.SerialNumber = strings.TrimSpace(fields.SerialNumber)
	fields.ControllerVersion = strings.TrimSpace(fields.ControllerVersion)
	fields.InstallationDate = strings.TrimSpace(fields.InstallationDate)
	fields.BillDate = strings.TrimSpace(fields.BillDate)
	fields.WarrantyDate = strings.TrimSpace(fields.WarrantyDate)
	fields.HourlyCost = strings.TrimSpace(fields.HourlyCost)
	fields.AMC = strings.TrimSpace(fields.AMC)
	fields.MachineSystemType = strings.TrimSpace(fields.MachineSystemType)
	fields.SecurityMode = strings.TrimSpace(fields.SecurityMode)
	fields.LocalIP = strings.TrimSpace(fields.LocalIP)
	fields.LiveIP = strings.TrimSpace(fields.LiveIP)
	fields.Username = strings.TrimSpace(fields.Username)
	fields.Longitude = strings.TrimSpace(fields.Longitude)
	fields.Latitude = strings.TrimSpace(fields.Latitude)
	fields.HealthIDs = strings.ReplaceAll(strings.TrimSpace(fields.HealthIDs), " ", "")
	if fields.MachineSystemType == "" {
		fields.MachineSystemType = MachineSystemTypeNewPlc
	}
	if fields.SecurityMode == "" {
		fields.SecurityMode = SecurityModeBasic256Sha256
	}
}

// validDate reports whether value is a real calendar date in YYYY-MM-DD form.
func validDate(value string) bool {
	_, err := time.Parse(dateLayout, value)
	return err == nil
}

func validIPv4(value string) bool {
	ip := net.ParseIP(value)
	return ip != nil && ip.To4() != nil && strings.Count(value, ".") == 3
}

// validCoordinate reports whether value is a decimal number between -limit and limit.
func validCoordinate(value string, limit float64) bool {
	number, err := strconv.ParseFloat(value, 64)
	return err == nil && number >= -limit && number <= limit && len(value) <= 50
}

// validHealthIDs accepts "", "All" or comma-separated health sensor ids ("1,5,6").
func validHealthIDs(value string) bool {
	if value == "" || value == MachineHealthAll {
		return true
	}
	return len(value) <= machineHealthIDsMax && healthIDsPattern.MatchString(value)
}
