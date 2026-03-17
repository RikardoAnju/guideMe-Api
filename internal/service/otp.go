package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type OTPClaims struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
	jwt.RegisteredClaims
}

type ResetClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// Generate 6 digit OTP secara kriptografis aman
func generateOTP() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// GenerateOTPToken — buat OTP + JWT yang menyimpan OTP (expire 10 menit)
// Return: otp (untuk dikirim email), otpToken (untuk dikirim ke FE), error
func GenerateOTPToken(email string) (string, string, error) {
	otp, err := generateOTP()
	if err != nil {
		return "", "", errors.New("failed to generate OTP")
	}

	claims := OTPClaims{
		Email: email,
		OTP:   otp,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", "", errors.New("failed to sign OTP token")
	}

	return otp, signed, nil
}

// VerifyOTPToken — validasi OTP input user vs OTP di dalam JWT
// Jika valid, return reset_token untuk dipakai ganti password
func VerifyOTPToken(otpInput, tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &OTPClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return "", errors.New("OTP token invalid atau sudah expired")
	}

	claims, ok := token.Claims.(*OTPClaims)
	if !ok || !token.Valid {
		return "", errors.New("OTP token tidak valid")
	}

	if claims.OTP != otpInput {
		return "", errors.New("kode OTP salah")
	}

	// OTP valid → generate reset token (expire 15 menit)
	resetToken, err := generateResetToken(claims.Email)
	if err != nil {
		return "", err
	}

	return resetToken, nil
}

// generateResetToken — JWT khusus untuk step ganti password
func generateResetToken(email string) (string, error) {
	claims := ResetClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", errors.New("failed to sign reset token")
	}

	return signed, nil
}

// ParseResetToken — parse reset token untuk dapat email saat ganti password
func ParseResetToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &ResetClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return "", errors.New("reset token invalid atau sudah expired")
	}

	claims, ok := token.Claims.(*ResetClaims)
	if !ok || !token.Valid {
		return "", errors.New("reset token tidak valid")
	}

	return claims.Email, nil
}