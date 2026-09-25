package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// ConnectDatabase membuat connection pool ke PostgreSQL.
func ConnectDatabase() (*pgxpool.Pool, error) {
	// Membaca file .env.
	err := godotenv.Load()
	if err != nil {
		fmt.Println("⚠️  .env tidak ditemukan, menggunakan environment system")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	sslMode := os.Getenv("DB_SSLMODE")

	if host == "" ||
		port == "" ||
		user == "" ||
		password == "" ||
		dbName == "" {
		return nil, fmt.Errorf("konfigurasi database belum lengkap")
	}

	// Format koneksi PostgreSQL.
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host,
		port,
		user,
		password,
		dbName,
		sslMode,
	)

	// Batasi waktu koneksi awal selama 5 detik.
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	pool, err := pgxpool.New(
		ctx,
		connectionString,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"gagal membuat PostgreSQL pool: %w",
			err,
		)
	}

	// Ping memastikan database benar-benar bisa diakses.
	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()

		return nil, fmt.Errorf(
			"gagal terhubung ke PostgreSQL: %w",
			err,
		)
	}

	return pool, nil
}
