package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Hilaladiii/aureus/internal/model"
	"github.com/Hilaladiii/aureus/internal/repository"
	"github.com/Hilaladiii/aureus/pkg/config"
	"github.com/Hilaladiii/aureus/pkg/exception"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type AuctionUsecaseItf interface {
	Create(ctx context.Context, req *model.AuctionCreateRequest, userID string) (model.AuctionResource, error)
	Update(ctx context.Context, req *model.AuctionUpdateRequest, auctionID string) (model.AuctionResource, error)
	Delete(ctx context.Context, auctionID string) (model.AuctionResource, error)
	GetAll(ctx context.Context) ([]model.AuctionResource, error)
	GetByID(ctx context.Context, auctionID string) (model.AuctionResource, error)
	BidAuction(ctx context.Context, req *model.AuctionBidRequest, auctionID string, userID string) (model.AuctionResource, error)
	GetBidHistory(ctx context.Context, auctionID string) ([]model.BidHistoryResource, error)
	FinalizeAuction(ctx context.Context, auctionID string) error
}

type AuctionUsecase struct {
	auctionRepo repository.AuctionRepoItf
	walletRepo  repository.WalletRepoItf
	bidRepo     repository.BidRepoItf
	seaweedFS   config.SeaweedFSStorageItf
	txManager   config.TxManagerItf
}

func NewAuctionUsecase(
	auctionRepo repository.AuctionRepoItf,
	walletRepo repository.WalletRepoItf,
	bidRepo repository.BidRepoItf,
	seaweedFS config.SeaweedFSStorageItf,
	txManager config.TxManagerItf,
) *AuctionUsecase {
	return &AuctionUsecase{auctionRepo, walletRepo, bidRepo, seaweedFS, txManager}
}

func (u *AuctionUsecase) Create(ctx context.Context, req *model.AuctionCreateRequest, userID string) (model.AuctionResource, error) {
	var uploadedImages []model.AuctionImage
	for _, fileHeader := range req.Images {
		if fileHeader.Size > 5*1024*1024 {
			return model.AuctionResource{}, exception.NewBadRequestError("file to large")
		}

		file, err := fileHeader.Open()
		if err != nil {
			return model.AuctionResource{}, err
		}
		defer file.Close()
		fileName := fmt.Sprintf("auctions/%s-%s", uuid.New().String(), fileHeader.Filename)

		imageUrl, err := u.seaweedFS.UploadFile(ctx, "images", fileName, file, fileHeader.Size, fileHeader.Header.Get("Content-Type"))
		if err != nil {
			return model.AuctionResource{}, err
		}

		uploadedImages = append(uploadedImages, model.AuctionImage{
			ImageUrl: imageUrl,
		})
	}
	newAuction := model.AuctionItem{
		Name:         req.Name,
		Description:  req.Description,
		StartPrice:   req.StartPrice,
		BidIncrement: req.BidIncrement,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		CategoryID:   req.CategoryID,
		AuctionImage: uploadedImages,
		AuctioneerID: userID,
	}

	err := u.auctionRepo.Create(ctx, &newAuction)
	if err != nil {
		return model.AuctionResource{}, err
	}

	return newAuction.Resource(), nil
}

func (u *AuctionUsecase) Update(ctx context.Context, req *model.AuctionUpdateRequest, auctionID string) (model.AuctionResource, error) {
	auction, err := u.auctionRepo.GetByID(ctx, auctionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.AuctionResource{}, exception.NewBadRequestError("invalid auction id")
		}
		return model.AuctionResource{}, err
	}

	if req.Name != "" {
		auction.Name = req.Name
	}

	if req.Description != nil && *req.Description != "" {
		auction.Description = req.Description
	}

	err = u.auctionRepo.Update(ctx, auction)
	if err != nil {
		return model.AuctionResource{}, err
	}
	return auction.Resource(), nil
}

func (u *AuctionUsecase) Delete(ctx context.Context, auctionID string) (model.AuctionResource, error) {
	auction, err := u.auctionRepo.GetByID(ctx, auctionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.AuctionResource{}, exception.NewBadRequestError("invalid auction id")
		}
		return model.AuctionResource{}, err
	}

	err = u.auctionRepo.Delete(ctx, auctionID)
	if err != nil {
		return model.AuctionResource{}, err
	}

	return auction.Resource(), nil
}

func (u *AuctionUsecase) GetAll(ctx context.Context) ([]model.AuctionResource, error) {
	auctions, err := u.auctionRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return model.AuctionResources(auctions), nil
}

func (u *AuctionUsecase) GetByID(ctx context.Context, auctionID string) (model.AuctionResource, error) {
	auction, err := u.auctionRepo.GetByID(ctx, auctionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.AuctionResource{}, exception.NewBadRequestError("auction not found")
		}
		return model.AuctionResource{}, err
	}
	return auction.Resource(), nil
}

