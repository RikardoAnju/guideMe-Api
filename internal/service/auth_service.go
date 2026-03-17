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
		ID:                      uuid.New().String(),
		FirstName:               req.FirstName,
		LastName:                req.LastName,
		Username:                req.Username,
		Email:                   req.Email,
		Password:                string(hashed),
		PhoneNumber:             req.PhoneNumber,
		Gender:                  req.Gender,
		Address:                 req.Address,
		Role:                    "user",
		EmailVerified:           false,
		EmailVerificationToken:  verificationToken,
		EmailVerificationExpiry: &expiry,
		CreatedAt:               time.Now(),
	}

	result := config.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	emailService := NewEmailService()
	if err := emailService.SendVerificationEmail(user.Email, user.FirstName, verificationToken); err != nil {
		// Rollback user jika email gagal terkirim
		config.DB.Delete(user)
		return nil, errors.New("failed to send verification email, please try again")
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

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
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
		return "", nil, errors.New("failed to generate token")
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

	if user.EmailVerificationExpiry != nil && time.Now().After(*user.EmailVerificationExpiry) {
		return errors.New("verification token has expired, please request a new one")
	}

	if err := config.DB.Model(&user).Updates(map[string]interface{}{
		"email_verified":          true,
		"email_verification_token": "",
		"email_verification_expiry": nil,
	}).Error; err != nil {
		return errors.New("failed to verify email")
	}

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
		fmt.Println("EMAIL ERROR:", err)
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
	if err := config.DB.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"first_name":   req.FirstName,
		"last_name":    req.LastName,
		"username":     req.Username,
		"phone_number": req.PhoneNumber,
		"gender":       req.Gender,
		"address":      req.Address,
	}).Error; err != nil {
		return nil, errors.New("failed to update profile")
	}

	return GetProfile(userID)
}

// ResetPassword — cek email, generate OTP, kirim ke email, return otp_token ke FE
func ResetPassword(req models.ResetPasswordRequest) (string, error) {
	var user models.User
	result := config.DB.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		return "", errors.New("email not found")
	}

	// Generate OTP + JWT token yang menyimpan OTP (expire 10 menit)
	otp, otpToken, err := GenerateOTPToken(user.Email)
	if err != nil {
		return "", errors.New("failed to generate OTP")
	}

	// Kirim OTP via email
	emailService := NewEmailService()
	if err := emailService.SendOTPEmail(user.Email, user.FirstName, otp); err != nil {
		fmt.Println("OTP EMAIL ERROR:", err)
		return "", errors.New("failed to send OTP email")
	}

	// otp_token dikirim ke FE untuk disimpan sementara
	return otpToken, nil
}

// VerifyOTP — validasi OTP input user vs OTP di dalam JWT, return reset_token
func VerifyOTP(req models.VerifyOTPRequest) (string, error) {
	resetToken, err := VerifyOTPToken(req.OTP, req.OTPToken)
	if err != nil {
		return "", err
	}
	return resetToken, nil
}

// ChangePassword — parse reset_token (stateless), update password di DB
func ChangePassword(req models.ChangePasswordRequest) error {
	// Parse reset token untuk dapat email tanpa query DB
	email, err := ParseResetToken(req.ResetToken)
	if err != nil {
		return err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to process password")
	}

	result := config.DB.Model(&models.User{}).
		Where("email = ?", email).
		Updates(map[string]interface{}{
			"password":             string(hashed),
			"last_password_update": time.Now(),
		})

	if result.Error != nil {
		return errors.New("failed to update password")
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}