package device

import (
	"context"
	"errors"
	"time"

	"github.com/art-injener/iot-platform/pkg/logger"
	"github.com/art-injener/iot-platform/pkg/models/device"
)

type DeviceInfoStorageManager struct {
	deviceInfoStorage        DeviceInfoStorage
	deviceInfoHistoryStorage DeviceInfoHistoryStorage
	logger                   *logger.Logger
}

func (d *DeviceInfoStorageManager) AddDeviceInfo(ctx context.Context, deviceID string, deviceInfo *device.ParamsModel) error {
	deviceInfo.CreatedAt = time.Now()
	deviceInfo.UpdatedAt = deviceInfo.CreatedAt
	go func() {
		err := d.deviceInfoHistoryStorage.AddDeviceInfo(ctx, deviceID, deviceInfo)
		if err != nil {
			d.logger.Error().Msgf("AddDeviceInfo in history error %v", err)
		}
	}()

	err := d.deviceInfoStorage.AddDeviceInfo(ctx, deviceID, deviceInfo)
	if err != nil {
		d.logger.Error().Msgf("AddDeviceInfo: add device info error %v", err)
		return err
	}

	return nil
}

func (d *DeviceInfoStorageManager) UpdateDeviceInfoByDeviceID(ctx context.Context, deviceID string, deviceInfo *device.ParamsModel) error {
	deviceInfo.UpdatedAt = time.Now()
	go func() {
		err := d.deviceInfoHistoryStorage.AddDeviceInfo(ctx, deviceID, deviceInfo)
		if err != nil {
			d.logger.Error().Msgf("AddDeviceInfo in history error %v", err)
		}
	}()
	err := d.deviceInfoStorage.UpdateDeviceInfoByDeviceID(ctx, deviceID, deviceInfo)
	if err != nil {
		d.logger.Error().Msgf("UpdateDeviceInfoByDeviceID: update device info error %v", err)
		return err
	}

	return nil
}

func (d *DeviceInfoStorageManager) GetDeviceInfoByDeviceID(ctx context.Context, deviceID string) (*device.ParamsModel, error) {
	deviceInfo, err := d.deviceInfoStorage.GetDeviceInfoByDeviceID(ctx, deviceID)
	if err != nil {
		//d.logger.Error().Msgf("GetDeviceInfoByDeviceID: get device info by device id error %v", err)
		return nil, err
	}

	return deviceInfo, nil
}

func (d *DeviceInfoStorageManager) GetAllDevicesInfo(ctx context.Context) ([]*device.ParamsModel, error) {
	devicesInfo, err := d.deviceInfoStorage.GerAllDevicesInfo(ctx)
	if err != nil {
		d.logger.Error().Msgf("GetAllDevicesInfo: get all devices info error %v", err)
		return nil, err
	}

	return devicesInfo, nil
}

func NewDeviceInfoStorageManager(
	deviceInfoStorage DeviceInfoStorage,
	deviceInfoHistoryStorage DeviceInfoHistoryStorage,
	logger *logger.Logger,
) (*DeviceInfoStorageManager, error) {
	if deviceInfoStorage == nil {
		return nil, errors.New("deviceInfoStorage is empty")
	}

	if deviceInfoHistoryStorage == nil {
		return nil, errors.New("deviceInfoHistoryStorage is empty")
	}

	if logger == nil {
		return nil, errors.New("logger is empty")
	}

	return &DeviceInfoStorageManager{
		deviceInfoStorage:        deviceInfoStorage,
		deviceInfoHistoryStorage: deviceInfoHistoryStorage,
		logger:                   logger,
	}, nil
}
