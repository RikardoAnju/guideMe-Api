// internal/models/event.go
package models

import "time"

type Event struct {
	ID             string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	CreatedAt      time.Time `json:"createdAt"`
	CreatedBy      string    `json:"createdBy" gorm:"type:varchar(255)"`
	Deskripsi      string    `json:"deskripsi" gorm:"type:text"`
	HargaTiket     int64     `json:"hargaTiket"`
	ImageURL       string    `json:"imageUrl" gorm:"type:text"`
	IsFree         bool      `json:"isFree"`
	Kategori       string    `json:"kategori" gorm:"type:varchar(100)"`
	Lokasi         string    `json:"lokasi" gorm:"type:varchar(255)"`
	NamaEvent      string    `json:"namaEvent" gorm:"type:varchar(255)"`
	TanggalMulai   string    `json:"tanggalMulai" gorm:"type:varchar(20)"`
	TanggalSelesai string    `json:"tanggalSelesai" gorm:"type:varchar(20)"`
	UpdatedAt      time.Time `json:"updatedAt"`
	URLMaps        string    `json:"urlMaps" gorm:"type:text"`
	WaktuMulai     string    `json:"waktuMulai" gorm:"type:varchar(10)"`
	WaktuSelesai   string    `json:"waktuSelesai" gorm:"type:varchar(10)"`
}

type CreateEventRequest struct {
	Deskripsi      string `json:"deskripsi" binding:"required"`
	HargaTiket     int64  `json:"hargaTiket"`
	ImageURL       string `json:"imageUrl" binding:"required"`
	IsFree         bool   `json:"isFree"`
	Kategori       string `json:"kategori" binding:"required"`
	Lokasi         string `json:"lokasi" binding:"required"`
	NamaEvent      string `json:"namaEvent" binding:"required"`
	TanggalMulai   string `json:"tanggalMulai" binding:"required"`
	TanggalSelesai string `json:"tanggalSelesai" binding:"required"`
	URLMaps        string `json:"urlMaps"`
	WaktuMulai     string `json:"waktuMulai" binding:"required"`
	WaktuSelesai   string `json:"waktuSelesai" binding:"required"`
}

type UpdateEventRequest struct {
	Deskripsi      string `json:"deskripsi"`
	HargaTiket     int64  `json:"hargaTiket"`
	ImageURL       string `json:"imageUrl"`
	IsFree         bool   `json:"isFree"`
	Kategori       string `json:"kategori"`
	Lokasi         string `json:"lokasi"`
	NamaEvent      string `json:"namaEvent"`
	TanggalMulai   string `json:"tanggalMulai"`
	TanggalSelesai string `json:"tanggalSelesai"`
	URLMaps        string `json:"urlMaps"`
	WaktuMulai     string `json:"waktuMulai"`
	WaktuSelesai   string `json:"waktuSelesai"`
}