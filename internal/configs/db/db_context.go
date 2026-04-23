package db

import (
	"fmt"
	"log"
	// "time"

	config "github.com/friedrichad/golang_web_api_demo/internal/configs"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	// "log"
)

var Instance *gorm.DB
var err error

func InitMysql() {
	var connectionString = viper.GetString("database.mysql.connection-string")
	Instance, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
}

func GenmodelFromDB(cfg *config.Config) {
	// 1. Xây dựng DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)

	// 2. Kết nối DB
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Không thể kết nối database: %v", err)
	}

	// 3. Khởi tạo generator
	g := gen.NewGenerator(gen.Config{
		OutPath: "internal/query", // thư mục sinh code
		Mode:    gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	g.UseDB(db)

	// 4. Sinh code cho tất cả bảng
	g.GenerateAllTable()

	// 5. Thực thi
	g.Execute()

	fmt.Println("--- model generated successfully from database ---")
}
