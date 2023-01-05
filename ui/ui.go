package ui

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type UIHandler struct {
	Locale  string
	LinkAPI string
	Router  *chi.Mux
}

func MakeUIHandler(locale string, api_url string) (*UIHandler, error) {
	var ui = &UIHandler{
		Locale:  locale,
		LinkAPI: api_url,
	}
	ui.Router = chi.NewRouter()
	ui.Router.Use(middleware.Logger)
	// Main route
	ui.Router.Get("/", ui.MainRoute)
	// Search route
	ui.Router.Get("/search", ui.SearchRoute)
	// Random route
	ui.Router.Get("/search", ui.RandomRoute)
	// Result route
	ui.Router.Get("/search", ui.ResultRoute)
	return ui, nil
}

func (u *UIHandler) Serve(addr string) error {
	return http.ListenAndServe(addr, u.Router)
}

// Main webpage route
func (u *UIHandler) MainRoute(w http.ResponseWriter, r *http.Request) {

}

// Search for articles
func (u *UIHandler) SearchRoute(w http.ResponseWriter, r *http.Request) {

}

// Random Articles
func (u *UIHandler) RandomRoute(w http.ResponseWriter, r *http.Request) {

}

// Results for path
func (u *UIHandler) ResultRoute(w http.ResponseWriter, r *http.Request) {

}
