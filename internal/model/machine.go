package model

import "time"

// Machine is a row of the machines table.
type Machine struct {
	ID                     int64      `db:"id" json:"id"`
	UserID                 *int64     `db:"user_id" json:"userId"`
	MachineTypeID          *int64     `db:"machine_type_id" json:"machineTypeId"`
	ProductID              *int64     `db:"product_id" json:"productId"`
	Name                   string     `db:"name" json:"name"`
	SerialNumber           string     `db:"serial_number" json:"serialNumber"`
	LiveIP                 string     `db:"live_ip" json:"liveIp"`
	Port                   int        `db:"port" json:"port"`
	Username               string     `db:"username" json:"username"`
	Password               string     `db:"password" json:"password"`
	LocalIP                string     `db:"local_ip" json:"localIp"`
	BrandID                *int64     `db:"brand_id" json:"brandId"`
	ControllerID           *int64     `db:"controller_id" json:"controllerId"`
	ControllerVersion      string     `db:"controller_version" json:"controllerVersion"`
	AccuracyID             *int64     `db:"accuracy_id" json:"accuracyId"`
	LineID                 *int64     `db:"line_id" json:"lineId"`
	InstallationDate       string     `db:"installation_date" json:"installationDate"`
	BillDate               string     `db:"bill_date" json:"billDate"`
	WarrantyDate           string     `db:"warranty_date" json:"warrantyDate"`
	AMC                    string     `db:"amc" json:"amc"`
	IsInsured              bool       `db:"is_insured" json:"isInsured"`
	OnlineStatus           bool       `db:"online_status" json:"onlineStatus"`
	MaintenanceRequired    bool       `db:"maintenance_required" json:"maintenanceRequired"`
	HourlyCost             string     `db:"hourly_cost" json:"hourlyCost"`
	Status                 string     `db:"status" json:"status"`
	IsMonitor              bool       `db:"is_monitor" json:"isMonitor"`
	OnlineDate             *string    `db:"online_date" json:"onlineDate"`
	OfflineDate            *time.Time `db:"offline_date" json:"offlineDate"`
	ModeID                 *int64     `db:"mode_id" json:"modeId"`
	ControllerVersionType  *int       `db:"controller_version_type" json:"controllerVersionType"`
	HealthIDs              string     `db:"health_ids" json:"healthIds"`
	Longitude              string     `db:"longitude" json:"longitude"`
	Latitude               string     `db:"latitude" json:"latitude"`
	CycleAnalysis          bool       `db:"cycle_analysis" json:"cycleAnalysis"`
	MachineSystemType      string     `db:"machine_system_type" json:"machineSystemType"`
	SecurityMode           string     `db:"security_mode" json:"securityMode"`
	IsDeleted              bool       `db:"is_deleted" json:"isDeleted"`
	DeletedAt              *time.Time `db:"deleted_at" json:"deletedAt"`
	CreatedAt              time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt              time.Time  `db:"updated_at" json:"updatedAt"`
}
