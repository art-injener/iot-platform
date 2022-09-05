package config

import (
	"flag"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/art-injener/iot-platform/pkg/logger"
	util "github.com/art-injener/iot-platform/util/helper"
)

const DebugLevel = "debug"

const (
	defaultDBScanPeriod      = 100 * time.Millisecond
	defaultMonitoringPeriod  = "3000"
	defaultTimeoutReadPeriod = "3000"
)

type Config struct {
	LogLevel         string
	Mode             string
	MonitoringPeriod time.Duration
	LaunchDuration   int
	WakeUpInterval   int
	Phones           []string
	Net              *NetworkConfig
	DB               *DBConfig
	Log              *logger.Logger
}

type NetworkConfig struct {
	Protocol      string
	Ip            string
	Port          uint16
	ReadTimeoutMs time.Duration
}

type DBConfig struct {
	Host         string `mapstructure:"DEVICES_DB_HOST"`
	Port         uint16 `mapstructure:"DEVICES_DB_PORT"`
	NameDB       string `mapstructure:"DEVICES_DB_DATABASE"`
	User         string `mapstructure:"DEVICES_DB_USERNAME"`
	Password     string `mapstructure:"DEVICES_DB_PASSWORD"`
	DBScanPeriod uint64
}
type SSHConfig struct {
	Host     string // SSH Server Hostname/IP
	Port     uint16 // SSH Port
	User     string // SSH Username
	Password string // Empty string for no password
}

func GetConfig(path string) (*Config, error) {
	cfg := Config{}
	cfg.Net = &NetworkConfig{}
	cfg.DB = &DBConfig{}

	// data base settings
	var portDB uint
	var basePhone string
	var countDev uint
	var period, timeout string
	var err error

	readEnvConfig(path, &cfg)

	// основные настройки
	flag.StringVar(&cfg.LogLevel, "loglevel", "release", "Log level variant : debug, release")

	flag.StringVar(&basePhone, "startid", "89991234567", "base phone number for generating devices ")
	flag.UintVar(&countDev, "countid", 10000, "number of devices to generate")
	flag.IntVar(&cfg.LaunchDuration, "launch", 20, "time in minute launching ")
	flag.IntVar(&cfg.WakeUpInterval, "wakeup", 20, "time in minute wakeup ")

	// настройки подключения к серверу
	var port uint
	// TODO : добавить валидацию переданных данных
	flag.StringVar(&cfg.Net.Protocol,
		"protocol",
		"tcp",
		"Type of network protocols, values variant : tcp, udp")

	cfg.Net.Protocol = strings.ToLower(cfg.Net.Protocol)
	if cfg.Net.Protocol != "tcp" && cfg.Net.Protocol != "udp" {
		return nil, errors.New("wrong network protocol passed")
	}
	flag.StringVar(&cfg.Net.Ip, "ip", "127.0.0.1", "Server IP address ")
	flag.UintVar(&port, "port", 9000, "Server port")
	flag.StringVar(&timeout, "timeout", defaultTimeoutReadPeriod, "timeout waite read from server")

	flag.StringVar(&period, "period", defaultMonitoringPeriod, "Period (in ms) of sending data to the server ")

	if cfg.MonitoringPeriod, err = time.ParseDuration(period + "ms"); err != nil {
		cfg.MonitoringPeriod = 5 * time.Second
	}

	if cfg.Net.ReadTimeoutMs, err = time.ParseDuration(timeout + "ms"); err != nil {
		cfg.Net.ReadTimeoutMs = 3 * time.Second
	}

	if cfg.DB != nil {
		flag.Uint64Var(&cfg.DB.DBScanPeriod, "dbperiod", uint64(defaultDBScanPeriod), "database scan period")
	}

	flag.Parse()

	if basePhone == "" {
		return nil, errors.New("Empty ids and iddev/countdev values")
	}

	cfg.Phones = util.GenerateIDs(basePhone, uint64(countDev))

	cfg.Net.Port = uint16(port & 0xffff)
	cfg.DB.Port = uint16(portDB & 0xffff)

	return &cfg, nil
}

func readEnvConfig(path string, cfg *Config) error {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	viper.ReadInConfig()

	cfgDb := DBConfig{}

	if err := viper.Unmarshal(&cfgDb); err != nil {
		return err
	}

	cfg.DB = &cfgDb

	return nil
}
