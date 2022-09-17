package imitation

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/art-injener/iot-platform/internal/config"

	"github.com/pkg/errors"

	"github.com/art-injener/iot-platform/internal/imitation/trackcache"
	"github.com/art-injener/iot-platform/internal/imitation/virtualdevice"
	"github.com/art-injener/iot-platform/internal/imitation/virtualdevice/beacon"
	"github.com/art-injener/iot-platform/pkg/logger"
	"github.com/art-injener/iot-platform/pkg/models/device"
	"github.com/art-injener/iot-platform/util"
)

const geoDataFilePath = "assets/data/geoData.json"

type Imitator struct {
	mu         sync.RWMutex
	log        *logger.Logger
	devices    []virtualdevice.VirtualDevice
	TrackCache *trackcache.Cache
	netConfig  *config.NetworkConfig
	DeviceInfo []device.ParamsModel
}

// NewImitator - конструктор имитатора
func NewImitator(log *logger.Logger) (*Imitator, error) {
	var polygons map[uint16]util.GeoPolygon
	var err error

	log.Info().Msg("[IMITATOR]::load geozone from file")
	if polygons, err = loadPolygonsFromFile(geoDataFilePath); err != nil {
		return nil, err
	}

	return &Imitator{
		log:        log,
		devices:    nil,
		TrackCache: trackcache.NewCache(polygons),
	}, nil
}

func (s *Imitator) CreateVirtualDevices(cfg *config.Config, deviceIDs []string) {
	N := len(deviceIDs)

	s.devices = make([]virtualdevice.VirtualDevice, 0, N)

	for i := 0; i < N; i++ {
		dev, err := NewVirtualDevice(deviceIDs[i], cfg)
		if err != nil {
			continue
		}
		s.devices = append(s.devices, dev)
	}

	if s.log != nil {
		s.log.Info().Msgf("%d devices were created for imitation", len(s.devices))
	}
	s.netConfig = cfg.NetworkConfig
}

// NewVirtualDevice - конструктор виртуального устройства
func NewVirtualDevice(deviceID string, cfg *config.Config) (virtualdevice.VirtualDevice, error) {
	if cfg == nil {
		return nil, errors.New("invalid config")
	}
	// создаем инстанс виртуального маяка
	if dev := beacon.NewVirtualBeacon(deviceID, cfg.LaunchDuration, cfg.WakeUpInterval, cfg.Log); dev != nil {
		return dev, nil
	}

	return nil, errors.New("error create new instance")
}

// StartImitation - запуск имитации, операция блокирующая
func (s *Imitator) StartImitation(finish chan struct{}, abort chan struct{}) error {
	if s == nil {
		return errors.New("StartImitation : StartImitation == nil")
	}
	const countRetrySend = 5

	var wg sync.WaitGroup
	wg.Add(1)

	go func(finish chan struct{}, wg *sync.WaitGroup) {
		defer wg.Done()
		N := len(s.devices)
		var waitGroup sync.WaitGroup
		s.DeviceInfo = make([]device.ParamsModel, len(s.devices))
		for {
			select {
			case <-finish:
				return
			case <-abort:
				return
			default:
				var counter uint32
				// каждое устройство выполнит подключение, оправку данных и отключение от сервера
				start := time.Now()
				for i := 0; i < N; i++ {
					if !s.devices[i].IsNeedWakeUp() {
						s.mu.Lock()
						s.DeviceInfo[i] = device.ParamsModel{}
						s.mu.Unlock()

						continue
					}

					key, err := strconv.ParseUint(s.devices[i].GetID(), 10, 64)
					if err != nil {
						continue
					}
					imitParamRecord, isFind := s.TrackCache.Get(key)

					if !isFind {
						continue
					}

					if err = s.devices[i].UpdateGeoParamByExternalImitator(&imitParamRecord); err != nil {
						continue
					}

					s.mu.Lock()
					v := s.devices[i].GetDeviceParameters()
					s.DeviceInfo[i] = v
					s.mu.Unlock()

					counter++
					waitGroup.Add(1)
					go func(dev virtualdevice.VirtualDevice, wg *sync.WaitGroup) {
						defer wg.Done()
						if netConnect, err := net.Dial(
							s.netConfig.Protocol,
							fmt.Sprintf("%s:%d", s.netConfig.Ip, s.netConfig.Port)); err == nil {
							defer netConnect.Close()

							for i := 0; i < countRetrySend; i++ {
								isSend, err := dev.Send(netConnect)

								// повторять отправку в течении минуты с интервалом 10 секунда, пока не отправим
								if err == nil && isSend {
									break
								}
								if err != nil {
									s.log.Error().Msgf("error send message %s/ try send late", err.Error())
								}
								time.Sleep(10 * time.Second)
							}
						} else {
							s.log.Error().Msgf("Can't connect to %s:%d. %s", s.netConfig.Ip, s.netConfig.Port, err.Error())
						}
					}(s.devices[i], &waitGroup)
				}
				waitGroup.Wait()
				s.log.Info().Msgf("Total device count %d, to update %d", len(s.devices), counter)
				time.Sleep(3*time.Second - time.Since(start))
			}
		}
	}(finish, &wg)
	wg.Wait()

	return nil
}

func (s *Imitator) GetDeviceImitInfo() []device.ParamsModel {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.DeviceInfo
}

func loadPolygonsFromFile(path string) (map[uint16]util.GeoPolygon, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	geoData := make([]util.GeoPolygon, 0)

	if err := json.Unmarshal(buf, &geoData); err != nil {
		return nil, err
	}

	mapGeoData := make(map[uint16]util.GeoPolygon)
	for i := 0; i < len(geoData); i++ {
		mapGeoData[geoData[i].RegCode] = geoData[i]
	}

	return mapGeoData, nil
}
