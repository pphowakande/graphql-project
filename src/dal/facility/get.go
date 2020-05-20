package facility

import (
	io "api/src/models"
	"context"
	"errors"

	"github.com/jinzhu/gorm"
)

// GetSportCategories ...
func (s *FacilityRepo) GetSportCategories(ctx context.Context, db *gorm.DB) (sportCategories io.AthSportCategories, err error) {
	if db == nil {
		return sportCategories, errors.New("invalid db passed in GetSportCategories")
	}
	var u io.AthSportCategories
	err = db.Where(io.AthSportCategories{}).Find(&u).Error
	return u, err
}

// GetAllSportCategories ...
func (s *FacilityRepo) GetAllSportCategories(ctx context.Context, db *gorm.DB) (venue []io.AthSportCategories, err error) {
	if db == nil {
		return venue, errors.New("invalid db passed in GetAllSportCategories")
	}
	var u []io.AthSportCategories
	err = db.Where(io.AthSportCategories{IsActive: false}).Find(&u).Error
	return u, err
}

// GetSportCategoryByID ...
func (s *FacilityRepo) GetSportCategoryByID(ctx context.Context, db *gorm.DB, id int32) (category io.AthSportCategories, err error) {
	if db == nil {
		return category, errors.New("invalid db passed in GetSportCategoryByID")
	}
	var u io.AthSportCategories
	err = db.Where(io.AthSportCategories{SportCategoryID: int(id), IsActive: false}).Find(&u).Error
	return u, err
}

// GetFacilityByID ...
func (s *FacilityRepo) GetFacilityByID(ctx context.Context, db *gorm.DB, ids []int) (facilities []io.AthFacilities, err error) {
	if db == nil {
		return facilities, errors.New("invalid db passed in GetFacilityByID")
	}
	var u []io.AthFacilities
	err = db.Where("facilityId IN (?) and isActive = true", ids).Find(&u).Error
	return u, err
}

// GetFacilitySlotsByID ...
func (s *FacilityRepo) GetFacilitySlotsByID(ctx context.Context, db *gorm.DB, id int) (venue []io.AthFacilitySlots, err error) {
	if db == nil {
		return venue, errors.New("invalid db passed in GetFacilitySlotsByID")
	}
	var u []io.AthFacilitySlots
	err = db.Where(io.AthFacilitySlots{FacilityID: int(id), IsActive: true}).Find(&u).Error
	return u, err
}

// GetFacilityForVenueID ...
func (s *FacilityRepo) GetFacilityForVenueID(ctx context.Context, db *gorm.DB, id int32) (facilityIds []io.AthFacilities, err error) {
	if db == nil {
		return facilityIds, errors.New("invalid db passed in GetFacilityForVenueID")
	}
	var u []io.AthFacilities
	err = db.Where(io.AthFacilities{VenueID: int(id), IsActive: true}).Find(&u).Error
	return u, err
}

//GetFacilityLastBookingsData ...
// func (s *FacilityRepo) GetFacilityLastBookingsData(ctx context.Context, db *gorm.DB, ids []string, startDate string, endDate string) (FacilityTotalBookingsData []io.FacilityLastBookingData, err error) {
// 	DB, ok := s.DBConnections[db]
// 	if !ok {
// 		return FacilityTotalBookingsData, errors.New("invalid db passed in GetFacilityLastBookingsData")
// 	}
// 	var finalData []io.FacilityLastBookingData
// 	var data io.FacilityLastBookingData
// 	var u []io.AthFacilityBookings

// 	for _, eachID := range ids {
// 		err = db.Where("facility_id IN (?)", eachID).Select("booking_date, facility_id").Order("created_at desc").Limit(1).Find(&u).Scan(&data).Error
// 		// if err != nil {
// 		// 	log.Printf("Failed to get facility stats GetFacilityLastBookingsData: %v", err)
// 		// }

// 		finalData = append(finalData, data)
// 	}

// 	return finalData, nil
// }

//GetFacilityTotalEarnings ...
// func (s *FacilityRepo) GetFacilityTotalEarnings(ctx context.Context, db *gorm.DB, ids []string, startDate string, endDate string) (totalEarn []io.FacilityTotalEarnings, err error) {
// 	DB, ok := s.DBConnections[db]
// 	if !ok {
// 		return totalEarn, errors.New("invalid db passed in GetFacilityTotalEarnings")
// 	}
// 	var u []io.AthFacilityBookings
// 	var earnData []io.FacilityTotalEarnings

