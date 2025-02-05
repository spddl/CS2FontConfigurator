package main

import (
	"log"
	"os"
	"path/filepath"
	"slices"

	"github.com/andygrunwald/vdf"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/sys/windows/registry"
)

func CSPath() (string, error) {
	k, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Valve\Steam`, registry.READ)
	if err != nil {
		return "", err
	}
	defer k.Close()

	s, _, err := k.GetStringValue("SteamPath")
	if err != nil {
		return "", err
	}

	f, err := os.Open(filepath.Join(s, "SteamApps", "libraryfolders.vdf"))
	if err != nil {
		panic(err)
	}

	p := vdf.NewParser(f)
	file, err := p.Parse()
	if err != nil {
		panic(err)
	}

	libraryfolders := GetVDFKey(file, "libraryfolders")
	for _, folders := range libraryfolders {
		apps := GetVDFKey(folders.(map[string]interface{}), "apps")
		for id := range apps {
			if id == "730" {
				if val, ok := folders.(map[string]interface{}); ok {
					if val, ok := val["path"]; ok {
						return filepath.Join(val.(string), "SteamApps", "common", "Counter-Strike Global Offensive"), nil
					}
				}
			}
		}
	}

	return "", nil
}

// https://golang.org/src/path/path.go?s=4371:4399#L158
func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func GetVDFKey(data map[string]interface{}, key string) map[string]interface{} {
	if val, ok := data[key]; ok {
		if val, ok := val.(map[string]interface{}); ok {
			return val
		}
	}
	return nil
}

func (c *Config) ClearUpFontsFolder() {
	panorama_fonts := filepath.Join(c.Path, "game", "csgo", "panorama", "fonts")
	entries, err := os.ReadDir(panorama_fonts)
	if err != nil {
		log.Fatal(err)
	}

	valveFonts := []string{
		"notomono-regular.ttf",
		"notosans-bold.ttf",
		"notosans-bolditalic.ttf",
		"notosans-italic.ttf",
		"notosans-regular.ttf",
		"notosansjp-bold.ttf",
		"notosansjp-light.ttf",
		"notosansjp-regular.ttf",
		"notosanskr-bold.ttf",
		"notosanskr-light.ttf",
		"notosanskr-regular.ttf",
		"notosanssc-bold.ttf",
		"notosanssc-light.ttf",
		"notosanssc-regular.ttf",
		"notosanssymbols-regular.ttf",
		"notosanstc-bold.ttf",
		"notosanstc-light.ttf",
		"notosanstc-regular.ttf",
		"notosansthai-bold.ttf",
		"notosansthai-light.ttf",
		"notosansthai-regular.ttf",
		"notoserif-bold.ttf",
		"notoserif-boliitalic.ttf",
		"notoserif-italic.ttf",
		"notoserif-regular.ttf",
	}

	for _, e := range entries {
		filename := e.Name()
		switch filepath.Ext(filename) {
		case ".ttf", ".otf":
			if !slices.Contains(valveFonts, filename) {
				if err := os.Remove(filepath.Join(panorama_fonts, filename)); err != nil {
					log.Println(err)
				}
			}
		case ".conf", ".uifont":
			// ignore
		case ".zip", ".7z", ".rar":
			// ignore
		default:
			if err := os.Remove(filepath.Join(panorama_fonts, filename)); err != nil {
				log.Println(err)
			}
		}
	}
}

func copyFile(in, out string) {
	data, err := os.ReadFile(in)
	if err != nil {
		log.Fatal(err)
	}

	if err = os.WriteFile(out, data, 0o644); err != nil {
		log.Fatal(err)
	}
}

func GetFontName(fontfilename string) string {
	data, err := os.ReadFile(fontfilename)
	if err != nil {
		log.Fatal(err)
	}
	f, err := opentype.Parse(data)
	if err != nil {
		// Windows Fallback
		return GetFontResourceInfo(fontfilename)
	} else {
		fontname, err := f.Name(nil, sfnt.NameIDFamily)
		if err != nil {
			log.Fatal(err)
		}
		return fontname
	}

}
