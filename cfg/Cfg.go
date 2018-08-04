package cfg

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct { // https://mholt.github.io/json-to-go/
	Path      string  `json:"path"`
	PixelSize float64 `json:"pixelsize"`
	Font      string  `json:"Font"`
	File      string  `json:"File"`
}

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// Init zumn initialisieren der Conifg
func Init() Config {
	found, e := Exists("./config.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	if !found {
		return NewConfig() // wenn es keine config.json gibt dann erstelle sie
	} else {
		return LoadConfigFile() // sonst lese sie
	}
}

// LoadConfigFile läd die Config
func LoadConfigFile() Config {
	file, e := ioutil.ReadFile("./config.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	var jsontype Config
	err := json.Unmarshal(file, &jsontype)
	if err != nil {
		panic(err)
	}
	return jsontype
}

// SaveConfigFile Speichert die Config
func SaveConfigFile(AppPath string, jsonBlob Config) error { // https://github.com/spddl/csgo-reporter/blob/master/Config/Config.go#L147
	bytes, err := json.Marshal(jsonBlob)
	if err != nil {
		return errors.New("Config konnte nicht gespeichert werden (json error)")
	} else {
		err = ioutil.WriteFile(filepath.Join(AppPath, "config.json"), bytes, 0644)
		if err == nil {
			return nil
		} else {
			return err
		}
	}
}

// NewConfig erstelle eine Config
func NewConfig() Config {
	jsonBlob := json.RawMessage(`{
		"path": "",
		"pixelsize": 1.2,
		"font": ""
	}`)

	var jsontype Config
	err := json.Unmarshal(jsonBlob, &jsontype)
	if err != nil {
		panic(err)
	}

	return jsontype
}
