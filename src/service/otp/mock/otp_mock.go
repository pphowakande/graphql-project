// Automatically generated by MockGen. DO NOT EDIT!
// Source: api/src/service/otp (interfaces: AthOtpService)

package mock

import (
	models "api/src/models"
	context "context"
	gomock "github.com/golang/mock/gomock"
	gorm "github.com/jinzhu/gorm"
)

// Mock of AthOtpService interface
type MockAthOtpService struct {
	ctrl     *gomock.Controller
	recorder *_MockAthOtpServiceRecorder
}

// Recorder for MockAthOtpService (not exported)
type _MockAthOtpServiceRecorder struct {
	mock *MockAthOtpService
}

func NewMockAthOtpService(ctrl *gomock.Controller) *MockAthOtpService {
	mock := &MockAthOtpService{ctrl: ctrl}
	mock.recorder = &_MockAthOtpServiceRecorder{mock}
	return mock
}

func (_m *MockAthOtpService) EXPECT() *_MockAthOtpServiceRecorder {
	return _m.recorder
}

func (_m *MockAthOtpService) CreateCode(_param0 context.Context, _param1 *gorm.DB, _param2 models.AthUserOTP) models.Response {
	ret := _m.ctrl.Call(_m, "CreateCode", _param0, _param1, _param2)
	ret0, _ := ret[0].(models.Response)
	return ret0
}

func (_mr *_MockAthOtpServiceRecorder) CreateCode(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateCode", arg0, arg1, arg2)
}

func (_m *MockAthOtpService) CreateOTP(_param0 context.Context, _param1 *gorm.DB, _param2 *gorm.DB, _param3 models.AthUserOTP) models.Response {
	ret := _m.ctrl.Call(_m, "CreateOTP", _param0, _param1, _param2, _param3)
	ret0, _ := ret[0].(models.Response)
	return ret0
}

func (_mr *_MockAthOtpServiceRecorder) CreateOTP(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateOTP", arg0, arg1, arg2, arg3)
}

func (_m *MockAthOtpService) EmailSend(_param0 context.Context, _param1 *gorm.DB, _param2 *gorm.DB, _param3 models.EmailSendReq) models.Response {
	ret := _m.ctrl.Call(_m, "EmailSend", _param0, _param1, _param2, _param3)
	ret0, _ := ret[0].(models.Response)
	return ret0
}

func (_mr *_MockAthOtpServiceRecorder) EmailSend(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "EmailSend", arg0, arg1, arg2, arg3)
}

func (_m *MockAthOtpService) VerifyOTP(_param0 context.Context, _param1 *gorm.DB, _param2 *gorm.DB, _param3 *gorm.DB, _param4 models.OTPVerify) models.Response {
	ret := _m.ctrl.Call(_m, "VerifyOTP", _param0, _param1, _param2, _param3, _param4)
	ret0, _ := ret[0].(models.Response)
	return ret0
}

func (_mr *_MockAthOtpServiceRecorder) VerifyOTP(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "VerifyOTP", arg0, arg1, arg2, arg3, arg4)
}

func (_m *MockAthOtpService) VerifyResetPasswordOTP(_param0 context.Context, _param1 models.OTPVerify) models.Response {
	ret := _m.ctrl.Call(_m, "VerifyResetPasswordOTP", _param0, _param1)
	ret0, _ := ret[0].(models.Response)
	return ret0
}

func (_mr *_MockAthOtpServiceRecorder) VerifyResetPasswordOTP(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "VerifyResetPasswordOTP", arg0, arg1)
}
