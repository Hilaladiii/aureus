package repository

import (
	"context"
	"errors"

	"github.com/Hilaladiii/aureus/internal/model"

	"gorm.io/gorm"
)

type UserRepoItf interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserById(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	CheckUserExist(ctx context.Context, email string) (bool, error)
}

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

// CheckUserExist implements [UserRepoItf].
func (ur *UserRepo) CheckUserExist(ctx context.Context, email string) (bool, error) {
	var count int64

	err := ur.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// CreateUser implements [UserRepoItf].
func (ur *UserRepo) CreateUser(ctx context.Context, user *model.User) error {
	err := ur.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return err
	}

	return nil
}

// GetUserById implements [UserRepoItf].
func (ur *UserRepo) GetUserById(ctx context.Context, id string) (*model.User, error) {
	var user model.User

	err := ur.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	err := ur.db.WithContext(ctx).First(&user, "email = ?", email).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
	}

	return &user, nil
}

// UpdateUser implements [UserRepoItf].
func (ur *UserRepo) UpdateUser(ctx context.Context, user *model.User) error {
	err := ur.db.WithContext(ctx).Save(user).Error
	if err != nil {
		return err
	}

	return nil
}
