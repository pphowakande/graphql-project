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

func (h *athVenueHandler) EditVenue(ctx context.Context, req *venuepb.EditVenueRequest) (*venuepb.GetVenueByIDRes, error) {

	loggedInUserID, err := utile.GetUserIDFromContext(ctx)
	if err != nil {
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"EditVenue",
			"edit venue request failed",
			gologger.ParseError, "", structs.Map(*req), true)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	h.logger.Log(gologger.Info,
		gologger.ExternalServices,
		"editVenue",
		"edit venue request parameters",
		"editVenueRequest", "", structs.Map(*req), true)
	venueRequest := io.AthVenues{
		VenueID:     int(req.VenueId),
		VenueName:   req.Name,
		Description: req.Description,
		Email:       req.Email,
		Address:     req.Address,
		Phone:       req.Phone,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Models: io.Models{
			UpdatedBy: loggedInUserID,
			UpdatedAt: int(time.Now().Unix()),
		}}

	masterTx, err := h.GetMasterDBTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		masterTx.RollbackUnlessCommitted()
	}()

	venueServiceRes := h.athVenueService.EditVenue(ctx, masterTx, venueRequest)

	storeErr := ""
	res := &venuepb.GetVenueByIDRes{}

	if venueServiceRes.Error != nil {
		storeErr = venueServiceRes.Error.Error()
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"editVenue",
			"edit venue request failed",
			gologger.ParseError, "", structs.Map(venueServiceRes), true)
		return nil, status.Errorf(codes.Internal, storeErr)
	}

	// save amenities for venue
	for _, amenityData := range req.Amenities {
		venueAmenitiesRequest := io.AthVenueAmenitiesData{
			AmenityID: int(amenityData.AmenityId),
			VenueID:   int(req.VenueId),
			Status:    amenityData.Status,
		}

		venueAmenityServiceRes := h.athVenueService.CreateVenueAmenity(ctx, masterTx, venueAmenitiesRequest)
		if venueAmenityServiceRes.Error != nil {
			storeErr = venueAmenityServiceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"EditVenue",
				"edit venue request failed",
				gologger.ParseError, "", structs.Map(venueAmenityServiceRes), true)
			return nil, status.Errorf(codes.Internal, "error in editvenue: %v", storeErr)
		}
	}

	//save hours of operation
	for _, Daydata := range req.HoursOfOperation {
		for _, timing := range Daydata.Timing {
			models := io.Models{}
			if timing.Status == true {
				models = io.Models{
					CreatedBy: loggedInUserID,
				}
			} else {
				models = io.Models{
					UpdatedBy: loggedInUserID,
					UpdatedAt: int(time.Now().Unix()),
				}
			}
			hoursOfOpeRequest := io.AthVenueHours{
				Day:         Daydata.Day,
				VenueID:     int(req.VenueId),
				OpeningTime: timing.OpeningTime,
				ClosingTime: timing.ClosingTime,
				IsActive:    timing.Status,
				Models:      models,
			}
			hoursOfOperationRes := h.athVenueService.SaveHoursOfOperation(ctx, masterTx, hoursOfOpeRequest)
			if hoursOfOperationRes.Error != nil {
				storeErr = hoursOfOperationRes.Error.Error()
				h.logger.Log(gologger.Errsev3,
					gologger.ExternalServices,
					"EditVenue",
					"edit venue request failed",
					gologger.ParseError, "", structs.Map(hoursOfOperationRes), true)
				return nil, status.Errorf(codes.Internal, storeErr)
			}
		}
	}

	// save holidays
	for _, val := range req.Holidays {
		models := io.Models{}
		if val.Status == true {
			models = io.Models{
				CreatedBy: loggedInUserID,
			}
		} else {
			models = io.Models{
				UpdatedBy: loggedInUserID,
				UpdatedAt: int(time.Now().Unix()),
			}
		}
		venueholidayRequest := io.AthVenueHolidays{
			Title:    val.Title,
			VenueID:  int(req.VenueId),
			Date:     val.Date,
			IsActive: val.Status,
			Models:   models,
		}
		// get venue id from venue service
		venueHolidayRes := h.athVenueService.CreateVenueHoliday(ctx, masterTx, venueholidayRequest)

		if venueHolidayRes.Error != nil {
			storeErr = venueHolidayRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"EditVenue",
				"edit venue request failed",
				gologger.ParseError, "", structs.Map(venueHolidayRes), true)
			return nil, status.Errorf(codes.Internal, storeErr)
		}
	}

	// upload and validate venue images
	if req.Images != nil {
		var galleryImgData []io.ImgData
		for _, val := range req.Images.GalleryImg {
			var eachGalleryImgData io.ImgData
			eachGalleryImgData.Image = val.ImgUrl
			eachGalleryImgData.Status = val.Status
			galleryImgData = append(galleryImgData, eachGalleryImgData)
		}

		var thumbImgData []io.ImgData
		for _, val := range req.Images.ThumbnailImg {
			var eachThumbImgData io.ImgData
			eachThumbImgData.Image = val.ImgUrl
			eachThumbImgData.Status = val.Status
			thumbImgData = append(thumbImgData, eachThumbImgData)
		}

		var headerImgData []io.ImgData
		for _, val := range req.Images.HeaderImg {
			var eachHeaderImgData io.ImgData
			eachHeaderImgData.Image = val.ImgUrl
			eachHeaderImgData.Status = val.Status
			headerImgData = append(headerImgData, eachHeaderImgData)
		}

		venueImagesReq := io.VenueImagesReq{
			HeaderImg:    headerImgData,
			ThumbnailImg: thumbImgData,
			GalleryImg:   galleryImgData,
			VenueID:      int(req.VenueId),
		}
		// first validate uploaded documents and check if document is already uploaded
		validateServiceRes := h.athUploadService.ValidateVenueImages(ctx, venueImagesReq)
		if validateServiceRes.Error != nil {
			storeErr = validateServiceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"EditVenue",
				"edit venue request failed",
				gologger.ParseError, "", structs.Map(validateServiceRes), true)
			return nil, status.Errorf(codes.Internal, "error in EditVenue :%v", storeErr)
		}

		validatedData := validateServiceRes.Data.(map[string]interface{})
		// upload documents to google cloud
		// if documents are validated , upload it to google cloud
		uploadServiceRes := h.athUploadService.VenueUploadToS3(ctx, validatedData)
		if uploadServiceRes.Error != nil {
			storeErr = uploadServiceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"EditVenue",
				"edit venue request failed",
				gologger.ParseError, "", structs.Map(uploadServiceRes), true)
			return nil, status.Errorf(codes.Internal, "error in EditVenue :%v", storeErr)
		}

		imgUploadedData := uploadServiceRes.Data.(io.VenueImagesReq)
		imgUploadedData.VenueID = int(req.VenueId)
		imgUploadedData.CreatedBy = loggedInUserID

		// save venue images
		venueAmenityServiceRes := h.athVenueService.CreateVenueImage(ctx, masterTx, imgUploadedData)
		if venueAmenityServiceRes.Error != nil {
			storeErr = venueAmenityServiceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"EditVenue",
				"edit venue request failed",
				gologger.ParseError, "", structs.Map(venueAmenityServiceRes), true)
			return nil, status.Errorf(codes.Internal, "error in EditVenue: %v", storeErr)
		}
	}

	// get updated data from database
	venueServiceRes = h.athVenueService.GetVenueByID(ctx, masterTx, int(req.VenueId))
	if venueServiceRes.Error != nil {
		storeErr = venueServiceRes.Error.Error()
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"EditVenue",
			"edit venue request failed",
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
		Images:           &imageData,
		Amenities:        VenueData["amenities"].([]*venuepb.AmenityData),
		Holidays:         VenueData["holidays"].([]*venuepb.HolidaysData),
		HoursOfOperation: VenueData["hoursOfOperation"].([]*venuepb.HoursOfOperationData),
		Address:          VenueData["address"].(string),
		Latitude:         VenueData["latitude"].(float32),
		Longitude:        VenueData["longitude"].(float32),
	}

	h.logger.Log(gologger.Info,
		gologger.ExternalServices,
		"EditVenue",
		"edit venue response body",
		"editVenueResponse", "", structs.Map(*res), true)
	return res, masterTx.Commit().Error
}
