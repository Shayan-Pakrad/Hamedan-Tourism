package resource

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/utilyre/xmate"
)

type PageResource struct {
	Pages *template.Template
	EH    xmate.ErrorHandler
}

type pageProps struct {
	Name string
	Data any
}

func (pr PageResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", pr.EH.HandleFunc(pr.home))

	return r
}

func (pr PageResource) home(w http.ResponseWriter, r *http.Request) error {
	return xmate.WriteHTML(w, pr.Pages, http.StatusOK, pageProps{Name: "home.html"})
}
