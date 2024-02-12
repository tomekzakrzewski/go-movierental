package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tomekzakrzewski/go-movierental/db"
	"github.com/tomekzakrzewski/go-movierental/types"
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
		return err
	}
	return c.JSON(rents)
}

func (h *RentHandler) HandlePostRent(c *fiber.Ctx) error {
	// TODO ERROR HANDLING
	var params types.CreateRentParams

	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}
	rent := types.NewRentFromParams(params)
	insertedRent, err := h.store.InsertRent(c.Context(), rent)
	if err != nil {
		return err
	}
	return c.JSON(insertedRent)
}
