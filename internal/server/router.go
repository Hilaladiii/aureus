package server

import (
	"github.com/Hilaladiii/aureus/internal/delivery/handler"
	"github.com/Hilaladiii/aureus/internal/delivery/middleware"
	"github.com/Hilaladiii/aureus/internal/model"

	"github.com/gofiber/fiber/v3"
)

type Router struct {
	UserHandler     *handler.UserHandler
	CategoryHandler *handler.CategoryHandler
	WalletHandler   *handler.WalletHandler
	AuctionHandler  *handler.AuctionHandler
	Middleware      middleware.MiddlewareItf
}

func NewRouter(
	userHandler *handler.UserHandler,
	categoryHandler *handler.CategoryHandler,
	wallethandler *handler.WalletHandler,
	auctionHandler *handler.AuctionHandler,
	middleware *middleware.Middleware,
) *Router {
	return &Router{
		UserHandler:     userHandler,
		CategoryHandler: categoryHandler,
		WalletHandler:   wallethandler,
		AuctionHandler:  auctionHandler,
		Middleware:      middleware,
	}
}

func (r *Router) Setup(app *fiber.App) {
	api := app.Group("/api/v1")

	// stripe webhook
	app.Post("/webhooks/stripe", r.WalletHandler.StripeWeebHook)

	// users route
	user := api.Group("/users")
	user.Post("/register", r.UserHandler.Register)
	user.Post("/login", r.UserHandler.Login)
	user.Use(r.Middleware.JwtMiddleware())
	user.Get("/me", r.UserHandler.GetProfile)

	// categories route
	category := api.Group("/categories")
	category.Use(r.Middleware.JwtMiddleware())
	category.Post("", r.CategoryHandler.CreateCategory)
	category.Put("/:categoryId", r.CategoryHandler.UpdateCategory)
	category.Delete("/:categoryId", r.CategoryHandler.DeleteCategory)
	category.Get("", r.CategoryHandler.GetAll)
	category.Get("/:categoryId", r.CategoryHandler.GetByID)

	// wallets route
	wallet := api.Group("/wallets")
	wallet.Use(r.Middleware.JwtMiddleware())
	wallet.Post("/top-up", r.WalletHandler.CreateTopUpSession)
	wallet.Post("", r.WalletHandler.Create)
	wallet.Post("/:walletId", r.WalletHandler.GetCurrentBalance)

	// auction route
	auction := api.Group("/auctions")
	auction.Use(r.Middleware.JwtMiddleware())
	auction.Post("", r.Middleware.RoleMiddleware(model.SELLER), r.AuctionHandler.Create)
	auction.Get("", r.AuctionHandler.GetAll)
	auction.Get("/:auctionId", r.AuctionHandler.GetByID)
	auction.Post("/:auctionId/bid", r.Middleware.RoleMiddleware(model.BIDDER), r.AuctionHandler.BidAuction)
	auction.Get("/:auctionId/leaderboard/stream", r.AuctionHandler.StreamLeaderboard)
}
