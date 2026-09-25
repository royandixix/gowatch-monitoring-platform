package service

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/royandi/gowatch/backend/internal/model"
)

type fakeMonitorRepository struct {
	createMonitor model.Monitor
	createError   error

	findMonitor model.Monitor
	findError   error

	listMonitors []model.Monitor
	listError    error

	updateMonitor model.Monitor
	updateError   error

	deleteResult bool
	deleteError  error

	lastCreateRequest model.CreateMonitorRequest
	lastCreateActive  bool
}

func (
	f *fakeMonitorRepository,
) Create(
	ctx context.Context,
	userID int64,
	request model.CreateMonitorRequest,
	isActive bool,
) (model.Monitor, error) {
	f.lastCreateRequest =
		request

	f.lastCreateActive =
		isActive

	if f.createError != nil {
		return model.Monitor{},
			f.createError
	}

	result :=
		f.createMonitor

	if result.ID == 0 {
		result =
			model.Monitor{
				ID: 1,

				UserID: userID,

				Name: request.Name,

				URL: request.URL,

				MonitorType: request.MonitorType,

				IntervalSeconds: request.IntervalSeconds,

				IsActive: isActive,
			}
	}

	return result, nil
}

func (
	f *fakeMonitorRepository,
) FindAllByUser(
	ctx context.Context,
	userID int64,
) ([]model.Monitor, error) {
	return f.listMonitors,
		f.listError
}

func (
	f *fakeMonitorRepository,
) FindByIDAndUser(
	ctx context.Context,
	id int64,
	userID int64,
) (model.Monitor, error) {
	if f.findError != nil {
		return model.Monitor{},
			f.findError
	}

	return f.findMonitor, nil
}

func (
	f *fakeMonitorRepository,
) Update(
	ctx context.Context,
	currentMonitor model.Monitor,
) (model.Monitor, error) {
	if f.updateError != nil {
		return model.Monitor{},
			f.updateError
	}

	if f.updateMonitor.ID !=
		0 {
		return f.updateMonitor,
			nil
	}

	return currentMonitor, nil
}

func (
	f *fakeMonitorRepository,
) Delete(
	ctx context.Context,
	id int64,
	userID int64,
) (bool, error) {
	return f.deleteResult,
		f.deleteError
}

func TestMonitorService_Create_DefaultValues(
	t *testing.T,
) {
	repository :=
		&fakeMonitorRepository{}

	service :=
		newMonitorServiceWithRepository(
			repository,
		)

	result, err :=
		service.Create(
			context.Background(),
			7,
			model.CreateMonitorRequest{
				Name: "  Cloudflare Test  ",

				URL: "https://1.1.1.1",
			},
		)

	if err != nil {
		t.Fatalf(
			"Create gagal: %v",
			err,
		)
	}

	if result.Name !=
		"Cloudflare Test" {
		t.Fatalf(
			"name = %q",
			result.Name,
		)
	}

	if result.MonitorType !=
		"website" {
		t.Fatalf(
			"monitor_type = %q, want website",
			result.MonitorType,
		)
	}

	if result.IntervalSeconds !=
		60 {
		t.Fatalf(
			"interval = %d, want 60",
			result.IntervalSeconds,
		)
	}

	if !result.IsActive {
		t.Fatal(
			"monitor baru seharusnya aktif secara default",
		)
	}

	if repository.lastCreateRequest.URL !=
		"https://1.1.1.1" {
		t.Fatalf(
			"URL tersimpan = %q",
			repository.lastCreateRequest.URL,
		)
	}
}

func TestMonitorService_Create_BlocksPrivateURL(
	t *testing.T,
) {
	repository :=
		&fakeMonitorRepository{}

	service :=
		newMonitorServiceWithRepository(
			repository,
		)

	_, err :=
		service.Create(
			context.Background(),
			1,
			model.CreateMonitorRequest{
				Name: "Local",

				URL: "http://127.0.0.1:8081",

				MonitorType: "website",

				IntervalSeconds: 60,
			},
		)

	if !errors.Is(
		err,
		ErrInvalidMonitorURL,
	) {
		t.Fatalf(
			"error = %v, want ErrInvalidMonitorURL",
			err,
		)
	}
}

func TestMonitorService_Create_InvalidInterval(
	t *testing.T,
) {
	repository :=
		&fakeMonitorRepository{}

	service :=
		newMonitorServiceWithRepository(
			repository,
		)

	_, err :=
		service.Create(
			context.Background(),
			1,
			model.CreateMonitorRequest{
				Name: "Test",

				URL: "https://1.1.1.1",

				MonitorType: "website",

				IntervalSeconds: 10,
			},
		)

	if !errors.Is(
		err,
		ErrInvalidMonitorInterval,
	) {
		t.Fatalf(
			"error = %v, want ErrInvalidMonitorInterval",
			err,
		)
	}
}

func TestMonitorService_Create_InvalidType(
	t *testing.T,
) {
	repository :=
		&fakeMonitorRepository{}

	service :=
		newMonitorServiceWithRepository(
			repository,
		)

	_, err :=
		service.Create(
			context.Background(),
			1,
			model.CreateMonitorRequest{
				Name: "Test",

				URL: "https://1.1.1.1",

				MonitorType: "database",

				IntervalSeconds: 60,
			},
		)

	if !errors.Is(
		err,
		ErrInvalidMonitorType,
	) {
		t.Fatalf(
			"error = %v, want ErrInvalidMonitorType",
			err,
		)
	}
}

func TestMonitorService_GetByID_NotFound(
	t *testing.T,
) {
	repository :=
		&fakeMonitorRepository{
			findError: pgx.ErrNoRows,
		}

	service :=
		newMonitorServiceWithRepository(
			repository,
		)

	_, err :=
		service.GetByID(
			context.Background(),
			99,
			1,
		)

	if !errors.Is(
		err,
		ErrMonitorNotFound,
	) {
		t.Fatalf(
			"error = %v, want ErrMonitorNotFound",
			err,
		)
	}
}

func TestMonitorService_Delete_NotFound(
	t *testing.T,
) {
	repository :=
		&fakeMonitorRepository{
			deleteResult: false,
		}

	service :=
		newMonitorServiceWithRepository(
			repository,
		)

	err :=
		service.Delete(
			context.Background(),
			99,
			1,
		)

	if !errors.Is(
		err,
		
		ErrMonitorNotFound,
	) {
		t.Fatalf(
			"error = %v, want ErrMonitorNotFound",
			err,
		)
	}
}

func TestMonitorService_Update_BlocksPrivateURL(
	t *testing.T,
) {
	repository :=
		&fakeMonitorRepository{
			findMonitor: model.Monitor{
				ID: 1,

				UserID: 1,

				Name: "Google",

				URL: "https://1.1.1.1",

				MonitorType: "website",

				IntervalSeconds: 60,

				IsActive: true,
			},
		}

	service :=
		newMonitorServiceWithRepository(
			repository,
		)

	target :=
		"http://192.168.1.1"

	_, err :=
		service.Update(
			context.Background(),
			1,
			1,
			model.UpdateMonitorRequest{
				URL: &target,
			},
		)

	if !errors.Is(
		err,
		ErrInvalidMonitorURL,
	) {
		t.Fatalf(
			"error = %v, want ErrInvalidMonitorURL",
			err,
		)
	}
}
