package handler

import (
	io "api/src/models"
	utile "api/src/utils"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fatih/structs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"stash.bms.bz/bms/gologger.git"

	merchantpb "stash.bms.bz/turf/generic-proto-files.git/merchant/v1"
)

func (h *athMerchantHandler) SignupMerchant(ctx context.Context, req *merchantpb.SignupRequest) (*merchantpb.SignupReply, error) {
	h.logger.Log(gologger.Info,
		gologger.ExternalServices,
		"SignupMerchant",
		"signup merchant request parameters",
		"SignupRequest", "", structs.Map(*req), true)

	storeErr := ""
	// check if required fields are present in request
	if (req.BusinessName != "") &&
		(req.Phone != "") &&
		(req.Email != "") &&
		(req.Password != "") &&
		(req.MerchantFullName != "") &&
		(req.UserSource != "") {
		// process request here
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
			FirstName:  FirstName,
			LastName:   LastName,
			Email:      req.Email,
			Password:   req.Password,
			UserSource: req.UserSource,
			Phone:      req.Phone,
		}

		masterTx, transactionTx, err := h.GetMasterAndTransactionDBTransaction(ctx)
		if err != nil {
			return nil, err
		}
		defer func() {
			masterTx.RollbackUnlessCommitted()
			transactionTx.RollbackUnlessCommitted()
		}()

		// check if merchant email already exists
		// save user details and get user id
		userServiceRes := h.athUserService.CreateUser(ctx, masterTx, userRequest, "merchant")
		if userServiceRes.Error != nil {
			storeErr = userServiceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"SignupMerchant",
				"signup merchant request failed",
				gologger.ParseError, "", structs.Map(userServiceRes), true)
			return nil, status.Errorf(codes.Internal, "error in SignupMerchant :%v", storeErr)
		}
		userData := userServiceRes.Data.(io.AthUser)
		userID := userData.UserID
		log.Printf("I am here 1:%v\n", userID)
		merchantRequest := io.AthMerchant{
			MerchantName: req.BusinessName,
			Phone:        req.Phone,
			Email:        req.Email,
			Models: io.Models{
				CreatedBy: userID,
				CreatedAt: int(time.Now().Unix()),
			},
		}

		// get merchant id from merchant service
		merchantServiceRes := h.athMerchantService.CreateMerchant(ctx, masterTx, merchantRequest)
		if merchantServiceRes.Error != nil {
			storeErr = merchantServiceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"SignupMerchant",
				"signup merchant request failed",
				gologger.ParseError, "", structs.Map(merchantServiceRes), true)
			return nil, status.Errorf(codes.Internal, "error in SignupMerchant :%v", storeErr)
		}
		merchantData := merchantServiceRes.Data.(map[string]interface{})
		if merchantData["Merchant_id"] == nil {
			return nil, status.Errorf(codes.Internal, "error in SignupMerchant :%v", storeErr)
		}
		merchantID := merchantData["Merchant_id"].(int)

		// add details to merchant_user table
		merchantUserReq := io.AthMerchantUser{
			MerchantID: merchantID,
			UserID:     userID,
			CreatedBy:  userID,
			CreatedAt:  int(time.Now().Unix()),
		}

		merchantUserServiceRes := h.athMerchantService.CreateMerchantUser(ctx, masterTx, merchantUserReq)
		if merchantUserServiceRes.Error != nil {
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"SignupMerchant",
				"signup merchant request failed",
				gologger.ParseError, "", structs.Map(merchantUserServiceRes), true)
			return nil, status.Errorf(codes.Internal, "error in SignupMerchant :%v", storeErr)
		}

		//send an OTP
		var otpRequest io.AthUserOTP
		otpRequest.UserID = int(userID)
		otpRequest.OTPNO = utile.RandomString(6)
		otpRequest.OTPType = "phone"
		otpRequest.ExpiredAt = int(time.Now().Unix()) + 600 // otp expires in 10 minutes
		otpRequest.CreatedAt = int(time.Now().Unix())

		fmt.Println("Calling createotp service-----------")
		otpServiceRes := h.athOTPService.CreateOTP(ctx, masterTx, transactionTx, otpRequest)
		fmt.Println("otpServiceRes : ", otpServiceRes)
		fmt.Println("otpServiceRes.error : ", otpServiceRes.Error)
		if otpServiceRes.Error != nil {
			storeErr = otpServiceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"SignupMerchant",
				"signup merchant request failed",
				gologger.ParseError, "", structs.Map(otpServiceRes), true)
			return nil, status.Errorf(codes.Internal, "error in SignupMerchant :%v", storeErr)
		}

		res := &merchantpb.SignupReply{
			AccountId: int32(userID),
		}

		h.logger.Log(gologger.Info,
			gologger.ExternalServices,
			"SignupMerchant",
			"signup merchant response body",
			"signupUserResponse", "", structs.Map(*res), true)

		if err = transactionTx.Commit().Error; err != nil {
			return nil, status.Errorf(codes.Internal, "error in SignupMerchant :%v", err.Error())
		}
		return res, masterTx.Commit().Error

	}

	storeErr = "Request parameter is missing or blank"
	h.logger.Log(gologger.Errsev3,
		gologger.ExternalServices,
		"SignupMerchant",
		"signup merchant request failed",
		gologger.ParseError, "", structs.Map(req), true)
	return nil, status.Errorf(codes.Internal, "error in SignupMerchant :%v", storeErr)

}

