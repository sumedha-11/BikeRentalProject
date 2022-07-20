package booking

import (
	"bikerentalProject/pkg/sql"
	"github.com/jinzhu/gorm"
)

type Booking struct {
	gorm.Model
	Bmodel string `json:"Bmodel" gorm:"unique;not null"`
	Date   int64  `json:"Date"`
	UserID int    `json:"UserID"`
	BikeID int    `json:"BikeID"`
}

func (booking *Booking) Get() error {
	return sql.DB.Model(booking).First(booking).Error
}

func (booking *Booking) GetAll() ([]Booking, error) {
	var bookings []Booking
	err := sql.DB.Model(booking).Find(&bookings).Error
	return bookings, err
}

func (booking *Booking) GetFromWhere() (int, error) {
	var cnt int64
	err := sql.DB.Model(booking).Where("user_id =? AND bike_id =?", booking.UserID, booking.BikeID).Count(&cnt).Error
	return int(cnt), err
}
func (booking *Booking) Create() (err error) {
	err = sql.DB.Model(booking).Create(booking).Error
	return
}

func (booking *Booking) Update() (err error) {
	err = sql.DB.Model(booking).Updates(booking).Error
	return
}
