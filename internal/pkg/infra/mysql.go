package infra

import (
	"log"
	"time"

	"github.com/tikivn/clickhousectl/internal/setting"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	model_mysql "github.com/tikivn/clickhousectl/internal/pkg/repo/mysql"
)

func NewMySqlSession(config *setting.MySql) (*gorm.DB, func(), error) {
	db, err := gorm.Open(mysql.Open(config.Addr()), &gorm.Config{})

	if err != nil {
		log.Fatalf("init mysql session failed: %v", err.Error())
		return nil, nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("init mysql pool failed: %v", err)
		return nil, nil, err
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(10)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	cleanup := func() {
		if err := sqlDB.Close(); err != nil {
			log.Fatalln(err)
		}
	}

	db.AutoMigrate(&model_mysql.ClickHouseServer{})
	db.AutoMigrate(&model_mysql.User{})

	return db, cleanup, nil
}
