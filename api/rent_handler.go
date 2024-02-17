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

func (h *RentHandler) HandleGetRents(c *fiber.Ctx) error {
	rents, err := h.store.GetRents(c.Context())
	if err != nil {
		return ErrResourceNotFound("Rents")
	}
	return c.JSON(rents)
}
