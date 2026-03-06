package handler

import (
	"github.com/Hilaladiii/aureus/internal/model"
	"github.com/Hilaladiii/aureus/internal/usecase"
	"github.com/Hilaladiii/aureus/pkg/response"
	"github.com/Hilaladiii/aureus/pkg/util"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	userUC    usecase.UserUsecaseItf
	validator *validator.Validate
}

func NewUserHandler(userUC usecase.UserUsecaseItf, validator *validator.Validate) *UserHandler {
	return &UserHandler{userUC, validator}
}

func (h *UserHandler) Register(c fiber.Ctx) error {
	var payloadUser model.UserRegisterRequest

	if err := c.Bind().Body(&payloadUser); err != nil {
		return err
	}

	if err := h.validator.Struct(&payloadUser); err != nil {
		return err
	}

	user, err := h.userUC.Register(c.Context(), &payloadUser)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessResponse(fiber.StatusCreated, user))
}

func (h *UserHandler) Login(c fiber.Ctx) error {
	var payloadUser model.UserLoginRequest

	if err := c.Bind().Body(&payloadUser); err != nil {
		return err
	}

	if err := h.validator.Struct(&payloadUser); err != nil {
		return err
	}

	token, err := h.userUC.Login(c.Context(), &payloadUser)
	if err != nil {
		if err.Error() == "invalid credentials" {
			return fiber.NewError(fiber.StatusBadRequest, "Email or Password is incorrect")
		}
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(fiber.StatusOK, token))
}

func (h *UserHandler) GetProfile(c fiber.Ctx) error {
	userID, err := util.GetJwtClaimLocals(c)
	if err != nil {
		return err
	}

	user, err := h.userUC.GetUserById(c.Context(), userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(fiber.StatusOK, user))
}

func (h *UserHandler) UpdateUser(c fiber.Ctx) error {
	var payloadUser model.UserUpdateRequest
	userID, err := util.GetJwtClaimLocals(c)
	if err != nil {
		return err
	}
	if err := c.Bind().Body(&payloadUser); err != nil {
		return err
	}

	if err := h.validator.Struct(&payloadUser); err != nil {
		return err
	}

	user, err := h.userUC.UpdateUser(c.Context(), &payloadUser, userID)

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(fiber.StatusOK, user))
}
