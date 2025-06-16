package db

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
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

	r.Get("/GetBookingByID/{id}", handler.GetBookingByID)
	r.Get("/GetComplaintByID/{id}", handler.GetComplaintByID)
	r.Get("/GetPaymentByID/{id}", handler.GetPaymentByID)
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

func (p *PsHandler) GetBookingByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error getting booking by id: %v", err)
	}
	booking, err := GetBookingByID(p.db, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("Error getting booking: %v", err)
	}
	if err := json.NewEncoder(w).Encode(booking); err != nil {
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

func (p *PsHandler) GetComplaintByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error getting complaint by id: %v", err)
	}
	complaint, err := GetComplaintByID(p.db, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("Error getting complaint: %v", err)
	}
	if err := json.NewEncoder(w).Encode(complaint); err != nil {
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

func (p *PsHandler) GetPaymentByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error getting payment by id: %v", err)
	}
	complaint, err := GetPaymentByID(p.db, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("Error getting payment: %v", err)
	}
	if err := json.NewEncoder(w).Encode(complaint); err != nil {
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