func (h *athMerchantHandler) UploadDoc(ctx context.Context, req *merchantpb.UploadDocRequest) (*merchantpb.GenericReply, error) {
	h.logger.Log(gologger.Info,
		gologger.ExternalServices,
		"UploadDoc",
		"Upload documents request parameters",
		"UploadDocRequest", "", structs.Map(*req), true)
	storeErr := ""

	// check if merchant id is supplied
	if req.AccountId != 0 && req.KycData != nil {

		MerchantKYC := io.MerchantKYC{
			PanNoFile:   req.KycData.PanNoFile,
			AddressFile: req.KycData.AddressFile,
			GstNoFile:   req.KycData.GstNoFile,
			DoLater:     req.KycData.DoLater,
			BankAccFile: req.KycData.BankAccFile,
			MerchantID:  int(req.AccountId),
		}

		// first validate uploaded documents and check if document is already uploaded
		validateServiceRes := h.athUploadService.ValidateMerchantDocs(ctx, MerchantKYC)
		if validateServiceRes.Error != nil {
			storeErr = validateServiceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"UploadDoc",
				"Upload documents request failed",
				gologger.ParseError, "", structs.Map(validateServiceRes), true)
			return nil, status.Errorf(codes.Internal, "error in UploadDoc :%v", storeErr)
		}

		validatedData := validateServiceRes.Data.(map[string]interface{})
		// upload documents to google cloud
		// if documents are validated , upload it to google cloud
		uploadServiceRes := h.athUploadService.UploadToS3(ctx, validatedData)
		if uploadServiceRes.Error != nil {
			storeErr = uploadServiceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"UploadDoc",
				"Upload documents request failed",
				gologger.ParseError, "", structs.Map(uploadServiceRes), true)
			return nil, status.Errorf(codes.Internal, "error in UploadDoc :%v", storeErr)
		}

		uploadData := uploadServiceRes.Data.(map[string]interface{})

		DoLater := false
		if req.KycData.PanNoFile != "" && req.KycData.AddressFile != "" && req.KycData.GstNoFile != "" && req.KycData.BankAccFile != "" {
			DoLater = false
		}
		DoLater = true

		// get google cloud links back and save it in database
		merchantReq := io.AthMerchant{
			Models: io.Models{
				CreatedBy: int(req.AccountId),
				UpdatedBy: int(req.AccountId),
			},
			DoLater:     DoLater,
			PanNoFile:   uploadData["PanNoFile"].(string),
			AddressFile: uploadData["AddressFile"].(string),
			BankAccFile: uploadData["BankAccFile"].(string),
			GstNoFile:   uploadData["GstNoFile"].(string),
		}

		masterTx, transactionTx, err := h.GetMasterAndTransactionDBTransaction(ctx)
		if err != nil {
			return nil, err
		}
		defer func() {
			masterTx.RollbackUnlessCommitted()
			transactionTx.RollbackUnlessCommitted()
		}()

		// write logic here
		editServiceRes := h.athMerchantService.EditMerchant(ctx, masterTx, merchantReq)
		if editServiceRes.Error != nil {
			storeErr = editServiceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"UploadDoc",
				"Upload documents request failed",
				gologger.ParseError, "", structs.Map(editServiceRes), true)
			return nil, status.Errorf(codes.Internal, "error in UploadDoc :%v", storeErr)
		}

		// check if user email is verified or not
		userData := editServiceRes.Data.(io.AthMerchant)
		verfifyOTP := utile.RandomString(6)
		// email verification link gets sent
		emailSendReq := io.EmailSendReq{
			EmailVerifyToken: verfifyOTP,
			Email:            userData.Email,
			Account:          "WEBIN",
			Subject:          "Please Verify your email",
			Body:             "Please click on link to verify your email - http://odp-turf-web-01.bigtree.org/sign-in?code=" + verfifyOTP + "&email=" + userData.Email,
			Eticket:          "N",
			Tid:              "0",
			OTPTtype:         "email",
		}

		otpServiceRes := h.athOTPService.EmailSend(ctx, masterTx, transactionTx, emailSendReq)
		if otpServiceRes.Error != nil {
			storeErr = otpServiceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"PhoneVerifyMerchant",
				"phone verify merchant request failed",
				gologger.ParseError, "", structs.Map(otpServiceRes), true)
			return nil, status.Errorf(codes.Internal, "error in PhoneVerifyMerchant :%v", storeErr)
		}

		res := &merchantpb.GenericReply{
			Status: otpServiceRes.Success,
		}
		h.logger.Log(gologger.Info,
			gologger.ExternalServices,
			"UploadDoc",
			"Upload documents response body",
			"UploadDocRequest", "", structs.Map(*res), true)
		if err = transactionTx.Commit().Error; err != nil {
			return nil, status.Errorf(codes.Internal, "error in UploadDoc :%v", err.Error())
		}
		return res, masterTx.Commit().Error
	}
	storeErr = "Request parameter is missing or blank"
	h.logger.Log(gologger.Errsev3,
		gologger.ExternalServices,
		"UploadDoc",
		"Upload documents request failed",
		gologger.ParseError, "", structs.Map(req), true)
	return nil, status.Errorf(codes.Internal, "error in UploadDoc :%v", storeErr)
}

