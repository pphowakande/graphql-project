package facility

import (
	"api/src/dal/facility"
	io "api/src/models"
	"context"
	"fmt"
	"strconv"
	"strings"

	facilitypb "stash.bms.bz/turf/generic-proto-files.git/facility/v1"

	"stash.bms.bz/bms/gologger.git"

	"github.com/fatih/structs"
	"github.com/jinzhu/gorm"
)

type AthFacilityService interface {
	CreateFacility(ctx context.Context, db *gorm.DB, u io.AthFacilities) (res io.Response)
	EditFacility(ctx context.Context, db *gorm.DB, u io.AthFacilities) (res io.Response)
	CreateFacilitySlots(ctx context.Context, db *gorm.DB, u io.AthFacilitySlots) (res io.Response)
	// BookFacility(ctx context.Context,  db *gorm.DB,u io.AthFacilityBookings) (res io.Response)
	// BookFacilitySlots(ctx context.Context, db *gorm.DB, u io.AthFacilityBookingSlots) (res io.Response)
	AddFacilityCustomRates(ctx context.Context, db *gorm.DB, u io.AthFacilityCustomRates) (res io.Response)
	GetCustomRatesForFacilityByID(ctx context.Context, db *gorm.DB, u io.CustomRatesForFacilityByID) (res io.Response)

	GetAllSportCategories(ctx context.Context, db *gorm.DB, u facilitypb.GetAllSportCategoriesReq) (res io.Response)
	GetSportCategoryByID(ctx context.Context, db *gorm.DB, u facilitypb.GetSportCategoryByIDReq) (res io.Response)
	GetFacilityByID(ctx context.Context, db *gorm.DB, u facilitypb.GetFacilityByIDReq) (res io.Response)
	GetFacilitySlotsByID(ctx context.Context, db *gorm.DB, u facilitypb.GetFacilitySlotsByIDReq) (res io.Response)
	GetFacilityForVenueID(ctx context.Context, db *gorm.DB, u facilitypb.GetFacilityForVenueIDReq) (res io.Response)
	// GetFacilityStats(ctx context.Context,  db *gorm.DB,u facilitypb.GetFacilityStatsReq) (res io.Response)
	DeleteFacilityByID(ctx context.Context, db *gorm.DB, db1 *gorm.DB, u io.DeleteFacilityByID) (res io.Response)
	CreateFacilitySportCategory(ctx context.Context, db *gorm.DB, u io.AthFacilitySportsData) (res io.Response)
}

type athFacilityService struct {
	Logger *gologger.Logger
	DbRepo facility.Repository
}

func NewBasicAthFacilityService(logger *gologger.Logger, DbRepo facility.Repository) AthFacilityService {
	return &athFacilityService{
		Logger: logger,
		DbRepo: DbRepo,
	}
}

func (b *athFacilityService) CreateFacility(ctx context.Context, db *gorm.DB, u io.AthFacilities) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"createFacility",
		"create facility request parameters",
		"createFacilityRequest", "", structs.Map(u), true)
	newFacility, err := b.DbRepo.CreateFacility(ctx, db, u)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error creating new facility")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"createFacility",
			"create facility request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}
	res = io.SuccessMessage(newFacility, "Facility created")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"createFacility",
		"create facility response body",
		"createFacilityResponse", "", structs.Map(res), true)
	return
}

func (b *athFacilityService) EditFacility(ctx context.Context, db *gorm.DB, u io.AthFacilities) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"editFacility",
		"edit facility request parameters",
		"editFacilityRequest", "", structs.Map(u), true)
	newFacility, err := b.DbRepo.EditFacility(ctx, db, u)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error editing new facility")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"editFacility",
			"edit facility request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}
	data := make(map[string]interface{})
	data["facility_id"] = newFacility.FacilityID
	res.Data = data
	res = io.SuccessMessage(data, "Facility edited")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"editFacility",
		"edit facility response body",
		"editFacilityResponse", "", structs.Map(res), true)
	return
}

func (b *athFacilityService) CreateFacilitySlots(ctx context.Context, db *gorm.DB, u io.AthFacilitySlots) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"createFacilitySlots",
		"create facility slots request parameters",
		"createFacilitySlotsRequest", "", structs.Map(u), true)
	newFacility, err := b.DbRepo.CreateFacilitySlots(ctx, db, u)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error creating new facility slot")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"createFacilitySlots",
			"create facility slots request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}
	data := make(map[string]interface{})
	data["facility_slot_id"] = newFacility.FacilityID
	res.Data = data
	res = io.SuccessMessage(data, "Facility slot created")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"createFacilitySlots",
		"create facility slots response body",
		"createFacilitySlotsResponse", "", structs.Map(res), true)
	return
}

