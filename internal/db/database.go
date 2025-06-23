package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"mis_kursach_backend/internal/models"
	"mis_kursach_backend/internal/services"
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

func GetBookingByID(dbpool *pgxpool.Pool, id int) (models.BookingResponse, error) {
	var booking models.BookingResponse
	err := pgxscan.Get(context.Background(), dbpool, &booking, `SELECT DISTINCT
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
			WHERE b.id = $1`, id)
	if err != nil {
		return booking, fmt.Errorf("error getting booking: %v", err)
	}
	return booking, nil
}

func CreateBooking(dbpool *pgxpool.Pool, b models.CreateBookingInput) error {

	var tariffs []models.Tariff
	var discounts []models.Discount
	var guests []models.Guest
	var discountAmount float64
	var discountID int
	var guestID int
	var bookingID int

	// Парсинг времени
	startDate, err := time.Parse("2006-01-02", b.StartDate)
	if err != nil {
		return fmt.Errorf("couldn't parse start_date: %v", err)
	}
	endDate, err := time.Parse("2006-01-02", b.EndDate)
	if err != nil {
		return fmt.Errorf("couldn't parse end_date: %v", err)
	}

	// Подсчёт ночей
	nights := int(endDate.Sub(startDate).Hours() / 24)
	if nights <= 0 {
		return fmt.Errorf("nights is zero or lower than zero")
	}
	err = pgxscan.Select(context.Background(), dbpool, &tariffs, `SELECT * FROM TARIFFS`)
	if err != nil {
		return fmt.Errorf("error fetching tariffs from db: %v", err)
	}

	err = pgxscan.Select(context.Background(), dbpool, &discounts, `SELECT * FROM DISCOUNTS`)
	if err != nil {
		return fmt.Errorf("error fetching discounts from db: %v", err)
	}

	err = pgxscan.Select(context.Background(), dbpool, &guests, `SELECT * FROM GUESTS G WHERE G.PASSPORT_NO = $1`, b.GuestPassportNumber)
	if len(guests) == 0 {
		// если гостя нет, то добавляем его в таблицу и сразу вытаскиваем айди
		err = pgxscan.Get(context.Background(), dbpool, &guestID,
			`INSERT INTO GUESTS(name, phone_number, passport_no) VALUES ($1, $2, $3) RETURNING ID`, b.GuestName, b.GuestPhoneNumber, b.GuestPassportNumber)
		if err != nil {
			return fmt.Errorf("error inserting guest: %v", err)
		}
	}
	// ретривим айди гостя
	err = pgxscan.Get(context.Background(), dbpool, &guestID, `SELECT G.ID FROM GUESTS G WHERE G.PASSPORT_NO = $1`, b.GuestPassportNumber)
	if err != nil {
		return fmt.Errorf("error fetching guest: %v", err)
	}
	if nights >= 3 && nights < 7 {
		discountID = 1
	} else if nights >= 7 && nights < 14 {
		discountID = 2
	} else if nights >= 14 {
		discountID = 3
	} else {
		discountID = 4
	}
	var basePrice float64
	for _, tariff := range tariffs {
		if tariff.CategoryCode == b.CategoryCode {
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
		b.StartDate, b.EndDate, b.CheckIn, b.CheckOut, b.BabyBed, bookingSum, discountID, totalSum)
	if err != nil {
		return fmt.Errorf("error inserting booking: %v", err)
	}
	_, err = dbpool.Exec(context.Background(),
		`INSERT INTO GUESTS_IN_BOOKINGS (GUEST_ID, BOOKING_ID, ROOM)
			VALUES ($1, $2, $3)`, guestID, bookingID, b.RoomNumber)
	err = CreatePayment(dbpool, b, totalSum, bookingID)
	if err != nil {
		return fmt.Errorf("error inserting payment: %v", err)
	}
	return err
}

func DeleteBooking(dbpool *pgxpool.Pool, bookingID int) error {
	tx, err := dbpool.Begin(context.Background())
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		return fmt.Errorf("error beginning transaction: %v", err)
	}

	_, err = tx.Exec(context.Background(), `DELETE FROM PAYMENTS WHERE BOOKING_ID = $1`, bookingID)
	if err != nil {
		log.Printf("Error deleting payment, rolling back: %v", err)
		err := tx.Rollback(context.Background())
		if err != nil {
			log.Printf("Error rolling back transaction: %v", err)
			return fmt.Errorf("error rolling back transaction: %v", err)
		}
		log.Printf("error deleting payment while deleting booking: %v", err)
		return fmt.Errorf("error deleting payment while deleting booking: %v", err)
	}

	_, err = tx.Exec(context.Background(), `DELETE FROM GUESTS_IN_BOOKINGS WHERE BOOKING_ID = $1`, bookingID)
	if err != nil {
		log.Printf("Error deleting guest in booking, rolling back: %v", err)
		err := tx.Rollback(context.Background())
		if err != nil {
			return fmt.Errorf("error rolling back transaction: %v", err)
		}
		log.Printf("Error deleting guest in booking while deleting booking: %v", err)
		return fmt.Errorf("error deleting guest in booking while deleting booking: %v", err)
	}
	_, err = tx.Exec(context.Background(), `DELETE FROM COMPLAINTS WHERE BOOKING_ID = $1`, bookingID)
	if err != nil {
		log.Printf("Error deleting complaint, rolling back: %v", err)
		err := tx.Rollback(context.Background())
		if err != nil {
			log.Printf("Error rolling back transaction: %v", err)
			return fmt.Errorf("error rolling back transaction: %v", err)
		}
		log.Printf("Error deleting complaint while deleting booking: %v", err)
		return fmt.Errorf("error deleting complaint while deleting booking: %v", err)
	}
	_, err = tx.Exec(context.Background(), `DELETE FROM BOOKINGS WHERE ID = $1`, bookingID)
	if err != nil {
		log.Printf("Error deleting booking, rolling back: %v", err)
		err := tx.Rollback(context.Background())
		if err != nil {
			log.Printf("Error rolling back transaction: %v", err)
			return fmt.Errorf("error rolling back transaction: %v", err)
		}
		log.Printf("Error deleting booking while deleting booking: %v", err)
		return fmt.Errorf("error deleting booking while deleting booking: %v", err)
	}
	err = tx.Commit(context.Background())
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		return fmt.Errorf("error commiting transaction while deleting booking: %v", err)
	}
	return nil
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

