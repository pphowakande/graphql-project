package otp

import (
	db "api/src/dal"
	"api/src/dal/otp"
	"api/src/dal/user"
	io "api/src/models"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	"stash.bms.bz/bms/gologger.git"

	"github.com/fatih/structs"
)

type AthOtpService interface {
	CreateOTP(ctx context.Context, db, transactionDB *gorm.DB, u io.AthUserOTP) (res io.Response)
	CreateCode(ctx context.Context, db *gorm.DB, u io.AthUserOTP) (res io.Response)
	VerifyOTP(ctx context.Context, masterdb, transactionDB *gorm.DB, db *gorm.DB, u io.OTPVerify) (res io.Response)
	VerifyResetPasswordOTP(ctx context.Context, u io.OTPVerify) (res io.Response)
	EmailSend(ctx context.Context, db, transactionDB *gorm.DB, emailSendReq io.EmailSendReq) (res io.Response)
}

type athOtpService struct {
	Logger   *gologger.Logger
	DbRepo   otp.Repository
	userRepo user.Repository
}

func NewBasicAthOtpService(logger *gologger.Logger, DbRepo otp.Repository, userRepo user.Repository) AthOtpService {
	return &athOtpService{
		Logger:   logger,
		DbRepo:   DbRepo,
		userRepo: userRepo,
	}
}

func (b *athOtpService) EmailSend(ctx context.Context, masterDB, transactionDB *gorm.DB, u io.EmailSendReq) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"EmailSend",
		"EmailSend service request parameters",
		"EmailSend service Request", "", structs.Map(u), true)

	// get user id from user email
	var emailData io.LoginRequest
	emailData.EmailParam = true
	emailData.PhoneParam = false
	emailData.Login = u.Email

	User, err := b.userRepo.GetUserByEmailORPhone(ctx, masterDB, emailData)
	fmt.Println("User : ", User)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error getting user details")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"EmailSend",
			"EmailSend service request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}

	if u.EmailVerifyToken != "" {
		fmt.Println("Inside if loop--------------")
		// save email OTP to database
		var otpRequest io.AthUserOTP
		otpRequest.UserID = User.UserID
		otpRequest.OTPNO = u.EmailVerifyToken
		otpRequest.OTPType = u.OTPTtype
		otpRequest.ExpiredAt = int(time.Now().Unix()) + 600 // otp expires in 10 minutes
		otpRequest.CreatedAt = int(time.Now().Unix())
		otpRequest.Email = u.Email
		fmt.Println("calling CreateOTP db function")
		err = b.DbRepo.CreateOTP(ctx, transactionDB, otpRequest)
		fmt.Println("err : ", err)
		if err != nil {
			res = io.FailureMessage(res.Error, "Error saving email OTP in database")
			res.Error = err
			b.Logger.Log(gologger.Errsev3,
				gologger.InternalServices,
				"EmailSend",
				"EmailSend service request failed",
				gologger.ParseError, "", structs.Map(res), true)
			return
		}
	}

	fmt.Println("lets send san email---------------")
	// send an email
	requestBody, err := json.Marshal(map[string]string{
		"email":   User.Email,
		"ac":      u.Account,
		"subject": u.Subject,
		"body":    u.Body,
		"eticket": u.Eticket,
		"tid":     u.Tid,
	})

	resp, err := http.Post("http://in-cs-app.sit.n3b.bookmyshow.org/send", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		// handle err
		res = io.FailureMessage(nil, "Error Sending email to given email address")
		res.Error = errors.New("Error Sending email to given email address")
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"EmailSend",
			"EmailSend service request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		res = io.SuccessMessage(nil, "Email has been successfully")
		b.Logger.Log(gologger.Info,
			gologger.InternalServices,
			"EmailSend",
			"EmailSend service response body",
			"EmailSend service Response", "", structs.Map(res), true)
		return
	}
	res = io.FailureMessage(nil, "Error Sending email to given email address")
	res.Error = err
	b.Logger.Log(gologger.Errsev3,
		gologger.InternalServices,
		"EmailSend",
		"EmailSend service request failed",
		gologger.ParseError, "", structs.Map(res), true)
	return
}

