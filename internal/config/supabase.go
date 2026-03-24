// internal/config/supabase.go
package config

import (
	"log"
	"os"

	supa "github.com/supabase-community/supabase-go"
)

var SupabaseClient *supa.Client
var SupabaseBucket string

func InitSupabase() {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")
	SupabaseBucket = os.Getenv("SUPABASE_BUCKET")

	client, err := supa.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		log.Fatalf("gagal inisialisasi Supabase: %v", err)
	}

	SupabaseClient = client
	log.Println("Supabase connected!")
}
