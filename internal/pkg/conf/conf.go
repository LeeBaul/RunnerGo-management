package conf

import (
	"flag"
	"fmt"

	"github.com/spf13/viper"
)

var Conf Config

type Config struct {
	Http  Http  `yaml:"http"`
	GRPC  GRPC  `yaml:"grpc"`
	MySQL MySQL `yaml:"mysql"`
	JWT   JWT   `yaml:"jwt"`
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
