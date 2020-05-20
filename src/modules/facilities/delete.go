package handler

import (
	io "api/src/models"
	utile "api/src/utils"
	"context"

	"github.com/fatih/structs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"stash.bms.bz/bms/gologger.git"

	facilitypb "stash.bms.bz/turf/generic-proto-files.git/facility/v1"
)

func (h *athFacilityHandler) DeleteFacilityByID(ctx context.Context, req *facilitypb.DeleteFacilityByIDReq) (*facilitypb.GenericReply, error) {
	loggedInUserID, err := utile.GetUserIDFromContext(ctx)
	if err != nil {
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"CreateFacility",
			"create facility  request failed",
			gologger.ParseError, "", structs.Map(*req), true)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	h.logger.Log(gologger.Info,
		gologger.ExternalServices,
		"DeleteFacilityByID",
		"delete facility by ID request parameters",
		"DeleteFacilityByIDRequest", "", structs.Map(*req), true)
	storeErr := ""

	res := &facilitypb.GenericReply{}
	if req.FacilityId == 0 || req.Type == "" {
		res = &facilitypb.GenericReply{}
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"DeleteFacilityByID",
			"delete facility by ID request failed",
			gologger.ParseError, "", structs.Map(*req), true)
		return nil, status.Errorf(codes.InvalidArgument, "Parameter is missing")
	}

	if req.Type != "delete" && req.Type != "unavailable" {
		res = &facilitypb.GenericReply{}
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"DeleteFacilityByID",
			"delete facility by ID request failed",
			gologger.ParseError, "", structs.Map(*req), true)
		return nil, status.Errorf(codes.InvalidArgument, "Wrong Type value passed")
	}

	masterTx, transactionTx, err := h.GetMasterAndTransactionDBTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		masterTx.RollbackUnlessCommitted()
		transactionTx.RollbackUnlessCommitted()
	}()

	var facilityReq io.DeleteFacilityByID
	facilityReq.FacilityID = req.FacilityId
	facilityReq.UserID = int32(loggedInUserID)
	facilityReq.Type = req.Type

	facilityServiceRes := h.athFacilityService.DeleteFacilityByID(ctx, masterTx, transactionTx, facilityReq)
	if facilityServiceRes.Error != nil {
		storeErr = facilityServiceRes.Error.Error()
		res.Status = facilityServiceRes.Success
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"DeleteFacilityByID",
			"delete facility by ID request failed",
			gologger.ParseError, "", structs.Map(facilityServiceRes), true)
		return nil, status.Errorf(codes.Internal, storeErr)
	}

	// get deleted facility
	//faciityData := facilityServiceRes.Data.(io.AthFacilities)
	res = &facilitypb.GenericReply{
		Status: facilityServiceRes.Success,
	}

	h.logger.Log(gologger.Info,
		gologger.ExternalServices,
		"DeleteFacilityByID",
		"delete facility by ID response body",
		"DeleteFacilityByIDResponse", "", structs.Map(*res), true)
	return res, masterTx.Commit().Error

}
