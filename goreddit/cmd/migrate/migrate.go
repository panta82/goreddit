package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/panta82/goreddit/lib"
	settingsLib "github.com/panta82/goreddit/settings"
)

var CLI struct {
	Up struct {
	} `cmd:"" help:"Migrate up."`

	Down struct {
	} `cmd:"" help:"Migrate down."`
}

func main() {
	ctx := kong.Parse(&CLI,
		kong.Name("migrate"),
		kong.Description("Migration util."),
		kong.UsageOnError())

	migrationsDir, err := lib.LookupRelativePath([]string{"../../migrations", "../migrations", "./migrations"}, nil)
	if err != nil {
		panic(err)
	}
	if migrationsDir == nil {
		panic("migrations directory not found")
	}
	fmt.Printf("Running migrations from: %s\n", *migrationsDir)

	settings := settingsLib.LoadOrDie()

	m, err := migrate.New(
		"file://"+*migrationsDir,
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s", settings.Database.User, settings.Database.Password, settings.Database.Host, settings.Database.Port, settings.Database.Name))
	if err != nil {
		panic(err)
	}

	switch ctx.Command() {
	case "up":
		up(m)
	case "down":
		down(m)
	default:
		panic(ctx.Command())
	}
}

func up(m *migrate.Migrate) {
	fmt.Printf("Migrating up...\n")
	if err := m.Up(); err != nil {
		panic(err)
	}
}
func down(m *migrate.Migrate) {
	fmt.Printf("Migrating down...\n")
	if err := m.Down(); err != nil {
		panic(err)
	}
}
