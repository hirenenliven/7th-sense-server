package dto

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"7th-sense-server/internal/response"
)

// AxisValuesMaxRawValues is the most stored values GET /machine-health/axis-values returns per
// axis as they are (every value, no buckets): about 12 hours of a busy machine, as many as a chart
// can still draw. A range holding more is cut into AxisValuesMaxPoints time buckets instead.
const AxisValuesMaxRawValues = 20000

// AxisValuesMaxPoints is the number of time buckets a range is cut into when it holds more than
// AxisValuesMaxRawValues values on an axis, so a chart never gets more points than it can draw.
const AxisValuesMaxPoints = 600

// AxisValuesMaxMinutes is the longest "last N minutes" range (MachineCycleMaxDays days).
const AxisValuesMaxMinutes = MachineCycleMaxDays * 24 * 60

// AxisSet is a group of X, Y, Z variables (cnc_variables codes) the Machine Health axis charts
// can show.
type AxisSet struct {
	Key   string
	Label string
	// Codes of the X, Y and Z variables, in that order
	Codes [3]string
}

// AxisSetMachinePosition is the default axis set.
const AxisSetMachinePosition = "machine_position"

// AxisSets are the axis sets the client may ask for; the codes are cnc_variables rows (migration
// 029: snapshot axis positions, ids 208–219; axis loads 35–37).
var AxisSets = []AxisSet{
	{Key: AxisSetMachinePosition, Label: "Machine position (MCS)", Codes: [3]string{"machPosX", "machPosY", "machPosZ"}},
	{Key: "work_position", Label: "Work position (WCS)", Codes: [3]string{"absPosX", "absPosY", "absPosZ"}},
	{Key: "relative_position", Label: "Relative position", Codes: [3]string{"relPosX", "relPosY", "relPosZ"}},
	{Key: "distance_to_go", Label: "Distance to go", Codes: [3]string{"distToGoX", "distToGoY", "distToGoZ"}},
	{Key: "axis_load", Label: "Axis load", Codes: [3]string{"axLoadX", "axLoadY", "axLoadZ"}},
}

// AxisNames label the three variables of a set.
var AxisNames = [3]string{"X", "Y", "Z"}

// FindAxisSet returns the set with key, or false.
func FindAxisSet(key string) (AxisSet, bool) {
	for _, set := range AxisSets {
		if set.Key == key {
			return set, true
		}
	}
	return AxisSet{}, false
}

// AxisValuesRequest is the query of GET /machine-health/axis-values: a machine, an axis set and
// either an exact window Start–End (RFC 3339 instants, the chart's zoom), the last Minutes up to
// now, or whole days From–To (like the mode timeline).
type AxisValuesRequest struct {
	MachineID int64
	AxisSet   string
	Start     string
	End       string
	Minutes   int
	From      string
	To        string

	startAt time.Time
	endAt   time.Time
}

// Validate checks the machine, the axis set (blank = machine position) and the range: at most
// AxisValuesMaxMinutes minutes, or valid days at most MachineCycleMaxDays apart.
func (axisRequest *AxisValuesRequest) Validate() error {
	if axisRequest.AxisSet == "" {
		axisRequest.AxisSet = AxisSetMachinePosition
	}
	if _, ok := FindAxisSet(axisRequest.AxisSet); !ok {
		return response.NewError(http.StatusBadRequest, "Please select valid axis set")
	}
	if axisRequest.Start != "" || axisRequest.End != "" {
		return axisRequest.checkWindow()
	}
	if axisRequest.Minutes < 0 || axisRequest.Minutes > AxisValuesMaxMinutes {
		return response.NewError(http.StatusBadRequest, fmt.Sprintf("Please select at most %d minutes", AxisValuesMaxMinutes))
	}
	listRequest := axisRequest.dayRequest()
	return listRequest.Validate()
}

