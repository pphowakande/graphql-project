package user

import (
	io "api/src/models"
	"context"
	"errors"

	"github.com/jinzhu/gorm"
	"stash.bms.bz/bms/gologger.git"
)

type Repository interface {
	GetDB(ctx context.Context, dbURL string) (*gorm.DB, error)
	BeginTransaction(ctx context.Context, dbURL string) (*gorm.DB, error)
	EditUser(ctx context.Context, db *gorm.DB, data io.AthUser, userType string) (user io.AthUser, err error)
	CreateUser(ctx context.Context, db *gorm.DB, data io.AthUser) (user io.AthUser, err error)
	UpdateLoginData(ctx context.Context, db *gorm.DB, UserID int) (user io.AthUser, err error)
	ValidateUser(ctx context.Context, db *gorm.DB, data io.ValidateUser) (user io.AthUser, err error)
	ResetPasswordUser(ctx context.Context, db *gorm.DB, data io.AthUser) (err error)
	GetUserByID(ctx context.Context, db *gorm.DB, id int, isActive bool) (user io.AthUser, err error)
	GetUserByEmailORPhone(ctx context.Context, db *gorm.DB, data io.LoginRequest) (user io.AthUser, err error)
	GetPrivilegesForUser(ctx context.Context, db *gorm.DB, UserID int) (privilegeData []io.AthVenueUser, err error)
	GetAccountTypeForUser(ctx context.Context, db *gorm.DB, UserID int) (accountType string, err error)
	PhoneOTPVerified(ctx context.Context, db *gorm.DB, UserID int) (err error)
	EmailOTPVerified(ctx context.Context, db *gorm.DB, UserID int) (err error)
}

type UserRepo struct {
	logger        *gologger.Logger
	DBConnections map[string]*gorm.DB
}

func NewUserRepo(logger *gologger.Logger, dbConnections map[string]*gorm.DB) Repository {
	return &UserRepo{logger: logger, DBConnections: dbConnections}
}

// BeginTransaction returns the object of transaction against the dbURL passed
func (s *UserRepo) BeginTransaction(ctx context.Context, dbURL string) (*gorm.DB, error) {
	DB, ok := s.DBConnections[dbURL]
	if !ok {
		return nil, errors.New("invalid dbURL passed in UserRepo")
	}
	return DB.Begin(), nil
}

// GetDB returns the DB object
func (s *UserRepo) GetDB(ctx context.Context, dbURL string) (*gorm.DB, error) {
	DB, ok := s.DBConnections[dbURL]
	if !ok {
		return nil, errors.New("invalid dbURL passed in UserRepo")
	}
	return DB, nil
}
