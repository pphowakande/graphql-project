package merchant

import (
	"api/src/dal/merchant"
	"api/src/dal/user"
	io "api/src/models"
	"context"
	"fmt"

	"github.com/jinzhu/gorm"

	"stash.bms.bz/bms/gologger.git"

	"github.com/fatih/structs"
	merchantpb "stash.bms.bz/turf/generic-proto-files.git/merchant/v1"
)

type AthMerchantService interface {
	CreateMerchant(ctx context.Context, db *gorm.DB, u io.AthMerchant) (res io.Response)
	EditMerchant(ctx context.Context, db *gorm.DB, u io.AthMerchant) (res io.Response)
	CreateMerchantUser(ctx context.Context, db *gorm.DB, u io.AthMerchantUser) (res io.Response)
	GetMerchantByID(ctx context.Context, db *gorm.DB, MerchantId int) (res io.Response)
	GetMerchantByUserID(ctx context.Context, db *gorm.DB, UserId int) (res io.Response)
	VerifyGetMerchant(ctx context.Context, db *gorm.DB, MerchantId int32) (res io.Response)
	AddTeamMember(ctx context.Context, db *gorm.DB, u io.TeamMemberReq) (res io.Response)
	UpdateTeamMemberPrivileges(ctx context.Context, db *gorm.DB, u io.AthVenueUser) (res io.Response)
	DeleteTeamMember(ctx context.Context, db *gorm.DB, u io.AthVenueUser) (res io.Response)
	GetTeamData(ctx context.Context, db *gorm.DB, UserID int, orderBy string) (res io.Response)
}

type athMerchantService struct {
	Logger   *gologger.Logger
	DbRepo   merchant.Repository
	userRepo user.Repository
	//logger     log.Logger
}

func NewBasicAthMerchantService(logger *gologger.Logger, DbRepo merchant.Repository, userRepo user.Repository) AthMerchantService {
	return &athMerchantService{
		Logger:   logger,
		DbRepo:   DbRepo,
		userRepo: userRepo,
	}
}

func (b *athMerchantService) CreateMerchant(ctx context.Context, db *gorm.DB, u io.AthMerchant) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"CreateMerchant",
		"create merchant service request parameters",
		"CreateMerchantRequest", "", structs.Map(u), true)
	newMerchant, err := b.DbRepo.CreateMerchant(ctx, db, u)
	if err != nil {
		res = io.FailureMessage(res.Error, "email already exists")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"createMerchant",
			"create merchant service request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}
	data := make(map[string]interface{})
	data["Merchant_id"] = newMerchant.MerchantID
	res.Data = data
	res = io.SuccessMessage(data, "merchant created")

	fmt.Println("response from service : ", res)
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"createMerchant",
		"create merchant service response body",
		"createMerchant service Response", "", structs.Map(res), true)
	return
}

func (b *athMerchantService) EditMerchant(ctx context.Context, db *gorm.DB, u io.AthMerchant) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"EditMerchant",
		"edit merchant request parameters",
		"EditMerchantRequest", "", structs.Map(u), true)

	merchantUpdated, err := b.DbRepo.EditMerchantByID(ctx, db, u)
	fmt.Println("Final updated merchantUpdated : ", merchantUpdated)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error updating merchant")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"editMerchant",
			"edit merchant request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}
	res = io.SuccessMessage(merchantUpdated, "merchant updated")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"EditMerchant",
		"edit merchant response body",
		"EditMerchantResponse", "", structs.Map(res), true)
	return
}

func (b *athMerchantService) CreateMerchantUser(ctx context.Context, db *gorm.DB, u io.AthMerchantUser) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"CreateMerchantUser",
		"create merchant user request parameters",
		"CreatemerchantUserRequest", "", structs.Map(u), true)
	err := b.DbRepo.CreateMerchantUser(ctx, db, u)
	if err != nil {
		res = io.FailureMessage(res.Error, "email already exists")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"CreateMerchantUser",
			"create merchant user request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}
	data := make(map[string]interface{})
	res.Data = data
	res = io.SuccessMessage(data, "merchant User has been added successfully")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"CreateMerchantUser",
		"create merchant user response body",
		"createMerchantUserResponse", "", structs.Map(res), true)
	return
}

