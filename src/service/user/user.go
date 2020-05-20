package user

import (
	"api/src/dal/user"
	io "api/src/models"
	utile "api/src/utils"
	"context"
	"fmt"
	"regexp"
	"time"
	"unicode"

	"stash.bms.bz/bms/gologger.git"

	"github.com/fatih/structs"
	"github.com/jinzhu/gorm"
)

type AthUserService interface {
	LoginUser(ctx context.Context, db *gorm.DB, u io.LoginRequest) (res io.Response)
	CreateUser(ctx context.Context, db *gorm.DB, u io.AthUser, accType string) (res io.Response)
	EditUser(ctx context.Context, db *gorm.DB, u io.AthUser, userType string) (res io.Response)
	GetUserByID(ctx context.Context, db *gorm.DB, userid int, flag bool) (res io.Response)
	ResetPasswordUser(ctx context.Context, db *gorm.DB, u io.ResetPasswordRequest) (res io.Response)
	GetAccountDetailsByID(ctx context.Context, db *gorm.DB, userid int) (res io.Response)
}

type athUserService struct {
	DbRepo user.Repository
	Logger *gologger.Logger
}

func NewBasicAthUserService(logger *gologger.Logger, DbRepo user.Repository) AthUserService {
	return &athUserService{
		Logger: logger,
		DbRepo: DbRepo,
	}
}

func (b *athUserService) LoginUser(ctx context.Context, db *gorm.DB, u io.LoginRequest) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"loginUser",
		"login user request parameters",
		"loginUserRequest", "", structs.Map(u), true)
	var userData io.AthUser

	// check if field passed is email
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if re.MatchString(u.Login) == false {
		u.EmailParam = false
		// check if field passed is phone
		if len(u.Login) == 0 || len(u.Login) != 13 {
			u.PhoneParam = false
		} else {
			u.PhoneParam = true
		}
	} else {
		u.EmailParam = true
	}

	if u.EmailParam == true || u.PhoneParam == true {
		// check if email/phone exists
		userData, res.Error = b.DbRepo.GetUserByEmailORPhone(ctx, db, u)
		if res.Error != nil {
			res = io.FailureMessage(res.Error)
			b.Logger.Log(gologger.Errsev3,
				gologger.InternalServices,
				"loginUser",
				"login user request failed",
				gologger.ParseError, "", structs.Map(res), true)
			return
		}

		// get account type for requested user
		accountType, err1 := b.DbRepo.GetAccountTypeForUser(ctx, db, userData.UserID)
		if err1 != nil {
			res = io.FailureMessage(err1)
			res.Error = err1
			b.Logger.Log(gologger.Errsev3,
				gologger.InternalServices,
				"loginUser",
				"login user request failed",
				gologger.ParseError, "", structs.Map(res), true)
			return
		}

		// if user is owner , then check if phone verification is completed or not
		if !userData.PhoneVerify {
			res = io.FailureMessage(fmt.Errorf(`You can not login because your contact number is not verified`))
			b.Logger.Log(gologger.Errsev3,
				gologger.InternalServices,
				"loginUser",
				"login user request failed",
				gologger.ValidationError, "", structs.Map(res), true)
			return
		}

		if !userData.EmailVerify {
			res = io.FailureMessage(fmt.Errorf(`You can not login because your email is not verified`))
			b.Logger.Log(gologger.Errsev3,
				gologger.InternalServices,
				"loginUser",
				"login user request failed",
				gologger.ValidationError, "", structs.Map(res), true)
			return
		}

		// check if credentials passed are correct
		if !utile.ComparePasswords(userData.Password, []byte(u.Password)) {
			res = io.FailureMessage(fmt.Errorf(`incorrect credentials`))
			//res.Error = res.Error
			b.Logger.Log(gologger.Errsev3,
				gologger.InternalServices,
				"loginUser",
				"login user request failed",
				gologger.ValidationError, "", structs.Map(res), true)
			return
		}

		//var privilegeData []io.AthVenueUser
		var listOfPrivileges []int32
		if accountType == "user" {
			// if user is any team member, then both email and phone should be verified
			if !userData.EmailVerify {
				res = io.FailureMessage(fmt.Errorf(`You can not login because your email is not verified`))
				b.Logger.Log(gologger.Errsev3,
					gologger.InternalServices,
					"loginUser",
					"login user request failed",
					gologger.ValidationError, "", structs.Map(res), true)
				return
			}
			// only user account has specific privileges
			// If account type is User, then get privileges for requested user
			privilegeData, err1 := b.DbRepo.GetPrivilegesForUser(ctx, db, userData.UserID)
			if err1 != nil {
				res = io.FailureMessage(err1)
				res.Error = err1
				b.Logger.Log(gologger.Errsev3,
					gologger.InternalServices,
					"loginUser",
					"login user request failed",
					gologger.ParseError, "", structs.Map(res), true)
				return
			}

			if len(privilegeData) > 0 {
				for _, value := range privilegeData {
					listOfPrivileges = append(listOfPrivileges, int32(value.VenueId))
				}
			}
		}

		// create login token, update lastlogin details and token to database
		loginData, err1 := b.DbRepo.UpdateLoginData(ctx, db, userData.UserID)
		if err1 != nil {
			res = io.FailureMessage(err1)
			res.Error = err1
			b.Logger.Log(gologger.Errsev3,
				gologger.InternalServices,
				"loginUser",
				"login user request failed",
				gologger.ParseError, "", structs.Map(res), true)
			return
		}
		data := make(map[string]interface{})
		data["AccountType"] = accountType
		data["AccessData"] = listOfPrivileges
		data["Email"] = userData.Email
		data["MerchantId"] = userData.UserID
		data["MerchantFullName"] = userData.FirstName + " " + userData.LastName
		data["Token"] = loginData.TokenValue
		data["EmailVerify"] = loginData.EmailVerify
		data["PhoneVerify"] = loginData.PhoneVerify
		data["BMSVerify"] = loginData.BmsActive
		data["LastLoginAt"] = loginData.LastLoginAt
		data["MerchantName"] = loginData.LastLoginAt

		res = io.SuccessMessage(data)
		b.Logger.Log(gologger.Info,
			gologger.InternalServices,
			"loginUser",
			"login user response body",
			"loginUserResponse", "", structs.Map(res), true)
		return
	}
	res.Error = fmt.Errorf("No valid phone or email passed")
	b.Logger.Log(gologger.Errsev3,
		gologger.InternalServices,
		"loginUser",
		"login user service request failed",
		gologger.ValidationError, "", structs.Map(res), true)
	return
}

