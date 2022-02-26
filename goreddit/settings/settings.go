package settings

import (
	"github.com/BurntSushi/toml"
	"github.com/panta82/goreddit/lib"
	"os"
)

type DatabaseSettings struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type Settings struct {
	Database DatabaseSettings
}

const LOCAL_SETTINGS_FILE = "settings.local.toml"

var SOURCE_PATHS = [...]string{
	"./" + LOCAL_SETTINGS_FILE,
	"../" + LOCAL_SETTINGS_FILE,
}

func Load() (*Settings, error) {
	localSettingsPath, err := lib.LookupRelativePath(SOURCE_PATHS[:], nil)
	if err != nil {
		return nil, err
	}

	localSettingsToml, err := os.ReadFile(*localSettingsPath)
	if err != nil {
		return nil, err
	}

	var settings Settings
	_, err = toml.Decode(string(localSettingsToml), &settings)
	if err != nil {
		return nil, err
	}

	return &settings, nil
}

func LoadOrDie() Settings {
	settings, err := Load()
	if err != nil {
		panic("Failed to load settings: " + err.Error())
	}
	return *settings
}
