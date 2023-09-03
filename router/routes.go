package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/tantowish/task-5-vix-btpns-tantowishahhanif/controllers"
	"github.com/tantowish/task-5-vix-btpns-tantowishahhanif/middlewares"
)

func SetupRoutes(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
	})

	r.POST("/users/login", controllers.Login)
	r.POST("/users/register", controllers.CreateUser)
	r.PUT("/users/:userId", controllers.UpdateUser)
	r.DELETE("/users/:userId", controllers.DeleteUser)

	r.GET("/photos", controllers.GetPhoto)

	secured := r.Group("/").Use(middlewares.Auth())
	{
		secured.POST("/photos", controllers.CreatePhoto)
		secured.PUT("/photos/:photoId", controllers.UpdatePhoto)
		secured.DELETE("/photos/:photoId", controllers.DeletePhoto)
	}

	return r
}
