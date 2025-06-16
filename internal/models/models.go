package models

import "time"

// Booking represents the bookings table
type Booking struct {
	ID         int           `gorm:"primaryKey;column:id" json:"id"`
	StatusCode int           `gorm:"column:status_code" json:"status_code"`
	StartDate  time.Time     `gorm:"column:start_date" json:"start_date"`
	EndDate    time.Time     `gorm:"column:end_date" json:"end_date"`
	CheckIn    *time.Time    `gorm:"column:check_in" json:"check_in"`
	CheckOut   *time.Time    `gorm:"column:check_out" json:"check_out"`
	BabyBed    bool          `gorm:"column:baby_bed" json:"baby_bed"`
	BookingSum float64       `gorm:"column:booking_sum" json:"booking_sum"`
	DiscountID int           `gorm:"column:discount_id" json:"discount_id"`
	TotalSum   float64       `gorm:"column:total_sum" json:"total_sum"`
	Status     BookingStatus `gorm:"foreignKey:StatusCode;references:StatusCode" json:"status"`
	Discount   Discount      `gorm:"foreignKey:DiscountID;references:ID" json:"discount"`
	Complaints []Complaint   `gorm:"foreignKey:BookingID" json:"complaints"`
	Guests     []Guest       `gorm:"many2many:guests_in_bookings;foreignKey:ID;joinForeignKey:BookingID;references:ID;joinReferences:GuestID" json:"guests"`
	Payments   []Payment     `gorm:"foreignKey:BookingID" json:"payments"`
}

// BookingStatus represents the booking_statuses table
type BookingStatus struct {
	StatusCode int    `gorm:"primaryKey;column:status_code" json:"status_code"`
	Name       string `gorm:"type:varchar(15);column:name" json:"name"`
}

// Complaint represents the complaints table
type Complaint struct {
	ID         int             `gorm:"primaryKey;column:id" json:"id"`
	Reason     string          `gorm:"type:varchar(1000);column:reason" json:"reason"`
	Commentary *string         `gorm:"type:varchar(1000);column:commentary" json:"commentary"`
	IssueDate  time.Time       `gorm:"column:issue_date" json:"issue_date"`
	BookingID  int             `gorm:"column:booking_id" json:"booking_id"`
	StatusCode int             `gorm:"column:status_code" json:"status_code"`
	Booking    Booking         `gorm:"foreignKey:BookingID;references:ID" json:"booking"`
	Status     ComplaintStatus `gorm:"foreignKey:StatusCode;references:StatusCode" json:"status"`
}

// ComplaintStatus represents the complaints_statuses table
type ComplaintStatus struct {
	StatusCode int    `gorm:"primaryKey;column:status_code" json:"status_code"`
	Name       string `gorm:"type:varchar(50);column:name" json:"name"`
}

// Discount represents the discounts table
type Discount struct {
	ID        int `gorm:"primaryKey;column:id" json:"id"`
	MinNights int `gorm:"column:min_nights" json:"min_nights"`
	Amount    int `gorm:"column:amount" json:"amount"`
}

// Guest represents the guests table
type Guest struct {
	ID          int       `gorm:"primaryKey;column:id" json:"id"`
	Name        string    `gorm:"type:varchar(100);column:name" json:"name"`
	PhoneNumber string    `gorm:"type:varchar(20);column:phone_number" json:"phone_number"`
	PassportNo  string    `gorm:"type:varchar(20);column:passport_no" json:"passport_no"`
	Bookings    []Booking `gorm:"many2many:guests_in_bookings;foreignKey:ID;joinForeignKey:GuestID;references:ID;joinReferences:BookingID" json:"bookings"`
}

// Payment represents the payments table
type Payment struct {
	ID         int           `gorm:"primaryKey;column:id" json:"id"`
	BookingID  int           `gorm:"column:booking_id" json:"booking_id"`
	PayDate    time.Time     `gorm:"column:pay_date" json:"pay_date"`
	Amount     float64       `gorm:"column:amount" json:"amount"`
	MethodCode int           `gorm:"column:method_code" json:"method_code"`
	StatusCode int           `gorm:"column:status_code" json:"status_code"`
	Booking    Booking       `gorm:"foreignKey:BookingID;references:ID" json:"booking"`
	Method     PaymentMethod `gorm:"foreignKey:MethodCode;references:Code" json:"method"`
	Status     PaymentStatus `gorm:"foreignKey:StatusCode;references:StatusCode" json:"status"`
}

// PaymentMethod represents the payment_methods table
type PaymentMethod struct {
	Code int    `gorm:"primaryKey;column:code" json:"code"`
	Name string `gorm:"type:varchar(30);column:name" json:"name"`
}

// PaymentStatus represents the payment_statuses table
type PaymentStatus struct {
	StatusCode int    `gorm:"primaryKey;column:status_code" json:"status_code"`
	Name       string `gorm:"type:varchar(30);column:name" json:"name"`
}

// Room represents the rooms table
type Room struct {
	Number       int          `gorm:"primaryKey;column:number" json:"number"`
	CategoryCode int          `gorm:"column:category_code" json:"category_code"`
	StateCode    int          `gorm:"column:state_code" json:"state_code"`
	Category     RoomCategory `gorm:"foreignKey:CategoryCode;references:Code" json:"category"`
	State        RoomState    `gorm:"foreignKey:StateCode;references:StateCode" json:"state"`
}

// RoomCategory represents the room_categories table
type RoomCategory struct {
	Code     int    `gorm:"primaryKey;column:code" json:"code"`
	Name     string `gorm:"type:varchar(15);column:name" json:"name"`
	Capacity int    `gorm:"column:capacity" json:"capacity"`
}

// RoomState represents the room_states table
type RoomState struct {
	StateCode int    `gorm:"primaryKey;column:state_code" json:"state_code"`
	Name      string `gorm:"type:varchar(20);column:name" json:"name"`
}

// Tariff represents the tariffs table
type Tariff struct {
	Code         int               `gorm:"primaryKey;column:code" json:"code"`
	CategoryCode int               `gorm:"column:category_code" json:"category_code"`
	BasePrice    float64           `gorm:"column:base_price" json:"base_price"`
	DayCode      int               `gorm:"column:day_code" json:"day_code"`
	Category     RoomCategory      `gorm:"foreignKey:CategoryCode;references:Code" json:"category"`
	Coefficient  TariffCoefficient `gorm:"foreignKey:DayCode;references:DayCode" json:"coefficient"`
}

// TariffCoefficient represents the tariff_coefficients table
type TariffCoefficient struct {
	DayCode     int     `gorm:"primaryKey;column:day_code" json:"day_code"`
	Coefficient float64 `gorm:"column:coefficient" json:"coefficient"`
}

// Holiday represents the holidays table
type Holiday struct {
	HolidayDate time.Time `gorm:"primaryKey;column:holiday_date" json:"holiday_date"`
	Name        string    `gorm:"type:varchar(100);column:name" json:"name"`
}
