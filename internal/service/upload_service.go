package service

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	storage_go "github.com/supabase-community/storage-go"

	"guide-me/internal/config"
)

func UploadImage(file multipart.File, header *multipart.FileHeader) (string, error) {
	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	// Read file safely
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("gagal membaca file: %w", err)
	}

	// Detect Content-Type
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}

	// Upload to Supabase Storage
	reader := bytes.NewReader(fileBytes)
	upsert := false

	_, err = config.SupabaseClient.Storage.UploadFile(
		config.SupabaseBucket,
		fileName,
		reader,
		storage_go.FileOptions{
			ContentType: &contentType,
			Upsert:      &upsert,
		},
	)
	if err != nil {
		return "", fmt.Errorf("gagal upload ke Supabase: %w", err)
	}

	// Build public URL
	supabaseURL := os.Getenv("SUPABASE_URL")
	if supabaseURL == "" {
		return "", fmt.Errorf("SUPABASE_URL environment variable tidak ditemukan")
	}

	publicURL := fmt.Sprintf(
		"%s/storage/v1/object/public/%s/%s",
		supabaseURL,
		config.SupabaseBucket,
		fileName,
	)

	return publicURL, nil
}