package database

import (
	"fmt"
	"github.com/go-errors/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

// GetGormConnection Возвращает клиент для работы с БД
func GetGormConnection(dbConfig DbConfig) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch dbConfig.Driver {
	case Postgres:
		dialector = postgres.Open(getConnectionString(dbConfig))
	}

	if dialector == nil {
		return nil, errors.New("Unknown db driver: " + dbConfig.Driver)
	}

	config := &gorm.Config{}

	if dbConfig.Logging {
		config.Logger = logger.Default.LogMode(logger.Info)
	}

	//if dbConfig.DisableAutomaticPing {
	config.DisableAutomaticPing = true
	//}

	gormDb, err := gorm.Open(dialector, config)

	if err != nil {
		return nil, err
	}

	db, err := gormDb.DB()

	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Hour)

	if dbConfig.MaxOpenConnections > 0 {
		db.SetMaxOpenConns(dbConfig.MaxOpenConnections)
	}

	if dbConfig.MaxIdleConnections > 0 {
		db.SetMaxIdleConns(dbConfig.MaxIdleConnections)
	}

	return gormDb, nil
}

// getConnectionString Возвращает строку (DSN) для создания соединения с БД
func getConnectionString(dbConfig DbConfig) string {
	sslMode := "disable"

	if dbConfig.SslMode {
		sslMode = "enable"
	}

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Db, sslMode)
}
