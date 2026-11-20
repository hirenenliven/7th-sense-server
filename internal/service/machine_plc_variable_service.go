package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"7th-sense-server/internal/cert"
	"7th-sense-server/internal/config"
	"7th-sense-server/internal/db"
	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/opcua"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/response"
)

// MachinePlcVariableService holds the rules of a machine's PLC node ids: each PLC variable reads
// its default node id until the superadmin gives the machine its own,
// and saving the default (or resetting) drops the machine's own again. The collector picks up
// every change through machine_variables' cnc_config_changed trigger. Each row also says what the
// machine's PLC reader does with it, and node ids can be read once to test them
// (machine_plc_node_probe.go).
type MachinePlcVariableService struct {
	dbPool                *pgxpool.Pool
	plcVariableRepository *repository.MachinePlcVariableRepository
	machineService        *MachineService
	diagnosticsRepository *repository.MachineDiagnosticsRepository // the machine's CNC engine
	plcConfig             config.PLCConfig
	plcBoard              *opcua.StatusBoard // the PLC workers' status; nil when PLC_ENABLED is off
	certificate           cert.Paths         // presented to controllers with security Basic256Sha256
}

// NewMachinePlcVariableService creates a MachinePlcVariableService. dbPool runs the reset
// transactions; machineService decides which machines the user may read; plcConfig, plcBoard and
// certificate are the PLC reader's settings, status and client certificate.
func NewMachinePlcVariableService(dbPool *pgxpool.Pool, plcVariableRepository *repository.MachinePlcVariableRepository, machineService *MachineService,
	diagnosticsRepository *repository.MachineDiagnosticsRepository, plcConfig config.PLCConfig, plcBoard *opcua.StatusBoard, certificate cert.Paths) *MachinePlcVariableService {
	return &MachinePlcVariableService{dbPool: dbPool, plcVariableRepository: plcVariableRepository, machineService: machineService,
		diagnosticsRepository: diagnosticsRepository, plcConfig: plcConfig, plcBoard: plcBoard, certificate: certificate}
}

// GetAll returns the machine's PLC variables with the node id each reads.
func (plcVariableService *MachinePlcVariableService) GetAll(ctx context.Context, actor model.AuthUser, machineID int64) (dto.MachinePlcVariablesResponse, error) {
	machine, err := plcVariableService.machineService.GetByID(ctx, actor, nil, machineID)
	if err != nil {
		return dto.MachinePlcVariablesResponse{}, err
	}
	variables, err := plcVariableService.plcVariableRepository.FindAll(ctx, machineID)
	if err != nil {
		return dto.MachinePlcVariablesResponse{}, err
	}
	return plcVariableService.withReader(ctx, dto.NewMachinePlcVariablesResponse(machine, variables))
}

// Update sets the machine's node id of one PLC variable. The variable's default node id is not
// stored as the machine's own: it resets the variable, so a later change of the default reaches it.
func (plcVariableService *MachinePlcVariableService) Update(ctx context.Context, actor model.AuthUser, machineID int64, variableID int, request dto.UpdatePlcNodeIDRequest) (dto.MachinePlcVariableResponse, error) {
	variable, err := plcVariableService.findVariable(ctx, actor, machineID, variableID)
	if err != nil {
		return dto.MachinePlcVariableResponse{}, err
	}
	if request.NodeID == variable.DefaultNodeID {
		return plcVariableService.reset(ctx, machineID, variable)
	}
	if err := plcVariableService.plcVariableRepository.SaveNodeID(ctx, machineID, variableID, request.NodeID); err != nil {
		return dto.MachinePlcVariableResponse{}, err
	}
	variable.CustomNodeID = request.NodeID
	return plcVariableService.withReaderNode(machineID, dto.NewMachinePlcVariableResponse(variable)), nil
}

// Reset makes one PLC variable of the machine read its default node id again.
func (plcVariableService *MachinePlcVariableService) Reset(ctx context.Context, actor model.AuthUser, machineID int64, variableID int) (dto.MachinePlcVariableResponse, error) {
	variable, err := plcVariableService.findVariable(ctx, actor, machineID, variableID)
	if err != nil {
		return dto.MachinePlcVariableResponse{}, err
	}
	return plcVariableService.reset(ctx, machineID, variable)
}

