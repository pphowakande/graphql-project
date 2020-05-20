// Automatically generated by MockGen. DO NOT EDIT!
// Source: api/src/service/venue (interfaces: AthVenueService)

package mock

import (
	models "api/src/models"
	context "context"
	gomock "github.com/golang/mock/gomock"
	gorm "github.com/jinzhu/gorm"
	v1 "stash.bms.bz/turf/generic-proto-files.git/venue/v1"
)

// Mock of AthVenueService interface
type MockAthVenueService struct {
	ctrl     *gomock.Controller
	recorder *_MockAthVenueServiceRecorder
}

// Recorder for MockAthVenueService (not exported)
type _MockAthVenueServiceRecorder struct {
	mock *MockAthVenueService
}

func NewMockAthVenueService(ctrl *gomock.Controller) *MockAthVenueService {
	mock := &MockAthVenueService{ctrl: ctrl}
	mock.recorder = &_MockAthVenueServiceRecorder{mock}
	return mock
}

func (_m *MockAthVenueService) EXPECT() *_MockAthVenueServiceRecorder {
	return _m.recorder
}

func (_m *MockAthVenueService) CreateVenue(_param0 context.Context, _param1 *gorm.DB, _param2 models.AthVenues) models.Response {
	ret := _m.ctrl.Call(_m, "CreateVenue", _param0, _param1, _param2)
	ret0, _ := ret[0].(models.Response)
	return ret0
}

func (_mr *_MockAthVenueServiceRecorder) CreateVenue(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateVenue", arg0, arg1, arg2)
}

func (_m *MockAthVenueService) CreateVenueAmenity(_param0 context.Context, _param1 *gorm.DB, _param2 models.AthVenueAmenitiesData) models.Response {
	ret := _m.ctrl.Call(_m, "CreateVenueAmenity", _param0, _param1, _param2)
	ret0, _ := ret[0].(models.Response)
	return ret0
}

func (_mr *_MockAthVenueServiceRecorder) CreateVenueAmenity(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateVenueAmenity", arg0, arg1, arg2)
}

func (_m *MockAthVenueService) CreateVenueHoliday(_param0 context.Context, _param1 *gorm.DB, _param2 models.AthVenueHolidays) models.Response {
	ret := _m.ctrl.Call(_m, "CreateVenueHoliday", _param0, _param1, _param2)
	ret0, _ := ret[0].(models.Response)
	return ret0
}

func (_mr *_MockAthVenueServiceRecorder) CreateVenueHoliday(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateVenueHoliday", arg0, arg1, arg2)
}

func (_m *MockAthVenueService) CreateVenueImage(_param0 context.Context, _param1 *gorm.DB, _param2 models.VenueImagesReq) models.Response {
	ret := _m.ctrl.Call(_m, "CreateVenueImage", _param0, _param1, _param2)
	ret0, _ := ret[0].(models.Response)
	return ret0
}

func (_mr *_MockAthVenueServiceRecorder) CreateVenueImage(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateVenueImage", arg0, arg1, arg2)
}

func (_m *MockAthVenueService) EditVenue(_param0 context.Context, _param1 *gorm.DB, _param2 models.AthVenues) models.Response {
	ret := _m.ctrl.Call(_m, "EditVenue", _param0, _param1, _param2)
	ret0, _ := ret[0].(models.Response)
	return ret0
}

func (_mr *_MockAthVenueServiceRecorder) EditVenue(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "EditVenue", arg0, arg1, arg2)
}

func (_m *MockAthVenueService) GetAllAmenities(_param0 context.Context, _param1 *gorm.DB, _param2 v1.GetAllAmenitiesReq) models.Response {
	ret := _m.ctrl.Call(_m, "GetAllAmenities", _param0, _param1, _param2)
	ret0, _ := ret[0].(models.Response)
	return ret0
}

func (_mr *_MockAthVenueServiceRecorder) GetAllAmenities(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetAllAmenities", arg0, arg1, arg2)
}

func (_m *MockAthVenueService) GetAmenitiesByID(_param0 context.Context, _param1 *gorm.DB, _param2 v1.GetAmenitiesByIDReq) models.Response {
	ret := _m.ctrl.Call(_m, "GetAmenitiesByID", _param0, _param1, _param2)
	ret0, _ := ret[0].(models.Response)
	return ret0
}

func (_mr *_MockAthVenueServiceRecorder) GetAmenitiesByID(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetAmenitiesByID", arg0, arg1, arg2)
}

func (_m *MockAthVenueService) GetAmenitiesForVenueByID(_param0 context.Context, _param1 *gorm.DB, _param2 int) models.Response {
	ret := _m.ctrl.Call(_m, "GetAmenitiesForVenueByID", _param0, _param1, _param2)
	ret0, _ := ret[0].(models.Response)
	return ret0
}

func (_mr *_MockAthVenueServiceRecorder) GetAmenitiesForVenueByID(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetAmenitiesForVenueByID", arg0, arg1, arg2)
}

func (_m *MockAthVenueService) GetListOfVenueByMerchantID(_param0 context.Context, _param1 *gorm.DB, _param2 v1.GetListOfVenueByMerchantIDReq) models.Response {
	ret := _m.ctrl.Call(_m, "GetListOfVenueByMerchantID", _param0, _param1, _param2)
	ret0, _ := ret[0].(models.Response)
	return ret0
}

func (_mr *_MockAthVenueServiceRecorder) GetListOfVenueByMerchantID(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetListOfVenueByMerchantID", arg0, arg1, arg2)
}

func (_m *MockAthVenueService) GetVenueByID(_param0 context.Context, _param1 *gorm.DB, _param2 int) models.Response {
	ret := _m.ctrl.Call(_m, "GetVenueByID", _param0, _param1, _param2)
	ret0, _ := ret[0].(models.Response)
	return ret0
}

func (_mr *_MockAthVenueServiceRecorder) GetVenueByID(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetVenueByID", arg0, arg1, arg2)
}

func (_m *MockAthVenueService) SaveHoursOfOperation(_param0 context.Context, _param1 *gorm.DB, _param2 models.AthVenueHours) models.Response {
	ret := _m.ctrl.Call(_m, "SaveHoursOfOperation", _param0, _param1, _param2)
	ret0, _ := ret[0].(models.Response)
	return ret0
}

func (_mr *_MockAthVenueServiceRecorder) SaveHoursOfOperation(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SaveHoursOfOperation", arg0, arg1, arg2)
}