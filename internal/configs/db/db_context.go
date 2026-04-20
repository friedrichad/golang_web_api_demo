package db

import (
	"fmt"
	"log"
	"time"

	config "github.com/friedrichad/golang_web_api_demo/internal/configs"
	"github.com/friedrichad/golang_web_api_demo/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB là biến instance toàn cục để các repository trong project sử dụng
var DB *gorm.DB

func ConnectDB(cfg *config.Config) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Không thể kết nối database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Không thể lấy sqlDB: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db

	fmt.Println("Database connected successfully")
}

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
		// Nhóm 1: Không phụ thuộc ai
		&models.Role{},
		&models.Warehouse{},
		&models.Componentcategory{},
		&models.Component{},
		&models.Customer{},

		// Nhóm 2: Phụ thuộc nhóm 1
		&models.User{},            // Many-to-Many với Role thông qua u_r
		&models.Bin{},             // Foreign Key đến Warehouse
		&models.Transferrequest{}, // Foreign Keys đến User, Warehouse, Customer

		// Nhóm 3: Bảng chi tiết giao dịch (Phụ thuộc nhóm 2)
		// GORM sẽ tự động tạo bảng trung gian cho many2many relationships:
		// - u_r (User-Role)
		// - c_c (Component-ComponentCategory)
		// - component_bins (Component-Bin)
		&models.Transferrequestcomponent{},
		&models.Finishedtransferrequestcomponent{},
	)
	if err != nil {
		log.Fatalf("Lỗi trong quá trình AutoMigrate: %v", err)
	}

	fmt.Println("--- Database Context initialized successfully ---")
}