func (h *athMerchantHandler) ResendCode(ctx context.Context, req *merchantpb.ResendCodeRequest) (*merchantpb.GenericReply, error) {
	fmt.Println("ResendCode module function---------")
	h.logger.Log(gologger.Info,
		gologger.ExternalServices,
		"resendCode",
		"resend code request parameters",
		"resendCodeRequest", "", structs.Map(*req), true)
	storeErr := ""

	// check if fields are not empty
	if req.AccountId != 0 {
		masterTx, transactionTx, err := h.GetMasterAndTransactionDBTransaction(ctx)
		if err != nil {
			return nil, err
		}
		defer func() {
			masterTx.RollbackUnlessCommitted()
			transactionTx.RollbackUnlessCommitted()
		}()
		// check if merchant id passed exists. If yes, get phone from database
		verifyServiceRes := h.athMerchantService.VerifyGetMerchant(ctx, masterTx, req.AccountId)
		if verifyServiceRes.Error != nil {
			storeErr = verifyServiceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"resendCode",
				"resend code request failed",
				gologger.ParseError, "", structs.Map(verifyServiceRes), true)
			return nil, status.Errorf(codes.Internal, "error in user resendCode :%v", storeErr)
		}

		MerchantData := verifyServiceRes.Data.(io.AthMerchant)

		// trigger create token service
		var otpRequest io.AthUserOTP
		otpRequest.UserID = MerchantData.CreatedBy
		otpRequest.OTPNO = utile.RandomString(6)
		otpRequest.OTPType = "phone"
		otpRequest.Phone = MerchantData.Phone
		otpRequest.ExpiredAt = int(time.Now().Unix()) + 600 // otp expires in 10 minutes
		otpRequest.CreatedAt = int(time.Now().Unix())

		otpServiceRes := h.athOTPService.CreateOTP(ctx, masterTx, transactionTx, otpRequest)
		if otpServiceRes.Error != nil {
			storeErr = otpServiceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"resendCode",
				"resend code request failed",
				gologger.ParseError, "", structs.Map(otpServiceRes), true)
			return nil, status.Errorf(codes.Internal, storeErr)
		}
		res := &merchantpb.GenericReply{
			Status: otpServiceRes.Success,
		}
		h.logger.Log(gologger.Info,
			gologger.ExternalServices,
			"resendCode",
			"resend code response body",
			"resendCodeResponse", "", structs.Map(*res), true)
		if err = transactionTx.Commit().Error; err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		return res, masterTx.Commit().Error
	}
	storeErr = "Request parameter is missing or blank"
	h.logger.Log(gologger.Errsev3,
		gologger.ExternalServices,
		"resendCode",
		"resend code request failed",
		gologger.ParseError, "", structs.Map(req), true)
	return nil, status.Errorf(codes.Internal, "error in resendCode :%v", storeErr)

}

