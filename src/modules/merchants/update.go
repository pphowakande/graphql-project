package handler

import (
	io "api/src/models"
	utile "api/src/utils"
	"context"
	"strings"
	"time"

	"github.com/fatih/structs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"stash.bms.bz/bms/gologger.git"

	merchantpb "stash.bms.bz/turf/generic-proto-files.git/merchant/v1"
)

func (h *athMerchantHandler) EditMerchant(ctx context.Context, req *merchantpb.EditMerchantRequest) (*merchantpb.GetMerchantByIDRes, error) {
	loggedInUserID, err := utile.GetUserIDFromContext(ctx)
	if err != nil {
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"EditMerchant",
			"edit merchant request failed",
			gologger.ParseError, "", structs.Map(*req), true)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	storeErr := ""
	// check if required fields are not empty
	if (req.MerchantFullName != "") && (req.BusinessName != "") && (req.Email != "") && (req.Phone != "") {
		h.logger.Log(gologger.Info,
			gologger.ExternalServices,
			"EditMerchant",
			"edit merchant request parameters",
			"EditMerchantRequest", "", structs.Map(*req), true)

		FirstName := ""
		LastName := ""

		if strings.Contains(req.MerchantFullName, " ") {
			splittedName := strings.Split(req.MerchantFullName, " ")
			FirstName = splittedName[0]
			LastName = splittedName[1]
			if len(splittedName) > 2 {
				LastName = LastName + " " + splittedName[2]
			}
		} else {
			FirstName = req.MerchantFullName
			LastName = ""
		}

		userRequest := io.AthUser{
			UserID:    loggedInUserID,
			Email:     req.Email,
			FirstName: FirstName,
			Phone:     req.Phone,
			LastName:  LastName,
			Models: io.Models{
				UpdatedBy: loggedInUserID,
				UpdatedAt: int(time.Now().Unix()),
			},
		}

		DoLater := true
		PanNoFile := ""
		BankAccFile := ""
		AddressFile := ""
		GstNoFile := ""

		masterTx, transactionTx, err := h.GetMasterAndTransactionDBTransaction(ctx)
		if err != nil {
			return nil, err
		}
		defer func() {
			masterTx.RollbackUnlessCommitted()
			transactionTx.RollbackUnlessCommitted()
		}()

		// update user details
		userServiceRes := h.athUserService.EditUser(ctx, masterTx, userRequest, "owner")
		if userServiceRes.Error != nil {
			storeErr = userServiceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"EditMerchant",
				"edit merchant request failed",
				gologger.ParseError, "", structs.Map(userServiceRes), true)
			return nil, status.Errorf(codes.Internal, storeErr)
		}

		// check if KYC documents are uploaded. If yes, update them on google cloud
		// check if KYC documents are already verified by bms,if yes dont allow merchant to change them
		if req.KycData != nil {
			MerchantKYC := io.MerchantKYC{
				PanNoFile:   req.KycData.PanNoFile,
				AddressFile: req.KycData.AddressFile,
				GstNoFile:   req.KycData.GstNoFile,
				DoLater:     req.KycData.DoLater,
				BankAccFile: req.KycData.BankAccFile,
				MerchantID:  loggedInUserID,
			}

			// first validate uploaded documents and check if document is already uploaded
			validateServiceRes := h.athUploadService.ValidateMerchantDocs(ctx, MerchantKYC)
			if validateServiceRes.Error != nil {
				storeErr = validateServiceRes.Error.Error()
				h.logger.Log(gologger.Errsev3,
					gologger.ExternalServices,
					"EditMerchant",
					"Edit merchant request failed",
					gologger.ParseError, "", structs.Map(validateServiceRes), true)
				return nil, status.Errorf(codes.Internal, "error in EditMerchant :%v", storeErr)
			}

			validatedData := validateServiceRes.Data.(map[string]interface{})
			// upload documents to google cloud
			// if documents are validated , upload it to google cloud
			uploadServiceRes := h.athUploadService.UploadToS3(ctx, validatedData)
			if uploadServiceRes.Error != nil {
				storeErr = uploadServiceRes.Error.Error()
				h.logger.Log(gologger.Errsev3,
					gologger.ExternalServices,
					"EditMerchant",
					"edit merchant request failed",
					gologger.ParseError, "", structs.Map(uploadServiceRes), true)
				return nil, status.Errorf(codes.Internal, "error in EditMerchant :%v", storeErr)
			}

			uploadData := uploadServiceRes.Data.(map[string]interface{})

			if uploadData["PanNoFile"].(string) != "" && uploadData["AddressFile"].(string) != "" && uploadData["BankAccFile"].(string) != "" && uploadData["GstNoFile"].(string) != "" {
				DoLater = false
			}
			DoLater = true
			// get google cloud links back and save it in database
			PanNoFile = uploadData["PanNoFile"].(string)
			AddressFile = uploadData["AddressFile"].(string)
			BankAccFile = uploadData["BankAccFile"].(string)
			GstNoFile = uploadData["GstNoFile"].(string)
		}

		// edit merchant data
		merchantRequest := io.AthMerchant{
			MerchantName: req.BusinessName,
			Phone:        req.Phone,
			Email:        req.Email,
			Models: io.Models{
				CreatedBy: loggedInUserID,
				UpdatedBy: loggedInUserID,
			},
			PanNoFile:   PanNoFile,
			AddressFile: AddressFile,
			BankAccFile: BankAccFile,
			GstNoFile:   GstNoFile,
			DoLater:     DoLater,
		}

		// update merchant data
		merchantServiceRes := h.athMerchantService.EditMerchant(ctx, masterTx, merchantRequest)
		if merchantServiceRes.Error != nil {
			storeErr = merchantServiceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"EditMerchant",
				"edit merchant request failed",
				gologger.ParseError, "", structs.Map(merchantServiceRes), true)
			return nil, status.Errorf(codes.Internal, storeErr)
		}

		merchantData := merchantServiceRes.Data.(io.AthMerchant)
		MerchantID := merchantData.MerchantID

		// get merchnat latest data from database
		merchantDataServiceRes := h.athMerchantService.GetMerchantByID(ctx, masterTx, MerchantID)
		if merchantDataServiceRes.Error != nil {
			storeErr = merchantDataServiceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"EditMerchant",
				"edit merchant request failed",
				gologger.ParseError, "", structs.Map(merchantDataServiceRes), true)
			return nil, status.Errorf(codes.Internal, storeErr)
		}

		merchantUpdatedData := merchantDataServiceRes.Data.(io.AthMerchant)

		KYCData := &merchantpb.KycData{
			BankAccFile: merchantUpdatedData.BankAccFile,
			GstNoFile:   merchantUpdatedData.GstNoFile,
			AddressFile: merchantUpdatedData.AddressFile,
			PanNoFile:   merchantUpdatedData.PanNoFile,
		}

		res := &merchantpb.GetMerchantByIDRes{
			AccountId:        int32(MerchantID),
			MerchantFullName: FirstName + " " + LastName,
			BusinessName:     merchantUpdatedData.MerchantName,
			Address:          merchantUpdatedData.Address,
			Phone:            merchantUpdatedData.Phone,
			Email:            merchantUpdatedData.Email,
			KycData:          KYCData,
		}

		userData := userServiceRes.Data.(map[string]interface{})

		// if email is updated, send email verification
		emailVerify := utile.RandomString(6)
		if userData["EmailVerify"].(bool) == false {
			// email verification link gets sent
			emailSendReq := io.EmailSendReq{
				EmailVerifyToken: emailVerify,
				Email:            userData["Email"].(string),
				Account:          "WEBIN",
				Subject:          "Please Verify your email",
				Body:             "Email verification OTP is - " + emailVerify,
				Eticket:          "N",
				Tid:              "0",
				OTPTtype:         "email",
			}

			otpServiceRes := h.athOTPService.EmailSend(ctx, masterTx, transactionTx, emailSendReq)
			if otpServiceRes.Error != nil {
				storeErr = otpServiceRes.Error.Error()
				h.logger.Log(gologger.Errsev3,
					gologger.ExternalServices,
					"EditMerchant",
					"edit merchant request failed",
					gologger.ParseError, "", structs.Map(otpServiceRes), true)
				return nil, status.Errorf(codes.Internal, "error in EditMerchant :%v", storeErr)
			}
		}

		// if phone is updated, send phone veification
		if userData["PhoneVerify"].(bool) == false {
			// send an OTP
			var otpRequest io.AthUserOTP
			otpRequest.UserID = userData["UserID"].(int)
			otpRequest.OTPNO = utile.RandomString(6)
			otpRequest.OTPType = "phone"
			otpRequest.ExpiredAt = int(time.Now().Unix()) + 600 // otp expires in 10 minutes
			otpRequest.CreatedAt = int(time.Now().Unix())

			otpServiceRes := h.athOTPService.CreateOTP(ctx, masterTx, transactionTx, otpRequest)
			if otpServiceRes.Error != nil {
				storeErr = otpServiceRes.Error.Error()
				h.logger.Log(gologger.Errsev3,
					gologger.ExternalServices,
					"EditMerchant",
					"edit merchant request failed",
					gologger.ParseError, "", structs.Map(otpServiceRes), true)
				return nil, status.Errorf(codes.Internal, "error in EditMerchant :%v", storeErr)
			}
		}
		h.logger.Log(gologger.Info,
			gologger.ExternalServices,
			"EditMerchant",
			"edit merchant response body",
			"EditMerchantResponse", "", structs.Map(*res), true)
		if err = transactionTx.Commit().Error; err != nil {
			return nil, status.Errorf(codes.Internal, "error in EditMerchant :%v", err.Error())
		}
		return res, masterTx.Commit().Error
	}

	storeErr = "Request parameter is missing or blank"
	h.logger.Log(gologger.Errsev3,
		gologger.ExternalServices,
		"EditMerchant",
		"edit merchant request failed",
		gologger.ParseError, "", structs.Map(req), true)
	return nil, status.Errorf(codes.Internal, "error in EditMerchant :%v", storeErr)

}

