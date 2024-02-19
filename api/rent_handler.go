package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tomekzakrzewski/go-movierental/db"
)

type RentHandler struct {
	store db.RentStore
}

func NewRentHandler(store db.RentStore) *RentHandler {
	return &RentHandler{
		store: store,
	}
}

// @Summary		Get all rents(user id, movie id, from, to)
// @Description	Handle getting all rents made by users
// @Tags			admin
// @Produce		json
// @Router			/rents [get]
func (h *RentHandler) HandleGetRents(c *fiber.Ctx) error {
	rents, err := h.store.GetRents(c.Context())
	if err != nil {
		return ErrResourceNotFound("Rents")
	}
	return c.JSON(rents)
}
