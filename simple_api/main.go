package main

import (
	"learn-go/simple_api/component"
	"learn-go/simple_api/middleware"
	"learn-go/simple_api/modules/restaurant/restauranttransport/ginrestaurant"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	dsn := viper.GetString("DBConnectionStr")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	if err := runService(db); err != nil {
		log.Fatalln(err)
	}
}

func runService(db *gorm.DB) error {
	appContext := component.NewAppContext(db)

	r := gin.Default()

	// NOTE: Sử dung middleware
	r.Use(middleware.Recover(appContext))

	r.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	restaurants := r.Group("/restaurants")
	{
		restaurants.POST("", ginrestaurant.CreateRestaurant(appContext))
		restaurants.GET("/:id", ginrestaurant.GetRestaurantById(appContext))
		restaurants.GET("", ginrestaurant.ListRestaurant(appContext))
		restaurants.PATCH("/:id", ginrestaurant.UpdateRestaurantById(appContext))
		restaurants.DELETE("/:id", ginrestaurant.DeleteRestaurant(appContext))
	}

	return r.Run()
}
