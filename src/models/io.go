package io

type Models struct {
	CreatedAt int `gorm:"column:createdAt" json:"createdAt,omitempty"`
	UpdatedAt int `gorm:"column:updatedAt" json:"updatedAt,omitempty"`
	//DeletedAt int `gorm:"column:deletedAt" json:"deletedAt,omitempty"`
	CreatedBy int `gorm:"column:createdBy" json:"createdBy"`
	UpdatedBy int `gorm:"column:updatedBy" json:"updatedBy"`
	//DeletedBy int `gorm:"column:deletedBy" json:"deletedBy"`
}

type Response struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Success bool        `json:"success" `
	Error   error       `json:"error"`
}

type LoginRequest struct {
	Login      string `gorm:"column:login" json:"login"`
	Password   string `gorm:"column:password" json:"password"`
	EmailParam bool   `gorm:"column:emailParam" json:"emailParam"`
	PhoneParam bool   `gorm:"column:phoneParam" json:"phoneParam"`
}

type ForgotPasswordRequest struct {
	Email string `gorm:"column:email" json:"email"`
}

type ResetPasswordRequest struct {
	Email    string `gorm:"column:email" json:"email"`
	Password string `gorm:"column:password" json:"password"`
	Token    string `gorm:"column:token" json:"token"`
}

type OtpRequest struct {
	OtpNo     string `gorm:"column:otpNo" json:"otpNo"`
	OtpType   string `gorm:"column:otpType" json:"otpType"`
	OtpExpiry string `gorm:"column:otpExpiry" json:"otpExpiry"`
	UserID    string `gorm:"column:userID" json:"userID"`
}

func SuccessMessage(data interface{}, msg ...string) Response {
	newMessage := "Success"
	if len(msg) > 0 {
		newMessage = msg[0]
	}
	return Response{
		Data:    data,
		Success: true,
		Message: newMessage,
	}
}
func FailureMessage(err error, msg ...string) Response {
	newMessage := "Failure"
	if len(msg) > 0 {
		newMessage = msg[0]
	}
	return Response{
		Success: false,
		Message: newMessage,
		Error:   err,
	}
}
