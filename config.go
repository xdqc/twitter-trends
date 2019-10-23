package tweet

import (
	"io/ioutil"
	"path/filepath"

	yaml "gopkg.in/yaml.v1"
)

type Config struct {
	APIKey       string `yaml:"API_Key"`
	APISecret    string `yaml:"API_Secret"`
	AccessKey    string `yaml:"Access_Key"`
	AccessSecret string `yaml:"Access_Secret"`
}

var (
	cfg         *Config
	cfgFilename = "./config.yml"
)

//SetConfigFile path
func SetConfigFile(path string) {
	cfgFilename = path
}

//GetConfig parse config
func GetConfig() *Config {
	if cfg != nil {
		return cfg
	}
	filename, _ := filepath.Abs(cfgFilename)
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}
	var c *Config
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		panic(err)
	}
	cfg = c
	return cfg
}
