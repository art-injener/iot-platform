package device

import (
	"context"

	"github.com/pkg/errors"

	"github.com/art-injener/iot-platform/pkg/models/device"
	"github.com/art-injener/iot-platform/pkg/storage/postgres"
)

const (
	addDeviceInfoHistorySQL = `
	INSERT INTO 
		device_info_history (device_id, created_at, updated_at, device)
	VALUES
		($1, $2, $3, $4) 
	ON CONFLICT DO NOTHING ;`
)

var _ DeviceInfoHistoryStorage = &deviceInfoHistoryStorage{}

type DeviceInfoHistoryStorage interface {
	AddDeviceInfo(context.Context, string, *device.ParamsModel) error
}

type deviceInfoHistoryStorage struct {
	postgres *postgres.Postgres
}

func (d *deviceInfoHistoryStorage) AddDeviceInfo(ctx context.Context, deviceID string, deviceInfo *device.ParamsModel) error {
	_, err := d.postgres.Exec(ctx, addDeviceInfoHistorySQL, deviceID, deviceInfo.CreatedAt, deviceInfo.UpdatedAt, deviceInfo)
	return err
}

func NewDeviceInfoHistoryStorage(postgres *postgres.Postgres) (DeviceInfoHistoryStorage, error) {
	if postgres == nil {
		return nil, errors.New("postgres connection is empty")
	}

	return &deviceInfoHistoryStorage{postgres: postgres}, nil
}
