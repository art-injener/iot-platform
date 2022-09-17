package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/art-injener/iot-platform/pkg/logger"
	"github.com/art-injener/iot-platform/pkg/storage/device"
)

type DeviceDataHandler struct {
	DeviceInfoManager *device.DeviceInfoStorageManager
	Logger            *logger.Logger
}

func (d *DeviceDataHandler) DeviceData(ctx *gin.Context) {
	type beaconInfo struct {
		BeaconID   string  `json:"beaconID"`
		Latitude   float64 `json:"lat"`
		Longitude  float64 `json:"lon"`
		Speed      int     `json:"speed"`
		Azim       float32 `json:"azim"`
		SeanceTime string  `json:"seanceTime"`
		CreatedAt  string  `json:"createdAt"`
	}

	infos, err := d.DeviceInfoManager.GetAllDevicesInfo(context.Background())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "internal server error")
	}

	dots := make([]*beaconInfo, 0, len(infos))
	for i := range infos {
		if infos[i].ID != "" {
			dots = append(dots, &beaconInfo{
				BeaconID:   infos[i].ID,
				Latitude:   infos[i].Latitude,
				Longitude:  infos[i].Longitude,
				Speed:      infos[i].Speed,
				Azim:       infos[i].DirectionMove,
				SeanceTime: time.Unix(infos[i].SystemTime, 0).Local().Format(time.RFC850),
				CreatedAt:  infos[i].CreatedAt.Format(time.RFC850),
			})
		}
	}
	ctx.JSON(http.StatusOK, dots)
}
