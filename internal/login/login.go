package login

import (
	"bikerentalProject/internal/user"
	"bikerentalProject/pkg/logger"
	"bikerentalProject/pkg/sql"
	"encoding/base64"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"math/rand"
	"net/http"
)

type data struct {
	gorm.Model
	Name     string
	Email    string
	Password string
}

func UserLogin(g *gin.Context) {
	g.HTML(http.StatusOK, "login.tmpl", gin.H{})
}

// RandToken generates a random @l length token.
func RandToken(l int) (string, error) {
	b := make([]byte, l)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}
func AuthHandler(g *gin.Context) {
	u := &data{Name: g.PostForm("Name"),
		Email:    g.PostForm("Email"),
		Password: g.PostForm("Password")}
	var rows []user.User
	err := sql.DB.Table("users").Where("Email=?", u.Email).Find(&rows).Error
	if err != nil {
		fmt.Println(err)
		logger.ErrorLogger.Printf("error in getting booked bike get method, err:%v", err)
		g.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"message": "Internal Error"})
		return
	}
	if len(rows) != 0 {
		if u.Password != rows[0].Password {
			fmt.Println("Password does not match")
			logger.ErrorLogger.Printf("Password does not match, err:%v", err)
			g.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"message": "Password does not match"})
			return
		}
	} else {
		B := &user.User{}
		B.Name = u.Name
		B.Email = u.Email
		B.Password = u.Password
		err := B.Create()
		if err != nil {
			logger.ErrorLogger.Printf("error in creating user, err:%v", err)
			g.JSON(http.StatusInternalServerError, gin.H{
				"data":     nil,
				"debugMsg": "Error in creating user",
			})
			return
		}
	}
	session := sessions.Default(g)
	session.Set("user-id", u.Email) //take it from html page
	err = session.Save()
	if err != nil {
		logger.ErrorLogger.Printf("error in saving from session, error:%v", err)
		g.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Error while saving session. Please try again."})
		return
	}
	g.Redirect(http.StatusFound, "/")
}

func AuthorizeRequest() gin.HandlerFunc {
	return func(g *gin.Context) {
		session := sessions.Default(g)
		v := session.Get("user-id")
		fmt.Println(v)
		if v == nil {
			logger.ErrorLogger.Printf("user is not valid, email:%s", v)
			g.HTML(http.StatusUnauthorized, "error.tmpl", gin.H{"message": "Please login."})
			g.Abort()
		}
		g.Next()
	}
}
func Isloggedin(g *gin.Context) string {
	session := sessions.Default(g)
	v := session.Get("user-id")
	k := ""
	if v == nil {
		logger.ErrorLogger.Printf("user is not valid, email:%s", v)
		g.Abort()
	} else {
		k = v.(string)
	}
	return k
}
func AdminAuthRequest() gin.HandlerFunc {
	return func(g *gin.Context) {
		session := sessions.Default(g)
		v := session.Get("user-id")
		if v == nil || (v != "saurabhk442@gmail.com" && v != "rajsumedha11@gmail.com") {
			logger.ErrorLogger.Printf("user is not admin, email:%s", v)
			g.HTML(http.StatusUnauthorized, "error.tmpl", gin.H{"message": "Admin login required."})
			g.Abort()
		}
		g.Next()
	}
}

func Logout(g *gin.Context) {
	session := sessions.Default(g)
	session.Clear()
	session.Save()
	g.Redirect(http.StatusFound, "/")
}