// checkWindow parses Start and End (both needed): End after Start, at most AxisValuesMaxMinutes
// apart. An End after now is cut to now by Range.
func (axisRequest *AxisValuesRequest) checkWindow() error {
	startAt, startErr := time.Parse(time.RFC3339Nano, axisRequest.Start)
	endAt, endErr := time.Parse(time.RFC3339Nano, axisRequest.End)
	switch {
	case axisRequest.MachineID <= 0:
		return response.SelectMachine
	case startErr != nil || endErr != nil:
		return response.NewError(http.StatusBadRequest, "Please enter valid start and end time")
	case !startAt.Before(endAt):
		return response.NewError(http.StatusBadRequest, "Start time must be before end time")
	case endAt.Sub(startAt) > AxisValuesMaxMinutes*time.Minute:
		return response.NewError(http.StatusBadRequest, fmt.Sprintf("Please select at most %d days", MachineCycleMaxDays))
	}
	axisRequest.startAt, axisRequest.endAt = startAt, endAt
	return nil
}

// Range returns the time span to draw: the Start–End window, else the last Minutes up to now, else
// the From day to the end of the To day; never after now. Call Validate first.
func (axisRequest *AxisValuesRequest) Range(now time.Time) (from, to time.Time) {
	if !axisRequest.startAt.IsZero() {
		return axisRequest.startAt, minTimeOf(axisRequest.endAt, now)
	}
	if axisRequest.Minutes > 0 {
		return now.Add(-time.Duration(axisRequest.Minutes) * time.Minute), now
	}
	listRequest := axisRequest.dayRequest()
	from, to = listRequest.Range(now)
	if to.After(now) {
		to = now
	}
	return from, to
}

func minTimeOf(first, second time.Time) time.Time {
	if first.Before(second) {
		return first
	}
	return second
}

func (axisRequest *AxisValuesRequest) dayRequest() MachineCycleListRequest {
	return MachineCycleListRequest{MachineID: axisRequest.MachineID, From: axisRequest.From, To: axisRequest.To}
}

// AxisBucketSeconds is the bucket length that cuts from–to into at most AxisValuesMaxPoints
// buckets (at least one second).
func AxisBucketSeconds(from, to time.Time) int64 {
	seconds := math.Ceil(to.Sub(from).Seconds() / AxisValuesMaxPoints)
	return int64(math.Max(1, seconds))
}

// AxisValuesResponse is the data of GET /machine-health/axis-values: one series per axis and how
// long the database took to produce them.
type AxisValuesResponse struct {
	MachineID   int64     `json:"machineId"`
	MachineName string    `json:"machineName"`
	AxisSet     string    `json:"axisSet"`
	AxisLabel   string    `json:"axisLabel"`
	From        time.Time `json:"from"`
	To          time.Time `json:"to"`
	// Raw: every stored value (BucketSec 0); else points are time buckets of BucketSec seconds
	Raw       bool         `json:"raw"`
	BucketSec int64        `json:"bucketSec"`
	Axes      []AxisSeries `json:"axes"`
	Stats     AxisStats    `json:"stats"`
}

// AxisSeries is one axis (X, Y or Z) of the set. StartValue is the value at From (the last one
// stored before it; null when the machine never reported one); Points are the stored values (raw)
// or the buckets with rows.
type AxisSeries struct {
	Axis       string      `json:"axis"`
	Code       string      `json:"code"`
	Name       string      `json:"name"`
	Unit       string      `json:"unit"`
	StartValue *float64    `json:"startValue"`
	Points     []AxisPoint `json:"points"`
}

// AxisPoint is one stored value (raw: Min = Max = Value, Count 1) or one time bucket: the last
// value in it (what the chart draws, values are stored only on change), the lowest and highest
// value and the number of stored rows.
type AxisPoint struct {
	At    time.Time `json:"at"`
	Value float64   `json:"value"`
	Min   float64   `json:"min"`
	Max   float64   `json:"max"`
	Count int64     `json:"count"`
}

// AxisStats says how much was read and how long it took: RowsRead stored rows (all axes) became
// PointCount points; QueryMs is the database time of the telemetry queries (the raw read, plus
// the bucket query when the range holds too many values).
type AxisStats struct {
	RowsRead   int64   `json:"rowsRead"`
	PointCount int     `json:"pointCount"`
	QueryMs    float64 `json:"queryMs"`
}
