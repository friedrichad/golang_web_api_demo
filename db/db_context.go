package db

import (
	"fmt"
	"log"
	"time"

	"github.com/friedrichad/golang_web_api_demo/configs"
	"github.com/friedrichad/golang_web_api_demo/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB là biến instance toàn cục để các repository trong project sử dụng
var DB *gorm.DB

// InitDB khởi tạo kết nối đến MySQL dựa trên cấu hình được load từ file json/env
func InitDB(cfg *config.Config) {
	var err error

	// 1. Xây dựng DSN (Data Source Name) từ struct config
	// parseTime=True rất quan trọng để mapping DATETIME của MySQL sang time.Time của Go
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)

	// 2. Mở kết nối GORM với cấu hình logger để dễ dàng debug SQL trong lúc phát triển
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Hiển thị các câu lệnh SQL trên console
	})

	if err != nil {
		log.Fatalf("Lỗi kết nối database: %v", err)
	}

	// 3. Cấu hình Connection Pool cho sql.DB phía dưới
	// Giúp tối ưu hóa hiệu năng và quản lý số lượng kết nối đồng thời
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Lỗi khi lấy sql.DB instance: %v", err)
	}

	// SetMaxIdleConns: Số lượng kết nối nhàn rỗi tối đa trong pool
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns: Số lượng kết nối mở tối đa đến database
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime: Thời gian tối đa một kết nối có thể được sử dụng lại
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 4. Auto Migration
	// GORM sẽ tự động tạo bảng hoặc thêm cột mới nếu struct Model có thay đổi
	// Thứ tự nên ưu tiên các bảng độc lập (bảng cha) trước
	err = DB.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Warehouse{},
		&models.Bin{},
		&models.ComponentCategory{},
		&models.Component{},
		&models.ComponentBin{},
		&models.Customer{},
		&models.TransferRequest{},
		&models.TransferRequestComponent{},
		&models.FinishedTransferRequestComponent{},
	)

	if err != nil {
		log.Fatalf("Lỗi trong quá trình AutoMigrate: %v", err)
	}

	fmt.Println("--- Database Context initialized successfully ---")
}