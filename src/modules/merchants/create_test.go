package handler

import (
	"errors"

	"github.com/jinzhu/gorm"

	//io "api/src/models"
	"context"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	merchantpb "stash.bms.bz/turf/generic-proto-files.git/merchant/v1"

	"stash.bms.bz/bms/gologger.git"

	merchandal "api/src/dal/merchant"
	io "api/src/models"
	mockmerchant "api/src/service/merchant/mock"
	mockotp "api/src/service/otp/mock"
	mockupload "api/src/service/upload/mock"
	mockuser "api/src/service/user/mock"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAthMerchantHandler_SignupMerchant(t *testing.T) {
	ctx := context.Background()
	logger := gologger.New("api", true)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sql, mockSQl, err := sqlmock.New()
	assert.NoError(t, err)
	db1, err := gorm.Open("mysql", sql)
	dbConnections := map[string]*gorm.DB{
		MasterDBConnectionName:      db1,
		TransactionDBConnectionName: db1,
	}
	mockDBRepo := merchandal.NewMerchantRepo(logger, dbConnections)
	mockMerchantService := mockmerchant.NewMockAthMerchantService(ctrl)
	mockUserService := mockuser.NewMockAthUserService(ctrl)
	mockOTPService := mockotp.NewMockAthOtpService(ctrl)
	handler := athMerchantHandler{
		logger:             logger,
		DbRepo:             mockDBRepo,
		athMerchantService: mockMerchantService,
		athUserService:     mockUserService,
		athOTPService:      mockOTPService,
	}

	type args struct {
		ctx context.Context
		req *merchantpb.SignupRequest
	}
	tests := []struct {
		name        string
		args        args
		mock        func()
		want        *merchantpb.SignupReply
		wantErr     bool
		expectedErr error
	}{
		{
			name: "when invalid business name is passed",
			args: args{
				ctx: ctx,
				req: &merchantpb.SignupRequest{},
			},
			mock:        func() {},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in SignupMerchant :%v", "Request parameter is missing or blank"),
		},
		{
			name: "when invalid phone is passed",
			args: args{
				ctx: ctx,
				req: &merchantpb.SignupRequest{},
			},
			mock:        func() {},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in SignupMerchant :%v", "Request parameter is missing or blank"),
		},
		{
			name: "when invalid email is passed",
			args: args{
				ctx: ctx,
				req: &merchantpb.SignupRequest{},
			},
			mock: func() {
				//mockMerchantService.
				//	EXPECT().
				//	CreateMerchant(ctx, gomock.Any(),gomock.Any()).
				//	Return(io.Response{Success: true})
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in SignupMerchant :%v", "Request parameter is missing or blank"),
		},
		{
			name: "when invalid password is passed",
			args: args{
				ctx: ctx,
				req: &merchantpb.SignupRequest{},
			},
			mock:        func() {},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in SignupMerchant :%v", "Request parameter is missing or blank"),
		},
		{
			name: "when invalid merchant full name is passed",
			args: args{
				ctx: ctx,
				req: &merchantpb.SignupRequest{},
			},
			mock:        func() {},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in SignupMerchant :%v", "Request parameter is missing or blank"),
		},
		{
			name: "when invalid user source is passed",
			args: args{
				ctx: ctx,
				req: &merchantpb.SignupRequest{},
			},
			mock:        func() {},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in SignupMerchant :%v", "Request parameter is missing or blank"),
		},
		{
			name: "when user service returns error while creating",
			args: args{
				ctx: ctx,
				req: &merchantpb.SignupRequest{
					Email:            "test@gmail.com",
					Password:         "test",
					BusinessName:     "test",
					MerchantFullName: "Test User",
					Phone:            "80000",
					UserSource:       "test",
				},
			},
			mock: func() {
				mockSQl.ExpectBegin()
				mockSQl.ExpectBegin()
				mockUserService.EXPECT().
					CreateUser(ctx, gomock.Any(), gomock.Any(), "merchant").
					Return(io.Response{Success: false, Error: errors.New("failed")})
					//mockMerchantService.
					//	EXPECT().
					//	CreateMerchant(ctx, gomock.Any(),gomock.Any()).
					//	Return(io.Response{Success: true, Data: map[string]interface{}{"Merchant_id": 1}})
					//mockMerchantService.
					//	EXPECT().
					//	CreateMerchantUser(ctx, gomock.Any(),gomock.Any()).
					//	Return(io.Response{Success: true})
				mockSQl.ExpectRollback()
				mockSQl.ExpectRollback()
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in SignupMerchant :%v", "failed"),
		}, {
			name: "when user service returns error while creating merchant",
			args: args{
				ctx: ctx,
				req: &merchantpb.SignupRequest{
					Email:            "test@gmail.com",
					Password:         "test",
					BusinessName:     "test",
					MerchantFullName: "Test User",
					Phone:            "80000",
					UserSource:       "test",
				},
			},
			mock: func() {
				mockSQl.ExpectBegin()
				mockSQl.ExpectBegin()
				mockUserService.EXPECT().
					CreateUser(ctx, gomock.Any(), gomock.Any(), "merchant").
					Return(io.Response{Success: true, Data: io.AthUser{
						UserID:              1,
						Email:               "test",
						Password:            "test",
						FirstName:           "test",
						LastName:            "test",
						Phone:               "test",
						Gender:              "test",
						DOB:                 "test",
						ProfileImage:        "test",
						LastPasswordResetAt: 0,
						UserSource:          "test",
						LastLoginAt:         0,
						LastLoginIP:         "test",
						IsActive:            true,
						EmailVerify:         true,
						PhoneVerify:         true,
						BmsActive:           true,
						TokenValue:          "test",
						TokenHash:           "test",
						TokenExpireAt:       0,
					}})
				mockMerchantService.
					EXPECT().
					CreateMerchant(ctx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: false, Error: errors.New("failed")})
				mockSQl.ExpectRollback()
				mockSQl.ExpectRollback()
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in SignupMerchant :%v", "failed"),
		}, {
			name: "when user service returns error while creating merchant user",
			args: args{
				ctx: ctx,
				req: &merchantpb.SignupRequest{
					Email:            "test@gmail.com",
					Password:         "test",
					BusinessName:     "test",
					MerchantFullName: "Test User",
					Phone:            "80000",
					UserSource:       "test",
				},
			},
			mock: func() {
				mockSQl.ExpectBegin()
				mockSQl.ExpectBegin()
				mockUserService.EXPECT().
					CreateUser(ctx, gomock.Any(), gomock.Any(), "merchant").
					Return(io.Response{Success: true, Data: io.AthUser{
						UserID:              1,
						Email:               "test",
						Password:            "test",
						FirstName:           "test",
						LastName:            "test",
						Phone:               "test",
						Gender:              "test",
						DOB:                 "test",
						ProfileImage:        "test",
						LastPasswordResetAt: 0,
						UserSource:          "test",
						LastLoginAt:         0,
						LastLoginIP:         "test",
						IsActive:            true,
						EmailVerify:         true,
						PhoneVerify:         true,
						BmsActive:           true,
						TokenValue:          "test",
						TokenHash:           "test",
						TokenExpireAt:       0,
					}})
				mockMerchantService.
					EXPECT().
					CreateMerchant(ctx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true, Data: map[string]interface{}{"Merchant_id": 1}})
				mockMerchantService.
					EXPECT().
					CreateMerchantUser(ctx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: false, Error: errors.New("failed")})
				mockSQl.ExpectRollback()
				mockSQl.ExpectRollback()
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in SignupMerchant :"),
		},
		{
			name: "when valid user with full name is passed",
			args: args{
				ctx: ctx,
				req: &merchantpb.SignupRequest{
					Email:            "test@gmail.com",
					Password:         "test",
					BusinessName:     "test",
					MerchantFullName: "Test User",
					Phone:            "80000",
					UserSource:       "test",
				},
			},
			mock: func() {
				mockSQl.ExpectBegin()
				mockSQl.ExpectBegin()
				mockUserService.EXPECT().
					CreateUser(ctx, gomock.Any(), gomock.Any(), "merchant").
					Return(io.Response{Success: true, Data: io.AthUser{
						UserID:              1,
						Email:               "test",
						Password:            "test",
						FirstName:           "test",
						LastName:            "test",
						Phone:               "test",
						Gender:              "test",
						DOB:                 "test",
						ProfileImage:        "test",
						LastPasswordResetAt: 0,
						UserSource:          "test",
						LastLoginAt:         0,
						LastLoginIP:         "test",
						IsActive:            true,
						EmailVerify:         true,
						PhoneVerify:         true,
						BmsActive:           true,
						TokenValue:          "test",
						TokenHash:           "test",
						TokenExpireAt:       0,
					}})
				mockMerchantService.
					EXPECT().
					CreateMerchant(ctx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true, Data: map[string]interface{}{"Merchant_id": 1}})
				mockMerchantService.
					EXPECT().
					CreateMerchantUser(ctx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true})
				mockOTPService.EXPECT().
					CreateOTP(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true})
				mockSQl.ExpectCommit()
				mockSQl.ExpectCommit()
			},
			want: &merchantpb.SignupReply{
				AccountId: int32(1),
			},
			wantErr: false,
		},
		{
			name: "when valid user with only first name is passed",
			args: args{
				ctx: ctx,
				req: &merchantpb.SignupRequest{
					Email:            "test@gmail.com",
					Password:         "test",
					BusinessName:     "test",
					MerchantFullName: "Test",
					Phone:            "80000",
					UserSource:       "test",
				},
			},
			mock: func() {
				mockSQl.ExpectBegin()
				mockSQl.ExpectBegin()
				mockUserService.EXPECT().
					CreateUser(ctx, gomock.Any(), gomock.Any(), "merchant").
					Return(io.Response{Success: true, Data: io.AthUser{
						UserID:              1,
						Email:               "test",
						Password:            "test",
						FirstName:           "test",
						LastName:            "test",
						Phone:               "test",
						Gender:              "test",
						DOB:                 "test",
						ProfileImage:        "test",
						LastPasswordResetAt: 0,
						UserSource:          "test",
						LastLoginAt:         0,
						LastLoginIP:         "test",
						IsActive:            true,
						EmailVerify:         true,
						PhoneVerify:         true,
						BmsActive:           true,
						TokenValue:          "test",
						TokenHash:           "test",
						TokenExpireAt:       0,
					}})
				mockMerchantService.
					EXPECT().
					CreateMerchant(ctx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true, Data: map[string]interface{}{"Merchant_id": 1}})
				mockMerchantService.
					EXPECT().
					CreateMerchantUser(ctx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true})
				mockOTPService.EXPECT().
					CreateOTP(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true})
				mockSQl.ExpectCommit()
				mockSQl.ExpectCommit()
			},
			want: &merchantpb.SignupReply{
				AccountId: int32(1),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := handler.SignupMerchant(tt.args.ctx, tt.args.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
				assert.Equal(t, tt.want, got)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAthMerchantHandler_PhoneVerifyMerchant(t *testing.T) {
	ctx := context.Background()
	logger := gologger.New("api", true)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sql, mockSQl, err := sqlmock.New()
	assert.NoError(t, err)
	db1, err := gorm.Open("mysql", sql)
	dbConnections := map[string]*gorm.DB{
		MasterDBConnectionName:      db1,
		TransactionDBConnectionName: db1,
	}
	mockDBRepo := merchandal.NewMerchantRepo(logger, dbConnections)
	mockMerchantService := mockmerchant.NewMockAthMerchantService(ctrl)
	mockUserService := mockuser.NewMockAthUserService(ctrl)
	mockOTPService := mockotp.NewMockAthOtpService(ctrl)
	handler := athMerchantHandler{
		logger:             logger,
		DbRepo:             mockDBRepo,
		athMerchantService: mockMerchantService,
		athUserService:     mockUserService,
		athOTPService:      mockOTPService,
	}

	type args struct {
		ctx context.Context
		req *merchantpb.PhoneVerifyRequest
	}
	tests := []struct {
		name        string
		args        args
		mock        func()
		want        *merchantpb.GenericReply
		wantErr     bool
		expectedErr error
	}{
		{
			name: "when invalid verification code request is passed",
			args: args{
				ctx: ctx,
				req: &merchantpb.PhoneVerifyRequest{
					VerificationCode: "",
				},
			},
			mock:        func() {},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in PhoneVerifyMerchant :%v", "Request parameter is missing or blank"),
		},
		{
			name: "when invalid verification type in request is passed",
			args: args{
				ctx: ctx,
				req: &merchantpb.PhoneVerifyRequest{
					VerificationType: "",
				},
			},
			mock:        func() {},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in PhoneVerifyMerchant :%v", "Request parameter is missing or blank"),
		},
		{
			name: "when verify OTP returns error",
			args: args{
				ctx: ctx,
				req: &merchantpb.PhoneVerifyRequest{
					VerificationType: "SMS",
					VerificationCode: "1234",
					AccountId:        1,
				},
			},
			mock: func() {
				mockSQl.ExpectBegin()
				mockSQl.ExpectBegin()
				mockOTPService.
					EXPECT().
					VerifyOTP(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: false, Error: errors.New("failed")})
				mockSQl.ExpectRollback()
				mockSQl.ExpectRollback()
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "failed"),
		}, {
			name: "when verify OTP returns success, but merchant email is verified",
			args: args{
				ctx: ctx,
				req: &merchantpb.PhoneVerifyRequest{
					VerificationType: "SMS",
					VerificationCode: "1234",
					AccountId:        1,
				},
			},
			mock: func() {
				mockSQl.ExpectBegin()
				mockSQl.ExpectBegin()
				mockOTPService.
					EXPECT().
					VerifyOTP(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true, Data: io.AthUser{
						UserID:              1,
						Email:               "test",
						Password:            "test",
						FirstName:           "test",
						LastName:            "test",
						Phone:               "test",
						Gender:              "test",
						DOB:                 "test",
						ProfileImage:        "test",
						LastPasswordResetAt: 0,
						UserSource:          "test",
						LastLoginAt:         0,
						LastLoginIP:         "test",
						IsActive:            false,
						EmailVerify:         true,
						PhoneVerify:         false,
						BmsActive:           false,
						TokenValue:          "test",
						TokenHash:           "test",
						TokenExpireAt:       0,
					}})
				mockSQl.ExpectCommit()
				mockSQl.ExpectCommit()
			},
			want:    &merchantpb.GenericReply{Status: true},
			wantErr: false,
		}, {
			name: "when verify OTP returns success, but merchant email is not verified, and email send succeed",
			args: args{
				ctx: ctx,
				req: &merchantpb.PhoneVerifyRequest{
					VerificationType: "SMS",
					VerificationCode: "1234",
					AccountId:        1,
				},
			},
			mock: func() {
				mockSQl.ExpectBegin()
				mockSQl.ExpectBegin()
				mockOTPService.
					EXPECT().
					VerifyOTP(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true, Data: io.AthUser{
						UserID:              1,
						Email:               "test",
						Password:            "test",
						FirstName:           "test",
						LastName:            "test",
						Phone:               "test",
						Gender:              "test",
						DOB:                 "test",
						ProfileImage:        "test",
						LastPasswordResetAt: 0,
						UserSource:          "test",
						LastLoginAt:         0,
						LastLoginIP:         "test",
						IsActive:            false,
						EmailVerify:         false,
						PhoneVerify:         false,
						BmsActive:           false,
						TokenValue:          "test",
						TokenHash:           "test",
						TokenExpireAt:       0,
					}})
				mockSQl.ExpectCommit()
				mockSQl.ExpectCommit()
			},
			want:    &merchantpb.GenericReply{Status: true},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := handler.PhoneVerifyMerchant(tt.args.ctx, tt.args.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
				assert.Equal(t, tt.want, got)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAthMerchantHandler_EmailVerifyMerchant(t *testing.T) {
	ctx := context.Background()
	logger := gologger.New("api", true)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sql, mockSQl, err := sqlmock.New()
	assert.NoError(t, err)
	db1, err := gorm.Open("mysql", sql)
	dbConnections := map[string]*gorm.DB{
		MasterDBConnectionName:      db1,
		TransactionDBConnectionName: db1,
	}
	mockDBRepo := merchandal.NewMerchantRepo(logger, dbConnections)
	mockMerchantService := mockmerchant.NewMockAthMerchantService(ctrl)
	mockUserService := mockuser.NewMockAthUserService(ctrl)
	mockOTPService := mockotp.NewMockAthOtpService(ctrl)
	handler := athMerchantHandler{
		logger:             logger,
		DbRepo:             mockDBRepo,
		athMerchantService: mockMerchantService,
		athUserService:     mockUserService,
		athOTPService:      mockOTPService,
	}

	type args struct {
		ctx context.Context
		req *merchantpb.EmailVerifyRequest
	}
	tests := []struct {
		name        string
		args        args
		mock        func()
		want        *merchantpb.GenericReply
		wantErr     bool
		expectedErr error
	}{
		{
			name: "when invalid verification code request is passed",
			args: args{
				ctx: ctx,
				req: &merchantpb.EmailVerifyRequest{
					VerificationCode: "",
				},
			},
			mock:        func() {},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in EmailVerifyMerchant :%v", "Request parameter is missing or blank"),
		},
		{
			name: "when invalid verification type in request is passed",
			args: args{
				ctx: ctx,
				req: &merchantpb.EmailVerifyRequest{
					VerificationCode: "SMS",
					VerificationType: "",
				},
			},
			mock:        func() {},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in EmailVerifyMerchant :%v", "Request parameter is missing or blank"),
		},
		{
			name: "when verify OTP returns error",
			args: args{
				ctx: ctx,
				req: &merchantpb.EmailVerifyRequest{
					VerificationType: "SMS",
					VerificationCode: "1234",
					Email:            "test@test.com",
				},
			},
			mock: func() {
				mockSQl.ExpectBegin()
				mockSQl.ExpectBegin()
				mockOTPService.
					EXPECT().
					VerifyOTP(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: false, Error: errors.New("failed")})
				mockSQl.ExpectRollback()
				mockSQl.ExpectRollback()
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "failed"),
		}, {
			name: "when verify OTP returns success",
			args: args{
				ctx: ctx,
				req: &merchantpb.EmailVerifyRequest{
					VerificationType: "SMS",
					VerificationCode: "1234",
					Email:            "test@test.com",
				},
			},
			mock: func() {
				mockSQl.ExpectBegin()
				mockSQl.ExpectBegin()
				mockOTPService.
					EXPECT().
					VerifyOTP(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true, Data: io.AthUser{
						UserID:              1,
						Email:               "test",
						Password:            "test",
						FirstName:           "test",
						LastName:            "test",
						Phone:               "test",
						Gender:              "test",
						DOB:                 "test",
						ProfileImage:        "test",
						LastPasswordResetAt: 0,
						UserSource:          "test",
						LastLoginAt:         0,
						LastLoginIP:         "test",
						IsActive:            false,
						EmailVerify:         false,
						PhoneVerify:         true,
						BmsActive:           false,
						TokenValue:          "test",
						TokenHash:           "test",
						TokenExpireAt:       0,
					}})
				mockSQl.ExpectCommit()
				mockSQl.ExpectCommit()
			},
			want:    &merchantpb.GenericReply{Status: true},
			wantErr: false,
		}, {
			name: "when verify OTP returns success, but phone is not verified create otp fails",
			args: args{
				ctx: ctx,
				req: &merchantpb.EmailVerifyRequest{
					VerificationType: "SMS",
					VerificationCode: "1234",
					Email:            "test@test.com",
				},
			},
			mock: func() {
				mockSQl.ExpectBegin()
				mockSQl.ExpectBegin()
				mockOTPService.
					EXPECT().
					VerifyOTP(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true, Data: io.AthUser{
						UserID:              1,
						Email:               "test",
						Password:            "test",
						FirstName:           "test",
						LastName:            "test",
						Phone:               "test",
						Gender:              "test",
						DOB:                 "test",
						ProfileImage:        "test",
						LastPasswordResetAt: 0,
						UserSource:          "test",
						LastLoginAt:         0,
						LastLoginIP:         "test",
						IsActive:            false,
						EmailVerify:         false,
						PhoneVerify:         false,
						BmsActive:           false,
						TokenValue:          "test",
						TokenHash:           "test",
						TokenExpireAt:       0,
					}})
				mockOTPService.
					EXPECT().
					CreateOTP(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: false, Error: errors.New("failed")})
				mockSQl.ExpectRollback()
				mockSQl.ExpectRollback()
			},
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in EmailVerifyMerchant :%v", "failed"),
		}, {
			name: "when verify OTP, and create OTP returns success",
			args: args{
				ctx: ctx,
				req: &merchantpb.EmailVerifyRequest{
					VerificationType: "SMS",
					VerificationCode: "1234",
					Email:            "test@test.com",
				},
			},
			mock: func() {
				mockSQl.ExpectBegin()
				mockSQl.ExpectBegin()
				mockOTPService.
					EXPECT().
					VerifyOTP(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true, Data: io.AthUser{
						UserID:              1,
						Email:               "test",
						Password:            "test",
						FirstName:           "test",
						LastName:            "test",
						Phone:               "test",
						Gender:              "test",
						DOB:                 "test",
						ProfileImage:        "test",
						LastPasswordResetAt: 0,
						UserSource:          "test",
						LastLoginAt:         0,
						LastLoginIP:         "test",
						IsActive:            false,
						EmailVerify:         false,
						PhoneVerify:         false,
						BmsActive:           false,
						TokenValue:          "test",
						TokenHash:           "test",
						TokenExpireAt:       0,
					}})
				mockOTPService.
					EXPECT().
					CreateOTP(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true})
				mockSQl.ExpectCommit()
				mockSQl.ExpectCommit()
			},
			want:    &merchantpb.GenericReply{Status: true},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := handler.EmailVerifyMerchant(tt.args.ctx, tt.args.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
				assert.Equal(t, tt.want, got)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAthMerchantHandler_ResendCode(t *testing.T) {
	ctx := context.Background()
	logger := gologger.New("api", true)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sql, mockSQl, err := sqlmock.New()
	assert.NoError(t, err)
	db1, err := gorm.Open("mysql", sql)
	dbConnections := map[string]*gorm.DB{
		MasterDBConnectionName:      db1,
		TransactionDBConnectionName: db1,
	}
	mockDBRepo := merchandal.NewMerchantRepo(logger, dbConnections)
	mockMerchantService := mockmerchant.NewMockAthMerchantService(ctrl)
	mockUserService := mockuser.NewMockAthUserService(ctrl)
	mockOTPService := mockotp.NewMockAthOtpService(ctrl)
	handler := athMerchantHandler{
		logger:             logger,
		DbRepo:             mockDBRepo,
		athMerchantService: mockMerchantService,
		athUserService:     mockUserService,
		athOTPService:      mockOTPService,
	}

	type args struct {
		ctx context.Context
		req *merchantpb.ResendCodeRequest
	}
	tests := []struct {
		name        string
		args        args
		mock        func()
		want        *merchantpb.GenericReply
		wantErr     bool
		expectedErr error
	}{
		{
			name: "when invalid verification code request is passed",
			args: args{
				ctx: ctx,
				req: &merchantpb.ResendCodeRequest{
					AccountId: 0,
				},
			},
			mock:        func() {},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in resendCode :%v", "Request parameter is missing or blank"),
		}, {
			name: "when create verify get merchant returns error",
			args: args{
				ctx: ctx,
				req: &merchantpb.ResendCodeRequest{
					AccountId: 1,
				},
			},
			mock: func() {
				mockSQl.ExpectBegin()
				mockSQl.ExpectBegin()
				mockMerchantService.
					EXPECT().
					VerifyGetMerchant(ctx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: false, Error: errors.New("failed")})
				mockSQl.ExpectRollback()
				mockSQl.ExpectRollback()
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in user resendCode :%v", "failed"),
		},
		{
			name: "when create OTP returns error",
			args: args{
				ctx: ctx,
				req: &merchantpb.ResendCodeRequest{
					AccountId: 1,
				},
			},
			mock: func() {
				mockSQl.ExpectBegin()
				mockSQl.ExpectBegin()
				mockMerchantService.
					EXPECT().
					VerifyGetMerchant(ctx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true, Data: io.AthMerchant{
						MerchantID:   1,
						MerchantName: "test",
						Address:      "test",
						Phone:        "test",
						Email:        "test",
						GstNoFile:    "test",
						PanNoFile:    "test",
						BankAccFile:  "test",
						AddressFile:  "test",
						IsActive:     true,
						DoLater:      false,
					}})
				mockOTPService.
					EXPECT().
					CreateOTP(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: false, Error: errors.New("failed")})
				mockSQl.ExpectRollback()
				mockSQl.ExpectRollback()
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "failed"),
		}, {
			name: "when create OTP returns success",
			args: args{
				ctx: ctx,
				req: &merchantpb.ResendCodeRequest{
					AccountId: 1,
				},
			},
			mock: func() {
				mockSQl.ExpectBegin()
				mockSQl.ExpectBegin()
				mockMerchantService.
					EXPECT().
					VerifyGetMerchant(ctx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true, Data: io.AthMerchant{
						MerchantID:   1,
						MerchantName: "test",
						Address:      "test",
						Phone:        "test",
						Email:        "test",
						GstNoFile:    "test",
						PanNoFile:    "test",
						BankAccFile:  "test",
						AddressFile:  "test",
						IsActive:     true,
						DoLater:      false,
					}})
				mockOTPService.
					EXPECT().
					CreateOTP(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true})
				mockSQl.ExpectCommit()
				mockSQl.ExpectCommit()
			},
			want:    &merchantpb.GenericReply{Status: true},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := handler.ResendCode(tt.args.ctx, tt.args.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
				assert.Equal(t, tt.want, got)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAthMerchantHandler_UploadDoc(t *testing.T) {
	ctx := context.Background()
	logger := gologger.New("api", true)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sql, mockSQl, err := sqlmock.New()
	assert.NoError(t, err)
	db1, err := gorm.Open("mysql", sql)
	dbConnections := map[string]*gorm.DB{
		MasterDBConnectionName:      db1,
		TransactionDBConnectionName: db1,
	}
	mockDBRepo := merchandal.NewMerchantRepo(logger, dbConnections)
	mockMerchantService := mockmerchant.NewMockAthMerchantService(ctrl)
	mockUserService := mockuser.NewMockAthUserService(ctrl)
	mockOTPService := mockotp.NewMockAthOtpService(ctrl)
	mockUploadService := mockupload.NewMockAthUploadService(ctrl)

	handler := athMerchantHandler{
		logger:             logger,
		DbRepo:             mockDBRepo,
		athMerchantService: mockMerchantService,
		athUserService:     mockUserService,
		athOTPService:      mockOTPService,
		athUploadService:   mockUploadService,
	}

	type args struct {
		ctx context.Context
		req *merchantpb.UploadDocRequest
	}
	tests := []struct {
		name        string
		args        args
		mock        func()
		want        *merchantpb.GenericReply
		wantErr     bool
		expectedErr error
	}{
		{
			name: "when invalid accountID is passed",
			args: args{
				ctx: ctx,
				req: &merchantpb.UploadDocRequest{
					KycData: &merchantpb.KycDataEdit{
						GstNoFile:   "test",
						PanNoFile:   "test",
						BankAccFile: "test",
						AddressFile: "test",
						DoLater:     false,
					},
				},
			},
			mock:        func() {},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in UploadDoc :%v", "Request parameter is missing or blank"),
		}, {
			name: "when invalid kyc data is passed",
			args: args{
				ctx: ctx,
				req: &merchantpb.UploadDocRequest{
					AccountId: 1,
				},
			},
			mock:        func() {},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in UploadDoc :%v", "Request parameter is missing or blank"),
		},
		{
			name: "when upload doc fails",
			args: args{
				ctx: ctx,
				req: &merchantpb.UploadDocRequest{
					AccountId: 1,
					KycData: &merchantpb.KycDataEdit{
						GstNoFile:   "test",
						PanNoFile:   "test",
						BankAccFile: "test",
						AddressFile: "test",
						DoLater:     false,
					},
				},
			},
			mock: func() {
				mockUploadService.
					EXPECT().
					ValidateMerchantDocs(ctx, gomock.Any()).
					Return(io.Response{Success: true, Data: map[string]interface{}{}})
				mockUploadService.
					EXPECT().
					UploadToS3(ctx, gomock.Any()).
					Return(io.Response{Success: false, Error: errors.New("failed")})
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in UploadDoc :%v", "failed"),
		}, {
			name: "when edit merchant service returns error",
			args: args{
				ctx: ctx,
				req: &merchantpb.UploadDocRequest{
					AccountId: 1,
					KycData: &merchantpb.KycDataEdit{
						GstNoFile:   "test",
						PanNoFile:   "test",
						BankAccFile: "test",
						AddressFile: "test",
						DoLater:     false,
					},
				},
			},
			mock: func() {
				mockSQl.ExpectBegin()
				mockSQl.ExpectBegin()
				mockUploadService.
					EXPECT().
					ValidateMerchantDocs(ctx, gomock.Any()).
					Return(io.Response{Success: true, Data: map[string]interface{}{}})
				mockUploadService.
					EXPECT().
					UploadToS3(ctx, gomock.Any()).
					Return(io.Response{Success: true, Data: map[string]interface{}{
						"PanNoFile":   "test",
						"AddressFile": "test",
						"BankAccFile": "test",
						"GstNoFile":   "test",
					}})
				mockMerchantService.
					EXPECT().
					EditMerchant(ctx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: false, Error: errors.New("failed")})
				mockSQl.ExpectRollback()
				mockSQl.ExpectRollback()
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in UploadDoc :%v", "failed"),
		}, {
			name: "when otp service returns error",
			args: args{
				ctx: ctx,
				req: &merchantpb.UploadDocRequest{
					AccountId: 1,
					KycData: &merchantpb.KycDataEdit{
						GstNoFile:   "test",
						PanNoFile:   "test",
						BankAccFile: "test",
						AddressFile: "test",
						DoLater:     false,
					},
				},
			},
			mock: func() {
				mockSQl.ExpectBegin()
				mockSQl.ExpectBegin()
				mockUploadService.
					EXPECT().
					ValidateMerchantDocs(ctx, gomock.Any()).
					Return(io.Response{Success: true, Data: map[string]interface{}{}})
				mockUploadService.
					EXPECT().
					UploadToS3(ctx, gomock.Any()).
					Return(io.Response{Success: true, Data: map[string]interface{}{
						"PanNoFile":   "test",
						"AddressFile": "test",
						"BankAccFile": "test",
						"GstNoFile":   "test",
					}})
				mockMerchantService.
					EXPECT().
					EditMerchant(ctx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true, Data: io.AthMerchant{
						MerchantID:       1,
						MerchantName:     "test",
						Address:          "test",
						Phone:            "test",
						Email:            "test",
						GstNoFile:        "test",
						PanNoFile:        "test",
						BankAccFile:      "test",
						AddressFile:      "test",
						GstFileVerified:  false,
						AddFileVerified:  false,
						BankFileVerified: false,
						PanFileVerified:  false,
						IsActive:         false,
						DoLater:          false,
					}})
				mockOTPService.
					EXPECT().
					EmailSend(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: false, Error: errors.New("failed")})
				mockSQl.ExpectRollback()
				mockSQl.ExpectRollback()
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in PhoneVerifyMerchant :%v", "failed"),
		}, {
			name: "when email service returns success",
			args: args{
				ctx: ctx,
				req: &merchantpb.UploadDocRequest{
					AccountId: 1,
					KycData: &merchantpb.KycDataEdit{
						GstNoFile:   "test",
						PanNoFile:   "test",
						BankAccFile: "test",
						AddressFile: "test",
						DoLater:     false,
					},
				},
			},
			mock: func() {
				mockSQl.ExpectBegin()
				mockSQl.ExpectBegin()
				mockUploadService.
					EXPECT().
					ValidateMerchantDocs(ctx, gomock.Any()).
					Return(io.Response{Success: true, Data: map[string]interface{}{}})
				mockUploadService.
					EXPECT().
					UploadToS3(ctx, gomock.Any()).
					Return(io.Response{Success: true, Data: map[string]interface{}{
						"PanNoFile":   "test",
						"AddressFile": "test",
						"BankAccFile": "test",
						"GstNoFile":   "test",
					}})
				mockMerchantService.
					EXPECT().
					EditMerchant(ctx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true, Data: io.AthMerchant{
						MerchantID:       1,
						MerchantName:     "test",
						Address:          "test",
						Phone:            "test",
						Email:            "test",
						GstNoFile:        "test",
						PanNoFile:        "test",
						BankAccFile:      "test",
						AddressFile:      "test",
						GstFileVerified:  false,
						AddFileVerified:  false,
						BankFileVerified: false,
						PanFileVerified:  false,
						IsActive:         false,
						DoLater:          false,
					}})
				mockOTPService.
					EXPECT().
					EmailSend(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true})
				mockSQl.ExpectCommit()
				mockSQl.ExpectCommit()
			},
			want:    &merchantpb.GenericReply{Status: true},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := handler.UploadDoc(tt.args.ctx, tt.args.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
				assert.Equal(t, tt.want, got)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAthMerchantHandler_ForgotPasswordMerchant(t *testing.T) {
	ctx := context.Background()
	logger := gologger.New("api", true)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sql, mockSQl, err := sqlmock.New()
	assert.NoError(t, err)
	db1, err := gorm.Open("mysql", sql)
	dbConnections := map[string]*gorm.DB{
		MasterDBConnectionName:      db1,
		TransactionDBConnectionName: db1,
	}
	mockDBRepo := merchandal.NewMerchantRepo(logger, dbConnections)
	mockMerchantService := mockmerchant.NewMockAthMerchantService(ctrl)
	mockUserService := mockuser.NewMockAthUserService(ctrl)
	mockOTPService := mockotp.NewMockAthOtpService(ctrl)
	handler := athMerchantHandler{
		logger:             logger,
		DbRepo:             mockDBRepo,
		athMerchantService: mockMerchantService,
		athUserService:     mockUserService,
		athOTPService:      mockOTPService,
	}

	type args struct {
		ctx context.Context
		req *merchantpb.ForgotPasswordRequest
	}
	tests := []struct {
		name        string
		args        args
		mock        func()
		want        *merchantpb.GenericReply
		wantErr     bool
		expectedErr error
	}{
		{
			name: "when invalid email is passed",
			args: args{
				ctx: ctx,
				req: &merchantpb.ForgotPasswordRequest{
					Email: "",
				},
			},
			mock:        func() {},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in ForgotPasswordMerchant :%v", "Request parameter is missing or blank"),
		},
		{
			name: "when forgot password user fails",
			args: args{
				ctx: ctx,
				req: &merchantpb.ForgotPasswordRequest{
					Email: "test@gmail.com",
				},
			},
			mock: func() {
				mockSQl.ExpectBegin()
				mockSQl.ExpectBegin()
				mockOTPService.
					EXPECT().
					EmailSend(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: false, Error: errors.New("failed")})
				mockSQl.ExpectRollback()
				mockSQl.ExpectRollback()
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in ForgotPasswordMerchant :%v", "failed"),
		},
		{
			name: "when forgot password user succeed",
			args: args{
				ctx: ctx,
				req: &merchantpb.ForgotPasswordRequest{
					Email: "test@gmail.com",
				},
			},
			mock: func() {
				mockSQl.ExpectBegin()
				mockSQl.ExpectBegin()
				mockOTPService.
					EXPECT().
					EmailSend(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true})
				mockSQl.ExpectCommit()
				mockSQl.ExpectCommit()
			},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := handler.ForgotPasswordMerchant(tt.args.ctx, tt.args.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
				assert.Equal(t, tt.want, got)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAthMerchantHandler_ResetPasswordMerchant(t *testing.T) {
	ctx := context.Background()
	logger := gologger.New("api", true)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sql, mockSQl, err := sqlmock.New()
	assert.NoError(t, err)
	db1, err := gorm.Open("mysql", sql)
	dbConnections := map[string]*gorm.DB{
		MasterDBConnectionName:      db1,
		TransactionDBConnectionName: db1,
	}
	mockDBRepo := merchandal.NewMerchantRepo(logger, dbConnections)
	mockMerchantService := mockmerchant.NewMockAthMerchantService(ctrl)
	mockUserService := mockuser.NewMockAthUserService(ctrl)
	mockOTPService := mockotp.NewMockAthOtpService(ctrl)
	handler := athMerchantHandler{
		logger:             logger,
		DbRepo:             mockDBRepo,
		athMerchantService: mockMerchantService,
		athUserService:     mockUserService,
		athOTPService:      mockOTPService,
	}

	type args struct {
		ctx context.Context
		req *merchantpb.ResetPasswordRequest
	}
	tests := []struct {
		name        string
		args        args
		mock        func()
		want        *merchantpb.GenericReply
		wantErr     bool
		expectedErr error
	}{
		{
			name: "when invalid new password is passed",
			args: args{
				ctx: ctx,
				req: &merchantpb.ResetPasswordRequest{
					NewPassword: "",
				},
			},
			mock:        func() {},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in ResetPasswordMerchant :%v", "Request parameter is missing or blank"),
		},
		{
			name: "when invalid password token is passed",
			args: args{
				ctx: ctx,
				req: &merchantpb.ResetPasswordRequest{
					NewPassword:        "asd",
					ResetPasswordToken: "",
				},
			},
			mock:        func() {},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in ResetPasswordMerchant :%v", "Request parameter is missing or blank"),
		}, {
			name: "when invalid email is passed",
			args: args{
				ctx: ctx,
				req: &merchantpb.ResetPasswordRequest{
					NewPassword:        "asd",
					ResetPasswordToken: "test",
					Email:              "",
				},
			},
			mock:        func() {},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in ResetPasswordMerchant :%v", "Request parameter is missing or blank"),
		},
		{
			name: "when verify OTP fails",
			args: args{
				ctx: ctx,
				req: &merchantpb.ResetPasswordRequest{
					Email:              "test@gmail.com",
					NewPassword:        "asd",
					ResetPasswordToken: "test",
				},
			},
			mock: func() {
				mockSQl.ExpectBegin()
				mockSQl.ExpectBegin()
				mockOTPService.
					EXPECT().
					VerifyOTP(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: false, Error: errors.New("failed")})
				mockSQl.ExpectRollback()
				mockSQl.ExpectRollback()
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "failed"),
		},
		{
			name: "when reset password user fails",
			args: args{
				ctx: ctx,
				req: &merchantpb.ResetPasswordRequest{
					Email:              "test@gmail.com",
					NewPassword:        "asd",
					ResetPasswordToken: "test",
				},
			},
			mock: func() {
				mockSQl.ExpectBegin()
				mockSQl.ExpectBegin()
				mockOTPService.
					EXPECT().
					VerifyOTP(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true})
				mockUserService.
					EXPECT().
					ResetPasswordUser(ctx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: false, Error: errors.New("failed")})
				mockSQl.ExpectRollback()
				mockSQl.ExpectRollback()
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "failed"),
		},
		{
			name: "when reset password user succeed",
			args: args{
				ctx: ctx,
				req: &merchantpb.ResetPasswordRequest{
					Email:              "test@gmail.com",
					NewPassword:        "asd",
					ResetPasswordToken: "test",
				},
			},
			mock: func() {
				mockSQl.ExpectBegin()
				mockSQl.ExpectBegin()
				mockOTPService.
					EXPECT().
					VerifyOTP(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true})
				mockUserService.
					EXPECT().
					ResetPasswordUser(ctx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true, Data: io.ResetPasswordRequest{
						Email:    "test@gmail.com",
						Password: "asd",
						Token:    "test",
					}})
				mockOTPService.
					EXPECT().
					EmailSend(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true})
				mockSQl.ExpectCommit()
				mockSQl.ExpectCommit()
			},
			want:    &merchantpb.GenericReply{Status: true},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := handler.ResetPasswordMerchant(tt.args.ctx, tt.args.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
				assert.Equal(t, tt.want, got)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAthMerchantHandler_AddTeamMember(t *testing.T) {
	ctx := context.Background()
	authCtx := context.WithValue(ctx, "CurrentUser", &io.AthUser{UserID: 1})
	logger := gologger.New("api", true)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sql, mockSQl, err := sqlmock.New()
	assert.NoError(t, err)
	db1, err := gorm.Open("mysql", sql)
	dbConnections := map[string]*gorm.DB{
		MasterDBConnectionName:      db1,
		TransactionDBConnectionName: db1,
	}
	mockDBRepo := merchandal.NewMerchantRepo(logger, dbConnections)
	mockMerchantService := mockmerchant.NewMockAthMerchantService(ctrl)
	mockUserService := mockuser.NewMockAthUserService(ctrl)
	mockOTPService := mockotp.NewMockAthOtpService(ctrl)
	handler := athMerchantHandler{
		logger:             logger,
		DbRepo:             mockDBRepo,
		athMerchantService: mockMerchantService,
		athUserService:     mockUserService,
		athOTPService:      mockOTPService,
	}

	type args struct {
		ctx context.Context
		req *merchantpb.AddTeamMemberRequest
	}
	tests := []struct {
		name        string
		args        args
		mock        func()
		want        *merchantpb.AddTeamMemberReply
		wantErr     bool
		expectedErr error
	}{
		{
			name: "when invalid first name is passed",
			args: args{
				ctx: authCtx,
				req: &merchantpb.AddTeamMemberRequest{
					FirstName: "",
					LastName:  "user",
					Email:     "test@gmail.com",
					Phone:     "1223",
				},
			},
			mock:        func() {},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in AddTeamMember :%v", "Request parameter is missing or blank"),
		},
		{
			name: "when invalid last name is passed",
			args: args{
				ctx: authCtx,
				req: &merchantpb.AddTeamMemberRequest{
					FirstName: "test",
					LastName:  "",
					Email:     "test@gmail.com",
					Phone:     "1223",
				},
			},
			mock:        func() {},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in AddTeamMember :%v", "Request parameter is missing or blank"),
		}, {
			name: "when invalid email is passed",
			args: args{
				ctx: authCtx,
				req: &merchantpb.AddTeamMemberRequest{
					FirstName: "test",
					LastName:  "user",
					Email:     "",
					Phone:     "1223",
				},
			},
			mock:        func() {},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in AddTeamMember :%v", "Request parameter is missing or blank"),
		}, {
			name: "when invalid phone is passed",
			args: args{
				ctx: authCtx,
				req: &merchantpb.AddTeamMemberRequest{
					FirstName: "test",
					LastName:  "user",
					Email:     "test@gmail.com",
					Phone:     "",
				},
			},
			mock:        func() {},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in AddTeamMember :%v", "Request parameter is missing or blank"),
		},
		{
			name: "when non-logged in user access",
			args: args{
				ctx: ctx,
				req: &merchantpb.AddTeamMemberRequest{
					FirstName:   "test",
					LastName:    "user",
					Email:       "test@gmail.com",
					Phone:       "13422",
					AccountType: "test",
					Privileges: &merchantpb.AccessData{
						VenueIds: []int32{1},
					},
				},
			},
			mock:    func() {},
			want:    nil,
			wantErr: true,
			expectedErr: status.Errorf(codes.Internal,
				status.Error(codes.Unauthenticated, "user must be logged in").Error()),
		},
		{
			name: "when create user fails",
			args: args{
				ctx: authCtx,
				req: &merchantpb.AddTeamMemberRequest{
					FirstName:   "test",
					LastName:    "user",
					Email:       "test@gmail.com",
					Phone:       "13422",
					AccountType: "test",
					Privileges: &merchantpb.AccessData{
						VenueIds: []int32{1},
					},
				},
			},
			mock: func() {
				mockSQl.ExpectBegin()
				mockSQl.ExpectBegin()
				mockUserService.
					EXPECT().
					CreateUser(authCtx, gomock.Any(), gomock.Any(), "teamMember").
					Return(io.Response{Success: false, Error: errors.New("failed")})
				mockSQl.ExpectRollback()
				mockSQl.ExpectRollback()
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in AddTeamMember :%v", "failed"),
		},
		{
			name: "when add team member fails",
			args: args{
				ctx: authCtx,
				req: &merchantpb.AddTeamMemberRequest{
					FirstName:   "test",
					LastName:    "user",
					Email:       "test@gmail.com",
					Phone:       "13422",
					AccountType: "test",
					Privileges: &merchantpb.AccessData{
						VenueIds: []int32{1},
					},
				},
			},
			mock: func() {
				mockSQl.ExpectBegin()
				mockSQl.ExpectBegin()
				mockUserService.
					EXPECT().
					CreateUser(authCtx, gomock.Any(), gomock.Any(), "teamMember").
					Return(io.Response{Success: true, Data: io.AthUser{
						UserID:              2,
						Email:               "test@gmail.com",
						Password:            "test",
						FirstName:           "test",
						LastName:            "test",
						Phone:               "test",
						Gender:              "test",
						DOB:                 "test",
						ProfileImage:        "test",
						LastPasswordResetAt: 0,
						UserSource:          "test",
						LastLoginAt:         0,
						LastLoginIP:         "test",
						IsActive:            true,
						EmailVerify:         true,
						PhoneVerify:         true,
						BmsActive:           true,
					}})
				mockMerchantService.
					EXPECT().
					AddTeamMember(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: false, Error: errors.New("failed")})
				mockSQl.ExpectRollback()
				mockSQl.ExpectRollback()
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in AddTeamMember :%v", "failed"),
		}, {
			name: "when email send fails",
			args: args{
				ctx: authCtx,
				req: &merchantpb.AddTeamMemberRequest{
					FirstName:   "test",
					LastName:    "user",
					Email:       "test@gmail.com",
					Phone:       "13422",
					AccountType: "test",
					Privileges: &merchantpb.AccessData{
						VenueIds: []int32{1},
					},
				},
			},
			mock: func() {
				mockSQl.ExpectBegin()
				mockSQl.ExpectBegin()
				mockUserService.
					EXPECT().
					CreateUser(authCtx, gomock.Any(), gomock.Any(), "teamMember").
					Return(io.Response{Success: true, Data: io.AthUser{
						UserID:              2,
						Email:               "test@gmail.com",
						Password:            "test",
						FirstName:           "test",
						LastName:            "test",
						Phone:               "test",
						Gender:              "test",
						DOB:                 "test",
						ProfileImage:        "test",
						LastPasswordResetAt: 0,
						UserSource:          "test",
						LastLoginAt:         0,
						LastLoginIP:         "test",
						IsActive:            true,
						EmailVerify:         true,
						PhoneVerify:         true,
						BmsActive:           true,
					}})
				mockMerchantService.
					EXPECT().
					AddTeamMember(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true})
				mockOTPService.
					EXPECT().
					EmailSend(authCtx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: false, Error: errors.New("failed")})
				mockSQl.ExpectRollback()
				mockSQl.ExpectRollback()
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in AddTeamMember :%v", "failed"),
		}, {
			name: "when add team member succeeds",
			args: args{
				ctx: authCtx,
				req: &merchantpb.AddTeamMemberRequest{
					FirstName:   "test",
					LastName:    "user",
					Email:       "test@gmail.com",
					Phone:       "13422",
					AccountType: "test",
					Privileges: &merchantpb.AccessData{
						VenueIds: []int32{1},
					},
				},
			},
			mock: func() {
				mockSQl.ExpectBegin()
				mockSQl.ExpectBegin()
				mockUserService.
					EXPECT().
					CreateUser(authCtx, gomock.Any(), gomock.Any(), "teamMember").
					Return(io.Response{Success: true, Data: io.AthUser{
						UserID:              2,
						Email:               "test@gmail.com",
						Password:            "test",
						FirstName:           "test",
						LastName:            "test",
						Phone:               "test",
						Gender:              "test",
						DOB:                 "test",
						ProfileImage:        "test",
						LastPasswordResetAt: 0,
						UserSource:          "test",
						LastLoginAt:         0,
						LastLoginIP:         "test",
						IsActive:            true,
						EmailVerify:         true,
						PhoneVerify:         true,
						BmsActive:           true,
					}})
				mockMerchantService.
					EXPECT().
					AddTeamMember(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true})
				mockOTPService.
					EXPECT().
					EmailSend(authCtx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true})
				mockUserService.
					EXPECT().
					GetAccountDetailsByID(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true, Data: map[string]interface{}{"AccountType": "test", "AccessData": []int32{1}}})
				mockSQl.ExpectCommit()
				mockSQl.ExpectCommit()
			},
			want: &merchantpb.AddTeamMemberReply{
				FirstName:   "test",
				LastName:    "user",
				Email:       "test@gmail.com",
				Phone:       "13422",
				AccountType: "test",
				Privileges: &merchantpb.AccessData{
					VenueIds: []int32{1},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := handler.AddTeamMember(tt.args.ctx, tt.args.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
				assert.Equal(t, tt.want, got)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAthMerchantHandler_PhoneVerifyTeam(t *testing.T) {
	ctx := context.Background()
	logger := gologger.New("api", true)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sql, mockSQl, err := sqlmock.New()
	assert.NoError(t, err)
	db1, err := gorm.Open("mysql", sql)
	dbConnections := map[string]*gorm.DB{
		MasterDBConnectionName:      db1,
		TransactionDBConnectionName: db1,
	}
	mockDBRepo := merchandal.NewMerchantRepo(logger, dbConnections)
	mockMerchantService := mockmerchant.NewMockAthMerchantService(ctrl)
	mockUserService := mockuser.NewMockAthUserService(ctrl)
	mockOTPService := mockotp.NewMockAthOtpService(ctrl)
	handler := athMerchantHandler{
		logger:             logger,
		DbRepo:             mockDBRepo,
		athMerchantService: mockMerchantService,
		athUserService:     mockUserService,
		athOTPService:      mockOTPService,
	}

	type args struct {
		ctx context.Context
		req *merchantpb.PhoneVerifyTeamRequest
	}
	tests := []struct {
		name        string
		args        args
		mock        func()
		want        *merchantpb.PhoneVerifyTeamReply
		wantErr     bool
		expectedErr error
	}{
		{
			name: "when invalid verification code request is passed",
			args: args{
				ctx: ctx,
				req: &merchantpb.PhoneVerifyTeamRequest{
					VerificationCode: "",
				},
			},
			mock:        func() {},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in PhoneVerifyTeam :%v", "Request parameter is missing or blank"),
		},
		{
			name: "when invalid verification type in request is passed",
			args: args{
				ctx: ctx,
				req: &merchantpb.PhoneVerifyTeamRequest{
					VerificationCode: "SMS",
					VerificationType: "",
					Email:            "",
				},
			},
			mock:        func() {},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in PhoneVerifyTeam :%v", "Request parameter is missing or blank"),
		}, {
			name: "when verify OTP fails",
			args: args{
				ctx: ctx,
				req: &merchantpb.PhoneVerifyTeamRequest{
					VerificationCode: "SMS",
					VerificationType: "1234",
					Email:            "test@test.com",
				},
			},
			mock: func() {
				mockSQl.ExpectBegin()
				mockSQl.ExpectBegin()
				mockOTPService.
					EXPECT().
					VerifyOTP(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: false, Error: errors.New("failed")})
				mockSQl.ExpectRollback()
				mockSQl.ExpectRollback()
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "failed"),
		}, {
			name: "when create code fails",
			args: args{
				ctx: ctx,
				req: &merchantpb.PhoneVerifyTeamRequest{
					VerificationCode: "SMS",
					VerificationType: "1234",
					Email:            "test@test.com",
				},
			},
			mock: func() {
				mockSQl.ExpectBegin()
				mockSQl.ExpectBegin()
				mockOTPService.
					EXPECT().
					VerifyOTP(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true, Data: io.AthUser{
						UserID:              1,
						Email:               "test",
						Password:            "",
						FirstName:           "test",
						LastName:            "test",
						Phone:               "test",
						Gender:              "test",
						DOB:                 "test",
						ProfileImage:        "test",
						LastPasswordResetAt: 0,
						UserSource:          "test",
						LastLoginAt:         0,
						LastLoginIP:         "test",
						IsActive:            false,
						EmailVerify:         false,
						PhoneVerify:         false,
						BmsActive:           false,
						TokenValue:          "test",
						TokenHash:           "test",
						TokenExpireAt:       0,
					}})
				mockOTPService.
					EXPECT().CreateCode(ctx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: false, Error: errors.New("failed")})
				mockSQl.ExpectRollback()
				mockSQl.ExpectRollback()

			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in PhoneVerifyTeam :%v", "failed"),
		},
		{
			name: "when request succeeds",
			args: args{
				ctx: ctx,
				req: &merchantpb.PhoneVerifyTeamRequest{
					VerificationType: "SMS",
					VerificationCode: "1234",
					Email:            "test@test.com",
				},
			},
			mock: func() {
				mockSQl.ExpectBegin()
				mockSQl.ExpectBegin()
				mockOTPService.
					EXPECT().
					VerifyOTP(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true, Data: io.AthUser{
						UserID:              1,
						Email:               "test",
						Password:            "",
						FirstName:           "test",
						LastName:            "test",
						Phone:               "test",
						Gender:              "test",
						DOB:                 "test",
						ProfileImage:        "test",
						LastPasswordResetAt: 0,
						UserSource:          "test",
						LastLoginAt:         0,
						LastLoginIP:         "test",
						IsActive:            false,
						EmailVerify:         false,
						PhoneVerify:         false,
						BmsActive:           false,
						TokenValue:          "test",
						TokenHash:           "test",
						TokenExpireAt:       0,
					}})
				mockOTPService.
					EXPECT().CreateCode(ctx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true})
				mockSQl.ExpectCommit()
				mockSQl.ExpectCommit()
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := handler.PhoneVerifyTeam(tt.args.ctx, tt.args.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}
		})
	}

}
