package settings

import (
	"github.com/BurntSushi/toml"
	"github.com/panta82/goreddit/lib"
	"github.com/panta82/goreddit/postgres"
	"github.com/panta82/goreddit/web"
	"os"
)

type Settings struct {
	Postgres postgres.Settings
	Web      web.Settings
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

	var settings = Settings{
		Postgres: *postgres.NewSettings(),
		Web:      *web.NewSettings(),
	}
	_, err = toml.Decode(string(localSettingsToml), &settings)
	if err != nil {
		return nil, err
	}

	return &settings, nil
}

func Must(settings *Settings, err error) Settings {
	if err != nil {
		panic("Failed to load settings: " + err.Error())
	}
	return *settings
}