// 	if startDate != "" && endDate != "" {
// 		err = db.Where("facility_id IN (?) and booking_date between (?) and (?) ", ids, startDate, endDate).Select("sum(booking_amount) as totalEarnings, facility_id ").Group("facility_id").Find(&u).Scan(&earnData).Error
// 		if err != nil {
// 			log.Printf("Failed to get facility stats GetFacilityTotalEarnings: %v", err)
// 		}
// 	} else {
// 		err = db.Where("facility_id IN (?)", ids).Select("sum(booking_amount) as totalEarnings, facility_id ").Group("facility_id").Find(&u).Scan(&earnData).Error
// 		if err != nil {
// 			log.Printf("Failed to get facility stats GetFacilityTotalEarnings: %v", err)
// 		}
// 	}
// 	return earnData, nil
// }

// GetFacilityTotalBookingsData ...
// func (s *FacilityRepo) GetFacilityTotalBookingsData(ctx context.Context, db *gorm.DB, ids []string, startDate string, endDate string) (FacilityTotalBookingsData []io.FacilityTotalBookingData, err error) {
// 	DB, ok := s.DBConnections[db]
// 	if !ok {
// 		return FacilityTotalBookingsData, errors.New("invalid db passed in GetFacilityTotalBookingsData")
// 	}
// 	var data []io.FacilityTotalBookingData
// 	var u []io.AthFacilityBookings

// 	if startDate != "" && endDate != "" {
// 		err = db.Where("facility_id IN (?) and booking_date between (?) and (?) ", ids, startDate, endDate).Select("count(*) as no_of_bookings,facility_id").Group("facility_id").Find(&u).Scan(&data).Error
// 	} else {
// 		err = db.Where("facility_id IN (?) ", ids).Select("count(*) as no_of_bookings,facility_id").Group("facility_id").Find(&u).Scan(&data).Error
// 	}

// 	return data, nil
// }

// GetFacilityTodayBookingsData ...
// func (s *FacilityRepo) GetFacilityTodayBookingsData(ctx context.Context, db *gorm.DB, ids []string, startDate string, endDate string) (FacilityTodayBookingsData []io.FacilityTodayStats, err error) {
// 	DB, ok := s.DBConnections[db]
// 	if !ok {
// 		return FacilityTodayBookingsData, errors.New("invalid db passed in GetFacilityTodayBookingsData")
// 	}
// 	var data []io.FacilityTodayStats
// 	var u []io.AthFacilityBookings

// 	currentTime := time.Now()
// 	todayDate := currentTime.Format("2006-01-02")
// 	err = db.Where("facility_id IN (?) and booking_date like '%"+todayDate+"%'", ids).Select("count(*) as no_of_bookings_today,sum(booking_amount) as earned_today, facility_id").Group("facility_id").Find(&u).Scan(&data).Error
// 	if err != nil {
// 		log.Printf("Failed to get GetFacilityTotalBookingsData: %v", err)
// 	}

// 	return data, nil
// }

// GetCategorisForFacilityByID ...
func (s *FacilityRepo) GetCategorisForFacilityByID(ctx context.Context, db *gorm.DB, id int) (categoryData []io.AthSportCategories, err error) {
	if db == nil {
		return categoryData, errors.New("invalid db passed in GetCategorisForFacilityByID")
	}
	// get list of category ids for a facility
	var categoryIds []io.AthFacilitySports
	var categoryList []int
	var u []io.AthSportCategories
	err = db.Table("").Where("facilityId IN (?)", id).Find(&categoryIds).Error
	if len(categoryIds) > 0 {
		for _, value := range categoryIds {
			categoryList = append(categoryList, value.SportCategoryID)
		}
		err = db.Where("sportCategoryId IN (?)", categoryList).Find(&u).Error
		return u, err
	}
	return u, err
}

// GetCustomRatesForFacilityByID ...
func (s *FacilityRepo) GetCustomRatesForFacilityByID(ctx context.Context, db *gorm.DB, data io.CustomRatesForFacilityByID) (customRatesData []io.AthFacilityCustomRates, err error) {
	if db == nil {
		return customRatesData, errors.New("invalid db passed in GetCustomRatesForFacilityByID")
	}
	// get list of custom rates for a facility
	var customData []io.AthFacilityCustomRates
	if data.FromDate != "" && data.ToDate != "" {
		err = db.Table("").Where("facilityId IN (?) and date between (?) and (?) and isActive=true", data.FacilityID, data.FromDate, data.ToDate).Find(&customData).Error
		return customData, err
	}

	err = db.Table("").Where("facilityId IN (?) and isActive=true", data.FacilityID).Find(&customData).Error
	return customData, err
}
