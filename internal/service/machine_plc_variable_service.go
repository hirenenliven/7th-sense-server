package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"7th-sense-server/internal/db"
	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/response"
)

// MachinePlcVariableService holds the rules of a machine's PLC node ids: each PLC variable reads
// its default node id until the superadmin gives the machine its own,
// and saving the default (or resetting) drops the machine's own again. The collector picks up
// every change through machine_variables' cnc_config_changed trigger.
type MachinePlcVariableService struct {
	dbPool                *pgxpool.Pool
	plcVariableRepository *repository.MachinePlcVariableRepository
	machineService        *MachineService
}

// NewMachinePlcVariableService creates a MachinePlcVariableService. dbPool runs the reset
// transactions; machineService decides which machines the user may read.
func NewMachinePlcVariableService(dbPool *pgxpool.Pool, plcVariableRepository *repository.MachinePlcVariableRepository, machineService *MachineService) *MachinePlcVariableService {
	return &MachinePlcVariableService{dbPool: dbPool, plcVariableRepository: plcVariableRepository, machineService: machineService}
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
	return dto.NewMachinePlcVariablesResponse(machine, variables), nil
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
	return dto.NewMachinePlcVariableResponse(variable), nil
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
	return dto.NewMachinePlcVariablesResponse(machine, variables), nil
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
	return dto.NewMachinePlcVariableResponse(variable), nil
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