// func (b *athFacilityService) BookFacility(ctx context.Context,db *gorm.DB, u io.AthFacilityBookings) (res io.Response) {
// 	b.Logger.Log(gologger.Info,
// 		gologger.InternalServices,
// 		"bookFacility",
// 		"book facility request parameters",
// 		"bookFacilityRequest", "", structs.Map(u), true)
// 	newFacilityBooking, err := b.DbRepo.BookFacility(ctx, db, u)
// 	if err != nil {
// 		res = io.FailureMessage(res.Error, "Error creating new facility booking")
// 		res.Error = err
// 		b.Logger.Log(gologger.Errsev3,
// 			gologger.InternalServices,
// 			"bookFacility",
// 			"book facility request failed",
// 			gologger.ParseError, "", structs.Map(res), true)
// 		return
// 	}
// 	data := make(map[string]interface{})
// 	data["facility_booking_id"] = newFacilityBooking.BookingID
// 	data["facility_booking_no"] = newFacilityBooking.BookingNo
// 	res.Data = data
// 	res = io.SuccessMessage(data, "Facility booking created")
// 	b.Logger.Log(gologger.Info,
// 		gologger.InternalServices,
// 		"bookFacility",
// 		"book facility response body",
// 		"bookFacilityResponse", "", structs.Map(res), true)
// 	return
// }

// func (b *athFacilityService) BookFacilitySlots(ctx context.Context, db *gorm.DB,u io.AthFacilityBookingSlots) (res io.Response) {
// 	b.Logger.Log(gologger.Info,
// 		gologger.InternalServices,
// 		"bookFacilitySlots",
// 		"book facility slots request parameters",
// 		"bookFacilitySlotsRequest", "", structs.Map(u), true)
// 	newFacilityBookingSlot, err := b.DbRepo.BookFacilitySlots(ctx, db, u)
// 	if err != nil {
// 		res = io.FailureMessage(res.Error, "Error creating new facility booking slot")
// 		res.Error = err
// 		b.Logger.Log(gologger.Errsev3,
// 			gologger.InternalServices,
// 			"bookFacilitySlots",
// 			"book facility slots request failed",
// 			gologger.ParseError, "", structs.Map(res), true)
// 		return
// 	}
// 	data := make(map[string]interface{})
// 	data["facility_booking_slot_id"] = newFacilityBookingSlot.BookingSlotID
// 	res.Data = data
// 	res = io.SuccessMessage(data, "Facility booking slot created")
// 	b.Logger.Log(gologger.Info,
// 		gologger.InternalServices,
// 		"bookFacilitySlots",
// 		"book facility slots response body",
// 		"bookFacilitySlotsResponse", "", structs.Map(res), true)
// 	return
// }

func (b *athFacilityService) AddFacilityCustomRates(ctx context.Context, db *gorm.DB, u io.AthFacilityCustomRates) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"addFacilityCustomRates",
		"add facility custom rates request parameters",
		"addFacilityCustomRatesRequest", "", structs.Map(u), true)
	_, err := b.DbRepo.AddFacilityCustomRates(ctx, db, u)
	if err != nil {
		res = io.FailureMessage(res.Error, "err adding custom rates for facility")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"addFacilityCustomRates",
			"add facility custom rates request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}
	res = io.SuccessMessage(nil, "Facility custom rate added")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"addFacilityCustomRates",
		"add facility custom rates response body",
		"addFacilityCustomRatesResponse", "", structs.Map(res), true)
	return
}

func (b *athFacilityService) GetAllSportCategories(ctx context.Context, db *gorm.DB, u facilitypb.GetAllSportCategoriesReq) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"getAllSportCategories",
		"getAll sport categories request parameters",
		"getAllSportCategoriesRequest", "", structs.Map(u), true)
	catgories, err := b.DbRepo.GetAllSportCategories(ctx, db)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error getting sport categories")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"getAllSportCategories",
			"getAll sport categories request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}
	categoryList := make([]*facilitypb.SportCategoryData, 0)
	for _, eachCat := range catgories {
		var catData facilitypb.SportCategoryData

		catData.CategoryId = int32(eachCat.SportCategoryID)
		catData.CategoryName = eachCat.CategoryName

		categoryList = append(categoryList, &catData)
	}
	res = io.SuccessMessage(categoryList, "Got all sport categories")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"getAllSportCategories",
		"getAll sport categories response body",
		"getAllSportCategoriesResponse", "", structs.Map(res), true)
	return
}

