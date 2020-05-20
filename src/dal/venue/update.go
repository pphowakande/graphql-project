package venue

import (
	io "api/src/models"
	"context"
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

// EditVenue ...
func (s *VenueRepo) EditVenue(ctx context.Context, db *gorm.DB, data io.AthVenues) (newVenue io.AthVenues, err error) {
	if db == nil {
		return newVenue, errors.New("invalid db passed in EditVenue")
	}
	var u io.AthVenues
	err = db.Where(io.AthVenues{VenueID: data.VenueID}).Find(&u).Error

	if data.VenueName != "" {
		u.VenueName = data.VenueName
	}

	if data.Description != "" {
		u.Description = data.Description
	}

	if data.Address != "" {
		u.Address = data.Address
	}

	if data.Phone != "" {
		u.Phone = data.Phone
	}

	if data.Email != "" {
		u.Email = data.Email
	}

	if data.Latitude != 0 {
		u.Latitude = data.Latitude
	}

	if data.Longitude != 0 {
		u.Longitude = data.Longitude
	}

	u.UpdatedBy = data.UpdatedBy
	u.UpdatedAt = int(time.Now().Unix())

	d := db.Save(&u)
	return data, d.Error
}
