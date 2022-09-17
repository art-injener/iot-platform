package util

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"

	geo "github.com/kellydunn/golang-geo"
	"github.com/pkg/errors"
)

// GeoPoint - структура гео-точки
type GeoPoint struct {
	Lat float64
	Lon float64
}

// Value - значение в которое должен возвращать драйвер БД
func (gp GeoPoint) Value() (driver.Value, error) {
	return fmt.Sprintf("(%f,%f)", gp.Lat, gp.Lon), nil
}

// Latitude - широта
func (gp GeoPoint) Latitude() float64 {
	return gp.Lat
}

// Longitude - долгота
func (gp GeoPoint) Longitude() float64 {
	return gp.Lon
}

// Scan - метод сканирования. Подходит для pgx, pgxpool, pq
// C ORM не проверялся
func (gp *GeoPoint) Scan(src interface{}) error {
	raw, ok := src.(string)
	if !ok {
		return errors.New("can't switch type")
	}
	start := strings.Index(raw, "(")
	end := strings.Index(raw, ")")
	if start < 0 || end < 0 {
		return errors.New("Index == -1 ")
	}

	coords := raw[start+1 : end]
	values := strings.Split(coords, " ")

	if len(values) < 2 {
		return errors.New("Too few values")
	}
	var err error

	if gp.Lat, err = strconv.ParseFloat(values[0], 64); err != nil {
		return errors.Wrap(err, "err when parse latitude")
	}
	if gp.Lon, err = strconv.ParseFloat(values[1], 64); err != nil {
		return errors.Wrap(err, "err when parse longitude")
	}
	return nil
}

// GeoPoint - структура гео-точки
type GeoPolygon struct {
	RegCode uint16
	Title   string
	Min     geo.Point
	Max     geo.Point
	Center  geo.Point
	Radius  float32
}

// Scan - метод сканирования. Подходит для pgx, pgxpool, pq
// C ORM не проверялся
func (gp *GeoPolygon) Scan(src interface{}) error {

	raw, ok := src.(string)
	if !ok {
		return errors.New("can't switch type")
	}
	start := strings.Index(raw, "(")
	end := strings.Index(raw, ")")
	if start < 0 || end < 0 {
		return errors.New("Index == -1 ")
	}

	coords := raw[start+1 : end]
	values := strings.Split(coords, ",")

	for i := 0; i < len(values); i++ {
		if strings.Contains(values[0], "(") {
			values[0] = values[0][1:]
		}
		dots := strings.Split(values[i], " ")

		if len(dots) < 2 {
			return errors.New("Too few values")
		}
		var err error
		var lat, lon float64

		if lat, err = strconv.ParseFloat(dots[1], 64); err != nil {
			return errors.Wrap(err, "err when parse latitude")
		}
		if lon, err = strconv.ParseFloat(dots[0], 64); err != nil {
			return errors.Wrap(err, "err when parse longitude")
		}

		p := geo.NewPoint(lat, lon)
		if i == 0 {
			gp.Min = *p
			gp.Max = *p
			continue
		}

		if gp.Min.Lat() > p.Lat() {
			gp.Min = *p
		}
		if gp.Max.Lat() < p.Lat() {
			gp.Max = *p
		}

		gp.Center = *gp.Min.MidpointTo(&gp.Max)
		gp.Radius = float32(gp.Min.GreatCircleDistance(&gp.Max) / 2)
	}
	return nil
}
