package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func NewDB() *sql.DB {
	user := GetEnv("DB_USER")
	pass := GetEnv("DB_PASSWORD")
	host := GetEnv("DB_HOST")
	port := GetEnv("DB_PORT")
	name := GetEnv("DB_NAME")

	log.Printf("DEBUG: Menghubungkan ke %s:%s sebagai %s", host, port, user)

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, name)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed conection in database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("database unreachable: %v", err)
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}

func NewRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: GetEnv("REDIS_ADDR"), // Akan mengambil 'redis_cache:6379' dari .env
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Println("Gagal koneksi Redis:", err)
	}

	return rdb
}
