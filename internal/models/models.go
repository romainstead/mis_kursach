package models

import "time"

// Booking represents the bookings table
type Booking struct {
	ID         int           `json:"id" db:"id"`
	StatusCode int           `json:"status_code" db:"status_code"`
	StartDate  time.Time     `json:"start_date" db:"start_date"`
	EndDate    time.Time     `json:"end_date" db:"end_date"`
	CheckIn    *time.Time    `json:"check_in" db:"check_in"`
	CheckOut   *time.Time    `json:"check_out" db:"check_out"`
	BabyBed    bool          `json:"baby_bed" db:"baby_bed"`
	BookingSum float64       `json:"booking_sum" db:"booking_sum"`
	DiscountID int           `json:"discount_id" db:"discount_id"`
	TotalSum   float64       `json:"total_sum" db:"total_sum"`
	Status     BookingStatus `json:"status" db:"status"`
	Discount   Discount      `json:"discount" db:"discount"`
	Complaints []Complaint   `json:"complaints" db:"complaints"`
	Guests     []Guest       `json:"guests" db:"guests"`
	Payments   []Payment     `json:"payments" db:"payments"`
}

type BookingResponse struct {
	ID             int        `json:"id"`
	StartDate      time.Time  `json:"start_date" db:"start_date"`
	EndDate        time.Time  `json:"end_date" db:"end_date"`
	CheckIn        *time.Time `json:"check_in" db:"check_in"`
	CheckOut       *time.Time `json:"check_out" db:"check_out"`
	BabyBed        bool       `json:"baby_bed" db:"baby_bed"`
	BookingSum     float64    `json:"booking_sum" db:"booking_sum"`
	TotalSum       float64    `json:"total_sum" db:"total_sum"`
	BookingStatus  string     `json:"booking_status"`
	DiscountAmount float64    `json:"discount_amount"`
	Room           int        `json:"room"`
}

// BookingStatus represents the booking_statuses table
type BookingStatus struct {
	StatusCode int    `json:"status_code" db:"status_code"`
	Name       string `json:"name" db:"name"`
}

// Complaint represents the complaints table
type Complaint struct {
	ID         int             `json:"id"`
	Reason     string          `json:"reason"`
	Commentary *string         `json:"commentary"`
	IssueDate  time.Time       `json:"issue_date"`
	BookingID  int             `json:"booking_id"`
	StatusCode int             `json:"status_code"`
	Booking    Booking         `json:"booking"`
	Status     ComplaintStatus `json:"status"`
}

type ComplaintResponse struct {
	ID         int       `json:"id"`
	Reason     string    `json:"reason"`
	Commentary *string   `json:"commentary"`
	IssueDate  time.Time `json:"issue_date"`
	BookingID  int       `json:"booking_id"`
	Status     string    `json:"status"`
	Room       int       `json:"room"`
}

// ComplaintStatus represents the complaints_statuses table
type ComplaintStatus struct {
	StatusCode int    `json:"status_code"`
	Name       string `json:"name"`
}

// Discount represents the discounts table
type Discount struct {
	ID        int     `json:"id" db:"id"`
	MinNights int     `json:"min_nights" db:"min_nights"`
	Amount    float64 `json:"amount" db:"amount"`
}

// Guest represents the guests table
type Guest struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	PhoneNumber string    `json:"phone_number"`
	PassportNo  string    `json:"passport_no"`
	Bookings    []Booking `json:"bookings"`
}

// Payment represents the payments table
type Payment struct {
	ID         int           `json:"id"`
	BookingID  int           `json:"booking_id"`
	PayDate    time.Time     `json:"pay_date"`
	Amount     float64       `json:"amount"`
	MethodCode int           `json:"method_code"`
	StatusCode int           `json:"status_code"`
	Booking    Booking       `json:"booking"`
	Method     PaymentMethod `json:"method"`
	Status     PaymentStatus `json:"status"`
}

type PaymentResponse struct {
	ID         int       `json:"id"`
	BookingID  int       `json:"booking_id"`
	PayDate    time.Time `json:"pay_date"`
	Amount     float64   `json:"amount"`
	MethodName string    `json:"method_name"`
	StatusName string    `json:"status_name"`
}

// PaymentMethod represents the payment_methods table
type PaymentMethod struct {
	Code int    `json:"code"`
	Name string `json:"name"`
}

// PaymentStatus represents the payment_statuses table
type PaymentStatus struct {
	StatusCode int    `json:"status_code"`
	Name       string `json:"name"`
}

// Room represents the rooms table
type Room struct {
	Number       int          `json:"number"`
	CategoryCode int          `json:"category_code"`
	StateCode    int          `json:"state_code"`
	Category     RoomCategory `json:"category"`
	State        RoomState    `json:"state"`
}

// RoomCategory represents the room_categories table
type RoomCategory struct {
	Code     int    `json:"code"`
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
}

// RoomState represents the room_states table
type RoomState struct {
	StateCode int    `json:"state_code"`
	Name      string `json:"name"`
}

// Tariff represents the tariffs table
type Tariff struct {
	Code         int               `json:"code"`
	CategoryCode int               `json:"category_code"`
	BasePrice    float64           `json:"base_price"`
	DayCode      int               `json:"day_code"`
	Category     RoomCategory      `json:"category"`
	Coefficient  TariffCoefficient `json:"coefficient"`
}

// TariffCoefficient represents the tariff_coefficients table
type TariffCoefficient struct {
	DayCode     int     `json:"day_code"`
	Coefficient float64 `json:"coefficient"`
}

// Holiday represents the holidays table
type Holiday struct {
	HolidayDate time.Time `json:"holiday_date"`
	Name        string    `json:"name"`
}
