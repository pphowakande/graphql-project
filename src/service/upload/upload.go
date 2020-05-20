package upload

import (
	db "api/src/dal"
	"api/src/dal/merchant"
	"api/src/dal/venue"
	io "api/src/models"
	utile "api/src/utils"
	"context"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"

	"github.com/fatih/structs"
	"stash.bms.bz/bms/gologger.git"
)

type AthUploadService interface {
	UploadToS3(ctx context.Context, u map[string]interface{}) (res io.Response)
	VenueUploadToS3(ctx context.Context, u map[string]interface{}) (res io.Response)

	ValidateMerchantDocs(ctx context.Context, u io.MerchantKYC) (res io.Response)
	ValidateVenueImages(ctx context.Context, u io.VenueImagesReq) (res io.Response)
}

type athUploadService struct {
	Logger    *gologger.Logger
	DbRepo    merchant.Repository
	VenueRepo venue.Repository
}

func NewBasicAthUploadService(logger *gologger.Logger, DbRepo merchant.Repository, VenueRepo venue.Repository) AthUploadService {
	return &athUploadService{
		Logger:    logger,
		DbRepo:    DbRepo,
		VenueRepo: VenueRepo,
	}
}

func (b *athUploadService) ValidateMerchantDocs(ctx context.Context, u io.MerchantKYC) (res io.Response) {
	// check if all documents are valid base64 string
	var gstFileData io.MerchantDocFileData
	var bankFileData io.MerchantDocFileData
	var addFileData io.MerchantDocFileData
	var panFileData io.MerchantDocFileData

	dal, err := b.DbRepo.GetDB(ctx, db.MasterDBConnectionName)
	if err != nil {
		res = io.FailureMessage(res.Error, "Can not update GST document.GST document is already uploaded and verified by BookMyShow")
		res.Error = errors.New("Can not update GST document.GST document is already uploaded and verified by BookMyShow")
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"ValidateMerchantDocs",
			"Validate merchant documents request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}

	if u.GstNoFile != "" {
		gstFileData.ValidatedBase64, gstFileData.ContentType, gstFileData.Verified = utile.ValidateBase64(u.GstNoFile, "merchant")
		if gstFileData.Verified == false {
			res = io.FailureMessage(errors.New("GST document Format is invalid"), "GST document Format is invalid")
			res.Error = errors.New("GST document Format is invalid")
			b.Logger.Log(gologger.Errsev3,
				gologger.InternalServices,
				"ValidateMerchantDocs",
				"Validate merchant documents request failed",
				gologger.ParseError, "", structs.Map(res), true)
			return
		}
	} else {
		gstFileData.Verified = true
	}
	if u.BankAccFile != "" {
		bankFileData.ValidatedBase64, bankFileData.ContentType, bankFileData.Verified = utile.ValidateBase64(u.BankAccFile, "merchant")
		if bankFileData.Verified == false {
			res = io.FailureMessage(errors.New("Bank document Format is invalid"), "Bank document Format is invalid")
			res.Error = errors.New("Bank document Format is invalid")
			b.Logger.Log(gologger.Errsev3,
				gologger.InternalServices,
				"ValidateMerchantDocs",
				"Validate merchant documents request failed",
				gologger.ParseError, "", structs.Map(res), true)
			return
		}
	} else {
		bankFileData.Verified = true
	}
	if u.AddressFile != "" {
		addFileData.ValidatedBase64, addFileData.ContentType, addFileData.Verified = utile.ValidateBase64(u.AddressFile, "merchant")
		if addFileData.Verified == false {
			res = io.FailureMessage(errors.New("Address document Format is invalid"), "Address document Format is invalid")
			res.Error = errors.New("Address document Format is invalid")
			b.Logger.Log(gologger.Errsev3,
				gologger.InternalServices,
				"ValidateMerchantDocs",
				"Validate merchant documents request failed",
				gologger.ParseError, "", structs.Map(res), true)
			return
		}
	} else {
		addFileData.Verified = true
	}
	if u.PanNoFile != "" {
		panFileData.ValidatedBase64, panFileData.ContentType, panFileData.Verified = utile.ValidateBase64(u.PanNoFile, "merchant")
		if panFileData.Verified == false {
			res = io.FailureMessage(errors.New("PAN document Format is invalid"), "PAN document Format is invalid")
			res.Error = errors.New("PAN document Format is invalid")
			b.Logger.Log(gologger.Errsev3,
				gologger.InternalServices,
				"ValidateMerchantDocs",
				"Validate merchant documents request failed",
				gologger.ParseError, "", structs.Map(res), true)
			return
		}
	} else {
		panFileData.Verified = true
	}

	merchantIdstr := strconv.Itoa(u.MerchantID)

	// all uploaded documents are verified now. Lets check
	// if document already exists in database and not verified by bms
	if gstFileData.ValidatedBase64 != "" {
		// check if image already exists in db
		imgExists, _ := b.DbRepo.CheckMerchantDocExists(ctx, dal, "gstNoFile", u.MerchantID)
		if imgExists == true {
			res = io.FailureMessage(res.Error, "Can not update GST document.GST document is already uploaded and verified by BookMyShow")
			res.Error = errors.New("Can not update GST document.GST document is already uploaded and verified by BookMyShow")
			b.Logger.Log(gologger.Errsev3,
				gologger.InternalServices,
				"ValidateMerchantDocs",
				"Validate merchant documents request failed",
				gologger.ParseError, "", structs.Map(res), true)
			return
		}
		// convert base64 image to blob
		blobImg, err1 := base64.StdEncoding.DecodeString(gstFileData.ValidatedBase64)
		if err1 != nil {
			panic(err1)
		}
		gstFileData.Blob = string(blobImg)
		gstFileData.Exists = imgExists
		fileExt := ""
		splittedContentType := strings.Split(gstFileData.ContentType, "/")
		fileExt = splittedContentType[1]
		gstFileData.FileName = "gstcertificate_0" + merchantIdstr + "." + fileExt
	}

	if bankFileData.ValidatedBase64 != "" {
		// convert base64 image to blob
		blobImg, err1 := base64.StdEncoding.DecodeString(bankFileData.ValidatedBase64)
		if err1 != nil {
			panic(err1)
		}
		// check if image already exists in db
		imgExists, _ := b.DbRepo.CheckMerchantDocExists(ctx, dal, "bankAccFile", u.MerchantID)
		if imgExists == true {
			res = io.FailureMessage(res.Error, "Can not update Bank document.Bank document is already uploaded and verified by BookMyShow")
			res.Error = errors.New("Can not update Bank document.Bank document is already uploaded and verified by BookMyShow")
			b.Logger.Log(gologger.Errsev3,
				gologger.InternalServices,
				"ValidateMerchantDocs",
				"Validate merchant documents request failed",
				gologger.ParseError, "", structs.Map(res), true)
			return
		}
		bankFileData.Blob = string(blobImg)
		bankFileData.Exists = imgExists
		fileExt := ""
		splittedContentType := strings.Split(bankFileData.ContentType, "/")
		fileExt = splittedContentType[1]
		bankFileData.FileName = "bankcertificate_0" + merchantIdstr + "." + fileExt

	}

	if addFileData.ValidatedBase64 != "" {
		// convert base64 image to blob
		blobImg, err1 := base64.StdEncoding.DecodeString(addFileData.ValidatedBase64)
		if err1 != nil {
			panic(err1)
		}
		// check if image already exists in db
		imgExists, _ := b.DbRepo.CheckMerchantDocExists(ctx, dal, "addressFile", u.MerchantID)
		if imgExists == true {
			res = io.FailureMessage(res.Error, "Can not update Address document.Address document is already uploaded and verified by BookMyShow")
			res.Error = errors.New("Can not update Address document.Address document is already uploaded and verified by BookMyShow")
			b.Logger.Log(gologger.Errsev3,
				gologger.InternalServices,
				"ValidateMerchantDocs",
				"Validate merchant documents request failed",
				gologger.ParseError, "", structs.Map(res), true)
			return
		}
		addFileData.Blob = string(blobImg)
		addFileData.Exists = imgExists
		fileExt := ""
		splittedContentType := strings.Split(addFileData.ContentType, "/")
		fileExt = splittedContentType[1]
		addFileData.FileName = "addcertificate_0" + merchantIdstr + "." + fileExt

	}

	if panFileData.ValidatedBase64 != "" {
		// convert base64 image to blob
		blobImg, err1 := base64.StdEncoding.DecodeString(panFileData.ValidatedBase64)
		if err1 != nil {
			panic(err1)
		}
		// check if image already exists in db
		imgExists, _ := b.DbRepo.CheckMerchantDocExists(ctx, dal, "panNoFile", u.MerchantID)
		if imgExists == true {
			res = io.FailureMessage(res.Error, "Can not update PAN document.PAN document is already uploaded and verified by BookMyShow")
			res.Error = errors.New("Can not update PAN document.PAN document is already uploaded and verified by BookMyShow")
			b.Logger.Log(gologger.Errsev3,
				gologger.InternalServices,
				"ValidateMerchantDocs",
				"Validate merchant documents request failed",
				gologger.ParseError, "", structs.Map(res), true)
			return
		}
		panFileData.Blob = string(blobImg)
		panFileData.Exists = imgExists
		fileExt := ""
		splittedContentType := strings.Split(panFileData.ContentType, "/")
		fileExt = splittedContentType[1]
		panFileData.FileName = "pancertificate_0" + merchantIdstr + "." + fileExt

	}

	data := make(map[string]interface{})
	data["gstFileData"] = gstFileData
	data["bankFileData"] = bankFileData
	data["addFileData"] = addFileData
	data["panFileData"] = panFileData

	res = io.SuccessMessage(data, "Uploaded validated")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"ValidateMerchantDocs",
		"Validate merchant documents response body",
		"ValidateMerchantDocsResponse", "", structs.Map(res), true)
	return
}

