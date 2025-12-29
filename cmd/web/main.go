package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	app "github.com/drownedsound/blackdog/internal/features/application"
	"github.com/drownedsound/blackdog/internal/infra"

	_ "github.com/mattn/go-sqlite3"
)

type config struct {
	addr      string
	staticDir string
	dbPath    string
}

func initSrvDependencies(cfg config) (*slog.Logger, app.Repository) {
	// TODO: Replicate command-line flags in a TOML file
	// TODO: Create command-line flag to direct log output to file or stdout
	// TODO: Create command-line flag to set log format
	// TODO: Use slog.NewTextHandler as default log format
	// TODO: Use slog.NewTextHandler for stdout
	// TODO: Use slog.NewJSONHandler for log files
	logger := slog.New(slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{
			// TODO: Create command-line flag to enable error/info/debug mode
			Level: slog.LevelDebug,
			// TODO: Enable AddSource on debug mode only
			AddSource: false,
		},
	))

	// TODO: Dynamically determine which repository to use based on flag
	// repo = infra.NewMockDb()
	// logger.Debug("Initialized in-memory database")

	db, err := infra.NewSQLiteConnection(cfg.dbPath)
	if err != nil {
		logger.Error("Failed to initialize database", slog.Any("error", err))
		os.Exit(1)
	}
	logger.Debug(
		"Initialized SQLite database",
		slog.String("path", cfg.dbPath))

	idGen, err := infra.NewSnowflakeIdGenerator(1)
	if err != nil {
		logger.Error("Failed to create id generator", slog.Any("error", err))
	}

	repo, err := infra.NewApplicationRepository(db, idGen)
	if err != nil {
		logger.Error("Failed to initialize repo", slog.Any("error", err))
		os.Exit(1)
	}

	return logger, repo
}

func loadConfig() config {
	var cfg config

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

	// Default to local file
	// Overridden by env vars or flags in production
	flag.StringVar(
		&cfg.dbPath,
		"db-path",
		"./data/blackdog.db",
		"Path to SQLite database file",
	)

	flag.Parse()
	return cfg
}

func main() {
	cfg := loadConfig()
	logger, repo := initSrvDependencies(cfg)

	logger.Debug(
		"Retrieved server configuration",
		slog.String("addr", cfg.addr),
		slog.String("static-dir", cfg.staticDir),
		slog.String("db-path", cfg.dbPath),
	)

	svc := app.NewService(repo, logger)
	mux := http.NewServeMux()
	handler := app.NewHandler(svc)
	handler.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:    cfg.addr,
		Handler: app.RequestLogger(logger)(mux),
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
