package application

import (
	"context"
	"errors"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"example.com/hamedan-tourism/model"
	"example.com/hamedan-tourism/resource"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/uptrace/bun"
	"github.com/utilyre/xmate"
)

type App struct {
	logger *slog.Logger
	root   string

	db       *bun.DB
	validate *validator.Validate

	pages      *template.Template
	components *template.Template

	router chi.Router
	eh     xmate.ErrorHandler
}

func New() *App {
	app := new(App)

	app.initLogger()
	app.initRoot()
	app.initDB()
	app.initValidate()
	app.initPages()
	app.initComponents()
	app.initRouter()
	app.initEH()

	return app
}

func (app *App) Setup() {
	app.registerRoutes()
	app.createTables()
}

func (app *App) Start() {
	srv := &http.Server{
		Addr:    os.Getenv("ADDR"),
		Handler: app.router,
	}

	app.logger.Info("starting to listen and serve", "address", srv.Addr)
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		app.logger.Error("failed to start http server", "error", err)
		os.Exit(1)
	}
}

func (app *App) createTables() {
	ctx := context.Background()

	if _, err := app.db.
		NewCreateTable().
		IfNotExists().
		Model((*model.Attraction)(nil)).
		Exec(ctx); err != nil {
		app.logger.Error("failed to create attractions table", "error", err)
		os.Exit(1)
	}

	if _, err := app.db.
		NewCreateTable().
		IfNotExists().
		Model((*model.Blog)(nil)).
		Exec(ctx); err != nil {
		app.logger.Error("failed to create blogs table", "error", err)
		os.Exit(1)
	}

	if _, err := app.db.
		NewCreateTable().
		IfNotExists().
		Model((*model.Event)(nil)).
		Exec(ctx); err != nil {
		app.logger.Error("failed to create events table", "error", err)
		os.Exit(1)
	}
}

func (app *App) registerRoutes() {
	app.router.Mount("/", resource.PageResource{
		Logger: app.logger,
		DB:     app.db,
		Pages:  app.pages,
		EH:     app.eh,
	}.Routes())

	app.router.Handle("/static/*", http.StripPrefix("/static/",
		http.FileServer(http.Dir(filepath.Join(app.root, "public")))))
}
