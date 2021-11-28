package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func LoadConfiguration() (Config, error) {
	//path := "E:\\Practice\\Golang\\CloudNative\\HttpServer\\src\\config.json"
	path := "./config.json"
	jsonConfig := &Config{}

	//读取json文件
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panicln(err)
		return *jsonConfig, err
	}

	//解析json
	err = json.Unmarshal(file, jsonConfig)
	if err != nil {
		log.Panicln(err)
		return *jsonConfig, err
	}

	return *jsonConfig, nil
}

type Config struct {
	Port int
}
