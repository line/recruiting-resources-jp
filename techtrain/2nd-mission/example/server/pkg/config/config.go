// Package config reads config file and initializes config variable
package config

import (
	"io/ioutil"
	"path/filepath"

	"github.com/go-openapi/loads"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Port     string `yaml:"port"`
		Loglevel string `yaml:"loglevel"`
	}
	Mysql struct {
		Endpoint string `yaml:"endpoint"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Db       string `yaml:"db"`
	}
	LineAPI    string `yaml:"lineAPI"`
	SchemaFile string `yaml:"schemaFile"`
	Schema     *loads.Document
}

var (
	conf     *Config
	FilePath string = "server/configs/config.yml"
)

func GetConf() Config {
	if conf == nil {
		FilePath, _ := filepath.Abs(FilePath)
		yamlFile, err := ioutil.ReadFile(FilePath)
		if err != nil {
			log.Fatal(err)
		}
		tmpConf := Config{}
		err = yaml.Unmarshal(yamlFile, &tmpConf)
		if err != nil {
			log.Fatalf("yaml parse fail: %v", err)
		}
		conf = &tmpConf
		if conf.Schema == nil {
			conf.Schema, err = loads.Spec(conf.SchemaFile)
			if err != nil {
				log.WithFields(log.Fields{"error": err.Error()}).
					Fatal("Failed to load API schema file!")
			}
		}
		log.WithFields(log.Fields{
			"serverConfig": conf.Server,
			"mysql":        conf.Mysql.Endpoint,
		}).Info("Load configuration")
	}
	return *conf
}
