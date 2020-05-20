package customer

import (
	db "api/src/dal"
	"api/src/dal/user"
	io "api/src/models"
	"context"
	"fmt"
	"log"
	"regexp"

	"stash.bms.bz/bms/gologger.git"

	"github.com/fatih/structs"
)

type AthCustomerService interface {
	CreateCustomer(ctx context.Context, u io.AthUser) (res io.Response)
}

type athCustomerService struct {
	Logger *gologger.Logger
	DbRepo user.Repository
	//logger     log.Logger
}

func NewBasicAthCustomerService(logger *gologger.Logger, DbRepo user.Repository) AthCustomerService {
	return &athCustomerService{
		Logger: logger,
		DbRepo: DbRepo,
	}
}

func (b *athCustomerService) CreateCustomer(ctx context.Context, u io.AthUser) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"createCustomer",
		"create customer request parameters",
		"createCustomerRequest", "", structs.Map(u), true)
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if re.MatchString(u.Email) == false {
		log.Printf("Validation error: %v", "incorrect email format")
		res.Error = fmt.Errorf("invalid email")
		res.Message = "invalid email address"
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"createCustomer",
			"create customer request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}
	masterDB, err := b.DbRepo.GetDB(ctx, db.MasterDBConnectionName)
	if err != nil {
		res.Error = fmt.Errorf("invalid DB")
		res.Message = "invalid DB address"
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"createCustomer",
			"create customer request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}
	newUser, err := b.DbRepo.CreateUser(ctx, masterDB, u)
	if err != nil {
		res = io.FailureMessage(res.Error, "email already exists")
		res.Success = false
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"createCustomer",
			"create customer request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}

	data := make(map[string]interface{})
	data["user_id"] = newUser.UserID
	res.Data = data
	res = io.SuccessMessage(data, "Customer has been saved")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"createCustomer",
		"create customer response body",
		"createCustomerResponse", "", structs.Map(res), true)
	return
}