func (h *athMerchantHandler) PhoneVerifyMerchant(ctx context.Context, req *merchantpb.PhoneVerifyRequest) (*merchantpb.GenericReply, error) {
	h.logger.Log(gologger.Info,
		gologger.ExternalServices,
		"PhoneVerifyMerchant",
		"phone verify merchant request parameters",
		"PhoneVerifyMerchantRequest", "", structs.Map(*req), true)
	// check if required fields are passed or not
	storeErr := ""
	if (req.VerificationCode != "") && (req.VerificationType != "") && (req.AccountId != 0) {
		// write logic here
		verifyRequest := io.OTPVerify{
			VerificationCode: req.VerificationCode,
			VerificationType: req.VerificationType,
			AccountId:        req.AccountId,
		}
		masterTx, transactionTx, err := h.GetMasterAndTransactionDBTransaction(ctx)
		if err != nil {
			return nil, err
		}
		defer func() {
			masterTx.RollbackUnlessCommitted()
			transactionTx.RollbackUnlessCommitted()
		}()

		db, err := h.DbRepo.GetDB(ctx, MasterDBConnectionName)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}

		serviceRes := h.athOTPService.VerifyOTP(ctx, masterTx, transactionTx, db, verifyRequest)
		if serviceRes.Error != nil {
			storeErr = serviceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"PhoneVerifyMerchant",
				"phone verify merchant request failed",
				gologger.ParseError, "", structs.Map(serviceRes), true)
			return nil, status.Errorf(codes.Internal, storeErr)
		}

		res := &merchantpb.GenericReply{
			Status: serviceRes.Success,
		}
		h.logger.Log(gologger.Info,
			gologger.ExternalServices,
			"PhoneVerifyMerchant",
			"phone verify merchant response body",
			"PhoneVerifyMerchantResponse", "", structs.Map(*res), true)
		if err = transactionTx.Commit().Error; err != nil {
			return nil, status.Errorf(codes.Internal, "error in AddTeamMember :%v", err.Error())
		}
		return res, masterTx.Commit().Error
	}
	storeErr = "Request parameter is missing or blank"
	h.logger.Log(gologger.Errsev3,
		gologger.ExternalServices,
		"PhoneVerifyMerchant",
		"phone verify merchant request failed",
		gologger.ParseError, "", structs.Map(req), true)
	return nil, status.Errorf(codes.Internal, "error in PhoneVerifyMerchant :%v", storeErr)
}

