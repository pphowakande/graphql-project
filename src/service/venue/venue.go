package venue

import (
	"api/src/dal/venue"
	io "api/src/models"
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	venuepb "stash.bms.bz/turf/generic-proto-files.git/venue/v1"

	"stash.bms.bz/bms/gologger.git"

	"github.com/fatih/structs"
	"github.com/jinzhu/gorm"
)

type AthVenueService interface {
	CreateVenue(ctx context.Context, db *gorm.DB, u io.AthVenues) (res io.Response)
	EditVenue(ctx context.Context, db *gorm.DB, u io.AthVenues) (res io.Response)
	CreateVenueHoliday(ctx context.Context, db *gorm.DB, u io.AthVenueHolidays) (res io.Response)
	SaveHoursOfOperation(ctx context.Context, db *gorm.DB, u io.AthVenueHours) (res io.Response)
	GetVenueByID(ctx context.Context, db *gorm.DB, venueID int) (res io.Response)
	GetAmenitiesByID(ctx context.Context, db *gorm.DB, u venuepb.GetAmenitiesByIDReq) (res io.Response)
	GetAllAmenities(ctx context.Context, db *gorm.DB, u venuepb.GetAllAmenitiesReq) (res io.Response)
	GetListOfVenueByMerchantID(ctx context.Context, db *gorm.DB, u venuepb.GetListOfVenueByMerchantIDReq) (res io.Response)
	CreateVenueAmenity(ctx context.Context, db *gorm.DB, u io.AthVenueAmenitiesData) (res io.Response)
	GetAmenitiesForVenueByID(ctx context.Context, db *gorm.DB, venueID int) (res io.Response)
	CreateVenueImage(ctx context.Context, db *gorm.DB, u io.VenueImagesReq) (res io.Response)
}

type athVenueService struct {
	Logger *gologger.Logger
	DbRepo venue.Repository
}

func NewBasicAthVenueService(logger *gologger.Logger, DbRepo venue.Repository) AthVenueService {
	return &athVenueService{
		Logger: logger,
		DbRepo: DbRepo,
	}
}

func (b *athVenueService) CreateVenue(ctx context.Context, db *gorm.DB, u io.AthVenues) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"createVenue",
		"create venue request parameters",
		"createVenueRequest", "", structs.Map(u), true)
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if re.MatchString(u.Email) == false {
		res.Error = fmt.Errorf("invalid email")
		res.Message = "invalid email address"
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"createVenue",
			"create venue request failed",
			gologger.ValidationError, "", structs.Map(res), true)
		return
	}
	// save venue details
	newVenue, err := b.DbRepo.CreateVenue(ctx, db, u)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error creating new venue profile")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"createVenue",
			"create venue request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}
	data := make(map[string]interface{})
	data["venue_id"] = newVenue.VenueID

	res.Data = data
	res = io.SuccessMessage(data, "Venue Profile created")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"createVenue",
		"create venue response body",
		"createVenueResponse", "", structs.Map(res), true)
	return
}

func (b *athVenueService) CreateVenueHoliday(ctx context.Context, db *gorm.DB, u io.AthVenueHolidays) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"createVenueHoliday",
		"create venue holiday request parameters",
		"createVenueHolidayRequest", "", structs.Map(u), true)
	err := b.DbRepo.CreateVenueHoliday(ctx, db, u)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error creating/updating new venue holiday")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"createVenueHoliday",
			"create venue holiday request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}
	res = io.SuccessMessage(nil, "Venue Holiday created")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"createVenueHoliday",
		"create venue holiday response body",
		"createVenueHolidayResponse", "", structs.Map(res), true)
	return
}

func (b *athVenueService) SaveHoursOfOperation(ctx context.Context, db *gorm.DB, u io.AthVenueHours) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"saveHoursOfOperation",
		"save hours of operation request parameters",
		"saveHoursOfOperationRequest", "", structs.Map(u), true)
	err := b.DbRepo.SaveHoursOfOperation(ctx, db, u)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error creating new venue hour")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"saveHoursOfOperation",
			"save hours of operation request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}
	res = io.SuccessMessage(nil, "Venue Hour Added")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"saveHoursOfOperation",
		"save hours of operation response body",
		"saveHoursOfOperationResponse", "", structs.Map(res), true)
	return
}

func (b *athVenueService) EditVenue(ctx context.Context, db *gorm.DB, u io.AthVenues) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"editVenue",
		"edit venue request parameters",
		"editVenueRequest", "", structs.Map(u), true)
	_, err := b.DbRepo.EditVenue(ctx, db, u)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error updating venue")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"editVenue",
			"edit venue request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}
	res = io.SuccessMessage(nil, "Venue updated")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"editVenue",
		"edit venue response body",
		"editVenueResponse", "", structs.Map(res), true)
	return
}

