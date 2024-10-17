package mysql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	DBName   string `json:"db_name"`
}

func WithClient(hand func(db *gorm.DB) error) error {
	conf := MysqlConfig{
		User:     "root",
		Password: "",
		Host:     "127.0.0.1",
		Port:     3306,
		DBName:   "test",
	}

	dsn := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		conf.User,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.DBName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	defer func() {
		dd, _ := db.DB()
		dd.Close()
	}()
	return hand(db)
}

func WithTx(hand func(db *gorm.DB) error) error {
	return WithClient(func(db *gorm.DB) error {
		tx := db.Begin()
		err := hand(tx)
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
		return nil
	})
}
