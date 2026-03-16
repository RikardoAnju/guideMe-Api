package service

import (
	"errors"
	"fmt"
	"time"

	"guide-me/internal/config"
	"guide-me/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Register(req models.RegisterRequest) (*models.User, error) {
	// Cek email sudah ada
	var count int64
	config.DB.Model(&models.User{}).Where("email = ?", req.Email).Count(&count)
	if count > 0 {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:          uuid.New().String(),
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Username:    req.Username,
		Email:       req.Email,
		Password:    string(hashed),
		PhoneNumber: req.PhoneNumber,
		Gender:      req.Gender,
		Address:     req.Address,
		Role:        "user",
		CreatedAt:   time.Now(),
	}

	result := config.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func Login(req models.LoginRequest) (string, *models.User, error) {
	var user models.User

	result := config.DB.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		return "", nil, errors.New("email or password incorrect")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", nil, errors.New("email or password incorrect")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(config.JWTExpire).Unix(),
	})

	tokenString, err := token.SignedString(config.JWTSecret)
	if err != nil {
		return "", nil, err
	}

	return tokenString, &user, nil
}

func GetProfile(userID string) (*models.User, error) {
	var user models.User

	result := config.DB.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

func UpdateProfile(userID string, req models.UpdateProfileRequest) (*models.User, error) {
	result := config.DB.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"first_name":   req.FirstName,
		"last_name":    req.LastName,
		"username":     req.Username,
		"phone_number": req.PhoneNumber,
		"gender":       req.Gender,
		"address":      req.Address,
	})
	if result.Error != nil {
		return nil, result.Error
	}

	return GetProfile(userID)
}

func ResetPassword(req models.ResetPasswordRequest) (string, error) {
	var user models.User

	result := config.DB.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		return "", errors.New("email not found")
	}

	token := fmt.Sprintf("%d", time.Now().UnixNano())
	config.DB.Model(&user).Updates(map[string]interface{}{
		"temporary_reset_token": token,
		"password_reset":        true,
	})

	return token, nil
}

func ChangePassword(req models.ChangePasswordRequest) error {
	var user models.User

	result := config.DB.Where("temporary_reset_token = ? AND password_reset = ?", req.Token, true).First(&user)
	if result.Error != nil {
		return errors.New("invalid or expired token")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	now := time.Now()
	config.DB.Model(&user).Updates(map[string]interface{}{
		"password":               string(hashed),
		"temporary_reset_token":  "",
		"password_reset":         false,
		"password_reset_success": true,
		"last_password_update":   now,
	})

	return nil
}