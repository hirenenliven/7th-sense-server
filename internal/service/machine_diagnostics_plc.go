package service

import (
	"fmt"
	"strings"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/opcua"
)

// PLC section states besides the worker's own (opcua.State*).
const (
	plcStateOff         = "off"          // PLC_ENABLED=false or not a Siemens controller
	plcStateNotReported = "not_reported" // no PLC worker reports the machine
)

// splitPLCVariables separates the variables the PLC worker reads (address in source plc) from
// the CNC DLL's, when the PLC reader reads this machine.
func splitPLCVariables(variables []model.DiagnosticVariable, readsPLC bool) (cncVariables, plcVariables []model.DiagnosticVariable) {
	for _, variable := range variables {
		if readsPLC && variable.Source == model.VariableSourcePLC {
			plcVariables = append(plcVariables, variable)
		} else {
			cncVariables = append(cncVariables, variable)
		}
	}
	return cncVariables, plcVariables
}

// plcSection builds the PLC part of the diagnostics from the worker's status (reported = false
// when no worker reports the machine) and the machine's PLC variables with their live values.
func plcSection(plcEnabled bool, engine *int, machine model.Machine, status opcua.MachineStatus, reported bool,
	variables []model.DiagnosticVariable, values map[string]any) dto.DiagnosticPLC {
	section := dto.DiagnosticPLC{
		Enabled:   plcEnabled,
		Supported: model.ReadsPLCOverOPCUA(engine),
		Endpoint:  opcua.Endpoint(machine.LiveIP, machine.Port),
		Variables: []dto.DiagnosticPLCVariable{},
	}
	if reported {
		section.State, section.Endpoint, section.SecurityMode = status.State, status.Endpoint, status.SecurityMode
		section.SecurityPolicy, section.ConnectedSince = status.SecurityPolicy, status.ConnectedSince
		section.LastError, section.LastErrorAt = status.LastError, status.LastErrorAt
		section.LastChangeAt, section.ChangeCount = status.LastChangeAt, status.ChangeCount
	}

	refused := 0
	for _, variable := range variables {
		node, subscribed := status.Nodes[variable.ID]
		row := dto.DiagnosticPLCVariable{
			ID: variable.ID, Code: variable.Code, Name: variable.Name, Unit: variable.Unit, Category: variable.Category,
			DataType: variable.DataType, NodeID: variable.Address, Enabled: variable.Enabled,
			Subscribed: subscribed && node.Status == "" && status.State == opcua.StateConnected,
			NodeStatus: node.Status, Value: values[variable.Code],
			StoredValue: variable.StoredValue, StoredAt: variable.StoredAt,
		}
		if node.Status != "" {
			refused++
		}
		section.Variables = append(section.Variables, row)
	}

	section.Status, section.Message, section.Hint = plcVerdict(section, reported, len(variables), refused)
	if !reported && section.State == "" {
		section.State = plcStateNotReported
		if !plcEnabled || !section.Supported {
			section.State = plcStateOff
		}
	}
	return section
}

// plcVerdict says, like a DiagnosticCheck, whether the PLC variables are read and what to do.
func plcVerdict(section dto.DiagnosticPLC, reported bool, variableCount, refused int) (status, message, hint string) {
	switch {
	case !section.Enabled:
		return dto.DiagnosticSkipped, "PLC reading is off on this server.", "Set PLC_ENABLED=true and restart the server."
	case !section.Supported:
		return dto.DiagnosticSkipped, "PLC variables are read over OPC UA from Siemens controllers only.",
			"This controller is not read over OPC UA; its PLC variables are not collected."
	case variableCount == 0:
		return dto.DiagnosticSkipped, "This machine has no PLC variable to read.", "Add PLC addresses to the variable catalog (cnc_variable_addresses)."
	case !reported:
		return dto.DiagnosticFailed, "No PLC reader reports this machine.",
			"The machine must be Active, and included in CNC_MACHINE_IDS when that is set."
	case section.State == opcua.StateNoVariables:
		return dto.DiagnosticSkipped, "No active PLC variable has a node ID.", "Check the node IDs on the PLC node IDs page."
	case section.State == opcua.StateConnecting:
		return dto.DiagnosticWarning, "Connecting to " + section.Endpoint + "…", "Wait a few seconds; the page reloads by itself."
	case section.State == opcua.StateFailed:
		return dto.DiagnosticFailed, "Cannot read the PLC: " + section.LastError, plcErrorHint(section.LastError, section.SecurityMode)
	case refused > 0:
		return dto.DiagnosticWarning, fmt.Sprintf("Connected; the controller refused %d of %d node IDs.", refused, variableCount),
			"Check the refused node IDs on the PLC node IDs page (BadNodeIdUnknown: the node does not exist on this controller)."
	case section.LastChangeAt == nil:
		return dto.DiagnosticOK, fmt.Sprintf("Connected and subscribed to %d node IDs; no value yet.", variableCount),
			"Values arrive when the controller publishes them."
	}
	return dto.DiagnosticOK, fmt.Sprintf("Connected and subscribed to %d node IDs.", variableCount), ""
}

// plcErrorHint says what to check for a connect or connection error of the PLC reader.
func plcErrorHint(lastError, securityMode string) string {
	text := strings.ToLower(lastError)
	switch {
	case strings.Contains(text, "baduseraccessdenied") || strings.Contains(text, "badidentitytoken"):
		return "The controller refused the login: check Username and Password on the machine form."
	case strings.Contains(text, "certificate") || strings.Contains(text, "badsecuritychecksfailed") || strings.Contains(text, "badrequestnotallowed"):
		return "The controller does not trust this server: add client_cert.der to the trusted certificates on the NCU (the same certificate the CNC DLL uses)."
	case strings.Contains(text, "no "+strings.ToLower(securityMode)+" endpoint"):
		return "The controller offers no " + securityMode + " endpoint: change Security Mode on the machine form."
	case strings.Contains(text, "refused") || strings.Contains(text, "timeout") || strings.Contains(text, "get endpoints") ||
		strings.Contains(text, "unreachable") || strings.Contains(text, "no such host"):
		return "Check the Live IP and Port on the machine form (Siemens OPC UA is usually 4840) and that the OPC UA server is enabled on the NCU."
	case strings.Contains(text, "connection lost") || strings.Contains(text, "server time") || strings.Contains(text, "publish failed"):
		return "The connection dropped; the server reconnects by itself. If it keeps dropping, check the network to the machine."
	}
	return "The server retries by itself; see logs/app.log (logType plc) for details."
}
