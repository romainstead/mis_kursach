package db

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"mis_kursach_backend/configs"
	"mis_kursach_backend/internal/models"
	"mis_kursach_backend/internal/services"
	"net/http"
	"strconv"
	"time"
)

type PsHandler struct {
	dbpool  *pgxpool.Pool
	jwtauth *jwtauth.JWTAuth
}

func PsRoutes(dbpool *pgxpool.Pool, config configs.Config) chi.Router {
	r := chi.NewRouter()
	tokenAuth := services.GenerateAuthToken(config)
	handler := &PsHandler{dbpool: dbpool, jwtauth: tokenAuth}
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))
		r.Get("/GetAllBookings", handler.GetAllBookings)
		r.Get("/GetAllComplaints", handler.GetAllComplaints)
		r.Get("/GetAllPayments", handler.GetAllPayments)
		r.Get("/GetAllRooms", handler.GetAllRooms)
		r.Get("/GetBookingByID/{id}", handler.GetBookingByID)
		r.Get("/GetComplaintByID/{id}", handler.GetComplaintByID)
		r.Get("/GetPaymentByID/{id}", handler.GetPaymentByID)
		r.Post("/CreateBooking", handler.CreateBooking)
		r.Get("/SetMetrics", handler.SetMetrics)
		r.Delete("/DeleteBooking/{id}", handler.DeleteBooking)
		r.Delete("/DeleteComplaint/{id}", handler.DeleteComplaint)
		r.Delete("/DeletePayment/{id}", handler.DeletePayment)
		r.Delete("/DeleteUser", handler.DeleteUser)
		r.Post("/CreateUser", handler.CreateUser)
		// TODO: UPDATE-функции
	})

	r.Group(func(r chi.Router) {
		r.Post("/CreateComplaint", handler.CreateComplaint)
		r.Put("/UpdateComplaint", handler.UpdateComplaint)
		r.Post("/login", handler.Login)
		r.Post("/logout", handler.Logout)
	})

	return r
}

func (p *PsHandler) SetMetrics(w http.ResponseWriter, r *http.Request) {
	time.Sleep(1 * time.Second)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	metrics, err := SetMetrics(p.dbpool)
	if err != nil {
		http.Error(w, `{"error": "failed to set metrics"}`, http.StatusInternalServerError)
		log.Printf("Error getting setting metrics: %v", err)
		return
	}
	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		http.Error(w, `{"error": "failed to encode metrics"}`, http.StatusInternalServerError)
		log.Printf("Error encoding metrics: %v", err)
	}
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

func (p *PsHandler) DeleteBooking(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, `{"error": "invalid id"}`, http.StatusBadRequest)
	}
	err = DeleteBooking(p.dbpool, id)
	if err != nil {
		http.Error(w, `{"error": "failed to delete booking"}`, http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	response := make(map[string]string)
	response["message"] = "success"
	json.NewEncoder(w).Encode(response)
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

func (p *PsHandler) DeleteComplaint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, `{"error": "invalid id"}`, http.StatusBadRequest)
	}
	err = DeleteComplaint(p.dbpool, id)
	if err != nil {
		http.Error(w, `{"error": "complaint not found"}`, http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	response := make(map[string]string)
	response["message"] = "success"
	json.NewEncoder(w).Encode(response)
}

func (p *PsHandler) UpdateComplaint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var c models.UpdateComplaintRequest
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		log.Printf("Error decoding complaint: %v", err)
		return
	}
	defer r.Body.Close()

	if c.Reason == "" {
		http.Error(w, `{"error": "reason cannot be empty"}`, http.StatusBadRequest)
		return
	}
	if c.Status == "" {
		http.Error(w, `{"error": "status is required"}`, http.StatusBadRequest)
		return
	}

	err := UpdateComplaint(p.dbpool, c)
	if err != nil {
		http.Error(w, `{"error": "failed to update complaint"}`, http.StatusBadRequest)
		log.Printf("Error updating complaint: %v", err)
		return
	}

	w.WriteHeader(200)
	response := make(map[string]string)
	response["message"] = "success"
	json.NewEncoder(w).Encode(response)
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

func (p *PsHandler) DeletePayment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, `{"error": "invalid id"}`, http.StatusBadRequest)
	}
	err = DeletePayment(p.dbpool, id)
	if err != nil {
		http.Error(w, `{"error": "payment not found"}`, http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	response := make(map[string]string)
	response["message"] = "success"
	json.NewEncoder(w).Encode(response)
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

func (p *PsHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userReqBody models.UserRequestBody
	if err := json.NewDecoder(r.Body).Decode(&userReqBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid request body"}`))
		return
	}
	UserID, err := CreateUser(p.dbpool, userReqBody)
	if err != nil {
		http.Error(w, `{"error": "failed to create user"}`, http.StatusInternalServerError)
		log.Printf("Error creating user: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	response := make(map[string]string)
	response["message"] = "user created successfully"
	response["user_id"] = strconv.Itoa(UserID)
	json.NewEncoder(w).Encode(response)
}

func (p *PsHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	username := r.URL.Query().Get("username")
	err := DeleteUser(p.dbpool, username)
	if err != nil {
		http.Error(w, `{"error": "failed to delete user"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	response := make(map[string]string)
	response["message"] = "success"
	json.NewEncoder(w).Encode(response)
}

func (p *PsHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userReqBody := new(models.UserRequestBody)
	if err := json.NewDecoder(r.Body).Decode(&userReqBody); err != nil {
		log.Printf("Invalid request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid request body"}`))
		log.Printf("Error decoding request body: %v", err)
		return
	}
	user, err := GetUser(p.dbpool, userReqBody)
	if err != nil {
		http.Error(w, `{"error": "failed to get user"}`, http.StatusBadRequest)
		log.Printf("Failed to get user: %v", err)
		return
	}
	if !services.CheckPasswordHash(userReqBody.Password, user.Hash) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "incorrect password"}`))
		return
	}
	claims := map[string]interface{}{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix()}
	_, tokenString, err := p.jwtauth.Encode(claims)
	if err != nil {
		http.Error(w, `{"error": "failed to generate token"}`, http.StatusInternalServerError)
	}
	response := make(map[string]string)
	response["token"] = tokenString
	response["username"] = user.Username
	json.NewEncoder(w).Encode(response)
}

func (p *PsHandler) Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (p *PsHandler) CreateGuest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var g models.Guest
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid request body"}`))
		log.Printf("Error decoding request body: %v", err)
		return
	}
	err := CreateGuest(p.dbpool, g)
	if err != nil {
		http.Error(w, `{"error": "failed to create guest"}`, http.StatusInternalServerError)
		log.Printf("Error creating guest: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	response := make(map[string]string)
	response["message"] = "success"
	json.NewEncoder(w).Encode(response)
}

func (p *PsHandler) GetAllGuests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	guests, err := GetAllGuests(p.dbpool)
	if err != nil {
		http.Error(w, `{"error": "failed to get guests"}`, http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(guests); err != nil {
		http.Error(w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
		log.Printf("Error encoding guests: %v", err)
		return
	}
}

func (p *PsHandler) DeleteGuest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Printf("Invalid id: %v", err)
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		return
	}
	err = DeleteGuest(p.dbpool, id)
	if err != nil {
		log.Printf("Failed to delete guest: %v", err)
		http.Error(w, `{"error": "failed to delete guest"}`, http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	response := make(map[string]string)
	response["message"] = "success"
	json.NewEncoder(w).Encode(response)
}
