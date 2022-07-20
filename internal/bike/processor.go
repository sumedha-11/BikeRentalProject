package bike

import (
	"bikerentalProject/internal/booking"
	"bikerentalProject/internal/login"
	"bikerentalProject/pkg/logger"
	"bikerentalProject/pkg/sql"
	"fmt"
	//"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func ViewBike(g *gin.Context) {
	var id int
	var err error
	ids := g.Param("id")
	id, err = strconv.Atoi(ids)
	if err != nil {
		logger.ErrorLogger.Printf("error in converting param id to integer, ids:%v err:%v", ids, err)
		g.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"message": "Internal Error."})
		return
	}
	b := &Bike{}
	b.ID = uint(id)
	err = b.Get()
	if err != nil {
		logger.ErrorLogger.Printf("error in Bike Get method, cId:%d err:%v", id, err)
		g.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"message": "Internal Error."})
		return
	}

	//TODO: hide form post url which include company id
	s := login.Isloggedin(g)

	g.HTML(http.StatusOK, "BikeView.html", gin.H{"Bike": b, "Bookings": b.Bookings, "Mail": s})
}

func ViewAllBike(g *gin.Context) {
	b := Bike{}
	cs, err := b.GetAll()
	if err != nil {
		logger.ErrorLogger.Printf("error in getting all bike get method, err:%v", err)
		g.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"message": "Internal Error."})
		return
	}
	s := login.Isloggedin(g)
	g.HTML(http.StatusOK, "bikeAll.tmpl", gin.H{"List": cs, "Mail": s})
}
func UpdateBike(g *gin.Context) {
	var id int
	var err error
	ids := g.Param("id")
	id, err = strconv.Atoi(ids)
	if err != nil {
		logger.ErrorLogger.Printf("error in converting parmas ids to integer, Id:%d, err:%v", id, err)
		g.JSON(http.StatusInternalServerError, gin.H{
			"data":     nil,
			"debugMsg": "Error converting id to integer",
		})
		return
	}
	c := &Bike{}
	err = g.ShouldBindJSON(c)
	if err != nil {
		logger.ErrorLogger.Printf("invalid json provided, Id:%d, err:%v", id, err)
		g.JSON(http.StatusUnprocessableEntity, gin.H{
			"data":     nil,
			"debugMsg": "Invalid json provided",
		})
		return
	}
	c.ID = uint(id)
	err = c.Update()
	if err != nil {
		logger.ErrorLogger.Printf("error in updating company, Id:%d, err:%v", id, err)
		g.JSON(http.StatusInternalServerError, gin.H{
			"data":     nil,
			"debugMsg": "Error in updating company",
		})
		return
	}
	g.JSON(http.StatusOK, gin.H{
		"data":     c,
		"debugMsg": "",
	})
}

func CreateBikeForm(g *gin.Context) {
	g.HTML(http.StatusOK, "BikeAdd.tmpl", gin.H{})
}
func CreateBike(g *gin.Context) {
	b := &Bike{
		Bmodel: g.PostForm("Bmodel"),
		Price:  g.PostForm("Price"),
	}
	err := b.Create()
	if err != nil {
		logger.ErrorLogger.Printf("error in creating bike, err:%v", err)
		g.JSON(http.StatusInternalServerError, gin.H{
			"data":     nil,
			"debugMsg": "Error in creating bike",
		})
		return
	}
	g.JSON(http.StatusOK, gin.H{
		"data":     b,
		"debugMsg": "",
	})
	return
}
func SearchBikes(g *gin.Context) {
	s := login.Isloggedin(g)
	g.HTML(http.StatusOK, "search.tmpl", gin.H{"Mail": s})
}
func GetUnbooked(g *gin.Context) {
	u := &BookDates{StartDate: g.PostForm("StartDate"),
		EndDate: g.PostForm("EndDate")}
	var booked, till int64
	t, _ := time.Parse("2006-01-02", u.StartDate)
	v, _ := time.Parse("2006-01-02", u.EndDate)
	booked = t.Unix()
	till = v.Unix()
	fmt.Println(booked)
	var rows []booking.Booking
	b := Bike{}
	cs, errr := b.GetAll()
	if errr != nil {
		logger.ErrorLogger.Printf("error in getting all bike get method, err:%v", errr)
		g.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"message": "Internal Error."})
		return
	}
	err := sql.DB.Table("bookings").Where("Date>=? AND Date<=?", booked, till).Find(&rows).Error
	if err != nil {
		fmt.Println(err)
		logger.ErrorLogger.Printf("error in getting booked bike get method, err:%v", err)
		g.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"message": "Internal Error"})
		return
	}
	m := make(map[string]int)
	for j := 0; j < len(rows); j++ {
		m[rows[j].Bmodel] = 1
	}
	var res []Bookit
	for j := 0; j < len(cs); j++ {
		if m[cs[j].Bmodel] != 1 {
			var v Bookit
			v.BikeID = cs[j].ID
			v.Booked = u.StartDate
			v.Till = u.EndDate
			v.Bmodel = cs[j].Bmodel
			v.Price = cs[j].Price
			res = append(res, v)
		}
	}
	s := login.Isloggedin(g)
	g.HTML(http.StatusOK, "bookbike.tmpl", gin.H{"List": res, "Mail": s})
}

func AddBooking(g *gin.Context) {

	Bmodel := g.Param("Bmodel")
	Booked := g.Param("Booked")
	Till := g.Param("Till")
	ID := g.Param("BikeID")
	id, errrr := strconv.Atoi(ID)
	//fmt.Println(ID)
	var booked int64
	var till int64
	t, e := time.Parse("2006-01-02", Booked)
	if e != nil {
		fmt.Println(e)
	}
	v, errrr := time.Parse("2006-01-02", Till)
	if errrr != nil {
		fmt.Println(errrr)
	}
	booked = t.Unix()
	till = v.Unix()
	var j int64
	for j = booked; j <= till; j = j + 86400 {
		B := &booking.Booking{}
		B.Bmodel = Bmodel
		B.Date = j
		B.BikeID = id
		err := B.Create()
		if err != nil {
			logger.ErrorLogger.Printf("error in creating booking, err:%v", err)
			g.JSON(http.StatusInternalServerError, gin.H{
				"data":     nil,
				"debugMsg": "Error in creating booking",
			})
			return
		}
	}
	g.HTML(http.StatusOK, "booked.tmpl", gin.H{})
}
