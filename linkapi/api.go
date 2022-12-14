package linkapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type ApiHandler struct {
	DB     *DatabaseHandler
	Search *SearchHandler
	Router *chi.Mux
}

func MakeApiHandler() (*ApiHandler, error) {
	var api = &ApiHandler{}
	var err error
	// Creating router
	api.Router = chi.NewRouter()
	api.Router.Use(middleware.Logger)
	api.Router.Get("/search", api.SearchRoute)

	// Creating database
	if api.DB, err = MakeDbHandler("bolt.db"); err != nil {
		return nil, err
	}
	if err := api.DB.CreateBuckets(); err != nil {
		return nil, err
	}

	// Creating search
	api.Search = MakeSearchHandler(api.DB)

	return api, nil
}

func (a *ApiHandler) Serve(addr string) error {
	defer a.DB.Close()
	return http.ListenAndServe(addr, a.Router)
}

type SearchResult struct {
	Error        string   `json:"error,omitempty"`
	ResultIds    []uint32 `json:"ids,omitempty"`
	ResultTitles []string `json:"titles,omitempty"`
}

func (a *ApiHandler) SearchRoute(w http.ResponseWriter, r *http.Request) {
	var res SearchResult
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	// Must have both params
	if end == "" || start == "" {
		res.Error = "must have 'start' and 'end' parameters!"
		render.JSON(w, r, res)
		return
	}
	// Finding start param
	startId, err := a.DB.GetId(start)
	if err != nil {
		res.Error = "start article not found!"
		render.JSON(w, r, res)
		return
	}
	// Finding end param
	endId, err := a.DB.GetId(end)
	if err != nil {
		res.Error = "end article not found!"
		render.JSON(w, r, res)
		return
	}

	// Finding path
	res.ResultIds, err = a.Search.ShortestPath(startId, endId, func(i int) {})
	if err != nil {
		res.Error = "could not find path!"
		render.JSON(w, r, res)
		return
	}

	// Finding names
	res.ResultTitles, err = a.DB.IdsToNames(res.ResultIds...)
	if err != nil {
		res.Error = "error parsing path!"
		render.JSON(w, r, res)
		return
	}
	render.JSON(w, r, res)
}
