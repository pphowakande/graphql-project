package handler

import (
	"context"
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

func (h *athFacilityHandler) CreateFacility(ctx context.Context, req *facilitypb.CreateFacilityRequest) (*facilitypb.FacilityData, error) {
	loggedInUserID, err := utile.GetUserIDFromContext(ctx)
	if err != nil {
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"CreateFacility",
			"create facility  request failed",
			gologger.ParseError, "", structs.Map(*req), true)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	h.logger.Log(gologger.Info,
		gologger.ExternalServices,
		"createFacility",
		"create facility request parameters",
		"createFacilityRequest", "", structs.Map(*req), true)
	facilityRequest := io.AthFacilities{
		VenueID:           int(req.VenueId),
		FacilityName:      req.FacilityName,
		FacilityBasePrice: req.DefaultRate,
		TimeSlot:          int(req.TimeSlot),
		IsActive:          true,
		Models: io.Models{
			CreatedBy: loggedInUserID,
			CreatedAt: int(time.Now().Unix()),
		},
	}

	masterTx, err := h.GetMasterDBTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		masterTx.RollbackUnlessCommitted()
	}()

	// get facility id from facility service
	facilityServiceRes := h.athFacilityService.CreateFacility(ctx, masterTx, facilityRequest)
	storeErr := ""
	res := &facilitypb.FacilityData{}

	if facilityServiceRes.Error != nil {
		storeErr = facilityServiceRes.Error.Error()
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"createFacility",
			"create facility request failed",
			gologger.ParseError, "", structs.Map(facilityServiceRes), true)
		return nil, status.Errorf(codes.Internal, storeErr)
	}

	faciityData := facilityServiceRes.Data.(io.AthFacilities)
	facilityID := faciityData.FacilityID

	// if week slots are provided
	for _, weekData := range req.WeekSlots {
		facilitySlotRequest := io.AthFacilitySlots{
			FacilityID:   facilityID,
			UserID:       loggedInUserID,
			SlotDays:     weekData.SlotDays,
			SlotType:     weekData.SlotType,
			SlotPrice:    weekData.SlotPrice,
			SlotFromTime: weekData.SlotStartTime,
			SlotToTime:   weekData.SlotEndTime,
			IsActive:     true,
			Models: io.Models{
				CreatedBy: loggedInUserID,
				CreatedAt: int(time.Now().Unix()),
			},
		}

		// add facility slots
		FacilitySlotsServiceRes := h.athFacilityService.CreateFacilitySlots(ctx, masterTx, facilitySlotRequest)
		if FacilitySlotsServiceRes.Error != nil {
			storeErr = FacilitySlotsServiceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"createFacility",
				"create facility request failed",
				gologger.ParseError, "", structs.Map(FacilitySlotsServiceRes), true)
			return nil, status.Errorf(codes.Internal, storeErr)
		}
	}

	// if custom rates are provided
	if req.CustomRates != nil {
		for _, eachRate := range req.CustomRates {
			facilityCustomRateRequest := io.AthFacilityCustomRates{
				FacilityID:   facilityID,
				SlotFromTime: eachRate.StartDate,
				SlotToTime:   eachRate.EndDate,
				UserID:       loggedInUserID,
				SlotPrice:    eachRate.Price,
				IsActive:     true,
				Available:    eachRate.Available,
				Models: io.Models{
					CreatedBy: loggedInUserID,
					CreatedAt: int(time.Now().Unix()),
				},
			}

			// store facility custom rates
			facilityCustomRatesServiceRes := h.athFacilityService.AddFacilityCustomRates(ctx, masterTx, facilityCustomRateRequest)

			if facilityCustomRatesServiceRes.Error != nil {
				storeErr = facilityCustomRatesServiceRes.Error.Error()
				h.logger.Log(gologger.Errsev3,
					gologger.ExternalServices,
					"CreateFacility",
					"create facility request failed",
					gologger.ParseError, "", structs.Map(facilityCustomRatesServiceRes), true)
				return nil, status.Errorf(codes.Internal, storeErr)
			}
		}
	}

	// save facility sport categories
	for _, val := range req.CategoryData {
		sportCategoryRequest := io.AthFacilitySportsData{
			SportCategoryID: int(val.CategoryId),
			FacilityID:      facilityID,
			Status:          true,
		}
		facilityCategoryServiceRes := h.athFacilityService.CreateFacilitySportCategory(ctx, masterTx, sportCategoryRequest)
		if facilityCategoryServiceRes.Error != nil {
			storeErr = facilityCategoryServiceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"createFacility",
				"create facility request failed",
				gologger.ParseError, "", structs.Map(facilityCategoryServiceRes), true)
			return nil, status.Errorf(codes.Internal, "error in CreateFacility: %v", storeErr)
		}
	}

	/** converting the int variable into a string using Itoa method */
	facilityIDStr := strconv.Itoa(facilityID)

	// get updated facility data from database
	facilitygetReq := facilitypb.GetFacilityByIDReq{
		FacilityIds: facilityIDStr,
	}
	facilitygetServiceRes := h.athFacilityService.GetFacilityByID(ctx, masterTx, facilitygetReq)
	if facilitygetServiceRes.Error != nil {
		storeErr = facilitygetServiceRes.Error.Error()
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"CreateFacility",
			"create facility request failed",
			gologger.ParseError, "", structs.Map(facilitygetServiceRes), true)
		return nil, status.Errorf(codes.Internal, storeErr)
	}

	facilityData := facilitygetServiceRes.Data.([]*facilitypb.FacilityData)
	if len(facilityData) > 0 {
		res = facilityData[0]
		h.logger.Log(gologger.Info,
			gologger.ExternalServices,
			"createFacility",
			"create facility response body",
			"createFacilityResponse", "", structs.Map(*res), true)
		return res, masterTx.Commit().Error
	}

	h.logger.Log(gologger.Errsev3,
		gologger.ExternalServices,
		"CreateFacility",
		"create facility request failed",
		gologger.ParseError, "", structs.Map(facilitygetServiceRes), true)
	return nil, status.Errorf(codes.Internal, storeErr)
}
