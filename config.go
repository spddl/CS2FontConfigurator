package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Path                     string `json:"path"`
	EmbeddedBitmap           bool   `json:"embeddedbitmap"`
	PreferOutline            bool   `json:"preferoutline"`
	DoSubstitutions          bool   `json:"dosubstitutions"`
	BitmapMonospace          bool   `json:"bitmapmonospace"`
	ForceAutohint            bool   `json:"forceautohint"`
	QtUseSubpixelPositioning bool   `json:"qtusesubpixelpositioning"`
	Dpi                      int    `json:"dpi"`

	Font      string  `json:"font"`
	Pixelsize float64 `json:"pixelsize"`
	Weight    int     `json:"weight"`

	Customize bool         `json:"customize"`
	Fonts     []FontStruct `json:"fonts"`
}

type FontStruct struct {
	ValveFont                string  `json:"valvefont"`
	ReplaceFont              string  `json:"replacefont"`
	Pixelsize                float64 `json:"pixelsize"`
	Weight                   int     `json:"weight"`
	EmbeddedBitmap           bool    `json:"embeddedbitmap"`
	PreferOutline            bool    `json:"prefer_outline"`
	DoSubstitutions          bool    `json:"do_substitutions"`
	BitmapMonospace          bool    `json:"bitmap_monospace"`
	ForceAutohint            bool    `json:"force_autohint"`
	Dpi                      int     `json:"dpi"`
	QtUseSubpixelPositioning bool    `json:"qt_use_subpixel_positioning"`
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

func (c *Config) Init() {
	found, e := Exists("./config.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	if !found {
		c.newConfig()
	} else {
		c.loadConfigFile()
	}
}

func (c *Config) loadConfigFile() {
	file, e := os.ReadFile("./config.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	err := json.Unmarshal(file, &c)
	if err != nil {
		log.Println(err)
		c.newConfig()
	}

	if c.Customize && len(c.Fonts) == 0 {
		c.Customize = false
	}
}

func (c *Config) SaveConfigFile(apppath string) error { // https://github.com/spddl/csgo-reporter/blob/master/Config/Config.go#L147
	if !c.Customize {
		c.Fonts = []FontStruct{}
	}

	bytes, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return errors.New("config.json could not be saved (JSON error)")
	} else {
		err = os.WriteFile(filepath.Join(apppath, "config.json"), bytes, 0644)
		if err == nil {
			return nil
		} else {
			return err
		}
	}
}

func (c *Config) newConfig() {
	jsonBlob := json.RawMessage(`{
		"path": "",
		"font": "",
		"embeddedbitmap": false,
		"preferoutline": true,
		"dosubstitutions": true,
		"bitmapmonospace": false,
		"forceautohint": false,
		"qtusesubpixelpositioning": false,
		"dpi": 96
	}`)

	if err := json.Unmarshal(jsonBlob, c); err != nil {
		panic(err)
	}
}
