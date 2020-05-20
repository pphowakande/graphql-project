package handler

import (
	"context"

	"github.com/fatih/structs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"stash.bms.bz/bms/gologger.git"

	facilitypb "stash.bms.bz/turf/generic-proto-files.git/facility/v1"
)

func (h *athFacilityHandler) GetAllSportCategories(ctx context.Context, req *facilitypb.GetAllSportCategoriesReq) (*facilitypb.GetAllSportCategoriesRes, error) {
	h.logger.Log(gologger.Info,
		gologger.ExternalServices,
		"getAllSportCategories",
		"get all sport categories request parameters",
		"getAllSportCategoriesRequest", "", structs.Map(*req), true)
	storeErr := ""
	masterDB, err := h.DbRepo.GetDB(ctx, MasterDBConnectionName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	res := &facilitypb.GetAllSportCategoriesRes{}
	facilityServiceRes := h.athFacilityService.GetAllSportCategories(ctx, masterDB, *req)
	if facilityServiceRes.Error != nil {
		storeErr = facilityServiceRes.Error.Error()
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"getAllSportCategories",
			"get all sport categories request failed",
			gologger.ParseError, "", structs.Map(facilityServiceRes), true)
		return nil, status.Errorf(codes.Internal, storeErr)
	}

	if facilityServiceRes.Data != nil {
		sportCategoryData := facilityServiceRes.Data.([]*facilitypb.SportCategoryData)
		res = &facilitypb.GetAllSportCategoriesRes{
			Data: sportCategoryData,
		}
	} else {
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"getAllSportCategories",
			"get all sport categories request failed",
			gologger.ParseError, "", structs.Map(*res), true)
		return nil, status.Errorf(codes.NotFound, "no data")
	}
	h.logger.Log(gologger.Info,
		gologger.ExternalServices,
		"getAllSportCategories",
		"get all sport categories response body",
		"getAllSportCategoriesResponse", "", structs.Map(*res), true)
	return res, nil

}

func (h *athFacilityHandler) GetFacilityByID(ctx context.Context, req *facilitypb.GetFacilityByIDReq) (*facilitypb.GetFacilityByIDRes, error) {
	h.logger.Log(gologger.Info,
		gologger.ExternalServices,
		"getFacilityByID",
		"get facility byID request parameters",
		"getFacilityByIDRequest", "", structs.Map(*req), true)
	storeErr := ""
	res := &facilitypb.GetFacilityByIDRes{}

	if req.FacilityIds == "" {
		res = &facilitypb.GetFacilityByIDRes{}
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"getFacilityByID",
			"get facility byID request failed",
			gologger.ParseError, "", structs.Map(*req), true)
		return nil, status.Errorf(codes.InvalidArgument, "Parameter is missing")
	}
	masterDB, err := h.DbRepo.GetDB(ctx, MasterDBConnectionName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	facilityServiceRes := h.athFacilityService.GetFacilityByID(ctx, masterDB, *req)

	if facilityServiceRes.Error != nil {
		storeErr = facilityServiceRes.Error.Error()
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"getFacilityByID",
			"get facility byID request failed",
			gologger.ParseError, "", structs.Map(facilityServiceRes), true)
		return nil, status.Errorf(codes.Internal, storeErr)
	}

	facilityFinalList := facilityServiceRes.Data.([]*facilitypb.FacilityData)
	res = &facilitypb.GetFacilityByIDRes{
		Data: facilityFinalList,
	}
	h.logger.Log(gologger.Info,
		gologger.ExternalServices,
		"getFacilityByID",
		"get facility byID response body",
		"getFacilityByIDResponse", "", structs.Map(*res), true)
	return res, nil
}

