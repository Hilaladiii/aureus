package handler

import (
	"github.com/Hilaladiii/aureus/internal/model"
	"github.com/Hilaladiii/aureus/internal/usecase"
	"github.com/Hilaladiii/aureus/pkg/config"
	"github.com/Hilaladiii/aureus/pkg/response"
	"github.com/Hilaladiii/aureus/pkg/util"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type WalletHandler struct {
	walletUc  usecase.WalletUsecaseItf
	validator *validator.Validate
	env       config.Env
}

func NewWalletHandler(walletUc usecase.WalletUsecaseItf, validator *validator.Validate, env config.Env) *WalletHandler {
	return &WalletHandler{
		walletUc:  walletUc,
		validator: validator,
		env:       env,
	}
}

func (h *WalletHandler) CreateTopUpSession(c fiber.Ctx) error {
	var walletPayload model.WalletCreateRequest
	userID, err := util.GetJwtClaimLocals(c)
	if err != nil {
		return err
	}
	if err := c.Bind().Body(&walletPayload); err != nil {
		return err
	}
	if err := h.validator.Struct(&walletPayload); err != nil {
		return err
	}

	url, err := h.walletUc.CreateTopUpSession(c.Context(), walletPayload.ActiveBalance.Abs(), userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(fiber.StatusOK, fiber.Map{
		"url": url,
	}))
}

func (h *WalletHandler) Create(c fiber.Ctx) error {
	userID, err := util.GetJwtClaimLocals(c)
	if err != nil {
		return err
	}

	var walletPayload model.WalletCreateRequest

	if err := c.Bind().Body(&walletPayload); err != nil {
		return err
	}

	if err := h.validator.Struct(&walletPayload); err != nil {
		return err
	}

	wallet, err := h.walletUc.Create(c.Context(), &walletPayload, userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(fiber.StatusOK, wallet))
}

func (h *WalletHandler) StripeWeebHook(c fiber.Ctx) error {
	payload := c.Body()
	signature := c.Get("Stripe-Signature")

	err := h.walletUc.StripeWebHook(c.Context(), payload, signature)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusOK)
}

func (h *WalletHandler) GetCurrentBalance(c fiber.Ctx) error {
	walletID := c.Params("walletId")
	if err := util.ValidateUUID(walletID); err != nil {
		return err
	}

	wallet, err := h.walletUc.GetCurrentBalance(c.Context(), walletID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(fiber.StatusOK, wallet))
}
