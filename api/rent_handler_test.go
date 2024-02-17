package api

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/tomekzakrzewski/go-movierental/db/fixtures"
	"github.com/tomekzakrzewski/go-movierental/types"
)

func TestGetRents(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	var (
		rentHandler  = NewRentHandler(tdb.Rent)
		movieHandler = NewMovieHandler(tdb.Store)
		movieAdded   = fixtures.AddMovie(tdb.Store, "The Matrix", []string{"Action"}, 120, 1999)
		userAdded    = fixtures.AddUser(tdb.Store, "tomek", "test", false)
		app          = fiber.New()
		apiv1        = app.Group("", JWTAuthentication(tdb.User))
	)
	apiv1.Get("/", rentHandler.HandleGetRents)
	apiv1.Put("/:id/rent", movieHandler.HandleRentMovie)
	token := CreateTokenFromUser(userAdded)
	req := httptest.NewRequest("PUT", "/"+movieAdded.ID.Hex()+"/rent", nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Api-Token", token)
	_, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Api-Token", token)
	resp, err := app.Test(req)

	var rents []types.Rent
	json.NewDecoder(resp.Body).Decode(&rents)
	if len(rents) != 1 {
		t.Errorf("expected 1 rent but got %d", len(rents))
	}
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status code 200 but got %d", resp.StatusCode)
	}
}
