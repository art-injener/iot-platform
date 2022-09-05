package params

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/art-injener/iot-platform/internal/models/imitator"
	"github.com/art-injener/iot-platform/pkg/models/device"
)

const (
	defaultMessageLen = 512
)

// Parameters - структура для хранения параметров устройства.
type Parameters struct {
	ConstParam   *device.SettingsModel
	MutableParam *device.ParamsModel
}

// NewParameters  - конструктор структуры параметров устройства.
func NewParameters(deviceID string, wakeUpInterval int) *Parameters {
	return &Parameters{
		ConstParam:   device.NewDeviceSettings(deviceID, wakeUpInterval),
		MutableParam: device.NewDeviceParams(deviceID),
	}
}

func (p *Parameters) Serialize() string {
	var concatRune = '&'
	builder := strings.Builder{}
	builder.Grow(defaultMessageLen)

	p.MutableParam.SystemTime = time.Now().UTC().Unix()

	p.serializeDeviceSettings(p.ConstParam, &builder, concatRune)
	p.serializeDeviceMutableParams(p.MutableParam, &builder, concatRune)
	builder.WriteRune(concatRune)

	return builder.String()
}

func (p *Parameters) UpdateGeoParamByExternalImitator(imitParam *imitator.GeoParamsModel) error {
	if p == nil || imitParam == nil {
		return errors.New("instance Parameters==nil")
	}

	p.MutableParam.UpdateAll()
	p.MutableParam.IntervalCounter++

	p.MutableParam.Longitude = imitParam.LongitudeBeacon
	p.MutableParam.Latitude = imitParam.LatitudeBeacon
	p.MutableParam.DirectionMove = imitParam.Azimuth
	p.MutableParam.Speed = imitParam.AverageSpeed

	return nil
}

func (p *Parameters) GetDeviceInfo() device.ParamsModel {
	v := *p.MutableParam
	v.SeanceTime += int64((time.Duration(p.ConstParam.WUI) * time.Minute).Seconds())

	return v
}

func (p *Parameters) serializeDeviceSettings(settingsModel *device.SettingsModel, builder *strings.Builder, specChar rune) {
	// ID
	builder.WriteString("ID=")
	builder.WriteString(settingsModel.ID)
	builder.WriteRune(specChar)

	// VER
	builder.WriteString("VER=")
	builder.WriteString(settingsModel.Version)
	builder.WriteRune(specChar)

	// TZ
	builder.WriteString("TZ=")
	builder.WriteString(strconv.Itoa(settingsModel.TZ))
	builder.WriteRune(specChar)

	// WUI
	builder.WriteString("WUI=")
	builder.WriteString(strconv.Itoa(settingsModel.WUI))
	builder.WriteRune(specChar)

	// GPST
	builder.WriteString("GPST=")
	builder.WriteString(strconv.Itoa(settingsModel.GPST))
	builder.WriteRune(specChar)
}

func (p *Parameters) serializeDeviceMutableParams(mutableParam *device.ParamsModel, builder *strings.Builder, specChar rune) {
	// STIME
	// формат описан в src/time/format.go
	builder.WriteString("STIME=")
	builder.WriteString(time.Unix(mutableParam.SystemTime, 0).Local().Format("060102150405"))
	builder.WriteRune(specChar)

	// BAL
	fmt.Fprintf(builder, "BAL=%.2f&", mutableParam.BalanceSim)

	// TE
	fmt.Fprintf(builder, "TE=%.1f&", mutableParam.Temperature)

	// VB "100%(5.95V)"
	builder.WriteString("VB=")
	builder.WriteString(mutableParam.BatteryChargeToString())
	builder.WriteRune(specChar)

	// IC
	builder.WriteString("IC=")
	builder.WriteString(strconv.FormatUint(mutableParam.IntervalCounter, 10))
	builder.WriteRune(specChar)

	// SQ
	builder.WriteString("SQ=")
	builder.WriteString(strconv.Itoa(mutableParam.SignalQuality))
	builder.WriteRune(specChar)

	// LA
	fmt.Fprintf(builder, "LA=%f&", mutableParam.Latitude)

	// LAD
	builder.WriteString("LAD=")
	builder.WriteString(mutableParam.LAD)
	builder.WriteRune(specChar)

	// LO
	fmt.Fprintf(builder, "LO=%f&", mutableParam.Longitude)

	// LOD
	builder.WriteString("LOD=")
	builder.WriteString(mutableParam.LOD)
	builder.WriteRune(specChar)

	// SPD
	builder.WriteString("SPD=")
	builder.WriteString(strconv.Itoa(mutableParam.Speed))
	builder.WriteRune(specChar)

	// DM
	fmt.Fprintf(builder, "DM=%.1f&", mutableParam.DirectionMove)

	// GT
	builder.WriteString("GT=")
	builder.WriteString(time.Unix(mutableParam.SystemTime, 0).UTC().Format("060102150405")) // формат описан в src/time/format.go
	builder.WriteRune(specChar)
}
