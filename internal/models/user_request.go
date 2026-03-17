// internal/models/auth.go

package models

type RegisterRequest struct {
    FirstName   string `json:"firstName" binding:"required"`
    LastName    string `json:"lastName" binding:"required"`
    Username    string `json:"username" binding:"required"`
    Email       string `json:"email" binding:"required,email"`
    Password    string `json:"password" binding:"required,min=6"`
    PhoneNumber string `json:"phoneNumber"`
    Gender      string `json:"gender"`
    Address     string `json:"address"`
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type UpdateProfileRequest struct {
    FirstName   string `json:"firstName"`
    LastName    string `json:"lastName"`
    Username    string `json:"username"`
    PhoneNumber string `json:"phoneNumber"`
    Gender      string `json:"gender"`
    Address     string `json:"address"`
}

// Request kirim OTP ke email
type ResetPasswordRequest struct {
    Email string `json:"email" binding:"required,email"`
}

// Request verifikasi OTP — BARU
type VerifyOTPRequest struct {
    OTP      string `json:"otp" binding:"required,len=6"`
    OTPToken string `json:"otp_token" binding:"required"`
}

// Request ganti password — update field Token -> ResetToken & NewPassword min 8
type ChangePasswordRequest struct {
    ResetToken  string `json:"reset_token" binding:"required"`
    NewPassword string `json:"newPassword" binding:"required,min=8"`
}