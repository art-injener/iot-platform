package device

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	defaultBalanceSize            float32 = 20
	defaultBalanceSubtractionSize float32 = 0.05
	tempLimit                     int     = 50
	BatteryFullChargeValue        float32 = 5.95
	sessionsNumber                int     = 1200
)

type ParamsModel struct {
	ID              string    `json:"id"`
	SystemTime      int64     `json:"system_time"`      // Местное время - параметр STIME
	BalanceSim      float32   `json:"balance_sim"`      // баланс - параметр BAL
	Temperature     float32   `json:"temperature"`      // температура  - параметр TE
	BatteryCharge   float32   `json:"battery_charge"`   // состояние батареи - параметр  VB=100%(5.95V)
	IntervalCounter uint64    `json:"interval_counter"` // номер интервального извещения - параметр IC. Увеличить на 1 при каждой успешной передаче.
	Latitude        float64   `json:"latitude"`         // широта, параметр LA
	LAD             string    `json:"lad"`              // полушарие, параметр LAD
	Longitude       float64   `json:"longitude"`        // долгота, параметр LO
	LOD             string    `json:"lod"`              // полушарие - параметр LOD
	SignalQuality   int       `json:"signal_quality"`   // параметр мощность сигнала
	Speed           int       `json:"speed"`            // параметр Speed=0 - скорость движения
	DirectionMove   float32   `json:"direction_move"`   // параметр DirectionMove - направление движения
	SeanceTime      int64     `json:"seance_time"`      // время следующего выхода на связь
	CreatedAt       time.Time `json:"-"`
	UpdatedAt       time.Time `json:"-"`
}

// NewDeviceParams - конструктор изменяемых параметров
func NewDeviceParams(deviceID string) *ParamsModel {
	return &ParamsModel{
		ID:              deviceID,
		SystemTime:      time.Now().Local().Unix(),
		BalanceSim:      20.0,
		Temperature:     25,
		BatteryCharge:   BatteryFullChargeValue,
		IntervalCounter: 1,
		Latitude:        55.758436,
		LAD:             "N",
		Longitude:       37.5510486,
		LOD:             "E",
		SignalQuality:   69,
		Speed:           0,
		DirectionMove:   0,
		SeanceTime:      time.Now().Local().Add(20 * time.Minute).Unix(),
	}
}

func (s *ParamsModel) UpdateTemperature() {
	if s == nil {
		return
	}

	if (s.IntervalCounter/uint64(tempLimit*2))%2 == 0 {
		if (s.IntervalCounter/uint64(tempLimit))%2 == 0 {
			s.Temperature++
		} else {
			s.Temperature--
		}
		return
	}

	if (s.IntervalCounter/uint64(tempLimit))%2 == 0 {
		s.Temperature--
	} else {
		s.Temperature++
	}
}

func (s *ParamsModel) RandomUpdateTemperature() {
	v := ((rand.Float32() * 5) + 5) / 10

	if time.Now().UnixNano()%2 == 0 {
		s.Temperature += v
		return
	}
	s.Temperature -= v
}

func (s *ParamsModel) UpdateBalance() {
	if s == nil {
		return
	}

	if s.BalanceSim > 0 {
		s.BalanceSim -= defaultBalanceSubtractionSize
	} else {
		s.BalanceSim = defaultBalanceSize
	}

}

func (s *ParamsModel) UpdateSignalQuality() {
	if s == nil {
		return
	}

	s.SignalQuality += rand.Intn(20)
	s.SignalQuality %= 100
}

func (s *ParamsModel) UpdateBatteryCharge() {
	if s == nil {
		return
	}

	s.BatteryCharge -= BatteryFullChargeValue / float32(sessionsNumber)
	if s.BatteryCharge < 0.01 {
		s.BatteryCharge = BatteryFullChargeValue
	}
}

// BatteryChargeToString - сериализация уровня заряда : "100%(5.95V)".
func (s *ParamsModel) BatteryChargeToString() string {
	if s == nil {
		return ""
	}

	per := 100.0 * s.BatteryCharge / BatteryFullChargeValue

	return fmt.Sprintf("%d%%(%.2fV)", int(per), s.BatteryCharge)
}

func (s *ParamsModel) UpdateAll() {
	if s == nil {
		return
	}
	s.UpdateTemperature()
	s.UpdateBalance()
	s.UpdateSignalQuality()
	s.UpdateBatteryCharge()
}
