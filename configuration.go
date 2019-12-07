package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const (
	configurationFilename = "settings.json"
)

// Описание конфигурации
type configuration struct {
	Login    string
	Password string
	Year     string
}

var (
	// Configuration - Конфигурация программы
	Configuration configuration
)

func loadConfiguration() {
	// Инициализация
	Configuration.Login = "some@university.ru"
	Configuration.Password = "PASSWORD"
	Configuration.Year = "2019"
	if fileExist(configurationFilename) {
		data, err := ioutil.ReadFile(configurationFilename)
		if err != nil {
			return
		}
		var jsonData map[string]interface{}
		if err := json.Unmarshal(data, &jsonData); err != nil {
			return
		}
		for key := range jsonData {
			switch key {
			case "login":
				Configuration.Login = jsonData[key].(string)
			case "password":
				Configuration.Password = jsonData[key].(string)
			case "year":
				Configuration.Year = jsonData[key].(string)
			}
		}
	}
}

// Сохранение конфигурации
func saveConfiguration() {
	jsonData := []byte(fmt.Sprintf("{\"login\":\"%s\", \"password\":\"%s\", \"year\":\"%s\"}", Configuration.Login, Configuration.Password, Configuration.Year))
	ioutil.WriteFile(configurationFilename, jsonData, 0644)
}
