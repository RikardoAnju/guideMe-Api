package models

import "time"

type GetUsersQuery struct {
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
	Search string `form:"search"` 
	Role   string `form:"role"`   
}

type UserListItem struct {
	ID            string    `json:"id"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	PhoneNumber   string    `json:"phoneNumber"`
	Role          string    `json:"role"`
	EmailVerified bool      `json:"emailVerified"`
	CreatedAt     time.Time `json:"createdAt"`
}

type PaginatedUsersResponse struct {
	Users      []UserListItem `json:"users"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"totalPages"`
}