package handler

import (
	"github.com/Hilaladiii/aureus/internal/model"
	"github.com/Hilaladiii/aureus/internal/usecase"
	"github.com/Hilaladiii/aureus/pkg/response"
	"github.com/Hilaladiii/aureus/pkg/util"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type CategoryHandler struct {
	categoryUC usecase.CategoryUsecaseItf
	validator  *validator.Validate
}

func NewCategoryHandler(categoryUC usecase.CategoryUsecaseItf, validator *validator.Validate) *CategoryHandler {
	return &CategoryHandler{categoryUC, validator}
}

func (h *CategoryHandler) CreateCategory(c fiber.Ctx) error {
	var categoryPayload model.CategoryCreateRequest

	if err := c.Bind().Body(&categoryPayload); err != nil {
		return err
	}

	if err := h.validator.Struct(&categoryPayload); err != nil {
		return err
	}

	category, err := h.categoryUC.CreateCategory(c.Context(), &categoryPayload)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessResponse(fiber.StatusCreated, category))
}

func (h *CategoryHandler) UpdateCategory(c fiber.Ctx) error {
	categoryID := c.Params("categoryId")
	if err := util.ValidateUUID(categoryID); err != nil {
		return err
	}

	var categoryPayload model.CategoryUpdateRequest

	if err := c.Bind().Body(&categoryPayload); err != nil {
		return err
	}

	if err := h.validator.Struct(&categoryPayload); err != nil {
		return err
	}

	category, err := h.categoryUC.UpdateCategory(c.Context(), &categoryPayload, categoryID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(fiber.StatusOK, category))
}

func (h *CategoryHandler) DeleteCategory(c fiber.Ctx) error {
	categoryID := c.Params("categoryId")
	if err := util.ValidateUUID(categoryID); err != nil {
		return err
	}

	category, err := h.categoryUC.DeleteCategory(c.Context(), categoryID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(fiber.StatusOK, category))
}

func (h *CategoryHandler) GetAll(c fiber.Ctx) error {
	categories, err := h.categoryUC.GetAll(c.Context())
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(fiber.StatusOK, categories))
}

func (h *CategoryHandler) GetByID(c fiber.Ctx) error {
	categoryID := c.Params("categoryId")
	if err := util.ValidateUUID(categoryID); err != nil {
		return err
	}

	category, err := h.categoryUC.GetByID(c.Context(), categoryID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(fiber.StatusOK, category))
}
