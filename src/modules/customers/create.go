package handler

import (
	"context"

	customerpb "stash.bms.bz/turf/generic-proto-files.git/customer/v1"
)

func (h *athCustomerHandler) CreateCustomer(ctx context.Context, req *customerpb.CreateCustomerRequest) (*customerpb.CreateCustomerReply, error) {
	/*
		loggedInUserID, err := utile.GetUserIDFromContext(ctx)
		if err != nil {
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"getUserByID",
				"get user byID  request failed",
				gologger.ParseError, "", structs.Map(*req), true)
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		h.logger.Log(gologger.Info,
			gologger.ExternalServices,
			"createCustomer",
			"create customer request parameters",
			"createCustomerRequest", "", structs.Map(*req), true)
		userRequest := io.AthUser{
			Email:      req.Email,
			UserSource: req.UserSource,
			Models: io.Models{
				CreatedBy: loggedInUserID,
			},
		}

		storeErr := ""
		res := &customerpb.CreateCustomerReply{}

		createCustomerServiceRes := h.athCustomerService.CreateCustomer(ctx, userRequest)

		if createCustomerServiceRes.Error != nil {
			storeErr = createCustomerServiceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"createCustomer",
				"create customer request failed",
				gologger.ParseError, "", structs.Map(createCustomerServiceRes), true)
			return nil, status.Errorf(codes.Internal, storeErr)
		}

		customerData := createCustomerServiceRes.Data.(map[string]interface{})
		userID := customerData["user_id"].(int)

		// profileRequest := io.AthUserProfile{
		// 	FirstName: req.FirstName,
		// 	LastName:  req.LastName,
		// 	UserID:    userID,
		// 	phone: req.Phone,
		// 	Models: io.Models{
		// 		CreatedBy: loggedInUserID,
		// 	},
		// }

		// profileServiceRes := h.athUserService.CreateUserProfile(ctx, profileRequest)
		// if profileServiceRes.Error != nil {
		// 	storeErr = profileServiceRes.Error.Error()
		// 	h.logger.Log(gologger.Errsev3,
		// 		gologger.ExternalServices,
		// 		"createCustomer",
		// 		"create customer request failed",
		// 		gologger.ParseError, "", structs.Map(profileServiceRes), true)
		// 	return nil, status.Errorf(codes.Internal, storeErr)
		// }

		// trigger create token service
		var otpRequest io.AthUserOTP

		otpRequest.UserID = userID
		otpRequest.OTPNO = utile.RandomString(6)
		//otpRequest.OTPExpiry = 3600
		otpRequest.OTPType = "phone"
		otpRequest.Phone = req.Phone

		otpServiceRes := h.athOtpService.CreateOTP(ctx, otpRequest)
		if otpServiceRes.Error != nil {
			storeErr = otpServiceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"createCustomer",
				"create customer request failed",
				gologger.ParseError, "", structs.Map(otpServiceRes), true)
			return nil, status.Errorf(codes.Internal, storeErr)
		}

		res = &customerpb.CreateCustomerReply{
			Data: &customerpb.CreateCustomerReplyData{
				CustomerId: int32(userID),
			},
		}
		h.logger.Log(gologger.Info,
			gologger.ExternalServices,
			"createCustomer",
			"create customer response body",
			"createCustomerResponse", "", structs.Map(*res), true)
		return res, nil
	*/
	return nil, nil
}