func (b *athFacilityService) GetSportCategoryByID(ctx context.Context, db *gorm.DB, u facilitypb.GetSportCategoryByIDReq) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"getSportCategoryByID",
		"get sport category byID request parameters",
		"getSportCategoryByIDRequest", "", structs.Map(u), true)
	Category, err := b.DbRepo.GetSportCategoryByID(ctx, db, u.CategoryId)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error getting sport category data")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"getSportCategoryByID",
			"get sport category byID  request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}
	categoryList := make([]*facilitypb.SportCategoryData, 0)
	var catData facilitypb.SportCategoryData
	catData.CategoryId = int32(Category.SportCategoryID)
	catData.CategoryName = Category.CategoryName
	categoryList = append(categoryList, &catData)

	res = io.SuccessMessage(categoryList, "Got sport categorydetails")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"getSportCategoryByID",
		"get sport category byID response body",
		"getSportCategoryByIDResponse", "", structs.Map(res), true)
	return
}

func (b *athFacilityService) GetFacilityByID(ctx context.Context, db *gorm.DB, u facilitypb.GetFacilityByIDReq) (res io.Response) {
	facilityList := make([]int, 0)
	if strings.Contains(u.FacilityIds, ",") {
		splittedIds := strings.Split(u.FacilityIds, ",")
		for _, eachFacility := range splittedIds {
			eachFacilityID, _ := strconv.Atoi(eachFacility)
			facilityList = append(facilityList, eachFacilityID)
		}
	} else {
		eachFacilityID, _ := strconv.Atoi(u.FacilityIds)
		facilityList = append(facilityList, eachFacilityID)
	}

	// get facility details
	Facility, err := b.DbRepo.GetFacilityByID(ctx, db, facilityList)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error getting facility data")
		res.Error = err
		return
	}
	facilityFinalList := make([]*facilitypb.FacilityData, 0)
	if len(Facility) > 0 {
		for _, eachFacility := range Facility {
			categoryData, err := b.DbRepo.GetCategorisForFacilityByID(ctx, db, eachFacility.FacilityID)
			if err != nil {
				res = io.FailureMessage(res.Error, "Error getting facility sport categories")
				res.Error = err
				b.Logger.Log(gologger.Errsev3,
					gologger.InternalServices,
					"GetFacilityByID",
					"get facility byID request failed",
					gologger.ParseError, "", structs.Map(res), true)
				return
			}

			var catDict []*facilitypb.SportCategoryData
			for _, eachCat := range categoryData {
				var eachCatDict facilitypb.SportCategoryData
				eachCatDict.CategoryId = int32(eachCat.SportCategoryID)
				eachCatDict.CategoryName = eachCat.CategoryName
				catDict = append(catDict, &eachCatDict)
			}

			fmt.Println("catDict : ", catDict)

			facilitycustomReq := io.CustomRatesForFacilityByID{
				FacilityID: int32(eachFacility.FacilityID),
			}
			customRatesData, err := b.DbRepo.GetCustomRatesForFacilityByID(ctx, db, facilitycustomReq)
			if err != nil {
				res = io.FailureMessage(res.Error, "Error getting facility custom rates")
				res.Error = err
				b.Logger.Log(gologger.Errsev3,
					gologger.InternalServices,
					"GetFacilityByID",
					"get facility byID request failed",
					gologger.ParseError, "", structs.Map(res), true)
				return
			}

			var finalCustomData []*facilitypb.CustomRates

			if len(customRatesData) > 0 {
				for _, eachData := range customRatesData {
					var eachCustomData facilitypb.CustomRates
					eachCustomData.StartDate = eachData.SlotFromTime
					eachCustomData.EndDate = eachData.SlotToTime
					eachCustomData.Price = eachData.SlotPrice
					eachCustomData.Available = eachData.IsActive
					finalCustomData = append(finalCustomData, &eachCustomData)
				}
			}

			fmt.Println("finalCustomData : ", finalCustomData)

			FacilitySlots, err := b.DbRepo.GetFacilitySlotsByID(ctx, db, eachFacility.FacilityID)
			if err != nil {
				res = io.FailureMessage(res.Error, "Error getting facility slot data")
				res.Error = err
				b.Logger.Log(gologger.Errsev3,
					gologger.InternalServices,
					"GetFacilityByID",
					"get facility  byID request failed",
					gologger.ParseError, "", structs.Map(res), true)
				return
			}
			timeSlotList := make([]*facilitypb.FacilitySlotData, 0)
			if len(FacilitySlots) > 0 {
				for _, eachslot := range FacilitySlots {
					fmt.Println("--------------------------")
					fmt.Println("eachslot : ", eachslot)
					fmt.Println("eachslot.SlotFromTime : ", eachslot.SlotFromTime)
					fmt.Println("eachslot.SlotToTime : ", eachslot.SlotToTime)
					var slotData facilitypb.FacilitySlotData
					slotData.SlotDays = eachslot.SlotDays
					slotData.SlotStartTime = eachslot.SlotFromTime
					slotData.SlotEndTime = eachslot.SlotToTime
					slotData.SlotType = eachslot.SlotType
					slotData.SlotPrice = eachslot.SlotPrice
					fmt.Println("slotData : ", slotData)
					timeSlotList = append(timeSlotList, &slotData)
				}
			}

			fmt.Println("timeSlotList : ", timeSlotList)

			var facData facilitypb.FacilityData

			facData.DefaultRate = eachFacility.FacilityBasePrice
			facData.VenueId = int32(eachFacility.VenueID)
			facData.FacilityName = eachFacility.FacilityName
			facData.TimeSlot = int32(eachFacility.TimeSlot)
			facData.FacilityId = int32(eachFacility.FacilityID)
			facData.CustomRates = finalCustomData
			facData.Available = eachFacility.IsActive
			facData.WeekSlots = timeSlotList
			facData.CategoryData = catDict

			facilityFinalList = append(facilityFinalList, &facData)

			fmt.Println("final facData : ", facData)
		}
	}

	fmt.Println("retunrning facilityFinalList : ", facilityFinalList)

	res = io.SuccessMessage(facilityFinalList, "Got facility details")
	return
}