func (u *AuctionUsecase) BidAuction(ctx context.Context, req *model.AuctionBidRequest, auctionID string, userID string) (model.AuctionResource, error) {
	var finalAuction model.AuctionResource
	err := u.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// get auction data
		auction, err := u.auctionRepo.GetByIDWithLock(txCtx, auctionID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return exception.NewNotFoundError("auction not found")
			}
			return err
		}

		// check auction timeline is valid
		now := time.Now()
		if !(now.After(auction.StartTime) && now.Before(auction.EndTime)) {
			return exception.NewBadRequestError("auction time has ended")
		}

		// check balance new bidder
		newBidderWallet, err := u.walletRepo.GetByUserID(txCtx, userID)
		if err != nil {
			return err
		}

		if newBidderWallet.ActiveBalance.LessThan(req.CurrentPrice) || newBidderWallet.ActiveBalance.LessThan(auction.StartPrice) {
			return exception.NewBadRequestError("insufficient balance")
		}

		// check is the requirement minimum bid valid
		var minimumBid decimal.Decimal
		if auction.CurrentPrice.IsZero() {
			minimumBid = auction.StartPrice
		} else {
			minimumBid = auction.BidIncrement.Add(auction.CurrentPrice)
		}

		isBidValid := minimumBid.LessThanOrEqual(req.CurrentPrice)

		if !isBidValid {
			return exception.NewBadRequestError(fmt.Sprintf("the bid is to low, the next minimum bid is %s", minimumBid.String()))
		}

		// update balance new bidder
		if auction.CurrentPrice.IsZero() {
			newBidderWallet.ActiveBalance = newBidderWallet.ActiveBalance.Sub(auction.StartPrice)
			newBidderWallet.HeldBalance = newBidderWallet.HeldBalance.Add(auction.StartPrice)
		} else {
			newBidderWallet.ActiveBalance = newBidderWallet.ActiveBalance.Sub(auction.CurrentPrice)
			newBidderWallet.HeldBalance = newBidderWallet.HeldBalance.Add(auction.CurrentPrice)
		}

		err = u.walletRepo.Update(txCtx, newBidderWallet, newBidderWallet.ID)
		if err != nil {
			return err
		}

		previousBidder, err := u.bidRepo.GetHighestBid(txCtx, auctionID)
		// update previous bidder wallet
		if err == nil && previousBidder != nil {
			if previousBidder.UserID == userID {
				return exception.NewBadRequestError("you are already the highest bidder at this time!")
			}
			previousUserID := previousBidder.UserID
			previousBidderWallet, err := u.walletRepo.GetByUserIDWithLock(txCtx, previousUserID)
			if err != nil {
				return err
			}

			refundAmount := previousBidder.BidAmount
			previousBidderWallet.HeldBalance = previousBidderWallet.HeldBalance.Sub(refundAmount)
			previousBidderWallet.ActiveBalance = previousBidderWallet.ActiveBalance.Add(refundAmount)

			err = u.walletRepo.Update(txCtx, previousBidderWallet, previousBidderWallet.ID)
			if err != nil {
				return err
			}

			// update bid history for previous bidder
			previousBidder.Status = "OUTBID"
			err = u.bidRepo.Update(txCtx, previousBidder, previousBidder.ID)
			if err != nil {
				return err
			}
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// new highest bid history
		newBid := model.BidHistory{
			BidAmount: req.CurrentPrice,
			UserID:    userID,
			Status:    "WINNING",
			ItemID:    auctionID,
		}
		err = u.bidRepo.Create(txCtx, &newBid)
		if err != nil {
			return err
		}

		// update auction current price
		auction.CurrentPrice = req.CurrentPrice
		err = u.auctionRepo.Update(txCtx, auction)
		if err != nil {
			return err
		}
		finalAuction = auction.Resource()
		return nil
	})
	if err != nil {
		return model.AuctionResource{}, err
	}

	return finalAuction, nil
}

func (u *AuctionUsecase) GetBidHistory(ctx context.Context, auctionID string) ([]model.BidHistoryResource, error) {
	histories, err := u.bidRepo.GetByAuctionID(ctx, auctionID)
	if err != nil {
		return nil, err
	}

	return model.BidResources(histories), nil
}

func (u *AuctionUsecase) FinalizeAuction(ctx context.Context, auctionID string) error {
	err := u.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		auction, err := u.auctionRepo.GetByIDWithLock(txCtx, auctionID)
		if err != nil {
			return err
		}

		highestBidder, err := u.bidRepo.GetHighestBid(txCtx, auctionID)
		if err != nil {
			return err
		}

		bidderWallet, err := u.walletRepo.GetByUserIDWithLock(txCtx, highestBidder.UserID)
		if err != nil {
			return err
		}

		bidderWallet.HeldBalance = bidderWallet.HeldBalance.Sub(auction.CurrentPrice)
		err = u.walletRepo.Update(txCtx, bidderWallet, bidderWallet.ID)
		if err != nil {
			return err
		}

		auctioneerWallet, err := u.walletRepo.GetByUserIDWithLock(txCtx, auction.AuctioneerID)
		if err != nil {
			return err
		}

		auctioneerWallet.ActiveBalance = auctioneerWallet.ActiveBalance.Add(auction.CurrentPrice)
		err = u.walletRepo.Update(txCtx, auctioneerWallet, auctioneerWallet.ID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