func (b *athUploadService) ValidateVenueImages(ctx context.Context, u io.VenueImagesReq) (res io.Response) {
	// check if all documents are valid base64 string
	var headerImgData []io.VenueImageData
	var thumbImgData []io.VenueImageData
	var galleryImgData []io.VenueImageData

	venueIdStr := strconv.Itoa(u.VenueID)
	dal, err := b.DbRepo.GetDB(ctx, db.MasterDBConnectionName)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error getting venue images details")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"getVenueHoursByID",
			"get venue hours byID request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}

	// get venue images for venue
	VenueImg, err := b.VenueRepo.GetVenueImagesById(ctx, dal, u.VenueID, false)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error getting venue images details")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"ValidateVenueImages",
			"validate venue images request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}
	totalimg := len(VenueImg)
	totalimgStr := strconv.Itoa(totalimg)

	if len(u.HeaderImg) > 0 {
		for _, eachHeaderImg := range u.HeaderImg {
			var eachHeaderImgData io.VenueImageData
			if eachHeaderImg.Status == true {
				// check if image url already exists in db. If yes, then dont do anything.
				// If url does not present, then consider it as base64
				// get venue images for venue
				exists, err := b.VenueRepo.CheckVenueImgExists(ctx, dal, eachHeaderImg.Image)
				if err != nil {
					res = io.FailureMessage(res.Error, "Error getting venue images details")
					res.Error = err
					b.Logger.Log(gologger.Errsev3,
						gologger.InternalServices,
						"ValidateVenueImages",
						"Validate venue images request failed",
						gologger.ParseError, "", structs.Map(res), true)
					return
				}
				if exists == false {
					eachHeaderImgData.ValidatedBase64, eachHeaderImgData.ContentType, eachHeaderImgData.Verified = utile.ValidateBase64(eachHeaderImg.Image, "venue")
					if eachHeaderImgData.Verified == false {
						res = io.FailureMessage(errors.New("Header Image Format is invalid"), "Header Image Format is invalid")
						res.Error = errors.New("Header Image Format is invalid")
						b.Logger.Log(gologger.Errsev3,
							gologger.InternalServices,
							"ValidateVenueImages",
							"Validate venue images request failed",
							gologger.ParseError, "", structs.Map(res), true)
						return
					}
					// convert base64 image to blob
					blobImg, err1 := base64.StdEncoding.DecodeString(eachHeaderImgData.ValidatedBase64)
					if err1 != nil {
						res = io.FailureMessage(errors.New("Error decoding image"), "Error decoding image")
						res.Error = errors.New("Error decoding image")
						b.Logger.Log(gologger.Errsev3,
							gologger.InternalServices,
							"ValidateVenueImages",
							"Validate venue images request failed",
							gologger.ParseError, "", structs.Map(res), true)
						return
					}
					eachHeaderImgData.Blob = string(blobImg)
					eachHeaderImgData.Exists = false
					fileExt := ""
					splittedContentType := strings.Split(eachHeaderImgData.ContentType, "/")
					fileExt = splittedContentType[1]
					eachHeaderImgData.FileName = "header_img_" + venueIdStr + "." + fileExt
					headerImgData = append(headerImgData, eachHeaderImgData)
				}
			} else {
				eachHeaderImgData.Verified = true
				eachHeaderImgData.ImageUrl = eachHeaderImg.Image
				eachHeaderImgData.Exists = true
				eachHeaderImgData.Blob = ""
				headerImgData = append(headerImgData, eachHeaderImgData)
			}
		}
	}

	if len(u.ThumbnailImg) > 0 {
		for _, eachThumbImg := range u.ThumbnailImg {
			var eachThumbImgData io.VenueImageData
			if eachThumbImg.Status == true {
				// check if image url already exists in db. If yes, then dont do anything.
				// If url does not present, then consider it as base64
				// get venue images for venue
				exists, err := b.VenueRepo.CheckVenueImgExists(ctx, dal, eachThumbImg.Image)
				if err != nil {
					res = io.FailureMessage(res.Error, "Error getting venue images details")
					res.Error = err
					b.Logger.Log(gologger.Errsev3,
						gologger.InternalServices,
						"ValidateVenueImages",
						"Validate venue images request failed",
						gologger.ParseError, "", structs.Map(res), true)
					return
				}
				if exists == false {
					eachThumbImgData.ValidatedBase64, eachThumbImgData.ContentType, eachThumbImgData.Verified = utile.ValidateBase64(eachThumbImg.Image, "venue")
					if eachThumbImgData.Verified == false {
						res = io.FailureMessage(errors.New("Thumbnail Image Format is invalid"), "Thumbnail Image Format is invalid")
						res.Error = errors.New("Thumbnail Image Format is invalid")
						b.Logger.Log(gologger.Errsev3,
							gologger.InternalServices,
							"ValidateVenueImages",
							"Validate venue images request failed",
							gologger.ParseError, "", structs.Map(res), true)
						return
					}
					// convert base64 image to blob
					blobImg, err1 := base64.StdEncoding.DecodeString(eachThumbImgData.ValidatedBase64)
					if err1 != nil {
						res = io.FailureMessage(errors.New("Error decoding image"), "Error decoding image")
						res.Error = errors.New("Error decoding image")
						b.Logger.Log(gologger.Errsev3,
							gologger.InternalServices,
							"ValidateVenueImages",
							"Validate venue images request failed",
							gologger.ParseError, "", structs.Map(res), true)
						return
					}
					eachThumbImgData.Blob = string(blobImg)
					eachThumbImgData.Exists = false
					fileExt := ""
					splittedContentType := strings.Split(eachThumbImgData.ContentType, "/")
					fileExt = splittedContentType[1]
					eachThumbImgData.FileName = "thumb_img_" + venueIdStr + "." + fileExt
					thumbImgData = append(thumbImgData, eachThumbImgData)
				}
			} else {
				eachThumbImgData.Verified = true
				eachThumbImgData.ImageUrl = eachThumbImg.Image
				eachThumbImgData.Exists = true
				eachThumbImgData.Blob = ""
				thumbImgData = append(thumbImgData, eachThumbImgData)
			}
		}
	}

	if len(u.GalleryImg) > 0 {
		counter := 0
		counterStr := ""
		for _, eachGalImg := range u.GalleryImg {
			var eachGalImgData io.VenueImageData
			if eachGalImg.Status == true {
				// check if image url already exists in db. If yes, then dont do anything.
				// If url does not present, then consider it as base64
				// get venue images for venue
				exists, err := b.VenueRepo.CheckVenueImgExists(ctx, dal, eachGalImg.Image)
				if err != nil {
					res = io.FailureMessage(res.Error, "Error getting venue images details")
					res.Error = err
					b.Logger.Log(gologger.Errsev3,
						gologger.InternalServices,
						"Validate venue images",
						"Validate venue images request failed",
						gologger.ParseError, "", structs.Map(res), true)
					return
				}
				if exists == false {
					counter = counter + 1
					counterStr = strconv.Itoa(counter)
					eachGalImgData.ValidatedBase64, eachGalImgData.ContentType, eachGalImgData.Verified = utile.ValidateBase64(eachGalImg.Image, "venue")
					if eachGalImgData.Verified == false {
						res = io.FailureMessage(errors.New("Gallery Image Format is invalid"), "Gallery Image Format is invalid")
						res.Error = errors.New("Gallery Image Format is invalid")
						b.Logger.Log(gologger.Errsev3,
							gologger.InternalServices,
							"ValidateMerchantDocs",
							"Validate merchant documents request failed",
							gologger.ParseError, "", structs.Map(res), true)
						return
					}
					// convert base64 image to blob
					blobImg, err1 := base64.StdEncoding.DecodeString(eachGalImgData.ValidatedBase64)
					if err1 != nil {
						res = io.FailureMessage(errors.New("Error decoding image"), "Error decoding image")
						res.Error = errors.New("Error decoding image")
						b.Logger.Log(gologger.Errsev3,
							gologger.InternalServices,
							"ValidateMerchantDocs",
							"Validate merchant documents request failed",
							gologger.ParseError, "", structs.Map(res), true)
						return
					}
					eachGalImgData.Blob = string(blobImg)
					eachGalImgData.Exists = false
					fileExt := ""
					splittedContentType := strings.Split(eachGalImgData.ContentType, "/")
					fileExt = splittedContentType[1]
					eachGalImgData.FileName = "gal_img_" + venueIdStr + "_" + totalimgStr + counterStr + "." + fileExt
					galleryImgData = append(galleryImgData, eachGalImgData)
				}
			} else {
				eachGalImgData.Verified = true
				eachGalImgData.ImageUrl = eachGalImg.Image
				eachGalImgData.Exists = true
				eachGalImgData.Blob = ""
				galleryImgData = append(galleryImgData, eachGalImgData)
			}
		}
	}

	data := make(map[string]interface{})
	data["headerImgData"] = headerImgData
	data["thumbImgData"] = thumbImgData
	data["galleryImgData"] = galleryImgData

	res = io.SuccessMessage(data, "Uploaded validated")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"Validate venue images",
		"Validate venue images response body",
		"Validate venue imagesResponse", "", structs.Map(res), true)
	return
}

