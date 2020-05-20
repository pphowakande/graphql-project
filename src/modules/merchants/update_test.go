package handler

import (
	io "api/src/models"
	"errors"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"

	//"errors"

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

func TestAthMerchantHandler_EditMerchant(t *testing.T) {
	ctx := context.Background()
	authCtx := context.WithValue(ctx, "CurrentUser", &io.AthUser{UserID: 1})
	logger := gologger.New("api", true)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sql, mockSQL, err := sqlmock.New()
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
		req *merchantpb.EditMerchantRequest
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
			name: "when non-logged in user is passed in context",
			args: args{
				ctx: ctx,
				req: &merchantpb.EditMerchantRequest{
					MerchantFullName: "",
				},
			},
			mock:    func() {},
			want:    nil,
			wantErr: true,
			expectedErr: status.Errorf(codes.Internal,
				status.Error(codes.Unauthenticated, "user must be logged in").Error()),
		},
		{
			name: "when invalid merchant full name is passed",
			args: args{
				ctx: authCtx,
				req: &merchantpb.EditMerchantRequest{
					MerchantFullName: "",
				},
			},
			mock:        func() {},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in EditMerchant :%v", "Request parameter is missing or blank"),
		},
		{
			name: "when invalid business name is passed",
			args: args{
				ctx: authCtx,
				req: &merchantpb.EditMerchantRequest{
					BusinessName: "",
				},
			},
			mock:        func() {},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in EditMerchant :%v", "Request parameter is missing or blank"),
		}, {
			name: "when invalid email is passed",
			args: args{
				ctx: authCtx,
				req: &merchantpb.EditMerchantRequest{
					Email: "",
				},
			},
			mock:        func() {},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in EditMerchant :%v", "Request parameter is missing or blank"),
		}, {
			name: "when invalid phone is passed",
			args: args{
				ctx: authCtx,
				req: &merchantpb.EditMerchantRequest{
					Phone: "",
				},
			},
			mock:        func() {},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in EditMerchant :%v", "Request parameter is missing or blank"),
		}, {
			name: "when user service returns error while editing",
			args: args{
				ctx: authCtx,
				req: &merchantpb.EditMerchantRequest{
					Email:            "test@gmail.com",
					BusinessName:     "test",
					MerchantFullName: "Test User",
					Phone:            "80000",
				},
			},
			mock: func() {
				mockSQL.ExpectBegin()
				mockSQL.ExpectBegin()
				mockUserService.EXPECT().
					EditUser(authCtx, gomock.Any(), gomock.Any(), "owner").
					Return(io.Response{Success: false, Error: errors.New("failed")})
					//mockMerchantService.
					//	EXPECT().
					//	CreateMerchant(ctx, gomock.Any()).
					//	Return(io.Response{Success: true, Data: map[string]interface{}{"Merchant_id": 1}})
					//mockMerchantService.
					//	EXPECT().
					//	CreateMerchantUser(ctx, gomock.Any()).
					//	Return(io.Response{Success: true})
				mockSQL.ExpectRollback()
				mockSQL.ExpectRollback()
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "failed"),
		}, {
			name: "when user service returns error while editing merchant",
			args: args{
				ctx: authCtx,
				req: &merchantpb.EditMerchantRequest{
					Email:            "test@gmail.com",
					BusinessName:     "test",
					MerchantFullName: "Test User",
					Phone:            "80000",
				},
			},
			mock: func() {
				mockSQL.ExpectBegin()
				mockSQL.ExpectBegin()
				mockUserService.EXPECT().
					EditUser(authCtx, gomock.Any(), gomock.Any(), "owner").
					Return(io.Response{Success: true, Data: map[string]interface{}{"UserID": 1, "Email": "test@gmail.com", "EmailVerify": false, "PhoneVerify": false}})
				mockMerchantService.
					EXPECT().
					EditMerchant(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: false, Error: errors.New("failed")})
				mockSQL.ExpectRollback()
				mockSQL.ExpectRollback()
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "failed"),
		}, {
			name: "when user service returns error while fetching merchant byID",
			args: args{
				ctx: authCtx,
				req: &merchantpb.EditMerchantRequest{
					Email:            "test@gmail.com",
					BusinessName:     "test",
					MerchantFullName: "Test User",
					Phone:            "80000",
				},
			},
			mock: func() {
				mockSQL.ExpectBegin()
				mockSQL.ExpectBegin()
				mockUserService.EXPECT().
					EditUser(authCtx, gomock.Any(), gomock.Any(), "owner").
					Return(io.Response{Success: true, Data: map[string]interface{}{"UserID": 1, "Email": "test@gmail.com", "EmailVerify": false, "PhoneVerify": false}})
				mockMerchantService.
					EXPECT().
					EditMerchant(authCtx, gomock.Any(), gomock.Any()).
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
				mockMerchantService.
					EXPECT().
					GetMerchantByID(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: false, Error: errors.New("failed")})
				mockSQL.ExpectRollback()
				mockSQL.ExpectRollback()
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "failed"),
		}, {
			name: "when user service returns error while editing merchant",
			args: args{
				ctx: authCtx,
				req: &merchantpb.EditMerchantRequest{
					Email:            "test@gmail.com",
					BusinessName:     "test",
					MerchantFullName: "Test User",
					Phone:            "80000",
				},
			},
			mock: func() {
				mockSQL.ExpectBegin()
				mockSQL.ExpectBegin()
				mockUserService.EXPECT().
					EditUser(authCtx, gomock.Any(), gomock.Any(), "owner").
					Return(io.Response{Success: true, Data: map[string]interface{}{"UserID": 1, "Email": "test@gmail.com", "EmailVerify": false, "PhoneVerify": false}})
				mockMerchantService.
					EXPECT().
					EditMerchant(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: false, Error: errors.New("failed")})
				mockSQL.ExpectRollback()
				mockSQL.ExpectRollback()
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "failed"),
		}, {
			name: "when email send fails",
			args: args{
				ctx: authCtx,
				req: &merchantpb.EditMerchantRequest{
					Email:            "test@gmail.com",
					BusinessName:     "test",
					MerchantFullName: "Test User",
					Phone:            "80000",
				},
			},
			mock: func() {
				mockSQL.ExpectBegin()
				mockSQL.ExpectBegin()
				mockUserService.EXPECT().
					EditUser(authCtx, gomock.Any(), gomock.Any(), "owner").
					Return(io.Response{Success: true, Data: map[string]interface{}{"UserID": 1, "Email": "test@gmail.com", "EmailVerify": false, "PhoneVerify": false}})
				mockMerchantService.
					EXPECT().
					EditMerchant(authCtx, gomock.Any(), gomock.Any()).
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
				mockMerchantService.
					EXPECT().
					GetMerchantByID(authCtx, gomock.Any(), gomock.Any()).
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
						DoLater:      false}})
				mockOTPService.
					EXPECT().
					EmailSend(authCtx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: false, Error: errors.New("failed")})
				mockSQL.ExpectRollback()
				mockSQL.ExpectRollback()
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in EditMerchant :%v", "failed"),
		}, {
			name: "when create OTP fails",
			args: args{
				ctx: authCtx,
				req: &merchantpb.EditMerchantRequest{
					Email:            "test@gmail.com",
					BusinessName:     "test",
					MerchantFullName: "Test User",
					Phone:            "80000",
				},
			},
			mock: func() {
				mockSQL.ExpectBegin()
				mockSQL.ExpectBegin()
				mockUserService.EXPECT().
					EditUser(authCtx, gomock.Any(), gomock.Any(), "owner").
					Return(io.Response{Success: true, Data: map[string]interface{}{"UserID": 1, "Email": "test@gmail.com", "EmailVerify": false, "PhoneVerify": false}})
				mockMerchantService.
					EXPECT().
					EditMerchant(authCtx, gomock.Any(), gomock.Any()).
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
				mockMerchantService.
					EXPECT().
					GetMerchantByID(authCtx, gomock.Any(), gomock.Any()).
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
						DoLater:      false}})
				mockOTPService.
					EXPECT().
					EmailSend(authCtx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true})
				mockOTPService.
					EXPECT().
					CreateOTP(authCtx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: false, Error: errors.New("failed")})
				mockSQL.ExpectRollback()
				mockSQL.ExpectRollback()
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in EditMerchant :%v", "failed"),
		},
		{
			name: "when valid merchant with full name is passed",
			args: args{
				ctx: authCtx,
				req: &merchantpb.EditMerchantRequest{
					Email:            "test@gmail.com",
					BusinessName:     "test",
					MerchantFullName: "Test User",
					Phone:            "80000",
				},
			},
			mock: func() {
				mockSQL.ExpectBegin()
				mockSQL.ExpectBegin()
				mockUserService.EXPECT().
					EditUser(authCtx, gomock.Any(), gomock.Any(), "owner").
					Return(io.Response{Success: true, Data: map[string]interface{}{"UserID": 1, "Email": "test@gmail.com", "EmailVerify": false, "PhoneVerify": false}})
				mockMerchantService.
					EXPECT().
					EditMerchant(authCtx, gomock.Any(), gomock.Any()).
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
				mockMerchantService.
					EXPECT().
					GetMerchantByID(authCtx, gomock.Any(), gomock.Any()).
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
					EmailSend(authCtx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true})
				mockOTPService.
					EXPECT().
					CreateOTP(authCtx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true})
				mockSQL.ExpectCommit()
				mockSQL.ExpectCommit()
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "when valid user with only first name is passed",
			args: args{
				ctx: authCtx,
				req: &merchantpb.EditMerchantRequest{
					Email:            "test@gmail.com",
					BusinessName:     "test",
					MerchantFullName: "Test",
					Phone:            "80000",
				},
			},
			mock: func() {
				mockSQL.ExpectBegin()
				mockSQL.ExpectBegin()
				mockUserService.EXPECT().
					EditUser(authCtx, gomock.Any(), gomock.Any(), "owner").
					Return(io.Response{Success: true, Data: map[string]interface{}{"UserID": 1, "Email": "test@gmail.com", "EmailVerify": false, "PhoneVerify": false}})
				mockMerchantService.
					EXPECT().
					EditMerchant(authCtx, gomock.Any(), gomock.Any()).
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
				mockMerchantService.
					EXPECT().
					GetMerchantByID(authCtx, gomock.Any(), gomock.Any()).
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
					EmailSend(authCtx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true})
				mockOTPService.
					EXPECT().
					CreateOTP(authCtx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true})
				mockSQL.ExpectCommit()
				mockSQL.ExpectCommit()
			},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := handler.EditMerchant(tt.args.ctx, tt.args.req)

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

