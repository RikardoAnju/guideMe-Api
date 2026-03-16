package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"

	"guide-me/internal/config"
	"guide-me/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mailersend/mailersend-go"
	"golang.org/x/crypto/bcrypt"
)

// ─── Helper ───────────────────────────────────────────────────────────────────

func generateToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func sendVerificationEmail(toEmail, toName, token string) error {
	ms := mailersend.NewMailersend(os.Getenv("MAILERSEND_API_KEY"))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	appURL := os.Getenv("APP_URL")
	verificationLink := fmt.Sprintf("%s/api/auth/verify-email?token=%s", appURL, token)

	from := mailersend.From{
		Name:  os.Getenv("MAILERSEND_FROM_NAME"),
		Email: os.Getenv("MAILERSEND_FROM_EMAIL"),
	}

	recipients := []mailersend.Recipient{
		{Name: toName, Email: toEmail},
	}

	htmlContent := fmt.Sprintf(`
		<div style="font-family:Arial,sans-serif;max-width:600px;margin:0 auto;">
			<h2>Hello, %s!</h2>
			<p>Thank you for registering. Please verify your email address by clicking the button below:</p>
			<a href="%s" style="
				background-color:#4F46E5;
				color:white;
				padding:12px 24px;
				text-decoration:none;
				border-radius:6px;
				display:inline-block;
				margin:16px 0;
			">Verify Email</a>
			<p>Or copy this link into your browser:</p>
			<p style="color:#6B7280;word-break:break-all;">%s</p>
			<p style="color:#EF4444;">This link will expire in <strong>24 hours</strong>.</p>
		</div>
	`, toName, verificationLink, verificationLink)

	message := ms.Email.NewMessage()
	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject("Verify Your Email Address")
	message.SetHTML(htmlContent)

	_, err := ms.Email.Send(ctx, message)
	return err
}

// ─── Auth Service ─────────────────────────────────────────────────────────────

func Register(req models.RegisterRequest) (*models.User, error) {
	// Cek email sudah terdaftar
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

	// Generate token verifikasi
	verifToken, err := generateToken(32)
	if err != nil {
		return nil, err
	}
	verifExpiry := time.Now().Add(24 * time.Hour)

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
		EmailVerificationToken:  verifToken,
		EmailVerificationExpiry: &verifExpiry,
		CreatedAt:               time.Now(),
	}

	if err := config.DB.Create(user).Error; err != nil {
		return nil, err
	}

	// Kirim email verifikasi (non-blocking)
	fullName := req.FirstName + " " + req.LastName
	go func() {
		if err := sendVerificationEmail(req.Email, fullName, verifToken); err != nil {
			fmt.Printf("[ERROR] Failed to send verification email: %v\n", err)
		}
	}()

	return user, nil
}

func VerifyEmail(token string) error {
	var user models.User

	result := config.DB.Where(
		"email_verification_token = ? AND email_verification_expiry > ?",
		token, time.Now(),
	).First(&user)

	if result.Error != nil {
		return errors.New("token invalid or expired")
	}

	return config.DB.Model(&user).Updates(map[string]interface{}{
		"email_verified":            true,
		"email_verification_token":  "",
		"email_verification_expiry": nil,
	}).Error
}

func ResendVerificationEmail(email string) error {
	var user models.User

	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return errors.New("email not found")
	}

	if user.EmailVerified {
		return errors.New("email already verified")
	}

	// Generate token baru
	verifToken, err := generateToken(32)
	if err != nil {
		return err
	}
	verifExpiry := time.Now().Add(24 * time.Hour)

	config.DB.Model(&user).Updates(map[string]interface{}{
		"email_verification_token":  verifToken,
		"email_verification_expiry": verifExpiry,
	})

	fullName := user.FirstName + " " + user.LastName
	return sendVerificationEmail(user.Email, fullName, verifToken)
}

func Login(req models.LoginRequest) (string, *models.User, error) {
	var user models.User

	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return "", nil, errors.New("email or password incorrect")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", nil, errors.New("email or password incorrect")
	}

	// Cek apakah email sudah diverifikasi
	if !user.EmailVerified {
		return "", nil, errors.New("please verify your email before logging in")
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
	if err := config.DB.Where("id = ?", userID).First(&user).Error; err != nil {
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
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
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
	result := config.DB.Where(
		"temporary_reset_token = ? AND password_reset = ?",
		req.Token, true,
	).First(&user)

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