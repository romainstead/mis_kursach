package db

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"mis_kursach_backend/internal/models"
	"time"
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
	err := dbpool.QueryRow(context.Background(),
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
			JOIN guests_in_bookings gib on gib.booking_id = b.id
			WHERE b.id = $1
`, id).Scan(&booking)
	if err != nil {
		return booking, fmt.Errorf("error getting booking: %v", err)
	}
	return booking, nil
}

func CreateBooking(dbpool *pgxpool.Pool, startDate time.Time,
	endDate time.Time, categoryCode int,
	checkIn *time.Time, checkOut *time.Time,
	roomNumber int, babyBed bool, guestName string, guestPassportNumber string, guestPhoneNumber string) error {
	var tariffs []models.Tariff
	var discounts []models.Discount
	var guests []models.Guest
	var discountAmount float64
	var discountID int
	var guestID int
	var bookingID int
	nights := int(endDate.Sub(startDate).Hours() / 24)
	if nights <= 0 {
		return fmt.Errorf("nights is zero or lower than zero")
	}

	err := pgxscan.Select(context.Background(), dbpool, &tariffs, `SELECT * FROM TARIFFS`)
	if err != nil {
		return fmt.Errorf("error fetching tariffs from db: %v", err)
	}

	err = pgxscan.Select(context.Background(), dbpool, &discounts, `SELECT * FROM DISCOUNTS`)
	if err != nil {
		return fmt.Errorf("error fetching discounts from db: %v", err)
	}

	err = pgxscan.Select(context.Background(), dbpool, &guests, `SELECT * FROM GUESTS G WHERE G.PASSPORT_NO = $1`, guestPassportNumber)
	if len(guests) == 0 {
		// если гостя нет, то добавляем его в таблицу и сразу вытаскиваем айди
		err = pgxscan.Get(context.Background(), dbpool, &guestID,
			`INSERT INTO GUESTS(name, phone_number, passport_no) VALUES ($1, $2, $3) RETURNING ID`, guestName, guestPhoneNumber, guestPassportNumber)
		if err != nil {
			return fmt.Errorf("error inserting guest: %v", err)
		}
	}
	// ретривим айди гостя
	err = pgxscan.Get(context.Background(), dbpool, &guestID, `SELECT G.ID FROM GUESTS G WHERE G.PASSPORT_NO = $1`, guestPassportNumber)
	if err != nil {
		return fmt.Errorf("error fetching guest: %v", err)
	}
	for _, discount := range discounts {
		if discount.MinNights == (1) {
			if nights < 3 {
				discountID = discount.ID
				discountAmount = discount.Amount
				break
			}
		} else if discount.MinNights == (2) {
			if nights < 7 {
				discountID = discount.ID
				discountAmount = discount.Amount
				break
			}
		} else if discount.MinNights == (3) {
			if nights < 14 {
				discountID = discount.ID
				discountAmount = discount.Amount
				break
			}
		}
	}
	var basePrice float64
	for _, tariff := range tariffs {
		if tariff.CategoryCode == categoryCode {
			basePrice = tariff.BasePrice
			break
		}
	}
	if basePrice == 0 {
		return fmt.Errorf("error fetching basePrice")
	}
	bookingSum := basePrice * float64(nights)
	totalSum := bookingSum * ((100 - discountAmount) / 100)

	err = pgxscan.Get(context.Background(), dbpool, &bookingID, `
					INSERT INTO BOOKINGS(
					status_code,
					start_date, end_date,
					check_in, check_out,
					baby_bed, booking_sum,
					discount_id, total_sum) VALUES (1, $1, $2, $3, $4, $5, $6, $7, $8) RETURNING ID`,
		startDate, endDate, checkIn, checkOut, babyBed, bookingSum, discountID, totalSum)
	if err != nil {
		return fmt.Errorf("error inserting booking: %v", err)
	}
	_, err = dbpool.Query(context.Background(),
		`INSERT INTO GUESTS_IN_BOOKINGS (GUEST_ID, BOOKING_ID, ROOM)
			VALUES ($1, $2, $3)`, guestID, bookingID, roomNumber)
	return err
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
	err := dbpool.QueryRow(context.Background(), `SELECT DISTINCT
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
					JOIN GUESTS_IN_BOOKINGS GIB ON B.ID = GIB.BOOKING_ID
				WHERE c.id=$1`, id).Scan(&complaint)
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
	err := dbpool.QueryRow(context.Background(),
		`SELECT 
				p.id, p.booking_id, p.amount, p.pay_date, 
				pm.name as method_name, 
				ps.name as status_name 
			FROM payments p
			JOIN payment_methods pm on p.method_code = pm.code
			JOIN payment_statuses ps on p.status_code = ps.status_code WHERE p.id=$1`, id).Scan(&payment)
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

//func SetMetrics(dbpool *pgxpool.Pool) (models.SetMetricsResponse, error) {
//	var metrics models.SetMetricsResponse
//	err := dbpool.QueryRow(context.Background(), `SELECT
//	ROUND(
//	(
//			SELECT
//				COUNT(G.PASSPORT_NO)
//			FROM
//				BOOKINGS B
//				JOIN BOOKING_STATUSES BS ON BS.STATUS_CODE = B.STATUS_CODE
//				JOIN GUESTS_IN_BOOKINGS GIB ON GIB.BOOKING_ID = B.ID
//				JOIN GUESTS G ON G.ID = GIB.GUEST_ID
//			WHERE
//				BS.STATUS_CODE = 1
//		) * 100.0
//		/
//		(
//			SELECT
//				SUM(RC.CAPACITY)
//			FROM
//				ROOMS R
//				JOIN ROOM_CATEGORIES RC ON R.CATEGORY_CODE = RC.CODE
//		) * 1.0
//	, 0) AS OCCUPANCY_RATIO`).Scan(&metrics.Occupancy)
//	if err != nil {
//		return metrics, fmt.Errorf("error setting metrics: %v", err)
//	}
//
//}
