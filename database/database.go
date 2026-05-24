package database

import (
	"fmt"

	"github.com/GemaSatya/LaundrySystem/env"
	"github.com/GemaSatya/LaundrySystem/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase(){

	connectionUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		env.GetEnv("DB_USERNAME"),
		env.GetEnv("DB_PASSWORD"),
		env.GetEnv("DB_HOST"),
		env.GetEnv("DB_PORT"),
		env.GetEnv("DB_NAME"))
	database, err := gorm.Open(mysql.Open(connectionUrl))
	if err != nil{
		panic(err)
	}

	database.AutoMigrate(&model.User{}, &model.Report{}, &model.Payment{}, &model.OrderDetail{}, &model.Order{}, &model.Customer{})

	DB = database
}