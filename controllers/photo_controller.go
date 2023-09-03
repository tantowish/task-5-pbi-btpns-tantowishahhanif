package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/tantowish/task-5-vix-btpns-tantowishahhanif/app"
	"github.com/tantowish/task-5-vix-btpns-tantowishahhanif/auth"
	"github.com/tantowish/task-5-vix-btpns-tantowishahhanif/formaterror"
	"github.com/tantowish/task-5-vix-btpns-tantowishahhanif/models"
)

func GetPhoto(c *gin.Context) {
	photos := []models.Photo{}

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Debug().Model(&models.Photo{}).Limit(100).Find(&photos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "F", "message": "photo not found", "data": nil})
		return
	}

	if len(photos) > 0 {
		for i := range photos {
			user := models.User{}
			err := db.Model(&models.User{}).Where("id = ?", photos[i].UserID).Take(&user).Error

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"status": "F", "message": err.Error(), "data": nil})
				return
			}

			photos[i].Author = app.Author{
				ID: user.ID, Username: user.Username, Email: user.Email,
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "T", "message": "success", "data": photos})
}

func CreatePhoto(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(401, gin.H{"error": "request does not contain an access token"})
		return
	}

	email, err := auth.GetEmail(strings.Split(tokenString, "Bearer ")[1])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": err.Error(), "data": nil})
	}

	var user_login models.User

	err = db.Debug().Where("email = ?", email).First(&user_login).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": "user not found", "data": nil})
		return
	}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
	}

	photo_input := models.Photo{}
	err = json.Unmarshal(body, &photo_input)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	photo_input.Initialize()
	photo_input.UserID = user_login.ID
	photo_input.Author = app.Author{
		ID:       user_login.ID,
		Username: user_login.Username,
		Email:    user_login.Email,
	}

	err = photo_input.Validate("upload")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	var old_photo models.Photo
	err = db.Debug().Model(&models.Photo{}).Where("user_id = ?", user_login.ID).Find(&old_photo).Error
	if err != nil {
		if err.Error() == "record not found" {
			err = db.Debug().Create(&photo_input).Error
			if err != nil {
				formattedError := formaterror.ErrorMessage(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"status": "F", "message": formattedError.Error(), "data": nil})
				return
			}
			c.JSON(http.StatusOK, gin.H{"status": "T", "message": "success upload photo", "data": photo_input})
			return
		}
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	photo_input.ID = old_photo.ID
	err = db.Debug().Model(&old_photo).Updates(&photo_input).Error
	if err != nil {
		formattedError := formaterror.ErrorMessage(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "F", "message": formattedError.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "T", "message": "success change photo", "data": photo_input})
}

func UpdatePhoto(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(401, gin.H{"error": "request does not contain an access token"})
		return
	}

	email, err := auth.GetEmail(strings.Split(tokenString, "Bearer ")[1])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": err.Error(), "data": nil})
	}

	var user_login models.User

	err = db.Debug().Where("email = ?", email).First(&user_login).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": "user not found", "data": nil})
		return
	}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
	}

	photo_input := models.Photo{}
	err = json.Unmarshal(body, &photo_input)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	err = photo_input.Validate("change")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	var photo models.Photo
	if err := db.Debug().Where("id = ?", c.Param("photoId")).First(&photo).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": "Photo not found", "data": nil})
		return
	}

	if user_login.ID != photo.UserID {
		c.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": "no access to change photo", "data": nil})
		return
	}

	err = db.Model(&photo).Updates(&photo_input).Error
	if err != nil {
		formattedError := formaterror.ErrorMessage(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "F", "message": formattedError.Error(), "data": nil})
		return
	}

	photo.Author = app.Author{
		ID:       user_login.ID,
		Username: user_login.Username,
		Email:    user_login.Email,
	}

	c.JSON(http.StatusOK, gin.H{"status": "T", "message": "success change photo", "data": photo})
}

func DeletePhoto(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(401, gin.H{"error": "request does not contain an access token"})
		return
	}

	email, err := auth.GetEmail(strings.Split(tokenString, "Bearer ")[1])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": err.Error(), "data": nil})
	}

	var user_login models.User
	if err := db.Debug().Where("email = ?", email).First(&user_login).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": "user not found", "data": nil})
		return
	}

	var photo models.Photo
	if err := db.Debug().Where("id = ?", c.Param("photoId")).First(&photo).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": "Photo not found", "data": nil})
		return
	}

	if user_login.ID != photo.UserID {
		c.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": "no access to delete photo", "data": nil})
		return
	}

	err = db.Debug().Delete(&photo).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "T", "message": "delete photo success", "data": nil})
}
