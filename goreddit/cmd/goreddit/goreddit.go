package main

import (
	"fmt"
	"github.com/panta82/goreddit/postgres"
	settingsLib "github.com/panta82/goreddit/settings"
	"github.com/panta82/goreddit/web"
	"log"
	"net/http"
)

func main() {
	settings := settingsLib.Must(settingsLib.Load())

	store, err := postgres.NewStore(settings.Postgres.ConnectionString())
	if err != nil {
		log.Fatalf("Failed to initialize store: %w", err)
	}

	h := web.NewHandler(store)

	err = http.ListenAndServe(fmt.Sprintf(":%d", settings.Web.Port), h)
	if err != nil {
		log.Fatalf("Web server failed: %w", err)
	}
}
