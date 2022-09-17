package webserver

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/art-injener/iot-platform/internal/config"
	"github.com/art-injener/iot-platform/internal/handlers"
	lg "github.com/art-injener/iot-platform/pkg/logger"
	"github.com/art-injener/iot-platform/pkg/storage/device"
	"github.com/art-injener/iot-platform/pkg/storage/postgres"
)

type App struct {
	cfg                      *config.Config
	router                   *gin.Engine
	deviceInfoStorageManager *device.DeviceInfoStorageManager
	postgres                 *postgres.Postgres
}

func (a *App) Initialize(ctx context.Context) {
	a.initConfig()
	a.initDB(ctx)
	a.initServices()
	a.initRoutes()
}

func (a *App) Run(ctx context.Context) {
	port := ":" + a.cfg.WebServerPort
	httpServer := &http.Server{
		Addr:    port,
		Handler: a.router,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		<-ctx.Done()
		a.cfg.Log.Info().Msg("Shutting down HTTP server...")
		tctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer func() {
			cancel()
			close(idleConnsClosed)
		}()
		if err := httpServer.Shutdown(tctx); err != nil {
			a.cfg.Log.Error().Msgf("Shutdown error: %v", err)
		}
	}()
	a.cfg.Log.Info().Msgf("Start HTTP server on %q", port)
	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		a.cfg.Log.Error().Msgf("Can't run HTTP server: %v", err)
	}
	<-idleConnsClosed
}

func (a *App) initConfig() {
	var err error
	a.cfg, err = config.GetConfig("configs")
	if err != nil {
		panic(err)
	}

	a.cfg.Log = lg.NewConsole(a.cfg.LogLevel == config.DebugLevel)
}

func (a *App) initDB(ctx context.Context) {
	var err error
	a.postgres, err = postgres.NewClient(ctx, a.cfg)
	if err != nil {
		panic(err)
	}
}

func (a *App) initServices() {
	var err error
	deviceInfoStorage, err := device.NewDeviceInfoStorage(a.postgres)
	if err != nil {
		panic(err)
	}

	deviceInfoHistoryStorage, err := device.NewDeviceInfoHistoryStorage(a.postgres)
	if err != nil {
		panic(err)
	}

	a.deviceInfoStorageManager, err = device.NewDeviceInfoStorageManager(deviceInfoStorage, deviceInfoHistoryStorage, a.cfg.Log)
	if err != nil {
		panic(err)
	}
}

func (a *App) initRoutes() {
	router := gin.Default()

	router.Use(
		gin.Recovery(),
	)
	router.LoadHTMLGlob("assets/templates/*")

	// корневйо маршрут
	router.GET("/", func(c *gin.Context) {
		c.HTML(
			http.StatusOK,
			"client_app.html",
			gin.H{},
		)
	})

	// получение данных о маяках
	deviceDataHandler := handlers.DeviceDataHandler{
		DeviceInfoManager: a.deviceInfoStorageManager,
		Logger:            a.cfg.Log,
	}
	router.GET("/data.json", deviceDataHandler.DeviceData)

	a.router = router
}
