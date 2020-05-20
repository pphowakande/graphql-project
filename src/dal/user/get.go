package user

import (
	io "api/src/models"
	"context"
	"errors"

	"github.com/jinzhu/gorm"
)

// GetUserByID ...
func (s *UserRepo) GetUserByID(ctx context.Context, db *gorm.DB, id int, flag bool) (newUser io.AthUser, err error) {
	if db == nil {
		return newUser, errors.New("invalid db passed in GetUserByID")
	}
	var u io.AthUser
	err = db.Where(io.AthUser{UserID: id, IsActive: flag}).Find(&u).Error
	return u, err
}

// GetUserByEmailORPhone ...
func (s *UserRepo) GetUserByEmailORPhone(ctx context.Context, db *gorm.DB, data io.LoginRequest) (newUser io.AthUser, err error) {
	if db == nil {
		return newUser, errors.New("invalid db passed in GetUserByEmailORPhone")
	}
	var u io.AthUser

	if data.EmailParam == true {
		err = db.Where(io.AthUser{Email: data.Login}).Find(&u).Error
	} else {
		err = db.Where(io.AthUser{Phone: data.Login}).Find(&u).Error
	}
	if u.UserID == 0 {
		return u, errors.New("Email does not exist")
	}
	return u, err
}

// GetPrivilegesForUser ...
func (s *UserRepo) GetPrivilegesForUser(ctx context.Context, db *gorm.DB, UserID int) (privilegeData []io.AthVenueUser, err error) {
	if db == nil {
		return nil, errors.New("invalid db passed in GetPrivilegesForUser")
	}
	var data []io.AthVenueUser
	// check if user has been added to any team member of any merchant proile
	err = db.Where(io.AthVenueUser{UserId: UserID, IsActive: true}).Find(&data).Error
	if err != nil {
		return nil, errors.New("Error getting privilege data of a user")
	}
	return data, err
}

// GetAccountTypeForUser ...
func (s *UserRepo) GetAccountTypeForUser(ctx context.Context, db *gorm.DB, UserID int) (accountType string, err error) {
	if db == nil {
		return "", errors.New("invalid db passed in GetAccountTypeForUser")
	}
	var accType string
	var userData io.AthMerchant
	// check if logged in user is owner of any company or not
	err = db.Where(io.AthMerchant{Models: io.Models{CreatedBy: UserID}}).Find(&userData).Error
	if userData.MerchantID != 0 {
		// logged in user is owner
		accType = "owner"
	} else {
		var memberData []io.AthVenueUser
		// check if user has been added to any team member of any merchant proile
		err = db.Where(io.AthVenueUser{UserId: UserID, IsActive: false}).Find(&memberData).Error
		if err != nil {
			return "", errors.New("Error checking if user is team member or not")
		}
		if len(memberData) > 0 {
			roleID := memberData[0].RoleId
			// get role name using role id
			var roleData io.AthRole
			err = db.Where(io.AthRole{RoleId: roleID}).Find(&roleData).Error
			if err != nil {
				return "", errors.New("Error getting role name using role id")
			}
			accType = roleData.RoleName
		}
	}
	return accType, err
}