func CreateComplaint(dbpool *pgxpool.Pool, complaint models.CreateComplaintInput) error {
	_, err := dbpool.Exec(context.Background(),
		`INSERT INTO COMPLAINTS(reason, commentary, issue_date, booking_id, status_code) VALUES ($1, $2, $3, $4, $5)`,
		complaint.Reason, complaint.Commentary, time.Now(), complaint.BookingID, 1)
	if err != nil {
		return fmt.Errorf("error inserting complaint: %v", err)
	}
	return nil
}

func DeleteComplaint(dbpool *pgxpool.Pool, complaintID int) error {
	_, err := dbpool.Exec(context.Background(), `DELETE FROM COMPLAINTS WHERE ID = $1`, complaintID)
	if err != nil {
		return fmt.Errorf("error deleting complaint: %v", err)
	}
	return nil
}

func UpdateComplaint(dbpool *pgxpool.Pool, c models.UpdateComplaintRequest) error {
	var statusCode int
	err := dbpool.QueryRow(context.Background(), `SELECT status_code FROM COMPLAINT_STATUSES WHERE name = $1`, c.Status).Scan(&statusCode)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("статус %s не найден", c.Status)
		}
		return fmt.Errorf("ошибка при поиске status_code: %v", err)
	}

	query := `
		UPDATE COMPLAINTS 
		SET reason = $1, commentary = $2, issue_date = $3, status_code = $4
		WHERE id = $5`
	result, err := dbpool.Exec(context.Background(), query, c.Reason, c.Commentary, c.IssueDate, statusCode, c.ID)
	if err != nil {
		log.Printf("Error updating complaint: %v", err)
		return fmt.Errorf("error updating complaint: %v", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("complaint with ID %d not found", c.ID)
	}

	return nil

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

func CreatePayment(dbpool *pgxpool.Pool, b models.CreateBookingInput, amount float64, bookingID int) error {
	_, err := dbpool.Exec(context.Background(),
		`INSERT INTO Payments(booking_id, pay_date, amount, method_code, status_code) VALUES ($1, $2, $3, $4, $5)`,
		bookingID, time.Now(), amount, b.MethodCode, 1)
	if err != nil {
		log.Printf("error inserting payment: %v", err)
		return fmt.Errorf("error inserting payment: %v", err)
	}
	return nil
}

func DeletePayment(dbpool *pgxpool.Pool, paymentID int) error {
	_, err := dbpool.Exec(context.Background(), `DELETE FROM Payments WHERE id = $1`, paymentID)
	if err != nil {
		return fmt.Errorf("error deleting payment: %v", err)
	}
	return nil
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

func SetMetrics(dbpool *pgxpool.Pool) (models.SetMetricsResponse, error) {
	var metrics models.SetMetricsResponse
	err := dbpool.QueryRow(context.Background(), `SELECT
	ROUND(
	(
			SELECT
				COUNT(G.PASSPORT_NO)
			FROM
				BOOKINGS B
				JOIN BOOKING_STATUSES BS ON BS.STATUS_CODE = B.STATUS_CODE
				JOIN GUESTS_IN_BOOKINGS GIB ON GIB.BOOKING_ID = B.ID
				JOIN GUESTS G ON G.ID = GIB.GUEST_ID
			WHERE
				BS.STATUS_CODE = 1
		) * 100.0
		/
		(
			SELECT
				SUM(RC.CAPACITY)
			FROM
				ROOMS R
				JOIN ROOM_CATEGORIES RC ON R.CATEGORY_CODE = RC.CODE
		) * 1.0
	, 0) AS OCCUPANCY_RATIO`).Scan(&metrics.Occupancy)

	if err != nil {
		return metrics, fmt.Errorf("error setting metrics: %v", err)
	}
	err = dbpool.QueryRow(context.Background(),
		`SELECT
				COUNT(*) AS UNPAID_BOOKINGS
			FROM
				PAYMENTS P
			JOIN PAYMENT_STATUSES PS ON P.STATUS_CODE = PS.STATUS_CODE
			WHERE
				P.STATUS_CODE = 1`).Scan(&metrics.UnpaidBookings)
	if err != nil {
		return metrics, fmt.Errorf("error setting metrics: %v", err)
	}

	err = dbpool.QueryRow(context.Background(),
		`SELECT COUNT(ID) AS ACTIVE_BOOKINGS FROM BOOKINGS B WHERE B.STATUS_CODE = 1`).Scan(&metrics.CurrentBookings)
	if err != nil {
		return metrics, fmt.Errorf("error setting metrics: %v", err)
	}

	err = dbpool.QueryRow(context.Background(),
		`SELECT COUNT(ID) AS OPEN_COMPLAINTS FROM COMPLAINTS WHERE STATUS_CODE = 1`).Scan(&metrics.OpenComplaints)
	if err != nil {
		return metrics, fmt.Errorf("error setting metrics: %v", err)
	}

	err = dbpool.QueryRow(context.Background(),
		`SELECT
				COUNT(*)
			FROM
			(
				SELECT
					R.NUMBER
				FROM
					ROOMS R
				EXCEPT
				SELECT
					GIB.ROOM AS NUMBER
				FROM
					BOOKINGS B
				JOIN GUESTS_IN_BOOKINGS GIB ON B.ID = GIB.BOOKING_ID
				WHERE
					B.STATUS_CODE = 1
		) AS FREE_ROOMS`).Scan(&metrics.FreeRooms)
	if err != nil {
		return metrics, fmt.Errorf("error setting metrics: %v", err)
	}
	err = dbpool.QueryRow(context.Background(),
		`SELECT COUNT(R.NUMBER) 
			FROM ROOMS R 
    		JOIN ROOM_STATES RS ON R.STATE_CODE = RS.STATE_CODE 
			WHERE R.STATE_CODE = 3`).Scan(&metrics.RoomsUnderMaintenance)
	if err != nil {
		return metrics, fmt.Errorf("error setting metrics: %v", err)
	}

	err = dbpool.QueryRow(context.Background(),
		`SELECT
				COALESCE(SUM(AMOUNT), 0) AS REVENUE_7DAYS
			FROM
				PAYMENTS
			WHERE
				PAY_DATE BETWEEN NOW() - INTERVAL '7 DAYS' AND NOW()`).Scan(&metrics.Revenue7Days)
	if err != nil {
		return metrics, fmt.Errorf("error setting metrics: %v", err)
	}

	err = dbpool.QueryRow(context.Background(),
		`SELECT
	COALESCE(ROUND(
		(
			SELECT
				SUM(AMOUNT)
			FROM
				PAYMENTS P
		) / (
			SELECT
				COUNT(*)
			FROM
				ROOMS
		),
		2
	), 0) AS REVPAR`).Scan(&metrics.RevPar)
	if err != nil {
		return metrics, fmt.Errorf("error setting metrics: %v", err)
	}

	err = dbpool.QueryRow(context.Background(),
		`SELECT
				COUNT(*) AS NEW_GUESTS_7DAYS
			FROM
				BOOKINGS
			WHERE
				START_DATE < NOW() + INTERVAL '7 DAYS'
			AND START_DATE > NOW() - INTERVAL '7 DAYS'`).Scan(&metrics.NewGuests7Days)
	if err != nil {
		return metrics, fmt.Errorf("error setting metrics: %v", err)
	}

	err = dbpool.QueryRow(context.Background(),
		`SELECT
			COALESCE(SUM(AMOUNT), 0) AS REVPAC_7DAYS
		FROM
			PAYMENTS P
			JOIN BOOKINGS B ON P.BOOKING_ID = B.ID
			JOIN GUESTS_IN_BOOKINGS GIB ON GIB.BOOKING_ID = P.BOOKING_ID
			JOIN GUESTS G ON G.ID = GIB.GUEST_ID
		WHERE 
		    START_DATE < NOW() + INTERVAL '7 DAYS' `).Scan(&metrics.RevPac)
	if err != nil {
		return metrics, fmt.Errorf("error setting metrics: %v", err)
	}
	return metrics, nil
}

func CreateUser(dbpool *pgxpool.Pool, user models.UserRequestBody) (int, error) {
	var UserID int
	hashPassword, err := services.GetHashPassword(user.Password)
	if err != nil {
		return 0, fmt.Errorf("error getting hash password: %v", err)
	}
	err = pgxscan.Get(context.Background(), dbpool, &UserID,
		`INSERT INTO USERS (username, hash) VALUES ($1, $2) RETURNING ID`, user.Username, hashPassword)
	if err != nil {
		return 0, fmt.Errorf("error creating user in db: %v", err)
	}
	return UserID, nil
}

func GetUser(dbpool *pgxpool.Pool, user *models.UserRequestBody) (*models.User, error) {
	query := `SELECT * FROM USERS WHERE USERNAME = $1`
	var u models.User
	err := pgxscan.Get(context.Background(), dbpool, &u,
		query, user.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found in db")
		}
		return nil, fmt.Errorf("error getting user: %v", err)
	}
	return &u, nil
}

func DeleteUser(dbpool *pgxpool.Pool, usernameToDelete string) error {
	_, err := dbpool.Exec(context.Background(), `DELETE FROM USERS WHERE USERNAME = $1`, usernameToDelete)
	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}
	return nil
}