func (h *athMerchantHandler) EditTeamMember(ctx context.Context, req *merchantpb.EditTeamMemberRequest) (*merchantpb.EditTeamMemberReply, error) {
	storeErr := ""
	if req.AccountId != 0 {
		loggedInUserID, err := utile.GetUserIDFromContext(ctx)
		if err != nil {
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"EditTeamMember",
				"Edit team member request failed",
				gologger.ParseError, "", structs.Map(*req), true)
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		h.logger.Log(gologger.Info,
			gologger.ExternalServices,
			"EditTeamMember",
			"Edit team member request parameters",
			"EditTeamMemberRequest", "", structs.Map(*req), true)

		masterTx, transactionTx, err := h.GetMasterAndTransactionDBTransaction(ctx)
		if err != nil {
			return nil, err
		}
		defer func() {
			masterTx.RollbackUnlessCommitted()
			transactionTx.RollbackUnlessCommitted()
		}()

		// update team member privileges to table
		for _, privData := range req.Privileges {
			updatePrivReq := io.AthVenueUser{
				VenueId:  int(privData.VenueId),
				UserId:   int(req.AccountId),
				IsActive: privData.Status,
				Models: io.Models{
					UpdatedBy: int(loggedInUserID),
				},
			}

			teamMemberServiceRes := h.athMerchantService.UpdateTeamMemberPrivileges(ctx, masterTx, updatePrivReq)
			if teamMemberServiceRes.Error != nil {
				storeErr = teamMemberServiceRes.Error.Error()
				h.logger.Log(gologger.Errsev3,
					gologger.ExternalServices,
					"EditTeamMember",
					"Edit team member request failed",
					gologger.ParseError, "", structs.Map(teamMemberServiceRes), true)
				return nil, status.Errorf(codes.Internal, "error in EditTeamMember :%v", storeErr)
			}
		}

		// update basic user details
		userRequest := io.AthUser{
			UserID:    int(req.AccountId),
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Phone:     req.Phone,
			Models: io.Models{
				UpdatedBy: loggedInUserID,
			},
		}

		// add users details to athusers table
		//if email an phone are already verified, do not allow them to change
		userServiceRes := h.athUserService.EditUser(ctx, masterTx, userRequest, "teammember")
		if userServiceRes.Error != nil {
			storeErr = userServiceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"EditTeamMember",
				"Edit team member request failed",
				gologger.ParseError, "", structs.Map(userServiceRes), true)
			return nil, status.Errorf(codes.Internal, "error in EditTeamMember :%v", storeErr)
		}

		userData := userServiceRes.Data.(map[string]interface{})

		// if email is updated, send email verification
		if userData["EmailVerify"].(bool) == false {
			// email verification link gets sent
			emailVerify := utile.RandomString(6)
			emailSendReq := io.EmailSendReq{
				EmailVerifyToken: emailVerify,
				Email:            userData["Email"].(string),
				Account:          "WEBIN",
				Subject:          "Please Verify your email",
				Body:             "email verification OTP - " + emailVerify,
				Eticket:          "N",
				Tid:              "0",
				OTPTtype:         "email",
			}

			otpServiceRes := h.athOTPService.EmailSend(ctx, masterTx, transactionTx, emailSendReq)
			if otpServiceRes.Error != nil {
				storeErr = otpServiceRes.Error.Error()
				h.logger.Log(gologger.Errsev3,
					gologger.ExternalServices,
					"EditTeamMember",
					"edit team member request failed",
					gologger.ParseError, "", structs.Map(otpServiceRes), true)
				return nil, status.Errorf(codes.Internal, "error in EditTeamMember :%v", storeErr)
			}
		}

		// if phone is updated, send phone veification
		if userData["PhoneVerify"].(bool) == false {
			// send an OTP
			var otpRequest io.AthUserOTP
			otpRequest.UserID = userData["UserID"].(int)
			otpRequest.OTPNO = utile.RandomString(6)
			otpRequest.OTPType = "phone"
			otpRequest.ExpiredAt = int(time.Now().Unix()) + 600 // otp expires in 10 minutes
			otpRequest.CreatedAt = int(time.Now().Unix())

			otpServiceRes := h.athOTPService.CreateOTP(ctx, masterTx, transactionTx, otpRequest)
			if otpServiceRes.Error != nil {
				storeErr = otpServiceRes.Error.Error()
				h.logger.Log(gologger.Errsev3,
					gologger.ExternalServices,
					"EditTeamMember",
					"Edit team member request failed",
					gologger.ParseError, "", structs.Map(otpServiceRes), true)
				return nil, status.Errorf(codes.Internal, "error in EditTeamMember :%v", storeErr)
			}
		}

		// get updated user details
		res := &merchantpb.EditTeamMemberReply{
			FirstName:   userData["FirstName"].(string),
			LastName:    userData["LastName"].(string),
			Email:       userData["Email"].(string),
			Phone:       userData["Phone"].(string),
			AccountType: userData["AccountType"].(string),
			Privileges: &merchantpb.AccessData{
				VenueIds: userData["AccessData"].([]int32),
			},
			AccountId: int32(userData["AccountId"].(int)),
		}
		if err = transactionTx.Commit().Error; err != nil {
			return nil, status.Errorf(codes.Internal, "error in EditTeamMember :%v", err.Error())
		}
		return res, masterTx.Commit().Error
	}
	storeErr = "Request parameter is missing or blank"
	h.logger.Log(gologger.Errsev3,
		gologger.ExternalServices,
		"EditTeamMember",
		"Edit team member request failed",
		gologger.ParseError, "", structs.Map(req), true)
	return nil, status.Errorf(codes.Internal, "error in EditTeamMember :%v", storeErr)
}
