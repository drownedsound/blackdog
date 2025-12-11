package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	app "github.com/drownedsound/blackdog/internal/features/application"
	"github.com/drownedsound/blackdog/internal/infra"
)

type config struct {
	addr      string
	staticDir string
}

func initSrvDependencies() (logger *slog.Logger, repo app.Repository) {
	// TODO: Replicate command-line flags in a TOML file
	// TODO: Create command-line flag to direct log output to file or stdout
	// TODO: Create command-line flag to set log format
	// TODO: Use slog.NewTextHandler as default log format
	// TODO: Use slog.NewTextHandler for stdout
	// TODO: Use slog.NewJSONHandler for log files
	logger = slog.New(slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{
			// TODO: Create command-line flag to enable error/info/debug mode
			Level: slog.LevelDebug,
			// TODO: Enable AddSource on debug mode only
			AddSource: false,
		},
	))

	repo = infra.NewMockDb()
	logger.Debug("Initialized in-memory database")

	return
}

func loadConfig() (cfg config) {
	flag.StringVar(
		&cfg.addr,
		"addr",
		":4000",
		"HTTP network address",
	)

	flag.StringVar(
		&cfg.staticDir,
		"static-dir",
		"./ui/static/",
		"Path to static assets",
	)

	flag.Parse()
	return
}

func main() {
	logger, repo := initSrvDependencies()
	logger.Debug("Initialized server dependencies")

	cfg := loadConfig()
	logger.Debug(
		"Retrieved server configuration",
		slog.String("addr", cfg.addr),
		slog.String("static-dir", cfg.staticDir),
	)

	// TODO: Pass logger to svc as a dependency
	svc := app.NewService(repo, logger)
	mux := http.NewServeMux()
	handler := app.NewHandler(svc)
	handler.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:    cfg.addr,
		Handler: mux,
		// TODO: Create command-line flag to configure idle/read/write timeout
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Info("Starting HTTP server", slog.String("addr", cfg.addr))
	if err := srv.ListenAndServe(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
