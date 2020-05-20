package handler

import (
	"context"

	"github.com/fatih/structs"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"stash.bms.bz/bms/gologger.git"

	venuedal "api/src/dal/venue"
	"api/src/service/upload"
	"api/src/service/venue"

	venuepb "stash.bms.bz/turf/generic-proto-files.git/venue/v1"
)

const (
	// MasterDBConnectionName is the key for fetching the master DB connection
	MasterDBConnectionName = "master"

	// TransactionDBConnectionName is the key for fetching the Transactional DB connection
	TransactionDBConnectionName = "transactional"
)

type athVenueHandler struct {
	logger           *gologger.Logger
	DbRepo           venuedal.Repository
	athVenueService  venue.AthVenueService
	athUploadService upload.AthUploadService
}

func NewVenueHandler(logger *gologger.Logger, DbRepo venuedal.Repository, service venue.AthVenueService, athUploadService upload.AthUploadService) venuepb.VenueServer {
	return &athVenueHandler{
		logger:           logger,
		athVenueService:  service,
		DbRepo:           DbRepo,
		athUploadService: athUploadService,
	}
}

// GetMasterAndTransactionDBTransaction returns the master and transactionDB transaction objects
func (h *athVenueHandler) GetMasterDBTransaction(ctx context.Context) (*gorm.DB, error) {
	masterTx, err := h.DbRepo.BeginTransaction(ctx, MasterDBConnectionName)
	if err != nil {
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"Venue",
			"Venue request failed",
			gologger.ParseError, "",
			structs.Map(map[string]interface{}{"error": "error starting master transaction"}), true)
		return nil, status.Errorf(codes.Internal, "error in Venue handler :error starting transaction")
	}
	return masterTx, nil
}
