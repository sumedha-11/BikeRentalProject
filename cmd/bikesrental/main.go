package main

import (
	"bikerentalProject/internal/login"
	"bikerentalProject/pkg/logger"
	"flag"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"log"

	"bikerentalProject/internal/bike"
	"bikerentalProject/internal/config"
	"bikerentalProject/internal/user"
	"bikerentalProject/pkg/sql"
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

	sql.DB.SetLogger(logger.DBLogger)
	fmt.Println("DBCONNECTED")
	gin.DefaultWriter = logger.File
	router := gin.Default()

	token, err := login.RandToken(64)
	if err != nil {
		log.Fatal("unable to generate random token: ", err)
	}

	store := cookie.NewStore([]byte(token))
	store.Options(sessions.Options{
		Path:   "/",
		MaxAge: 86400 * 7,
	})

	router.Use(sessions.Sessions("bikerental", store))

	router.POST("/user/create", user.CreateUser)
	router.GET("/login", login.UserLogin)
	router.POST("/auth", login.AuthHandler)
	router.GET("/", bike.ViewAllBike)
	router.GET("/bike/view/:id", bike.ViewBike)

	authorized := router.Group("/")
	authorized.Use(login.AuthorizeRequest())
	{
		authorized.GET("/user/view/:id", user.ViewUser)
		authorized.GET("/logout", login.Logout)
		authorized.GET("/book/:Bmodel/:Booked/:Till/:BikeID", bike.AddBooking)
		authorized.GET("/booking/search", bike.SearchBikes)
		authorized.POST("/booking/available", bike.GetUnbooked)
	}
	adminAuth := router.Group("/")
	adminAuth.Use(login.AdminAuthRequest())
	{
		adminAuth.GET("/admin/bike/create", bike.CreateBikeForm)
		adminAuth.POST("/admin/bike/create", bike.CreateBike)
	}
	router.LoadHTMLGlob("templates/*")
	router.Static("/static/", "./static/")

	logger.InfoLogger.Printf("msg:%v", "server starting....")

	err = router.Run(":" + config.Config.Common.ServerPort)
	if err != nil {
		log.Fatalf("error in starting server , err:%v", err)
	}
}
