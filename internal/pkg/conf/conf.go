package conf

import (
	"flag"
	"fmt"

	"github.com/spf13/viper"
)

var Conf Config

type Config struct {
	Base       Base       `yaml:"base"`
	Http       Http       `yaml:"http"`
	GRPC       GRPC       `yaml:"grpc"`
	MySQL      MySQL      `yaml:"mysql"`
	JWT        JWT        `yaml:"jwt"`
	MongoDB    MongoDB    `yaml:"mongodb"`
	Prometheus Prometheus `yaml:"prometheus"`
	Kafka      Kafka      `yaml:"kafka"`
	ES         ES         `yaml:"es"`
	Clients    Clients    `yaml:"clients"`
	Proof      Proof      `yaml:"proof"`
}

type Base struct {
	IsDebug bool `mapstructure:"is_debug"`
}

type Http struct {
	Port int `yaml:"port"`
}

type GRPC struct {
	Port int `yaml:"port"`
}

type MySQL struct {
	Username string `yaml:"username"`
	Passport string `yaml:"passport"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DBName   string `yaml:"dbname"`
	Charset  string `yaml:"charset"`
}

type JWT struct {
	Issuer string `yaml:"issuer"`
	Secret string `yaml:"secret"`
}

type MongoDB struct {
	DSN      string `yaml:"dsn"`
	Database string `yaml:"database"`
	PoolSize uint64 `mapstructure:"pool_size"`
}

type Prometheus struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Kafka struct {
	Host  string `yaml:"host"`
	Topic string `yaml:"topic"`
}

type ES struct {
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Clients struct {
	Runner Runner
}

type Runner struct {
	RunAPI   string `mapstructure:"run_api"`
	RunScene string `mapstructure:"run_scene"`
	RunPlan  string `mapstructure:"run_plan"`
}

type Proof struct {
	InfoLog string `mapstructure:"info_log"`
	ErrLog  string `mapstructure:"err_log"`
}

func MustInitConf() {
	var configFile string
	flag.StringVar(&configFile, "c", "./configs/dev.yaml", "app config file.")
	if !flag.Parsed() {
		flag.Parse()
	}

	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if err := viper.Unmarshal(&Conf); err != nil {
		panic(fmt.Errorf("unmarshal error config file: %w", err))
	}

	fmt.Println("config initialized")
}
