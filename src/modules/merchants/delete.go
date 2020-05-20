package handler

import (
	io "api/src/models"
	utile "api/src/utils"
	"context"
	"time"

	"github.com/fatih/structs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"stash.bms.bz/bms/gologger.git"

	merchantpb "stash.bms.bz/turf/generic-proto-files.git/merchant/v1"
)

func (h *athMerchantHandler) DeleteTeamMember(ctx context.Context, req *merchantpb.DeleteTeamMemberRequest) (*merchantpb.GenericReply, error) {
	loggedInUserID, err := utile.GetUserIDFromContext(ctx)
	if err != nil {
		h.logger.Log(gologger.Errsev3,
			gologger.ExternalServices,
			"DeleteTeamMember",
			"delete team member request failed",
			gologger.ParseError, "", structs.Map(*req), true)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	storeErr := ""
	// check if required fields are not empty
	if req.AccountId != 0 {
		h.logger.Log(gologger.Info,
			gologger.ExternalServices,
			"DeleteTeamMember",
			"delete team member request parameters",
			"DeleteTeamMemberRequest", "", structs.Map(*req), true)

		venueUserRequest := io.AthVenueUser{
			UserId: int(req.AccountId),
			Models: io.Models{
				CreatedBy: loggedInUserID,
				UpdatedBy: loggedInUserID,
				UpdatedAt: int(time.Now().Unix()),
				//DeletedBy: loggedInUserID,
				//DeletedAt: int(time.Now().Unix()),
			},
		}
		masterDB, err := h.DbRepo.GetDB(ctx, MasterDBConnectionName)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}

		// delete team member
		deletTeamMemberServiceRes := h.athMerchantService.DeleteTeamMember(ctx, masterDB, venueUserRequest)
		if deletTeamMemberServiceRes.Error != nil {
			storeErr = deletTeamMemberServiceRes.Error.Error()
			h.logger.Log(gologger.Errsev3,
				gologger.ExternalServices,
				"DeleteTeamMember",
				"delete team member request failed",
				gologger.ParseError, "", structs.Map(deletTeamMemberServiceRes), true)
			return nil, status.Errorf(codes.Internal, storeErr)
		}
		res := &merchantpb.GenericReply{
			Status: deletTeamMemberServiceRes.Success,
		}
		h.logger.Log(gologger.Info,
			gologger.ExternalServices,
			"DeleteTeamMember",
			"delete team memberresponse body",
			"DeleteTeamMemberResponse", "", structs.Map(*res), true)
		return res, nil
	}
	storeErr = "Request parameter is missing or blank"
	h.logger.Log(gologger.Errsev3,
		gologger.ExternalServices,
		"DeleteTeamMember",
		"delete team member request failed",
		gologger.ParseError, "", structs.Map(req), true)
	return nil, status.Errorf(codes.Internal, "error in DeleteTeamMember :%v", storeErr)

}
