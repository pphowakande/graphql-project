package io

type AthVenues struct {
	VenueID     int     `gorm:"column:venueId;primary_key" json:venueId"`
	MerchantID  int     `gorm:"column:merchantId;" json:"merchantId"`
	VenueName   string  `gorm:"column:venueName" json:"venueName"`
	Email       string  `gorm:"column:email" json:"email"`
	Phone       string  `gorm:"column:phone" json:"phone"`
	Address     string  `gorm:"column:address" json:"address"`
	Latitude    float32 `gorm:"column:latitude" json:"latitude"`
	Longitude   float32 `gorm:"column:longitude" json:"longitude"`
	Description string  `gorm:"column:description" json:"description"`
	IsActive    bool    `gorm:"column:isActive" json:"isActive"`
	Models
}

// set User's table name
func (AthVenues) TableName() string {
	return "athVenues"
}

type AthVenueImages struct {
	VenueImageID int    `gorm:"column:venueImageId;primary_key" json:"venueImageId"`
	VenueID      int    `gorm:"column:venueId;" json:"venueId"`
	ImageTitle   string `gorm:"column:imageTitle" json:"imageTitle"`
	ImageUrl     string `gorm:"column:imageUrl" json:"imageUrl"`
	ImageType    string `gorm:"column:imageType" json:"imageType"`
	IsActive     bool   `gorm:"column:isActive" json:"isActive"`
	Models
}

// set User's table name
func (AthVenueImages) TableName() string {
	return "athVenueImages"
}

type AthVenueHours struct {
	HourID      int    `gorm:"column:hourId;primary_key" json:"hourId"`
	VenueID     int    `gorm:"column:venueId;" json:"venueId"`
	Day         string `gorm:"column:day" json:"day"`
	OpeningTime string `gorm:"column:openingTime" json:"openingTime"`
	ClosingTime string `gorm:"column:closingTime" json:"closingTime"`
	IsActive    bool   `gorm:"column:isActive" json:"isActive"`
	Models
}

// set User's table name
func (AthVenueHours) TableName() string {
	return "athVenueHours"
}

type AthVenueHolidays struct {
	HolidayID int    `gorm:"column:holidayId;primary_key" json:"holidayId"`
	VenueID   int    `gorm:"column:venueId;" json:"venueId"`
	Title     string `gorm:"column:title" json:"title"`
	Date      string `gorm:"column:date" json:"date"`
	IsActive  bool   `gorm:"column:isActive" json:"isActive"`
	Models
}

// set User's table name
func (AthVenueHolidays) TableName() string {
	return "athVenueHolidays"
}

type AthAmenities struct {
	AmenityID   int    `gorm:"column:amenityId;primary_key" json:"amenityId"`
	AmenityName string `gorm:"column:amenityName;" json:"amenityName"`
	IsActive    bool   `gorm:"column:isActive" json:"isActive"`
	Models
}

// set User's table name
func (AthAmenities) TableName() string {
	return "athAmenities"
}

type AthVenueAmenities struct {
	AmenityID int `gorm:"column:amenityId;" json:"amenityId"`
	VenueID   int `gorm:"column:venueId;" json:"venueId"`
}

// set User's table name
func (AthVenueAmenities) TableName() string {
	return "athVenueAmenities"
}

type AthVenueAmenitiesData struct {
	AmenityID int  `json:"amenityId"`
	VenueID   int  `json:"venueId"`
	Status    bool `json:"status"`
}

type VenueImageData struct {
	ContentType     string `gorm:"column:contentType" json:"contentType"`
	ValidatedBase64 string `gorm:"column:validatedBase64" json:"validatedBase64"`
	Blob            string `gorm:"column:blob" json:"blob"`
	Verified        bool   `gorm:"column:verified" json:"verified"`
	Exists          bool   `gorm:"column:exists" json:"exists"`
	FileName        string `gorm:"column:fileName" json:"fileName"`
	ImageUrl        string `gorm:"column:imageUrl" json:"imageUrl"`
}

type VenueImagesReq struct {
	HeaderImg    []ImgData `json:"headerImg"`
	ThumbnailImg []ImgData `json:"thumbnailImg"`
	GalleryImg   []ImgData `json:"galleryImg"`
	VenueID      int       `json:"venueId"`
	CreatedBy    int       `json:"createdBy"`
}

type ImgData struct {
	Image  string `json:"image`
	Status bool   `json:"status"`
}