func (h *athFacilityHandler) GetFacilityForVenueID(ctx context.Context, req *facilitypb.GetFacilityForVenueIDReq) (*facilitypb.GetFacilityForVenueIDRes, error) {
	h.logger.Log(gologger.Info,
		gologger.ExternalServices,
		"getFacilityForVenueID",
		"get facility for venueID request parameters",
		"getFacilityForVenueIDRequest", "", structs.Map(*req), true)
	storeErr := ""

	res := &facilitypb.GetFacilityForVenueIDRes{}

	if int(req.VenueId) == 0 {
		res = &facilitypb.GetFacilityForVenueIDRes{}
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"getFacilityForVenueID",
			"get facility for venueID request failed",
			gologger.ParseError, "", structs.Map(*res), true)
		return nil, status.Errorf(codes.InvalidArgument, "Parameter is missing")
	}
	masterDB, err := h.DbRepo.GetDB(ctx, MasterDBConnectionName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	facilityServiceRes := h.athFacilityService.GetFacilityForVenueID(ctx, masterDB, *req)

	if facilityServiceRes.Error != nil {
		storeErr = facilityServiceRes.Error.Error()
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"getFacilityForVenueID",
			"get facility for venueID request failed",
			gologger.ParseError, "", structs.Map(facilityServiceRes), true)
		return nil, status.Errorf(codes.Internal, storeErr)
	}

	if facilityServiceRes.Data != nil {
		FacilityIdsData := facilityServiceRes.Data.([]*facilitypb.FacilityIDData)
		res = &facilitypb.GetFacilityForVenueIDRes{
			Data: FacilityIdsData,
		}
	} else {
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"getFacilityForVenueID",
			"get facility for venueID request failed",
			gologger.ParseError, "", structs.Map(*res), true)
		return nil, status.Errorf(codes.NotFound, "no data")
	}
	h.logger.Log(gologger.Info,
		gologger.ExternalServices,
		"getFacilityForVenueID",
		"get facility for venueID response body",
		"getFacilityForVenueIDResponse", "", structs.Map(*res), true)
	return res, nil

}

// func (h *athFacilityHandler) GetFacilityStats(ctx context.Context, req *facilitypb.GetFacilityStatsReq) (*facilitypb.GetFacilityStatsRes, error) {
// 	fmt.Println("Inside GetFacilityStats handler function--------------")
// 	fmt.Println("req : ", req)
// 	h.logger.Log(gologger.Info,
// 		gologger.ExternalServices,
// 		"getFacilityStats",
// 		"get facility stats request parameters",
// 		"getFacilityStatsRequest", "", structs.Map(*req), true)
// 	storeErr := ""

// 	res := &facilitypb.GetFacilityStatsRes{}

// 	if req.FacilityIds == "" {
// 		res = &facilitypb.GetFacilityStatsRes{}
// 		h.logger.Log(gologger.Errsev3,
// 			gologger.ExternalServices,
// 			"getFacilityStats",
// 			"get facility stats request failed",
// 			gologger.ParseError, "", structs.Map(*req), true)
// 		return nil, status.Errorf(codes.InvalidArgument, "Parameter is missing")
// 	}
// masterDB, err := h.DbRepo.GetDB(ctx, MasterDBConnectionName)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, err.Error())
// 	}
// 	facilityServiceRes := h.athFacilityService.GetFacilityStats(ctx,masterDB, *req)

// 	if facilityServiceRes.Error != nil {
// 		storeErr = facilityServiceRes.Error.Error()
// 		h.logger.Log(gologger.Errsev3,
// 			gologger.ExternalServices,
// 			"getFacilityStats",
// 			"get facility stats request failed",
// 			gologger.ParseError, "", structs.Map(facilityServiceRes), true)
// 		return nil, status.Errorf(codes.Internal, storeErr)
// 	}

// 	if facilityServiceRes.Data != nil {
// 		FacilityStatsData := facilityServiceRes.Data.([]*facilitypb.FinalStats)
// 		res = &facilitypb.GetFacilityStatsRes{
// 			Data: FacilityStatsData,
// 		}
// 	} else {
// 		res = &facilitypb.GetFacilityStatsRes{
// 			Data: nil,
// 		}
// 	}
// 	h.logger.Log(gologger.Info,
// 		gologger.ExternalServices,
// 		"getFacilityStats",
// 		"get facility stats response body",
// 		"getFacilityStatsResponse", "", structs.Map(*res), true)
// 	return res, nil
// }
