package handler

import (
	io "api/src/models"
	utile "api/src/utils"
	"context"
	"time"

	"github.com/fatih/structs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"stash.bms.bz/bms/gologger.git"

	venuepb "stash.bms.bz/turf/generic-proto-files.git/venue/v1"
)

func (h *athVenueHandler) CreateVenue(ctx context.Context, req *venuepb.CreateVenueRequest) (*venuepb.GetVenueByIDRes, error) {
	loggedInUserID, err := utile.GetUserIDFromContext(ctx)
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
		"createVenue",
		"create venue request parameters",
		"createVenueRequest", "", structs.Map(*req), true)

	storeErr := ""
	res := &venuepb.GetVenueByIDRes{}

	venueRequest := io.AthVenues{
		MerchantID:  loggedInUserID,
		VenueName:   req.Name,
		Description: req.Description,
		Email:       req.Email,
		Phone:       req.Phone,
		Address:     req.Address,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		IsActive:    true,
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

	// get venue id from venue service
	venueServiceRes := h.athVenueService.CreateVenue(ctx, masterTx, venueRequest)

	if venueServiceRes.Error != nil {
		storeErr = venueServiceRes.Error.Error()
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"createVenue",
			"create venue request failed",
			gologger.ParseError, "", structs.Map(venueServiceRes), true)
		return nil, status.Errorf(codes.Internal, "error in create venue: %v", storeErr)
	}

	venueData := venueServiceRes.Data.(map[string]interface{})
	venueID := venueData["venue_id"].(int)

	// upload and validate venue images
	if req.Images != nil {
		var galleryImgData []io.ImgData
		for _, val := range req.Images.GalleryImg {
			var eachGalleryImgData io.ImgData
			eachGalleryImgData.Image = val.ImgUrl
			eachGalleryImgData.Status = true
			galleryImgData = append(galleryImgData, eachGalleryImgData)
		}

		var thumbImgData []io.ImgData
		for _, val := range req.Images.ThumbnailImg {
			var eachThumbImgData io.ImgData
			eachThumbImgData.Image = val.ImgUrl
			eachThumbImgData.Status = true
			thumbImgData = append(thumbImgData, eachThumbImgData)
		}

		var headerImgData []io.ImgData
		for _, val := range req.Images.HeaderImg {
			var eachHeaderImgData io.ImgData
			eachHeaderImgData.Image = val.ImgUrl
			eachHeaderImgData.Status = true
			headerImgData = append(headerImgData, eachHeaderImgData)
		}

		venueImagesReq := io.VenueImagesReq{
			HeaderImg:    headerImgData,
			ThumbnailImg: thumbImgData,
			GalleryImg:   galleryImgData,
			VenueID:      venueID,
		}
		// first validate uploaded documents and check if document is already uploaded
		validateServiceRes := h.athUploadService.ValidateVenueImages(ctx, venueImagesReq)
		if validateServiceRes.Error != nil {
			storeErr = validateServiceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"create venue ",
				"create venue request failed",
				gologger.ParseError, "", structs.Map(validateServiceRes), true)
			return nil, status.Errorf(codes.Internal, "error in createVenue  :%v", storeErr)
		}

		validatedData := validateServiceRes.Data.(map[string]interface{})
		// upload documents to google cloud
		// if documents are validated , upload it to google cloud
		uploadServiceRes := h.athUploadService.VenueUploadToS3(ctx, validatedData)
		if uploadServiceRes.Error != nil {
			storeErr = uploadServiceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"createVenue",
				"create venue  request failed",
				gologger.ParseError, "", structs.Map(uploadServiceRes), true)
			return nil, status.Errorf(codes.Internal, "error in createVenue  :%v", storeErr)
		}

		imgUploadedData := uploadServiceRes.Data.(io.VenueImagesReq)
		imgUploadedData.VenueID = venueID
		imgUploadedData.CreatedBy = loggedInUserID

		// save venue images
		venueAmenityServiceRes := h.athVenueService.CreateVenueImage(ctx, masterTx, imgUploadedData)
		if venueAmenityServiceRes.Error != nil {
			storeErr = venueAmenityServiceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"createVenue",
				"create venue request failed",
				gologger.ParseError, "", structs.Map(venueAmenityServiceRes), true)
			return nil, status.Errorf(codes.Internal, "error in createVenue: %v", storeErr)
		}
	}

	// save amenities for venue
	for _, amenityData := range req.Amenities {
		venueAmenitiesRequest := io.AthVenueAmenitiesData{
			AmenityID: int(amenityData.AmenityId),
			VenueID:   venueID,
			Status:    true,
		}

		venueAmenityServiceRes := h.athVenueService.CreateVenueAmenity(ctx, masterTx, venueAmenitiesRequest)
		if venueAmenityServiceRes.Error != nil {
			storeErr = venueAmenityServiceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"createVenue",
				"create venue request failed",
				gologger.ParseError, "", structs.Map(venueAmenityServiceRes), true)
			return nil, status.Errorf(codes.Internal, "error in create venue: %v", storeErr)
		}
	}

	//save hours of operation
	for _, Daydata := range req.HoursOfOperation {
		for _, timing := range Daydata.Timing {
			hoursOfOpeRequest := io.AthVenueHours{
				Day:         Daydata.Day,
				VenueID:     venueID,
				OpeningTime: timing.OpeningTime,
				ClosingTime: timing.ClosingTime,
				IsActive:    true,
				Models: io.Models{
					CreatedBy: loggedInUserID,
				},
			}
			hoursOfOperationRes := h.athVenueService.SaveHoursOfOperation(ctx, masterTx, hoursOfOpeRequest)
			if hoursOfOperationRes.Error != nil {
				storeErr = hoursOfOperationRes.Error.Error()
				h.logger.Log(gologger.Errsev3,
					gologger.ExternalServices,
					"createVenue",
					"create venue request failed",
					gologger.ParseError, "", structs.Map(hoursOfOperationRes), true)
				return nil, status.Errorf(codes.Internal, storeErr)
			}
		}
	}

	// save holidays
	for _, val := range req.Holidays {
		venueholidayRequest := io.AthVenueHolidays{
			Title:    val.Title,
			VenueID:  venueID,
			Date:     val.Date,
			IsActive: true,
		}
		// get venue id from venue service
		venueHolidayRes := h.athVenueService.CreateVenueHoliday(ctx, masterTx, venueholidayRequest)

		if venueHolidayRes.Error != nil {
			storeErr = venueHolidayRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"createVenue",
				"create venue request failed",
				gologger.ParseError, "", structs.Map(venueHolidayRes), true)
			return nil, status.Errorf(codes.Internal, storeErr)
		}
	}

	// get updated data from database
	venueServiceRes = h.athVenueService.GetVenueByID(ctx, masterTx, venueID)
	if venueServiceRes.Error != nil {
		storeErr = venueServiceRes.Error.Error()
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"createVenue",
			"create venue request failed",
			gologger.ParseError, "", structs.Map(venueServiceRes), true)
		return nil, status.Errorf(codes.Internal, storeErr)
	}

	VenueData := venueServiceRes.Data.(map[string]interface{})
	imageData := VenueData["images"].(venuepb.CreateImageData)
	res = &venuepb.GetVenueByIDRes{
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
		"createVenue",
		"create venue response body",
		"createVenueResponse", "", structs.Map(*res), true)
	return res, masterTx.Commit().Error
}
