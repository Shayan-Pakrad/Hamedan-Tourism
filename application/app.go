package application

import (
	"bytes"
	"database/sql"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

type App struct {
	Logger     *slog.Logger
	Root       string
	DB         *bun.DB
	Pages      *template.Template
	Components *template.Template
	Echo       *echo.Echo
}

func New() *App {
	app := new(App)

	app.Logger = newLogger()

	if r, ok := os.LookupEnv("ROOT"); ok {
		app.Root = r
	} else {
		r, err := os.Getwd()
		if err != nil {
			app.Logger.Error("failed to get working directory", "error", err)
			os.Exit(1)
		}

		app.Root = r
	}

	db, err := newDB()
	if err != nil {
		app.Logger.Error("failed to open database connection", "error", err)
		os.Exit(1)
	}
	app.DB = db

	pages, err := newPages()
	if err != nil {
		app.Logger.Error("failed to parse pages template", "error", err)
		os.Exit(1)
	}
	app.Pages = pages

	components, err := newComponents()
	if err != nil {
		app.Logger.Error("failed to parse components template", "error", err)
		os.Exit(1)
	}
	app.Components = components

	app.Echo = newEcho(app.Logger)

	return app
}

func newLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}

func newDB() (*bun.DB, error) {
	sqldb, err := sql.Open(sqliteshim.ShimName, "./data.db") // TODO: use Root
	if err != nil {
		return nil, err
	}

	return bun.NewDB(sqldb, sqlitedialect.New()), nil
}

func newPages() (*template.Template, error) {
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

	return pages.ParseGlob("./pages/*.html") // TODO: use Root
}

func newComponents() (*template.Template, error) {
	components := template.New("components")
	return components.ParseGlob("./components/*.html") // TODO: use Root
}

type customValidator struct {
	validate *validator.Validate
}

func (cv customValidator) Validate(s any) error {
	if err := cv.validate.Struct(s); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

func newEcho(logger *slog.Logger) *echo.Echo {
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true
	e.Validator = customValidator{validate: validator.New()}

	/*
		e.HTTPErrorHandler = func(err error, c echo.Context) {
			if httpErr := new(echo.HTTPError); errors.As(err, &httpErr) {
				c.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
				c.Response().Header().Set("X-Content-Type-Options", "nosniff")
				c.Response().WriteHeader(httpErr.Code)

				_ = errorPage.Execute(c.Response(), TODO: pass http status and message)
			}

			logger.Warn("failed to run http handler",
				slog.String("method", c.Request().Method),
				slog.String("path", c.Request().URL.Path),
				slog.String("error", err.Error()),
			)
		}
	*/

	return e
}
