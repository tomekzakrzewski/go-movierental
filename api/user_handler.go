package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tomekzakrzewski/go-movierental/db"
	"github.com/tomekzakrzewski/go-movierental/types"
)

type UserHandler struct {
	store db.UserStore
}

func NewUserHandler(store db.UserStore) *UserHandler {
	return &UserHandler{
		store: store,
	}
}

// @Summary		Get users
// @Description	Handle getting users
// @Tags			admin
// @Produce		json
// @Router			/users [get]
func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.store.GetUsers(c.Context())
	if err != nil {
		return ErrResourceNotFound("Users")
	}
	return c.JSON(users)
}

// @Summary		Post user
// @Description	Handle posting user
// @Tags			user
// @Produce		json
// @Router			/users [post]
func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return ErrBadRequest()
	}
	insertedUser, err := h.store.InsertUser(c.Context(), user)
	if err != nil {
		return ErrBadRequest()
	}
	return c.JSON(insertedUser)
}

// @Summary		Get user by id
// @Description	Handle getting user by id
// @Tags			admin
// @Produce		json
// @Router			/users/:id [get]
func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	var (
		id = c.Params("id")
	)
	users, err := h.store.GetUserByID(c.Context(), id)
	if err != nil {
		return ErrResourceNotFound("User")
	}
	return c.JSON(users)
}

// @Summary		Delete user by id
// @Description	Handle deleting user by id
// @Tags			admin
// @Produce		json
// @Router			/users/:id [delete]
func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	var (
		id = c.Params("id")
	)
	if err := h.store.DeleteUser(c.Context(), id); err != nil {
		return ErrResourceNotFound("User")
	}
	return c.JSON(map[string]string{"deleted": id})
}
