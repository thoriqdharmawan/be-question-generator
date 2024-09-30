package db

import (
	"fmt"

	"github.com/thoriqdharmawan/be-question-generator/api/v1/models/entity"
	C "github.com/thoriqdharmawan/be-question-generator/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Postgre *gorm.DB

func GetPostgresURL() string {
	dbHost := C.Conf.PostgresHost
	dbPort := C.Conf.PostgresPort
	dbUser := C.Conf.PostgresUser
	dbPass := C.Conf.PostgresPassword
	dbName := C.Conf.PostgresDB

	if C.Conf.PostgresSSLMode == "disable" {
		return fmt.Sprintf("host=%s port=%s user=%s "+
			"password=%s dbname=%s sslmode=disable",
			dbHost, dbPort, dbUser, dbPass, dbName)
	} else {
		return fmt.Sprintf("host=%s port=%s user=%s "+
			"password=%s dbname=%s sslmode=%s sslrootcert=%s",
			dbHost, dbPort, dbUser, dbPass, dbName, C.Conf.PostgresSSLMode, C.Conf.PostgresRootCertLoc)
	}
}

func Init() error {
	db, err := gorm.Open(postgres.Open(GetPostgresURL()), &gorm.Config{})

	if err != nil {
		return fmt.Errorf("error opening database connection: %w", err)
	}

	db.AutoMigrate(&entity.User{})

	fmt.Println("Databse Migrated")

	Postgre = db

	return nil
}
