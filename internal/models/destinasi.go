// internal/models/destinasi.go
package models

import "time"

type Destinasi struct {
	ID            string    `json:"id"            gorm:"primaryKey;type:uuid"`
	CreatedAt     time.Time `json:"createdAt"`
	CreatedBy     string    `json:"createdBy"`
	Deskripsi     string    `json:"deskripsi"`
	HargaTiket    int64     `json:"hargaTiket"`
	ImageURL      string    `json:"imageUrl"`
	IsFree        bool      `json:"isFree"`
	JamBuka       string    `json:"jamBuka"`
	JamTutup      string    `json:"jamTutup"`
	Kategori      string    `json:"kategori"`
	Lokasi        string    `json:"lokasi"`
	NamaDestinasi string    `json:"namaDestinasi"`
	Rating        float64   `json:"rating"`
	RatingCount   int64     `json:"ratingCount"`
	UpdatedAt     time.Time `json:"updatedAt"`
	URLMaps       string    `json:"urlMaps"`
}

// CreateDestinasiRequest digunakan saat menambahkan destinasi baru.
type CreateDestinasiRequest struct {
	Deskripsi     string `json:"deskripsi"     binding:"required"`
	HargaTiket    int64  `json:"hargaTiket"`
	ImageURL      string `json:"imageUrl"`
	IsFree        bool   `json:"isFree"`
	JamBuka       string `json:"jamBuka"       binding:"required"`
	JamTutup      string `json:"jamTutup"      binding:"required"`
	Kategori      string `json:"kategori"      binding:"required"`
	Lokasi        string `json:"lokasi"        binding:"required"`
	NamaDestinasi string `json:"namaDestinasi" binding:"required"`
	URLMaps       string `json:"urlMaps"`
}

// UpdateDestinasiRequest digunakan saat mengupdate data destinasi (semua field opsional).
type UpdateDestinasiRequest struct {
	Deskripsi     string `json:"deskripsi"`
	HargaTiket    int64  `json:"hargaTiket"`
	ImageURL      string `json:"imageUrl"`
	IsFree        bool   `json:"isFree"`
	JamBuka       string `json:"jamBuka"`
	JamTutup      string `json:"jamTutup"`
	Kategori      string `json:"kategori"`
	Lokasi        string `json:"lokasi"`
	NamaDestinasi string `json:"namaDestinasi"`
	URLMaps       string `json:"urlMaps"`
}

// RatingRequest digunakan saat user memberikan rating pada destinasi.
type RatingRequest struct {
	Rating float64 `json:"rating" binding:"required,min=1,max=5"`
}
