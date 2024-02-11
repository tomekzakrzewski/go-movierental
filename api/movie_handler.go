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
	insertedHotel, err := h.store.InsertHotel(c.Context(), movie)
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
