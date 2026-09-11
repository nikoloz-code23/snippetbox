package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/lib/pq"
	"snippetbox.alexedwards.net/internal/models"
)

type dbData struct {
	pgPassword	string
	pgUser 			string
	pgDb 				string
}

type config struct {
	addr 			string
	staticDir string
	db 				dbData
}

type application struct {
	logger *slog.Logger
	snippets *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	var cfg config
	var app *application

	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP Network Address")
	cfg.staticDir = *flag.String("static-dir", "./ui/static/", "Static Directory Path")

	cfg.db.pgDb 			= os.Getenv("DB_NAME")
	cfg.db.pgPassword = os.Getenv("DB_PASSWORD")
	cfg.db.pgUser 		= os.Getenv("DB_USER")

	// First arg is the destination
	// Second arg is a pointer to the handler options struct. Can be nil if not needed
	loggerHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level: slog.LevelDebug,
	})
	logger := slog.New(loggerHandler)

	app = &application{
		logger: logger,
	}

	// Without this, we will always use the default value that we
	// give to a flag variable.
	flag.Parse()

	// Config for the database.
	dbCfg := pq.Config{
		Host:           "localhost",
		Port:           5432,
		User:           cfg.db.pgUser,
		Password: 			cfg.db.pgPassword,
		Database: 			cfg.db.pgDb,
		SSLMode: 				pq.SSLModeDisable,
		ConnectTimeout: 5 * time.Second,
	}

	dbConnector, err := pq.NewConnectorConfig(dbCfg)
	if err != nil {
		app.logger.Error(err.Error())
	}

	// We use the connector to return an *sql.DB
	db, err := openDB(dbConnector)
	if err != nil {
		app.logger.Error(err.Error())
	}
	defer db.Close() // Make the DB close at the end of the main() block.
	
	// Template Cache Init
	templateCache, err := newTemplateCache()
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
	
	app.templateCache = templateCache;
	app.snippets = &models.SnippetModel{DB: db}

	app.logger.Info("Started Server", slog.String("addr", "http://localhost" + cfg.addr))

	err = http.ListenAndServe(cfg.addr, app.routes(&cfg))
	// Will only crash the entire program if `err` is true.
	app.logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dbConnector *pq.Connector) (*sql.DB, error) {
	db := sql.OpenDB(dbConnector)

	// Make sure it works
	err := db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}