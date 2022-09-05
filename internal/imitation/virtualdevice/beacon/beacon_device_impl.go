package beacon

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/art-injener/iot-platform/internal/imitation/virtualdevice/beacon/params"
	"github.com/art-injener/iot-platform/internal/models/imitator"
	"github.com/art-injener/iot-platform/pkg/logger"
	"github.com/art-injener/iot-platform/pkg/models/device"
	util "github.com/art-injener/iot-platform/util/helper"
)

const MaxPacketSize = 1 * 1024

type VirtualBeaconImpl struct {
	DevParam   *params.Parameters
	log        *logger.Logger
	seanceTime time.Time
}

func NewVirtualBeacon(deviceID string, launch, wakeUpInterval int, log *logger.Logger) *VirtualBeaconImpl {
	seance := time.Now().Add(time.Duration(rand.Intn(launch)) * time.Minute)
	seance = seance.Add(time.Duration(rand.Intn(60)) * time.Second)
	seance = seance.Add(time.Duration(rand.Intn(1000)) * time.Millisecond)

	return &VirtualBeaconImpl{
		DevParam:   params.NewParameters(deviceID, wakeUpInterval),
		log:        log,
		seanceTime: seance,
	}
}

func (v *VirtualBeaconImpl) GetID() string {
	return v.DevParam.ConstParam.ID
}

func (v *VirtualBeaconImpl) IsNeedWakeUp() bool {
	return time.Since(v.seanceTime) > 0
}

func (v *VirtualBeaconImpl) GetSeanceTime() time.Time {
	return v.seanceTime
}

// Send - функция отправки данных по сети. Для корректной работы требуется наличие установленного соединения
// Первоначально выполняется упаковка параметров устройства
// затем отправка и чтение ответа.
func (v *VirtualBeaconImpl) Send(conn net.Conn) (bool, error) { // выполняем сериализация параметров устройства
	str, _, err := v.Serialize()
	if err != nil {
		return false, errors.Wrap(err, "error serialize data for send in Send()")
	}

	// TODO: добавить проверку isConnected()
	// отправка пакета данных на сервер
	_, err = conn.Write([]byte(str))
	if err != nil {
		return false, errors.Wrap(err, "error write data to network in Send()")
	}
	if v.log != nil {
		v.log.Debug().Msgf("Send device data: \n\t%s\t", str)
	}

	////// читаем ответ от сервака
	//buf := make([]byte, MaxPacketSize)
	//rb, err := conn.Read(buf)
	//if err != nil {
	//	return false, errors.Wrap(err, "error read data from network in Send()")
	//}
	//if v.log != nil {
	//	v.log.Debug().Msgf("Get response : \t%s\t", string(buf[:rb]))
	//}
	//
	//if v.DevParam.ConstParam.WUI == 0 {
	//	v.DevParam.ConstParam.WUI = 1
	//}
	//
	//v.seanceTime = v.seanceTime.Add(time.Duration(v.DevParam.ConstParam.WUI) * time.Minute)
	//v.DevParam.MutableParam.SeanceTime = v.seanceTime.UTC().Unix()
	//
	//isSend, err := v.checkResponseOnSuccessState(string(buf[:rb]), crc)
	//if err != nil {
	//	return false, errors.Wrap(err, "error in validation server response")
	//}
	//
	//if isSend {
	//	if err = v.wakeUpIntervalParser(string(buf[:rb])); err != nil {
	//		return false, err
	//	}
	//	if err = v.satelliteSearchTimeParser(string(buf[:rb])); err != nil {
	//		return false, err
	//	}
	//}
	//buf = nil

	return true, nil
}

func (v *VirtualBeaconImpl) UpdateGeoParamByExternalImitator(imitParam *imitator.GeoParamsModel) error {
	return v.DevParam.UpdateGeoParamByExternalImitator(imitParam)
}

// Serialize - выполняем сериализацию параметров устройства
// формат сериализации :
// PARAM_NAME_1=VALUE_1&PARAM_NAME_2=VALUE_2&...PARAM_NAME_N=VALUE_N&*CRC
// последним символом, по которому счиатется crc является & (некторые устройства в конце могут присылать &&)
func (v *VirtualBeaconImpl) Serialize() (message string, msgCRC uint8, err error) {
	s := v.DevParam.Serialize()
	crc := v.crc(s)
	s += fmt.Sprintf("*%d", crc)

	return s, crc, nil
}

// checkResponseOnSuccessState - проверяем ответ от сервера  на отсутствие ошибок
// пример успешного ответа : RE=0&CRC=20&STIME=210729140544&P=0000&&*180
// пример ответа с ошибкой : ERR&CRC=196&&*123
func (v *VirtualBeaconImpl) checkResponseOnSuccessState(response string, _ uint8) (bool, error) {
	// посчитаем crc ответа
	indx := strings.IndexRune(response, '*')
	msg := response[:indx]
	rc := response[indx+1:]
	responseCRC, err := strconv.Atoi(rc)
	if err != nil {
		return false, errors.Wrap(err, " error Atoi CRC")
	}
	if v.crc(msg) != uint8(responseCRC) {
		return false, errors.New("CRC does not match")
	}
	if ret := strings.HasPrefix(response, "ERR"); ret {
		return false, errors.New(fmt.Sprintf("Server answer error : %s ", response))
	}

	// TODO : добавить проверку на RE=0
	return true, nil
}

func (v *VirtualBeaconImpl) GetDeviceParameters() device.ParamsModel {
	devInf := v.DevParam.GetDeviceInfo()
	return devInf
}

func (v *VirtualBeaconImpl) String() string {
	return fmt.Sprintf("Beacon: \n\tID=%s\n", v.DevParam.ConstParam.ID)
}

// crc - рассчет контрольной суммы строки
func (v *VirtualBeaconImpl) crc(s string) uint8 {
	if v == nil {
		return 0
	}
	var checksum uint8
	b := []byte(s)
	for i := 0; i < len(b); i++ {
		checksum += b[i]
		if checksum < b[i] {
			checksum++
		}
	}

	return checksum
}

// wakeUpIntervalParser - Временной интервал просыпания
func (v *VirtualBeaconImpl) wakeUpIntervalParser(data string) (err error) {

	const wuiToken string = "WUI=" // Интервал извещения

	value := util.ValueExtractor(data, wuiToken)
	if len(value) == 0 {
		return nil
	}

	valueWUI, err := strconv.Atoi(value)
	if err != nil {
		return err
	}

	v.DevParam.ConstParam.WUI = valueWUI

	if v.log != nil {
		v.log.Info().Msgf("setting a new parameter value WUI = %d ", v.DevParam.ConstParam.WUI)
	}
	return nil
}

// satelliteSearchTimeParser - Парсинг параметра "лимит поиска спутниковых координат"
func (v *VirtualBeaconImpl) satelliteSearchTimeParser(data string) (err error) {

	const token string = "GPST=" // Лимит поиска спутниковых координат

	value := util.ValueExtractor(data, token)
	if len(value) == 0 {
		return nil
	}

	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return err
	}

	v.DevParam.ConstParam.GPST = valueInt

	if v.log != nil {
		v.log.Info().Msgf("setting a new parameter value GPST = %d ", v.DevParam.ConstParam.GPST)
	}

	return nil
}
