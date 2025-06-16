package main

import (
	"fmt"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"mis_kursach_backend/configs"
	"mis_kursach_backend/internal/db"
	"net/http"
)

func main() {
	// Инициализация конфига
	configs.InitConfig()
	config := configs.NewConfig()

	// Подключение к БД
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		config.DBConfig.Username, config.DBConfig.Password, config.DBConfig.Host, config.DBConfig.Port,
		config.DBConfig.Name)
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	// После завершения работы программы закрываем соединение с БД
	defer func() {
		dbInstance, _ := database.DB()
		dbInstance.Close()
	}()
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Use(middleware.Logger)
	r.Mount("/api", db.PsRoutes(database))

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Unable to start server: %v", err)
	}
	log.Print("Listening and serving HTTP on port 8080")
}