func (b *athUserService) verifyPassword(s string) (sevenOrMore, number, upper, special bool) {
	letters := 0
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
			letters++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		case unicode.IsLetter(c) || c == ' ':
			letters++
		default:
			//return false, false, false, false
		}
	}
	//sevenOrMore = letters >= 7
	sevenOrMore = len(s) >= 6
	return
}

func (b *athUserService) CreateUser(ctx context.Context, db *gorm.DB, u io.AthUser, accType string) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"CreateUser",
		"create user service request parameters",
		"CreateUserService", "", structs.Map(u), true)
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	//re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if re.MatchString(u.Email) == false {
		res.Error = fmt.Errorf("invalid email address")
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"CreateUser",
			"create user service request failed",
			gologger.ValidationError, "", structs.Map(res), true)
		return
	}

	if accType != "teamMember" {
		if len(u.Password) == 0 || len(u.Password) < 8 {
			res.Error = fmt.Errorf("password length must be greater or equal to 8")
			b.Logger.Log(gologger.Errsev3,
				gologger.InternalServices,
				"CreateUser",
				"create user service request failed",
				gologger.ValidationError, "", structs.Map(res), true)
			return
		}

		if len(u.Password) > 0 {
			sevenOrMore, number, upper, special := b.verifyPassword(u.Password)
			if !sevenOrMore || !number || !upper || !special {
				res.Error = fmt.Errorf("password must contain atleast a number,uppercase alphabet and special character")
				res = io.FailureMessage(res.Error)
				b.Logger.Log(gologger.Errsev3,
					gologger.InternalServices,
					"CreateUser",
					"create user service request failed",
					gologger.ValidationError, "", structs.Map(res), true)
				return
			}
		}

		if u.Password, res.Error = utile.HashAndSalt([]byte(u.Password)); res.Error != nil {
			res = io.FailureMessage(res.Error)
			b.Logger.Log(gologger.Errsev3,
				gologger.InternalServices,
				"CreateUser",
				"create user service request failed",
				gologger.ParseError, "", structs.Map(res), true)
			return
		}
	}

	if len(u.Phone) == 0 || len(u.Phone) != 13 {
		res.Error = fmt.Errorf("Invaid Phone number")
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"CreateUser",
			"create user service request failed",
			gologger.ValidationError, "", structs.Map(res), true)
		return
	}

	u.CreatedAt = int(time.Now().Unix())

	newUser, err := b.DbRepo.CreateUser(ctx, db, u)
	if err != nil {
		res = io.FailureMessage(res.Error, "email already exists")
		res.Success = false
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"CreateUser",
			"create user service failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}
	res = io.SuccessMessage(newUser, "User has been saved")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"CreateUser",
		"create user service response body",
		"CreateUser response ", "", structs.Map(res), true)
	return
}

