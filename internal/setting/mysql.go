package setting

import (
	"fmt"
	"os"
)

func NewMySqlConfig() *MySql {
	return &MySql{
		MYSQL_USER: os.Getenv("MYSQL_USER"),
		MYSQL_PASS: os.Getenv("MYSQL_PASS"),
		MYSQL_HOST: os.Getenv("MYSQL_HOST"),
		MYSQL_PORT: os.Getenv("MYSQL_PORT"),
		MYSQL_DB:   os.Getenv("MYSQL_DB"),
	}
}

type MySql struct {
	MYSQL_USER string
	MYSQL_PASS string
	MYSQL_HOST string
	MYSQL_PORT string
	MYSQL_DB   string
}

func (c *MySql) Addr() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		c.MYSQL_USER, c.MYSQL_PASS, c.MYSQL_HOST, c.MYSQL_PORT, c.MYSQL_DB)
}
