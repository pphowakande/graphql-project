package handler

import (
	"context"

	facilitydal "api/src/dal/facility"

	"github.com/fatih/structs"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"stash.bms.bz/bms/gologger.git"

	"api/src/service/facility"

	facilitypb "stash.bms.bz/turf/generic-proto-files.git/facility/v1"
)

type athFacilityHandler struct {
	logger             *gologger.Logger
	athFacilityService facility.AthFacilityService
	DbRepo             facilitydal.Repository
	//	athOTPService      otp.AthOtpService
}

const (
	// MasterDBConnectionName is the key for fetching the master DB connection
	MasterDBConnectionName = "master"

	// TransactionDBConnectionName is the key for fetching the Transactional DB connection
	TransactionDBConnectionName = "transactional"
)

func NewFacilityHandler(logger *gologger.Logger, DbRepo facilitydal.Repository, service facility.AthFacilityService) facilitypb.FacilityServer {
	return &athFacilityHandler{
		logger:             logger,
		athFacilityService: service,
		DbRepo:             DbRepo,
		//athOTPService:      athOTPService,
	}
}

// GetMasterAndTransactionDBTransaction returns the master and transactionDB transaction objects
func (h *athFacilityHandler) GetMasterDBTransaction(ctx context.Context) (*gorm.DB, error) {
	masterTx, err := h.DbRepo.BeginTransaction(ctx, MasterDBConnectionName)
	if err != nil {
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"Facility",
			"Facility request failed",
			gologger.ParseError, "",
			structs.Map(map[string]interface{}{"error": "error starting master transaction"}), true)
		return nil, status.Errorf(codes.Internal, "error in Facility handler :error starting transaction")
	}
	return masterTx, nil
}

// GetMasterAndTransactionDBTransaction returns the master and transactionDB transaction objects
func (h *athFacilityHandler) GetMasterAndTransactionDBTransaction(ctx context.Context) (*gorm.DB, *gorm.DB, error) {
	masterTx, err := h.DbRepo.BeginTransaction(ctx, MasterDBConnectionName)
	if err != nil {
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"Facility",
			"Facility request failed",
			gologger.ParseError, "",
			structs.Map(map[string]interface{}{"error": "error starting master transaction"}), true)
		return nil, nil, status.Errorf(codes.Internal, "error in Facility handler :error starting transaction")
	}
	transactionTx, err := h.DbRepo.BeginTransaction(ctx, TransactionDBConnectionName)
	if err != nil {
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"Merchant",
			"Merchant request failed",
			gologger.ParseError, "",
			structs.Map(map[string]interface{}{"error": "error starting transaction"}), true)
		return nil, nil, status.Errorf(codes.Internal, "error in Merchant :error starting transaction")
	}
	return masterTx, transactionTx, nil
}
