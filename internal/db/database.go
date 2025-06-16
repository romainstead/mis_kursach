package db

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"mis_kursach_backend/internal/models"
)

func GetAllBookings(db *gorm.DB) ([]models.Booking, error) {
	var bookings []models.Booking
	err := db.Preload("Discount").Preload("Status").Find(&bookings).Error
	if err != nil {
		return nil, fmt.Errorf("error getting all bookings: %v", err)
	}
	return bookings, nil
}

func GetBookingByID(db *gorm.DB, id int) (models.Booking, error) {
	var booking models.Booking
	result := db.Preload("Discount").Preload("Status").First(&booking, id)
	if result.Error != nil {
		if errors.Is(gorm.ErrRecordNotFound, result.Error) {
			return models.Booking{}, fmt.Errorf("booking with id %v not found", id)
		} else {
			return models.Booking{}, fmt.Errorf("error getting booking: %v", result.Error)
		}
	}
	return booking, nil
}

func GetAllComplaints(db *gorm.DB) ([]models.Complaint, error) {
	var complaints []models.Complaint
	err := db.Preload("Booking").Preload("Status").Find(&complaints).Error
	if err != nil {
		return nil, fmt.Errorf("error getting all complaints: %v", err)
	}
	return complaints, nil
}

func GetComplaintByID(db *gorm.DB, id int) (models.Complaint, error) {
	var complaint models.Complaint
	result := db.Preload("Booking").Preload("Status").First(&complaint, id)
	if result.Error != nil {
		if errors.Is(gorm.ErrRecordNotFound, result.Error) {
			return models.Complaint{}, fmt.Errorf("complaint with id %v not found", id)
		} else {
			return models.Complaint{}, fmt.Errorf("error getting complaint: %v", result.Error)
		}
	}
	return complaint, nil
}

func GetAllPayments(db *gorm.DB) ([]models.Payment, error) {
	var payments []models.Payment
	err := db.Preload("Booking").Preload("Method").Preload("Status").Find(&payments).Error
	if err != nil {
		return nil, fmt.Errorf("error getting all complaints: %v", err)
	}
	return payments, nil
}

func GetPaymentByID(db *gorm.DB, id int) (models.Payment, error) {
	var payment models.Payment
	result := db.Preload("Booking").Preload("Method").Preload("Status").First(&payment, id)
	if result.Error != nil {
		if errors.Is(gorm.ErrRecordNotFound, result.Error) {
			return models.Payment{}, fmt.Errorf("payment with id %v not found", id)
		} else {
			return models.Payment{}, fmt.Errorf("error getting payment: %v", result.Error)
		}
	}
	return payment, nil
}

type GetAllRoomsResult struct {
	Number       int    `gorm:"column:number" json:"number"`
	CategoryName string `gorm:"column:category_name" json:"category_name"`
	StateName    string `gorm:"column:state_name" json:"state_name"`
	Capacity     int    `gorm:"column:capacity" json:"capacity"`
}

func GetAllRooms(db *gorm.DB) ([]GetAllRoomsResult, error) {
	var rooms []GetAllRoomsResult
	db.Raw(`
        SELECT r.number, rc.name AS category_name, rs.name AS state_name, rc.capacity FROM rooms r
		JOIN room_categories rc ON rc.code = r.category_code
		JOIN room_states rs ON rs.state_code = r.state_code
    `).Scan(&rooms)
	return rooms, nil
}