func (b *athFacilityService) GetFacilitySlotsByID(ctx context.Context, db *gorm.DB, u facilitypb.GetFacilitySlotsByIDReq) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"getFacilitySlotsByID",
		"get facility slots byID request parameters",
		"getFacilitySlotsByIDRequest", "", structs.Map(u), true)
	FacilitySlots, err := b.DbRepo.GetFacilitySlotsByID(ctx, db, int(u.FacilityId))
	if err != nil {
		res = io.FailureMessage(res.Error, "Error getting facility slot data")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"getFacilitySlotsByID",
			"get facility slots byID request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}
	res = io.SuccessMessage(FacilitySlots, "Got facility slot details")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"getFacilitySlotsByID",
		"get facility slots byID response body",
		"getFacilitySlotsByIDResponse", "", structs.Map(res), true)
	return
}

func (b *athFacilityService) GetFacilityForVenueID(ctx context.Context, db *gorm.DB, u facilitypb.GetFacilityForVenueIDReq) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"getFacilityForVenueID",
		"get facility for venueID request parameters",
		"getFacilityForVenueIDRequest", "", structs.Map(u), true)
	FaciliyIds, err := b.DbRepo.GetFacilityForVenueID(ctx, db, u.VenueId)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error getting facility ids for a venue")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"getFacilityForVenueID",
			"get facility for venueID request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}

	var facilityIdsList []*facilitypb.FacilityIDData
	if len(FaciliyIds) > 0 {
		for _, eachFacIDData := range FaciliyIds {
			var eachData facilitypb.FacilityIDData
			eachData.FacilityId = int32(eachFacIDData.FacilityID)
			eachData.FacilityName = eachFacIDData.FacilityName
			facilityIdsList = append(facilityIdsList, &eachData)
		}
	}
	res = io.SuccessMessage(facilityIdsList, "Got facility ids for a venue")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"getFacilityForVenueID",
		"get facility for venueID response body",
		"getFacilityForVenueIDResponse", "", structs.Map(res), true)
	return
}

