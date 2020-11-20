package utils

import (
	"cybervein.org/CyberveinDB/logger"
	"github.com/creasty/defaults"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"time"
)

type Configuration struct {
	App        AppConfig        `yaml:"app-config"`
	Rpc        RpcConfig        `yaml:"rpc-config"`
	Redis      RedisConfig      `yaml:"redis-config"`
	Tendermint TendermintConfig `yaml:"tendermint-config"`
	HttpServer HttpServerConfig `yaml:"http-server-config"`
}

type RedisConfig struct {
	Url      string `yaml:"url"`
	Db       int    `yaml:"db"`
	Password string `yaml:"password"`
	RDBPath  string `yaml:"rdb_path"`
}

type TendermintConfig struct {
	Url string `yaml:"url"`
}

type HttpServerConfig struct {
	RunMode      string        `yaml:"run_mode"`
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

type RpcConfig struct {
	Port int `yaml:"port"`
}

type AppConfig struct {
	Name          string `yaml:"name"`
	DbPassword    string `yaml:"db_password"`
	Plugin        string `yaml:"plugin" default:"default"`
}

var Config Configuration

func InitConfig() {
	absPath, err := filepath.Abs("../conf/configuration.yaml")
	if err != nil {
	}
	yamlFile, err := ioutil.ReadFile(absPath)
	if err != nil {
		logger.Log.Error(err)
		return
	}
	yaml.Unmarshal(yamlFile, &Config)
	setDefaults(&Config)
}

func setDefaults(v interface{}) {
	if err := defaults.Set(v); err != nil {
		panic(err)
	}
}