func (b *athMerchantService) GetMerchantByID(ctx context.Context, db *gorm.DB, MerchantID int) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"GetMerchantByID",
		"get merchant by ID request parameters",
		"GetMerchantByIDRequest", "", structs.Map(struct{ MerchantID int }{MerchantID: int(MerchantID)}), true)
	Merchant, err := b.DbRepo.GetMerchantByID(ctx, db, MerchantID)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error getting merchant details")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"GetMerchantByID",
			"get merchant by ID request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}
	res = io.SuccessMessage(Merchant, "Got merchant details")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"GetMerchantByID",
		"get merchant by ID response body",
		"GetMerchantByIDResponse", "", structs.Map(res), true)
	return
}

func (b *athMerchantService) VerifyGetMerchant(ctx context.Context, db *gorm.DB, MerchantID int32) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"VerifyGetMerchant",
		"verify merchant by ID request parameters",
		"VerifyGetMerchantRequest", "", structs.Map(struct{ MerchantID int }{MerchantID: int(MerchantID)}), true)
	Merchant, err := b.DbRepo.GetMerchantByUserID(ctx, db, int(MerchantID))
	if err != nil {
		res = io.FailureMessage(res.Error, "Error getting merchant details")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"VerifyGetMerchant",
			"verify merchant by ID request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}
	res = io.SuccessMessage(Merchant, "Got merchant details")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"VerifyGetMerchant",
		"verify merchant by ID response body",
		"VerifyGetMerchantResponse", "", structs.Map(res), true)
	return
}

func (b *athMerchantService) GetMerchantByUserID(ctx context.Context, db *gorm.DB, UserId int) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"GetMerchantByUserID",
		"get merchant by user ID request parameters",
		"GetMerchantByUserIDRequest", "", structs.Map(struct{ UserId int }{UserId: int(UserId)}), true)
	Merchant, err := b.DbRepo.GetMerchantByUserID(ctx, db, UserId)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error getting merchant details")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"GetMerchantByUserID",
			"get merchant by user ID request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}
	res = io.SuccessMessage(Merchant, "Got merchant details")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"GetMerchantByUserID",
		"get merchant by user ID response body",
		"GetMerchantByUserIDResponse", "", structs.Map(res), true)
	return
}

func (b *athMerchantService) AddTeamMember(ctx context.Context, db *gorm.DB, u io.TeamMemberReq) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"AddTeamMember",
		"add team member service request parameters",
		"AddTeamMemberRequest", "", structs.Map(u), true)
	if len(u.AccessData) > 0 {
		// add user account details
		for i := 0; i < len(u.AccessData); i++ {
			memberdata := io.MemberData{
				UserID:      u.UserID,
				CreatedBy:   u.CreatedBy,
				VenueID:     u.AccessData[i],
				AccountType: u.AccountType,
			}
			// add team member to db
			err := b.DbRepo.AddTeamMember(ctx, db, memberdata)
			if err != nil {
				res = io.FailureMessage(res.Error, "Error creating team member")
				res.Error = err
				b.Logger.Log(gologger.Errsev3,
					gologger.InternalServices,
					"AddTeamMember",
					"add team member service request failed",
					gologger.ParseError, "", structs.Map(res), true)
				return
			}
		}
	} else {
		// add admin account details
		memberdata := io.MemberData{
			UserID:      u.UserID,
			CreatedBy:   u.CreatedBy,
			AccountType: u.AccountType,
		}
		// add team member to db
		err := b.DbRepo.AddTeamMember(ctx, db, memberdata)
		if err != nil {
			res = io.FailureMessage(res.Error, "Error creating team member")
			res.Error = err
			b.Logger.Log(gologger.Errsev3,
				gologger.InternalServices,
				"AddTeamMember",
				"add team member service request failed",
				gologger.ParseError, "", structs.Map(res), true)
			return
		}
	}
	fmt.Println("response from service : ", res)
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"AddTeamMember",
		"add team member service response body",
		"AddTeamMember service Response", "", structs.Map(res), true)
	return

}

