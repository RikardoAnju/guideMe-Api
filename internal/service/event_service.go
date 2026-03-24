// internal/service/event_service.go
package service

import (
	"errors"
	"time"

	"guide-me/internal/config"
	"guide-me/internal/models"

	"github.com/google/uuid"
)

func GetAllEvent() ([]models.Event, error) {
	var events []models.Event
	result := config.DB.Find(&events)
	if result.Error != nil {
		return nil, errors.New("failed to fetch events")
	}
	return events, nil
}

func GetEventByID(id string) (*models.Event, error) {
	var event models.Event
	result := config.DB.Where("id = ?", id).First(&event)
	if result.Error != nil {
		return nil, errors.New("event tidak ditemukan")
	}
	return &event, nil
}

func CreateEvent(req models.CreateEventRequest, createdBy string) (*models.Event, error) {
	event := &models.Event{
		ID:             uuid.New().String(),
		CreatedAt:      time.Now(),
		CreatedBy:      createdBy,
		Deskripsi:      req.Deskripsi,
		HargaTiket:     req.HargaTiket,
		ImageURL:       req.ImageURL,
		IsFree:         req.IsFree,
		Kategori:       req.Kategori,
		Lokasi:         req.Lokasi,
		NamaEvent:      req.NamaEvent,
		TanggalMulai:   req.TanggalMulai,
		TanggalSelesai: req.TanggalSelesai,
		UpdatedAt:      time.Now(),
		URLMaps:        req.URLMaps,
		WaktuMulai:     req.WaktuMulai,
		WaktuSelesai:   req.WaktuSelesai,
	}

	result := config.DB.Create(event)
	if result.Error != nil {
		return nil, result.Error
	}

	return event, nil
}

func UpdateEvent(id string, req models.UpdateEventRequest) (*models.Event, error) {
	var event models.Event
	if err := config.DB.Where("id = ?", id).First(&event).Error; err != nil {
		return nil, errors.New("event tidak ditemukan")
	}

	if err := config.DB.Model(&event).Updates(map[string]interface{}{
		"deskripsi":       req.Deskripsi,
		"harga_tiket":     req.HargaTiket,
		"image_url":       req.ImageURL,
		"is_free":         req.IsFree,
		"kategori":        req.Kategori,
		"lokasi":          req.Lokasi,
		"nama_event":      req.NamaEvent,
		"tanggal_mulai":   req.TanggalMulai,
		"tanggal_selesai": req.TanggalSelesai,
		"url_maps":        req.URLMaps,
		"waktu_mulai":     req.WaktuMulai,
		"waktu_selesai":   req.WaktuSelesai,
		"updated_at":      time.Now(),
	}).Error; err != nil {
		return nil, errors.New("failed to update event")
	}

	return GetEventByID(id)
}

func DeleteEvent(id string) error {
	result := config.DB.Where("id = ?", id).Delete(&models.Event{})
	if result.Error != nil {
		return errors.New("failed to delete event")
	}
	if result.RowsAffected == 0 {
		return errors.New("event tidak ditemukan")
	}
	return nil
}