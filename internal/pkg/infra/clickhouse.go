package infra

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewClickHouseConnection(host string, port string) *clickHouseClientImpl {
	return &clickHouseClientImpl{
		host:     host,
		port:     port,
		username: "default",
		password: "",
	}
}

type clickHouseClientImpl struct {
	host     string
	port     string
	username string
	password string
}

func (c *clickHouseClientImpl) endpoint() string {
	if c.username != "" && c.password != "" {
		return fmt.Sprintf("tcp://%s:%s?username=%s&password=%s&debug=true", c.host, c.port, c.username, c.password)
	} else if c.username != "" {
		return fmt.Sprintf("tcp://%s:%s?username=%s&debug=true", c.host, c.port, c.username)
	} else {
		return fmt.Sprintf("tcp://%s:%s?debug=true", c.host, c.port)
	}
}

func (c *clickHouseClientImpl) WithAuthentication(username string, password string) *clickHouseClientImpl {
	c.username = username
	c.password = password
	return c
}

func NewClickhouseSession(client *clickHouseClientImpl) (*gorm.DB, func(), error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
		},
	)

	db, err := gorm.Open(clickhouse.Open(client.endpoint()), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Fatalf("init clickhouse session failed: %v", err.Error())
		return nil, nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("init clickhouse pool failed: %v", err)
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

	return db, cleanup, nil
}
