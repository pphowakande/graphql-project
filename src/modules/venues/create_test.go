package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/assert"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"stash.bms.bz/bms/gologger.git"

	venuedal "api/src/dal/venue"
	io "api/src/models"

	mockvenue "api/src/service/venue/mock"

	venuepb "stash.bms.bz/turf/generic-proto-files.git/venue/v1"
)

func TestAthVenueHandler_CreateVenue(t *testing.T) {
	ctx := context.Background()
	authCtx := context.WithValue(ctx, "CurrentUser", &io.AthUser{UserID: 1})
	logger := gologger.New("api", true)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sql, mockSQL, err := sqlmock.New()
	assert.NoError(t, err)
	db1, err := gorm.Open("mysql", sql)
	dbConnections := map[string]*gorm.DB{
		MasterDBConnectionName: db1,
	}
	mockDBRepo := venuedal.NewVenueRepo(logger, dbConnections)
	mockVenueService := mockvenue.NewMockAthVenueService(ctrl)
	//mockUserService := mockuser.NewMockAthUserService(ctrl)
	//mockOTPService := mockotp.NewMockAthOtpService(ctrl)
	//
	//mockMerchantService := mockmerchant.NewMockAthMerchantService(ctrl)
	//mockUserService := mockuser.NewMockAthUserService(ctrl)
	//mockOTPService := mockotp.NewMockAthOtpService(ctrl)
	handler := athVenueHandler{
		logger:          logger,
		DbRepo:          mockDBRepo,
		athVenueService: mockVenueService,
	}

	type args struct {
		ctx context.Context
		req *venuepb.CreateVenueRequest
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
			name: "negative 1: when non-logged in user is passed in context",
			args: args{
				ctx: ctx,
				req: &venuepb.CreateVenueRequest{
					Name:             "test",
					Description:      "test",
					Address:          "test",
					Phone:            "test",
					Email:            "test",
					Latitude:         23.0,
					Longitude:        23.0,
					Amenities:        nil,
					HoursOfOperation: nil,
					Holidays:         nil,
					Images:           nil,
				},
			},
			mock:    func() {},
			want:    nil,
			wantErr: true,
			expectedErr: status.Errorf(codes.Internal,
				status.Error(codes.Unauthenticated, "user must be logged in").Error()),
		},
		{
			name: "negative 2: when create venue fails",
			args: args{
				ctx: authCtx,
				req: &venuepb.CreateVenueRequest{
					Name:             "test",
					Description:      "test",
					Address:          "test",
					Phone:            "test",
					Email:            "test",
					Latitude:         23.0,
					Longitude:        23.0,
					Amenities:        nil,
					HoursOfOperation: nil,
					Holidays:         nil,
					Images:           nil,
				},
			},
			mock: func() {
				mockSQL.ExpectBegin()
				mockVenueService.EXPECT().
					CreateVenue(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: false, Error: errors.New("failed")})
				mockSQL.ExpectRollback()
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in create venue: %v", "failed"),
		},
		{
			name: "negative 3: when create venue amenity fails",
			args: args{
				ctx: authCtx,
				req: &venuepb.CreateVenueRequest{
					Name:        "test",
					Description: "test",
					Address:     "test",
					Phone:       "test",
					Email:       "test",
					Latitude:    23.0,
					Longitude:   23.0,
					Amenities: []*venuepb.CreateAmenitiesData{
						{
							AmenityId: 1,
						},
					},
					HoursOfOperation: nil,
					Holidays:         nil,
					Images:           nil,
				},
			},
			mock: func() {
				mockSQL.ExpectBegin()
				mockVenueService.EXPECT().
					CreateVenue(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true, Data: map[string]interface{}{
						"venue_id": 1,
					}})
				mockVenueService.EXPECT().
					CreateVenueAmenity(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: false, Error: errors.New("failed")})
				mockSQL.ExpectRollback()
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "error in create venue: %v", "failed"),
		}, {
			name: "negative 4: when saving hours of operation fails",
			args: args{
				ctx: authCtx,
				req: &venuepb.CreateVenueRequest{
					Name:        "test",
					Description: "test",
					Address:     "test",
					Phone:       "test",
					Email:       "test",
					Latitude:    23.0,
					Longitude:   23.0,
					Amenities: []*venuepb.CreateAmenitiesData{
						{
							AmenityId: 1,
						},
					},
					HoursOfOperation: []*venuepb.HoursOfOperationData{
						{
							Day: "test",
							Timing: []*venuepb.Timing{
								{
									OpeningTime: "test",
									ClosingTime: "test",
								},
							},
						},
					},
					Holidays: nil,
					Images:   nil,
				},
			},
			mock: func() {
				mockSQL.ExpectBegin()
				mockVenueService.EXPECT().
					CreateVenue(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true, Data: map[string]interface{}{
						"venue_id": 1,
					}})
				mockVenueService.EXPECT().
					CreateVenueAmenity(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true})
				mockVenueService.EXPECT().
					SaveHoursOfOperation(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: false, Error: errors.New("failed")})
				mockSQL.ExpectRollback()
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "failed"),
		}, {
			name: "negative 5: when create venue holiday fails",
			args: args{
				ctx: authCtx,
				req: &venuepb.CreateVenueRequest{
					Name:        "test",
					Description: "test",
					Address:     "test",
					Phone:       "test",
					Email:       "test",
					Latitude:    23.0,
					Longitude:   23.0,
					Amenities: []*venuepb.CreateAmenitiesData{
						{
							AmenityId: 1,
						},
					},
					HoursOfOperation: []*venuepb.HoursOfOperationData{
						{
							Day: "test",
							Timing: []*venuepb.Timing{
								{
									OpeningTime: "test",
									ClosingTime: "test",
								},
							},
						},
					},
					Holidays: []*venuepb.CreateVenueHolidayRequest{
						{
							Title: "test",
							Date:  "test",
						},
					},
					Images: nil,
				},
			},
			mock: func() {
				mockSQL.ExpectBegin()
				mockVenueService.EXPECT().
					CreateVenue(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true, Data: map[string]interface{}{
						"venue_id": 1,
					}})
				mockVenueService.EXPECT().
					CreateVenueAmenity(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true})
				mockVenueService.EXPECT().
					SaveHoursOfOperation(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true})
				mockVenueService.EXPECT().
					CreateVenueHoliday(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: false, Error: errors.New("failed")})
				mockSQL.ExpectRollback()
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "failed"),
		}, {
			name: "negative 6: when get venue byID fails",
			args: args{
				ctx: authCtx,
				req: &venuepb.CreateVenueRequest{
					Name:        "test",
					Description: "test",
					Address:     "test",
					Phone:       "test",
					Email:       "test",
					Latitude:    23.0,
					Longitude:   23.0,
					Amenities: []*venuepb.CreateAmenitiesData{
						{
							AmenityId: 1,
						},
					},
					HoursOfOperation: []*venuepb.HoursOfOperationData{
						{
							Day: "test",
							Timing: []*venuepb.Timing{
								{
									OpeningTime: "test",
									ClosingTime: "test",
								},
							},
						},
					},
					Holidays: []*venuepb.CreateVenueHolidayRequest{
						{
							Title: "test",
							Date:  "test",
						},
					},
					Images: nil,
				},
			},
			mock: func() {
				mockSQL.ExpectBegin()
				mockVenueService.EXPECT().
					CreateVenue(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true, Data: map[string]interface{}{
						"venue_id": 1,
					}})
				mockVenueService.EXPECT().
					CreateVenueAmenity(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true})
				mockVenueService.EXPECT().
					SaveHoursOfOperation(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true})
				mockVenueService.EXPECT().
					CreateVenueHoliday(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true})
				mockVenueService.EXPECT().
					GetVenueByID(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: false, Error: errors.New("failed")})
				mockSQL.ExpectRollback()
			},
			want:        nil,
			wantErr:     true,
			expectedErr: status.Errorf(codes.Internal, "failed"),
		}, {
			name: "positive: successfully create venue request",
			args: args{
				ctx: authCtx,
				req: &venuepb.CreateVenueRequest{
					Name:        "test",
					Description: "test",
					Address:     "test",
					Phone:       "test",
					Email:       "test",
					Latitude:    23.0,
					Longitude:   23.0,
					Amenities: []*venuepb.CreateAmenitiesData{
						{
							AmenityId: 1,
						},
					},
					HoursOfOperation: []*venuepb.HoursOfOperationData{
						{
							Day: "test",
							Timing: []*venuepb.Timing{
								{
									OpeningTime: "test",
									ClosingTime: "test",
								},
							},
						},
					},
					Holidays: []*venuepb.CreateVenueHolidayRequest{
						{
							Title: "test",
							Date:  "test",
						},
					},
					Images: nil,
				},
			},
			mock: func() {
				mockSQL.ExpectBegin()
				mockVenueService.EXPECT().
					CreateVenue(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true, Data: map[string]interface{}{
						"venue_id": 1,
					}})
				mockVenueService.EXPECT().
					CreateVenueAmenity(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true})
				mockVenueService.EXPECT().
					SaveHoursOfOperation(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true})
				mockVenueService.EXPECT().
					CreateVenueHoliday(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true})
				mockVenueService.EXPECT().
					GetVenueByID(authCtx, gomock.Any(), gomock.Any()).
					Return(io.Response{Success: true, Data: map[string]interface{}{
						"venueId":     1,
						"name":        "test",
						"description": "test",
						"email":       "test",
						"phone":       "test",
						"address":     "test",
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
					}})
				mockSQL.ExpectCommit()
			},
			want: &venuepb.GetVenueByIDRes{
				VenueId:          1,
				Name:             "test",
				Description:      "test",
				Address:          "test",
				Phone:            "test",
				Email:            "test",
				Latitude:         float32(23.0),
				Longitude:        float32(23.0),
				Amenities:        []*venuepb.AmenityData{},
				Holidays:         []*venuepb.HolidaysData{},
				HoursOfOperation: []*venuepb.HoursOfOperationData{},
				Images: &venuepb.CreateImageData{
					HeaderImg:    []*venuepb.CreateImgData{},
					ThumbnailImg: []*venuepb.CreateImgData{},
					GalleryImg:   []*venuepb.CreateImgData{},
				},
			},
			wantErr:     false,
			expectedErr: status.Errorf(codes.Internal, "failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := handler.CreateVenue(tt.args.ctx, tt.args.req)

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
