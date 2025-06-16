package db

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type PsHandler struct {
	db *gorm.DB
}

func PsRoutes(db *gorm.DB) chi.Router {
	r := chi.NewRouter()
	handler := &PsHandler{db: db}

	r.Get("/GetAllBookings", handler.GetAllBookings)
	r.Get("/GetAllComplaints", handler.GetAllComplaints)
	r.Get("/GetAllPayments", handler.GetAllPayments)
	r.Get("/GetAllRooms", handler.GetAllRooms)
	return r
}

func (p *PsHandler) GetAllBookings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	bookings, err := GetAllBookings(p.db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error getting all bookings: %v", err)
	}
	if len(bookings) == 0 {
		w.WriteHeader(http.StatusNotFound)
	}
	if err := json.NewEncoder(w).Encode(bookings); err != nil {
		http.Error(w, `{"error": "no bookings found"}`, http.StatusNotFound)
	}
}

func (p *PsHandler) GetAllComplaints(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	complaints, err := GetAllComplaints(p.db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error getting all complaints: %v", err)
	}
	if len(complaints) == 0 {
		w.WriteHeader(http.StatusNotFound)
	}
	if err := json.NewEncoder(w).Encode(complaints); err != nil {
		http.Error(w, `{"error": "no complaints found"}`, http.StatusNotFound)
	}
}

func (p *PsHandler) GetAllPayments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	payments, err := GetAllPayments(p.db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error getting all payments: %v", err)
	}
	if len(payments) == 0 {
		w.WriteHeader(http.StatusNotFound)
	}
	if err := json.NewEncoder(w).Encode(payments); err != nil {
		http.Error(w, `{"error": "no payments found"}`, http.StatusNotFound)
	}
}

func (p *PsHandler) GetAllRooms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rooms, err := GetAllRooms(p.db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error getting all rooms: %v", err)
	}
	if len(rooms) == 0 {
		w.WriteHeader(http.StatusNotFound)
	}
	if err := json.NewEncoder(w).Encode(rooms); err != nil {
		http.Error(w, `{"error": "no rooms found"}`, http.StatusNotFound)
	}
}
