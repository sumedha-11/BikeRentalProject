package main

import (
	"bikerentalProject/internal/bike"
	"bikerentalProject/internal/booking"
	"bikerentalProject/internal/config"
	"bikerentalProject/internal/user"
	"bikerentalProject/pkg/sql"
	"flag"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

var configFile = flag.String("config", "config/common_lol.json", "config file")

func main() {
	flag.Parse()
	config.ReadConfig(*configFile)
	err := sql.DBConn(config.Config.DBConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer sql.DB.Close()
	sql.DB.AutoMigrate(&user.User{}, &booking.Booking{}, &bike.Bike{})
}