func CreateGuest(dbpool *pgxpool.Pool, guest models.Guest) error {
	_, err := dbpool.Exec(context.Background(),
		`INSERT INTO GUESTS (name, phone_number, passport_no) VALUES ($1, $2, $3)`,
		guest.Name, guest.PhoneNumber, guest.PassportNo)
	if err != nil {
		return fmt.Errorf("error inserting guest: %v", err)
	}
	return nil
}

func GetAllGuests(dbpool *pgxpool.Pool) ([]models.Guest, error) {
	var guests []models.Guest
	err := pgxscan.Select(context.Background(), dbpool, &guests, "SELECT * FROM GUESTS ORDER BY id")
	if err != nil {
		log.Printf("error getting all guests: %v", err)
		return nil, fmt.Errorf("error getting all guests: %v", err)
	}
	return guests, nil
}

func DeleteGuest(dbpool *pgxpool.Pool, id int) error {
	_, err := dbpool.Exec(context.Background(), `DELETE FROM GUESTS WHERE ID = $1`, id)
	if err != nil {
		return fmt.Errorf("error deleting guest: %v", err)
	}
	return nil
}

func GetRoomCategories(dbpool *pgxpool.Pool) ([]models.RoomCategory, error) {
	var categories []models.RoomCategory
	err := pgxscan.Select(context.Background(), dbpool, &categories, "SELECT * FROM ROOM_CATEGORIES")
	if err != nil {
		log.Printf("error getting categories: %v", err)
		return nil, fmt.Errorf("error getting categories: %v", err)
	}
	return categories, nil
}