func (b *athOtpService) CreateOTP(ctx context.Context, masterDB, transactionDB *gorm.DB, u io.AthUserOTP) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"createOTP",
		"create OTP service request parameters",
		"createOTP service Request", "", structs.Map(u), true)
	// first check if valid userid is passed or not

	// get user email from user id
	User, err := b.userRepo.GetUserByID(ctx, masterDB, u.UserID, false)
	fmt.Println("User : ", User)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error getting user details")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"createOTP",
			"create OTP service request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}

	u.Phone = User.Phone
	if User.UserID != 0 {
		fmt.Println("create otp in database")
		fmt.Println("transactionDB : ", transactionDB)
		err := b.DbRepo.CreateOTP(ctx, transactionDB, u)
		if err != nil {
			res = io.FailureMessage(res.Error, "Error saving OTP in database")
			res.Error = err
			b.Logger.Log(gologger.Errsev3,
				gologger.InternalServices,
				"createOTP",
				"create OTP service request failed",
				gologger.ParseError, "", structs.Map(res), true)
			return
		}

		uri, err := url.ParseRequestURI("http://in-cs-app.sit.n3b.bookmyshow.org/sendsms")
		if err != nil {
			res = io.FailureMessage(nil, "Error Sending OTP to given mobile number")
			res.Error = errors.New("Error Sending OTP to given mobile number")
			b.Logger.Log(gologger.Errsev3,
				gologger.InternalServices,
				"createOTP",
				"create OTP service request failed",
				gologger.ParseError, "", structs.Map(res), true)
			return
			// handle err
		}

		//var testMessage string
		testMessage := "OTP to verify your mobile number is " + u.OTPNO

		splittedPhone := strings.Split(u.Phone, "+91")
		phone := ""
		if len(splittedPhone) > 0 {
			phone = splittedPhone[1]
		} else {
			phone = u.Phone
		}

		values := uri.Query()
		values.Set("to", phone)
		values.Set("type", "Confirmation")
		values.Set("message", testMessage)
		values.Set("country", "91")
		values.Set("medium", "test")
		values.Set("refcode", "test")

		uri.RawQuery = values.Encode()

		req, err := http.NewRequest(http.MethodGet, uri.String(), nil)
		if err != nil {
			// handle err
			res = io.FailureMessage(nil, "Error Sending OTP to given mobile number")
			res.Error = errors.New("Error Sending OTP to given mobile number")
			b.Logger.Log(gologger.Errsev3,
				gologger.InternalServices,
				"createOTP",
				"create OTP service request failed",
				gologger.ParseError, "", structs.Map(res), true)
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			// handle err
			res = io.FailureMessage(nil, "Error Sending OTP to given mobile number")
			res.Error = errors.New("Error Sending OTP to given mobile number")
			b.Logger.Log(gologger.Errsev3,
				gologger.InternalServices,
				"createOTP",
				"create OTP service request failed",
				gologger.ParseError, "", structs.Map(res), true)
			return
		}

		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			res = io.FailureMessage(nil, "Error Sending OTP to given mobile number")
			res.Error = errors.New("Error Sending OTP to given mobile number")
			b.Logger.Log(gologger.Errsev3,
				gologger.InternalServices,
				"createOTP",
				"create OTP service request failed",
				gologger.ParseError, "", structs.Map(res), true)
		}

		data := make(map[string]interface{})
		data["verify_token"] = u.OTPNO

		res = io.SuccessMessage(data, "OTP has been sent to given mobile number")
		b.Logger.Log(gologger.Info,
			gologger.InternalServices,
			"createOTP",
			"create OTP service response body",
			"createOTP service Response", "", structs.Map(res), true)
		return
	} else {
		res = io.FailureMessage(res.Error, "invalid user")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"createOTP",
			"create OTP service request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}
}

func (b *athOtpService) CreateCode(ctx context.Context, db *gorm.DB, u io.AthUserOTP) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"CreateCode",
		"create code service request parameters",
		"CreateCode service Request", "", structs.Map(u), true)

	// save OTP code to database TransactionDBConnectionName
	err := b.DbRepo.CreateOTP(ctx, db, u)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error saving OTP in database")
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"CreateCode",
			"create code service request failed",
			gologger.ParseError, "", structs.Map(res), true)
		return
	}
	return
}