// func (b *athFacilityService) GetFacilityStats(ctx context.Context,db *gorm.DB, u facilitypb.GetFacilityStatsReq) (res io.Response) {
// 	b.Logger.Log(gologger.Info,
// 		gologger.InternalServices,
// 		"getFacilityStats",
// 		"get facility stats request parameters",
// 		"getFacilityStatsRequest", "", structs.Map(u), true)
// 	facilityList := make([]string, 0)
// 	if strings.Contains(u.FacilityIds, ",") {
// 		splittedIds := strings.Split(u.FacilityIds, ",")
// 		for _, eachFacility := range splittedIds {
// 			facilityList = append(facilityList, eachFacility)
// 		}
// 	} else {
// 		facilityList = append(facilityList, u.FacilityIds)
// 	}

// 	FacilityTotalBookingsData, err := b.DbRepo.GetFacilityTotalBookingsData(ctx, db, facilityList, u.StartDate, u.EndDate)
// 	// fmt.Println("FacilityTotalBookingsData : ", FacilityTotalBookingsData)
// 	// fmt.Println("type : ", reflect.TypeOf(FacilityTotalBookingsData))
// 	if err != nil {
// 		res = io.FailureMessage(res.Error, "Error getting FacilityTotalBookingsData")
// 		res.Error = err
// 		return
// 	}

// 	FacilityLastBookingsData, err := b.DbRepo.GetFacilityLastBookingsData(ctx, db, facilityList, u.StartDate, u.EndDate)
// 	// fmt.Println("FacilityLastBookingsData : ", FacilityLastBookingsData)
// 	// fmt.Println("type : ", reflect.TypeOf(FacilityLastBookingsData))
// 	if err != nil {
// 		res = io.FailureMessage(res.Error, "Error getting FacilityLastBookingsData")
// 		res.Error = err
// 		b.Logger.Log(gologger.Errsev3,
// 			gologger.InternalServices,
// 			"getFacilityStats",
// 			"get facility stats request failed",
// 			gologger.ParseError, "", structs.Map(res), true)
// 		return
// 	}

// 	FacilityTotalEarningsData, err := b.DbRepo.GetFacilityTotalEarnings(ctx, db, facilityList, u.StartDate, u.EndDate)
// 	// fmt.Println("FacilityTotalEarningsData : ", FacilityTotalEarningsData)
// 	// fmt.Println("type : ", reflect.TypeOf(FacilityTotalEarningsData))
// 	if err != nil {
// 		res = io.FailureMessage(res.Error, "Error getting FacilityTotalEarningsData")
// 		res.Error = err
// 		b.Logger.Log(gologger.Errsev3,
// 			gologger.InternalServices,
// 			"getFacilityStats",
// 			"get facility stats request failed",
// 			gologger.ParseError, "", structs.Map(res), true)
// 		return
// 	}

// 	FacilityTodayBookingsData, err := b.DbRepo.GetFacilityTodayBookingsData(ctx, db, facilityList, u.StartDate, u.EndDate)
// 	fmt.Println("FacilityTodayBookingsData : ", FacilityTodayBookingsData)
// 	fmt.Println("len(FacilityTodayBookingsData) : ", len(FacilityTodayBookingsData))
// 	if err != nil {
// 		res = io.FailureMessage(res.Error, "Error getting FacilityTodayBookingsData")
// 		res.Error = err
// 		b.Logger.Log(gologger.Errsev3,
// 			gologger.InternalServices,
// 			"getFacilityStats",
// 			"get facility stats request failed",
// 			gologger.ParseError, "", structs.Map(res), true)
// 		return
// 	}

// 	// rearranging data
// 	finalData := make(map[string]*facilitypb.FacilityStats)

// 	for _, facilityID := range facilityList {

// 		/** converting the str1 variable into an int using Atoi method */
// 		facilityIDInt, _ := strconv.Atoi(facilityID)
// 		var stats facilitypb.FacilityStats

// 		fmt.Println("stats before  : ", stats)
// 		// stats.LastBookingDate = ""
// 		// stats.TotalEarnings = 0
// 		// stats.TotalEarningsToday = 0
// 		// stats.NoOfBookings = 0
// 		// stats.NoOfBookingsToday = 0

// 		if len(FacilityTodayBookingsData) > 0 {
// 			for _, each_today_booking_data := range FacilityTodayBookingsData {
// 				fmt.Println("each_today_booking_data : ", each_today_booking_data)
// 				if int(each_today_booking_data.FacilityID) == facilityIDInt {
// 					stats.TotalEarningsToday = each_today_booking_data.EarnedToday
// 					stats.NoOfBookingsToday = each_today_booking_data.NoOfBookingsToday
// 				}
// 			}
// 		}

