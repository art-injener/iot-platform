package trackcache

import (
	"math/rand"
	"sync"
	"time"

	geo "github.com/kellydunn/golang-geo"

	"github.com/art-injener/iot-platform/internal/models/imitator"
	util "github.com/art-injener/iot-platform/util/helper"
)

type TypeImitTrack uint8

type DeviceID uint64
type ImitParamIndex uint32

// Cache struct cache
type Cache struct {
	mu               sync.RWMutex
	imitationTrack   []*imitator.GeoParamsModel
	imitationPoligon map[uint16]util.GeoPolygon
	devices          map[DeviceID]ImitParamIndex // device_hash - index in track
	changeInterval   time.Duration
	abort            chan struct{}
}

// NewCache  - конструктор кэша
func NewCache(poligons map[uint16]util.GeoPolygon) *Cache {
	// cache item
	cache := Cache{
		mu:               sync.RWMutex{},
		imitationTrack:   make([]*imitator.GeoParamsModel, 5),
		devices:          make(map[DeviceID]ImitParamIndex, 1000),
		changeInterval:   1 * time.Second, // Должно быть 1 минута
		imitationPoligon: poligons,
		abort:            make(chan struct{}),
	}

	return &cache
}

// Get - получение информации.
func (c *Cache) Get(key uint64) (imitator.GeoParamsModel, bool) {
	c.mu.Lock()

	defer c.mu.Unlock()

	var value ImitParamIndex
	var ok bool
	var p = imitator.GeoParamsModel{}
	rand.Seed(time.Now().UnixNano())

	randZone := uint16(rand.Intn(99)) //nolint:gosec
	if randZone == 0 {
		randZone++
	}
	centralDot, ok := c.imitationPoligon[randZone]
	if !ok {
		return imitator.GeoParamsModel{}, false
	}
	radius := int(centralDot.Radius)
	if radius <= 0 {
		radius = 20
	}
	if value, ok = c.devices[DeviceID(key)]; !ok {
		p = imitator.GeoParamsModel{}

		p.AverageSpeed = rand.Intn(100)
		p.Azimuth = float32(rand.Intn(360))

		newPoint := centralDot.Center.PointAtDistanceAndBearing(float64(rand.Intn(radius)), //nolint:gosec
			float64(p.Azimuth))
		p.LatitudeBeacon = newPoint.Lat()
		p.LongitudeBeacon = newPoint.Lng()

		c.imitationTrack = append(c.imitationTrack, &p)

		value = ImitParamIndex(len(c.imitationTrack)) - 1
		c.devices[DeviceID(key)] = value

		return p, true
	}

	item := c.imitationTrack[value]
	// cache not found
	if item == nil {
		return imitator.GeoParamsModel{}, false
	}

	point := geo.NewPoint(item.LatitudeBeacon, item.LongitudeBeacon)
	newPoint := point.PointAtDistanceAndBearing(float64(rand.Intn(10)), //nolint:gosec
		float64(item.Azimuth))

	if newPoint.GreatCircleDistance(&centralDot.Center) > float64(centralDot.Radius) {
		newPoint = centralDot.Center.PointAtDistanceAndBearing(float64(rand.Intn(radius)), //nolint:gosec
			float64(p.Azimuth))
	}
	item.LatitudeBeacon = newPoint.Lat()
	item.LongitudeBeacon = newPoint.Lng()

	item.AverageSpeed = rand.Intn(100)
	item.Azimuth += 0.5
	if item.Azimuth > 360 {
		item.Azimuth = 0
	}

	return *item, true
}

func (c *Cache) GetAllRegion() (items map[uint16]util.GeoPolygon) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	items = make(map[uint16]util.GeoPolygon)

	for k, v := range c.imitationPoligon {
		items[k] = v
	}

	return items
}
