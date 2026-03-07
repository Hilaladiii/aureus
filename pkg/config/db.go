package config

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(cfg Env) *gorm.DB {
	psqlInfo := fmt.Sprintf(
		"host=%s user=%s dbname=%s port=%s password=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBName, cfg.DBPort, cfg.DBPassword)

	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	return db
}

type (
	txKey        struct{}
	TxManagerItf interface {
		WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
	}
)

type txManager struct {
	db *gorm.DB
}

func NewTxManager(db *gorm.DB) TxManagerItf {
	return &txManager{db}
}

func (t *txManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, txKey{}, tx)
		return fn(txCtx)
	})
}

func ExtractTx(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return defaultDB
}
