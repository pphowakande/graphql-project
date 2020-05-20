package facility

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
	CreateFacility(ctx context.Context, db *gorm.DB, data io.AthFacilities) (newFacility io.AthFacilities, err error)
	CreateFacilitySlots(ctx context.Context, db *gorm.DB, data io.AthFacilitySlots) (newFacilitySlots io.AthFacilitySlots, err error)
	// BookFacility(ctx context.Context, db *gorm.DB, data io.AthFacilityBookings) (newFacilityBookings io.AthFacilityBookings, err error)
	// BookFacilitySlots(ctx context.Context, db *gorm.DB, data io.AthFacilityBookingSlots) (newFacilityBookingSlot io.AthFacilityBookingSlots, err error)
	AddFacilityCustomRates(ctx context.Context, db *gorm.DB, data io.AthFacilityCustomRates) (customRates io.AthFacilityCustomRates, err error)

	EditFacility(ctx context.Context, db *gorm.DB, data io.AthFacilities) (newFacility io.AthFacilities, err error)

	GetSportCategories(ctx context.Context, db *gorm.DB) (sportCategories io.AthSportCategories, err error)
	GetAllSportCategories(ctx context.Context, db *gorm.DB) (venue []io.AthSportCategories, err error)
	GetSportCategoryByID(ctx context.Context, db *gorm.DB, id int32) (category io.AthSportCategories, err error)
	GetFacilityByID(ctx context.Context, db *gorm.DB, ids []int) (facilities []io.AthFacilities, err error)
	GetFacilitySlotsByID(ctx context.Context, db *gorm.DB, id int) (venue []io.AthFacilitySlots, err error)
	GetFacilityForVenueID(ctx context.Context, db *gorm.DB, id int32) (facilityIds []io.AthFacilities, err error)
	// GetFacilityLastBookingsData(ctx context.Context, db *gorm.DB, ids []string, startDate string, endDate string) (FacilityTotalBookingsData []io.FacilityLastBookingData, err error)
	// GetFacilityTotalEarnings(ctx context.Context, db *gorm.DB, ids []string, startDate string, endDate string) (totalEarn []io.FacilityTotalEarnings, err error)
	// GetFacilityTotalBookingsData(ctx context.Context, db *gorm.DB, ids []string, startDate string, endDate string) (FacilityTotalBookingsData []io.FacilityTotalBookingData, err error)
	// GetFacilityTodayBookingsData(ctx context.Context, db *gorm.DB, ids []string, startDate string, endDate string) (FacilityTodayBookingsData []io.FacilityTodayStats, err error)

	CreateFacilitySportCategory(ctx context.Context, db *gorm.DB, data io.AthFacilitySportsData) (err error)
	GetCategorisForFacilityByID(ctx context.Context, db *gorm.DB, id int) (categoryData []io.AthSportCategories, err error)
	GetCustomRatesForFacilityByID(ctx context.Context, db *gorm.DB, data io.CustomRatesForFacilityByID) (customRatesData []io.AthFacilityCustomRates, err error)
	DeleteFacilityByID(ctx context.Context, db *gorm.DB, db1 *gorm.DB, data io.DeleteFacilityByID) (facility io.AthFacilities, err error)
}

type FacilityRepo struct {
	logger        *gologger.Logger
	DBConnections map[string]*gorm.DB
}

func NewFacilityRepo(logger *gologger.Logger, dbConnections map[string]*gorm.DB) Repository {
	return &FacilityRepo{logger: logger, DBConnections: dbConnections}
}

// BeginTransaction returns the object of transaction against the dbURL passed
func (s *FacilityRepo) BeginTransaction(ctx context.Context, dbURL string) (*gorm.DB, error) {
	DB, ok := s.DBConnections[dbURL]
	if !ok {
		return nil, errors.New("invalid dbURL passed in FacilityRepo")
	}
	return DB.Begin(), nil
}

// GetDB returns the DB object
func (s *FacilityRepo) GetDB(ctx context.Context, dbURL string) (*gorm.DB, error) {
	DB, ok := s.DBConnections[dbURL]
	if !ok {
		return nil, errors.New("invalid dbURL passed in FacilityRepo")
	}
	return DB, nil
}
