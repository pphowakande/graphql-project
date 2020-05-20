package merchant

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
	CreateMerchant(ctx context.Context, db *gorm.DB, data io.AthMerchant) (merchant io.AthMerchant, err error)
	EditMerchantByID(ctx context.Context, db *gorm.DB, data io.AthMerchant) (merchant io.AthMerchant, err error)
	CreateMerchantUser(ctx context.Context, db *gorm.DB, data io.AthMerchantUser) error
	GetMerchantByID(ctx context.Context, db *gorm.DB, id int) (merchant io.AthMerchant, err error)
	GetMerchantByUserID(ctx context.Context, db *gorm.DB, id int) (merchant io.AthMerchant, err error)
	AddTeamMember(ctx context.Context, db *gorm.DB, data io.MemberData) (err error)
	UpdateTeamMemberPrivileges(ctx context.Context, db *gorm.DB, data io.AthVenueUser) (err error)
	DeleteTeamMember(ctx context.Context, db *gorm.DB, data io.AthVenueUser) (err error)
	GetTeamData(ctx context.Context, db *gorm.DB, userId int, orderBy string) (teamdata []io.AthVenueUser, err error)
	GetRoleByRoleID(ctx context.Context, db *gorm.DB, roleId int) (role io.AthRole, err error)
	CheckMerchantDocExists(ctx context.Context, db *gorm.DB, docType string, MerchantID int) (exists bool, err error)
}

type MerchantRepo struct {
	logger        *gologger.Logger
	DBConnections map[string]*gorm.DB
}

func NewMerchantRepo(logger *gologger.Logger, dbConnections map[string]*gorm.DB) Repository {
	return &MerchantRepo{logger: logger, DBConnections: dbConnections}
}

// BeginTransaction returns the object of transaction against the dbURL passed
func (s *MerchantRepo) BeginTransaction(ctx context.Context, dbURL string) (*gorm.DB, error) {
	DB, ok := s.DBConnections[dbURL]
	if !ok {
		return nil, errors.New("invalid dbURL passed in MerchantRepo")
	}
	return DB.Begin(), nil
}

// GetDB returns the DB object
func (s *MerchantRepo) GetDB(ctx context.Context, dbURL string) (*gorm.DB, error) {
	DB, ok := s.DBConnections[dbURL]
	if !ok {
		return nil, errors.New("invalid dbURL passed in MerchantRepo")
	}
	return DB, nil
}
