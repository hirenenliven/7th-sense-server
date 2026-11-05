package dto

import (
	"time"

	"7th-sense-server/internal/model"
)

// The cnc_variables codes the Virtual Screen reads from machine_variable_values.
const (
	VariableSpindleLoad     = "SpindleLoad"
	VariableSpindleSpeed    = "SpindleRpm"
	VariableSpindleOverride = "speedOvr"
	VariableFeedOverride    = "feedRateOvr"
	VariableFeedRate        = "actFeedRate"
	VariableFeedUnit        = "FeedRev"
	VariableAxisLoadX       = "axLoadX"
	VariableAxisLoadY       = "axLoadY"
	VariableAxisLoadZ       = "axLoadZ"
	VariableToolNumber      = "toolNumber"
	VariableEmergency       = "isEmergency"
)

// VirtualScreenVariableCodes are the variables GET /productivity/virtual-screen reads.
var VirtualScreenVariableCodes = []string{
	VariableSpindleLoad, VariableSpindleSpeed, VariableSpindleOverride, VariableFeedOverride,
	VariableFeedRate, VariableFeedUnit, VariableAxisLoadX, VariableAxisLoadY, VariableAxisLoadZ,
	VariableToolNumber, VariableEmergency,
}

// feedUnits names the FeedRev values (0 or unknown = mm/min).
var feedUnits = map[int]string{0: "mm/min", 1: "mm/rev", 2: "inch/min", 3: "inch/rev"}

// VirtualScreenResponse is GET /productivity/virtual-screen: what a machine is doing now. Every
// live value is null when the collector never read it.
type VirtualScreenResponse struct {
	MachineID   int64               `json:"machineId"`
	ServerTime  time.Time           `json:"serverTime"`
	Online      bool                `json:"online"`
	ModeID      *int                `json:"modeId"`
	ModeName    string              `json:"modeName"`
	ModeSince   *time.Time          `json:"modeSince"`
	ProgramName string              `json:"programName"`
	OpMode      string              `json:"opMode"`
	Live        VirtualScreenLive   `json:"live"`
	Alarm       *model.MachineAlarm `json:"alarm"`
	Cycle       *VirtualScreenCycle `json:"cycle"`
	Job         *model.JobProgress  `json:"job"`
}

// VirtualScreenLive is the machine's latest CNC values: loads and overrides in %, spindle speed
// in rpm, feed rate in FeedUnit. UpdatedAt is when the newest of them was read.
type VirtualScreenLive struct {
	SpindleLoad     *float64   `json:"spindleLoad"`
	SpindleSpeed    *float64   `json:"spindleSpeed"`
	SpindleOverride *float64   `json:"spindleOverride"`
	FeedOverride    *float64   `json:"feedOverride"`
	FeedRate        *float64   `json:"feedRate"`
	FeedUnit        string     `json:"feedUnit"`
	AxisLoadX       *float64   `json:"axisLoadX"`
	AxisLoadY       *float64   `json:"axisLoadY"`
	AxisLoadZ       *float64   `json:"axisLoadZ"`
	ToolNumber      *float64   `json:"toolNumber"`
	Emergency       bool       `json:"emergency"`
	UpdatedAt       *time.Time `json:"updatedAt"`
}

// NewVirtualScreenLive maps the variables read (by code) to the live values.
func NewVirtualScreenLive(variables []model.LiveVariable) VirtualScreenLive {
	values := map[string]*float64{}
	live := VirtualScreenLive{FeedUnit: feedUnits[0]}
	for _, variable := range variables {
		values[variable.Code] = variable.Value
		if live.UpdatedAt == nil || variable.UpdatedAt.After(*live.UpdatedAt) {
			updatedAt := variable.UpdatedAt
			live.UpdatedAt = &updatedAt
		}
	}
	live.SpindleLoad, live.SpindleSpeed = values[VariableSpindleLoad], values[VariableSpindleSpeed]
	live.SpindleOverride, live.FeedOverride = values[VariableSpindleOverride], values[VariableFeedOverride]
	live.FeedRate = values[VariableFeedRate]
	live.AxisLoadX, live.AxisLoadY, live.AxisLoadZ = values[VariableAxisLoadX], values[VariableAxisLoadY], values[VariableAxisLoadZ]
	live.ToolNumber = values[VariableToolNumber]
	if unit := values[VariableFeedUnit]; unit != nil {
		if name, found := feedUnits[int(*unit)]; found {
			live.FeedUnit = name
		}
	}
	if emergency := values[VariableEmergency]; emergency != nil {
		live.Emergency = *emergency != 0
	}
	return live
}

// VirtualScreenCycle is the machine's running cycle: how long it has run and its planned time
// (null when no part process matched it).
type VirtualScreenCycle struct {
	ProgramName      string    `json:"programName"`
	StartedAt        time.Time `json:"startedAt"`
	ElapsedSec       int       `json:"elapsedSec"`
	StandardCycleSec *int      `json:"standardCycleSec"`
}
