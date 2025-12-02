package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"

	"github.com/drownedsound/blackdog/common"
	"github.com/drownedsound/blackdog/web"
)

type config struct {
	addr      string
	staticDir string
}

func main() {
	logger := slog.New(slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		},
	))

	svc := &common.Services{
		Logger: logger,
	}

	var cfg config
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(
		&cfg.staticDir, "static-dir", "./ui/static/", "Path to static assets",
	)
	flag.Parse()
	logger.Debug(
		"Retrieving server configuration",
		slog.String("addr", cfg.addr),
		slog.String("static-dir", cfg.staticDir),
	)

	logger.Info("Starting HTTP server", slog.String("addr", cfg.addr))
	err := http.ListenAndServe(cfg.addr, web.Routes(cfg.staticDir, svc))
	logger.Error(err.Error())
	os.Exit(1)
}
