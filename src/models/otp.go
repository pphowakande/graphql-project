package io

/*

	create table turfdb_transactional.athUserOtps (
		otpId	integer(11)	auto_increment primary key not null,
		userId	int(11) ,
		otpNo	varchar(100)	not	NULL,
		otpType	varchar(100)	not	NULL,
		createdAt	int(11)	default	NULL,
		updatedAt	int(11)	default	NULL,
		expiredAt	int(11)	default	NULL,
		isActive	tinyint(1)	default	1,

		CONSTRAINT FK_ath_user_otps FOREIGN KEY (userId) REFERENCES turfdb_masters.athUsers (userId)
	);

*/
type AthUserOTP struct {
	OtpID     int    `gorm:"column:otpId;primary_key" json:"otpId"`
	UserID    int    `gorm:"column:userId;not null" json:"userId"`
	OTPNO     string `gorm:"column:otpNo" json:"otpNo"`
	Phone     string `gorm:"column:phone" json:"phone"`
	Email     string `gorm:"column:email" json:"email"`
	OTPType   string `gorm:"column:otpType" json:"otpType"`
	CreatedAt int    `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt int    `gorm:"column:updatedAt" json:"updatedAt"`
	ExpiredAt int    `gorm:"column:expiredAt" json:"expiredAt"`
	IsActive  bool   `gorm:"column:isActive" json:"isActive"`
}

func (AthUserOTP) TableName() string {
	return "athUserOtps"
}
