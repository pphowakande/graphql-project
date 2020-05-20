package otp

import (
	io "api/src/models"
	"context"
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
)

// CreateOTP ...
func (s *OTPRepo) CreateOTP(ctx context.Context, db *gorm.DB, data io.AthUserOTP) (err error) {
	fmt.Println("inside CreateOTP db function")
	if db == nil {
		return errors.New("invalid db passed in CreateOTP")
	}
	fmt.Println("data : ", data)
	d := db.Save(&data)
	fmt.Println("otp saved---------")
	return d.Error
}
