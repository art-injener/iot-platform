package device

import (
	"context"

	"github.com/pkg/errors"

	"github.com/art-injener/iot-platform/pkg/models/device"
	"github.com/art-injener/iot-platform/pkg/storage/postgres"
)

const (
	addDeviceInfoSQL = `
	INSERT INTO 
		device_info (device_id, created_at, updated_at, device)
	VALUES
		($1, $2, $3, $4) 
	ON CONFLICT DO NOTHING ;`

	getDeviceInfoByDeviceIDSQL = `
	SELECT 
		device, created_at, updated_at
	FROM
		device_info
	WHERE device_id = $1 ;`

	updateDeviceInfoByDeviceIDSQL = `
	UPDATE 
		device_info
	SET 
		device = $1, updated_at = $2
	WHERE 
		device_id = $3 ;`

	getAllDevicesInfoSQL = `
	SELECT 
		device_id,created_at, device
	FROM
		device_info ;`
)

var _ DeviceInfoStorage = &deviceInfoStorage{}

type DeviceInfoStorage interface {
	AddDeviceInfo(context.Context, string, *device.ParamsModel) error
	GetDeviceInfoByDeviceID(context.Context, string) (*device.ParamsModel, error)
	UpdateDeviceInfoByDeviceID(context.Context, string, *device.ParamsModel) error
	GerAllDevicesInfo(context.Context) ([]*device.ParamsModel, error)
}

type deviceInfoStorage struct {
	postgres *postgres.Postgres
}

func (d *deviceInfoStorage) AddDeviceInfo(ctx context.Context, deviceID string, deviceInfo *device.ParamsModel) error {
	_, err := d.postgres.Exec(ctx, addDeviceInfoSQL, deviceID, deviceInfo.CreatedAt, deviceInfo.UpdatedAt, deviceInfo)
	return err
}

func (d *deviceInfoStorage) GetDeviceInfoByDeviceID(ctx context.Context, deviceID string) (*device.ParamsModel, error) {
	deviceInfo := &device.ParamsModel{}
	row := d.postgres.QueryRow(ctx, getDeviceInfoByDeviceIDSQL, deviceID)
	if err := row.Scan(deviceInfo, &deviceInfo.CreatedAt, &deviceInfo.UpdatedAt); err != nil {
		return deviceInfo, err
	}
	return deviceInfo, nil
}

func (d *deviceInfoStorage) UpdateDeviceInfoByDeviceID(ctx context.Context, deviceID string, deviceInfo *device.ParamsModel) error {
	_, err := d.postgres.Exec(ctx, updateDeviceInfoByDeviceIDSQL, deviceInfo, deviceInfo.UpdatedAt, deviceID)
	return err
}

func (d *deviceInfoStorage) GerAllDevicesInfo(ctx context.Context) ([]*device.ParamsModel, error) {
	rows, err := d.postgres.Query(ctx, getAllDevicesInfoSQL)
	if err != nil {
		return nil, err
	}

	devicesInfo := make([]*device.ParamsModel, 0, 1000)
	for rows.Next() {
		info := &device.ParamsModel{}
		err := rows.Scan(&info.ID, &info.CreatedAt, &info)
		if err != nil {
			continue
		}
		devicesInfo = append(devicesInfo, info)
	}
	return devicesInfo, nil
}

func NewDeviceInfoStorage(postgres *postgres.Postgres) (DeviceInfoStorage, error) {
	if postgres == nil {
		return nil, errors.New("postgres connection is empty")
	}

	return &deviceInfoStorage{postgres: postgres}, nil
}
