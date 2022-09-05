package virtualdevice

import (
	"net"

	"github.com/art-injener/iot-platform/internal/models/imitator"
	"github.com/art-injener/iot-platform/pkg/models/device"
)

type VirtualDevice interface {
	GetID() string
	Send(conn net.Conn) (bool, error)
	IsNeedWakeUp() bool
	Serialize() (message string, msgCRC uint8, err error)
	GetDeviceParameters() device.ParamsModel
	UpdateGeoParamByExternalImitator(imitParam *imitator.GeoParamsModel) error
}
