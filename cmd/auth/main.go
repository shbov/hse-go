package main

import (
	"context"
	"flag"
	"github.com/anonimpopov/hw4/internal/logger"
	"github.com/juju/zaputil/zapctx"
	"log"

	"github.com/anonimpopov/hw4/internal/app"
)

var DSN = "https://7112fabeed3af34b9b72f7879b870b0a@o4506218171793408.ingest.sentry.io/4506218174545920"

func getConfigPath() string {
	var configPath string

	flag.StringVar(&configPath, "c", "../../.config/auth.yaml", "path to config file")
	flag.Parse()

	return configPath
}

func main() {
	lg, err := logger.GetLogger(true, DSN, "development")
	if err != nil {
		log.Fatal(err.Error())
	}

	path := getConfigPath()
	config, err := app.NewConfig(path)
	if err != nil {
		lg.Fatal(err.Error())
	}

	ctx := zapctx.WithLogger(context.Background(), lg)
	a, err := app.New(ctx, config)
	if err != nil {
		lg.Fatal(err.Error())
	}

	if err := a.Serve(ctx); err != nil {
		lg.Fatal(err.Error())
	}
}
