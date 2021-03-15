package main

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	HTTP httpConf
	TCP  tcpConf
	DATA dataConf
}

type httpConf struct {
	Addr         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	MaxHeaderMIB int
}

type tcpConf struct {
	Addr        string
	Port        string
	BufferLimit int
	AwaitConn   time.Duration
}

type dataConf struct {
	availableFields map[string]string
	message         string
	timeDuration    int64
}

func setConfigs(filepath string, cfg *Config) error {
	err := parseConfigFile(filepath)
	if err != nil {
		return err
	}

	cfg.HTTP.Addr = viper.GetString("http.addr")
	cfg.HTTP.Port = viper.GetString("http.port")
	cfg.HTTP.ReadTimeout = viper.GetDuration("http.readTimeout")
	cfg.HTTP.WriteTimeout = viper.GetDuration("http.writeTimeout")
	cfg.HTTP.MaxHeaderMIB = viper.GetInt("http.maxHeaderMIB")

	cfg.TCP.Addr = viper.GetString("tcp.addr")
	cfg.TCP.Port = viper.GetString("tcp.port")
	cfg.TCP.AwaitConn = viper.GetDuration("tcp.awaitConn")
	cfg.TCP.BufferLimit = viper.GetInt("tcp.bufLimitMIB")

	cfg.DATA.availableFields = viper.GetStringMapString("data.availableFields")
	cfg.DATA.message = viper.GetString("data.message")
	cfg.DATA.timeDuration = viper.GetInt64("data.timeDuration")
	return nil
}

func parseConfigFile(filepath string) error {
	path := strings.Split(filepath, "/")

	viper.AddConfigPath(path[0]) // folder
	viper.SetConfigName(path[1]) // config file name
	return viper.ReadInConfig()
}
