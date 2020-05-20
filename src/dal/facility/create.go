package facility

import (
	io "api/src/models"
	"context"
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

func (s *FacilityRepo) CreateFacility(ctx context.Context, db *gorm.DB, data io.AthFacilities) (newFacility io.AthFacilities, err error) {
	if db == nil {
		return newFacility, errors.New("invalid db passed in CreateFacility")
	}
	d := db.Save(&data)
	return data, d.Error
}

//CreateFacilitySlots ...
func (s *FacilityRepo) CreateFacilitySlots(ctx context.Context, db *gorm.DB, data io.AthFacilitySlots) (newFacilitySlots io.AthFacilitySlots, err error) {
	if db == nil {
		return newFacilitySlots, errors.New("invalid db passed in CreateFacilitySlots")
	}
	// check if slot already present
	var u io.AthFacilitySlots
	db.Where(io.AthFacilitySlots{FacilityID: data.FacilityID, SlotDays: data.SlotDays, SlotType: data.SlotType, SlotFromTime: data.SlotFromTime, SlotToTime: data.SlotToTime}).Find(&u)
	if u.FacilitySlotID != 0 {
		if data.IsActive == true {
			// facility slot already exists
			// update facility slot price
			u.SlotPrice = data.SlotPrice
			u.UpdatedAt = int(time.Now().Unix())
			u.UpdatedBy = data.UpdatedBy
			d := db.Save(&u)
			return data, d.Error
		}
		// delete facility slot
		// update deleted at and deleted by field too
		u.UpdatedAt = int(time.Now().Unix())
		u.UpdatedBy = data.UpdatedBy
		u.IsActive = false
		d := db.Save(&u)
		return u, d.Error
	}
	// slot does not exists
	if data.IsActive == true {
		// add new slot
		d := db.Save(&data)
		return data, d.Error
	}
	return data, nil
}

// BookFacility ...
// func (s *FacilityRepo) BookFacility(ctx context.Context, db *gorm.DB, data io.AthFacilityBookings) (newFacilityBookings io.AthFacilityBookings, err error) {
// 	DB, ok := s.DBConnections[db]
// 	if !ok {
// 		return newFacilityBookings, errors.New("invalid db passed in BookFacility")
// 	}
// 	d := db.Save(&data)
// 	if d.Error != nil {
// 		log.Printf("Failed to save error: %v", d.Error)
// 		return data, d.Error
// 	}
// 	return data, nil
// }

// BookFacilitySlots ...
// func (s *FacilityRepo) BookFacilitySlots(ctx context.Context, db *gorm.DB, data io.AthFacilityBookingSlots) (newFacilityBookingSlot io.AthFacilityBookingSlots, err error) {
// 	DB, ok := s.DBConnections[db]
// 	if !ok {
// 		return newFacilityBookingSlot, errors.New("invalid db passed in BookFacilitySlots")
// 	}
// 	d := db.Save(&data)
// 	if d.Error != nil {
// 		log.Printf("Failed to save error: %v", d.Error)
// 		return data, d.Error
// 	}
// 	return data, nil
// }

// AddFacilityCustomRates ...
func (s *FacilityRepo) AddFacilityCustomRates(ctx context.Context, db *gorm.DB, data io.AthFacilityCustomRates) (customRates io.AthFacilityCustomRates, err error) {
	if db == nil {
		return customRates, errors.New("invalid db passed in AddFacilityCustomRates")
	}
	// check if custom rate for same date slot already exists.
	var u io.AthFacilityCustomRates
	err = db.Where(io.AthFacilityCustomRates{SlotFromTime: data.SlotFromTime, SlotToTime: data.SlotToTime, FacilityID: data.FacilityID}).Find(&u).Error
	if u.RateID != 0 {
		if data.IsActive == true {
			// If yes, override rates
			u.SlotPrice = data.SlotPrice
			u.Available = data.Available
			u.UpdatedAt = int(time.Now().Unix())
			d := db.Save(&u)
			return u, d.Error
		}
		// delete custom rate slot
		u.IsActive = false
		u.Available = data.Available
		d := db.Save(&u)
		return u, d.Error
		// update deleted at and deleted by fields
	}
	//custom rate does not present , save the new record
	d := db.Save(&data)
	return data, d.Error
}

func (s *FacilityRepo) CreateFacilitySportCategory(ctx context.Context, db *gorm.DB, data io.AthFacilitySportsData) (err error) {
	if db == nil {
		return errors.New("invalid db passed in CreateFacilitySportCategory")
	}
	// if flag is true, save the data
	var categoryData io.AthFacilitySports
	// check if data already present
	err = db.Where(io.AthFacilitySports{SportCategoryID: data.SportCategoryID, FacilityID: data.FacilityID}).Find(&categoryData).Error
	if categoryData.FacilityID != 0 {
		if data.Status == false {
			err = db.Where(io.AthFacilitySports{SportCategoryID: data.SportCategoryID, FacilityID: data.FacilityID}).Delete(io.AthFacilitySports{}).Error
			return err
		}
	} else {
		if data.Status == true {
			saveReq := io.AthFacilitySports{
				SportCategoryID: data.SportCategoryID,
				FacilityID:      data.FacilityID,
			}
			d := db.Save(&saveReq)
			return d.Error
		}
	}
	return nil
}
