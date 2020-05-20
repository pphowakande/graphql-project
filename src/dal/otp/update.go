package otp

import (
	io "api/src/models"
	"context"
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

// VerifyOTP ...
func (s *OTPRepo) VerifyOTP(ctx context.Context, db *gorm.DB, data io.OTPVerify) (otp io.AthUserOTP, err error) {
	var u io.AthUserOTP
	if db == nil {
		return u, errors.New("invalid db passed in VerifyOTP")
	}
	// check if otp exists
	dbres := db.Where(io.AthUserOTP{OTPNO: data.VerificationCode, UserID: int(data.AccountId), OTPType: data.VerificationType, IsActive: false}).Find(&u).Error
	if dbres != nil {
		return u, errors.New("Invalid OTP number")
	}

	if u.OtpID != 0 {
		if u.OTPNO != data.VerificationCode {
			return u, errors.New("Invalid OTP number")
		}
		currenttime := int(time.Now().Unix())
		if u.ExpiredAt > currenttime {
			u.IsActive = true
			err = db.Save(&u).Error
			if err != nil {
				return u, err
			}
			return u, err
		}
		return u, errors.New("OTP is expired")
	}
	return u, errors.New("Invalid OTP number")
}
