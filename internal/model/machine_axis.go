package model

import "time"

// AxisVariable is a cnc_variables row an axis chart draws (e.g. machPosX, "X axis machine
// position (MCS)").
type AxisVariable struct {
	ID   int64  `db:"id"`
	Code string `db:"code"`
	Name string `db:"name"`
	Unit string `db:"unit"`
}

// AxisBucket is one time bucket of one variable's machine_telementry rows: the last value in the
// bucket, its lowest and highest value and how many rows it summarises.
type AxisBucket struct {
	VariableID int64     `db:"variable_id"`
	BucketAt   time.Time `db:"bucket_at"`
	LastValue  float64   `db:"last_value"`
	MinValue   float64   `db:"min_value"`
	MaxValue   float64   `db:"max_value"`
	RowCount   int64     `db:"row_count"`
}

// AxisRow is one stored machine_telementry value of an axis variable.
type AxisRow struct {
	VariableID int64     `db:"variable_id"`
	RecordedAt time.Time `db:"recorded_at"`
	Value      float64   `db:"value"`
}

// AxisStartValue is a variable's last machine_telementry value before a chart's range: values are
// stored only when they change, so this is the value at the range start.
type AxisStartValue struct {
	VariableID int64     `db:"variable_id"`
	Value      float64   `db:"value"`
	RecordedAt time.Time `db:"recorded_at"`
}
