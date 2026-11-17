package model

import "time"

// ProductionSchedule is a row of the production_schedules table, plus the display names
// of its customer and machine and the number of its items (read with joins, never
// written back). Dates are "YYYY-MM-DD". Items is filled only when one schedule is read.
type ProductionSchedule struct {
	ID           int64      `db:"id" json:"id"`
	UserID       int64      `db:"user_id" json:"userId"`
	MachineID    int64      `db:"machine_id" json:"machineId"`
	ScheduleType string     `db:"schedule_type" json:"scheduleType"`
	StartDate    string     `db:"start_date" json:"startDate"`
	EndDate      string     `db:"end_date" json:"endDate"`
	Status       string     `db:"status" json:"status"`
	IsDeleted    bool       `db:"is_deleted" json:"isDeleted"`
	DeletedAt    *time.Time `db:"deleted_at" json:"deletedAt"`
	CreatedAt    time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updatedAt"`

	// Read-only (joins)
	CustomerName string `db:"customer_name" json:"customerName"`
	CompanyName  string `db:"company_name" json:"companyName"`
	MachineName  string `db:"machine_name" json:"machineName"`
	ItemCount    int    `db:"item_count" json:"itemCount"`
	TotalQty     int    `db:"total_quantity" json:"totalQuantity"`

	Items []ProductionScheduleItem `db:"-" json:"items,omitempty"`
}

// ProductionScheduleItem is a row of the production_schedule_items table, plus the
// display names of what it refers to. Times are "HH:MM".
type ProductionScheduleItem struct {
	ID           int64      `db:"id" json:"id"`
	ScheduleID   int64      `db:"schedule_id" json:"scheduleId"`
	SortOrder    int        `db:"sort_order" json:"sortOrder"`
	ScheduleDate string     `db:"schedule_date" json:"scheduleDate"`
	JobID        string     `db:"job_id" json:"jobId"`
	SupervisorID int64      `db:"supervisor_id" json:"supervisorId"`
	OperatorID   int64      `db:"operator_id" json:"operatorId"`
	ShiftID      int64      `db:"shift_id" json:"shiftId"`
	PartID       int64      `db:"part_id" json:"partId"`
	ProcessID    int64      `db:"process_id" json:"processId"`
	Quantity     int        `db:"quantity" json:"quantity"`
	StartTime    string     `db:"start_time" json:"startTime"`
	EndTime      string     `db:"end_time" json:"endTime"`
	TotalMinutes int        `db:"total_minutes" json:"totalMinutes"`
	IsDeleted    bool       `db:"is_deleted" json:"isDeleted"`
	DeletedAt    *time.Time `db:"deleted_at" json:"deletedAt"`
	CreatedAt    time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updatedAt"`

	// Read-only (joins)
	SupervisorName string `db:"supervisor_name" json:"supervisorName"`
	OperatorName   string `db:"operator_name" json:"operatorName"`
	ShiftName      string `db:"shift_name" json:"shiftName"`
	PartName       string `db:"part_name" json:"partName"`
	ProcessName    string `db:"process_name" json:"processName"`
}

// MachineQuantity is the quantity planned for one machine (production schedule jobs of a day).
type MachineQuantity struct {
	MachineID int64 `db:"machine_id"`
	Quantity  int   `db:"quantity"`
}

// OperatorJob is one production schedule item seen from its operator: the job with the names of
// its operator, shift, machine (and model), part, process and program, and its schedule's customer.
// Read only (Operator Wise Production Schedule). Dates are "YYYY-MM-DD", times "HH:MM".
type OperatorJob struct {
	ItemID       int64  `db:"item_id" json:"itemId"`
	ScheduleID   int64  `db:"schedule_id" json:"scheduleId"`
	UserID       int64  `db:"user_id" json:"userId"`
	CustomerName string `db:"customer_name" json:"customerName"`
	ScheduleDate string `db:"schedule_date" json:"scheduleDate"`
	OperatorID   int64  `db:"operator_id" json:"operatorId"`
	OperatorName string `db:"operator_name" json:"operatorName"`
	ShiftID      int64  `db:"shift_id" json:"shiftId"`
	ShiftName    string `db:"shift_name" json:"shiftName"`
	JobID        string `db:"job_id" json:"jobId"`
	Quantity     int    `db:"quantity" json:"quantity"`
	MachineID    int64  `db:"machine_id" json:"machineId"`
	MachineName  string `db:"machine_name" json:"machineName"`
	ModelName    string `db:"model_name" json:"modelName"`
	StartTime    string `db:"start_time" json:"startTime"`
	EndTime      string `db:"end_time" json:"endTime"`
	PartID       int64  `db:"part_id" json:"partId"`
	PartNo       string `db:"part_no" json:"partNo"`
	PartName     string `db:"part_name" json:"partName"`
	ProcessName  string `db:"process_name" json:"processName"`
	ProgramName  string `db:"program_name" json:"programName"`
}