func (h *athMerchantHandler) PhoneVerifyTeam(ctx context.Context, req *merchantpb.PhoneVerifyTeamRequest) (*merchantpb.PhoneVerifyTeamReply, error) {
	h.logger.Log(gologger.Info,
		gologger.ExternalServices,
		"PhoneVerifyTeam",
		"phone verify team request parameters",
		"PhoneVerifyTeamRequest", "", structs.Map(*req), true)
	// check if required fields are passed or not
	storeErr := ""
	if (req.VerificationCode != "") && (req.VerificationType != "") && (req.Email != "") {
		// write logic here
		verifyRequest := io.OTPVerify{
			VerificationCode: req.VerificationCode,
			VerificationType: req.VerificationType,
			Email:            req.Email,
		}
		masterTx, transactionTx, err := h.GetMasterAndTransactionDBTransaction(ctx)
		if err != nil {
			return nil, err
		}
		defer func() {
			masterTx.RollbackUnlessCommitted()
			transactionTx.RollbackUnlessCommitted()
		}()
		db, err := h.DbRepo.GetDB(ctx, MasterDBConnectionName)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		serviceRes := h.athOTPService.VerifyOTP(ctx, masterTx, transactionTx, db, verifyRequest)
		if serviceRes.Error != nil {
			storeErr = serviceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"PhoneVerifyTeam",
				"phone verify team request failed",
				gologger.ParseError, "", structs.Map(serviceRes), true)
			return nil, status.Errorf(codes.Internal, storeErr)
		}
		userData := serviceRes.Data.(io.AthUser)
		otpNo := utile.RandomString(6)
		if userData.Password == "" {
			// create code to reset password and store it in db
			phoneOTPReq := io.AthUserOTP{
				UserID:    userData.UserID,
				Phone:     userData.Phone,
				OTPType:   "reset_password",
				OTPNO:     otpNo,
				ExpiredAt: int(time.Now().Unix()) + 600,
				CreatedAt: int(time.Now().Unix()),
			}

			// add team member details to athusers table
			otpServiceRes := h.athOTPService.CreateCode(ctx, transactionTx, phoneOTPReq)
			if otpServiceRes.Error != nil {
				storeErr = otpServiceRes.Error.Error()
				h.logger.Log(gologger.Errsev3,
					gologger.ExternalServices,
					"PhoneVerifyTeam",
					"phone verify team request failed",
					gologger.ParseError, "", structs.Map(otpServiceRes), true)
				return nil, status.Errorf(codes.Internal, "error in PhoneVerifyTeam :%v", storeErr)
			}
		}
		res := &merchantpb.PhoneVerifyTeamReply{
			ResetPasswordCode: otpNo,
		}
		h.logger.Log(gologger.Info,
			gologger.ExternalServices,
			"PhoneVerifyTeam",
			"phone verify team response body",
			"PhoneVerifyTeamResponse", "", structs.Map(*res), true)
		if err = transactionTx.Commit().Error; err != nil {
			return nil, status.Errorf(codes.Internal, "error in AddTeamMember :%v", err.Error())
		}
		return res, masterTx.Commit().Error
	}
	storeErr = "Request parameter is missing or blank"
	h.logger.Log(gologger.Errsev3,
		gologger.ExternalServices,
		"PhoneVerifyTeam",
		"phone verify team request failed",
		gologger.ParseError, "", structs.Map(req), true)
	return nil, status.Errorf(codes.Internal, "error in PhoneVerifyTeam :%v", storeErr)
}

func (h *athMerchantHandler) EmailVerifyMerchant(ctx context.Context, req *merchantpb.EmailVerifyRequest) (*merchantpb.GenericReply, error) {
	h.logger.Log(gologger.Info,
		gologger.ExternalServices,
		"EmailVerifyMerchant",
		"email verify merchant request parameters",
		"EmailVerifyMerchantRequest", "", structs.Map(*req), true)
	// check if required fields are passed or not
	storeErr := ""
	if (req.VerificationCode != "") && (req.VerificationType != "") && (req.Email != "") {
		verifyRequest := io.OTPVerify{
			VerificationType: req.VerificationType,
			VerificationCode: req.VerificationCode,
			Email:            req.Email,
		}

		masterTx, transactionTx, err := h.GetMasterAndTransactionDBTransaction(ctx)
		if err != nil {
			return nil, err
		}
		defer func() {
			masterTx.RollbackUnlessCommitted()
			transactionTx.RollbackUnlessCommitted()
		}()
		db, err := h.DbRepo.GetDB(ctx, MasterDBConnectionName)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}

		serviceRes := h.athOTPService.VerifyOTP(ctx, masterTx, transactionTx, db, verifyRequest)
		if serviceRes.Error != nil {
			storeErr = serviceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"EmailVerifyMerchant",
				"email verify merchant request failed",
				gologger.ParseError, "", structs.Map(serviceRes), true)
			return nil, status.Errorf(codes.Internal, storeErr)
		}

		// check if user phone is verified or not
		userData := serviceRes.Data.(io.AthUser)
		if userData.PhoneVerify != true {
			// phone is not verified. Lets verify it
			phoneOTPReq := io.AthUserOTP{
				UserID:    userData.UserID,
				Phone:     userData.Phone,
				OTPType:   "phone",
				OTPNO:     utile.RandomString(6),
				ExpiredAt: int(time.Now().Unix()) + 600,
				CreatedAt: int(time.Now().Unix()),
			}

			// add team member details to athusers table
			otpServiceRes := h.athOTPService.CreateOTP(ctx, masterTx, transactionTx, phoneOTPReq)
			if otpServiceRes.Error != nil {
				storeErr = otpServiceRes.Error.Error()
				h.logger.Log(gologger.Errsev3,
					gologger.ExternalServices,
					"EmailVerifyMerchant",
					"email verify merchant request failed",
					gologger.ParseError, "", structs.Map(otpServiceRes), true)
				return nil, status.Errorf(codes.Internal, "error in EmailVerifyMerchant :%v", storeErr)
			}
		}

		res := &merchantpb.GenericReply{
			Status: serviceRes.Success,
		}
		h.logger.Log(gologger.Info,
			gologger.ExternalServices,
			"EmailVerifyMerchant",
			"email verify merchant response body",
			"EmailVerifyMerchantResponse", "", structs.Map(*res), true)
		if err = transactionTx.Commit().Error; err != nil {
			return nil, status.Errorf(codes.Internal, "error in AddTeamMember :%v", err.Error())
		}
		return res, masterTx.Commit().Error
	}

	storeErr = "Request parameter is missing or blank"
	h.logger.Log(gologger.Errsev3,
		gologger.ExternalServices,
		"EmailVerifyMerchant",
		"email verify merchant request failed",
		gologger.ParseError, "", structs.Map(req), true)
	return nil, status.Errorf(codes.Internal, "error in EmailVerifyMerchant :%v", storeErr)
}