func (b *athUploadService) UploadToS3(ctx context.Context, u map[string]interface{}) (res io.Response) {

	data := make(map[string]interface{})
	data["GstNoFile"] = ""
	data["BankAccFile"] = ""
	data["AddressFile"] = ""
	data["PanNoFile"] = ""

	gstFileData := u["gstFileData"].(io.MerchantDocFileData)
	if gstFileData.Blob != "" {
		uploaderr := utile.UploadToS3([]byte(gstFileData.Blob), gstFileData.ContentType, gstFileData.FileName)
		if uploaderr != nil {
			res = io.FailureMessage(res.Error, "Error Uploading file to google cloud")
			res.Error = uploaderr
			b.Logger.Log(gologger.Errsev3,
				gologger.InternalServices,
				"UploadToS3",
				"Upload documents service request failed",
				gologger.ParseError, "", structs.Map(res), true)
			return
		}
		data["GstNoFile"] = "https://storage.cloud.google.com/sit-turf/account/" + gstFileData.FileName
	}

	bankFileData := u["bankFileData"].(io.MerchantDocFileData)
	if bankFileData.Blob != "" {
		uploaderr := utile.UploadToS3([]byte(bankFileData.Blob), bankFileData.ContentType, bankFileData.FileName)
		if uploaderr != nil {
			res = io.FailureMessage(res.Error, "Error Uploading file to google cloud")
			res.Error = uploaderr
			b.Logger.Log(gologger.Errsev3,
				gologger.InternalServices,
				"UploadToS3",
				"Upload documents service request failed",
				gologger.ParseError, "", structs.Map(res), true)
			return
		}
		data["BankAccFile"] = "https://storage.cloud.google.com/sit-turf/account/" + bankFileData.FileName
	}

	addFileData := u["addFileData"].(io.MerchantDocFileData)
	if addFileData.Blob != "" {
		uploaderr := utile.UploadToS3([]byte(addFileData.Blob), addFileData.ContentType, addFileData.FileName)
		if uploaderr != nil {
			res = io.FailureMessage(res.Error, "Error Uploading file to google cloud")
			res.Error = uploaderr
			b.Logger.Log(gologger.Errsev3,
				gologger.InternalServices,
				"UploadToS3",
				"Upload documents service request failed",
				gologger.ParseError, "", structs.Map(res), true)
			return
		}
		data["AddressFile"] = "https://storage.cloud.google.com/sit-turf/account/" + addFileData.FileName
	}

	panFileData := u["panFileData"].(io.MerchantDocFileData)
	if panFileData.Blob != "" {
		uploaderr := utile.UploadToS3([]byte(panFileData.Blob), panFileData.ContentType, panFileData.FileName)
		if uploaderr != nil {
			res = io.FailureMessage(res.Error, "Error Uploading file to google cloud")
			res.Error = uploaderr
			b.Logger.Log(gologger.Errsev3,
				gologger.InternalServices,
				"UploadToS3",
				"Upload documents service request failed",
				gologger.ParseError, "", structs.Map(res), true)
			return
		}
		data["PanNoFile"] = "https://storage.cloud.google.com/sit-turf/account/" + panFileData.FileName
	}

	res.Data = data
	res = io.SuccessMessage(data, "Uploaded documents")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"UploadToS3",
		"Upload documents response body",
		"UploadToS3Response", "", structs.Map(res), true)
	return
}

