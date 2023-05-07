package main

import (
	"learn-go/simple_api/component"
	"learn-go/simple_api/component/uploadprovider"
	"learn-go/simple_api/middleware"
	"learn-go/simple_api/modules/restaurant/restauranttransport/ginrestaurant"
	"learn-go/simple_api/modules/upload/uploadtransport/ginupload"
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

	s3BucketName := viper.GetString("S3BucketName")
	s3Region := viper.GetString("S3Region")
	s3APIKey := viper.GetString("S3APIKey")
	s3SecretKey := viper.GetString("S3SecretKey")
	s3Domain := viper.GetString("S3Domain")

	s3Provider := uploadprovider.NewS3Provider(s3BucketName, s3Region, s3APIKey, s3SecretKey, s3Domain)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	if err := runService(db, s3Provider); err != nil {
		log.Fatalln(err)
	}
}

func runService(db *gorm.DB, uploadProvider uploadprovider.UploadProvider) error {
	appContext := component.NewAppContext(db, uploadProvider)

	r := gin.Default()

	// NOTE: Sử dung middleware
	r.Use(middleware.Recover(appContext))

	r.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/upload", ginupload.Upload(appContext))

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
