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
	var count int64
	config.DB.Model(&models.User{}).Where("email = ?", req.Email).Count(&count)
	if count > 0 {
		return nil, errors.New("email already registered")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	verificationToken := uuid.New().String()
	expiry := time.Now().Add(24 * time.Hour)

	user := &models.User{
		ID:                     uuid.New().String(),
		FirstName:              req.FirstName,
		LastName:               req.LastName,
		Username:               req.Username,
		Email:                  req.Email,
		Password:               string(hashed),
		PhoneNumber:            req.PhoneNumber,
		Gender:                 req.Gender,
		Address:                req.Address,
		Role:                   "user",
		EmailVerified:          false,
		EmailVerificationToken: verificationToken,
		EmailVerificationExpiry: &expiry,
		CreatedAt:              time.Now(),
	}

	result := config.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	emailService := NewEmailService()
	if err := emailService.SendVerificationEmail(user.Email, user.FirstName, verificationToken); err != nil {
		return nil, errors.New("user registered but failed to send verification email")
	}

	return user, nil
}

func Login(req models.LoginRequest) (string, *models.User, error) {
	var user models.User
	result := config.DB.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		return "", nil, errors.New("email or password incorrect")
	}

	if !user.EmailVerified {
		return "", nil, errors.New("email not verified, please check your inbox")
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

func VerifyEmail(token string) error {
	var user models.User
	result := config.DB.Where("email_verification_token = ?", token).First(&user)
	if result.Error != nil {
		return errors.New("invalid or expired verification token")
	}

	if user.EmailVerified {
		return errors.New("email already verified")
	}

	// Cek expiry
	if user.EmailVerificationExpiry != nil && time.Now().After(*user.EmailVerificationExpiry) {
		return errors.New("verification token has expired, please request a new one")
	}

	config.DB.Model(&user).Updates(map[string]interface{}{
		"email_verified":           true,
		"email_verification_token": "",
		"email_verification_expiry": nil,
	})

	return nil
}

func ResendVerificationEmail(email string) error {
	var user models.User
	result := config.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return errors.New("email not found")
	}

	if user.EmailVerified {
		return errors.New("email already verified")
	}

	newToken := uuid.New().String()
	expiry := time.Now().Add(24 * time.Hour)

	if err := config.DB.Model(&user).Updates(map[string]interface{}{
		"email_verification_token":  newToken,
		"email_verification_expiry": expiry,
	}).Error; err != nil {
		return errors.New("failed to update verification token")
	}

	emailService := NewEmailService()
	if err := emailService.SendVerificationEmail(user.Email, user.FirstName, newToken); err != nil {
		return errors.New("failed to send verification email")
	}

	return nil
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