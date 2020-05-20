package facility

import (
	io "api/src/models"
	"context"
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

// DeleteFacilityByID ...
func (s *FacilityRepo) DeleteFacilityByID(ctx context.Context, db *gorm.DB, db1 *gorm.DB, data io.DeleteFacilityByID) (facility io.AthFacilities, err error) {
	if db == nil {
		return facility, errors.New("invalid db passed in DeleteFacilityByID")
	}
	// check if facility present
	db.Where(io.AthFacilities{FacilityID: int(data.FacilityID), IsActive: true}).Find(&facility)
	if facility.FacilityID != 0 {
		// check if person deleting turf, is a owner
		loggedinuserId := int(data.UserID)
		if facility.CreatedBy == loggedinuserId {
			if data.Type == "delete" {
				var fb io.AthFacilityBookings
				// check if this facility have any active bookings
				todayDate := time.Now().Format("yyyy-mm-dd hh:mm:ss")
				db1.Where("facilityId IN (?) and bookingDate >= '"+todayDate+"'", data.FacilityID).Select("*").Find(&fb)
				if err != nil {
					return facility, err
				}
				if fb.FacilityID == 0 {
					// If not, then delete turf
					facility.IsActive = false
					err = db.Save(&facility).Error
					return facility, err
				}
				return facility, errors.New("Can not elete/unavailable a facility which has future bookings")
			}
			// If not, then delete turf
			facility.IsActive = false
			err = db.Save(&facility).Error
			return facility, err
		}
		return facility, errors.New("Only owner can delete/unavailable a facility")
	}
	return facility, errors.New("Facility does not exists")
}
