package main

import (
	"log"
	"os"
	"time"

	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	ID            string    `gorm:"primaryKey;type:varchar(255)"`
	Name          string    `gorm:"type:varchar(255)"`
	Email         string    `gorm:"type:varchar(255);uniqueIndex"`
	EmailVerified time.Time `gorm:"type:timestamptz"`
	Image         string    `gorm:"type:varchar(255)"`
	Accounts      []Account `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Sessions      []Session `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Posts         []Post    `gorm:"foreignKey:CreatedByID;constraint:OnDelete:CASCADE"`
}

type Post struct {
	ID          int    `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"type:varchar(255);index"`
	Content     string
	CreatedAt   time.Time `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"`
	CreatedBy   User      `gorm:"foreignKey:CreatedByID"`
	CreatedByID string    `gorm:"type:varchar(255)"`
}

type Account struct {
	ID                    string `gorm:"primaryKey;type:varchar(255)"`
	UserID                string `gorm:"type:varchar(255);index"`
	Type                  string `gorm:"type:varchar(255)"`
	Provider              string `gorm:"type:varchar(255)"`
	ProviderAccountID     string `gorm:"type:varchar(255)"`
	RefreshToken          string `gorm:"type:text"`
	AccessToken           string `gorm:"type:text"`
	ExpiresAt             int64
	TokenType             string `gorm:"type:varchar(255)"`
	Scope                 string `gorm:"type:varchar(255)"`
	IDToken               string `gorm:"type:text"`
	SessionState          string `gorm:"type:varchar(255)"`
	RefreshTokenExpiresIn int64
}

type Session struct {
	ID           string    `gorm:"primaryKey;type:varchar(255)"`
	SessionToken string    `gorm:"type:varchar(255);uniqueIndex"`
	UserID       string    `gorm:"type:varchar(255);index"`
	Expires      time.Time `gorm:"type:timestamptz"`
}

type VerificationToken struct {
	Identifier string    `gorm:"type:varchar(255);uniqueIndex:verificationtoken_identifier_token"`
	Token      string    `gorm:"primaryKey;type:varchar(255)"`
	Expires    time.Time `gorm:"type:timestamptz"`
}

func main() {
	POSTGRES_USER := os.Getenv("POSTGRES_USER")
	POSTGRES_PASSWORD := os.Getenv("POSTGRES_PASSWORD")
	POSTGRES_DATABASE := os.Getenv("POSTGRES_DATABASE")

	sysDSN := fmt.Sprintf(
		"host=db user=%s password=%s dbname=postgres port=5432 sslmode=disable",
		POSTGRES_USER, POSTGRES_PASSWORD,
	)
	sysDB, err := gorm.Open(postgres.Open(sysDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("Failed to connect to system database: %v", err)
	}

	if err := ensureDatabase(sysDB, POSTGRES_DATABASE); err != nil {
		log.Fatalf("Database creation failed: %v", err)
	}

	sqlDB, _ := sysDB.DB()
	sqlDB.Close()

	appDSN := fmt.Sprintf(
		"host=db user=%s password=%s dbname=%s port=5432 sslmode=disable",
		POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DATABASE,
	)
	db, err := gorm.Open(postgres.Open(appDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to application database: %v", err)
	}

	models := []any{
		&User{},
		&Post{},
		&Account{},
		&Session{},
		&VerificationToken{},
	}

	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			log.Fatalf("Failed to migrate %T: %v", model, err)
		}
	}

	log.Println("Database migration completed successfully")
}

func ensureDatabase(db *gorm.DB, dbName string) error {
	var exists bool
	err := db.Raw(
		"SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = ?)",
		dbName,
	).Scan(&exists).Error
	if err != nil {
		return err
	}

	if exists {
		log.Printf("Database %s already exists", dbName)
		return nil
	}

	if err := db.Exec("CREATE DATABASE " + dbName).Error; err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	log.Printf("Database %s created", dbName)
	return nil
}
