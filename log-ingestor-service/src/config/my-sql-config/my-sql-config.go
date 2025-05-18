package mysqlconfig

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/yatender-pareek/log-ingestor-service/src/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Container struct {
	db *gorm.DB
	mu sync.RWMutex
}

var container *Container
var once sync.Once
var initErr error // Store initialization error

func Init() error {
	once.Do(func() {
		container = &Container{}
		container.db, initErr = newDB()
	})
	return initErr
}

func newDB() (*gorm.DB, error) {
	rootDSN := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
	)
	log.Println("======Mysql dsn=======", rootDSN)

	dbTemp, err := sql.Open("mysql", rootDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL server: %v", err)
	}
	defer dbTemp.Close()

	dbName := os.Getenv("MYSQL_DBNAME")
	_, err = dbTemp.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to create database %s: %v", dbName, err)
	}

	var dbCount int
	err = dbTemp.QueryRow("SELECT COUNT(*) FROM information_schema.schemata WHERE schema_name = ?", dbName).Scan(&dbCount)
	if err != nil {
		return nil, fmt.Errorf("failed to verify database %s existence: %v", dbName, err)
	}
	if dbCount == 0 {
		return nil, fmt.Errorf("database %s does not exist after creation attempt", dbName)
	}
	fmt.Printf("Database %s exists\n", dbName)

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		dbName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL database %s: %v", dbName, err)
	}

	var dbNameCheck string
	if err := db.Raw("SELECT DATABASE()").Scan(&dbNameCheck).Error; err != nil {
		return nil, fmt.Errorf("failed to verify active database: %v", err)
	}
	if dbNameCheck != dbName {
		return nil, fmt.Errorf("connected to wrong database: expected %s, got %s", dbName, dbNameCheck)
	}
	fmt.Printf("Successfully connected to database %s\n", dbNameCheck)

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	if err := db.AutoMigrate(models.GetAllModels()...); err != nil {
		return nil, fmt.Errorf("failed to migrate tables: %v", err)
	}

	return db, nil
}

func GetDB() *gorm.DB {
	if container == nil {
		panic("database container not initialized; call mysqlconfig.Init() first")
	}
	container.mu.RLock()
	defer container.mu.RUnlock()
	if container.db == nil {
		panic("database not initialized; Init() failed")
	}
	return container.db
}