func (b *athUploadService) VenueUploadToS3(ctx context.Context, u map[string]interface{}) (res io.Response) {
	var data io.VenueImagesReq
	headerFileData := u["headerImgData"].([]io.VenueImageData)
	for _, eachHeaderFileData := range headerFileData {
		var eachHeaderImg io.ImgData
		if eachHeaderFileData.Exists == false {
			//upload image to google cloud
			uploaderr := utile.VenueUploadToS3([]byte(eachHeaderFileData.Blob), eachHeaderFileData.ContentType, eachHeaderFileData.FileName)
			if uploaderr != nil {
				res = io.FailureMessage(res.Error, "Error Uploading header file to google cloud")
				res.Error = uploaderr
				b.Logger.Log(gologger.Errsev3,
					gologger.InternalServices,
					"UploadToS3",
					"Upload documents service request failed",
					gologger.ParseError, "", structs.Map(res), true)
				return
			}
			eachHeaderImg.Image = "https://storage.cloud.google.com/sit-turf-images/venue/" + eachHeaderFileData.FileName
			eachHeaderImg.Status = true
		} else {
			// delete image from db
			eachHeaderImg.Image = eachHeaderFileData.ImageUrl
			eachHeaderImg.Status = false
		}
		data.HeaderImg = append(data.HeaderImg, eachHeaderImg)
	}

	thumbFileData := u["thumbImgData"].([]io.VenueImageData)
	for _, eachThumbFileData := range thumbFileData {
		var eachThumbImg io.ImgData
		if eachThumbFileData.Exists == false {
			//upload image to google cloud
			uploaderr := utile.VenueUploadToS3([]byte(eachThumbFileData.Blob), eachThumbFileData.ContentType, eachThumbFileData.FileName)
			if uploaderr != nil {
				res = io.FailureMessage(res.Error, "Error Uploading thumbnail file to google cloud")
				res.Error = uploaderr
				b.Logger.Log(gologger.Errsev3,
					gologger.InternalServices,
					"VenueUploadToS3",
					"venue Upload documents service request failed",
					gologger.ParseError, "", structs.Map(res), true)
				return
			}
			eachThumbImg.Image = "https://storage.cloud.google.com/sit-turf-images/venue/" + eachThumbFileData.FileName
			eachThumbImg.Status = true
		} else {
			// delete image from db
			eachThumbImg.Image = eachThumbFileData.ImageUrl
			eachThumbImg.Status = false
		}
		data.ThumbnailImg = append(data.ThumbnailImg, eachThumbImg)
	}

	galleryFileData := u["galleryImgData"].([]io.VenueImageData)
	for _, eachGalFileData := range galleryFileData {
		var eachGalImg io.ImgData
		if eachGalFileData.Exists == false {
			//upload image to google cloud
			uploaderr := utile.VenueUploadToS3([]byte(eachGalFileData.Blob), eachGalFileData.ContentType, eachGalFileData.FileName)
			if uploaderr != nil {
				res = io.FailureMessage(res.Error, "Error Uploading Gallery file to google cloud")
				res.Error = uploaderr
				b.Logger.Log(gologger.Errsev3,
					gologger.InternalServices,
					"UploadToS3",
					"Upload documents service request failed",
					gologger.ParseError, "", structs.Map(res), true)
				return
			}
			eachGalImg.Image = "https://storage.cloud.google.com/sit-turf-images/venue/" + eachGalFileData.FileName
			eachGalImg.Status = true
		} else {
			// delete image from db
			eachGalImg.Image = eachGalFileData.ImageUrl
			eachGalImg.Status = false
		}
		data.GalleryImg = append(data.GalleryImg, eachGalImg)
	}

	res.Data = data
	res = io.SuccessMessage(data, "Uploaded documents")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"VenueUploadToS3",
		"venue Upload documents response body",
		"VenueUploadToS3Response", "", structs.Map(res), true)
	return
}
