package app

import (
	"context"

	"github.com/art-injener/iot-platform/internal/config"
	device_amqp "github.com/art-injener/iot-platform/internal/rabbitmq"
	lg "github.com/art-injener/iot-platform/pkg/logger"
	"github.com/art-injener/iot-platform/pkg/rabbitmq"
	"github.com/art-injener/iot-platform/pkg/storage/device"
	"github.com/art-injener/iot-platform/pkg/storage/postgres"
)

type App struct {
	config                   *config.Config
	dbConn                   *postgres.Postgres
	deviceInfoStorageManager *device.DeviceInfoStorageManager
	deviceSettingsStorage    device.DeviceSettingsStorage
	deviceInfoConsumer       *device_amqp.DeviceInfoConsumer
}

func NewApp() *App {
	return &App{}
}

func (a *App) Initialization(ctx context.Context) {
	a.initConfig()
	a.initDB(ctx)
	a.initDeviceInfoStorageManager()
	a.initDeviceSettingsStorage()
	a.initDeviceInfoConsumer()
}

func (a *App) Run(ctx context.Context) {
	defer a.deviceInfoConsumer.Stop()

	a.deviceInfoConsumer.Poll(ctx)

	<-ctx.Done()
	a.config.Log.Info().Msgf("Shutdown consumer...")
}

func (a *App) initConfig() {
	var err error
	a.config, err = config.GetConfig("configs")
	if err != nil {
		panic(err)
	}

	a.config.Log = lg.NewConsole(a.config.LogLevel == config.DebugLevel)
}

func (a *App) initDB(ctx context.Context) {
	var err error
	a.dbConn, err = postgres.NewClient(ctx, a.config)
	if err != nil {
		panic(err)
	}
}

func (a *App) initDeviceInfoStorageManager() {
	deviceInfoStorage, err := device.NewDeviceInfoStorage(a.dbConn)
	if err != nil {
		panic(err)
	}

	deviceInfoHistoryStorage, err := device.NewDeviceInfoHistoryStorage(a.dbConn)
	if err != nil {
		panic(err)
	}

	a.deviceInfoStorageManager, err = device.NewDeviceInfoStorageManager(deviceInfoStorage, deviceInfoHistoryStorage, a.config.Log)
	if err != nil {
		panic(err)
	}
}

func (a *App) initDeviceSettingsStorage() {
	var err error
	a.deviceSettingsStorage, err = device.NewDeviceSettingsStorage(a.dbConn)
	if err != nil {
		panic(err)
	}
}

func (a *App) initDeviceInfoConsumer() {
	consumer := &rabbitmq.Consumer{
		Config: a.config.RabbitConfig,
		Logger: a.config.Log,
	}
	err := consumer.Initialize()
	if err != nil {
		panic(err)
	}

	a.deviceInfoConsumer, err = device_amqp.NewDeviceInfoConsumer(consumer, a.deviceInfoStorageManager, a.deviceSettingsStorage, a.config.Log)
	if err != nil {
		panic(err)
	}
}
