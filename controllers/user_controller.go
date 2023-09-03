package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/tantowish/task-5-pbi-btpns-tantowishahhanif/app"
	"github.com/tantowish/task-5-pbi-btpns-tantowishahhanif/app/auth"
	"github.com/tantowish/task-5-pbi-btpns-tantowishahhanif/helpers/formaterror"
	"github.com/tantowish/task-5-pbi-btpns-tantowishahhanif/helpers/hash"
	"github.com/tantowish/task-5-pbi-btpns-tantowishahhanif/models"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
	}

	user_input := models.User{}
	err = json.Unmarshal(body, &user_input)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	user_input.Initialize()

	err = user_input.Validate("update")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	err = user_input.HashPassword()
	if err != nil {
		log.Fatal(err)
	}

	err = db.Debug().Create(&user_input).Error
	if err != nil {
		formattedError := formaterror.ErrorMessage(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "F", "message": formattedError.Error(), "data": nil})
		return
	}

	data := app.UserRegister{
		ID:        user_input.ID,
		Username:  user_input.Username,
		Email:     user_input.Email,
		CreatedAt: user_input.CreatedAt,
		UpdatedAt: user_input.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{"status": "T", "message": "register user success", "data": data})
}

func UpdateUser(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	var user models.User
	err := db.Debug().Where("id = ?", c.Param("userId")).First(&user).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": "user not found", "data": nil})
		return
	}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
	}

	user_input := models.User{}
	user_input.ID = user.ID
	err = json.Unmarshal(body, &user_input)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	err = user_input.Validate("update")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	err = user_input.HashPassword()
	if err != nil {
		log.Fatal(err)
	}

	err = db.Debug().Model(&user).Updates(&user_input).Error
	if err != nil {
		formattedError := formaterror.ErrorMessage(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "F", "message": formattedError.Error(), "data": nil})
		return
	}

	data := app.UserRegister{
		ID:        user_input.ID,
		Username:  user_input.Username,
		Email:     user_input.Email,
		CreatedAt: user_input.CreatedAt,
		UpdatedAt: user_input.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{"status": "T", "message": "update user success", "data": data})
}

func DeleteUser(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	var user models.User

	err := db.Debug().Where("id = ?", c.Param("userId")).First(&user).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": "user not found", "data": nil})
		return
	}

	err = db.Debug().Delete(&user).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "T", "message": "delete user success", "data": nil})
}

func Login(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	user_input := models.User{}
	err = json.Unmarshal(body, &user_input)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	user_input.Initialize()
	err = user_input.Validate("login")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	var user_login app.UserLogin

	err = db.Debug().Table("users").Select("*").Joins("left join photos on photos.user_id = users.id").
		Where("users.email = ?", user_input.Email).Find(&user_login).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": "user not found", "data": nil})
		return
	}

	err = hash.VerifyPassword(user_login.Password, user_input.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		formattedError := formaterror.ErrorMessage(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": formattedError.Error(), "data": nil})
		return
	}

	token, err := auth.GenerateJWT(user_login.Email, user_login.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	data := app.DataUser{
		ID: user_login.ID, Username: user_login.Username, Email: user_login.Email, Token: token,
		Photos: app.Photo{Title: user_login.Title, Caption: user_login.Caption, PhotoUrl: user_login.PhotoUrl},
	}

	c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "T", "message": "login success", "data": data})
}
