package dal

import (
	"fmt"

	"kp-management/internal/pkg/conf"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	d *gorm.DB
)

const dsnTemplate = "%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local"

func MustInitMySQL() {
	c := conf.Conf
	dsn := fmt.Sprintf(dsnTemplate, c.MySQL.Username, c.MySQL.Passport, c.MySQL.IP, c.MySQL.Port, c.MySQL.DBName, c.MySQL.Charset)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("fatal error db init: %w", err))
	}

	d = db

	fmt.Println("mysql initialized")
}

func DB() *gorm.DB {
	return d
}
