package wikilinkapi

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
)

type ApiHandler struct {
	DB     *DatabaseHandler
	Search *SearchHandler
	Router *chi.Mux
	Logger zerolog.Logger
}

func MakeApiHandler(db_path string, logLevel zerolog.Level, writer io.Writer) (*ApiHandler, error) {
	var api = &ApiHandler{}
	var err error

	// Creating logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logwriter := io.MultiWriter(os.Stdout, writer)
	api.Logger = zerolog.New(logwriter).With().Str("service", "linkapi").Timestamp().Logger().Level(logLevel)

	// Creating router
	api.Router = chi.NewRouter()
	api.Router.Get("/search", api.SearchRoute)
	api.Logger.Debug().Msg("created router")

	// Checking DB file exists
	if _, err := os.Stat(db_path); err != nil {
		api.Logger.Error().Msgf("database not found! %s", db_path)
		return nil, fmt.Errorf("database file does not exist")
	}

	// Creating database handler
	if api.DB, err = MakeDbHandler(db_path); err != nil {
		api.Logger.Error().Msgf("error creating db handler! %s", err)
		return nil, err
	}
	if err := api.DB.CreateBuckets(); err != nil {
		api.Logger.Error().Msgf("error creating db buckers! %s", err)
		return nil, err
	}

	// Creating search handler
	api.Search = MakeSearchHandler(api.DB)
	api.Logger.Info().Msg("created handler succesfully!")

	return api, nil
}

func (a *ApiHandler) Serve(addr string) error {
	defer a.DB.Close()
	a.Logger.Info().Msgf("serving linkapi on %s", addr)
	return http.ListenAndServe(addr, a.Router)
}

type SearchResult struct {
	Error        string   `json:"error,omitempty"`
	ResultIds    []uint32 `json:"ids,omitempty"`
	ResultTitles []string `json:"titles,omitempty"`
}

func (a *ApiHandler) SearchRoute(w http.ResponseWriter, r *http.Request) {
	var res SearchResult
	sTime := time.Now()
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	// Must have both params
	if end == "" || start == "" {
		res.Error = "must have 'start' and 'end' parameters!"
		render.JSON(w, r, res)
		a.Logger.Info().Str("took", timeToMs(time.Since(sTime))).Str("src", r.Host).Str("dst", "/search").Msg("not enough params in request")
		return
	}
	// Finding start param
	startId, err := a.DB.GetId(start)
	if err != nil {
		res.Error = "start article not found!"
		render.JSON(w, r, res)
		a.Logger.Info().Str("took", timeToMs(time.Since(sTime))).Str("src", r.Host).Str("dst", "/search").Msgf("article %s not found", start)
		return
	}
	// Finding end param
	endId, err := a.DB.GetId(end)
	if err != nil {
		res.Error = "end article not found!"
		render.JSON(w, r, res)
		a.Logger.Info().Str("took", timeToMs(time.Since(sTime))).Str("src", r.Host).Str("dst", "/search").Msgf("article %s not found", end)
		return
	}

	// Finding path
	res.ResultIds, err = a.Search.ShortestPath(startId, endId, func(i int) {})
	if err != nil {
		res.Error = "could not find path!"
		render.JSON(w, r, res)
		a.Logger.Warn().Str("took", timeToMs(time.Since(sTime))).Str("src", r.Host).Str("dst", "/search").Msgf("path not found %s - %s", start, end)
		return
	}

	// Finding names
	res.ResultTitles, err = a.DB.IdsToNames(res.ResultIds...)
	if err != nil {
		res.Error = "error parsing path!"
		render.JSON(w, r, res)
		a.Logger.Warn().Str("took", timeToMs(time.Since(sTime))).Str("src", r.Host).Str("dst", "/search").Msgf("error parsing path %s - %s", start, end)
		return
	}

	// Checking for no path
	if len(res.ResultIds) == 0 {
		res.Error = "no path found!"
		render.JSON(w, r, res)
		a.Logger.Info().Str("took", timeToMs(time.Since(sTime))).Str("src", r.Host).Str("dst", "/search").Msgf("no path %s - %s", start, end)
		return
	}
	render.JSON(w, r, res)
	a.Logger.Info().Str("took", timeToMs(time.Since(sTime))).Str("src", r.Host).Str("dst", "/search").Msg("success")
}

func timeToMs(t time.Duration) string {
	return fmt.Sprintf("%dms", t/time.Millisecond)
}
