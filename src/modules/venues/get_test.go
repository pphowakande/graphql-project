package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/mock/gomock"

	"stash.bms.bz/bms/gologger.git"

	venuedal "api/src/dal/venue"
	io "api/src/models"
	mockvenue "api/src/service/venue/mock"

	venuepb "stash.bms.bz/turf/generic-proto-files.git/venue/v1"
)

func TestAthVenueHandler_GetVenueByID(t *testing.T) {
	ctx := context.Background()
	logger := gologger.New("api", true)
	authCtx := context.WithValue(ctx, "CurrentUser", &io.AthUser{UserID: 1})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sql, _, err := sqlmock.New()
	assert.NoError(t, err)
	db1, err := gorm.Open("mysql", sql)
	dbConnections := map[string]*gorm.DB{
		MasterDBConnectionName: db1,
	}
	mockDBRepo := venuedal.NewVenueRepo(logger, dbConnections)
	mockVenueService := mockvenue.NewMockAthVenueService(ctrl)
	handler := athVenueHandler{
		logger:          logger,
		DbRepo:          mockDBRepo,
		athVenueService: mockVenueService,
	}

	type args struct {
		ctx context.Context
		req *venuepb.GetVenueByIDReq
	}
	tests := []struct {
		name        string
		args        args
		mock        func()
		want        *venuepb.GetVenueByIDRes
		wantErr     bool
		expectedErr error
	}{
		{
			name: "negative 1: when invalid venue_id is passed",
			args: args{
				ctx: authCtx,
				req: &venuepb.GetVenueByIDReq{},
			},
			mock:        func() {},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.InvalidArgument, "Parameter is missing"),
		},
		{
			name: "negative 2: when get venue by Id service fails",
			args: args{
				ctx: authCtx,
				req: &venuepb.GetVenueByIDReq{
					VenueId: 1,
				},
			},
			mock: func() {
				mockVenueService.EXPECT().
					GetVenueByID(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: false, Error: errors.New("failed")})
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "failed"),
		}, {
			name: "positive: when get venue by ID succeed",
			args: args{
				ctx: authCtx,
				req: &venuepb.GetVenueByIDReq{
					VenueId: 1,
				},
			},
			mock: func() {
				mockVenueService.EXPECT().
					GetVenueByID(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true, Data: map[string]interface{}{
						"venueId":     1,
						"name":        "test",
						"description": "test",
						"email":       "test",
						"phone":       "test",
						"images": venuepb.CreateImageData{
							HeaderImg:    []*venuepb.CreateImgData{},
							ThumbnailImg: []*venuepb.CreateImgData{},
							GalleryImg:   []*venuepb.CreateImgData{},
						},
						"amenities":        []*venuepb.AmenityData{},
						"holidays":         []*venuepb.HolidaysData{},
						"hoursOfOperation": []*venuepb.HoursOfOperationData{},
						"addressLine":      "test",
						"landmark":         "test",
						"city":             "test",
						"state":            "test",
						"pincode":          "test",
						"latitude":         float32(23.0),
						"longitude":        float32(23.0),
						"address":          "test",
					}})
			},
			want: &venuepb.GetVenueByIDRes{
				VenueId:          1,
				Name:             "test",
				Description:      "test",
				Address:          "test",
				Phone:            "test",
				Email:            "test",
				Latitude:         23.0,
				Longitude:        23.0,
				Amenities:        []*venuepb.AmenityData{},
				Holidays:         []*venuepb.HolidaysData{},
				HoursOfOperation: []*venuepb.HoursOfOperationData{},
				Images: &venuepb.CreateImageData{
					HeaderImg:    []*venuepb.CreateImgData{},
					ThumbnailImg: []*venuepb.CreateImgData{},
					GalleryImg:   []*venuepb.CreateImgData{},
				}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := handler.GetVenueByID(tt.args.ctx, tt.args.req)

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

func TestAthVenueHandler_GetAllAmenities(t *testing.T) {
	ctx := context.Background()
	logger := gologger.New("api", true)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authCtx := context.WithValue(ctx, "CurrentUser", &io.AthUser{UserID: 1})

	sql, _, err := sqlmock.New()
	assert.NoError(t, err)
	db1, err := gorm.Open("mysql", sql)
	dbConnections := map[string]*gorm.DB{
		MasterDBConnectionName: db1,
	}
	mockDBRepo := venuedal.NewVenueRepo(logger, dbConnections)
	mockVenueService := mockvenue.NewMockAthVenueService(ctrl)
	handler := athVenueHandler{
		logger:          logger,
		DbRepo:          mockDBRepo,
		athVenueService: mockVenueService,
	}

	type args struct {
		ctx context.Context
		req *venuepb.GetAllAmenitiesReq
	}
	tests := []struct {
		name        string
		args        args
		mock        func()
		want        *venuepb.GetAllAmenitiesRes
		wantErr     bool
		expectedErr error
	}{
		{
			name: "negative 1: when get all amenities service call fails",
			args: args{
				ctx: authCtx,
				req: &venuepb.GetAllAmenitiesReq{},
			},
			mock: func() {
				mockVenueService.EXPECT().
					GetAllAmenities(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: false, Error: errors.New("failed")})
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "failed"),
		},
		{
			name: "negative 2: when no amenities are found",
			args: args{
				ctx: authCtx,
				req: &venuepb.GetAllAmenitiesReq{},
			},
			mock: func() {
				mockVenueService.EXPECT().
					GetAllAmenities(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true, Data: nil})
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.NotFound, "No amenities present"),
		}, {
			name: "positive: when getAllAmenities succeeds",
			args: args{
				ctx: authCtx,
				req: &venuepb.GetAllAmenitiesReq{},
			},
			mock: func() {
				mockVenueService.EXPECT().
					GetAllAmenities(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true, Data: []*venuepb.AmenityData{
						{
							AmenityId:   1,
							AmenityName: "test",
						},
					},
					})
			},
			want: &venuepb.GetAllAmenitiesRes{
				AmenityData: []*venuepb.AmenityData{
					{
						AmenityId:   1,
						AmenityName: "test",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := handler.GetAllAmenities(tt.args.ctx, tt.args.req)

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

func TestAthVenueHandler_GetListOfVenueByMerchantID(t *testing.T) {
	ctx := context.Background()
	logger := gologger.New("api", true)
	authCtx := context.WithValue(ctx, "CurrentUser", &io.AthUser{UserID: 1})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sql, _, err := sqlmock.New()
	assert.NoError(t, err)
	db1, err := gorm.Open("mysql", sql)
	dbConnections := map[string]*gorm.DB{
		MasterDBConnectionName: db1,
	}
	mockDBRepo := venuedal.NewVenueRepo(logger, dbConnections)
	mockVenueService := mockvenue.NewMockAthVenueService(ctrl)
	handler := athVenueHandler{
		logger:          logger,
		DbRepo:          mockDBRepo,
		athVenueService: mockVenueService,
	}

	type args struct {
		ctx context.Context
		req *venuepb.GetListOfVenueByMerchantIDReq
	}
	tests := []struct {
		name        string
		args        args
		mock        func()
		want        *venuepb.GetListOfVenueByMerchantIDRes
		wantErr     bool
		expectedErr error
	}{
		{
			name: "negative 1: when get list of venue by merchantID service call fails",
			args: args{
				ctx: authCtx,
				req: &venuepb.GetListOfVenueByMerchantIDReq{},
			},
			mock: func() {
				mockVenueService.EXPECT().
					GetListOfVenueByMerchantID(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: false, Error: errors.New("failed")})
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "failed"),
		},
		{
			name: "negative 2: when no venue are found",
			args: args{
				ctx: authCtx,
				req: &venuepb.GetListOfVenueByMerchantIDReq{},
			},
			mock: func() {
				mockVenueService.EXPECT().
					GetListOfVenueByMerchantID(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true, Data: nil})
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.NotFound, "No venues present"),
		}, {
			name: "positive: when getAllAmenities succeeds",
			args: args{
				ctx: authCtx,
				req: &venuepb.GetListOfVenueByMerchantIDReq{},
			},
			mock: func() {
				mockVenueService.EXPECT().
					GetListOfVenueByMerchantID(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true, Data: []io.AthVenues{
						{
							VenueID:   1,
							VenueName: "test",
							IsActive:  true,
						},
					},
					})
			},
			want: &venuepb.GetListOfVenueByMerchantIDRes{
				VenueData: []*venuepb.VenueList{
					{
						VenueId:   1,
						VenueName: "test",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := handler.GetListOfVenueByMerchantID(tt.args.ctx, tt.args.req)

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
