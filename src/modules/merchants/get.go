package handler

import (
	io "api/src/models"
	utile "api/src/utils"
	"context"

	"github.com/fatih/structs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"stash.bms.bz/bms/gologger.git"

	merchantpb "stash.bms.bz/turf/generic-proto-files.git/merchant/v1"
)

func (h *athMerchantHandler) GetMerchantByID(ctx context.Context, req *merchantpb.GenericRequest) (*merchantpb.GetMerchantByIDRes, error) {
	// get userid from auth token
	loggedInUserID, err := utile.GetUserIDFromContext(ctx)
	if err != nil {
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"GetMerchantByID",
			"get Merchant byID  request failed",
			gologger.ParseError, "", structs.Map(*req), true)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	h.logger.Log(gologger.Info,
		gologger.ExternalServices,
		"getMerchantByID",
		"get Merchant byID request parameters",
		"getMerchantByIDRequest", "", structs.Map(*req), true)
	storeErr := ""

	res := &merchantpb.GetMerchantByIDRes{}
	masterDB, err := h.DbRepo.GetDB(ctx, MasterDBConnectionName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	// get userid from auth token and get merchant id from user id
	serviceRes := h.athMerchantService.GetMerchantByUserID(ctx, masterDB, int(loggedInUserID))
	if serviceRes.Error != nil {
		storeErr = serviceRes.Error.Error()
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"getMerchantByID",
			"get Merchant byID request failed",
			gologger.ParseError, "", structs.Map(serviceRes), true)
		return nil, status.Errorf(codes.Internal, storeErr)
	}

	MerchantData := serviceRes.Data.(io.AthMerchant)

	// get user by id
	UserServiceRes := h.athUserService.GetUserByID(ctx, masterDB, loggedInUserID, false)
	if UserServiceRes.Error != nil {
		storeErr = UserServiceRes.Error.Error()
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"getMerchantByID",
			"get Merchant byID request failed",
			gologger.ParseError, "", structs.Map(UserServiceRes), true)
		return nil, status.Errorf(codes.Internal, storeErr)
	}

	UserData := UserServiceRes.Data.(io.AthUser)
	KYCData := &merchantpb.KycData{
		BankAccFile: MerchantData.BankAccFile,
		GstNoFile:   MerchantData.GstNoFile,
		AddressFile: MerchantData.AddressFile,
		PanNoFile:   MerchantData.PanNoFile,
	}
	res = &merchantpb.GetMerchantByIDRes{
		MerchantFullName: UserData.FirstName + " " + UserData.LastName,
		AccountId:        int32(MerchantData.MerchantID),
		BusinessName:     MerchantData.MerchantName,
		Address:          MerchantData.Address,
		Phone:            UserData.Phone,
		Email:            UserData.Email,
		KycData:          KYCData,
		Password:         UserData.Password,
	}
	h.logger.Log(gologger.Info,
		gologger.ExternalServices,
		"getMerchantByID",
		"get Merchant byID response body",
		"getMerchantByIDResponse", "", structs.Map(*res), true)
	return res, nil
}

func (h *athMerchantHandler) LoginMerchant(ctx context.Context, req *merchantpb.LoginRequest) (*merchantpb.LoginReply, error) {
	h.logger.Log(gologger.Info,
		gologger.ExternalServices,
		"LoginMerchant",
		"login merchant request parameters",
		"LoginMerchantRequest", "", structs.Map(*req), true)
	storeErr := ""
	// check if fields are not empty
	if (req.Login != "") && (req.Password != "") {
		request := io.LoginRequest{
			Login:    req.Login,
			Password: req.Password,
		}
		masterDB, err := h.DbRepo.GetDB(ctx, MasterDBConnectionName)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		// get basic details
		serviceRes := h.athUserService.LoginUser(ctx, masterDB, request)
		if serviceRes.Error != nil {
			storeErr = serviceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"LoginMerchant",
				"login merchant request failed",
				gologger.ParseError, "", structs.Map(serviceRes), true)
			return nil, status.Errorf(codes.Internal, storeErr)
		}

		res := &merchantpb.LoginReply{}
		loginData := serviceRes.Data.(map[string]interface{})

		// if logged in user is merchant only then check for KYC documents
		doLater := false
		BusinessName := ""
		if loginData["AccountType"].(string) == "owner" {
			// get userid from auth token and get merchant id from user id
			serviceRes = h.athMerchantService.GetMerchantByUserID(ctx, masterDB, loginData["MerchantId"].(int))
			if serviceRes.Error != nil {
				storeErr = serviceRes.Error.Error()
				h.logger.Log(gologger.Errsev3,
					gologger.ExternalServices,
					"LoginMerchant",
					"login merchant request failed",
					gologger.ParseError, "", structs.Map(serviceRes), true)
				return nil, status.Errorf(codes.Internal, storeErr)
			}

			MerchantData := serviceRes.Data.(io.AthMerchant)
			doLater = MerchantData.DoLater
			BusinessName = MerchantData.MerchantName
		}

		res = &merchantpb.LoginReply{
			AccountId:        int32(loginData["MerchantId"].(int)),
			MerchantFullName: loginData["MerchantFullName"].(string),
			BusinessName:     BusinessName,
			AccountType:      loginData["AccountType"].(string),
			Privileges: &merchantpb.AccessData{
				VenueIds: loginData["AccessData"].([]int32),
			},
			EmailVerify: loginData["EmailVerify"].(bool),
			BmsVerify:   loginData["BMSVerify"].(bool),
			PhoneVerify: loginData["PhoneVerify"].(bool),
			Token:       loginData["Token"].(string),
			LastLoginAt: int32(loginData["LastLoginAt"].(int)),
			KycPending:  doLater,
		}

		h.logger.Log(gologger.Info,
			gologger.ExternalServices,
			"LoginMerchant",
			"LoginMerchant response body",
			"LoginMerchantResponse", "", structs.Map(*res), true)
		return res, nil
	}

	storeErr = "Request parameter is missing or blank"
	h.logger.Log(gologger.Errsev3,
		gologger.ExternalServices,
		"LoginMerchant",
		"login merchant request failed",
		gologger.ParseError, "", structs.Map(req), true)
	return nil, status.Errorf(codes.Internal, "error in LoginMerchant :%v", storeErr)
}

func (h *athMerchantHandler) GetMerchantTeam(ctx context.Context, req *merchantpb.GetMerchantTeamRequest) (*merchantpb.GetMerchantTeamRes, error) {
	loggedInUserID, err := utile.GetUserIDFromContext(ctx)
	if err != nil {
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"GetMerchantTeam",
			"Get Merchant team request failed",
			gologger.ParseError, "", structs.Map(*req), true)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	storeErr := ""
	masterDB, err := h.DbRepo.GetDB(ctx, MasterDBConnectionName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	// get team member  basic details
	serviceRes := h.athMerchantService.GetTeamData(ctx, masterDB, loggedInUserID, req.OrderBy)
	if serviceRes.Error != nil {
		storeErr = serviceRes.Error.Error()
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"GetMerchantTeam",
			"Get Merchant team request failed",
			gologger.ParseError, "", structs.Map(serviceRes), true)
		return nil, status.Errorf(codes.Internal, storeErr)
	}

	data := serviceRes.Data.([]*merchantpb.TeamMemberData)
	res := &merchantpb.GetMerchantTeamRes{
		TeamData: data,
	}
	return res, nil
}
