package service

import (
	"errors"
	"math"
	"strings"

	"guide-me/internal/config"
	"guide-me/internal/models"
)

func ToggleUserActive(id string, callerID string) (*models.UserListItem, error) {
	if id == callerID {
		return nil, errors.New("tidak bisa menonaktifkan akun sendiri")
	}

	var u models.User
	if err := config.DB.Where("id = ?", id).First(&u).Error; err != nil {
		return nil, errors.New("user not found")
	}

	if u.Role == "admin" {
		return nil, errors.New("tidak bisa menonaktifkan sesama admin")
	}

	newStatus := !u.IsActive
	if err := config.DB.Model(&u).Update("is_active", newStatus).Error; err != nil {
		return nil, err
	}

	u.IsActive = newStatus
	item := toUserListItem(u)
	return &item, nil
}

func GetAllUsers(query models.GetUsersQuery) (*models.PaginatedUsersResponse, error) {
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Limit <= 0 {
		query.Limit = 10
	}
	offset := (query.Page - 1) * query.Limit

	db := config.DB.Model(&models.User{})

	if query.Search != "" {
		s := "%" + strings.ToLower(query.Search) + "%"
		db = db.Where(
			"LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ? OR LOWER(username) LIKE ? OR LOWER(email) LIKE ?",
			s, s, s, s,
		)
	}

	if query.Role != "" {
		db = db.Where("role = ?", query.Role)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	var users []models.User
	if err := db.
		Offset(offset).
		Limit(query.Limit).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		return nil, err
	}

	items := make([]models.UserListItem, len(users))
	for i, u := range users {
		items[i] = toUserListItem(u) // ← pakai toUserListItem agar konsisten
	}

	return &models.PaginatedUsersResponse{
		Users:      items,
		Total:      total,
		Page:       query.Page,
		Limit:      query.Limit,
		TotalPages: int(math.Ceil(float64(total) / float64(query.Limit))),
	}, nil
}

func GetUserByID(id string) (*models.UserListItem, error) {
	var u models.User
	if err := config.DB.Where("id = ?", id).First(&u).Error; err != nil {
		return nil, errors.New("user not found")
	}

	item := toUserListItem(u)
	return &item, nil
}

// ── Helper ────────────────────────────────────────────────────────────────────

func toUserListItem(u models.User) models.UserListItem {
	return models.UserListItem{
		ID:            u.ID,
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		Username:      u.Username,
		Email:         u.Email,
		PhoneNumber:   u.PhoneNumber,
		Role:          u.Role,
		EmailVerified: u.EmailVerified,
		CreatedAt:     u.CreatedAt,
		IsActive:      u.IsActive, // ← tambahan
	}
}