func (b *athVenueService) GetVenueByID(ctx context.Context, db *gorm.DB, venueId int) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"getVenueByID",
		"get venue byID request parameters",
		"getVenueByIDRequest", "", structs.Map(struct{ venueId int }{venueId: int(venueId)}), true)
	// get basic details about venue
	Venue, err := b.DbRepo.GetVenueByID(ctx, db, venueId)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error getting venue details")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"getVenueByID",
			"get venue byID request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}

	// get venue amenity details using venue id
	amenityData, err := b.DbRepo.GetAmenitiesForVenueByID(ctx, db, int(venueId))
	if err != nil {
		res = io.FailureMessage(res.Error, "Error getting venue amenities")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"getVenueByID",
			"get venue byID request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}

	// get venue holidays
	holidaysData, err := b.DbRepo.GetVenueHolidaysByID(ctx, db, venueId)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error getting venue holidays")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"getVenueByID",
			"get venue byID request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}

	// get hours of operation for venue
	VenueHour, err := b.DbRepo.GetVenueHoursByID(ctx, db, venueId)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error getting venue hour details")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"getVenueByID",
			"get venue byID request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}

	// get venue images for venue
	VenueImg, err := b.DbRepo.GetVenueImagesById(ctx, db, venueId, true)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error getting venue images details")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"getVenueByID",
			"get venue byID request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}

	VenueData := make(map[string]interface{})

	VenueData["venueId"] = Venue.VenueID
	VenueData["name"] = Venue.VenueName
	VenueData["description"] = Venue.Description
	VenueData["address"] = Venue.Address
	VenueData["phone"] = Venue.Phone
	VenueData["email"] = Venue.Email
	VenueData["latitude"] = Venue.Latitude
	VenueData["longitude"] = Venue.Longitude

	// process amenities
	var amenityDataFinal []*venuepb.AmenityData
	for _, val := range amenityData {
		var eachAmenity venuepb.AmenityData
		eachAmenity.AmenityId = int32(val.AmenityID)
		eachAmenity.AmenityName = val.AmenityName
		amenityDataFinal = append(amenityDataFinal, &eachAmenity)
	}

	// arrange image data
	var finalImgData venuepb.CreateImageData
	var headerImgData []*venuepb.CreateImgData
	var thumbImgData []*venuepb.CreateImgData
	var galleryImgData []*venuepb.CreateImgData
	for _, val := range VenueImg {
		var eachImg venuepb.CreateImgData
		eachImg.ImgUrl = val.ImageUrl
		if val.ImageType == "header" {
			headerImgData = append(headerImgData, &eachImg)
		} else if val.ImageType == "thumbnail" {
			thumbImgData = append(thumbImgData, &eachImg)
		} else if val.ImageType == "gallery" {
			galleryImgData = append(galleryImgData, &eachImg)
		}
	}

	finalImgData.HeaderImg = headerImgData
	finalImgData.ThumbnailImg = thumbImgData
	finalImgData.GalleryImg = galleryImgData
	// process holidays
	var holidayDataFinal []*venuepb.HolidaysData
	for _, val := range holidaysData {
		var eachAHoliday venuepb.HolidaysData
		eachAHoliday.Title = val.Title
		eachAHoliday.Date = val.Date
		holidayDataFinal = append(holidayDataFinal, &eachAHoliday)
	}

	// process hours of operation
	var hourDataFinal []*venuepb.HoursOfOperationData
	for _, val := range VenueHour {
		keyExists := false
		for _, val1 := range hourDataFinal {
			if val.Day == val1.Day {
				// day already exists, get its value
				var eachTiming venuepb.Timing
				eachTiming.OpeningTime = val.OpeningTime
				eachTiming.ClosingTime = val.ClosingTime
				val1.Timing = append(val1.Timing, &eachTiming)
				keyExists = true
				break
			}
		}
		if keyExists == false {
			var eachHour venuepb.HoursOfOperationData
			var eachTiming venuepb.Timing
			eachTiming.OpeningTime = val.OpeningTime
			eachTiming.ClosingTime = val.ClosingTime
			eachHour.Day = val.Day
			eachHour.Timing = append(eachHour.Timing, &eachTiming)
			hourDataFinal = append(hourDataFinal, &eachHour)
		}
	}
	VenueData["amenities"] = amenityDataFinal
	VenueData["holidays"] = holidayDataFinal
	VenueData["hoursOfOperation"] = hourDataFinal
	VenueData["images"] = finalImgData

	res = io.SuccessMessage(VenueData, "Got venue profile details")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"getVenueByID",
		"get venue byID response body",
		"getVenueByIDResponse", "", structs.Map(res), true)
	return
}

func (b *athVenueService) GetAmenitiesByID(ctx context.Context, db *gorm.DB, u venuepb.GetAmenitiesByIDReq) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"getAmenitiesByID",
		"get amenities byID request parameters",
		"getAmenitiesByIDRequest", "", structs.Map(u), true)
	amenityList := make([]int, 0)
	if strings.Contains(u.AmenityIds, ",") {
		splittedIds := strings.Split(u.AmenityIds, ",")
		for _, eachAmenity := range splittedIds {
			eachAmenityId, _ := strconv.Atoi(eachAmenity)
			amenityList = append(amenityList, eachAmenityId)
		}
	} else {
		eachAmenityId, _ := strconv.Atoi(u.AmenityIds)
		amenityList = append(amenityList, eachAmenityId)
	}

	amenity, err := b.DbRepo.GetAmenitiesByID(ctx, db, amenityList)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error getting amenity details")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"getAmenitiesByID",
			"get amenities byID request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}
	res = io.SuccessMessage(amenity, "Got amenity details")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"getAmenitiesByID",
		"get amenities byID response body",
		"getAmenitiesByIDResponse", "", structs.Map(res), true)
	return
}

