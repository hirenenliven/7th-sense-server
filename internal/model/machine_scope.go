package model

// MachineScope narrows the machines a query returns. A nil field does not narrow,
// so the zero value means every machine.
type MachineScope struct {
	OwnerID        *int64 // only the machines of this customer (machines.user_id)
	AssignedUserID *int64 // only the machines this user may access (user_machine_access or user_machine_filters)
}
