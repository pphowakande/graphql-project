package user

import (
	io "api/src/models"
	utile "api/src/utils"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dchest/uniuri"
	"github.com/jinzhu/gorm"
)

// EditUser ... Edit
func (s *UserRepo) EditUser(ctx context.Context, db *gorm.DB, data io.AthUser, userType string) (newUser io.AthUser, err error) {
	if db == nil {
		return newUser, errors.New("invalid db passed in EditUser")
	}
	var u io.AthUser
	err = db.Where(io.AthUser{UserID: data.UserID}).Find(&u).Error

	emailVerifyChange := false
	phoneVerifyChange := false

	if u.UserID != 0 {

		// check if passed email is same as in db
		if u.Email != data.Email {
			// check if same email address already exists for another user in db
			var uemail io.AthUser
			err = db.Where(io.AthUser{Email: data.Email}).Find(&uemail).Error
			if uemail.UserID == 0 {
				if userType == "teammember" {
					// check if email is already verified or not
					if u.EmailVerify == true {
						// updating verified email address is not allowed
						return data, errors.New("updating verified email address is not allowed")
					}
					emailVerifyChange = true
				}
			} else {
				return data, errors.New("Email already exists for another user")
			}
			emailVerifyChange = true
		}

		// check if passed phone number is same as in db
		if u.Phone != data.Phone {
			if emailVerifyChange == true {
				return data, errors.New("You can change either email or phone at a time")
			}
			// check if same phone number already exists for another user in db
			var uphone io.AthUser
			err = db.Where(io.AthUser{Phone: data.Phone}).Find(&uphone).Error
			if uphone.UserID == 0 {
				if userType == "teammember" {
					// check if phone is already verified or not
					if u.PhoneVerify == true {
						// updating verified email address is not allowed
						return data, errors.New("updating verified phone number is not allowed")
					}
					phoneVerifyChange = true
				}
			} else {
				return data, errors.New("Contact number already exists for another user")
			}
			phoneVerifyChange = true
		}

		if emailVerifyChange == true {
			u.EmailVerify = false
		}
		if phoneVerifyChange == true {
			u.PhoneVerify = false
		}

		u.Email = data.Email
		u.UpdatedBy = data.UpdatedBy
		u.FirstName = data.FirstName
		u.LastName = data.LastName
		u.Phone = data.Phone

		d := db.Save(&u)
		return u, d.Error
	}
	return data, errors.New("User does not exists")
}

// ResetPasswordUser ...
func (s *UserRepo) ResetPasswordUser(ctx context.Context, db *gorm.DB, data io.AthUser) (err error) {
	if db == nil {
		return errors.New("invalid db passed in ResetPasswordUser")
	}
	var user io.AthUser
	fmt.Println("Fecthing user details")
	err = db.Where(io.AthUser{Email: data.Email, IsActive: false}).Find(&user).Error
	fmt.Println("err : ", err)
	if err != nil {
		return err
	}
	if user.UserID != 0 {
		user.Password = data.Password
		user.UpdatedAt = int(time.Now().Unix())
		user.LastPasswordResetAt = int(time.Now().Unix())

		fmt.Println("Updated user is  --- ", user)
		db.Save(&user)
		fmt.Println("Updated user saved--------")
		return
	}
	return errors.New("Email does not exists")
}

// UpdateLoginData ...
func (s *UserRepo) UpdateLoginData(ctx context.Context, db *gorm.DB, UserID int) (user io.AthUser, err error) {
	if db == nil {
		return user, errors.New("invalid db passed in UpdateLoginData")
	}
	err = db.Where(io.AthUser{UserID: UserID}).Find(&user).Error
	if err == nil {
		user.LastLoginAt = int(time.Now().Unix())
		user.TokenValue = uniuri.NewLen(15)
		user.TokenHash, err = utile.Hash(user.TokenValue)
		if err != nil {
			return
		}
		user.TokenValue = fmt.Sprintf("%v:%v", user.TokenValue, user.UserID)
		user.TokenExpireAt = int(time.Now().AddDate(0, 0, 60).Unix())
		db.Save(&user)
	}
	return
}

// ValidateUser ...
func (s *UserRepo) ValidateUser(ctx context.Context, db *gorm.DB, ValidateUser io.ValidateUser) (user io.AthUser, err error) {
	if db == nil {
		return user, errors.New("invalid db passed in ValidateUser")
	}
	err = db.Where(io.AthUser{UserID: ValidateUser.UserID}).Find(&user).Error
	if err == nil {
		user.UpdatedAt = int(time.Now().Unix())
		user.UpdatedBy = ValidateUser.UserID
		if ValidateUser.EmailVerify == true {
			user.EmailVerify = ValidateUser.EmailVerify
		}
		if ValidateUser.PhoneVerify == true {
			user.PhoneVerify = ValidateUser.PhoneVerify
		}
		if ValidateUser.IsActive == true {
			user.IsActive = ValidateUser.IsActive
		}
		db.Save(&user)
	}
	return
}

// PhoneOTPVerified ...
func (s *UserRepo) PhoneOTPVerified(ctx context.Context, db *gorm.DB, UserID int) (err error) {
	if db == nil {
		return errors.New("invalid db passed in PhoneOTPVerified")
	}
	var user io.AthUser
	err = db.Where(io.AthUser{UserID: UserID}).Find(&user).Error
	if err == nil {
		user.UpdatedAt = int(time.Now().Unix())
		user.UpdatedBy = UserID
		user.PhoneVerify = true
		db.Save(&user)
	}
	return
}

// EmailOTPVerified ...
func (s *UserRepo) EmailOTPVerified(ctx context.Context, db *gorm.DB, UserID int) (err error) {
	if db == nil {
		return errors.New("invalid db passed in EmailOTPVerified")
	}
	var user io.AthUser
	err = db.Where(io.AthUser{UserID: UserID}).Find(&user).Error
	if err == nil {
		user.UpdatedAt = int(time.Now().Unix())
		user.UpdatedBy = UserID
		user.EmailVerify = true
		db.Save(&user)
	}
	return
}
