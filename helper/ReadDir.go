package helper

import (
	"os"
	"strings"
)

func ReadDir(root, suffix string) ([]string, error) { // https://stackoverflow.com/questions/14668850/list-directory-in-go/49196644#49196644
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
		if suffix != "" {
			name := file.Name()
			if strings.HasSuffix(name, suffix) {
				files = append(files, name)
			}
		} else {
			files = append(files, file.Name())
		}

	}
	return files, nil
}
