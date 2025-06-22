package services

import (
	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"
	"mis_kursach_backend/configs"
)

func GenerateAuthToken(config configs.Config) *jwtauth.JWTAuth {
	tokenAuth := jwtauth.New("HS256", []byte(config.JWTConfig.Secret), nil)
	return tokenAuth
}

func GetHashPassword(password string) (string, error) {
	bytePassword := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
