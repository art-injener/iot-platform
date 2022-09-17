package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	pgx_5 "github.com/jackc/pgx/v5"

	"github.com/art-injener/iot-platform/pkg/logger"
	"github.com/art-injener/iot-platform/pkg/models/rmq"
	"github.com/art-injener/iot-platform/pkg/rabbitmq"
	"github.com/art-injener/iot-platform/pkg/storage/device"
)

type DeviceInfoConsumer struct {
	consumer              *rabbitmq.Consumer
	deviceInfoManager     *device.DeviceInfoStorageManager
	deviceSettingsStorage device.DeviceSettingsStorage
	logger                *logger.Logger
}

func (d *DeviceInfoConsumer) Poll(ctx context.Context) {
	go func() {
		d.logger.Info().Msgf("Consumer ready, PID: %d", os.Getpid())
		for msg := range d.consumer.MessageChannel {
			dm := rmq.DeviceMessageModel{}

			err := json.Unmarshal(msg.Body, &dm)
			if err != nil {
				return
			}
			d.logger.Debug().Msgf("Received a message: %+v", dm)

			deviceInfo, err := d.deviceInfoManager.GetDeviceInfoByDeviceID(ctx, dm.DevSettings.ID)
			if err != nil {
				if errors.Is(err, pgx_5.ErrNoRows) {
					err := d.deviceInfoManager.AddDeviceInfo(ctx, dm.DevSettings.ID, &dm.DevParameters)
					if err != nil {
						d.logger.Error().Msgf("AddDeviceInfo err %v", err)
						continue
					}
				}
				continue
			}

			dm.DevParameters.CreatedAt = deviceInfo.CreatedAt
			err = d.deviceInfoManager.UpdateDeviceInfoByDeviceID(ctx, dm.DevSettings.ID, &dm.DevParameters)
			if err != nil {
				d.logger.Error().Msgf("UpdateDeviceInfoByDeviceID err %v", err)
				continue
			}

			if err := msg.Ack(false); err != nil {
				d.logger.Error().Msgf("Error acknowledging message : %s", err)
			}
		}
	}()
}

func (d *DeviceInfoConsumer) Stop() {
	d.consumer.Stop()
}

func NewDeviceInfoConsumer(
	consumer *rabbitmq.Consumer,
	deviceInfoManager *device.DeviceInfoStorageManager,
	deviceSettingsStorage device.DeviceSettingsStorage,
	logger *logger.Logger,
) (*DeviceInfoConsumer, error) {
	if consumer == nil {
		return nil, errors.New("RabbitMQ consumer is empty")
	}

	if deviceInfoManager == nil {
		return nil, errors.New("deviceInfoManager is empty")
	}

	if deviceSettingsStorage == nil {
		return nil, errors.New("deviceSettingsStorage is empty")
	}

	if logger == nil {
		return nil, errors.New("logger is empty")
	}

	return &DeviceInfoConsumer{
		consumer:              consumer,
		deviceInfoManager:     deviceInfoManager,
		deviceSettingsStorage: deviceSettingsStorage,
		logger:                logger,
	}, nil
}
