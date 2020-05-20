package venue

import (
	io "api/src/models"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

// CreateVenue ...
func (s *VenueRepo) CreateVenue(ctx context.Context, db *gorm.DB, data io.AthVenues) (newVenue io.AthVenues, err error) {
	if db == nil {
		return newVenue, errors.New("invalid db passed in CreateVenue")
	}
	var u io.AthVenues
	err = db.Where(io.AthVenues{Email: data.Email}).Find(&u).Error
	if u.VenueID != 0 {
		err = fmt.Errorf(`email address already exist`)
		return
	}
	var u1 io.AthVenues
	err = db.Where(io.AthVenues{Phone: data.Phone}).Find(&u1).Error
	if u1.VenueID != 0 {
		err = fmt.Errorf(`phone already exist`)
		return
	}
	d := db.Save(&data)
	return data, d.Error
}

// CreateVenueHoliday ...
func (s *VenueRepo) CreateVenueHoliday(ctx context.Context, db *gorm.DB, data io.AthVenueHolidays) (err error) {
	if db == nil {
		return errors.New("invalid db passed in CreateVenueHoliday")
	}
	// search the record in database
	var u io.AthVenueHolidays
	err = db.Where(io.AthVenueHolidays{Title: data.Title, Date: data.Date, VenueID: data.VenueID}).Find(&u).Error
	if u.HolidayID != 0 {
		// holiday alreasy exists
		if data.IsActive == false {
			// delete existing holiday
			u.IsActive = false
			err = db.Save(&u).Error
			return err
		}
		if data.IsActive == true {
			// if it already exists. but deleted. Lets reactive the same holiday
			if u.IsActive == false {
				u.IsActive = true
				err = db.Save(&u).Error
				return err
			}
		}
		return err
	}
	// holiday not present
	if data.IsActive == true {
		// add new holiday
		d := db.Save(&data)
		return d.Error
	}
	return nil
}

// SaveHoursOfOperation ...
func (s *VenueRepo) SaveHoursOfOperation(ctx context.Context, db *gorm.DB, data io.AthVenueHours) (err error) {
	if db == nil {
		return errors.New("invalid db passed in AthVenueHours")
	}
	// search the record in database
	var u io.AthVenueHours
	err = db.Where(io.AthVenueHours{VenueID: data.VenueID, Day: data.Day, OpeningTime: data.OpeningTime, ClosingTime: data.ClosingTime}).Find(&u).Error
	if u.HourID != 0 {
		// hour data alreasy exists
		if data.IsActive == false {
			// delete existing hour data
			u.IsActive = false
			err = db.Save(&u).Error
			return err
		}
		if data.IsActive == true {
			// if it already exists. but deleted. Lets reactive the same hour data
			if u.IsActive == false {
				u.IsActive = true
				err = db.Save(&u).Error
				return err
			}
		}
		return err
	}
	// hour data not present
	if data.IsActive == true {
		// add new hour data
		d := db.Save(&data)
		return d.Error
	}
	return nil
}

func (s *VenueRepo) CreateVenueAmenity(ctx context.Context, db *gorm.DB, data io.AthVenueAmenitiesData) (err error) {
	if db == nil {
		return errors.New("invalid db passed in CreateVenueAmenity")
	}
	// if flag is true, save the data
	var amenityData io.AthVenueAmenities
	err = db.Where(io.AthVenueAmenities{AmenityID: data.AmenityID, VenueID: data.VenueID}).Find(&amenityData).Error
	// check if data already present
	if amenityData.AmenityID != 0 {
		if data.Status == false {
			err = db.Where(io.AthVenueAmenities{AmenityID: data.AmenityID, VenueID: data.VenueID}).Delete(io.AthVenueAmenities{}).Error
			return err
		}
	} else {
		if data.Status == true {
			saveReq := io.AthVenueAmenities{
				AmenityID: data.AmenityID,
				VenueID:   data.VenueID,
			}
			d := db.Save(&saveReq)
			return d.Error
		}
	}
	return nil
}

func (s *VenueRepo) CreateVenueImage(ctx context.Context, db *gorm.DB, data io.VenueImagesReq) (err error) {
	if db == nil {
		return errors.New("invalid db passed in CreateVenueImage")
	}
	//save header images
	if len(data.HeaderImg) > 0 {
		for _, eachimg := range data.HeaderImg {
			if eachimg.Status == true {
				// save new image
				venueImg := io.AthVenueImages{
					VenueID:   data.VenueID,
					ImageType: "header",
					ImageUrl:  eachimg.Image,
					IsActive:  eachimg.Status,
					Models: io.Models{
						CreatedAt: int(time.Now().Unix()),
						CreatedBy: data.CreatedBy,
					},
				}
				d := db.Save(&venueImg)
				if d.Error != nil {
					return d.Error
				}
			} else {
				// delete image
				var venueImg io.AthVenueImages
				_ = db.Where(io.AthVenueImages{VenueID: data.VenueID, ImageUrl: eachimg.Image}).Find(&venueImg)
				if venueImg.VenueImageID != 0 {
					venueImg.IsActive = false
					// update deleted at and deleted by fields
					//venueImg.Models = io.Models{}
					d := db.Save(&venueImg)
					if d.Error != nil {
						return d.Error
					}
				}
			}
		}
	}
	// save thumbnail images
	if len(data.ThumbnailImg) > 0 {
		for _, eachimg := range data.ThumbnailImg {
			if eachimg.Status == true {
				// save new image
				venueImg := io.AthVenueImages{
					VenueID:   data.VenueID,
					ImageType: "thumbnail",
					ImageUrl:  eachimg.Image,
					IsActive:  eachimg.Status,
					Models: io.Models{
						CreatedAt: int(time.Now().Unix()),
						CreatedBy: data.CreatedBy,
					},
				}
				d := db.Save(&venueImg)
				if d.Error != nil {
					return d.Error
				}
			} else {
				// delete image
				var venueImg io.AthVenueImages
				_ = db.Where(io.AthVenueImages{VenueID: data.VenueID, ImageUrl: eachimg.Image}).Find(&venueImg)
				if venueImg.VenueImageID != 0 {
					venueImg.IsActive = false
					// update deleted at and deleted by fields
					//venueImg.Models = io.Models{}
					d := db.Save(&venueImg)
					if d.Error != nil {
						return d.Error
					}
				}
			}
		}
	}

	// save gallery images
	if len(data.GalleryImg) > 0 {
		for _, eachimg := range data.GalleryImg {
			if eachimg.Status == true {
				// save new image
				venueImg := io.AthVenueImages{
					VenueID:   data.VenueID,
					ImageType: "gallery",
					ImageUrl:  eachimg.Image,
					IsActive:  eachimg.Status,
					Models: io.Models{
						CreatedAt: int(time.Now().Unix()),
						CreatedBy: data.CreatedBy,
					},
				}
				d := db.Save(&venueImg)
				if d.Error != nil {
					return d.Error
				}
			} else {
				// delete image
				var venueImg io.AthVenueImages
				_ = db.Where(io.AthVenueImages{VenueID: data.VenueID, ImageUrl: eachimg.Image}).Find(&venueImg)
				if venueImg.VenueImageID != 0 {
					venueImg.IsActive = false
					// update deleted at and deleted by fields
					//venueImg.Models = io.Models{}
					d := db.Save(&venueImg)
					if d.Error != nil {
						return d.Error
					}
				}
			}
		}
	}
	return nil
}
