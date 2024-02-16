package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/tomekzakrzewski/go-movierental/db/fixtures"
	"github.com/tomekzakrzewski/go-movierental/types"
)

func TestPostMovie(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	movieHandler := NewMovieHandler(tdb.Store)
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
	fixtures.AddMovie(tdb.Store, "The Matrix", []string{"Action"}, 120, 1999)
	app := fiber.New()
	movieHandler := NewMovieHandler(tdb.Store)
	app.Get("/", movieHandler.HandleGetMovies)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	var respo ResourceResp
	if err := json.NewDecoder(resp.Body).Decode(&respo); err != nil {
		t.Fatal(err)
	}

	if respo.Results == 0 {
		t.Errorf("expecting a movie")
	}
	if respo.Data == nil {
		t.Errorf("expecting movies")
	}
}

func TestGetMovieByID(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	movieAdded := fixtures.AddMovie(tdb.Store, "The Matrix", []string{"Action"}, 120, 1999)
	app := fiber.New()
	movieHandler := NewMovieHandler(tdb.Store)
	app.Get("/:id", movieHandler.HandleGetMovieByID)
	req := httptest.NewRequest("GET", "/"+movieAdded.ID.Hex(), nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	var movieResp types.Movie
	json.NewDecoder(resp.Body).Decode(&movieResp)
	if movieResp.ID != movieAdded.ID {
		t.Errorf("expected movie id to be %s but got %s", movieAdded.ID, movieResp.ID)
	}
	if movieResp.Length != movieAdded.Length {
		t.Errorf("expected movie length to be %d but got %d", movieAdded.Length, movieResp.Length)
	}
	if movieResp.Year != movieAdded.Year {
		t.Errorf("expected movie year to be %d but got %d", movieAdded.Year, movieResp.Year)
	}
	if movieResp.Title != movieAdded.Title {
		t.Errorf("expected movie title to be %s but got %s", movieAdded.Title, movieResp.Title)
	}
	if movieResp.Genre[0] != movieAdded.Genre[0] {
		t.Errorf("expected movie genre to be %s but got %s", movieAdded.Genre, movieResp.Genre)
	}
}

func TestHandleUpdateMovieRating(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	movieHandler := NewMovieHandler(tdb.Store)
	app.Put("/:id/rate", movieHandler.HandleUpdateMovieRating)
	app.Get("/:id", movieHandler.HandleGetMovieByID)

	movieAdded := fixtures.AddMovie(tdb.Store, "The Matrix", []string{"Action"}, 120, 1999)

	type Rating struct {
		Rating int `json:"rating"`
	}
	rating := Rating{
		Rating: 5,
	}

	b, _ := json.Marshal(rating)
	req := httptest.NewRequest("PUT", "/"+movieAdded.ID.Hex()+"/rate", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status code 200 but got %d", resp.StatusCode)
	}
	req = httptest.NewRequest("GET", "/"+movieAdded.ID.Hex(), nil)
	resp, err = app.Test(req)
	if err != nil {
		t.Error(err)
	}
	var movieNewRating types.Movie
	json.NewDecoder(resp.Body).Decode(&movieNewRating)

	if movieNewRating.Rating != rating.Rating {
		t.Errorf("expected movie rating to be %d but got %d", rating.Rating, movieNewRating.Rating)
	}
}
