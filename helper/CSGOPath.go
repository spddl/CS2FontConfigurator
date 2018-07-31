package helper

import (
	"strings"

	"golang.org/x/sys/windows/registry"
)

func CSGOPath() (string, error) {
	var defaultPath string
	k, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\\Valve\\Steam`, registry.READ)
	if err != nil {
		return "", err
	}
	defer k.Close()

	s, _, err := k.GetStringValue("SteamPath")
	if err != nil {
		return "", err
	}

	s = strings.Replace(s, "/", "\\", -1)
	defaultPath = s + "\\SteamApps\\common\\Counter-Strike Global Offensive\\cstrike.exe"

	return defaultPath, nil
}
