package device

import (
	"context"
	"strconv"

	"github.com/pkg/errors"

	"github.com/art-injener/iot-platform/pkg/models/device"
	"github.com/art-injener/iot-platform/pkg/storage/postgres"
)

const (
	addDeviceSettingsSQL = `
	INSERT INTO
		device_settings (
		device_id, 
		version, 
		timezone,
		wui,
		gpst
	)
	VALUES (
		$1,
		$2,
		$3,
		$4,
		$5 
	)`

	getSettingsByDeviceIDSQL = `
	SELECT
		device_id, 
		version, 
		timezone,
		wui,
		gpst
	FROM
		device_settings 
	WHERE
		device_id = $1 ;`

	updateVersionByDeviceIDSQL = `
	UPDATE 
		device_settings
	SET
		version=$1
	WHERE 
		device_id = $2 ;`
)

type DeviceSettingsStorage interface {
	AddDeviceSettings(context.Context, *device.SettingsModel) error
	GetSettingsByDeviceID(context.Context, int) (*device.SettingsModel, error)
	UpdateVersionByDeviceID(context.Context, string, int) error
}

type deviceSettingsStorage struct {
	postgres *postgres.Postgres
}

var _ DeviceSettingsStorage = &deviceSettingsStorage{}

func (d *deviceSettingsStorage) AddDeviceSettings(ctx context.Context, settings *device.SettingsModel) error {
	_, err := d.postgres.Exec(ctx, addDeviceSettingsSQL, settings.ID, settings.Version, settings.TZ, settings.WUI, settings.GPST)
	if err != nil {
		return err
	}
	return nil
}

func (d *deviceSettingsStorage) GetSettingsByDeviceID(ctx context.Context, deviceID int) (*device.SettingsModel, error) {
	settings := device.SettingsModel{}
	row := d.postgres.QueryRow(ctx, getSettingsByDeviceIDSQL, deviceID)
	var intID int
	if err := row.Scan(&intID, &settings.Version, &settings.TZ, &settings.WUI, &settings.GPST); err != nil {
		return nil, err
	}
	settings.ID = strconv.Itoa(intID)
	return &settings, nil
}

func (d *deviceSettingsStorage) UpdateVersionByDeviceID(ctx context.Context, version string, deviceID int) error {
	_, err := d.postgres.Exec(ctx, updateVersionByDeviceIDSQL, version, deviceID)
	if err != nil {
		return err
	}
	return nil
}

func NewDeviceSettingsStorage(postgres *postgres.Postgres) (DeviceSettingsStorage, error) {
	if postgres == nil {
		return nil, errors.New("postgres connection is empty")
	}

	return &deviceSettingsStorage{postgres: postgres}, nil
}
