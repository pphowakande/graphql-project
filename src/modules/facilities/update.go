package handler

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/fatih/structs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"stash.bms.bz/bms/gologger.git"

	io "api/src/models"
	utile "api/src/utils"

	facilitypb "stash.bms.bz/turf/generic-proto-files.git/facility/v1"
)

func (h *athFacilityHandler) EditFacility(ctx context.Context, req *facilitypb.EditFacilityRequest) (*facilitypb.FacilityData, error) {
	loggedInUserID, err := utile.GetUserIDFromContext(ctx)
	if err != nil {
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"EditFacility",
			"edit facility  request failed",
			gologger.ParseError, "", structs.Map(*req), true)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	h.logger.Log(gologger.Info,
		gologger.ExternalServices,
		"editFacility",
		"edit facility request parameters",
		"editFacilityRequest", "", structs.Map(*req), true)
	facilityRequest := io.AthFacilities{
		FacilityName:      req.FacilityName,
		FacilityBasePrice: req.DefaultRate,
		FacilityID:        int(req.FacilityId),
		Models: io.Models{
			UpdatedBy: loggedInUserID,
			UpdatedAt: int(time.Now().Unix()),
		},
	}

	masterTx, err := h.GetMasterDBTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		masterTx.RollbackUnlessCommitted()
	}()

	// edit facility service
	facilityServiceRes := h.athFacilityService.EditFacility(ctx, masterTx, facilityRequest)
	storeErr := ""
	res := &facilitypb.FacilityData{}

	if facilityServiceRes.Error != nil {
		storeErr = facilityServiceRes.Error.Error()
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"editFacility",
			"edit facility request failed",
			gologger.ParseError, "", structs.Map(facilityServiceRes), true)
		return nil, status.Errorf(codes.Internal, storeErr)
	}

	// save/update facility sport categories
	for _, val := range req.CategoryData {
		sportCategoryRequest := io.AthFacilitySportsData{
			SportCategoryID: int(val.CategoryId),
			FacilityID:      int(req.FacilityId),
			Status:          val.Status,
		}

		facilityCategoryServiceRes := h.athFacilityService.CreateFacilitySportCategory(ctx, masterTx, sportCategoryRequest)
		if facilityCategoryServiceRes.Error != nil {
			storeErr = facilityCategoryServiceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"editFacility",
				"edit facility request failed",
				gologger.ParseError, "", structs.Map(facilityCategoryServiceRes), true)
			return nil, status.Errorf(codes.Internal, "error in create facility sport category: %v", storeErr)
		}
	}

	// if custom rates are provided
	if req.CustomRates != nil {
		for _, eachRate := range req.CustomRates {
			facilityCustomRateRequest := io.AthFacilityCustomRates{
				FacilityID:   int(req.FacilityId),
				SlotFromTime: eachRate.StartDate,
				SlotToTime:   eachRate.EndDate,
				UserID:       loggedInUserID,
				SlotPrice:    eachRate.Price,
				IsActive:     eachRate.Status,
				Available:    eachRate.Available,
				Models: io.Models{
					UpdatedBy: loggedInUserID,
					UpdatedAt: int(time.Now().Unix()),
				},
			}

			// store facility custom rates
			facilityCustomRatesServiceRes := h.athFacilityService.AddFacilityCustomRates(ctx, masterTx, facilityCustomRateRequest)

			if facilityCustomRatesServiceRes.Error != nil {
				storeErr = facilityCustomRatesServiceRes.Error.Error()
				h.logger.Log(gologger.Errsev3,
					gologger.ExternalServices,
					"editFacility",
					"edit facility request failed",
					gologger.ParseError, "", structs.Map(facilityCustomRatesServiceRes), true)
				return nil, status.Errorf(codes.Internal, storeErr)
			}
		}
	}

	//save/update facility week data
	for _, weekData := range req.WeekSlots {
		facilitySlotRequest := io.AthFacilitySlots{
			FacilityID:   int(req.FacilityId),
			UserID:       loggedInUserID,
			SlotDays:     weekData.SlotDays,
			SlotType:     weekData.SlotType,
			SlotPrice:    weekData.SlotPrice,
			SlotFromTime: weekData.SlotStartTime,
			IsActive:     weekData.Status,
			SlotToTime:   weekData.SlotEndTime,
			Models: io.Models{
				CreatedBy: loggedInUserID,
				UpdatedBy: loggedInUserID,
				UpdatedAt: int(time.Now().Unix()),
			},
		}

		// add facility slots
		FacilitySlotsServiceRes := h.athFacilityService.CreateFacilitySlots(ctx, masterTx, facilitySlotRequest)
		if FacilitySlotsServiceRes.Error != nil {
			storeErr = FacilitySlotsServiceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"editFacility",
				"edit facility request failed",
				gologger.ParseError, "", structs.Map(FacilitySlotsServiceRes), true)
			return nil, status.Errorf(codes.Internal, storeErr)
		}
	}

	/** converting the i1 variable into a string using Itoa method */
	facilityIdInt := strconv.Itoa(int(req.FacilityId))

	// get updated facility data from database
	facilitygetReq := facilitypb.GetFacilityByIDReq{
		FacilityIds: facilityIdInt,
	}

	fmt.Println("Calling GetFacilityByID service ---------------")
	facilitygetServiceRes := h.athFacilityService.GetFacilityByID(ctx, masterTx, facilitygetReq)

	if facilitygetServiceRes.Error != nil {
		storeErr = facilitygetServiceRes.Error.Error()
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"EditFacility",
			"edit facility request failed",
			gologger.ParseError, "", structs.Map(facilitygetServiceRes), true)
		return nil, status.Errorf(codes.Internal, storeErr)
	}

	facilityData := facilitygetServiceRes.Data.([]*facilitypb.FacilityData)
	fmt.Println("facilityData : ", facilityData)
	if len(facilityData) > 0 {
		res = facilityData[0]
		h.logger.Log(gologger.Info,
			gologger.ExternalServices,
			"EditFacility",
			"edit facility response body",
			"EditFacilityResponse", "", structs.Map(*res), true)
		return res, masterTx.Commit().Error
	}

	h.logger.Log(gologger.Errsev3,
		gologger.ExternalServices,
		"EditFacility",
		"edit facility request failed",
		gologger.ParseError, "", structs.Map(facilitygetServiceRes), true)
	return nil, status.Errorf(codes.Internal, "error in EditFacility :%v", storeErr)
}
