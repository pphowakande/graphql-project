package user

import (
	io "api/src/models"
	"context"
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
)

// CreateUser ... create merchant user
func (s *UserRepo) CreateUser(ctx context.Context, db *gorm.DB, data io.AthUser) (newUser io.AthUser, err error) {
	if db == nil {
		return newUser, errors.New("invalid db passed in CreateUser")
	}
	var u io.AthUser

	// check if email already exists
	err = db.Where(io.AthUser{Email: data.Email}).Find(&u).Error
	if u.UserID != 0 {
		err = fmt.Errorf(`User email already exist`)
		return
	}
	var u1 io.AthUser
	// check if phone already exists
	err = db.Where(io.AthUser{Phone: data.Phone}).Find(&u1).Error
	if u1.UserID != 0 {
		err = fmt.Errorf(`User phone already exist`)
		return
	}
	d := db.Save(&data)
	return data, d.Error
}
