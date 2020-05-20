package handler

import (
	io "api/src/models"
	"errors"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"

	//io "api/src/models"
	"context"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	merchandal "api/src/dal/merchant"

	merchantpb "stash.bms.bz/turf/generic-proto-files.git/merchant/v1"

	"stash.bms.bz/bms/gologger.git"

	mockmerchant "api/src/service/merchant/mock"
	mockotp "api/src/service/otp/mock"
	mockuser "api/src/service/user/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAthMerchantHandler_LoginMerchant(t *testing.T) {
	ctx := context.Background()
	logger := gologger.New("api", true)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sql, _, err := sqlmock.New()
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
		req *merchantpb.LoginRequest
	}
	tests := []struct {
		name        string
		args        args
		mock        func()
		want        *merchantpb.LoginReply
		wantErr     bool
		expectedErr error
	}{
		{
			name: "when invalid email is passed",
			args: args{
				ctx: ctx,
				req: &merchantpb.LoginRequest{
					Login: "",
				},
			},
			mock:        func() {},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in LoginMerchant :%v", "Request parameter is missing or blank"),
		}, {
			name: "when invalid password is passed",
			args: args{
				ctx: ctx,
				req: &merchantpb.LoginRequest{
					Password: "",
				},
			},
			mock:        func() {},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in LoginMerchant :%v", "Request parameter is missing or blank"),
		}, {
			name: "when login user fails",
			args: args{
				ctx: ctx,
				req: &merchantpb.LoginRequest{
					Login:    "test@gmail.com",
					Password: "123",
				},
			},
			mock: func() {
				mockUserService.
					EXPECT().
					LoginUser(ctx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: false, Error: errors.New("failed")})
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "failed"),
		}, {
			name: "when login user succeed",
			args: args{
				ctx: ctx,
				req: &merchantpb.LoginRequest{
					Login:    "test@gmail.com",
					Password: "123",
				},
			},
			mock: func() {
				mockUserService.
					EXPECT().
					LoginUser(ctx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true, Data: map[string]interface{}{
						"AccountType":      "test",
						"MerchantId":       1,
						"MerchantFullName": "test",
						"AccessData":       []int32{1},
						"Token":            "test",
						"LastLoginAt":      0,
						"BMSVerify":        true,
						"EmailVerify":      true,
						"PhoneVerify":      true}})
			},
			want: &merchantpb.LoginReply{
				AccountId:        1,
				MerchantFullName: "test",
				Token:            "test",
				AccountType:      "test",
				EmailVerify:      true,
				PhoneVerify:      true,
				BmsVerify:        true,
				LastLoginAt:      0,
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

			got, err := handler.LoginMerchant(tt.args.ctx, tt.args.req)

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

func TestAthMerchantHandler_GetMerchantByID(t *testing.T) {
	ctx := context.Background()
	authCtx := context.WithValue(ctx, "CurrentUser", &io.AthUser{UserID: 1})

	logger := gologger.New("api", true)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sql, _, err := sqlmock.New()
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
		req *merchantpb.GenericRequest
	}
	tests := []struct {
		name        string
		args        args
		mock        func()
		want        *merchantpb.GetMerchantByIDRes
		wantErr     bool
		expectedErr error
	}{
		{
			name: "when non-logged in user access",
			args: args{
				ctx: ctx,
				req: &merchantpb.GenericRequest{},
			},
			mock:    func() {},
			want:    nil,
			wantErr: true,
			expectedErr: status.Errorf(codes.Internal,
				status.Error(codes.Unauthenticated, "user must be logged in").Error()),
		},
		{
			name: "when get merchantByID req fails",
			args: args{
				ctx: authCtx,
				req: &merchantpb.GenericRequest{},
			},
			mock: func() {
				mockMerchantService.
					EXPECT().
					GetMerchantByUserID(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: false, Error: errors.New("failed")})
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "failed"),
		}, {
			name: "when get userByID req fails",
			args: args{
				ctx: authCtx,
				req: &merchantpb.GenericRequest{},
			},
			mock: func() {
				mockMerchantService.
					EXPECT().
					GetMerchantByUserID(authCtx, gomock.Any(), gomock.Any()).
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
						Models:       io.Models{},
					}})

				mockUserService.
					EXPECT().
					GetUserByID(authCtx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: false, Error: errors.New("failed")})
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "failed"),
		}, {
			name: "when get userByID req succeed",
			args: args{
				ctx: authCtx,
				req: &merchantpb.GenericRequest{},
			},
			mock: func() {
				mockMerchantService.
					EXPECT().
					GetMerchantByUserID(authCtx, gomock.Any(), gomock.Any()).
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
						Models:       io.Models{},
					}})

				mockUserService.
					EXPECT().
					GetUserByID(authCtx, gomock.Any(), gomock.Any(), gomock.Any()).
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
						BmsActive:           false,
						TokenValue:          "test",
						TokenHash:           "test",
						TokenExpireAt:       0,
						Models:              io.Models{},
					}})
			},
			want: &merchantpb.GetMerchantByIDRes{
				MerchantFullName: "test",
				AccountId:        1,
				BusinessName:     "test",
				Address:          "test",
				Phone:            "test",
				Email:            "test",
				KycData: &merchantpb.KycData{
					BankAccFile: "test",
					GstNoFile:   "test",
					AddressFile: "test",
					PanNoFile:   "test",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := handler.GetMerchantByID(tt.args.ctx, tt.args.req)

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

func TestAthMerchantHandler_GetMerchantTeam(t *testing.T) {
	ctx := context.Background()
	authCtx := context.WithValue(ctx, "CurrentUser", &io.AthUser{UserID: 1})
	logger := gologger.New("api", true)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sql, _, err := sqlmock.New()
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
		req *merchantpb.GetMerchantTeamRequest
	}
	tests := []struct {
		name        string
		args        args
		mock        func()
		want        *merchantpb.GetMerchantTeamRes
		wantErr     bool
		expectedErr error
	}{
		{
			name: "when non-logged in user makes request",
			args: args{
				ctx: ctx,
				req: &merchantpb.GetMerchantTeamRequest{},
			},
			mock:    func() {},
			want:    nil,
			wantErr: true,
			expectedErr: status.Errorf(codes.Internal,
				status.Error(codes.Unauthenticated, "user must be logged in").Error()),
		}, {
			name: "when get team fails",
			args: args{
				ctx: authCtx,
				req: &merchantpb.GetMerchantTeamRequest{},
			},
			mock: func() {
				mockMerchantService.
					EXPECT().
					GetTeamData(authCtx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: false, Error: errors.New("failed")})
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "failed"),
		}, {
			name: "when request succeed",
			args: args{
				ctx: authCtx,
				req: &merchantpb.GetMerchantTeamRequest{},
			},
			mock: func() {
				mockMerchantService.
					EXPECT().
					GetTeamData(authCtx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true, Data: []*merchantpb.TeamMemberData{
						{
							AccountId:   1,
							FullName:    "test",
							AccountType: "test",
							Email:       "test@gmail.com",
							Phone:       "test",
							CreatedAt:   0,
							CreatedBy:   "test",
						},
					}})
			},
			want: &merchantpb.GetMerchantTeamRes{
				TeamData: []*merchantpb.TeamMemberData{
					{
						AccountId:   1,
						FullName:    "test",
						AccountType: "test",
						Email:       "test@gmail.com",
						Phone:       "test",
						CreatedAt:   0,
						CreatedBy:   "test",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := handler.GetMerchantTeam(tt.args.ctx, tt.args.req)

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
