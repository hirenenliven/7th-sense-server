package model

// VariableSourcePLC is the cnc_variable_addresses / machine_variables source of a PLC address,
// an OPC UA node id such as "ns=2;s=/Random/Plc/DB3300.DBX3.0".
const VariableSourcePLC = "plc"

// MachinePlcVariable is one PLC variable (a cnc_variables row with a plc address for any engine)
// of a machine: its default node id (the best ranked plc address in cnc_variable_addresses: the
// machine's engine and controller version first) and the machine's own node id
// (machine_variables, source plc), empty while the default is used.
type MachinePlcVariable struct {
	ID            int    `db:"id"`
	Code          string `db:"code"`
	Name          string `db:"name"`
	Unit          string `db:"unit"`
	Category      string `db:"category"`
	DataType      string `db:"data_type"`
	Status        string `db:"status"`
	DefaultNodeID string `db:"default_node_id"`
	CustomNodeID  string `db:"custom_node_id"`
}
