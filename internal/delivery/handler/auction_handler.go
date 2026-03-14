package handler

import (
	"bufio"
	"context"
	"fmt"
	"time"

	"github.com/Hilaladiii/aureus/internal/model"
	"github.com/Hilaladiii/aureus/internal/usecase"
	"github.com/Hilaladiii/aureus/pkg/response"
	"github.com/Hilaladiii/aureus/pkg/util"
	"github.com/bytedance/sonic"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/valyala/fasthttp"
)

type AuctionHandler struct {
	auctionUc usecase.AuctionUsecaseItf
	validator *validator.Validate
}

func NewAuctionHandler(auctionUc usecase.AuctionUsecaseItf, validator *validator.Validate) *AuctionHandler {
	return &AuctionHandler{auctionUc, validator}
}

func (h *AuctionHandler) Create(c fiber.Ctx) error {
	userID, err := util.GetJwtClaimLocals(c)
	if err != nil {
		return err
	}

	var auctionPayload model.AuctionCreateRequest
	if err := c.Bind().Form(&auctionPayload); err != nil {
		return err
	}

	if err := h.validator.Struct(&auctionPayload); err != nil {
		return err
	}

	auction, err := h.auctionUc.Create(c.Context(), &auctionPayload, userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessResponse(fiber.StatusCreated, auction))
}

func (h *AuctionHandler) BidAuction(c fiber.Ctx) error {
	auctionID := c.Params("auctionID")
	if err := util.ValidateUUID(auctionID); err != nil {
		return err
	}

	userID, err := util.GetJwtClaimLocals(c)
	if err != nil {
		return err
	}

	var auctionPayload model.AuctionBidRequest
	if err := c.Bind().Body(&auctionPayload); err != nil {
		return err
	}

	if err := h.validator.Struct(&auctionPayload); err != nil {
		return err
	}

	auction, err := h.auctionUc.BidAuction(c.Context(), &auctionPayload, auctionID, userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(fiber.StatusOK, auction))
}

func (h *AuctionHandler) GetAll(c fiber.Ctx) error {
	auctions, err := h.auctionUc.GetAll(c.Context())
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(fiber.StatusOK, auctions))
}

func (h *AuctionHandler) GetByID(c fiber.Ctx) error {
	auctionID := c.Params("auctionId")
	if err := util.ValidateUUID(auctionID); err != nil {
		return err
	}

	auctions, err := h.auctionUc.GetByID(c.Context(), auctionID)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(fiber.StatusOK, auctions))
}

func (h *AuctionHandler) GetByAuctioneerID(c fiber.Ctx) error {
	auctionID := c.Params("auctionId")
	if err := util.ValidateUUID(auctionID); err != nil {
		return err
	}

	auctions, err := h.auctionUc.GetByAuctioneerID(c.Context(), auctionID)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(fiber.StatusOK, auctions))
}

func (h *AuctionHandler) StreamLeaderboard(c fiber.Ctx) error {
	auctionID := c.Params("auctionId")
	err := util.ValidateUUID(auctionID)
	if err != nil {
		return err
	}

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	c.Status(fiber.StatusOK).RequestCtx().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			topBidders, err := h.auctionUc.GetLeaderboard(context.Background(), auctionID)
			if err != nil || len(topBidders) == 0 {
				continue
			}

			jsonData, _ := sonic.Marshal(topBidders)
			eventString := fmt.Sprintf("data: %s\n\n", string(jsonData))
			if _, writeErr := w.WriteString(eventString); writeErr != nil {
				return
			}
			w.Flush()
		}
	}))

	return nil
}
