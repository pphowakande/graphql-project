package handler

import (
	io "api/src/models"
	utile "api/src/utils"
	"context"
	"fmt"

	"github.com/fatih/structs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"stash.bms.bz/bms/gologger.git"

	venuepb "stash.bms.bz/turf/generic-proto-files.git/venue/v1"
)

func (h *athVenueHandler) GetVenueByID(ctx context.Context, req *venuepb.GetVenueByIDReq) (*venuepb.GetVenueByIDRes, error) {
	_, err := utile.GetUserIDFromContext(ctx)
	if err != nil {
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"CreateVenue",
			"create venue request failed",
			gologger.ParseError, "", structs.Map(*req), true)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	h.logger.Log(gologger.Info,
		gologger.ExternalServices,
		"getVenueByID",
		"get venue byID request parameters",
		"getVenueByIDRequest", "", structs.Map(*req), true)
	storeErr := ""

	if int(req.VenueId) == 0 {
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"getVenueByID",
			"get venue byID request failed",
			gologger.ParseError, "", structs.Map(*req), true)
		return nil, status.Errorf(codes.InvalidArgument, "Parameter is missing")
	}

	masterDB, err := h.DbRepo.GetDB(ctx, MasterDBConnectionName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Parameter is missing")
	}

	venueServiceRes := h.athVenueService.GetVenueByID(ctx, masterDB, int(req.VenueId))
	if venueServiceRes.Error != nil {
		storeErr = venueServiceRes.Error.Error()
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"getVenueByID",
			"get venue byID request failed",
			gologger.ParseError, "", structs.Map(venueServiceRes), true)
		return nil, status.Errorf(codes.Internal, storeErr)
	}

	VenueData := venueServiceRes.Data.(map[string]interface{})
	imageData := VenueData["images"].(venuepb.CreateImageData)
	res := &venuepb.GetVenueByIDRes{
		VenueId:          int32(VenueData["venueId"].(int)),
		Name:             VenueData["name"].(string),
		Description:      VenueData["description"].(string),
		Email:            VenueData["email"].(string),
		Phone:            VenueData["phone"].(string),
		Amenities:        VenueData["amenities"].([]*venuepb.AmenityData),
		Images:           &imageData,
		Holidays:         VenueData["holidays"].([]*venuepb.HolidaysData),
		HoursOfOperation: VenueData["hoursOfOperation"].([]*venuepb.HoursOfOperationData),
		Address:          VenueData["address"].(string),
		Latitude:         VenueData["latitude"].(float32),
		Longitude:        VenueData["longitude"].(float32),
	}

	h.logger.Log(gologger.Info,
		gologger.ExternalServices,
		"getVenueByID",
		"get venue byID response body",
		"getVenueByIDResponse", "", structs.Map(*res), true)
	return res, nil
}

func (h *athVenueHandler) GetAllAmenities(ctx context.Context, req *venuepb.GetAllAmenitiesReq) (*venuepb.GetAllAmenitiesRes, error) {
	_, err := utile.GetUserIDFromContext(ctx)
	if err != nil {
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"CreateVenue",
			"create venue request failed",
			gologger.ParseError, "", structs.Map(*req), true)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	h.logger.Log(gologger.Info,
		gologger.ExternalServices,
		"GetAllAmenities",
		"get all amenities request parameters",
		"getAllAmenitiesByIDRequest", "", structs.Map(*req), true)
	storeErr := ""

	res := &venuepb.GetAllAmenitiesRes{}

	fmt.Println("req : ", req)

	fmt.Println("MasterDBConnectionName : ", MasterDBConnectionName)
	fmt.Println("h : ", h)
	fmt.Println("h.repo : ", h.DbRepo)
	masterDB, err := h.DbRepo.GetDB(ctx, MasterDBConnectionName)
	fmt.Println("err : ", err)
	fmt.Println("masterDB : ", masterDB)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	venueServiceRes := h.athVenueService.GetAllAmenities(ctx, masterDB, *req)
	fmt.Println("venueServiceRes : ", venueServiceRes)
	if venueServiceRes.Error != nil {
		storeErr = venueServiceRes.Error.Error()
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"GetAllAmenities",
			"get all amenities request failed",
			gologger.ParseError, "", structs.Map(venueServiceRes), true)
		return nil, status.Errorf(codes.Internal, storeErr)
	}

	if venueServiceRes.Data != nil {
		amenityList := venueServiceRes.Data.([]*venuepb.AmenityData)
		fmt.Println("amenityList returning : ", amenityList)
		res = &venuepb.GetAllAmenitiesRes{
			AmenityData: amenityList,
		}
	} else {
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"GetAllAmenities",
			"get all amenities request failed",
			gologger.ParseError, "", structs.Map(*res), true)
		return nil, status.Errorf(codes.NotFound, "No amenities present")
	}
	h.logger.Log(gologger.Info,
		gologger.ExternalServices,
		"GetAllAmenities",
		"get all amenities response body",
		"GetAllAmenitiesResponse", "", structs.Map(*res), true)
	return res, nil
}

func (h *athVenueHandler) GetListOfVenueByMerchantID(ctx context.Context, req *venuepb.GetListOfVenueByMerchantIDReq) (*venuepb.GetListOfVenueByMerchantIDRes, error) {
	h.logger.Log(gologger.Info,
		gologger.ExternalServices,
		"GetListOfVenueByMerchantID",
		"get all veneues for merchant request parameters",
		"GetListOfVenueByMerchantIDRequest", "", structs.Map(*req), true)
	storeErr := ""

	res := &venuepb.GetListOfVenueByMerchantIDRes{}
	masterDB, err := h.DbRepo.GetDB(ctx, MasterDBConnectionName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Parameter is missing")
	}

	venueServiceRes := h.athVenueService.GetListOfVenueByMerchantID(ctx, masterDB, *req)
	if venueServiceRes.Error != nil {
		storeErr = venueServiceRes.Error.Error()
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"GetListOfVenueByMerchantID",
			"get all veneues for merchant request failed",
			gologger.ParseError, "", structs.Map(venueServiceRes), true)
		return nil, status.Errorf(codes.Internal, storeErr)
	}

	if venueServiceRes.Data != nil {
		venuesData := venueServiceRes.Data.([]io.AthVenues)
		if len(venuesData) > 0 {
			venueList := make([]*venuepb.VenueList, 0)
			for _, eachVenue := range venuesData {
				var venueData venuepb.VenueList

				venueData.VenueId = int32(eachVenue.VenueID)
				venueData.VenueName = eachVenue.VenueName
				venueList = append(venueList, &venueData)
			}
			res = &venuepb.GetListOfVenueByMerchantIDRes{
				VenueData: venueList,
			}
		}
	} else {
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"GetListOfVenueByMerchantID",
			"get all veneues for merchant request failed",
			gologger.ParseError, "", structs.Map(*res), true)
		return nil, status.Errorf(codes.NotFound, "No venues present")
	}
	h.logger.Log(gologger.Info,
		gologger.ExternalServices,
		"GetListOfVenueByMerchantID",
		"get all veneues for merchant response body",
		"GetListOfVenueByMerchantIDResponse", "", structs.Map(*res), true)
	return res, nil
}
