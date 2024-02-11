package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tomekzakrzewski/go-movierental/db"
	"github.com/tomekzakrzewski/go-movierental/types"
)

type MovieHandler struct {
	store db.MovieStore
}

func NewMovieHandler(store db.MovieStore) *MovieHandler {
	return &MovieHandler{
		store: store,
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
	var (
		rating  int
		movieID = c.Params("id")
	)
	if err := c.BodyParser(&rating); err != nil {
		return err
	}

	// TODO ERROR HANDLING
	if rating < 0 || rating > 10 {
		return c.Status(400).JSON(map[string]string{"error": "invalid rating"})
	}
	if err := h.store.UpdateRating(c.Context(), movieID, rating); err != nil {
		return err
	}

	return c.JSON(map[string]string{"updated": movieID})
}
