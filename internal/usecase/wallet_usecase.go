package usecase

import (
	"context"
	"errors"

	"github.com/Hilaladiii/aureus/internal/model"
	"github.com/Hilaladiii/aureus/internal/repository"
	"github.com/Hilaladiii/aureus/pkg/config"
	"github.com/Hilaladiii/aureus/pkg/exception"
	"github.com/stripe/stripe-go/v84"
	"gorm.io/gorm"
)

type WalletUsecaseItf interface {
	Create(ctx context.Context, req *model.WalletCreateRequest, userID string) (model.WalletResource, error)
	CreateTopUpSession(ctx context.Context, amount float64, userID string) (string, error)
	TopUpBalance(ctx context.Context, amount float64, walletID string) error
	GetCurrentBalance(ctx context.Context, walletID string) (model.WalletResource, error)
	GetByID(ctx context.Context, walletID string) (model.WalletResource, error)
}

type WalletUsecase struct {
	walletRepo repository.WalletRepoItf
	env        config.Env
}

func NewWalletUsecase(walletRepo repository.WalletRepoItf, env config.Env) *WalletUsecase {
	return &WalletUsecase{
		walletRepo: walletRepo,
		env:        env,
	}
}

func (u *WalletUsecase) Create(ctx context.Context, req *model.WalletCreateRequest, userID string) (model.WalletResource, error) {
	wallet, err := u.walletRepo.GetByUserID(ctx, userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return model.WalletResource{}, err
	}
	if wallet != nil {
		return model.WalletResource{}, exception.NewBadRequestError("wallet for this user already exists")
	}

	newWallet := model.Wallet{
		ActiveBalance: req.ActiveBalance,
		UserID:        userID,
	}
	err = u.walletRepo.Create(ctx, &newWallet)
	if err != nil {
		return model.WalletResource{}, nil
	}

	return newWallet.Resource(), nil
}

func (u *WalletUsecase) CreateTopUpSession(ctx context.Context, amount float64, userID string) (string, error) {
	wallet, err := u.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", exception.NewNotFoundError("wallet not found")
		}
		return "", nil
	}

	sc := stripe.NewClient(u.env.StripeSecretKey)

	params := &stripe.CheckoutSessionCreateParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		Mode:               stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:         stripe.String(u.env.StripeSuccessUrl),
		CancelURL:          stripe.String(u.env.StripeCancelUrl),
		ClientReferenceID:  stripe.String(wallet.ID),
		LineItems: []*stripe.CheckoutSessionCreateLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionCreateLineItemPriceDataParams{
					Currency: stripe.String("IDR"),
					ProductData: &stripe.CheckoutSessionCreateLineItemPriceDataProductDataParams{
						Name: stripe.String("Top up wallet balance"),
					},
					UnitAmount: stripe.Int64(int64(amount) * 100),
				},
				Quantity: stripe.Int64(1),
			},
		},
	}

	s, err := sc.V1CheckoutSessions.Create(ctx, params)
	if err != nil {
		return "", err
	}
	return s.URL, nil
}

func (u *WalletUsecase) TopUpBalance(ctx context.Context, amount float64, walletID string) error {
	_, err := u.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exception.NewNotFoundError("wallet not found")
		}
		return err
	}

	err = u.walletRepo.AddBalance(ctx, amount, walletID)
	if err != nil {
		return err
	}

	return nil
}

func (u *WalletUsecase) GetCurrentBalance(ctx context.Context, walletID string) (model.WalletResource, error) {
	wallet, err := u.walletRepo.GetBalance(ctx, walletID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.WalletResource{}, exception.NewNotFoundError("wallet not found")
		}
		return model.WalletResource{}, err
	}

	return wallet.Resource(), nil
}

func (u *WalletUsecase) GetByID(ctx context.Context, walletID string) (model.WalletResource, error) {
	wallet, err := u.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.WalletResource{}, exception.NewNotFoundError("wallet not found")
		}
		return model.WalletResource{}, err
	}

	return wallet.Resource(), nil
}
