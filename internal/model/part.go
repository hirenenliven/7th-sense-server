package model

import "time"

// Part is a row of the parts table, plus the display names of its customer and
// clientele and the number of its processes (read with joins, never written back).
// Processes is filled only when one part is read (GET /parts/{id}). Times are whole seconds.
type Part struct {
	ID             int64      `db:"id" json:"id"`
	UserID         int64      `db:"user_id" json:"userId"`
	ClienteleID    *int64     `db:"clientele_id" json:"clienteleId"`
	PartNo         string     `db:"part_no" json:"partNo"`
	PartName       string     `db:"part_name" json:"partName"`
	PartCount      int        `db:"part_count" json:"partCount"`
	ProductionTime int        `db:"production_time" json:"productionTime"`
	Status         string     `db:"status" json:"status"`
	IsDeleted      bool       `db:"is_deleted" json:"isDeleted"`
	DeletedAt      *time.Time `db:"deleted_at" json:"deletedAt"`
	CreatedAt      time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updatedAt"`

	// Read-only (joins)
	CustomerName  string `db:"customer_name" json:"customerName"`
	ProcessCount  int    `db:"process_count" json:"processCount"`
	ClienteleName string `db:"clientele_name" json:"clienteleName"`

	Processes []PartProcess `db:"-" json:"processes,omitempty"`
}

// PartProcess is a row of the part_processes table. MachineIDs are its applicable
// machines (part_process_machines); SubPrograms is filled by the service.
type PartProcess struct {
	ID                int64      `db:"id" json:"id"`
	PartID            int64      `db:"part_id" json:"partId"`
	SortOrder         int        `db:"sort_order" json:"sortOrder"`
	Process           string     `db:"process" json:"process"`
	ProgramName       string     `db:"program_name" json:"programName"`
	DrawingNo         string     `db:"drawing_no" json:"drawingNo"`
	DrawingRevisionNo string     `db:"drawing_revision_no" json:"drawingRevisionNo"`
	DrawingFile       string     `db:"drawing_file" json:"drawingFile"`
	DrawingName       string     `db:"drawing_name" json:"drawingName"`
	AccuracyID        int64      `db:"accuracy_id" json:"accuracyId"`
	SetupTime         int        `db:"setup_time" json:"setupTime"`
	CycleTime         int        `db:"cycle_time" json:"cycleTime"`
	LoadingTime       int        `db:"loading_time" json:"loadingTime"`
	ProductionTime    int        `db:"production_time" json:"productionTime"`
	HasSubProgram     bool       `db:"has_sub_program" json:"hasSubProgram"`
	MachineIDs        []int64    `db:"machine_ids" json:"machineIds"`
	IsDeleted         bool       `db:"is_deleted" json:"isDeleted"`
	DeletedAt         *time.Time `db:"deleted_at" json:"deletedAt"`
	CreatedAt         time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt         time.Time  `db:"updated_at" json:"updatedAt"`

	SubPrograms []PartSubProgram `db:"-" json:"subPrograms"`
}

// PartSubProgram is a row of the part_sub_programs table.
type PartSubProgram struct {
	ID        int64      `db:"id" json:"id"`
	ProcessID int64      `db:"process_id" json:"processId"`
	SortOrder int        `db:"sort_order" json:"sortOrder"`
	Name      string     `db:"name" json:"name"`
	CycleTime int        `db:"cycle_time" json:"cycleTime"`
	IsDeleted bool       `db:"is_deleted" json:"isDeleted"`
	DeletedAt *time.Time `db:"deleted_at" json:"deletedAt"`
	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time  `db:"updated_at" json:"updatedAt"`
}

// PartDrawing is an uploaded part drawing PDF: File is the stored name a process
// refers to (part_processes.drawing_file), Name the name it was uploaded with.
type PartDrawing struct {
	File string `json:"file"`
	Name string `json:"name"`
}
