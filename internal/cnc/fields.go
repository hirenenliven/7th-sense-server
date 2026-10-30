package cnc

import "fmt"

// AxisCount is the number of axes the DLL snapshot holds positions for (its [5] arrays).
const AxisCount = 5

// The CncMasterSnapshot fields a variable can be read from (cnc_variable_addresses.address
// with source 'snapshot'), with the DLL's spelling. Reading.Values is indexed by these.
const (
	FieldMachineMode = iota
	FieldProcStatus
	FieldIsEmergency
	FieldProgName
	FieldActiveProgPath
	FieldCurProgLine
	FieldTotalProgLines
	FieldProgramPartCount
	FieldRequiredPartCount
	FieldSpindleSpeed
	FieldSpindleLoad
	FieldSpindleOverride
	FieldFeedRate
	FieldFeedOverride
	FieldAxisLoadX
	FieldAxisLoadY
	FieldAxisLoadZ
	FieldToolNumber
	FieldAlarmNo
	FieldAlarmMsg
	FieldTimerMinute // NC timer, whole minutes (which timer the DLL fills is not documented)
	FieldTimerMsec   // NC timer, milliseconds past TimerMinute

	// Axis positions take AxisCount fields each, axis 1 first: "MachPos[0]" … "MachPos[4]".
	FieldMachPos
	FieldAbsPos   = FieldMachPos + AxisCount
	FieldRelPos   = FieldAbsPos + AxisCount
	FieldDistToGo = FieldRelPos + AxisCount

	FieldCount = FieldDistToGo + AxisCount
)

var fieldNames = buildFieldNames()

// buildFieldNames names every field; an axis field is its array name and index, e.g. "AbsPos[2]".
func buildFieldNames() [FieldCount]string {
	names := [FieldCount]string{
		FieldMachineMode:       "MachineMode",
		FieldProcStatus:        "ProcStatus",
		FieldIsEmergency:       "IsEmergency",
		FieldProgName:          "ProgName",
		FieldActiveProgPath:    "ActiveProgPath",
		FieldCurProgLine:       "CurProgLine",
		FieldTotalProgLines:    "TotalProgLines",
		FieldProgramPartCount:  "ProgramPartCount",
		FieldRequiredPartCount: "RequiredPartCount",
		FieldSpindleSpeed:      "SpindeOverrideActual",
		FieldSpindleLoad:       "SpindleLoadPercetnage",
		FieldSpindleOverride:   "SpindeOverridePercentage",
		FieldFeedRate:          "FeedRateActual",
		FieldFeedOverride:      "FeedRatePercetnage",
		FieldAxisLoadX:         "AxisLoadXPercetnage",
		FieldAxisLoadY:         "AxisLoadYPercetnage",
		FieldAxisLoadZ:         "AxisLoadZPercetnage",
		FieldToolNumber:        "ToolNumber",
		FieldAlarmNo:           "AlarmNo",
		FieldAlarmMsg:          "AlarmMsg",
		FieldTimerMinute:       "TimerMinute",
		FieldTimerMsec:         "TimerMsec",
	}
	for axis := range AxisCount {
		names[FieldMachPos+axis] = fmt.Sprintf("MachPos[%d]", axis)
		names[FieldAbsPos+axis] = fmt.Sprintf("AbsPos[%d]", axis)
		names[FieldRelPos+axis] = fmt.Sprintf("RelPos[%d]", axis)
		names[FieldDistToGo+axis] = fmt.Sprintf("DistToGo[%d]", axis)
	}
	return names
}

// FieldIndex returns the index in Reading.Values of a snapshot field name, or false when the
// snapshot has no such field.
func FieldIndex(name string) (int, bool) {
	for index, fieldName := range fieldNames {
		if fieldName == name {
			return index, true
		}
	}
	return 0, false
}
