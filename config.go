package main

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	HTTP     httpConf
	TCP      tcpConf
	DATA     dataConf
	SECURITY securityConf
}

type httpConf struct {
	Addr         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	MaxHeaderMIB int
	RefreshURL   string
}

type tcpConf struct {
	Addr        string
	Port        string
	BufferLimit int
	AwaitConn   time.Duration
}

type dataConf struct {
	directionTravel            map[string][]string
	tagPlatform                string
	multipleViolationField     string
	bindMultipleViolationField string
	availableFields            map[string]string
	violationTypes             []string
	violationValue             map[string][]string
	violationName              map[string]string
	trackThrustes              []string
	numberTC                   []string
	message                    string
	timeDuration               int64
}

type securityConf struct {
	deadline time.Time
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
	cfg.HTTP.RefreshURL = viper.GetString("http.refreshURL")

	cfg.TCP.Addr = viper.GetString("tcp.addr")
	cfg.TCP.Port = viper.GetString("tcp.port")
	cfg.TCP.AwaitConn = viper.GetDuration("tcp.awaitConn")
	cfg.TCP.BufferLimit = viper.GetInt("tcp.bufLimitMIB")

	cfg.DATA.tagPlatform = viper.GetString("data.tagPlatform")
	cfg.DATA.multipleViolationField = viper.GetString("data.multipleViolationField")
	cfg.DATA.bindMultipleViolationField = viper.GetString("data.bindMultipleViolationField")
	cfg.DATA.violationTypes = viper.GetStringSlice("data.violationTypes")
	cfg.DATA.violationValue = viper.GetStringMapStringSlice("data.violationValue")
	cfg.DATA.violationName = viper.GetStringMapString("data.violationName")
	cfg.DATA.directionTravel = viper.GetStringMapStringSlice("data.directionTravel")
	cfg.DATA.availableFields = viper.GetStringMapString("data.availableFields")
	cfg.DATA.trackThrustes = viper.GetStringSlice("data.trackThrustes")
	cfg.DATA.numberTC = viper.GetStringSlice("data.numberTC")
	cfg.DATA.message = viper.GetString("data.message")
	cfg.DATA.timeDuration = viper.GetInt64("data.timeDuration")

	// Setting deadline
	cfg.SECURITY.deadline, _ = time.Parse("2006-Jan-02", "2021-May-15")
	return nil
}

func parseConfigFile(filepath string) error {
	path := strings.Split(filepath, "/")

	viper.AddConfigPath(path[0]) // folder
	viper.SetConfigName(path[1]) // config file name
	return viper.ReadInConfig()
}