func (h *athMerchantHandler) ForgotPasswordMerchant(ctx context.Context, req *merchantpb.ForgotPasswordRequest) (*merchantpb.GenericReply, error) {
	h.logger.Log(gologger.Info,
		gologger.ExternalServices,
		"ForgotPasswordMerchant",
		"forgot password merchant request parameters",
		"ForgotPasswordMerchantRequest", "", structs.Map(*req), true)
	storeErr := ""
	if req.Email != "" {
		token := utile.RandomString(6)
		// email verification link gets sent
		emailSendReq := io.EmailSendReq{
			EmailVerifyToken: token,
			Email:            req.Email,
			Account:          "WEBIN",
			Subject:          "Please reset your password",
			Body:             "Please click on link to reset password - http://odp-turf-web-01.bigtree.org/change-password?code=" + token + "&email=" + req.Email,
			Eticket:          "N",
			Tid:              "0",
			OTPTtype:         "reset_password",
		}
		masterTx, transactionTx, err := h.GetMasterAndTransactionDBTransaction(ctx)
		if err != nil {
			return nil, err
		}
		defer func() {
			masterTx.RollbackUnlessCommitted()
			transactionTx.RollbackUnlessCommitted()
		}()

		// add team member details to athusers table
		otpServiceRes := h.athOTPService.EmailSend(ctx, masterTx, transactionTx, emailSendReq)
		if otpServiceRes.Error != nil {
			storeErr = otpServiceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"ForgotPasswordMerchant",
				"forgot password merchant request failed",
				gologger.ParseError, "", structs.Map(otpServiceRes), true)
			return nil, status.Errorf(codes.Internal, "error in ForgotPasswordMerchant :%v", storeErr)
		}

		res := &merchantpb.GenericReply{
			Status: otpServiceRes.Success,
		}
		h.logger.Log(gologger.Info,
			gologger.ExternalServices,
			"ForgotPasswordMerchant",
			"forgot password merchant response body",
			"ForgotPasswordMerchantResponse", "", structs.Map(*res), true)
		if err = transactionTx.Commit().Error; err != nil {
			return nil, status.Errorf(codes.Internal, "error in AddTeamMember :%v", err.Error())
		}
		return res, masterTx.Commit().Error
	}
	storeErr = "Request parameter is missing or blank"
	h.logger.Log(gologger.Errsev3,
		gologger.ExternalServices,
		"ForgotPasswordMerchant",
		"forgot password merchant request failed",
		gologger.ParseError, "", structs.Map(req), true)
	return nil, status.Errorf(codes.Internal, "error in ForgotPasswordMerchant :%v", storeErr)

}

