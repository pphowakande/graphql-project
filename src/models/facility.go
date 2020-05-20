package io

type FacilityTotalBookingData struct {
	NoOfBookings int32 `gorm:"column:noOfBookings" json:"noOfBookings"`
	FacilityID   int32 `gorm:"column:facilityId" json:"facilityId"`
}

type FacilityLastBookingData struct {
	LastBookingDate string `gorm:"column:lastBookingDate" json:"lastBookingDate"`
	FacilityID      int32  `gorm:"column:facilityId" json:"facilityId"`
}

type FacilityTotalEarnings struct {
	TotalEarnings float32 `gorm:"column:totalEarnings" json:"totalEarnings"`
	FacilityID    int32   `gorm:"column:facilityId" json:"facilityId"`
}

type FacilityTodayStats struct {
	FacilityID        int32   `gorm:"column:facilityId" json:"facilityId"`
	NoOfBookingsToday int32   `gorm:"column:noOfBookingsToday" json:"noOfBookingsToday"`
	EarnedToday       float32 `gorm:"column:earnedToday" json:"earnedToday"`
}

type DeleteFacilityByID struct {
	FacilityID int32  `gorm:"column:facilityID" json:"facilityID"`
	Type       string `gorm:"column:type" json:"type"`
	UserID     int32  `gorm:"column:userId" json:"userId"`
}

type CustomRatesForFacilityByID struct {
	FacilityID int32  `gorm:"column:facilityId" json:"facilityId"`
	FromDate   string `gorm:"column:fromDate" json:"fromDate"`
	ToDate     string `gorm:"column:toDate" json:"toDate"`
}

type AthFacilities struct {
	FacilityID        int     `gorm:"column:facilityId;primary_key" json:"facilityId"`
	VenueID           int     `gorm:"column:venueId;" json:"venueId"`
	FacilityName      string  `gorm:"column:facilityName" json:"facilityName"`
	FacilityBasePrice float32 `gorm:"column:facilityBasePrice" json:"facilityBasePrice"`
	TimeSlot          int     `gorm:"column:timeSlot" json:"timeSlot"`
	IsActive          bool    `gorm:"column:isActive" json:"isActive"`
	Models
}

func (AthFacilities) TableName() string {
	return "athFacilities"
}

type AthFacilityBookings struct {
	BookingID       int     `gorm:"column:bookingID;primary_key" json:"bookingID"`
	UserID          int     `gorm:"column:userID;" json:"userID"`
	FacilityID      int     `gorm:"column:facilityID;" json:"facilityID"`
	BookingNo       string  `gorm:"column:bookingNo" json:"bookingNo"`
	BookingDate     int     `gorm:"column:bookingDate" json:"bookingDate"`
	BaseTotalAmount float32 `gorm:"column:baseTotalAmount" json:"baseTotalAmount"`
	DiscountAmount  float32 `gorm:"column:discountAmount" json:"discountAmount"`
	BookingAmount   float32 `gorm:"column:bookingAmount" json:"bookingAmount"`
	BookingFee      float32 `gorm:"column:bookingFee" json:"bookingFee"`
	DiscountID      int     `gorm:"column:discountID" json:"discountID"`
	Models
}

func (AthFacilityBookings) TableName() string {
	return "athFacilityBookings"
}

type AthFacilitySlots struct {
	FacilitySlotID int     `gorm:"column:facilitySlotId;primary_key" json:"facilitySlotId"`
	FacilityID     int     `gorm:"column:facilityId;" json:"facilityId"`
	UserID         int     `gorm:"column:userId" json:"userId"`
	SlotDays       string  `gorm:"column:slotDays" json:"slotDays"`
	SlotType       string  `gorm:"column:slotType" json:"slotType"`
	SlotFromTime   string  `gorm:"column:slotFromTime" json:"slotFromTime"`
	SlotToTime     string  `gorm:"column:slotTotime" json:"slotTotime"`
	SlotPrice      float32 `gorm:"column:slotPrice" json:"slotPrice"`
	IsActive       bool    `gorm:"column:isActive" json:"isActive"`
	Models
}

func (AthFacilitySlots) TableName() string {
	return "athFacilitySlots"
}

type AthSportCategories struct {
	SportCategoryID int    `gorm:"column:sportCategoryId;primary_key" json:"sportCategoryId"`
	CategoryName    string `gorm:"column:categoryName;" json:"categoryName"`
	IsActive        bool   `gorm:"column:isActive" json:"isActive"`
	Models
}

func (AthSportCategories) TableName() string {
	return "athSportCategories"
}

type AthFacilityCustomRates struct {
	RateID       int     `gorm:"column:rateId;primary_key" json:"rateId"`
	FacilityID   int     `gorm:"column:facilityId;" json:"facilityId"`
	UserID       int     `gorm:"column:userId" json:"userId"`
	SlotFromTime string  `gorm:"column:slotFromTime" json:"slotFromTime"`
	SlotToTime   string  `gorm:"column:slotToTime" json:"slotToTime"`
	SlotPrice    float32 `gorm:"column:slotPrice" json:"slotPrice"`
	IsActive     bool    `gorm:"column:isActive" json:"isActive"`
	Available    bool    `gorm:"column:available" json:"available"`
	Models
}

func (AthFacilityCustomRates) TableName() string {
	return "athFacilityCustomRates"
}

type AthFacilitySports struct {
	SportCategoryID int `gorm:"column:sportCategoryId;" json:"sportCategoryId"`
	FacilityID      int `gorm:"column:facilityId;" json:"facilityId"`
}

func (AthFacilitySports) TableName() string {
	return "athFacilitySports"
}

type AthFacilitySportsData struct {
	SportCategoryID int  `json:"sportCategoryId"`
	FacilityID      int  `json:"facilityId"`
	Status          bool `json:"status"`
}
