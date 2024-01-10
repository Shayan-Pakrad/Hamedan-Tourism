package application

import (
	"bytes"
	"database/sql"
	"errors"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
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

func (app *App) Start() {
	srv := &http.Server{
		Addr: os.Getenv("ADDR"),
		Handler: app.router,
	}

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		app.logger.Error("failed to start http server", "error", err)
		os.Exit(1)
	}
}

func (app *App) initLogger() {
	app.logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
}

func (app *App) initRoot() {
	if r, ok := os.LookupEnv("ROOT"); ok {
		app.root = r
		return
	}

	r, err := os.Getwd()
	if err != nil {
		app.logger.Error("failed to get working directory", "error", err)
		os.Exit(1)
	}

	app.root = r
}

func (app *App) initValidate() {
	app.validate = &validator.Validate{}
}

func (app *App) initDB() {
	sqldb, err := sql.Open(sqliteshim.ShimName, filepath.Join(app.root, "data.db"))
	if err != nil {
		app.logger.Error("failed to open database connection", "error", err)
		os.Exit(1)
	}

	app.db = bun.NewDB(sqldb, sqlitedialect.New())
}

func (app *App) initPages() {
	pages := template.New("pages")
	pages.Funcs(template.FuncMap{
		"partial": func(name string, data any) (template.HTML, error) {
			buf := new(bytes.Buffer)
			if err := pages.ExecuteTemplate(buf, name, data); err != nil {
				return "", err
			}

			return template.HTML(buf.String()), nil
		},
	})

	_, err := pages.ParseGlob(filepath.Join(app.root, "pages", "*.html"))
	if err != nil {
		app.logger.Error("failed to parse pages template", "error", err)
		os.Exit(1)
	}

	app.pages = pages
}

func (app *App) initComponents() {
	components := template.New("components")

	_, err := components.ParseGlob(filepath.Join(app.root, "components", "*.html"))
	if err != nil {
		app.logger.Error("failed to parse components template", "error", err)
		os.Exit(1)
	}

	app.components = components
}

func (app *App) initRouter() {
	app.router = chi.NewRouter()
}

func (app *App) initEH() {
	app.eh = func(w http.ResponseWriter, r *http.Request) {
		err := r.Context().Value(xmate.ErrorKey{}).(error)

		if httpErr := new(xmate.HTTPError); errors.As(err, &httpErr) {
			_ = xmate.WriteText(w, httpErr.Code, httpErr.Message)
			return
		}

		app.logger.Warn("failed to run http handler",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("error", err.Error()),
		)

		_ = xmate.WriteText(w, http.StatusInternalServerError, "Internal Server Error")
	}
}
