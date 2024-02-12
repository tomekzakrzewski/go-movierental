package api

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tomekzakrzewski/go-movierental/db"
	"github.com/tomekzakrzewski/go-movierental/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MovieHandler struct {
	store     db.MovieStore
	rentStore db.RentStore
}

func NewMovieHandler(store db.MovieStore, rentStore db.RentStore) *MovieHandler {
	return &MovieHandler{
		store:     store,
		rentStore: rentStore,
	}
}

func (h *MovieHandler) HandlePostMovie(c *fiber.Ctx) error {
	var params types.CreateMovieParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	validate := types.Validate(params)
	if len(validate) > 0 {
		return c.JSON(validate)
	}
	movie := types.NewMovieFromParams(params)
	insertedHotel, err := h.store.InsertMovie(c.Context(), movie)
	if err != nil {
		return err
	}
	return c.JSON(insertedHotel)
}

func (h *MovieHandler) HandleGetMovie(c *fiber.Ctx) error {
	movies, err := h.store.GetMovies(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(movies)
}

func (h *MovieHandler) HandleUpdateMovie(c *fiber.Ctx) error {
	var (
		params  types.UpdateMovieParams
		movieID = c.Params("id")
	)
	// TODO ERROR HANDLING
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if err := h.store.PutMovie(c.Context(), movieID, params); err != nil {
		return err
	}

	return c.JSON(map[string]string{"updated": movieID})
}

func (h *MovieHandler) HandleDeleteMovie(c *fiber.Ctx) error {
	movieID := c.Params("id")
	if err := h.store.DeleteMovie(c.Context(), movieID); err != nil {
		return err
	}
	return c.JSON(map[string]string{"deleted": movieID})
}

func (h *MovieHandler) HandleGetMovieByID(c *fiber.Ctx) error {
	movieID := c.Params("id")
	movie, err := h.store.GetMovieByID(c.Context(), movieID)
	if err != nil {
		return err
	}
	return c.JSON(movie)
}

func (h *MovieHandler) HandleUpdateMovieRating(c *fiber.Ctx) error {
	type Rating struct {
		Rating int `json:"rating"`
	}
	var (
		movieID = c.Params("id")
		rating  Rating
	)
	if err := c.BodyParser(&rating); err != nil {
		return err
	}

	// TODO ERROR HANDLING
	if rating.Rating < 0 || rating.Rating > 10 {
		return c.Status(400).JSON(map[string]string{"error": "invalid rating"})
	}
	if err := h.store.UpdateRating(c.Context(), movieID, rating.Rating); err != nil {
		return err
	}

	return c.JSON(map[string]string{"updated": movieID})
}

func (h *MovieHandler) HandleRentMovie(c *fiber.Ctx) error {
	// TODO ERROR HANDLING
	movieID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		fmt.Println("here")
		return c.Status(401).JSON(map[string]string{"error": "unauthorized"})
	}

	if err := h.rentStore.CheckRent(c.Context(), types.CheckRentParams{
		UserID:  user.ID,
		MovieID: movieID,
		From:    time.Now(),
		To:      time.Now().Add(time.Hour * 24),
	}); err != nil {
		return err
	}
	params := types.CreateRentParams{
		UserID:  user.ID,
		MovieID: movieID,
	}
	rent := types.NewRentFromParams(params)
	insertedRent, err := h.rentStore.InsertRent(c.Context(), rent)
	if err != nil {
		return err
	}
	return c.JSON(insertedRent)
}

func (h *MovieHandler) HandleGetRentedMovies(c *fiber.Ctx) error {
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(401).JSON(map[string]string{"error": "unauthorized"})
	}

	movies, err := h.rentStore.GetRentsByUser(c.Context(), user.ID.Hex())
	if err != nil {
		return err
	}
	return c.JSON(movies)
}
