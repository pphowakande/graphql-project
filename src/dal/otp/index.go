package otp

import (
	io "api/src/models"
	"context"

	"github.com/jinzhu/gorm"
	"stash.bms.bz/bms/gologger.git"
)

type Repository interface {
	CreateOTP(ctx context.Context, db *gorm.DB, data io.AthUserOTP) (err error)
	VerifyOTP(ctx context.Context, db *gorm.DB, data io.OTPVerify) (otpData io.AthUserOTP, err error)
}

type OTPRepo struct {
	logger        *gologger.Logger
	DBConnections map[string]*gorm.DB
}

func NewOTPRepo(logger *gologger.Logger, dbConnections map[string]*gorm.DB) Repository {
	return &OTPRepo{logger: logger, DBConnections: dbConnections}
}
