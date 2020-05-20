package handler

import (
	"context"

	customerpb "stash.bms.bz/turf/generic-proto-files.git/customer/v1"
)

func (h *athCustomerHandler) EditCustomer(ctx context.Context, req *customerpb.EditCustomerRequest) (*customerpb.GenericReply, error) {
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
			"editCustomer",
			"edit customer request parameters",
			"editCustomerRequest", "", structs.Map(*req), true)
		storeErr := ""
		res := &customerpb.GenericReply{}

		// profileRequest := io.AthUserProfile{
		// 	FirstName: req.FirstName,
		// 	LastName:  req.LastName,
		// 	phone: req.Phone,
		// 	UserID:    loggedInUserID,
		// 	Models: io.Models{
		// 		UpdatedBy: loggedInUserID,
		// 	},
		// }

		// profileServiceRes := h.athUserService.EditUserProfile(ctx, profileRequest)
		// if profileServiceRes.Error != nil {
		// 	storeErr = profileServiceRes.Error.Error()
		// 	res.Status = profileServiceRes.Success
		// 	h.logger.Log(gologger.Errsev3,
		// 		gologger.ExternalServices,
		// 		"editCustomer",
		// 		"edit customer request failed",
		// 		gologger.ParseError, "", structs.Map(profileServiceRes), true)
		// 	return res, status.Errorf(codes.Internal, storeErr)
		// }

		if req.Phone != "" {
			// trigger create token service
			var otpRequest io.AthUserOTP

			otpRequest.UserID = loggedInUserID
			otpRequest.OTPNO = utile.RandomString(6)
			//otpRequest.OTPExpiry = 3600
			otpRequest.OTPType = "phone"
			otpRequest.Phone = req.Phone

			otpServiceRes := h.athOtpService.CreateOTP(ctx, otpRequest)
			if otpServiceRes.Error != nil {
				storeErr = otpServiceRes.Error.Error()
				res.Status = otpServiceRes.Success
				h.logger.Log(gologger.Errsev3,
					gologger.ExternalServices,
					"editCustomer",
					"edit customer request failed",
					gologger.ParseError, "", structs.Map(otpServiceRes), true)
				return res, status.Errorf(codes.Internal, storeErr)
			}
		}

		res = &customerpb.GenericReply{
			//Status: otpServiceRes.Success,
		}
		h.logger.Log(gologger.Info,
			gologger.ExternalServices,
			"editCustomer",
			"edit customer response body",
			"editCustomerResponse", "", structs.Map(*res), true)
		return res, nil
	*/
	return nil, nil
}