// ResetAll makes every PLC variable of the machine read its default node id again and returns the
// list as GetAll does.
func (plcVariableService *MachinePlcVariableService) ResetAll(ctx context.Context, actor model.AuthUser, machineID int64) (dto.MachinePlcVariablesResponse, error) {
	machine, err := plcVariableService.machineService.GetByID(ctx, actor, nil, machineID)
	if err != nil {
		return dto.MachinePlcVariablesResponse{}, err
	}
	variables, err := plcVariableService.plcVariableRepository.FindAll(ctx, machineID)
	if err != nil {
		return dto.MachinePlcVariablesResponse{}, err
	}

	var updatedIDs []int
	for _, variable := range variables {
		if variable.CustomNodeID != "" {
			updatedIDs = append(updatedIDs, variable.ID)
		}
	}
	if len(updatedIDs) > 0 {
		if err := plcVariableService.deleteNodeIDs(ctx, machineID, updatedIDs); err != nil {
			return dto.MachinePlcVariablesResponse{}, err
		}
	}
	for i := range variables {
		variables[i].CustomNodeID = ""
	}
	return plcVariableService.withReader(ctx, dto.NewMachinePlcVariablesResponse(machine, variables))
}

// findVariable returns the PLC variable of a machine the user may read, else
// response.MachineNotFound or response.PlcVariableNotFound.
func (plcVariableService *MachinePlcVariableService) findVariable(ctx context.Context, actor model.AuthUser, machineID int64, variableID int) (model.MachinePlcVariable, error) {
	if _, err := plcVariableService.machineService.GetByID(ctx, actor, nil, machineID); err != nil {
		return model.MachinePlcVariable{}, err
	}
	variables, err := plcVariableService.plcVariableRepository.FindAll(ctx, machineID)
	if err != nil {
		return model.MachinePlcVariable{}, err
	}
	for _, variable := range variables {
		if variable.ID == variableID {
			return variable, nil
		}
	}
	return model.MachinePlcVariable{}, toMachinePlcVariableResponseError(repository.ErrNotFound)
}

func (plcVariableService *MachinePlcVariableService) reset(ctx context.Context, machineID int64, variable model.MachinePlcVariable) (dto.MachinePlcVariableResponse, error) {
	if variable.CustomNodeID != "" {
		if err := plcVariableService.deleteNodeIDs(ctx, machineID, []int{variable.ID}); err != nil {
			return dto.MachinePlcVariableResponse{}, err
		}
	}
	variable.CustomNodeID = ""
	return plcVariableService.withReaderNode(machineID, dto.NewMachinePlcVariableResponse(variable)), nil
}

// withReader adds what the machine's PLC reader does to list: its state and, per variable, the
// node id it subscribed to and the controller's answer.
func (plcVariableService *MachinePlcVariableService) withReader(ctx context.Context, list dto.MachinePlcVariablesResponse) (dto.MachinePlcVariablesResponse, error) {
	engine, err := plcVariableService.diagnosticsRepository.FindEngine(ctx, list.MachineID)
	if err != nil {
		return dto.MachinePlcVariablesResponse{}, err
	}
	status, reported := plcVariableService.plcBoard.Get(list.MachineID)
	list.Reader = dto.PlcReaderResponse{Enabled: plcVariableService.plcConfig.Enabled, Supported: model.ReadsPLCOverOPCUA(engine)}
	if reported {
		list.Reader.State, list.Reader.LastError = status.State, status.LastError
	}
	for i, row := range list.Variables {
		node := status.Nodes[row.VariableID]
		list.Variables[i].ReaderNodeID, list.Variables[i].ReaderStatus = node.NodeID, node.Status
	}
	return list, nil
}

// withReaderNode adds the node id the PLC reader subscribed to for row and the controller's answer.
func (plcVariableService *MachinePlcVariableService) withReaderNode(machineID int64, row dto.MachinePlcVariableResponse) dto.MachinePlcVariableResponse {
	status, _ := plcVariableService.plcBoard.Get(machineID)
	node := status.Nodes[row.VariableID]
	row.ReaderNodeID, row.ReaderStatus = node.NodeID, node.Status
	return row
}

func (plcVariableService *MachinePlcVariableService) deleteNodeIDs(ctx context.Context, machineID int64, variableIDs []int) error {
	return db.WithTransaction(ctx, plcVariableService.dbPool, func(tx pgx.Tx) error {
		return repository.NewMachinePlcVariableRepository(tx).DeleteNodeIDs(ctx, machineID, variableIDs)
	})
}

// toMachinePlcVariableResponseError maps repository errors to response errors.
func toMachinePlcVariableResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.PlcVariableNotFound
	default:
		return err
	}
}
