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

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.User)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		Username:  "tomek_test",
		Email:     "tomek@test.com",
		FirstName: "tomek",
		LastName:  "test",
		Password:  "tomektestpass",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)
	if len(user.ID) == 0 {
		t.Errorf("expecting a user id to be set")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Errorf("expecting the EncryptedPassword not to be included in the json response")
	}
	if user.FirstName != params.FirstName {
		t.Errorf("expected firstname %s but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected last name %s but got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected email %s but got %s", params.Email, user.Email)
	}
}

func TestGetUsers(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	fixtures.AddUser(tdb.Store, "tomek", "test", false)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.User)
	app.Get("/", userHandler.HandleGetUsers)

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	var users []types.User
	json.NewDecoder(resp.Body).Decode(&users)
	if len(users) != 1 {
		t.Errorf("expected 1 user but got %d", len(users))
	}
}

// fix
func TestDeleteUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.User)
	app.Post("/", userHandler.HandlePostUser)
	app.Delete("/:id", userHandler.HandleDeleteUser)

	params := types.CreateUserParams{
		Username:  "tomek_test",
		Email:     "tomek@test.com",
		FirstName: "tomek",
		LastName:  "test",
		Password:  "tomektestpass",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))

	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)

	app.Delete("/:id", userHandler.HandleDeleteUser)
	req = httptest.NewRequest("DELETE", "/"+user.ID.Hex(), nil)
	_, err = app.Test(req)
	if err != nil {
		t.Error(err)
	}
}
