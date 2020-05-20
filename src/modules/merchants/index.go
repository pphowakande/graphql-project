package handler

import (
	merchandal "api/src/dal/merchant"
	"api/src/service/merchant"
	"api/src/service/otp"
	"api/src/service/upload"
	"api/src/service/user"
	"context"

	"github.com/fatih/structs"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"stash.bms.bz/bms/gologger.git"

	merchantpb "stash.bms.bz/turf/generic-proto-files.git/merchant/v1"
)

const (
	// MasterDBConnectionName is the key for fetching the master DB connection
	MasterDBConnectionName = "master"

	// TransactionDBConnectionName is the key for fetching the Transactional DB connection
	TransactionDBConnectionName = "transactional"
)

type athMerchantHandler struct {
	logger             *gologger.Logger
	DbRepo             merchandal.Repository
	athMerchantService merchant.AthMerchantService
	athUserService     user.AthUserService
	athOTPService      otp.AthOtpService
	athUploadService   upload.AthUploadService
}

func NewMerchantHandler(logger *gologger.Logger, DbRepo merchandal.Repository, service merchant.AthMerchantService,
	athUserService user.AthUserService,
	athOTPService otp.AthOtpService,
	athUploadService upload.AthUploadService) merchantpb.MerchantServer {
	return &athMerchantHandler{
		logger:             logger,
		DbRepo:             DbRepo,
		athMerchantService: service,
		athUserService:     athUserService,
		athOTPService:      athOTPService,
		athUploadService:   athUploadService,
	}
}

// GetMasterAndTransactionDBTransaction returns the master and transactionDB transaction objects
func (h *athMerchantHandler) GetMasterAndTransactionDBTransaction(ctx context.Context) (*gorm.DB, *gorm.DB, error) {
	masterTx, err := h.DbRepo.BeginTransaction(ctx, MasterDBConnectionName)
	if err != nil {
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"Merchant",
			"Merchant request failed",
			gologger.ParseError, "",
			structs.Map(map[string]interface{}{"error": "error starting master transaction"}), true)
		return nil, nil, status.Errorf(codes.Internal, "error in Merchant handler :error starting transaction")
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
