package mysqlconfig

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/yatender-pareek/threat-analyzer-service/src/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Container struct {
	db    *gorm.DB
	sqlDb *sql.DB
	mu    sync.RWMutex
}

var container *Container
var once sync.Once

func Init() error {
	var err error
	once.Do(func() {
		container = &Container{}
		container.db, container.sqlDb, err = newDB()
	})
	return err
}

func newDB() (*gorm.DB, *sql.DB, error) {
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
		return nil, nil, fmt.Errorf("failed to connect to MySQL server: %v", err)
	}
	defer dbTemp.Close()

	dbName := os.Getenv("MYSQL_DBNAME")
	_, err = dbTemp.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create database %s: %v", dbName, err)
	}

	var dbCount int
	err = dbTemp.QueryRow("SELECT COUNT(*) FROM information_schema.schemata WHERE schema_name = ?", dbName).Scan(&dbCount)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to verify database %s existence: %v", dbName, err)
	}
	if dbCount == 0 {
		return nil, nil, fmt.Errorf("database %s does not exist after creation attempt", dbName)
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
	// db, err := gorm.Open(mysql.Open(dsn),&gorm.Config{
	// Logger: logger.Default.LogMode(logger.Info)})
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to MySQL database %s: %v", dbName, err)
	}

	var dbNameCheck string
	if err := db.Raw("SELECT DATABASE()").Scan(&dbNameCheck).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to verify active database: %v", err)
	}
	if dbNameCheck != dbName {
		return nil, nil, fmt.Errorf("connected to wrong database: expected %s, got %s", dbName, dbNameCheck)
	}
	fmt.Printf("Successfully connected to database %s\n", dbNameCheck)

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	if err := db.AutoMigrate(models.GetAllModels()...); err != nil {
		return nil, nil, fmt.Errorf("failed to migrate tables: %v", err)
	}

	return db, sqlDB, nil
}

func GetDB() *gorm.DB {
	container.mu.RLock()
	defer container.mu.RUnlock()
	if container == nil || container.db == nil {
		panic("Container not initialized. Call Init first.")
	}
	return container.db
}

func GeSqltDB() *sql.DB {
	container.mu.RLock()
	defer container.mu.RUnlock()
	if container == nil || container.sqlDb == nil {
		panic("Container not initialized. Call Init first.")
	}
	return container.sqlDb
}
