package venue

import (
	io "api/src/models"
	"context"
	"errors"

	"github.com/jinzhu/gorm"
)

// GetVenueByID ...
func (s *VenueRepo) GetVenueByID(ctx context.Context, db *gorm.DB, id int) (venue io.AthVenues, err error) {
	if db == nil {
		return venue, errors.New("invalid db passed in GetVenueByID")
	}
	var u io.AthVenues
	err = db.Where(io.AthVenues{VenueID: int(id), IsActive: true}).Find(&u).Error
	return u, err
}

// GetVenueHolidaysByID ...
func (s *VenueRepo) GetVenueHolidaysByID(ctx context.Context, db *gorm.DB, id int) (venue []io.AthVenueHolidays, err error) {
	if db == nil {
		return venue, errors.New("invalid db passed in GetVenueHolidaysByID")
	}
	var u []io.AthVenueHolidays
	err = db.Where(io.AthVenueHolidays{VenueID: int(id), IsActive: true}).Find(&u).Error
	return u, err
}

// GetVenueHoursByID ...
func (s *VenueRepo) GetVenueHoursByID(ctx context.Context, db *gorm.DB, id int) (venue []io.AthVenueHours, err error) {
	if db == nil {
		return venue, errors.New("invalid db passed in GetVenueHoursByID")
	}
	var u []io.AthVenueHours
	err = db.Where(io.AthVenueHours{VenueID: int(id), IsActive: true}).Find(&u).Error
	return u, err
}

// GetAmenitiesByID ...
func (s *VenueRepo) GetAmenitiesByID(ctx context.Context, db *gorm.DB, ids []int) (venue []io.AthAmenities, err error) {
	if db == nil {
		return venue, errors.New("invalid db passed in GetAmenitiesByID")
	}
	var u []io.AthAmenities

	if len(ids) > 1 {
		err = db.Where(io.AthAmenities{IsActive: false}).Attrs(ids).Find(&u).Error
	} else {
		err = db.Where(io.AthAmenities{AmenityID: ids[0], IsActive: false}).Find(&u).Error
	}
	return u, err
}

// GetListOfVenueByMerchantID ...
func (s *VenueRepo) GetListOfVenueByMerchantID(ctx context.Context, db *gorm.DB, merchantID int) (venues []io.AthVenues, err error) {
	if db == nil {
		return venues, errors.New("invalid db passed in GetListOfVenueByMerchantID")
	}
	var u []io.AthVenues
	err = db.Where(io.AthVenues{MerchantID: merchantID}).Find(&u).Error
	return u, err
}

// GetAllAmenities ...
func (s *VenueRepo) GetAllAmenities(ctx context.Context, db *gorm.DB) (venue []io.AthAmenities, err error) {
	if db == nil {
		return venue, errors.New("invalid db passed in GetAllAmenities")
	}
	var u []io.AthAmenities
	err = db.Where(io.AthAmenities{IsActive: false}).Find(&u).Error
	return u, err
}

// GetAmenitiesForVenueByID ...
func (s *VenueRepo) GetAmenitiesForVenueByID(ctx context.Context, db *gorm.DB, VenueID int) (amenityData []io.AthAmenities, err error) {
	if db == nil {
		return amenityData, errors.New("invalid db passed in GetAmenitiesForVenueByID")
	}
	// get list of amenity ids for a venue
	var amentiyIds []io.AthVenueAmenities
	var amenityList []int
	var u []io.AthAmenities
	err = db.Where(io.AthVenueAmenities{VenueID: VenueID}).Find(&amentiyIds).Error
	if len(amentiyIds) > 0 {
		for _, value := range amentiyIds {
			amenityList = append(amenityList, value.AmenityID)
		}
		err = db.Where("amenityId IN (?)", amenityList).Find(&u).Error
		return u, err
	}
	return nil, err
}

// GetVenueImagesById ...
func (s *VenueRepo) GetVenueImagesById(ctx context.Context, db *gorm.DB, id int, IsActive bool) (venue []io.AthVenueImages, err error) {
	if db == nil {
		return venue, errors.New("invalid db passed in GetVenueImagesById")
	}
	var u []io.AthVenueImages
	err = db.Where(io.AthVenueImages{VenueID: int(id), IsActive: IsActive}).Find(&u).Error
	return u, err
}

// CheckVenueImgExists ...
func (s *VenueRepo) CheckVenueImgExists(ctx context.Context, db *gorm.DB, ImageUrl string) (exists bool, err error) {
	if db == nil {
		return false, errors.New("invalid db passed in CheckVenueImgExists")
	}
	var u io.AthVenueImages
	db.Where(io.AthVenueImages{ImageUrl: ImageUrl}).Find(&u)
	if u.VenueImageID != 0 {
		return true, nil
	}
	return false, nil
}