func (b *athOtpService) VerifyOTP(ctx context.Context, masterDB, transactionDB *gorm.DB, db *gorm.DB, u io.OTPVerify) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"VerifyOTP",
		"Verify OTP service request parameters",
		"VerifyOTPRequest", "", structs.Map(u), true)

	if u.Email != "" {
		var emailReq io.LoginRequest
		emailReq.Login = u.Email
		emailReq.EmailParam = true
		// if email address is passed, get user id from email address
		userData, err := b.userRepo.GetUserByEmailORPhone(ctx, db, emailReq)
		if err != nil {
			res = io.FailureMessage(err)
			b.Logger.Log(gologger.Errsev3,
				gologger.InternalServices,
				"loginUser",
				"login user request failed",
				gologger.ParseError, "", structs.Map(res), true)
			return
		}
		if userData.UserID != 0 {
			u.AccountId = int32(userData.UserID)
		} else {
			res = io.FailureMessage(err, "No such user exists")
			res.Success = false
			res.Error = err
			b.Logger.Log(gologger.Errsev3,
				gologger.InternalServices,
				"VerifyOTP",
				"Verify OTP service request failed",
				gologger.ValidationError, "", structs.Map(res), true)
			return
		}
	}

	otpData, err := b.DbRepo.VerifyOTP(ctx, transactionDB, u)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error updating OTP isactive flag")
		res.Success = false
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"VerifyOTP",
			"Verify OTP service failed",
			gologger.ValidationError, "", structs.Map(res), true)
		return
	}

	var validateUser io.ValidateUser
	validateUser.UserID = otpData.UserID
	if u.VerificationType == "phone" {
		// update phoneverify flag from user table
		validateUser.PhoneVerify = true
	} else if u.VerificationType == "email" {
		// update emailverify flag from user table
		validateUser.EmailVerify = true
	}

	// update user isactive , emailVerify and phoneVerifys flag to true
	ValidatedUser, err := b.userRepo.ValidateUser(ctx, masterDB, validateUser)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error updating user verification flags")
		res.Success = false
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"VerifyOTP",
			"Verify OTP service request failed",
			gologger.ValidationError, "", structs.Map(res), true)
		return
	}

	res.Data = ValidatedUser
	res = io.SuccessMessage(ValidatedUser, "verified successfully")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"VerifyOTP",
		"Verify OTP service response body",
		"VerifyOTPResponse", "", structs.Map(res), true)
	return

}

func (b *athOtpService) VerifyResetPasswordOTP(ctx context.Context, u io.OTPVerify) (res io.Response) {
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"VerifyResetPasswordOTP",
		"VerifyResetPasswordOTP service request parameters",
		"VerifyOTPRequest", "", structs.Map(u), true)
	tansactionDB, err := b.userRepo.GetDB(ctx, db.TransactionDBConnectionName)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error updating OTP isactive flag")
		res.Success = false
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"VerifyResetPasswordOTP",
			"VerifyResetPasswordOTP service failed",
			gologger.ValidationError, "", structs.Map(res), true)
		return
	}
	_, err = b.DbRepo.VerifyOTP(ctx, tansactionDB, u)
	if err != nil {
		res = io.FailureMessage(res.Error, "Error updating OTP isactive flag")
		res.Success = false
		res.Error = err
		b.Logger.Log(gologger.Errsev3,
			gologger.InternalServices,
			"VerifyResetPasswordOTP",
			"VerifyResetPasswordOTP service failed",
			gologger.ValidationError, "", structs.Map(res), true)
		return
	}
	res = io.SuccessMessage(nil, "verified successfully")
	b.Logger.Log(gologger.Info,
		gologger.InternalServices,
		"VerifyResetPasswordOTP",
		"VerifyResetPasswordOTP service response body",
		"VerifyResetPasswordOTPResponse", "", structs.Map(res), true)
	return
}