func (h *athMerchantHandler) ResetPasswordMerchant(ctx context.Context, req *merchantpb.ResetPasswordRequest) (*merchantpb.GenericReply, error) {
	h.logger.Log(gologger.Info,
		gologger.ExternalServices,
		"ResetPasswordMerchant",
		"reset password user request parameters",
		"ResetPasswordMerchantRequest", "", structs.Map(*req), true)
	storeErr := ""
	if (req.NewPassword != "") && (req.ResetPasswordToken != "") && (req.Email != "") {
		// check if OTP is valid
		verifyRequest := io.OTPVerify{
			VerificationCode: req.ResetPasswordToken,
			VerificationType: "reset_password",
			Email:            req.Email,
		}
		masterTx, transactionTx, err := h.GetMasterAndTransactionDBTransaction(ctx)
		if err != nil {
			return nil, err
		}
		defer func() {
			masterTx.RollbackUnlessCommitted()
			transactionTx.RollbackUnlessCommitted()
		}()

		db, err := h.DbRepo.GetDB(ctx, MasterDBConnectionName)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		serviceRes := h.athOTPService.VerifyOTP(ctx, masterTx, transactionTx, db, verifyRequest)
		if serviceRes.Error != nil {
			storeErr = serviceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"ResetPasswordMerchant",
				"reset password user request failed",
				gologger.ParseError, "", structs.Map(serviceRes), true)
			return nil, status.Errorf(codes.Internal, storeErr)
		}

		// check if new password is valid
		request := io.ResetPasswordRequest{
			Email:    req.Email,
			Password: req.NewPassword,
			Token:    req.ResetPasswordToken,
		}
		fmt.Println("calling reset password user service function---")
		userServiceRes := h.athUserService.ResetPasswordUser(ctx, masterTx, request)
		fmt.Println("userServiceRes : ", userServiceRes)
		if userServiceRes.Error != nil {
			storeErr = userServiceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"ResetPasswordMerchant",
				"reset password merchant request failed",
				gologger.ParseError, "", structs.Map(userServiceRes), true)
			return nil, status.Errorf(codes.Internal, storeErr)
		}

		userData := userServiceRes.Data.(io.ResetPasswordRequest)
		fmt.Println("userData : ", userData)

		//send an email saying - Password has been reset successfully
		emailSendReq := io.EmailSendReq{
			EmailVerifyToken: "",
			Email:            userData.Email,
			Account:          "WEBIN",
			Subject:          "Reset Password successfully",
			Body:             "Password has been reset successfully",
			Eticket:          "N",
			Tid:              "0",
			OTPTtype:         "email",
		}
		fmt.Println("emailSendReq : ", emailSendReq)
		otpServiceRes := h.athOTPService.EmailSend(ctx, masterTx, transactionTx, emailSendReq)
		if otpServiceRes.Error != nil {
			storeErr = otpServiceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"ResetPasswordMerchant",
				"reset password merchant request failed",
				gologger.ParseError, "", structs.Map(otpServiceRes), true)
			return nil, status.Errorf(codes.Internal, "error in ResetPasswordMerchant :%v", storeErr)
		}

		res := &merchantpb.GenericReply{
			Status: userServiceRes.Success,
		}
		h.logger.Log(gologger.Info,
			gologger.ExternalServices,
			"ResetPasswordMerchant",
			"reset password merchant response body",
			"ResetPasswordMerchantResponse", "", structs.Map(*res), true)
		if err = transactionTx.Commit().Error; err != nil {
			return nil, status.Errorf(codes.Internal, "error in AddTeamMember :%v", err.Error())
		}
		return res, masterTx.Commit().Error
	}
	storeErr = "Request parameter is missing or blank"
	h.logger.Log(gologger.Errsev3,
		gologger.ExternalServices,
		"ResetPasswordMerchant",
		"reset password merchant request failed",
		gologger.ParseError, "", structs.Map(req), true)
	return nil, status.Errorf(codes.Internal, "error in ResetPasswordMerchant :%v", storeErr)
}

