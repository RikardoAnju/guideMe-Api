// internal/service/destinasi_service.go
package service

import (
	"errors"
	"time"

	"guide-me/internal/config"
	"guide-me/internal/models"

	"github.com/google/uuid"
)

func GetAllDestinasi() ([]models.Destinasi, error) {
	var destinasi []models.Destinasi
	result := config.DB.Find(&destinasi)
	if result.Error != nil {
		return nil, errors.New("failed to fetch destinasi")
	}
	return destinasi, nil
}

func GetDestinasiByID(id string) (*models.Destinasi, error) {
	var destinasi models.Destinasi
	result := config.DB.Where("id = ?", id).First(&destinasi)
	if result.Error != nil {
		return nil, errors.New("destinasi tidak ditemukan")
	}
	return &destinasi, nil
}

func CreateDestinasi(req models.CreateDestinasiRequest, createdBy string) (*models.Destinasi, error) {
	destinasi := &models.Destinasi{
		ID:            uuid.New().String(),
		CreatedAt:     time.Now(),
		CreatedBy:     createdBy,
		Deskripsi:     req.Deskripsi,
		HargaTiket:    req.HargaTiket,
		ImageURL:      req.ImageURL,
		IsFree:        req.IsFree,
		JamBuka:       req.JamBuka,
		JamTutup:      req.JamTutup,
		Kategori:      req.Kategori,
		Lokasi:        req.Lokasi,
		NamaDestinasi: req.NamaDestinasi,
		Rating:        0,
		RatingCount:   0,
		UpdatedAt:     time.Now(),
		URLMaps:       req.URLMaps,
	}

	result := config.DB.Create(destinasi)
	if result.Error != nil {
		return nil, result.Error
	}

	return destinasi, nil
}

func UpdateDestinasi(id string, req models.UpdateDestinasiRequest) (*models.Destinasi, error) {
	var destinasi models.Destinasi
	if err := config.DB.Where("id = ?", id).First(&destinasi).Error; err != nil {
		return nil, errors.New("destinasi tidak ditemukan")
	}

	if err := config.DB.Model(&destinasi).Updates(map[string]interface{}{
		"deskripsi":      req.Deskripsi,
		"harga_tiket":    req.HargaTiket,
		"image_url":      req.ImageURL,
		"is_free":        req.IsFree,
		"jam_buka":       req.JamBuka,
		"jam_tutup":      req.JamTutup,
		"kategori":       req.Kategori,
		"lokasi":         req.Lokasi,
		"nama_destinasi": req.NamaDestinasi,
		"url_maps":       req.URLMaps,
		"updated_at":     time.Now(),
	}).Error; err != nil {
		return nil, errors.New("failed to update destinasi")
	}

	return GetDestinasiByID(id)
}

func DeleteDestinasi(id string) error {
	result := config.DB.Where("id = ?", id).Delete(&models.Destinasi{})
	if result.Error != nil {
		return errors.New("failed to delete destinasi")
	}
	if result.RowsAffected == 0 {
		return errors.New("destinasi tidak ditemukan")
	}
	return nil
}

func SubmitRating(id string, newRating float64) (*models.Destinasi, error) {
	var destinasi models.Destinasi
	if err := config.DB.Where("id = ?", id).First(&destinasi).Error; err != nil {
		return nil, errors.New("destinasi tidak ditemukan")
	}

	totalScore := destinasi.Rating*float64(destinasi.RatingCount) + newRating
	destinasi.RatingCount++
	destinasi.Rating = totalScore / float64(destinasi.RatingCount)
	destinasi.UpdatedAt = time.Now()

	if err := config.DB.Model(&destinasi).Updates(map[string]interface{}{
		"rating":       destinasi.Rating,
		"rating_count": destinasi.RatingCount,
		"updated_at":   destinasi.UpdatedAt,
	}).Error; err != nil {
		return nil, errors.New("failed to submit rating")
	}

	return &destinasi, nil
}
