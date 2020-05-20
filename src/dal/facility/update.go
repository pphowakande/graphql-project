package facility

import (
	io "api/src/models"
	"context"
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

// EditFacility ...
func (s *FacilityRepo) EditFacility(ctx context.Context, db *gorm.DB, data io.AthFacilities) (newFacility io.AthFacilities, err error) {
	if db == nil {
		return newFacility, errors.New("invalid db passed in EditFacility")
	}
	var u io.AthFacilities
	err = db.Where(io.AthFacilities{FacilityID: data.FacilityID}).Find(&u).Error
	if err != nil {
		return u, err
	}

	if data.FacilityName != "" {
		u.FacilityName = data.FacilityName
	}

	if data.FacilityBasePrice != 0 {
		u.FacilityBasePrice = data.FacilityBasePrice
	}

	if data.UpdatedBy != 0 {
		u.UpdatedBy = data.UpdatedBy
	}

	u.UpdatedAt = int(time.Now().Unix())
	err = db.Save(&u).Error
	return u, err
}
