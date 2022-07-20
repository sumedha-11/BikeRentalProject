package bike

import (
	"bikerentalProject/internal/booking"
	"bikerentalProject/pkg/sql"
	"github.com/jinzhu/gorm"
)

type Bike struct {
	gorm.Model
	Bmodel   string            `json:"name" gorm:"unique;not null"`
	Price    string            `json:"price"`
	Bookings []booking.Booking `gorm:"foreignkey:BikeID"`
}

type Bookit struct {
	gorm.Model
	Bmodel   string            `json:"name" gorm:"unique;not null"`
	Price    string            `json:"price"`
	Booked   string            `json:"Booked"`
	Till     string            `json:"Till"`
	BikeID   uint              `json:"BikeID"`
	Bookings []booking.Booking `gorm:"foreignkey:BikeID"`
}
type BookDates struct {
	gorm.Model
	StartDate string
	EndDate   string
}

func (b *Bike) Get() error {
	err := sql.DB.Model(b).First(b).Error
	if err != nil {
		return err
	}
	return sql.DB.Model(b).Related(&b.Bookings).Error
}

func (b *Bike) GetAll() ([]Bike, error) {
	var cs []Bike
	err := sql.DB.Model(b).Find(&cs).Error
	return cs, err
}

func (b *Bike) Create() (err error) {
	err = sql.DB.Model(b).Create(b).Error
	return
}

func (b *Bike) Update() (err error) {
	err = sql.DB.Model(b).Updates(b).Error
	return
}
