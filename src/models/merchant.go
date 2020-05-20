package io

type AthMerchant struct {
	MerchantID   int    `gorm:"column:merchantId;primary_key" json:"merchantId"`
	MerchantName string `gorm:"column:merchantName" json:"merchantName"`
	Address      string `gorm:"column:address" json:"address"`
	Phone        string `gorm:"column:phone" json:"phone"`
	Email        string `gorm:"column:email" json:"email"`
	GstNoFile    string `gorm:"column:gstNoFile" json:"gstNoFile"`
	PanNoFile    string `gorm:"column:panNoFile" json:"panNoFile"`
	BankAccFile  string `gorm:"column:bankAccFile" json:"bankAccFile"`
	AddressFile  string `gorm:"column:addressFile" json:"addressFile"`

	GstFileVerified  bool `gorm:"column:gstFileVerified" json:"gstFileVerified"`
	AddFileVerified  bool `gorm:"column:addFileVerified" json:"addFileVerified"`
	BankFileVerified bool `gorm:"column:bankFileVerified" json:"bankFileVerified"`
	PanFileVerified  bool `gorm:"column:panFileVerified" json:"panFileVerified"`
	IsActive         bool `gorm:"column:isActive" json:"isActive"`
	DoLater          bool `gorm:"column:doLater" json:"doLater"`
	Models
}

func (AthMerchant) TableName() string {
	return "athMerchants"
}

type AthMerchantUser struct {
	MerchantUserID int `gorm:"column:merchantUserId;primary_key" json:"merchantUserId"`
	MerchantID     int `gorm:"column:merchantId" json:"merchantId"`
	UserID         int `gorm:"column:userId" json:"userId"`
	CreatedAt      int `gorm:"column:createdAt" json:"createdAt,omitempty"`
	CreatedBy      int `gorm:"column:createdBy" json:"createdBy"`
	DeletedAt      int `gorm:"column:deletedAt" json:"deletedAt,omitempty"`
	DeletedBy      int `gorm:"column:deletedBy" json:"deletedBy"`
}

func (AthMerchantUser) TableName() string {
	return "athMerchantUsers"
}

type AthUser struct {
	UserID              int    `gorm:"column:userId;primary_key" json:"userId"`
	Email               string `gorm:"column:email;unique;not null" json:"email"`
	Password            string `gorm:"column:password" json:"password"`
	FirstName           string `gorm:"column:firstName" json:"firstName"`
	LastName            string `gorm:"column:lastName" json:"lastName"`
	Phone               string `gorm:"column:phone" json:"phone"`
	Gender              string `gorm:"column:gender" json:"gender"`
	DOB                 string `gorm:"column:dob" json:"dob"`
	ProfileImage        string `gorm:"column:profileImage" json:"profileImage"`
	LastPasswordResetAt int    `gorm:"column:lastPasswordResetAt" json:"lastPasswordResetAt"`
	UserSource          string `gorm:"column:userSource" json:"userSource"`
	LastLoginAt         int    `gorm:"column:lastLoginAt" json:"lastLoginAt"`
	LastLoginIP         string `gorm:"column:lastLoginIP" json:"lastLoginIP"`
	IsActive            bool   `gorm:"column:isActive" json:"isActive"`
	EmailVerify         bool   `gorm:"column:emailVerify" json:"emailVerify"`
	PhoneVerify         bool   `gorm:"column:phoneVerify" json:"phoneVerify"`
	BmsActive           bool   `gorm:"column:bmsActive" json:"bmsActive"`
	TokenValue          string `gorm:"-" json:"-"`
	TokenHash           string `gorm:"column:tokenHash" json:"tokenHash"`
	TokenExpireAt       int    `gorm:"column:tokenExpireAt" json:"tokenExpireAt"`
	Models
}

// set User's table name
func (AthUser) TableName() string {
	return "athUsers"
}

type OTPVerify struct {
	VerificationType string `json:"verificationType"`
	VerificationCode string `json:"verificationCode"`
	AccountId        int32  `json:"accountId"`
	Email            string `json:"email"`
}
type ValidateUser struct {
	EmailVerify bool `gorm:"column:emailVerify" json:"emailVerify"`
	PhoneVerify bool `gorm:"column:phoneVerify" json:"phoneVerify"`
	IsActive    bool `gorm:"column:isActive" json:"isActive"`
	UserID      int  `gorm:"column:userId" json:"userId"`
}

type EmailSendReq struct {
	Email            string `json:"email"`
	Account          string `json:"account"`
	Subject          string `json:"subject"`
	Body             string `json:"body"`
	Eticket          string `json:"eTicket"`
	Tid              string `json:"tid"`
	EmailVerifyToken string `json:"emailVerifyToken"`
	OTPTtype         string `json:"otpType"`
}

type TeamMemberReq struct {
	AccountType string  `json:"accountType"`
	UserID      int     `json:"userId"`
	AccessData  []int32 `json:"accessData"`
	CreatedBy   int     `json:"createdBy"`
	UpdatedBy   int     `json:"updatedBy"`
}

type MerchantKYC struct {
	GstNoFile   string `json:"gstNoFile"`
	PanNoFile   string `json:"panNoFile"`
	AddressFile string `json:"addressFile"`
	BankAccFile string `json:"bankAccFile"`
	DoLater     bool   `json:"doLater"`
	MerchantID  int    `json:"merchantId"`
}

type MemberData struct {
	AccountType string `json:"accountType"`
	UserID      int    `json:"userId"`
	VenueID     int32  `json:"venueId"`
	CreatedBy   int    `json:"createdBy"`
}

type AthVenueUser struct {
	VenueUserId int  `gorm:"column:venueUserId;primary_key" json:"venueUserId"`
	VenueId     int  `gorm:"column:venueId" json:"venueId"`
	UserId      int  `gorm:"column:userId" json:"userId"`
	RoleId      int  `gorm:"column:roleId" json:"roleId"`
	IsActive    bool `gorm:"column:isActive" json:"isActive"`
	Models
}

// set User's table name
func (AthVenueUser) TableName() string {
	return "athVenueUsers"
}

type AthRole struct {
	RoleId   int    `gorm:"column:roleId" json:"roleId"`
	RoleName string `gorm:"column:roleName" json:"roleName"`
	IsActive bool   `gorm:"column:isActive" json:"isActive"`
	Models
}

// set User's table name
func (AthRole) TableName() string {
	return "athRoles"
}

type TeamData struct {
	AccountId   int    `gorm:"column:accountId" json:"accountId"`
	FullName    string `gorm:"column:fullName" json:"fullName"`
	AccountType string `gorm:"column:accountType" json:"accountType"`
	Email       string `gorm:"column:email" json:"email"`
	Phone       string `gorm:"column:phone" json:"phone"`
	CreatedAt   int32  `gorm:"column:createdAt" json:"createdAt"`
	CreatedBy   int32  `gorm:"column:createdBy" json:"createdBy"`
}

type MerchantDocFileData struct {
	ContentType     string `gorm:"column:contentType" json:"contentType"`
	ValidatedBase64 string `gorm:"column:validatedBase64" json:"validatedBase64"`
	Blob            string `gorm:"column:blob" json:"blob"`
	Verified        bool   `gorm:"column:verified" json:"verified"`
	Exists          bool   `gorm:"column:exists" json:"exists"`
	FileName        string `gorm:"column:fileName" json:"fileName"`
}
