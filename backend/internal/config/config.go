package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	JWTSecret   string
	Port        string
	CORSOrigins []string // comma-separated in env var CORS_ORIGIN
}

func Load() Config {
	raw := getEnv("CORS_ORIGIN", "http://localhost:5173")
	var origins []string
	for _, o := range strings.Split(raw, ",") {
		o = strings.TrimSpace(o)
		if o != "" {
			origins = append(origins, strings.TrimRight(o, "/"))
		}
	}
	return Config{
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "3306"),
		DBUser:      getEnv("DB_USER", "root"),
		DBPassword:  getEnv("DB_PASSWORD", ""),
		DBName:      getEnv("DB_NAME", "md_note"),
		JWTSecret:   getEnv("JWT_SECRET", "dev-secret-change-me"),
		Port:        getEnv("PORT", "8080"),
		CORSOrigins: origins,
	}
}

func (c Config) MySQLDSN() string {
	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
	// fmt.Println(conn)
	return conn
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
