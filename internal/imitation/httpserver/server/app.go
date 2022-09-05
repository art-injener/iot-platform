package server

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/art-injener/iot-platform/internal/config"
	"github.com/art-injener/iot-platform/internal/imitation"
)

type WebApp struct {
	httpServer *http.Server
	cfg        *config.Config
	imitator   *imitation.Imitator
}

// NewApp - конструктор для имитатора с поддержкой вывода инфы на web интерфейс
func NewApp(cfg *config.Config) (*WebApp, error) {
	imitator, err := imitation.NewImitator(cfg.Log)

	if err != nil {
		return nil, err
	}

	return &WebApp{
		cfg:      cfg,
		imitator: imitator,
	}, nil
}

func (a *WebApp) Run(finish chan struct{}) error {

	var abort = make(chan struct{})

	router := gin.Default()
	router.Use(
		gin.Recovery(),
	)
	router.LoadHTMLGlob("assets/templates/*")

	// корневйо маршрут
	router.GET("/", func(c *gin.Context) {
		c.HTML(
			http.StatusOK,
			"index.html",
			gin.H{},
		)
	})

	// получение данных о маяках
	router.GET("/data.json", func(c *gin.Context) {
		type beaconInfo struct {
			BeaconID   string  `json:"beaconID"`
			Latitude   float64 `json:"lat"`
			Longitude  float64 `json:"lon"`
			Speed      int     `json:"speed"`
			Azim       float32 `json:"azim"`
			SeanceTime string  `json:"seanceTime"`
		}

		if a.imitator == nil {
			c.JSON(http.StatusOK, []*beaconInfo{})
			return
		}

		info := a.imitator.GetDeviceImitInfo()
		dots := make([]*beaconInfo, 0, len(info))
		for i := 0; i < len(info); i++ {
			if info[i].ID != "" {
				dots = append(dots, &beaconInfo{
					BeaconID:   info[i].ID,
					Latitude:   info[i].Latitude,
					Longitude:  info[i].Longitude,
					Speed:      info[i].Speed,
					Azim:       info[i].DirectionMove,
					SeanceTime: time.Unix(info[i].SystemTime, 0).Local().Format(time.RFC1123),
				})
			}
		}
		c.JSON(http.StatusOK, dots)
	})

	// получение данных о маяках
	router.GET("/region.json", func(c *gin.Context) {
		type regionInfo struct {
			RegionID  string  `json:"regionID"`
			Latitude  float64 `json:"lat"`
			Longitude float64 `json:"lon"`
			Radius    float32 `json:"radius"`
		}

		if a.imitator == nil {
			c.JSON(http.StatusOK, []*regionInfo{})
			return
		}

		data := a.imitator.TrackCache.GetAllRegion()

		dots := make([]*regionInfo, 0, len(data))
		for k, v := range data {
			dots = append(dots, &regionInfo{
				RegionID:  strconv.Itoa(int(k)),
				Latitude:  v.Center.Lat(),
				Longitude: v.Center.Lng(),
				Radius:    v.Radius,
			})
		}
		c.JSON(http.StatusOK, dots)
	})

	a.imitator.CreateVirtualDevices(a.cfg, a.cfg.Phones)
	go func() {
		a.imitator.StartImitation(finish, abort)
	}()

	a.httpServer = &http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	return nil
}

func (a *WebApp) Stop() {
	if a == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := a.httpServer.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
