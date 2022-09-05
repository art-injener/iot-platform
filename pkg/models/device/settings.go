package device

// SettingsModel - параметры устройства, значения которых константны или меняются только по команде из панели
type SettingsModel struct {
	ID      string `json:"id"`      // Номер устройства - параметр ID
	Version string `json:"version"` // Версия устройства - параметр VER
	TZ      int    `json:"tz"`      // таймзона - параметр TZ. 220 - Москва. [-14, +14] - смещение от UTC
	WUI     int    `json:"wui"`     // параметр WUI - Временной интервал выхода на связь [10 - 10080]. Время в минутах.
	GPST    int    `json:"gpst"`    // параметр GPST=3 - Максимальное время поиска спутников в минутах
}

// NewDeviceSettings staticParameters - конструктор статичных параметров
func NewDeviceSettings(deviceID string, wakeUpInterval int) *SettingsModel {
	return &SettingsModel{
		ID:      deviceID,
		Version: "0.0.1",
		TZ:      220,
		WUI:     wakeUpInterval,
		GPST:    3,
	}
}
