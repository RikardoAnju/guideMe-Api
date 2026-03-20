package models

import "time"

type TokenBlacklist struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Token     string    `gorm:"column:token;uniqueIndex"`
	ExpiredAt time.Time `gorm:"column:expired_at"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}