func TestAthMerchantHandler_EditTeamMember(t *testing.T) {
	ctx := context.Background()
	authCtx := context.WithValue(ctx, "CurrentUser", &io.AthUser{UserID: 1})
	logger := gologger.New("api", true)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sql, mockSQL, err := sqlmock.New()
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
		req *merchantpb.EditTeamMemberRequest
	}
	tests := []struct {
		name        string
		args        args
		mock        func()
		want        *merchantpb.EditTeamMemberReply
		wantErr     bool
		expectedErr error
	}{
		{
			name: "when invalid accountID is passed",
			args: args{
				ctx: authCtx,
				req: &merchantpb.EditTeamMemberRequest{
					AccountId: 0,
				},
			},
			mock:        func() {},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in EditTeamMember :%v", "Request parameter is missing or blank"),
		}, {
			name: "when non-logged in user makes request",
			args: args{
				ctx: ctx,
				req: &merchantpb.EditTeamMemberRequest{
					AccountId: 1,
				},
			},
			mock:    func() {},
			want:    nil,
			wantErr: true,
			expectedErr: status.Errorf(codes.Internal,
				status.Error(codes.Unauthenticated, "user must be logged in").Error()),
		}, {
			name: "when update team member privileges fails",
			args: args{
				ctx: authCtx,
				req: &merchantpb.EditTeamMemberRequest{
					FirstName: "test",
					LastName:  "test",
					Email:     "test",
					Phone:     "test",
					AccountId: 1,
					Privileges: []*merchantpb.EditAccessData{
						{
							VenueId: 1,
							Status:  true,
						},
					},
				},
			},
			mock: func() {
				mockSQL.ExpectBegin()
				mockSQL.ExpectBegin()
				mockMerchantService.
					EXPECT().
					UpdateTeamMemberPrivileges(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: false, Error: errors.New("failed")})
				mockSQL.ExpectRollback()
				mockSQL.ExpectRollback()
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in EditTeamMember :%v", "failed"),
		}, {
			name: "when edit user fails",
			args: args{
				ctx: authCtx,
				req: &merchantpb.EditTeamMemberRequest{
					FirstName: "test",
					LastName:  "test",
					Email:     "test",
					Phone:     "test",
					AccountId: 1,
					Privileges: []*merchantpb.EditAccessData{
						{
							VenueId: 1,
							Status:  true,
						},
					},
				},
			},
			mock: func() {
				mockSQL.ExpectBegin()
				mockSQL.ExpectBegin()
				mockMerchantService.
					EXPECT().
					UpdateTeamMemberPrivileges(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true})
				mockUserService.
					EXPECT().
					EditUser(authCtx, gomock.Any(), gomock.Any(), "teammember").
					Return(io.Response{Success: false, Error: errors.New("failed")})
				mockSQL.ExpectRollback()
				mockSQL.ExpectRollback()
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in EditTeamMember :%v", "failed"),
		}, {
			name: "when email send fails",
			args: args{
				ctx: authCtx,
				req: &merchantpb.EditTeamMemberRequest{
					FirstName: "test",
					LastName:  "test",
					Email:     "test",
					Phone:     "test",
					AccountId: 1,
					Privileges: []*merchantpb.EditAccessData{
						{
							VenueId: 1,
							Status:  true,
						},
					},
				},
			},
			mock: func() {
				mockSQL.ExpectBegin()
				mockSQL.ExpectBegin()
				mockMerchantService.
					EXPECT().
					UpdateTeamMemberPrivileges(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true})
				mockUserService.
					EXPECT().
					EditUser(authCtx, gomock.Any(), gomock.Any(), "teammember").
					Return(io.Response{Success: true, Data: map[string]interface{}{
						"PhoneVerify": false,
						"EmailVerify": false,
						"Email":       "test@gmail.com",
						"UserID":      1,
						"FirstName":   "test",
						"LastName":    "test",
						"Phone":       "test",
						"AccountType": "member",
						"AccessData":  []int32{1},
						"AccountId":   1,
					}})
				mockOTPService.
					EXPECT().
					EmailSend(authCtx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: false, Error: errors.New("failed")})
				mockSQL.ExpectRollback()
				mockSQL.ExpectRollback()
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in EditTeamMember :%v", "failed"),
		}, {
			name: "when create OTP fails",
			args: args{
				ctx: authCtx,
				req: &merchantpb.EditTeamMemberRequest{
					FirstName: "test",
					LastName:  "test",
					Email:     "test",
					Phone:     "test",
					AccountId: 1,
					Privileges: []*merchantpb.EditAccessData{
						{
							VenueId: 1,
							Status:  true,
						},
					},
				},
			},
			mock: func() {
				mockSQL.ExpectBegin()
				mockSQL.ExpectBegin()
				mockMerchantService.
					EXPECT().
					UpdateTeamMemberPrivileges(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true})
				mockUserService.
					EXPECT().
					EditUser(authCtx, gomock.Any(), gomock.Any(), "teammember").
					Return(io.Response{Success: true, Data: map[string]interface{}{
						"PhoneVerify": false,
						"EmailVerify": false,
						"Email":       "test@gmail.com",
						"UserID":      1,
						"FirstName":   "test",
						"LastName":    "test",
						"Phone":       "test",
						"AccountType": "member",
						"AccessData":  []int32{1},
						"AccountId":   1,
					}})
				mockOTPService.
					EXPECT().
					EmailSend(authCtx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true})
				mockOTPService.
					EXPECT().
					CreateOTP(authCtx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: false, Error: errors.New("failed")})
				mockSQL.ExpectRollback()
				mockSQL.ExpectRollback()
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in EditTeamMember :%v", "failed"),
		}, {
			name: "when request succeeds",
			args: args{
				ctx: authCtx,
				req: &merchantpb.EditTeamMemberRequest{
					FirstName: "test",
					LastName:  "test",
					Email:     "test",
					Phone:     "test",
					AccountId: 1,
					Privileges: []*merchantpb.EditAccessData{
						{
							VenueId: 1,
							Status:  true,
						},
					},
				},
			},
			mock: func() {
				mockSQL.ExpectBegin()
				mockSQL.ExpectBegin()
				mockMerchantService.
					EXPECT().
					UpdateTeamMemberPrivileges(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true})
				mockUserService.
					EXPECT().
					EditUser(authCtx, gomock.Any(), gomock.Any(), "teammember").
					Return(io.Response{Success: true, Data: map[string]interface{}{
						"PhoneVerify": false,
						"EmailVerify": false,
						"Email":       "test@gmail.com",
						"UserID":      1,
						"FirstName":   "test",
						"LastName":    "test",
						"Phone":       "test",
						"AccountType": "member",
						"AccessData":  []int32{1},
						"AccountId":   1,
					}})
				mockOTPService.
					EXPECT().
					EmailSend(authCtx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true})
				mockOTPService.
					EXPECT().
					CreateOTP(authCtx, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true})
				mockSQL.ExpectCommit()
				mockSQL.ExpectCommit()
			},
			want: &merchantpb.EditTeamMemberReply{
				FirstName:   "test",
				LastName:    "test",
				Email:       "test@gmail.com",
				Phone:       "test",
				AccountType: "member",
				Privileges: &merchantpb.AccessData{
					VenueIds: []int32{1},
				},
				AccountId: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := handler.EditTeamMember(tt.args.ctx, tt.args.req)

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
