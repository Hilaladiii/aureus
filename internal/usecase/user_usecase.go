package usecase

import (
	"context"
	"errors"

	"github.com/Hilaladiii/aureus/internal/model"
	"github.com/Hilaladiii/aureus/internal/repository"
	"github.com/Hilaladiii/aureus/pkg/exception"
	"github.com/Hilaladiii/aureus/pkg/jwt"
	"github.com/Hilaladiii/aureus/pkg/util"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecaseItf interface {
	Register(ctx context.Context, req *model.UserRegisterRequest) (model.UserResource, error)
	Login(ctx context.Context, req *model.UserLoginRequest) (string, error)
	GetUserById(ctx context.Context, id string) (model.UserResource, error)
	UpdateUser(ctx context.Context, req *model.UserUpdateRequest, userId string) (model.UserResource, error)
}

type UserUsecase struct {
	userRepo repository.UserRepoItf
	jwt      jwt.JwtItf
}

func NewUserUsecase(ur repository.UserRepoItf, jwt jwt.JwtItf) *UserUsecase {
	return &UserUsecase{
		userRepo: ur,
		jwt:      jwt,
	}
}

// GetUserById implements [UserUsecaseItf].
func (u *UserUsecase) GetUserById(ctx context.Context, id string) (model.UserResource, error) {
	user, err := u.userRepo.GetUserById(ctx, id)
	if err != nil {
		return model.UserResource{}, err
	}

	return user.Resource(), nil
}

// Login implements [UserUsecaseItf].
func (u *UserUsecase) Login(ctx context.Context, req *model.UserLoginRequest) (string, error) {
	user, err := u.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := u.jwt.CreateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Register implements [UserUsecaseItf].
func (u *UserUsecase) Register(ctx context.Context, req *model.UserRegisterRequest) (model.UserResource, error) {
	exist, err := u.userRepo.CheckUserExist(ctx, req.Email)
	if err != nil {
		return model.UserResource{}, err
	}

	if exist {
		return model.UserResource{}, exception.NewBadRequestError("email already exists")
	}

	if err := util.ValidateStrongPassword(req.Password); err != nil {
		return model.UserResource{}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.UserResource{}, err
	}

	newUser := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	err = u.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		return model.UserResource{}, err
	}

	return newUser.Resource(), nil
}

// UpdateUser implements [UserUsecaseItf].
func (u *UserUsecase) UpdateUser(ctx context.Context, req *model.UserUpdateRequest, userId string) (model.UserResource, error) {
	user, err := u.userRepo.GetUserById(ctx, userId)
	if err != nil {
		return model.UserResource{}, err
	}

	if req.Email != "" {
		user.Email = req.Email
	}

	if req.Username != "" {
		user.Username = req.Username
	}

	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return model.UserResource{}, err
		}

		user.Password = string(hashedPassword)
	}

	if err = u.userRepo.UpdateUser(ctx, user); err != nil {
		return model.UserResource{}, err
	}

	return user.Resource(), nil
}