func (b *athUserService) EditUser(ctx context.Context, db *gorm.DB, u io.AthUser, userType string) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"editUser",
		"edit user request parameters",
		"editUserRequest", "", structs.Map(u), true)
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if re.MatchString(u.Email) == false {
		res.Error = fmt.Errorf("invalid email")
		res.Message = "invalid email address"
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"editUser",
			"edit user request failed",
			gologger.ValidationError, "", structs.Map(res), true)
		return
	}

	updatedUser, err := b.DbRepo.EditUser(ctx, db, u, userType)
	if err != nil {
		res = io.FailureMessage(res.Error, "error updating user")
		res.Success = false
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"editUser",
			"edit user request failed",
			gologger.ValidationError, "", structs.Map(res), true)
		return
	}

	// get account type for requested user
	accountType, err1 := b.DbRepo.GetAccountTypeForUser(ctx, db, updatedUser.UserID)
	if err1 != nil {
		res = io.FailureMessage(err1)
		res.Error = err1
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"editUser",
			"edit user request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}

	//var privilegeData []io.AthVenueUser
	var listOfPrivileges []int32
	if accountType == "user" {
		// only user account has specific privileges
		// If account type is User, then get privileges for requested user
		privilegeData, err1 := b.DbRepo.GetPrivilegesForUser(ctx, db, updatedUser.UserID)
		if err1 != nil {
			res = io.FailureMessage(err1)
			res.Error = err1
			b.Logger.Log(gologger.Errsev3,
				gologger.InternalServices,
				"editUser",
				"edit user request failed",
				gologger.ParseError, "", structs.Map(res), true)
			return
		}

		if len(privilegeData) > 0 {
			for _, value := range privilegeData {
				listOfPrivileges = append(listOfPrivileges, int32(value.VenueId))
			}
		}
	}

	data := make(map[string]interface{})
	data["AccountType"] = accountType
	data["AccessData"] = listOfPrivileges
	data["Email"] = updatedUser.Email
	data["FirstName"] = updatedUser.FirstName
	data["LastName"] = updatedUser.LastName
	data["Phone"] = updatedUser.Phone
	data["AccountId"] = updatedUser.UserID
	data["AccountType"] = accountType
	data["EmailVerify"] = updatedUser.EmailVerify
	data["PhoneVerify"] = updatedUser.PhoneVerify
	data["UserID"] = updatedUser.UserID

	res = io.SuccessMessage(data, "User has been updated")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"editUser",
		"edit user response body",
		"editUserResponse", "", structs.Map(res), true)
	return
}

func (b *athUserService) GetUserByID(ctx context.Context, db *gorm.DB, userid int, flag bool) (res io.Response) {
	// b.Logger.Log(gologger.Info,
	// 	gologger.InternalServices,
	// 	"getUserByID",
	// 	"get user byID request parameters",
	// 	"getUserByIDRequest", "", structs.Map(userid), true)
	User, err := b.DbRepo.GetUserByID(ctx, db, userid, flag)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error getting user details")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"getUserByID",
			"get user byID request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}
	res = io.SuccessMessage(User, "Got user details")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"getUserByID",
		"get user byID response body",
		"getUserByIDResponse", "", structs.Map(res), true)
	return
}

