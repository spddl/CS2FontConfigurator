package main

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/ConradIrwin/font/sfnt"
	"github.com/andygrunwald/vdf"
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
func FileAndExt(path string) (string, string) {
	for i := len(path) - 1; i >= 0 && path[i] != '/'; i-- {
		if path[i] == '.' {
			return path[:i], path[i:]
		}
	}
	return "", ""
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

func IndexOf(word string, data []string) int {
	for k, v := range data {
		if word == v {
			return k
		}
	}
	return -1
}

// FontData describes a font file and the various metadata associated with it.
type FontData struct { // https://github.com/ConradIrwin/font/blob/d797009a8098ca7f6c36a29d0c132a3d39bc4212/sfnt/table_name.go#L77
	Name     string
	Family   string
	Location string
	FileName string
	Metadata map[sfnt.NameID]string
	Data     []byte
}

// FontExtensions is a list of file extensions that denote fonts.
// Only files ending with these extensions will be installed.
var FontExtensions = map[string]bool{
	".otf": true,
	".ttf": true,
}

// NewFontData creates a new FontData struct.
// fileName is the font's file name, and data is a byte slice containing the font file data.
// It returns a FontData struct describing the font, or an error.
func NewFontData(location, fileName string, data []byte) (fontData *FontData, err error) {
	if _, ok := FontExtensions[path.Ext(fileName)]; !ok {
		return nil, fmt.Errorf("Not a font: %v", fileName)
	}

	fontData = &FontData{
		FileName: fileName,
		Metadata: make(map[sfnt.NameID]string),
		Data:     data,
		Location: location,
	}

	font, err := sfnt.Parse(bytes.NewReader(fontData.Data))
	if err != nil {
		return nil, err
	}

	if font.HasTable(sfnt.TagName) == false {
		return nil, fmt.Errorf("Font %v has no name table", fileName)
	}

	nameTable, err := font.NameTable()
	if err != nil {
		return nil, err
	}

	for _, nameEntry := range nameTable.List() {
		fontData.Metadata[nameEntry.NameID] = nameEntry.String()
	}

	fontData.Name = strings.ReplaceAll(fontData.Metadata[sfnt.NameFull], "\u0000", "")
	fontData.Family = fontData.Metadata[sfnt.NamePreferredFamily]
	if fontData.Family == "" {
		if v, ok := fontData.Metadata[sfnt.NameFontFamily]; ok {
			fontData.Family = v
		} else {
			fmt.Printf("Font %v has no font family!\n", fontData.Name)
		}
	}

	if fontData.Name == "" {
		fmt.Printf("Font %v has no name! Using file name instead.\n", fileName)
		fontData.Name = fileName
	}

	return
}

type LocalWindowsFonts struct {
	Folder   string
	Filename string
}

func MustReadDir(root string) []*FontData { // https://stackoverflow.com/questions/14668850/list-directory-in-go/49196644#49196644
	var files []*FontData

	f, err := os.Open(root)
	if err != nil {
		panic(err)
	}
	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		panic(err)
	}

	for _, file := range fileInfo {
		filename := file.Name()

		if _, ok := FontExtensions[path.Ext(filename)]; ok {
			// files = append(files, name)
			// log.Println(filepath.Join(root, filename))
			data, err := os.ReadFile(filepath.Join(root, filename))
			if err != nil {
				continue
			}

			fontData, err := NewFontData(root, filename, data)
			if err != nil {
				continue
			}

			files = append(files, fontData)

		}
	}

	return files
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func GetVDFKey(data map[string]interface{}, key string) map[string]interface{} {
	if val, ok := data[key]; ok {
		if val, ok := val.(map[string]interface{}); ok {
			return val
		}
	}
	return nil
}
