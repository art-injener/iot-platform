package rmq

import (
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/art-injener/iot-platform/pkg/models/device"
	"github.com/art-injener/iot-platform/util"
)

type DeviceMessageModel struct {
	ID            uuid.UUID            `json:"id"`
	SagaID        uuid.UUID            `json:"saga_id"`
	DevSettings   device.SettingsModel `json:"dev_settings"`
	DevParameters device.ParamsModel   `json:"dev_parameters"`
}

func NewDeviceMessage(str string) (*DeviceMessageModel, error) {
	dm := DeviceMessageModel{
		ID:     uuid.New(),
		SagaID: uuid.New(),
	}
	err := dm.unserialize(str)
	if err != nil {
		return nil, err
	}

	return &dm, nil
}

// ID=89991234650&VER=0.0.1&TZ=220&WUI=20&GPST=3&STIME=220828003254&BAL=19.85&TE=28.0&VB=99%(5.94V)&IC=4&SQ=86&LA=59.418731&LAD=N&LO=105.700373&LOD=E&SPD=50&DM=102.0&GT=220827213254&&*150
func (d *DeviceMessageModel) unserialize(str string) error {
	if err := d.idParser(str); err != nil {
		return err
	}

	if err := d.geoParamParser(str); err != nil {
		return err
	}

	if err := d.trafficParamParser(str); err != nil {
		return err
	}
	return nil
}

// idParser - ID устройства
func (d *DeviceMessageModel) idParser(data string) error {

	const idToken string = "ID=" // ID устройства

	value := util.ValueExtractor(data, idToken)
	if len(value) == 0 {
		return nil
	}
	d.DevSettings.ID = value
	d.DevParameters.ID = value
	return nil
}

// balanceParser - Парсер баланса
func (d *DeviceMessageModel) balanceParser(data string) error {

	const token string = "BAL=" // ID устройства

	v := util.ValueExtractor(data, token)
	if len(v) == 0 {
		return nil
	}

	value, err := strconv.ParseFloat(v, 32)
	if err != nil {
		return err
	}

	d.DevParameters.BalanceSim = float32(value)
	return nil
}

func (d *DeviceMessageModel) geoParamParser(data string) error {

	const latToken string = "LA="
	const latDToken string = "LAD="

	const lonToken string = "LO="
	const lonDToken string = "LOD="

	var errLat, errLon error
	var valueLat, valueLon float64

	// широта
	v := util.ValueExtractor(data, latToken)
	if valueLat, errLat = strconv.ParseFloat(v, 64); errLat != nil {
		return errLat
	}

	// определение полушария
	if v = util.ValueExtractor(data, latDToken); strings.Contains(v, "S") {
		valueLat *= -1
	}
	d.DevParameters.Latitude = valueLat
	d.DevParameters.LAD = v

	// долгота
	v = util.ValueExtractor(data, lonToken)

	if valueLon, errLon = strconv.ParseFloat(v, 64); errLon != nil {
		return errLon
	}

	// определение полушария
	if v = util.ValueExtractor(data, lonDToken); strings.Contains(v, "W") {
		valueLon *= -1
	}
	d.DevParameters.Longitude = valueLon
	d.DevParameters.LOD = v

	return nil
}

// &SD=14&MO=28.01&
func (d *DeviceMessageModel) trafficParamParser(data string) error {

	const speedToken string = "SPD="  // средняя скорость
	const azimuthToken string = "DM=" // направление движения

	var errSpeed, errAzim error
	v := util.ValueExtractor(data, speedToken)
	valueSpeed, errSpeed := strconv.Atoi(v)

	v = util.ValueExtractor(data, azimuthToken)
	valueAzim, errAzim := strconv.ParseFloat(v, 32)

	if errSpeed != nil || errAzim != nil {
		return errors.New("error parsing data")
	}

	d.DevParameters.Speed = valueSpeed
	d.DevParameters.DirectionMove = float32(valueAzim)

	return nil
}
