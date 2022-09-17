package config

import (
	"errors"
	"flag"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/art-injener/iot-platform/pkg/logger"
	"github.com/art-injener/iot-platform/util"
)

const DebugLevel = "debug"

const (
	defaultDBScanPeriod      = 100 * time.Millisecond
	defaultMonitoringPeriod  = 3 * time.Second
	defaultTimeoutReadPeriod = 3 * time.Second
)

type Config struct {
	// Log level variant : debug, release
	LogLevel string `mapstructure:"LOG_LEVEL"`
	// Base phone number for generating devices
	BasePhone string `mapstructure:"BASE_PHONE"`
	// Number of devices to generate
	CountDevices uint64 `mapstructure:"COUNT_DEVICES"`
	// Period (in ms) of sending data to the server
	MonitoringPeriod int64 `mapstructure:"MONITORING_PERIOD"`
	// Time in minute launching
	LaunchDuration int `mapstructure:"LAUNCH_DURATION"`
	// Time in minute wakeup
	WakeUpInterval int `mapstructure:"WAKE_UP_INTERVAL"`
	// Web server port
	WebServerPort string `mapstructure:"WEB_SERVER_PORT"`
	*NetworkConfig
	*DBConfig
	*RabbitConfig
	Phones []string
	Log    *logger.Logger
}

type NetworkConfig struct {
	// Type of network protocols, values variant : tcp, udp
	Protocol string `mapstructure:"NETWORK_PROTOCOL"`
	// Server IP address
	Ip string `mapstructure:"NETWORK_IP"`
	// Server port
	Port uint16 `mapstructure:"NETWORK_PORT"`
	// Timeout wait read from server
	ReadTimeoutMs int64 `mapstructure:"NETWORK_TIMEOUT_READ"`
}

type DBConfig struct {
	Host         string `mapstructure:"DEVICES_DB_HOST"`
	Port         uint16 `mapstructure:"DEVICES_DB_PORT"`
	NameDB       string `mapstructure:"DEVICES_DB_DATABASE"`
	User         string `mapstructure:"DEVICES_DB_USERNAME"`
	Password     string `mapstructure:"DEVICES_DB_PASSWORD"`
	ExecTimeout  int    `mapstructure:"DEVICES_DB_EXEC_TIMEOUT"`
	DBScanPeriod uint64 `mapstructure:"DEVICES_DB_SCAN_PERIOD"`
}

type SSHConfig struct {
	Host     string // SSH Server Hostname/IP
	Port     uint16 // SSH Port
	User     string // SSH Username
	Password string // Empty string for no password
}

type RabbitConfig struct {
	Url      string `mapstructure:"RABBIT_URL"`
	Queue    Queue
	Qos      Qos
	Consumer Consumer
}

type Queue struct {
	QueueName  string `mapstructure:"RABBIT_QUEUE_NAME"`
	Durable    bool   `mapstructure:"RABBIT_QUEUE_DURABLE"`
	AutoDelete bool   `mapstructure:"RABBIT_QUEUE_AUTO_DELETE"`
	Exclusive  bool   `mapstructure:"RABBIT_QUEUE_EXCLUSIVE"`
	NoWait     bool   `mapstructure:"RABBIT_QUEUE_NO_WAIT"`
}

type Qos struct {
	PrefetchCount int  `mapstructure:"RABBIT_QOS_PREFETCH_COUNT"`
	PrefetchSize  int  `mapstructure:"RABBIT_QOS_PREFETCH_SIZE"`
	Global        bool `mapstructure:"RABBIT_QOS_GLOBAL"`
}

type Consumer struct {
	Tag       string `mapstructure:"RABBIT_CONSUMER_TAG"`
	AutoAck   bool   `mapstructure:"RABBIT_CONSUMER_AUTO_ACK"`
	Exclusive bool   `mapstructure:"RABBIT_CONSUMER_EXCLUSIVE"`
	NoLocal   bool   `mapstructure:"RABBIT_CONSUMER_NO_LOCAL"`
	NoWait    bool   `mapstructure:"RABBIT_CONSUMER_NO_WAIT"`
}

func GetConfig(path string) (*Config, error) {
	cfg := Config{}
	cfg.NetworkConfig = &NetworkConfig{}
	cfg.DBConfig = &DBConfig{}
	cfg.RabbitConfig = &RabbitConfig{}

	err := readEnvConfig(path, &cfg)
	if err != nil {
		return nil, err
	}

	cfg.NetworkConfig.Protocol = strings.ToLower(cfg.NetworkConfig.Protocol)
	if cfg.NetworkConfig.Protocol != "tcp" && cfg.NetworkConfig.Protocol != "udp" {
		return nil, errors.New("wrong network protocol passed")
	}

	if cfg.MonitoringPeriod == 0 {
		cfg.MonitoringPeriod = int64(defaultMonitoringPeriod)
	}

	if cfg.NetworkConfig.ReadTimeoutMs == 0 {
		cfg.NetworkConfig.ReadTimeoutMs = int64(defaultTimeoutReadPeriod)
	}

	if cfg.DBConfig != nil && cfg.DBConfig.DBScanPeriod == 0 {
		cfg.DBConfig.DBScanPeriod = uint64(defaultDBScanPeriod)
	}

	if cfg.BasePhone == "" || cfg.CountDevices == 0 {
		return nil, errors.New("empty ids and iddev/countdev values")
	}

	cfg.Phones = util.GenerateIDs(cfg.BasePhone, cfg.CountDevices)

	flag.StringVar(&cfg.LogLevel, "loglevel", "release", "Log level variant : debug, release")

	flag.StringVar(&cfg.BasePhone, "startid", "89991234567", "base phone number for generating devices ")

	return &cfg, nil
}

func readEnvConfig(path string, cfg *Config) error {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return err
	}

	if err := viper.Unmarshal(cfg.DBConfig); err != nil {
		return err
	}

	if err := viper.Unmarshal(cfg.NetworkConfig); err != nil {
		return err
	}

	if err := viper.Unmarshal(cfg.RabbitConfig); err != nil {
		return err
	}

	return nil
}
