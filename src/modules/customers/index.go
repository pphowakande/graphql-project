package handler

import (
	"stash.bms.bz/bms/gologger.git"

	"api/src/service/customer"
	"api/src/service/otp"

	customerpb "stash.bms.bz/turf/generic-proto-files.git/customer/v1"
)

type athCustomerHandler struct {
	logger             *gologger.Logger
	athCustomerService customer.AthCustomerService
	athOtpService      otp.AthOtpService
}

func NewCustomerHandler(logger *gologger.Logger, service customer.AthCustomerService,
	athOtpService otp.AthOtpService) customerpb.CustomerServer {
	return &athCustomerHandler{
		logger:             logger,
		athCustomerService: service,
		athOtpService:      athOtpService,
	}
}
