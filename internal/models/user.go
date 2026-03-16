package models

import "time"

type User struct {
	ID                   string     `json:"id" gorm:"primaryKey"`
	FirstName            string     `json:"firstName" gorm:"column:first_name"`
	LastName             string     `json:"lastName" gorm:"column:last_name"`
	Username             string     `json:"username" gorm:"column:username"`
	Email                string     `json:"email" gorm:"uniqueIndex"`
	Password             string     `json:"-" gorm:"column:password"`
	PhoneNumber          string     `json:"phoneNumber" gorm:"column:phone_number"`
	Gender               string     `json:"gender"`
	Address              string     `json:"address"`
	Role                 string     `json:"role" gorm:"default:user"`
	EmailVerified        bool       `json:"emailVerified" gorm:"column:email_verified;default:false"`
	PasswordReset        bool       `json:"passwordReset" gorm:"column:password_reset;default:false"`
	PasswordResetSuccess bool       `json:"passwordResetSuccess" gorm:"column:password_reset_success;default:false"`
	TemporaryResetToken  string     `json:"temporaryResetToken" gorm:"column:temporary_reset_token;default:''"`
	LastPasswordUpdate   *time.Time `json:"lastPasswordUpdate" gorm:"column:last_password_update"`
	CreatedAt            time.Time  `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
}