func (b *athVenueService) CreateVenueAmenity(ctx context.Context, db *gorm.DB, u io.AthVenueAmenitiesData) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"createVenueAmenity",
		"create venue amenity request parameters",
		"createVenueAmenityRequest", "", structs.Map(u), true)
	err := b.DbRepo.CreateVenueAmenity(ctx, db, u)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error creating/updating new venue amenity")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"createVenueAmenity",
			"create venue amenity request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}
	res = io.SuccessMessage(nil, "Venue amenity created/updated")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"createVenueAmenity",
		"create venue amenity response body",
		"createVenueAmenityResponse", "", structs.Map(res), true)
	return
}

func (b *athVenueService) CreateVenueImage(ctx context.Context, db *gorm.DB, u io.VenueImagesReq) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"CreateVenueImage",
		"create venue image request parameters",
		"createVenueimageRequest", "", structs.Map(u), true)
	err := b.DbRepo.CreateVenueImage(ctx, db, u)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error creating/updating new venue image")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"CreateVenueImage",
			"create venue image request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}
	res = io.SuccessMessage(nil, "Venue image created/updated")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"CreateVenueImage",
		"create venue image response body",
		"CreateVenueImageResponse", "", structs.Map(res), true)
	return
}

func (b *athVenueService) GetAmenitiesForVenueByID(ctx context.Context, db *gorm.DB, venueID int) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"getAmenitiesForVenueByID",
		"get amenities for venue byID request parameters",
		"getAmenitiesForVenueByIDRequest", "", map[string]interface{}{"venueID": venueID}, true)
	amenity_data, err := b.DbRepo.GetAmenitiesForVenueByID(ctx, db, venueID)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error getting venue amenities")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"getAmenitiesForVenueByID",
			"get amenities for venue byID request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}

	amenityList := make([]*venuepb.AmenityData, 0)

	if len(amenity_data) > 0 {
		for _, eachAmenity := range amenity_data {
			var amenityData venuepb.AmenityData

			amenityData.AmenityId = int32(eachAmenity.AmenityID)
			amenityData.AmenityName = eachAmenity.AmenityName
			amenityList = append(amenityList, &amenityData)
		}

		res = io.SuccessMessage(amenityList, "got Venue amenity data")
		b.Logger.Log(gologger.Info,
			gologger.InternalServices,
			"getAmenitiesForVenueByID",
			"get amenities for venue byID response body",
			"getAmenitiesForVenueByIDyResponse", "", structs.Map(res), true)
		return
	}
	res = io.SuccessMessage(nil, "got Venue amenity data")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"getAmenitiesForVenueByID",
		"get amenities for venue byID response body",
		"getAmenitiesForVenueByIDResponse", "", structs.Map(res), true)
	return
}

func (b *athVenueService) GetAllAmenities(ctx context.Context, db *gorm.DB, u venuepb.GetAllAmenitiesReq) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"GetAllAmenities",
		"get all amenities request parameters",
		"GetAllAmenitiesRequest", "", structs.Map(u), true)

	amenity, err := b.DbRepo.GetAllAmenities(ctx, db)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error getting amenity lists")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"GetAllAmenities",
			"get all amenities request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}
	amenityList := make([]*venuepb.AmenityData, 0)
	if len(amenity) > 0 {
		for _, eachAmenity := range amenity {
			var amenityData venuepb.AmenityData
			amenityData.AmenityId = int32(eachAmenity.AmenityID)
			amenityData.AmenityName = eachAmenity.AmenityName
			amenityList = append(amenityList, &amenityData)
		}
	}
	res = io.SuccessMessage(amenityList, "Got amenity details")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"GetAllAmenities",
		"get all amenities response body",
		"GetAllAmenitiesResponse", "", structs.Map(res), true)
	return
}

func (b *athVenueService) GetListOfVenueByMerchantID(ctx context.Context, db *gorm.DB, u venuepb.GetListOfVenueByMerchantIDReq) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"GetListOfVenueByMerchantID",
		"get list of venues for merchant request parameters",
		"GetListOfVenueByMerchantIDRequest", "", structs.Map(u), true)

	venues, err := b.DbRepo.GetListOfVenueByMerchantID(ctx, db, int(u.AccountId))
	if err != nil {
		res = io.FailureMessage(res.Error, "Error getting venue lists")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"GetListOfVenueByMerchantID",
			"get list of venues for merchant request parameters failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}
	res = io.SuccessMessage(venues, "Got venue list")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"GetListOfVenueByMerchantID",
		"get list of venues for merchant response body",
		"GetListOfVenueByMerchantIDResponse", "", structs.Map(res), true)
	return
}
