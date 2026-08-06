package handler

import (
	"Merchant-POS-backend/module/user/dto"
	"Merchant-POS-backend/module/user/ports"
	"Merchant-POS-backend/pkg/response"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type userHandler struct {
	service ports.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(service ports.UserService) *userHandler {
	return &userHandler{service: service}
}

// CreateUser godoc
// @Summary      Create a new user
// @Description  Create a new user with username, email, and password
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreateUserReq  true  "Create User Request"
// @Success      201      {object}  response.SuccessResponse{data=dto.UserResponse}
// @Failure      400      {object}  response.ErrorResponse
// @Router       /users [post]
func (h *userHandler) CreateUser(c *fiber.Ctx) error {
	var req dto.CreateUserReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(fiber.StatusBadRequest, "Invalid request body", nil))
	}

	user, err := h.service.CreateUser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(fiber.StatusBadRequest, err.Error(), nil))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(fiber.StatusCreated, "User created successfully", user))
}

// GetUser godoc
// @Summary      Get a user by ID
// @Description  Get detailed information of a user by their ID
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  response.SuccessResponse{data=dto.UserResponse}
// @Failure      404  {object}  response.ErrorResponse
// @Router       /users/{id} [get]
func (h *userHandler) GetUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(fiber.StatusBadRequest, "Invalid user ID", nil))
	}

	user, err := h.service.GetUser(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.Error(fiber.StatusNotFound, err.Error(), nil))
	}

	return c.JSON(response.Success(fiber.StatusOK, "User retrieved successfully", user))
}

// GetAllUsers godoc
// @Summary      Get all users
// @Description  List all users in the system
// @Tags         users
// @Produce      json
// @Success      200  {object}  response.SuccessResponse{data=[]dto.UserResponse}
// @Router       /users [get]
func (h *userHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.service.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(fiber.StatusInternalServerError, err.Error(), nil))
	}

	return c.JSON(response.Success(fiber.StatusOK, "Users retrieved successfully", users))
}

// UpdateUser godoc
// @Summary      Update a user
// @Description  Update a user's details
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id       path      int                true  "User ID"
// @Param        request  body      dto.UpdateUserReq  true  "Update User Request"
// @Success      200      {object}  response.SuccessResponse{data=dto.UserResponse}
// @Failure      400      {object}  response.ErrorResponse
// @Failure      404      {object}  response.ErrorResponse
// @Router       /users/{id} [put]
func (h *userHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(fiber.StatusBadRequest, "Invalid user ID", nil))
	}

	var req dto.UpdateUserReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(fiber.StatusBadRequest, "Invalid request body", nil))
	}

	user, err := h.service.UpdateUser(uint(id), &req)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(response.Error(fiber.StatusNotFound, err.Error(), nil))
		}
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(fiber.StatusBadRequest, err.Error(), nil))
	}

	return c.JSON(response.Success(fiber.StatusOK, "User updated successfully", user))
}

// DeleteUser godoc
// @Summary      Delete a user
// @Description  Delete a user from the system
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  response.SuccessResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Router       /users/{id} [delete]
func (h *userHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(fiber.StatusBadRequest, "Invalid user ID", nil))
	}

	if err := h.service.DeleteUser(uint(id)); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.Error(fiber.StatusNotFound, err.Error(), nil))
	}

	return c.JSON(response.Success(fiber.StatusOK, "User deleted successfully", nil))
}
