package venue

import (
	io "api/src/models"
	"context"
	"errors"

	"github.com/jinzhu/gorm"
	"stash.bms.bz/bms/gologger.git"
)

type Repository interface {
	GetDB(ctx context.Context, dbURL string) (*gorm.DB, error)
	BeginTransaction(ctx context.Context, dbURL string) (*gorm.DB, error)
	CreateVenue(ctx context.Context, db *gorm.DB, data io.AthVenues) (venue io.AthVenues, err error)
	EditVenue(ctx context.Context, db *gorm.DB, data io.AthVenues) (venue io.AthVenues, err error)
	CreateVenueHoliday(ctx context.Context, db *gorm.DB, data io.AthVenueHolidays) (err error)

	GetAmenitiesForVenueByID(ctx context.Context, db *gorm.DB, venueID int) (amenities []io.AthAmenities, err error)
	GetVenueByID(ctx context.Context, db *gorm.DB, id int) (venue io.AthVenues, err error)
	GetVenueHolidaysByID(ctx context.Context, db *gorm.DB, id int) (venueHoliday []io.AthVenueHolidays, err error)
	GetVenueHoursByID(ctx context.Context, db *gorm.DB, id int) (venueHour []io.AthVenueHours, err error)
	GetAmenitiesByID(ctx context.Context, db *gorm.DB, ids []int) (amenity []io.AthAmenities, err error)
	GetAllAmenities(ctx context.Context, db *gorm.DB) (amenity []io.AthAmenities, err error)
	GetListOfVenueByMerchantID(ctx context.Context, db *gorm.DB, merchantID int) (venues []io.AthVenues, err error)

	CreateVenueAmenity(ctx context.Context, db *gorm.DB, data io.AthVenueAmenitiesData) (err error)
	CreateVenueImage(ctx context.Context, db *gorm.DB, data io.VenueImagesReq) (err error)
	SaveHoursOfOperation(ctx context.Context, db *gorm.DB, data io.AthVenueHours) (err error)
	GetVenueImagesById(ctx context.Context, db *gorm.DB, venueId int, isactive bool) (venueImages []io.AthVenueImages, err error)
	CheckVenueImgExists(ctx context.Context, db *gorm.DB, imgUrL string) (exists bool, err error)
}

type VenueRepo struct {
	logger        *gologger.Logger
	DBConnections map[string]*gorm.DB
}

func NewVenueRepo(logger *gologger.Logger, dbConnections map[string]*gorm.DB) Repository {
	return &VenueRepo{logger: logger, DBConnections: dbConnections}
}

// BeginTransaction returns the object of transaction against the dbURL passed
func (s *VenueRepo) BeginTransaction(ctx context.Context, dbURL string) (*gorm.DB, error) {
	DB, ok := s.DBConnections[dbURL]
	if !ok {
		return nil, errors.New("invalid dbURL passed in VenueRepo")
	}
	return DB.Begin(), nil
}

// GetDB returns the DB object
func (s *VenueRepo) GetDB(ctx context.Context, dbURL string) (*gorm.DB, error) {
	DB, ok := s.DBConnections[dbURL]
	if !ok {
		return nil, errors.New("invalid dbURL passed in VenueRepo")
	}
	return DB, nil
}