func (b *athUserService) ResetPasswordUser(ctx context.Context, db *gorm.DB, u io.ResetPasswordRequest) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"resetPasswordUser",
		"reset password user request parameters",
		"resetPasswordUserRequest", "", structs.Map(u), true)
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if re.MatchString(u.Email) == false {
		res.Error = fmt.Errorf("invalid email")
		res.Message = "invalid email address"
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"resetPasswordUser",
			"reset password user request failed",
			gologger.ValidationError, "", structs.Map(res), true)
		return
	}
	if len(u.Password) == 0 || len(u.Password) < 8 {
		res.Error = fmt.Errorf("password length must be greater or equal to 8")
		res.Message = "invalid password"
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"resetPasswordUser",
			"reset password user request failed",
			gologger.ValidationError, "", structs.Map(res), true)
		return

	}
	if len(u.Password) > 0 {
		sevenOrMore, number, upper, special := b.verifyPassword(u.Password)
		if !sevenOrMore || !number || !upper || !special {
			res.Error = fmt.Errorf("password must contain atleast a number,uppercase alphabet and special character")
			res.Message = "invalid password"
			res = io.FailureMessage(res.Error)
			b.Logger.Log(gologger.Errsev3,
				gologger.InternalServices,
				"resetPasswordUser",
				"reset password user request failed",
				gologger.ValidationError, "", structs.Map(res), true)
			return
		}
	}
	if u.Password, res.Error = utile.HashAndSalt([]byte(u.Password)); res.Error != nil {
		res.Message = "invalid password"
		res = io.FailureMessage(res.Error)
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"resetPasswordUser",
			"reset password user request failed",
			gologger.ValidationError, "", structs.Map(res), true)
		return
	}
	fmt.Println("Calling ResetPasswordUser DB function--------------")
	err := b.DbRepo.ResetPasswordUser(ctx, db, io.AthUser{Email: u.Email, Password: u.Password})
	if err != nil {
		res = io.FailureMessage(res.Error, err.Error())
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"resetPasswordUser",
			"reset password user request failed",
			gologger.ValidationError, "", structs.Map(res), true)
		return
	}

	res = io.SuccessMessage(u, "Password has been reset successfully")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"resetPasswordUser",
		"reset password user response body",
		"resetPasswordUserResponse", "", structs.Map(res), true)
	return
}

func (b *athUserService) GetAccountDetailsByID(ctx context.Context, db *gorm.DB, accountId int) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"GetAccountDetailsByID",
		"Get account details by ID request parameters",
		"GetAccountDetailsByIDRequest", "", structs.Map(struct{ AccountId int }{AccountId: int(accountId)}), true)

	// get account type for requested user
	accountType, err1 := b.DbRepo.GetAccountTypeForUser(ctx, db, accountId)
	if err1 != nil {
		res = io.FailureMessage(err1)
		res.Error = err1
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"GetAccountDetailsByID",
			"Get account details by ID request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return

	}

	//var privilegeData []io.AthVenueUser
	var listOfPrivileges []int32
	if accountType == "user" {
		// only user account has specific privileges
		// If account type is User, then get privileges for requested user
		privilegeData, err1 := b.DbRepo.GetPrivilegesForUser(ctx, db, accountId)
		fmt.Println("privilegeData : ", privilegeData)
		fmt.Println("len(privilegeData) : ", len(privilegeData))
		if err1 != nil {
			res = io.FailureMessage(err1)
			res.Error = err1
			b.Logger.Log(gologger.Errsev3,
				gologger.InternalServices,
				"GetAccountDetailsByID",
				"Get account details by ID request failed",
				gologger.ParseError, "", structs.Map(res), true)
			return
		}

		if len(privilegeData) > 0 {
			for _, value := range privilegeData {
				listOfPrivileges = append(listOfPrivileges, int32(value.VenueId))
			}
		}
	}

	data := make(map[string]interface{})
	data["AccountType"] = accountType
	data["AccessData"] = listOfPrivileges

	res = io.SuccessMessage(data)
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"GetAccountDetailsByID",
		"Get account details by ID response body",
		"GetAccountDetailsByIDResponse", "", structs.Map(res), true)
	return

}
