package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/tomekzakrzewski/go-movierental/types"
)

func TestPostMovie(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	movieHandler := NewMovieHandler(tdb.Movie, tdb.Rent)
	app.Post("/", movieHandler.HandlePostMovie)

	params := types.CreateMovieParams{
		Title:  "The Matrix",
		Length: 120,
		Year:   1999,
		Genre:  []string{"Action"},
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	var movie types.Movie
	json.NewDecoder(resp.Body).Decode(&movie)
	if len(movie.ID) == 0 {
		t.Errorf("expecting a movie id to be set")
	}
	if movie.Title != params.Title {
		t.Errorf("expected title %s but got %s", params.Title, movie.Title)
	}
	if movie.Length != params.Length {
		t.Errorf("expected length %d but got %d", params.Length, movie.Length)
	}
	if movie.Year != params.Year {
		t.Errorf("expected year %d but got %d", params.Year, movie.Year)

	}
}

func TestGetMovies(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	movieHandler := NewMovieHandler(tdb.Movie, tdb.Rent)
	app.Post("/", movieHandler.HandlePostMovie)
	app.Get("/", movieHandler.HandleGetMovie)
	params := types.CreateMovieParams{
		Title:  "The Matrix",
		Length: 120,
		Year:   1999,
		Genre:  []string{"Action"},
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	req = httptest.NewRequest("GET", "/", nil)
	resp, err = app.Test(req)
	if err != nil {
		t.Error(err)
	}
	var movies []types.Movie
	json.NewDecoder(resp.Body).Decode(&movies)
	if len(movies) == 0 {
		t.Errorf("expecting a movie")
	}
}
