package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/royandi/gowatch/backend/internal/model"
	monitorchecker "github.com/royandi/gowatch/backend/internal/monitor"
	"github.com/royandi/gowatch/backend/internal/repository"
)

var (
	ErrMonitorNotFound = errors.New(
		"monitor tidak ditemukan",
	)

	ErrInvalidMonitorURL = errors.New(
		"URL monitor tidak valid",
	)

	ErrInvalidMonitorInterval = errors.New(
		"interval minimal 30 detik",
	)

	ErrInvalidMonitorType = errors.New(
		"monitor_type harus website atau api",
	)
)

type MonitorRepository interface {
	Create(
		ctx context.Context,
		userID int64,
		request model.CreateMonitorRequest,
		isActive bool,
	) (model.Monitor, error)

	FindAllByUser(
		ctx context.Context,
		userID int64,
	) ([]model.Monitor, error)

	FindByIDAndUser(
		ctx context.Context,
		id int64,
		userID int64,
	) (model.Monitor, error)

	Update(
		ctx context.Context,
		monitor model.Monitor,
	) (model.Monitor, error)

	Delete(
		ctx context.Context,
		id int64,
		userID int64,
	) (bool, error)
}

type MonitorService struct {
	monitorRepository MonitorRepository
}

func NewMonitorService(
	monitorRepository *repository.MonitorRepository,
) *MonitorService {
	return &MonitorService{
		monitorRepository: monitorRepository,
	}
}

func newMonitorServiceWithRepository(
	monitorRepository MonitorRepository,
) *MonitorService {
	return &MonitorService{
		monitorRepository: monitorRepository,
	}
}

func (s *MonitorService) Create(
	ctx context.Context,
	userID int64,
	request model.CreateMonitorRequest,
) (model.Monitor, error) {
	request.Name =
		strings.TrimSpace(
			request.Name,
		)

	request.MonitorType =
		strings.ToLower(
			strings.TrimSpace(
				request.MonitorType,
			),
		)

	if request.MonitorType == "" {
		request.MonitorType =
			"website"
	}

	if request.IntervalSeconds == 0 {
		request.IntervalSeconds =
			60
	}

	normalizedURL, err :=
		monitorchecker.NormalizeAndValidateURL(
			ctx,
			request.URL,
		)

	if err != nil {
		return model.Monitor{},
			fmt.Errorf(
				"%w: %v",
				ErrInvalidMonitorURL,
				err,
			)
	}

	request.URL =
		normalizedURL

	if request.IntervalSeconds <
		30 {
		return model.Monitor{},
			ErrInvalidMonitorInterval
	}

	if request.MonitorType !=
		"website" &&
		request.MonitorType !=
			"api" {
		return model.Monitor{},
			ErrInvalidMonitorType
	}

	isActive := true

	if request.IsActive != nil {
		isActive =
			*request.IsActive
	}

	return s.monitorRepository.Create(
		ctx,
		userID,
		request,
		isActive,
	)
}

func (s *MonitorService) List(
	ctx context.Context,
	userID int64,
) ([]model.Monitor, error) {
	return s.monitorRepository.FindAllByUser(
		ctx,
		userID,
	)
}

func (s *MonitorService) GetByID(
	ctx context.Context,
	id int64,
	userID int64,
) (model.Monitor, error) {
	currentMonitor, err :=
		s.monitorRepository.FindByIDAndUser(
			ctx,
			id,
			userID,
		)

	if errors.Is(
		err,
		pgx.ErrNoRows,
	) {
		return model.Monitor{},
			ErrMonitorNotFound
	}

	return currentMonitor,
		err
}

func (s *MonitorService) Update(
	ctx context.Context,
	id int64,
	userID int64,
	request model.UpdateMonitorRequest,
) (model.Monitor, error) {
	currentMonitor, err :=
		s.GetByID(
			ctx,
			id,
			userID,
		)

	if err != nil {
		return model.Monitor{},
			err
	}

	if request.Name != nil {
		currentMonitor.Name =
			strings.TrimSpace(
				*request.Name,
			)
	}

	if request.URL != nil {
		normalizedURL, err :=
			monitorchecker.NormalizeAndValidateURL(
				ctx,
				*request.URL,
			)

		if err != nil {
			return model.Monitor{},
				fmt.Errorf(
					"%w: %v",
					ErrInvalidMonitorURL,
					err,
				)
		}

		currentMonitor.URL =
			normalizedURL
	}

	if request.MonitorType != nil {
		monitorType :=
			strings.ToLower(
				strings.TrimSpace(
					*request.MonitorType,
				),
			)

		if monitorType !=
			"website" &&
			monitorType !=
				"api" {
			return model.Monitor{},
				ErrInvalidMonitorType
		}

		currentMonitor.MonitorType =
			monitorType
	}

	if request.IntervalSeconds != nil {
		if *request.IntervalSeconds <
			30 {
			return model.Monitor{},
				ErrInvalidMonitorInterval
		}

		currentMonitor.IntervalSeconds =
			*request.IntervalSeconds
	}

	if request.IsActive != nil {
		currentMonitor.IsActive =
			*request.IsActive
	}

	if strings.TrimSpace(
		currentMonitor.Name,
	) == "" {
		return model.Monitor{},
			errors.New(
				"nama monitor wajib diisi",
			)
	}

	return s.monitorRepository.Update(
		ctx,
		currentMonitor,
	)
}

func (s *MonitorService) Delete(
	ctx context.Context,
	id int64,
	userID int64,
) error {
	deleted, err :=
		s.monitorRepository.Delete(
			ctx,
			id,
			userID,
		)

	if err != nil {
		return err
	}

	if !deleted {
		return ErrMonitorNotFound
	}

	return nil
}