// 		if len(FacilityTotalEarningsData) > 0 {
// 			for _, each_total_earning_data := range FacilityTotalEarningsData {
// 				//fmt.Println("each_total_earning_data : ", each_total_earning_data)
// 				if int(each_total_earning_data.FacilityID) == facilityIDInt {
// 					stats.TotalEarnings = each_total_earning_data.TotalEarnings
// 				}
// 			}
// 		}

// 		if len(FacilityLastBookingsData) > 0 {
// 			//fmt.Println("if loop")
// 			for _, each_last_booking_data := range FacilityLastBookingsData {
// 				//fmt.Println("each_last_booking_data : ", each_last_booking_data)
// 				if int(each_last_booking_data.FacilityID) == facilityIDInt {
// 					stats.LastBookingDate = each_last_booking_data.LastBookingDate
// 				}
// 			}
// 		}

// 		if len(FacilityTotalBookingsData) > 0 {
// 			for _, each_total_booking_data := range FacilityTotalBookingsData {
// 				//fmt.Println("each_total_booking_data : ", each_total_booking_data)
// 				if int(each_total_booking_data.FacilityID) == facilityIDInt {
// 					stats.NoOfBookings = each_total_booking_data.NoOfBookings
// 				}
// 			}
// 		}

// 		fmt.Println("stats after : ", stats)
// 		fmt.Println("type of stats : ", reflect.TypeOf(stats))

// 		finalData[facilityID] = &stats
// 		fmt.Println("finalData : ", finalData)
// 	}

// 	res = io.SuccessMessage(finalData, "Got facility stats")
// 	fmt.Println("res : ", res)
// 	fmt.Println("res.data  : ", res.Data)

// 	b.Logger.Log(gologger.Info,
// 		gologger.InternalServices,
// 		"getFacilityStats",
// 		"get facility stats response body",
// 		"getFacilityStatsResponse", "", structs.Map(res), true)
// 	return
// }

func (b *athFacilityService) CreateFacilitySportCategory(ctx context.Context, db *gorm.DB, u io.AthFacilitySportsData) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"createFacilitySportCategory",
		"create facility sport category request parameters",
		"createFacilitySportCategoryRequest", "", structs.Map(u), true)
	err := b.DbRepo.CreateFacilitySportCategory(ctx, db, u)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error creating/updating new facility sport category")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"createFacilitySportCategory",
			"create facility sport category request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}
	res = io.SuccessMessage(nil, "Facility sport category created/updated")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"createFacilitySportCategory",
		"create facility sport category response body",
		"createFacilitySportCategoryResponse", "", structs.Map(res), true)
	return
}

func (b *athFacilityService) GetCustomRatesForFacilityByID(ctx context.Context, db *gorm.DB, u io.CustomRatesForFacilityByID) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"GetCustomRatesForFacilityByID",
		"get custom rates for facility byID request parameters",
		"GetCustomRatesForFacilityByIDRequest", "", structs.Map(u), true)

	customRatesData, err := b.DbRepo.GetCustomRatesForFacilityByID(ctx, db, u)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error getting facility custom rates")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"GetCustomRatesForFacilityByID",
			"get custom rates for facility byID request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}

	if len(customRatesData) > 0 {
		res = io.SuccessMessage(customRatesData, "got Facility custom rates")
		b.Logger.Log(gologger.Info,
			gologger.InternalServices,
			"GetCustomRatesForFacilityByID",
			"get custom rates for facility byID response body",
			"GetCustomRatesForFacilityByIDResponse", "", structs.Map(res), true)
		return
	}
	res = io.SuccessMessage(nil, "got Facility custom rates")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"GetCustomRatesForFacilityByID",
		"get custom rates for facility byID response body",
		"GetCustomRatesForFacilityByIDResponse", "", structs.Map(res), true)

	return
}

func (b *athFacilityService) DeleteFacilityByID(ctx context.Context, db *gorm.DB, db1 *gorm.DB, u io.DeleteFacilityByID) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"DeleteFacilityByID",
		"delete facility by ID request parameters",
		"DeleteFacilityByIDRequest", "", structs.Map(u), true)

	deletedFacility, err := b.DbRepo.DeleteFacilityByID(ctx, db, db1, u)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error deleting facility by id")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"DeleteFacilityByID",
			"delete facility by ID request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}
	res = io.SuccessMessage(deletedFacility, "Facility Deleted/unavailable")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"DeleteFacilityByID",
		"delete facility by ID response body",
		"DeleteFacilityByIDResponse", "", structs.Map(res), true)
	return
}
