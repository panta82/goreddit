package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/panta82/goreddit/lib"
	settingsLib "github.com/panta82/goreddit/settings"
	"log"
)

var CLI struct {
	Up struct {
	} `cmd:"" help:"Migrate up."`

	Down struct {
	} `cmd:"" help:"Migrate down."`
}

type migrateLogger struct{}

func (m migrateLogger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}
func (m migrateLogger) Verbose() bool {
	return true
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

	settings := settingsLib.Must(settingsLib.Load())

	m, err := migrate.New("file://"+*migrationsDir, settings.Postgres.ConnectionString())
	if err != nil {
		panic(err)
	}
	m.Log = migrateLogger{}

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
		if err.Error() == "no change" {
			fmt.Printf("No change.\n")
		} else {
			panic(err)
		}
	}
}
func down(m *migrate.Migrate) {
	fmt.Printf("Migrating down...\n")
	if err := m.Down(); err != nil {
		if err.Error() == "no change" {
			fmt.Printf("No change.\n")
		} else {
			panic(err)
		}
	}
}
