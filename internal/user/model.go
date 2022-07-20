package user

import (
	"bikerentalProject/internal/booking"
	"bikerentalProject/pkg/sql"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique;not null"`
	Active   bool
	Password string            `json:"password" gorm:"not null"`
	Bookings []booking.Booking `gorm:"foreignkey:UserID"`
}

func (u *User) Get() (err error) {
	return sql.DB.Model(u).First(u).Error
}

func (u *User) GetByEmail() (err error) {
	return sql.DB.Model(u).Where("email = ?", u.Email).First(&u).Error
}

func (u *User) Create() (err error) {
	err = sql.DB.Model(u).Create(u).Error
	return
}
