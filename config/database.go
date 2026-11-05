package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDatabase() {
	dbUser := getEnv("DB_USER", "root")
	dbPass := getEnv("DB_PASS", "root")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME", "course_management")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		dbUser, dbPass, dbHost, dbPort, dbName)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(10)
	DB.SetConnMaxLifetime(5 * time.Minute)

	if err = DB.Ping(); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	log.Println("✅ Database connected successfully")
	createTables()
}

func createTables() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			user_id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			email VARCHAR(150) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			role ENUM('ADMIN','LEARNER') DEFAULT 'LEARNER',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS courses (
			course_id INT AUTO_INCREMENT PRIMARY KEY,
			course_name VARCHAR(100) NOT NULL,
			description TEXT,
			duration VARCHAR(50),
			fees DECIMAL(10,2) DEFAULT 0.00,
			is_active BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS subjects (
			subject_id INT AUTO_INCREMENT PRIMARY KEY,
			subject_name VARCHAR(100) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS course_subject (
			course_subject_id INT AUTO_INCREMENT PRIMARY KEY,
			course_id INT NOT NULL,
			subject_id INT NOT NULL,
			FOREIGN KEY (course_id) REFERENCES courses(course_id) ON DELETE CASCADE,
			FOREIGN KEY (subject_id) REFERENCES subjects(subject_id) ON DELETE CASCADE,
			UNIQUE KEY unique_course_subject (course_id, subject_id)
		)`,
		`CREATE TABLE IF NOT EXISTS enrollments (
			enrollment_id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT NOT NULL,
			course_id INT NOT NULL,
			enrollment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			status VARCHAR(50) DEFAULT 'ACTIVE',
			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
			FOREIGN KEY (course_id) REFERENCES courses(course_id) ON DELETE CASCADE,
			UNIQUE KEY unique_enrollment (user_id, course_id)
		)`,
		`CREATE TABLE IF NOT EXISTS audit_logs (
			log_id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT,
			action VARCHAR(255),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, q := range queries {
		if _, err := DB.Exec(q); err != nil {
			log.Fatalf("Failed to create table: %v\nQuery: %s", err, q)
		}
	}

	log.Println("✅ All tables created/verified")
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
