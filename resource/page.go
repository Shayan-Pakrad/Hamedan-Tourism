package resource

import (
	"html/template"
	"log/slog"
	"net/http"
	"strconv"

	"example.com/hamedan-tourism/model"
	"github.com/go-chi/chi/v5"
	"github.com/uptrace/bun"
	"github.com/utilyre/xmate"
)

type PageResource struct {
	Logger *slog.Logger
	DB     *bun.DB
	Pages  *template.Template
	EH     xmate.ErrorHandler
}

type pageProps struct {
	Name string
	Data any
}

func (pr PageResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", pr.EH.HandleFunc(pr.home))
	r.Get("/attractions", pr.EH.HandleFunc(pr.attractions))
	r.Get("/attractions/{id:[0-9]+}", pr.EH.HandleFunc(pr.attraction))
	r.Get("/blogs", pr.EH.HandleFunc(pr.blogs))
	r.Get("/blogs/{id:[0-9]+}", pr.EH.HandleFunc(pr.blog))
	r.Get("/events", pr.EH.HandleFunc(pr.events))
	r.Get("/criticism", pr.EH.HandleFunc(pr.criticism))
	r.Get("/login", pr.EH.HandleFunc(pr.login))
	r.Get("/signup", pr.EH.HandleFunc(pr.signup))

	return r
}

func (pr PageResource) home(w http.ResponseWriter, r *http.Request) error {
	return xmate.WriteHTML(w, pr.Pages, http.StatusOK, pageProps{Name: "home.html"})
}

func (pr PageResource) attractions(w http.ResponseWriter, r *http.Request) error {
	attractions := []model.Attraction{}
	if err := pr.DB.NewSelect().Model(&attractions).Scan(r.Context()); err != nil {
		return err
	}

	return xmate.WriteHTML(w, pr.Pages, http.StatusOK, pageProps{
		Name: "attractions.html",
		Data: attractions,
	})
}

type Attraction struct {
	ID       int64         `bun:"id"`
	ImageURL string        `bun:"image_url"`
	Title    string        `bun:"title"`
	Brief    string        `bun:"brief"`
	Content  template.HTML `bun:"content"`
}

func (pr PageResource) attraction(w http.ResponseWriter, r *http.Request) error {
	strID := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		return xmate.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	attraction := new(Attraction)
	if err := pr.DB.NewSelect().Model((*model.Attraction)(nil)).Where("id = ?", id).Scan(r.Context(), attraction); err != nil {
		return err
	}

	return xmate.WriteHTML(w, pr.Pages, http.StatusOK, pageProps{
		Name: "attraction.html",
		Data: attraction,
	})
}

func (pr PageResource) blogs(w http.ResponseWriter, r *http.Request) error {
	blogs := []model.Blog{}
	if err := pr.DB.NewSelect().Model(&blogs).Scan(r.Context()); err != nil {
		return err
	}

	return xmate.WriteHTML(w, pr.Pages, http.StatusOK, pageProps{
		Name: "blogs.html",
		Data: blogs,
	})
}

type Blog struct {
	ID      int64         `bun:"id"`
	Title   string        `bun:"title"`
	Content template.HTML `bun:"content"`
	Brief   string        `bun:"brief"`
	Views   int           `bun:"views,default:0"`
}

func (pr PageResource) blog(w http.ResponseWriter, r *http.Request) error {
	strID := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		return xmate.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	blog := new(Blog)
	if err := pr.DB.NewSelect().Model((*model.Blog)(nil)).Where("id = ?", id).Scan(r.Context(), blog); err != nil {
		return err
	}

	return xmate.WriteHTML(w, pr.Pages, http.StatusOK, pageProps{
		Name: "blog.html",
		Data: blog,
	})
}

func (pr PageResource) events(w http.ResponseWriter, r *http.Request) error {
	events := []model.Event{}
	if err := pr.DB.NewSelect().Model(&events).Scan(r.Context()); err != nil {
		return err
	}

	return xmate.WriteHTML(w, pr.Pages, http.StatusOK, pageProps{
		Name: "events.html",
		Data: events,
	})
}

func (pr PageResource) criticism(w http.ResponseWriter, r *http.Request) error {
	return xmate.WriteHTML(w, pr.Pages, http.StatusOK, pageProps{
		Name: "criticism.html",
	})
}

func (pr PageResource) login(w http.ResponseWriter, r *http.Request) error {
	return xmate.WriteHTML(w, pr.Pages, http.StatusOK, pageProps{
		Name: "sign_in.html",
	})
}

func (pr PageResource) signup(w http.ResponseWriter, r *http.Request) error {
	return xmate.WriteHTML(w, pr.Pages, http.StatusOK, pageProps{
		Name: "sign_up.html",
	})
}
