package config

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type Logging struct {
	LogFile string `toml:"log_file"`
}

type DBConfig struct {
	Driver string `toml:"driver"`
	Host   string `toml:"host"`
	Port   int    `toml:"port"`
	Name   string `toml:"name"`
	ETC    string `toml:"etc"`
}

type Web struct {
	Port int `toml:"port"`
}

type Config struct {
	Log     Logging  `toml:"logging"`
	DB      DBConfig `toml:"db"`
	Web     Web      `toml:"web"`
}

const fileName = "config.toml"

var Cfg Config

func init() {
	log.Println("Reading config file...")
	_, err := toml.DecodeFile(fileName, &Cfg)
	if err != nil {
		log.Printf("Failed to read file: %v", err)
		os.Exit(1)
	}
}
