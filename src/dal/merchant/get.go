package merchant

import (
	io "api/src/models"
	"context"
	"errors"

	"github.com/jinzhu/gorm"
)

// GetMerchantByID ...
func (s *MerchantRepo) GetMerchantByID(ctx context.Context, db *gorm.DB, id int) (merchant io.AthMerchant, err error) {
	if db == nil {
		return merchant, errors.New("invalid db passed in GetMerchantByID")
	}
	var u io.AthMerchant
	err = db.Where(io.AthMerchant{MerchantID: int(id), IsActive: false}).Find(&u).Error
	if u.MerchantID != 0 {
		return u, err
	}
	return u, errors.New("Merchant does not exists")
}

// GetMerchantByUserID ...
func (s *MerchantRepo) GetMerchantByUserID(ctx context.Context, db *gorm.DB, id int) (merchant io.AthMerchant, err error) {
	if db == nil {
		return merchant, errors.New("invalid db passed in GetMerchantByUserID")
	}
	var u io.AthMerchant
	err = db.Where(io.AthMerchant{Models: io.Models{CreatedBy: int(id)}, IsActive: false}).Find(&u).Error
	if u.MerchantID != 0 {
		return u, err
	}
	return u, errors.New("Merchant does not exists")
}

// GetTeamData ...
func (s *MerchantRepo) GetTeamData(ctx context.Context, db *gorm.DB, id int, orderBy string) (teamMembers []io.AthVenueUser, err error) {
	if db == nil {
		return teamMembers, errors.New("invalid db passed in GetTeamData")
	}
	var u []io.AthVenueUser
	if orderBy == "asc" {
		err = db.Where(io.AthVenueUser{Models: io.Models{CreatedBy: int(id)}, IsActive: true}).Order("venueUserId asc").Find(&u).Error
	} else {
		err = db.Where(io.AthVenueUser{Models: io.Models{CreatedBy: int(id)}, IsActive: true}).Order("venueUserId desc").Find(&u).Error
	}
	if len(u) > 0 {
		return u, err
	}
	return u, errors.New("Team members does not present")
}

// GetRoleByRoleID ...
func (s *MerchantRepo) GetRoleByRoleID(ctx context.Context, db *gorm.DB, id int) (role io.AthRole, err error) {
	if db == nil {
		return role, errors.New("invalid db passed in GetRoleByRoleID")
	}
	var u io.AthRole
	err = db.Where(io.AthRole{RoleId: id}).Find(&u).Error
	if u.RoleId != 0 {
		return u, err
	}
	return u, errors.New("Role Does not exist")
}

// CheckMerchantDocExists ...
func (s *MerchantRepo) CheckMerchantDocExists(ctx context.Context, db *gorm.DB, docType string, MerchantID int) (exists bool, err error) {
	if db == nil {
		return false, errors.New("invalid db passed in CheckMerchantDocExists")
	}
	var u io.AthMerchant
	if docType == "gstNoFile" {
		err = db.Where(io.AthMerchant{MerchantID: MerchantID, GstFileVerified: true}).Find(&u).Error
	}
	if docType == "panNoFile" {
		err = db.Where(io.AthMerchant{MerchantID: MerchantID, PanFileVerified: true}).Find(&u).Error
	}
	if docType == "addressFile" {
		err = db.Where(io.AthMerchant{MerchantID: MerchantID, AddFileVerified: true}).Find(&u).Error
	}
	if docType == "bankAccFile" {
		err = db.Where(io.AthMerchant{MerchantID: MerchantID, BankFileVerified: true}).Find(&u).Error
	}
	if u.MerchantID != 0 {
		return true, err
	}
	return false, errors.New("document does not exist or not verified")
}
