package service

import (
	"context"
	"errors"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/response"
)

// MachineService holds the business rules for machines.
type MachineService struct {
	machineRepository *repository.MachineRepository
}

// NewMachineService creates a MachineService backed by machineRepository.
func NewMachineService(machineRepository *repository.MachineRepository) *MachineService {
	return &MachineService{machineRepository: machineRepository}
}

func (machineService *MachineService) GetAll(ctx context.Context) ([]model.Machine, error) {
	return machineService.machineRepository.FindAll(ctx)
}

func (machineService *MachineService) GetByID(ctx context.Context, id int64) (model.Machine, error) {
	machine, err := machineService.machineRepository.FindByID(ctx, id)
	return machine, toMachineResponseError(err)
}

func (machineService *MachineService) Create(ctx context.Context, machineRequest dto.CreateMachineRequest) (model.Machine, error) {
	machine, err := machineService.machineRepository.Create(ctx, model.Machine{
		UserID:                machineRequest.UserID,
		MachineTypeID:         machineRequest.MachineTypeID,
		ProductID:             machineRequest.ProductID,
		Name:                  machineRequest.Name,
		SerialNumber:          machineRequest.SerialNumber,
		LiveIP:                machineRequest.LiveIP,
		Port:                  machineRequest.Port,
		Username:              machineRequest.Username,
		Password:              machineRequest.Password,
		LocalIP:               machineRequest.LocalIP,
		BrandID:               machineRequest.BrandID,
		ControllerID:          machineRequest.ControllerID,
		ControllerVersion:     machineRequest.ControllerVersion,
		AccuracyID:            machineRequest.AccuracyID,
		LineID:                machineRequest.LineID,
		InstallationDate:      machineRequest.InstallationDate,
		BillDate:              machineRequest.BillDate,
		WarrantyDate:          machineRequest.WarrantyDate,
		AMC:                   machineRequest.AMC,
		IsInsured:             machineRequest.IsInsured,
		OnlineStatus:          machineRequest.OnlineStatus,
		MaintenanceRequired:   machineRequest.MaintenanceRequired,
		HourlyCost:            machineRequest.HourlyCost,
		Status:                machineRequest.Status,
		IsMonitor:             machineRequest.IsMonitor,
		OnlineDate:            machineRequest.OnlineDate,
		ModeID:                machineRequest.ModeID,
		ControllerVersionType: machineRequest.ControllerVersionType,
		HealthIDs:             machineRequest.HealthIDs,
		Longitude:             machineRequest.Longitude,
		Latitude:              machineRequest.Latitude,
		CycleAnalysis:         machineRequest.CycleAnalysis,
		MachineSystemType:     machineRequest.MachineSystemType,
		SecurityMode:          machineRequest.SecurityMode,
	})
	return machine, toMachineResponseError(err)
}

func (machineService *MachineService) Update(ctx context.Context, id int64, machineRequest dto.UpdateMachineRequest) (model.Machine, error) {
	machine, err := machineService.machineRepository.Update(ctx, model.Machine{
		ID:                    id,
		UserID:                machineRequest.UserID,
		MachineTypeID:         machineRequest.MachineTypeID,
		ProductID:             machineRequest.ProductID,
		Name:                  machineRequest.Name,
		SerialNumber:          machineRequest.SerialNumber,
		LiveIP:                machineRequest.LiveIP,
		Port:                  machineRequest.Port,
		Username:              machineRequest.Username,
		Password:              machineRequest.Password,
		LocalIP:               machineRequest.LocalIP,
		BrandID:               machineRequest.BrandID,
		ControllerID:          machineRequest.ControllerID,
		ControllerVersion:     machineRequest.ControllerVersion,
		AccuracyID:            machineRequest.AccuracyID,
		LineID:                machineRequest.LineID,
		InstallationDate:      machineRequest.InstallationDate,
		BillDate:              machineRequest.BillDate,
		WarrantyDate:          machineRequest.WarrantyDate,
		AMC:                   machineRequest.AMC,
		IsInsured:             machineRequest.IsInsured,
		OnlineStatus:          machineRequest.OnlineStatus,
		MaintenanceRequired:   machineRequest.MaintenanceRequired,
		HourlyCost:            machineRequest.HourlyCost,
		Status:                machineRequest.Status,
		IsMonitor:             machineRequest.IsMonitor,
		OnlineDate:            machineRequest.OnlineDate,
		ModeID:                machineRequest.ModeID,
		ControllerVersionType: machineRequest.ControllerVersionType,
		HealthIDs:             machineRequest.HealthIDs,
		Longitude:             machineRequest.Longitude,
		Latitude:              machineRequest.Latitude,
		CycleAnalysis:         machineRequest.CycleAnalysis,
		MachineSystemType:     machineRequest.MachineSystemType,
		SecurityMode:          machineRequest.SecurityMode,
	})
	return machine, toMachineResponseError(err)
}

func (machineService *MachineService) Delete(ctx context.Context, id int64) error {
	return toMachineResponseError(machineService.machineRepository.Delete(ctx, id))
}

// toMachineResponseError maps repository errors to machine responses. Other errors
// pass through and are sent as ServerError.
func toMachineResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.MachineNotFound
	case errors.Is(err, repository.ErrDuplicateName):
		return response.NameExists
	default:
		return err
	}
}
