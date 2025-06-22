package db

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"mis_kursach_backend/internal/models"
	"net/http"
	"strconv"
)

type PsHandler struct {
	dbpool *pgxpool.Pool
}

func PsRoutes(dbpool *pgxpool.Pool) chi.Router {
	r := chi.NewRouter()
	handler := &PsHandler{dbpool: dbpool}
	r.Get("/GetAllBookings", handler.GetAllBookings)
	r.Get("/GetAllComplaints", handler.GetAllComplaints)
	r.Get("/GetAllPayments", handler.GetAllPayments)
	r.Get("/GetAllRooms", handler.GetAllRooms)
	r.Get("/GetBookingByID/{id}", handler.GetBookingByID)
	r.Get("/GetComplaintByID/{id}", handler.GetComplaintByID)
	r.Get("/GetPaymentByID/{id}", handler.GetPaymentByID)
	r.Post("/CreateBooking", handler.CreateBooking)
	r.Post("/CreateComplaint", handler.CreateComplaint)
	return r
}

func (p *PsHandler) GetAllBookings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bookings, err := GetAllBookings(p.dbpool)
	if err != nil {
		http.Error(w, `{"error": "failed to get bookings"}`, http.StatusInternalServerError)
		log.Printf("Error getting all bookings: %v", err)
		return
	}
	if len(bookings) == 0 {
		http.Error(w, `{"error": "no bookings found"}`, http.StatusNotFound)
		return
	}
	if err := json.NewEncoder(w).Encode(bookings); err != nil {
		http.Error(w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
		log.Printf("Error encoding bookings: %v", err)
	}
}

func (p *PsHandler) GetBookingByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, `{"error": "invalid id"}`, http.StatusBadRequest)
		log.Printf("Error parsing booking id: %v", err)
		return
	}
	booking, err := GetBookingByID(p.dbpool, id)
	if err != nil {
		http.Error(w, `{"error": "booking not found"}`, http.StatusNotFound)
		log.Printf("Error getting booking: %v", err)
		return
	}
	if err := json.NewEncoder(w).Encode(booking); err != nil {
		http.Error(w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
		log.Printf("Error encoding booking: %v", err)
	}
}

func (p *PsHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var b models.CreateBookingInput
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		log.Printf("Error decoding booking input: %v", err)
		return
	}
	defer r.Body.Close()

	err := CreateBooking(p.dbpool, b)
	if err != nil {
		http.Error(w, `{"error": "failed to create booking"}`, http.StatusBadRequest)
		log.Printf("Error creating booking: %v", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Booking created successfully"})
}

func (p *PsHandler) GetAllComplaints(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	complaints, err := GetAllComplaints(p.dbpool)
	if err != nil {
		http.Error(w, `{"error": "failed to get complaints"}`, http.StatusInternalServerError)
		log.Printf("Error getting all complaints: %v", err)
		return
	}
	if len(complaints) == 0 {
		http.Error(w, `{"error": "no complaints found"}`, http.StatusNotFound)
		return
	}
	if err := json.NewEncoder(w).Encode(complaints); err != nil {
		http.Error(w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
		log.Printf("Error encoding complaints: %v", err)
	}
}

func (p *PsHandler) GetComplaintByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, `{"error": "invalid id"}`, http.StatusBadRequest)
		log.Printf("Error parsing complaint id: %v", err)
		return
	}
	complaint, err := GetComplaintByID(p.dbpool, id)
	if err != nil {
		http.Error(w, `{"error": "complaint not found"}`, http.StatusNotFound)
		log.Printf("Error getting complaint: %v", err)
		return
	}
	if err := json.NewEncoder(w).Encode(complaint); err != nil {
		http.Error(w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
		log.Printf("Error encoding complaint: %v", err)
	}
}

func (p *PsHandler) CreateComplaint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var c models.CreateComplaintInput
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		log.Printf("Error decoding complaint: %v", err)
		return
	}
	defer r.Body.Close()
	err := CreateComplaint(p.dbpool, c)
	if err != nil {
		http.Error(w, `{"error": "failed to create complaint"}`, http.StatusBadRequest)
		log.Printf("Error creating complaint: %v", err)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Complaint created successfully"})
}

func (p *PsHandler) GetAllPayments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	payments, err := GetAllPayments(p.dbpool)
	if err != nil {
		http.Error(w, `{"error": "failed to get payments"}`, http.StatusInternalServerError)
		log.Printf("Error getting all payments: %v", err)
		return
	}
	if len(payments) == 0 {
		http.Error(w, `{"error": "no payments found"}`, http.StatusNotFound)
		return
	}
	if err := json.NewEncoder(w).Encode(payments); err != nil {
		http.Error(w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
		log.Printf("Error encoding payments: %v", err)
	}
}

func (p *PsHandler) GetPaymentByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, `{"error": "invalid id"}`, http.StatusBadRequest)
		log.Printf("Error parsing payment id: %v", err)
		return
	}
	payment, err := GetPaymentByID(p.dbpool, id)
	if err != nil {
		http.Error(w, `{"error": "payment not found"}`, http.StatusNotFound)
		log.Printf("Error getting payment: %v", err)
		return
	}
	if err := json.NewEncoder(w).Encode(payment); err != nil {
		http.Error(w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
		log.Printf("Error encoding payment: %v", err)
	}
}

func (p *PsHandler) GetAllRooms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rooms, err := GetAllRooms(p.dbpool)
	if err != nil {
		http.Error(w, `{"error": "failed to get rooms"}`, http.StatusInternalServerError)
		log.Printf("Error getting all rooms: %v", err)
		return
	}
	if len(rooms) == 0 {
		http.Error(w, `{"error": "no rooms found"}`, http.StatusNotFound)
		return
	}
	if err := json.NewEncoder(w).Encode(rooms); err != nil {
		http.Error(w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
		log.Printf("Error encoding rooms: %v", err)
	}
}
