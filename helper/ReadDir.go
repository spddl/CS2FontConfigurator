package helper

import (
	"os"
	"path"
)

func ReadDir(root string) ([]string, error) { // https://stackoverflow.com/questions/14668850/list-directory-in-go/49196644#49196644
	var files []string
	f, err := os.Open(root)
	if err != nil {
		return files, err
	}
	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		name := file.Name()
		_, ok := FontExtensions[path.Ext(name)]
		if ok {
			files = append(files, name)
		}
	}
	return files, nil
}
