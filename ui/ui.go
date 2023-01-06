package ui

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

const WikiSearchEndpoint = `https://he.wikipedia.org/w/api.php?action=query&list=search&srnamespace=0&srlimit=5&prop=info&utf8=&format=json&origin=*&srsearch=`
const WikiRandomEndpoint = "https://he.wikipedia.org/w/api.php?action=query&generator=random&grnnamespace=0&grnlimit=1&prop=info|extracts&exlimit=1&explaintext=true&exsentences=1&utf8=&format=json&origin=*"

//go:embed static/styles.css
var styles []byte

//go:embed static/main.js
var mainjs []byte

//go:embed static/index.html
var indexhtml []byte

type UIHandler struct {
	Locale  string
	LinkAPI string
	Client  *http.Client
	Router  *chi.Mux
}

func MakeUIHandler(locale string, api_url string) (*UIHandler, error) {
	var ui = &UIHandler{
		Locale:  locale,
		LinkAPI: api_url,
		Client:  http.DefaultClient,
	}
	ui.Router = chi.NewRouter()
	ui.Router.Use(middleware.Logger)
	// Main route
	ui.Router.Get("/", ui.MainRoute)
	// Search route
	ui.Router.Get("/search", ui.SearchRoute)
	// Random route
	ui.Router.Get("/random", ui.RandomRoute)
	// Result route
	ui.Router.Get("/result", ui.ResultRoute)

	// Static files
	ui.Router.Get("/main.js", func(w http.ResponseWriter, r *http.Request) {
		w.Write(mainjs)
	})
	ui.Router.Get("/styles.css", func(w http.ResponseWriter, r *http.Request) {
		w.Write(styles)
	})
	return ui, nil
}

func (u *UIHandler) Serve(addr string) error {
	return http.ListenAndServe(addr, u.Router)
}

// Main webpage route
func (u *UIHandler) MainRoute(w http.ResponseWriter, r *http.Request) {
	w.Write(indexhtml)
}

// Search for articles
func (u *UIHandler) SearchRoute(w http.ResponseWriter, r *http.Request) {
	var res = &SearchResult{}
	query := r.URL.Query().Get("q")
	if query == "" {
		w.Write([]byte("error, please provide a query"))
		return
	}
	fmt.Println(WikiSearchEndpoint + url.QueryEscape(query))
	resp, err := u.Client.Get(WikiSearchEndpoint + url.QueryEscape(query))
	if err != nil {
		w.Write([]byte("error, search failed"))
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	if err != nil {
		w.Write([]byte("error, search failed"))
		fmt.Println(err)
		return
	}

	if err := json.NewDecoder(resp.Body).Decode(res); err != nil {
		w.Write([]byte("error, search failed"))
		fmt.Println(err)
		return
	}

	render.JSON(w, r, res)
}

// Random Articles
func (u *UIHandler) RandomRoute(w http.ResponseWriter, r *http.Request) {

}

// Results for path
func (u *UIHandler) ResultRoute(w http.ResponseWriter, r *http.Request) {

}
