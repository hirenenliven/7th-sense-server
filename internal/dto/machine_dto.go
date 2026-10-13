package dto

import (
	"regexp"
	"strings"

	"7th-sense-server/internal/response"
)

var (
	ipPattern              = regexp.MustCompile(`^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`)
	namePattern            = regexp.MustCompile(`^[a-zA-Z0-9_\s-]{1,50}$`)
	serialPattern          = regexp.MustCompile(`^[a-zA-Z0-9-]{1,50}$`)
	machineUsernamePattern = regexp.MustCompile(`^[a-zA-Z0-9_.]{1,15}$`)
	datePattern            = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	decimalPattern         = regexp.MustCompile(`^\d+(\.\d{1,2})?$`)
	coordinatePattern      = regexp.MustCompile(`^-?\d+\.?\d*$`)
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

// CreateMachineRequest is the body for POST /machines.
type CreateMachineRequest struct {
	UserID                *int64  `json:"userId"`
	MachineTypeID         *int64  `json:"machineTypeId"`
	ProductID             *int64  `json:"productId"`
	Name                  string  `json:"name"`
	SerialNumber          string  `json:"serialNumber"`
	LiveIP                string  `json:"liveIp"`
	Port                  int     `json:"port"`
	Username              string  `json:"username"`
	Password              string  `json:"password"`
	LocalIP               string  `json:"localIp"`
	BrandID               *int64  `json:"brandId"`
	ControllerID          *int64  `json:"controllerId"`
	ControllerVersion     string  `json:"controllerVersion"`
	AccuracyID            *int64  `json:"accuracyId"`
	LineID                *int64  `json:"lineId"`
	InstallationDate      string  `json:"installationDate"`
	BillDate              string  `json:"billDate"`
	WarrantyDate          string  `json:"warrantyDate"`
	AMC                   string  `json:"amc"`
	IsInsured             bool    `json:"isInsured"`
	OnlineStatus          bool    `json:"onlineStatus"`
	MaintenanceRequired   bool    `json:"maintenanceRequired"`
	HourlyCost            string  `json:"hourlyCost"`
	Status                string  `json:"status"`
	IsMonitor             bool    `json:"isMonitor"`
	OnlineDate            *string `json:"onlineDate"`
	ModeID                *int64  `json:"modeId"`
	ControllerVersionType *int    `json:"controllerVersionType"`
	HealthIDs             string  `json:"healthIds"`
	Longitude             string  `json:"longitude"`
	Latitude              string  `json:"latitude"`
	CycleAnalysis         bool    `json:"cycleAnalysis"`
	MachineSystemType     string  `json:"machineSystemType"`
	SecurityMode          string  `json:"securityMode"`
}

// Validate trims the fields and checks they are usable.
func (request *CreateMachineRequest) Validate() error {
	request.Name = strings.TrimSpace(request.Name)
	request.SerialNumber = strings.TrimSpace(request.SerialNumber)
	request.LiveIP = strings.TrimSpace(request.LiveIP)
	request.Username = strings.TrimSpace(request.Username)
	request.LocalIP = strings.TrimSpace(request.LocalIP)
	request.ControllerVersion = strings.TrimSpace(request.ControllerVersion)
	request.InstallationDate = strings.TrimSpace(request.InstallationDate)
	request.BillDate = strings.TrimSpace(request.BillDate)
	request.WarrantyDate = strings.TrimSpace(request.WarrantyDate)
	request.AMC = strings.TrimSpace(request.AMC)
	request.HourlyCost = strings.TrimSpace(request.HourlyCost)
	request.Status = strings.ToLower(strings.TrimSpace(request.Status))
	request.HealthIDs = strings.TrimSpace(request.HealthIDs)
	request.Longitude = strings.TrimSpace(request.Longitude)
	request.Latitude = strings.TrimSpace(request.Latitude)
	request.MachineSystemType = strings.TrimSpace(request.MachineSystemType)
	request.SecurityMode = strings.TrimSpace(request.SecurityMode)

	// Validate required fields
	if !namePattern.MatchString(request.Name) {
		return response.ValidName
	}
	if !serialPattern.MatchString(request.SerialNumber) {
		return response.ValidSerialNumber
	}
	if !ipPattern.MatchString(request.LiveIP) {
		return response.ValidIP
	}
	if request.Port <= 0 || request.Port > 65535 {
		return response.ValidPort
	}
	if !machineUsernamePattern.MatchString(request.Username) {
		return response.ValidUserName
	}
	if len(request.Password) == 0 {
		return response.ValidPassword
	}
	if !ipPattern.MatchString(request.LocalIP) {
		return response.ValidIP
	}
	if len(request.ControllerVersion) > 10 {
		return response.ValidControllerVersion
	}
	if !datePattern.MatchString(request.InstallationDate) {
		return response.ValidDate
	}
	if !datePattern.MatchString(request.BillDate) {
		return response.ValidDate
	}
	if !datePattern.MatchString(request.WarrantyDate) {
		return response.ValidDate
	}
	if !decimalPattern.MatchString(request.HourlyCost) {
		return response.ValidPrice
	}
	if request.Status != "" && request.Status != MachineStatusActive && request.Status != MachineStatusInactive {
		return response.ValidStatus
	}
	if request.HealthIDs != "" && len(request.HealthIDs) > 100 {
		return response.ValidHealthIDs
	}
	if request.Longitude != "" && !coordinatePattern.MatchString(request.Longitude) {
		return response.ValidCoordinate
	}
	if request.Latitude != "" && !coordinatePattern.MatchString(request.Latitude) {
		return response.ValidCoordinate
	}
	if request.MachineSystemType != "" && request.MachineSystemType != MachineSystemTypeNewPlc && request.MachineSystemType != MachineSystemTypeOldPlc {
		return response.ValidMachineSystemType
	}
	if request.SecurityMode != "" && request.SecurityMode != SecurityModeNone && request.SecurityMode != SecurityModeBasic256Sha256 {
		return response.ValidSecurityMode
	}

	return nil
}

// UpdateMachineRequest is the body for PUT /machines/{id}.
type UpdateMachineRequest struct {
	UserID                *int64  `json:"userId"`
	MachineTypeID         *int64  `json:"machineTypeId"`
	ProductID             *int64  `json:"productId"`
	Name                  string  `json:"name"`
	SerialNumber          string  `json:"serialNumber"`
	LiveIP                string  `json:"liveIp"`
	Port                  int     `json:"port"`
	Username              string  `json:"username"`
	Password              string  `json:"password"`
	LocalIP               string  `json:"localIp"`
	BrandID               *int64  `json:"brandId"`
	ControllerID          *int64  `json:"controllerId"`
	ControllerVersion     string  `json:"controllerVersion"`
	AccuracyID            *int64  `json:"accuracyId"`
	LineID                *int64  `json:"lineId"`
	InstallationDate      string  `json:"installationDate"`
	BillDate              string  `json:"billDate"`
	WarrantyDate          string  `json:"warrantyDate"`
	AMC                   string  `json:"amc"`
	IsInsured             bool    `json:"isInsured"`
	OnlineStatus          bool    `json:"onlineStatus"`
	MaintenanceRequired   bool    `json:"maintenanceRequired"`
	HourlyCost            string  `json:"hourlyCost"`
	Status                string  `json:"status"`
	IsMonitor             bool    `json:"isMonitor"`
	OnlineDate            *string `json:"onlineDate"`
	ModeID                *int64  `json:"modeId"`
	ControllerVersionType *int    `json:"controllerVersionType"`
	HealthIDs             string  `json:"healthIds"`
	Longitude             string  `json:"longitude"`
	Latitude              string  `json:"latitude"`
	CycleAnalysis         bool    `json:"cycleAnalysis"`
	MachineSystemType     string  `json:"machineSystemType"`
	SecurityMode          string  `json:"securityMode"`
}

// Validate trims the fields and checks they are usable.
func (request *UpdateMachineRequest) Validate() error {
	request.Name = strings.TrimSpace(request.Name)
	request.SerialNumber = strings.TrimSpace(request.SerialNumber)
	request.LiveIP = strings.TrimSpace(request.LiveIP)
	request.Username = strings.TrimSpace(request.Username)
	request.LocalIP = strings.TrimSpace(request.LocalIP)
	request.ControllerVersion = strings.TrimSpace(request.ControllerVersion)
	request.InstallationDate = strings.TrimSpace(request.InstallationDate)
	request.BillDate = strings.TrimSpace(request.BillDate)
	request.WarrantyDate = strings.TrimSpace(request.WarrantyDate)
	request.AMC = strings.TrimSpace(request.AMC)
	request.HourlyCost = strings.TrimSpace(request.HourlyCost)
	request.Status = strings.ToLower(strings.TrimSpace(request.Status))
	request.HealthIDs = strings.TrimSpace(request.HealthIDs)
	request.Longitude = strings.TrimSpace(request.Longitude)
	request.Latitude = strings.TrimSpace(request.Latitude)
	request.MachineSystemType = strings.TrimSpace(request.MachineSystemType)
	request.SecurityMode = strings.TrimSpace(request.SecurityMode)

	// Validate fields if provided
	if request.Name != "" && !namePattern.MatchString(request.Name) {
		return response.ValidName
	}
	if request.SerialNumber != "" && !serialPattern.MatchString(request.SerialNumber) {
		return response.ValidSerialNumber
	}
	if request.LiveIP != "" && !ipPattern.MatchString(request.LiveIP) {
		return response.ValidIP
	}
	if request.Port != 0 && (request.Port <= 0 || request.Port > 65535) {
		return response.ValidPort
	}
	if request.Username != "" && !machineUsernamePattern.MatchString(request.Username) {
		return response.ValidUserName
	}
	// Password is optional in update
	if request.LocalIP != "" && !ipPattern.MatchString(request.LocalIP) {
		return response.ValidIP
	}
	if request.ControllerVersion != "" && len(request.ControllerVersion) > 10 {
		return response.ValidControllerVersion
	}
	if request.InstallationDate != "" && !datePattern.MatchString(request.InstallationDate) {
		return response.ValidDate
	}
	if request.BillDate != "" && !datePattern.MatchString(request.BillDate) {
		return response.ValidDate
	}
	if request.WarrantyDate != "" && !datePattern.MatchString(request.WarrantyDate) {
		return response.ValidDate
	}
	if request.HourlyCost != "" && !decimalPattern.MatchString(request.HourlyCost) {
		return response.ValidPrice
	}
	if request.Status != "" && request.Status != MachineStatusActive && request.Status != MachineStatusInactive {
		return response.ValidStatus
	}
	if request.HealthIDs != "" && len(request.HealthIDs) > 100 {
		return response.ValidHealthIDs
	}
	if request.Longitude != "" && !coordinatePattern.MatchString(request.Longitude) {
		return response.ValidCoordinate
	}
	if request.Latitude != "" && !coordinatePattern.MatchString(request.Latitude) {
		return response.ValidCoordinate
	}
	if request.MachineSystemType != "" && request.MachineSystemType != MachineSystemTypeNewPlc && request.MachineSystemType != MachineSystemTypeOldPlc {
		return response.ValidMachineSystemType
	}
	if request.SecurityMode != "" && request.SecurityMode != SecurityModeNone && request.SecurityMode != SecurityModeBasic256Sha256 {
		return response.ValidSecurityMode
	}

	return nil
}