func GetPaymentMethods(dbpool *pgxpool.Pool) ([]models.PaymentMethod, error) {
	var paymentMethods []models.PaymentMethod
	err := pgxscan.Select(context.Background(), dbpool, &paymentMethods, "SELECT * FROM Payment_Methods")
	if err != nil {
		log.Printf("error getting payment methods: %v", err)
		return nil, fmt.Errorf("error getting payment methods: %v", err)
	}
	return paymentMethods, nil
}

func GetFreeRooms(dbpool *pgxpool.Pool, start string, end string, categoryCode int) ([]int, error) {
	var freeRooms []int
	err := pgxscan.Select(context.Background(), dbpool, &freeRooms,
		`SELECT r.number
		FROM 
		    ROOMS r
		WHERE r.number NOT IN (
			SELECT g.room
			FROM BOOKINGS b
			JOIN GUESTS_IN_BOOKINGS g ON g.booking_id = b.id
			WHERE NOT ($2 <= b.start_date OR $1 >= b.end_date)
		) AND STATE_CODE = 1 AND R.CATEGORY_CODE = $3`, start, end, categoryCode)
	if err != nil {
		log.Printf("error getting free rooms: %v", err)
		return nil, fmt.Errorf("error getting free rooms: %v", err)
	}
	return freeRooms, nil
}

func ConfirmBooking(dbpool *pgxpool.Pool, id int) error {
	_, err := dbpool.Exec(context.Background(), `UPDATE BOOKINGS SET STATUS_CODE = 1 WHERE ID = $1`, id)
	if err != nil {
		log.Printf("error updating booking: %v", err)
		return fmt.Errorf("error updating booking: %v", err)
	}
	_, err = dbpool.Exec(context.Background(), `UPDATE ROOMS R
        SET STATE_CODE = 2
        FROM GUESTS_IN_BOOKINGS GIB
        JOIN BOOKINGS B ON GIB.BOOKING_ID = B.ID
        WHERE GIB.ROOM = R.NUMBER AND B.ID = $1`, id)
	if err != nil {
		log.Printf("error updating room: %v", err)
		return fmt.Errorf("error updating room: %v", err)
	}
	return nil
}
