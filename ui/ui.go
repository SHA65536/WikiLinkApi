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
	// Input search for starting article
	ui.Router.Get("/searchstart", ui.SearchStartRoute)
	// Selection for starting article
	ui.Router.Get("/selectstart", ui.SelectStartRoute)
	// Input search for ending article
	ui.Router.Get("/searchend", ui.SearchEndRoute)
	// Selection for ending article
	ui.Router.Get("/selectend", ui.SelectEndRoute)
	// Results display
	ui.Router.Get("/results", ui.ResultsDisplayRoute)
	return ui, nil
}

func (u *UIHandler) Serve(addr string) error {
	return http.ListenAndServe(addr, u.Router)
}

// Input search for starting article
func (u *UIHandler) SearchStartRoute(w http.ResponseWriter, r *http.Request) {

}

// Selection for starting article
func (u *UIHandler) SelectStartRoute(w http.ResponseWriter, r *http.Request) {

}

// Input search for ending article
func (u *UIHandler) SearchEndRoute(w http.ResponseWriter, r *http.Request) {

}

// Selection for ending article
func (u *UIHandler) SelectEndRoute(w http.ResponseWriter, r *http.Request) {

}

// Results display
func (u *UIHandler) ResultsDisplayRoute(w http.ResponseWriter, r *http.Request) {

}
