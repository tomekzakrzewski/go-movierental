package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tomekzakrzewski/go-movierental/db"
	"github.com/tomekzakrzewski/go-movierental/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MovieHandler struct {
	store *db.Store
}

func NewMovieHandler(store *db.Store) *MovieHandler {
	return &MovieHandler{
		store: store,
	}
}

//	@Summary		Add movie
//	@Description	Handle posting movie to database
//	@Tags			admin
//	@Accept			json
//	@Produce		json
//	@Router			/movies [post]
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
	insertedHotel, err := h.store.Movie.InsertMovie(c.Context(), movie)
	if err != nil {
		return err
	}
	return c.JSON(insertedHotel)
}

type ResourceResp struct {
	Results int `json:"results"`
	Data    any `json:"data"`
	Page    int `json:"page"`
}

type MovieQueryParams struct {
	db.Pagination
	Rating int
}

//	@Summary		Get all movies
//	@Description	Handle getting all movies from database
//	@Tags			user
//	@Produce		json
//	@Router			/movies [get]
func (h *MovieHandler) HandleGetMovies(c *fiber.Ctx) error {
	var params MovieQueryParams
	if err := c.QueryParser(&params); err != nil {
		return ErrBadRequest()
	}
	filter := map[string]any{
		"rating": params.Rating,
	}
	movies, err := h.store.Movie.GetMovies(c.Context(), filter, &params.Pagination)
	if err != nil {
		return ErrResourceNotFound("Movies")
	}

	res := ResourceResp{
		Results: len(movies),
		Data:    movies,
		Page:    params.Page,
	}
	return c.JSON(res)
}

//	@Summary		Update movie
//	@Description	Handle updating movie
//	@Tags			admin
//	@Produce		json
//	@Router			/movies/:id [put]
func (h *MovieHandler) HandleUpdateMovie(c *fiber.Ctx) error {
	var (
		params  types.UpdateMovieParams
		movieID = c.Params("id")
	)
	if err := c.BodyParser(&params); err != nil {
		return ErrInvalidID()
	}
	if err := h.store.Movie.PutMovie(c.Context(), movieID, params); err != nil {
		return ErrResourceNotFound("Movie")
	}

	return c.JSON(map[string]string{"updated": movieID})
}

//	@Summary		Delete movie
//	@Description	Handle deleting movie
//	@Tags			admin
//	@Produce		json
//	@Router			/movies/:id [delete]
func (h *MovieHandler) HandleDeleteMovie(c *fiber.Ctx) error {
	movieID := c.Params("id")
	if err := h.store.Movie.DeleteMovie(c.Context(), movieID); err != nil {
		return ErrInvalidID()
	}
	return c.JSON(map[string]string{"deleted": movieID})
}

//	@Summary		Get movie by ID
//	@Description	Handle getting movie by id
//	@Tags			user
//	@Produce		json
//	@Router			/movies/:id [get]
func (h *MovieHandler) HandleGetMovieByID(c *fiber.Ctx) error {
	movieID := c.Params("id")
	movie, err := h.store.Movie.GetMovieByID(c.Context(), movieID)
	if err != nil {
		return ErrInvalidID()
	}
	return c.JSON(movie)
}

//	@Summary		Update movie movie rating
//	@Description	Handle updating movie rating
//	@Tags			user
//	@Produce		json
//	@Router			/movies/:id/rate [put]
func (h *MovieHandler) HandleUpdateMovieRating(c *fiber.Ctx) error {
	type Rating struct {
		Rating int `json:"rating"`
	}

	var (
		movieID = c.Params("id")
		rating  Rating
	)
	if err := c.BodyParser(&rating); err != nil {
		return ErrInvalidID()
	}

	if rating.Rating < 0 || rating.Rating > 10 {
		return ErrBadRequest()
	}
	if err := h.store.Movie.UpdateRating(c.Context(), movieID, rating.Rating); err != nil {
		return ErrResourceNotFound("Movie")
	}

	return c.JSON(map[string]string{"updated": movieID})
}

//	@Summary		Rent a movie
//	@Description	Handle renting movie
//	@Tags			user
//	@Produce		json
//	@Router			/movies/:id/rent [post]
func (h *MovieHandler) HandleRentMovie(c *fiber.Ctx) error {
	movieID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return ErrInvalidID()
	}
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return ErrUnAuthorized()
	}

	if err := h.store.Rent.CheckRent(c.Context(), types.CheckRentParams{
		UserID:  user.ID,
		MovieID: movieID,
		From:    time.Now(),
		To:      time.Now().Add(time.Hour * 24),
	}); err != nil {
		return c.Status(http.StatusBadRequest).JSON(genericResp{
			Type: "error",
			Msg:  fmt.Sprintf("Movie already rented, id: %s", movieID),
		})
	}
	params := types.CreateRentParams{
		UserID:  user.ID,
		MovieID: movieID,
	}
	rent := types.NewRentFromParams(params)
	insertedRent, err := h.store.Rent.InsertRent(c.Context(), rent)
	if err != nil {
		return ErrBadRequest()
	}
	return c.JSON(insertedRent)
}

//	@Summary		Get movies rented by user
//	@Description	Handle getting movies rented by user
//	@Tags			user
//	@Produce		json
//	@Router			/movies/rented [post]
func (h *MovieHandler) HandleGetRentedMovies(c *fiber.Ctx) error {
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return ErrBadRequest()
	}

	movies, err := h.store.Rent.GetRentsByUser(c.Context(), user.ID.Hex())
	if err != nil {
		return ErrResourceNotFound("rented movies")
	}
	return c.JSON(movies)
}
