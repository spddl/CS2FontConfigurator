package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
)

type Config struct { // https://mholt.github.io/json-to-go/
	Path      string  `json:"path"`
	Font      string  `json:"font"`
	Pixelsize float64 `json:"pixelsize"`
	TestCase  string  `json:"TestCase"`
}

func (c *Config) Init() {
	found, e := FileExists("./config.json")
	if e != nil {
		log.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	if !found {
		c.newConfig()
	} else {
		c.LoadConfigFile()
	}
}

func (c *Config) LoadConfigFile() {
	file, e := os.ReadFile("./config.json")
	if e != nil {
		log.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	err := json.Unmarshal(file, &c)
	if err != nil {
		log.Println(err)
		c.newConfig()
	}
}

func (c *Config) SaveConfigFile(apppath string) error {
	bytes, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return errors.New("config.json could not be saved (JSON error)")
	} else {
		err = os.WriteFile(filepath.Join(apppath, "config.json"), bytes, 0o644)
		if err == nil {
			return nil
		} else {
			return err
		}
	}
}

// newConfig create a Config
func (c *Config) newConfig() {
	jsonBlob := json.RawMessage(`{
		"path": "",
		"font": "",
		"pixelsize": 1.0,
		"TestCase": "ABC 100 PLAY DUST II 98 in 4"
	}`)

	if err := json.Unmarshal(jsonBlob, c); err != nil {
		panic(err)
	}
}
