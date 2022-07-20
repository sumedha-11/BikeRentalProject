package config

import (
	sql "bikerentalProject/pkg/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var Config *AppConfig

type AppConfig struct {
	Common   *CommonConfig `json:"CommonConfig"`
	DBConfig *sql.DBConfig `json:"DBConfig"`
}

type CommonConfig struct {
	ServerPort string `json:"serverPort"`
}

func ReadConfig(configFile string) {
	Config = &AppConfig{}
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Printf("File error: %v\n", err)

		os.Exit(1)
	}
	if err := json.Unmarshal(file, &Config); err != nil {
		log.Fatalf("unable to marshal config data")
		return
	}
	fmt.Println("config loaded ", Config)
}
