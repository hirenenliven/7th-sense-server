package cnc

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

	FieldCount
)

var fieldNames = [FieldCount]string{
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
