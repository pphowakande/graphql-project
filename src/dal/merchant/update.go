package merchant

import (
	io "api/src/models"
	"context"
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

// EditMerchantByID ...
func (s *MerchantRepo) EditMerchantByID(ctx context.Context, db *gorm.DB, data io.AthMerchant) (merchant io.AthMerchant, err error) {
	if db == nil {
		return merchant, errors.New("invalid db passed in EditMerchantByID")
	}

	err = db.Where(io.AthMerchant{Models: io.Models{CreatedBy: data.CreatedBy}}).Find(&merchant).Error

	if merchant.MerchantID != 0 {

		if data.MerchantName != "" {
			merchant.MerchantName = data.MerchantName
		}

		if data.Phone != "" {
			merchant.Phone = data.Phone
		}

		if data.Email != "" {
			merchant.Email = data.Email
		}

		if data.GstNoFile != "" {
			merchant.GstNoFile = data.GstNoFile
		}

		if data.PanNoFile != "" {
			merchant.PanNoFile = data.PanNoFile
		}

		if data.AddressFile != "" {
			merchant.AddressFile = data.AddressFile
		}

		if data.BankAccFile != "" {
			merchant.BankAccFile = data.BankAccFile
		}

		if data.Models.UpdatedBy != 0 {
			merchant.UpdatedBy = data.UpdatedBy
		}

		merchant.DoLater = data.DoLater
		merchant.UpdatedAt = int(time.Now().Unix())

		d := db.Save(&merchant)
		return merchant, d.Error
	}
	return merchant, errors.New("Merchant does not exists")
}

// UpdateTeamMemberPrivileges ...
func (s *MerchantRepo) UpdateTeamMemberPrivileges(ctx context.Context, db *gorm.DB, data io.AthVenueUser) (err error) {
	if db == nil {
		return errors.New("invalid db passed in UpdateTeamMemberPrivileges")
	}
	var u io.AthVenueUser
	err = db.Where(io.AthVenueUser{VenueId: int(data.VenueId), UserId: data.UserId}).Find(&u).Error
	if err == nil {
		// Team member already added to venue
		// check status
		if data.IsActive == false {
			// update flag
			u.IsActive = data.IsActive
			d := db.Save(&u)
			return d.Error
		}
		if data.IsActive == true {
			// update flag
			u.IsActive = data.IsActive
			d := db.Save(&u)
			return d.Error
		}
	} else if data.IsActive == true {
		// add new venue user
		// get role id of the user
		var u1 io.AthVenueUser
		err = db.Where(io.AthVenueUser{UserId: data.UserId}).Find(&u1).Error

		var userSave io.AthVenueUser
		userSave.VenueId = int(data.VenueId)
		userSave.RoleId = u1.RoleId
		userSave.UserId = data.UserId
		userSave.IsActive = data.IsActive

		d := db.Save(&userSave)
		return d.Error
	}
	return nil
}

// DeleteTeamMember ...
func (s *MerchantRepo) DeleteTeamMember(ctx context.Context, db *gorm.DB, data io.AthVenueUser) (err error) {
	if db == nil {
		return errors.New("invalid db passed in DeleteTeamMember")
	}
	var u []io.AthVenueUser
	err = db.Where(io.AthVenueUser{UserId: int(data.UserId), IsActive: true, Models: io.Models{CreatedBy: data.CreatedBy}}).Find(&u).Error
	if err == nil {
		// Team member already added to venue
		if len(u) > 0 {
			for _, val := range u {
				val.IsActive = false
				val.UpdatedBy = data.UserId
				val.UpdatedAt = data.UpdatedAt
				d := db.Save(&val)
				if d.Error != nil {
					return d.Error
				}
			}
			return nil
		}
		return nil
	}
	return err
}
