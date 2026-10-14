package model

import "time"

// Machine is a row of the machines table, plus the display names of its customer,
// model, brand, controller and line read with LEFT JOINs. Those names are
// never written back.
//
// Password is what the data collector logs in to the machine with. It is never sent
// to clients (json:"-"); HasPassword tells the edit form one is saved.
type Machine struct {
	ID                    int64      `db:"id" json:"id"`
	UserID                *int64     `db:"user_id" json:"userId"`
	MachineTypeID         *int64     `db:"machine_type_id" json:"machineTypeId"`
	ModelID               *int64     `db:"model_id" json:"modelId"`
	Name                  string     `db:"name" json:"name"`
	SerialNumber          string     `db:"serial_number" json:"serialNumber"`
	LiveIP                string     `db:"live_ip" json:"liveIp"`
	Port                  int        `db:"port" json:"port"`
	Username              string     `db:"username" json:"username"`
	Password              string     `db:"password" json:"-"`
	HasPassword           bool       `db:"has_password" json:"hasPassword"`
	LocalIP               string     `db:"local_ip" json:"localIp"`
	BrandID               *int64     `db:"brand_id" json:"brandId"`
	ControllerID          *int64     `db:"controller_id" json:"controllerId"`
	ControllerVersion     string     `db:"controller_version" json:"controllerVersion"`
	AccuracyID            *int64     `db:"accuracy_id" json:"accuracyId"`
	LineID                *int64     `db:"line_id" json:"lineId"`
	InstallationDate      string     `db:"installation_date" json:"installationDate"`
	BillDate              *string    `db:"bill_date" json:"billDate"`
	WarrantyDate          *string    `db:"warranty_date" json:"warrantyDate"`
	AMC                   string     `db:"amc" json:"amc"`
	IsInsured             bool       `db:"is_insured" json:"isInsured"`
	OnlineStatus          bool       `db:"online_status" json:"onlineStatus"`
	MaintenanceRequired   bool       `db:"maintenance_required" json:"maintenanceRequired"`
	HourlyCost            string     `db:"hourly_cost" json:"hourlyCost"`
	Status                string     `db:"status" json:"status"`
	IsMonitor             bool       `db:"is_monitor" json:"isMonitor"`
	OnlineDate            *string    `db:"online_date" json:"onlineDate"`
	OfflineDate           *time.Time `db:"offline_date" json:"offlineDate"`
	ModeID                *int64     `db:"mode_id" json:"modeId"`
	ControllerVersionType *int       `db:"controller_version_type" json:"controllerVersionType"`
	HealthIDs             string     `db:"health_ids" json:"healthIds"`
	Longitude             string     `db:"longitude" json:"longitude"`
	Latitude              string     `db:"latitude" json:"latitude"`
	CycleAnalysis         bool       `db:"cycle_analysis" json:"cycleAnalysis"`
	MachineSystemType     string     `db:"machine_system_type" json:"machineSystemType"`
	SecurityMode          string     `db:"security_mode" json:"securityMode"`
	IsDeleted             bool       `db:"is_deleted" json:"isDeleted"`
	DeletedAt             *time.Time `db:"deleted_at" json:"deletedAt"`
	CreatedAt             time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt             time.Time  `db:"updated_at" json:"updatedAt"`

	// Display names (LEFT JOIN, read-only)
	CustomerName   string `db:"customer_name" json:"customerName"`
	ModelName      string `db:"model_name" json:"modelName"`
	BrandName      string `db:"brand_name" json:"brandName"`
	ControllerName string `db:"controller_name" json:"controllerName"`
	LineName       string `db:"line_name" json:"lineName"`
}
