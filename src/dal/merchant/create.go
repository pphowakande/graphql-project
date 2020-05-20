package merchant

import (
	io "api/src/models"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

// CreateMerchant ... Created merchant account for a merchant user
func (s *MerchantRepo) CreateMerchant(ctx context.Context, db *gorm.DB, data io.AthMerchant) (newMerchant io.AthMerchant, err error) {
	if db == nil {
		return newMerchant, errors.New("invalid db passed in CreateMerchant")
	}
	err = db.Save(&data).Error
	return data, err
}

// CreateMerchantUser ... Relate merchant with merchant user
func (s *MerchantRepo) CreateMerchantUser(ctx context.Context, db *gorm.DB, data io.AthMerchantUser) (err error) {
	if db == nil {
		return errors.New("invalid db passed in CreateMerchantUser")
	}
	err = db.Save(&data).Error
	return err
}

// AddTeamMember ... Add team member id against venue ID
func (s *MerchantRepo) AddTeamMember(ctx context.Context, db *gorm.DB, data io.MemberData) (err error) {
	if db == nil {
		return errors.New("invalid db passed in AddTeamMember")
	}

	// get account type id using account type
	var roleData io.AthRole
	err = db.Where(io.AthRole{RoleName: data.AccountType}).Find(&roleData).Error
	if err != nil {
		err = fmt.Errorf(`Role does not exists`)
		return
	}

	if roleData.RoleId != 0 {
		// role present.. go ahead
		var u io.AthVenueUser
		err = db.Where(io.AthVenueUser{VenueId: int(data.VenueID), UserId: data.UserID, RoleId: roleData.RoleId}).Find(&u).Error
		if err == nil {
			err = fmt.Errorf(`Team member already added to venue`)
			return
		}

		var userSave io.AthVenueUser
		userSave.VenueId = int(data.VenueID)
		userSave.RoleId = roleData.RoleId
		userSave.UserId = data.UserID
		userSave.IsActive = true
		userSave.CreatedBy = data.CreatedBy
		userSave.CreatedAt = int(time.Now().Unix())

		d := db.Save(&userSave)
		return d.Error
	}
	err = fmt.Errorf(`Role does not exists`)
	return err
}