func (h *athMerchantHandler) AddTeamMember(ctx context.Context, req *merchantpb.AddTeamMemberRequest) (*merchantpb.AddTeamMemberReply, error) {
	storeErr := ""
	if (req.FirstName != "") && (req.LastName != "") && (req.Email != "") && (req.Phone != "") && (req.AccountType != "") {
		loggedInUserID, err := utile.GetUserIDFromContext(ctx)
		if err != nil {
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"AddTeamMember",
				"Add team member request failed",
				gologger.ParseError, "", structs.Map(*req), true)
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		h.logger.Log(gologger.Info,
			gologger.ExternalServices,
			"AddTeamMember",
			"Add team member request parameters",
			"AddTeamMemberRequest", "", structs.Map(*req), true)
		masterTx, transactionTx, err := h.GetMasterAndTransactionDBTransaction(ctx)
		if err != nil {
			return nil, err
		}
		defer func() {
			masterTx.RollbackUnlessCommitted()
			transactionTx.RollbackUnlessCommitted()
		}()

		userRequest := io.AthUser{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Phone:     req.Phone,
			Models: io.Models{
				CreatedBy: loggedInUserID,
			},
			UserSource: "merchant app",
		}

		// add users details to athusers table
		userServiceRes := h.athUserService.CreateUser(ctx, masterTx, userRequest, "teamMember")
		if userServiceRes.Error != nil {
			storeErr = userServiceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"AddTeamMember",
				"Add team member request failed",
				gologger.ParseError, "", structs.Map(userServiceRes), true)
			return nil, status.Errorf(codes.Internal, "error in AddTeamMember :%v", storeErr)
		}

		userData := userServiceRes.Data.(io.AthUser)
		userID := userData.UserID
		// add team member privileges to table
		addTeamMemberReq := io.TeamMemberReq{
			UserID:      userID,
			AccountType: req.AccountType,
			AccessData:  req.Privileges.VenueIds,
			CreatedBy:   loggedInUserID,
			UpdatedBy:   loggedInUserID,
		}

		// add team member details to athusers table
		teamMemberServiceRes := h.athMerchantService.AddTeamMember(ctx, masterTx, addTeamMemberReq)
		if teamMemberServiceRes.Error != nil {
			storeErr = teamMemberServiceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"AddTeamMember",
				"Add team member request failed",
				gologger.ParseError, "", structs.Map(teamMemberServiceRes), true)
			return nil, status.Errorf(codes.Internal, "error in AddTeamMember :%v", storeErr)
		}
		emailVerify := utile.RandomString(6)
		// email verification link gets sent
		emailSendReq := io.EmailSendReq{
			EmailVerifyToken: emailVerify,
			Email:            req.Email,
			Account:          "WEBIN",
			Subject:          "Please Verify your email",
			Body:             "Email verification token is - " + emailVerify,
			Eticket:          "N",
			Tid:              "0",
			OTPTtype:         "email",
		}

		// add team member details to athusers table
		otpServiceRes := h.athOTPService.EmailSend(ctx, masterTx, transactionTx, emailSendReq)
		if otpServiceRes.Error != nil {
			storeErr = otpServiceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"AddTeamMember",
				"Add team member request failed",
				gologger.ParseError, "", structs.Map(otpServiceRes), true)
			return nil, status.Errorf(codes.Internal, "error in AddTeamMember :%v", storeErr)
		}

		// get basic details
		serviceRes := h.athUserService.GetAccountDetailsByID(ctx, masterTx, userID)
		if serviceRes.Error != nil {
			storeErr = serviceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"LoginMerchant",
				"login merchant request failed",
				gologger.ParseError, "", structs.Map(serviceRes), true)
			return nil, status.Errorf(codes.Internal, storeErr)
		}

		accountData := serviceRes.Data.(map[string]interface{})

		res := &merchantpb.AddTeamMemberReply{
			FirstName:   userData.FirstName,
			LastName:    userData.LastName,
			AccountType: accountData["AccountType"].(string),
			Privileges: &merchantpb.AccessData{
				VenueIds: accountData["AccessData"].([]int32),
			},
			Phone:     userData.Phone,
			Email:     userData.Email,
			AccountId: int32(userID),
		}
		h.logger.Log(gologger.Info,
			gologger.ExternalServices,
			"AddTeamMember",
			"Add team member response body",
			"AddTeamMemberResponse", "", structs.Map(*res), true)
		if err = transactionTx.Commit().Error; err != nil {
			return nil, status.Errorf(codes.Internal, "error in AddTeamMember :%v", err.Error())
		}
		return res, masterTx.Commit().Error
	}
	storeErr = "Request parameter is missing or blank"
	h.logger.Log(gologger.Errsev3,
		gologger.ExternalServices,
		"AddTeamMember",
		"Add team member request failed",
		gologger.ParseError, "", structs.Map(req), true)
	return nil, status.Errorf(codes.Internal, "error in AddTeamMember :%v", storeErr)
}
