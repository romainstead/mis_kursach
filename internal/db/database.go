package db

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"mis_kursach_backend/internal/models"
)

func GetAllBookings(dbpool *pgxpool.Pool) ([]*models.BookingResponse, error) {
	var bookings []*models.BookingResponse
	err := pgxscan.Select(context.Background(), dbpool, &bookings,
		`SELECT DISTINCT
				b.id AS "id", 
				b.start_date, 
				b.end_date, 
				b.check_in, 
				b.check_out, 
				b.baby_bed, 
				b.booking_sum, 
				b.total_sum, 
				bs.name AS "booking_status", 
				d.amount AS "discount_amount",
				gib.room
			FROM 
				bookings b
			JOIN 
				booking_statuses bs ON bs.status_code = b.status_code
			LEFT JOIN 
				discounts d ON d.id = b.discount_id
			JOIN guests_in_bookings gib on gib.booking_id = b.id`)
	if err != nil {
		return nil, fmt.Errorf("error getting all bookings: %v", err)
	}
	return bookings, nil
}

func GetBookingByID(dbpool *pgxpool.Pool, id int) (models.Booking, error) {
	var booking models.Booking
	err := dbpool.QueryRow(context.Background(), `SELECT * FROM bookings WHERE id=$1`, id).Scan(&booking)
	if err != nil {
		return booking, fmt.Errorf("error getting booking: %v", err)
	}
	return booking, nil
}

func GetAllComplaints(dbpool *pgxpool.Pool) ([]models.ComplaintResponse, error) {
	var complaints []models.ComplaintResponse
	err := pgxscan.Select(context.Background(), dbpool, &complaints,
		`SELECT DISTINCT
					C.ID,
					C.REASON,
					C.COMMENTARY,
					C.ISSUE_DATE,
					C.BOOKING_ID,
					CS.NAME AS STATUS,
					GIB.ROOM
				FROM
					COMPLAINTS C
					JOIN COMPLAINT_STATUSES CS ON C.STATUS_CODE = CS.STATUS_CODE
					JOIN BOOKINGS B ON C.BOOKING_ID = B.ID
					JOIN GUESTS_IN_BOOKINGS GIB ON B.ID = GIB.BOOKING_ID`)
	if err != nil {
		return nil, fmt.Errorf("error getting all complaints: %v", err)
	}
	return complaints, nil
}

func GetComplaintByID(dbpool *pgxpool.Pool, id int) (models.Complaint, error) {
	var complaint models.Complaint
	err := dbpool.QueryRow(context.Background(), `SELECT * FROM complaints WHERE id=$1`, id).Scan(&complaint)
	if err != nil {
		return complaint, fmt.Errorf("error getting complaint: %v", err)
	}
	return complaint, nil
}

func GetAllPayments(dbpool *pgxpool.Pool) ([]models.PaymentResponse, error) {
	var payments []models.PaymentResponse
	err := pgxscan.Select(context.Background(), dbpool, &payments,
		`select p.id, p.booking_id, p.amount, p.pay_date, pm.name as method_name, ps.name as status_name from payments p
				join payment_methods pm on p.method_code = pm.code
				join payment_statuses ps on p.status_code = ps.status_code`)
	if err != nil {
		return nil, fmt.Errorf("error getting all payments: %v", err)
	}
	return payments, nil
}

func GetPaymentByID(dbpool *pgxpool.Pool, id int) (models.Payment, error) {
	var payment models.Payment
	err := dbpool.QueryRow(context.Background(), `SELECT * FROM payments WHERE id=$1`, id).Scan(&payment)
	if err != nil {
		return payment, fmt.Errorf("error getting payment: %v", err)
	}
	return payment, nil
}

type GetAllRoomsResult struct {
	Number       int    `gorm:"column:number" json:"number"`
	CategoryName string `gorm:"column:category_name" json:"category_name"`
	StateName    string `gorm:"column:state_name" json:"state_name"`
	Capacity     int    `gorm:"column:capacity" json:"capacity"`
}

func GetAllRooms(dbpool *pgxpool.Pool) ([]GetAllRoomsResult, error) {
	var rooms []GetAllRoomsResult
	err := pgxscan.Select(context.Background(), dbpool, &rooms,
		`SELECT r.number, rc.name AS category_name, rs.name AS state_name, rc.capacity 
		FROM rooms r
		JOIN room_categories rc ON rc.code = r.category_code
		JOIN room_states rs ON rs.state_code = r.state_code`)
	if err != nil {
		return nil, fmt.Errorf("error getting all rooms: %v", err)
	}
	return rooms, nil
}

//
//func CreateBooking(db *gorm.DB, input models.Booking) (models.Booking, error) {
//	var booking models.Booking
//}
