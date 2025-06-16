package db

import (
	"fmt"
	"gorm.io/gorm"
	"mis_kursach_backend/internal/models"
)

func GetAllBookings(db *gorm.DB) ([]models.Booking, error) {
	var bookings []models.Booking
	err := db.Preload("Discount").Find(&bookings).Error
	if err != nil {
		return nil, fmt.Errorf("Error getting all bookings: %v", err)
	}
	return bookings, nil
}

func GetAllComplaints(db *gorm.DB) ([]models.Complaint, error) {
	var complaints []models.Complaint
	err := db.Preload("Booking").Preload("Status").Find(&complaints).Error
	if err != nil {
		return nil, fmt.Errorf("Error getting all complaints: %v", err)
	}
	return complaints, nil
}

func GetAllPayments(db *gorm.DB) ([]models.Payment, error) {
	var payments []models.Payment
	err := db.Preload("Booking").Preload("Method").Preload("Status").Find(&payments).Error
	if err != nil {
		return nil, fmt.Errorf("Error getting all complaints: %v", err)
	}
	return payments, nil
}

type GetAllRoomsResult struct {
	Number       int    `gorm:"column:number" json:"number"`
	CategoryName string `gorm:"column:category_name" json:"category_name"`
	StateName    string `gorm:"column:state_name" json:"state_name"`
	Capacity     int    `gorm:"column:capacity" json:"capacity"`
	Prices       string `gorm:"type:json;column:prices" json:"prices"`
}

func GetAllRooms(db *gorm.DB) ([]GetAllRoomsResult, error) {
	var rooms []GetAllRoomsResult
	db.Raw(`
        SELECT 
            r.number,
            c.name AS category_name,
            s.name AS state_name,
            c.capacity,
            array_agg(
                json_build_object(
                    'day_code', t.day_code,
                    'price', t.base_price
                )
            ) AS prices
        FROM rooms r
        JOIN room_categories c ON r.category_code = c.code
        JOIN room_states s ON r.state_code = s.state_code
        JOIN tariffs t ON t.category_code = c.code
        GROUP BY r.number, c.name, s.name, c.capacity
    `).Scan(&rooms)
	return rooms, nil
}