func (b *athMerchantService) UpdateTeamMemberPrivileges(ctx context.Context, db *gorm.DB, u io.AthVenueUser) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"UpdateTeamMemberPrivileges",
		"update team member privileges request parameters",
		"UpdateTeamMemberPrivilegesRequest", "", structs.Map(u), true)
	// update team member privileges
	err := b.DbRepo.UpdateTeamMemberPrivileges(ctx, db, u)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error updating team member")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"UpdateTeamMemberPrivileges",
			"update team member privileges request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}
	res = io.SuccessMessage(nil, "team member privileges updated")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"UpdateTeamMemberPrivileges",
		"update team member privileges response body",
		"UpdateTeamMemberPrivileges service Response", "", structs.Map(res), true)
	return
}

func (b *athMerchantService) DeleteTeamMember(ctx context.Context, db *gorm.DB, u io.AthVenueUser) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"DeleteTeamMember",
		"delete team member request parameters",
		"DeleteTeamMemberRequest", "", structs.Map(u), true)
	// delete team member
	err := b.DbRepo.DeleteTeamMember(ctx, db, u)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error deleting team member")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"DeleteTeamMember",
			"delete team member request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}
	res = io.SuccessMessage(nil, "team member deleted")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"DeleteTeamMember",
		"delete team member response body",
		"DeleteTeamMember service Response", "", structs.Map(res), true)
	return
}

func (b *athMerchantService) GetTeamData(ctx context.Context, db *gorm.DB, UserId int, orderBy string) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"GetTeamData",
		"get team member request parameters",
		"GetTeamDataRequest", "", structs.Map(struct{ UserId int }{UserId: int(UserId)}), true)
	// get team member
	teamData, err := b.DbRepo.GetTeamData(ctx, db, UserId, orderBy)
	fmt.Println("teamData : ", teamData)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error getting team member data")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"GetTeamData",
			"get team member request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}

	var FinalTeamData []*merchantpb.TeamMemberData
	// get details about each user
	for _, eachUser := range teamData {
		fmt.Println("eachUser : ", eachUser)
		var eachUserData merchantpb.TeamMemberData
		UserData, err := b.userRepo.GetUserByID(ctx, db, eachUser.UserId, false)
		if err != nil {
			res = io.FailureMessage(res.Error, "Error getting user details")
			res.Error = err
			b.Logger.Log(gologger.Errsev3,
				gologger.InternalServices,
				"GetTeamData",
				"get team member request failed",
				gologger.ParseError, "", structs.Map(res), true)
			return
		}

		// get role data for each user
		roleData, err := b.DbRepo.GetRoleByRoleID(ctx, db, eachUser.RoleId)
		if err != nil {
			res = io.FailureMessage(res.Error, "Error getting role details")
			res.Error = err
			b.Logger.Log(gologger.Errsev3,
				gologger.InternalServices,
				"GetTeamData",
				"get team member request failed",
				gologger.ParseError, "", structs.Map(res), true)
			return

		}

		UserData1, err := b.userRepo.GetUserByID(ctx, db, eachUser.CreatedBy, false)
		if err != nil {
			res = io.FailureMessage(res.Error, "Error getting user details")
			res.Error = err
			b.Logger.Log(gologger.Errsev3,
				gologger.InternalServices,
				"GetTeamData",
				"get team member request failed",
				gologger.ParseError, "", structs.Map(res), true)
			return
		}
		eachUserData.CreatedAt = int32(eachUser.CreatedAt)
		eachUserData.CreatedBy = UserData1.FirstName + " " + UserData1.LastName
		eachUserData.Email = UserData.Email
		eachUserData.Phone = UserData.Phone
		eachUserData.AccountId = int32(UserData.UserID)
		eachUserData.FullName = UserData.FirstName + " " + UserData.LastName
		eachUserData.AccountType = roleData.RoleName
		FinalTeamData = append(FinalTeamData, &eachUserData)
	}

	res = io.SuccessMessage(FinalTeamData, "got team member details")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"GetTeamData",
		"get team member response body",
		"GetTeamData service Response", "", structs.Map(res), true)
	return
}
