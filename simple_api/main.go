package main

import (
	"learn-go/simple_api/component"
	"learn-go/simple_api/component/uploadprovider"
	"learn-go/simple_api/middleware"
	"learn-go/simple_api/modules/restaurant/restauranttransport/ginrestaurant"
	"learn-go/simple_api/modules/restaurantlike/restaurantliketransport/ginrestaurantlike"
	"learn-go/simple_api/modules/upload/uploadtransport/ginupload"
	"learn-go/simple_api/modules/user/usertransport/ginuser"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	jwtSecretKey := viper.GetString("JWTSecretKey")

	s3Provider := uploadprovider.NewS3Provider(s3BucketName, s3Region, s3APIKey, s3SecretKey, s3Domain)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,       // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Fatalln(err)
	}

	if err := runService(db, s3Provider, jwtSecretKey); err != nil {
		log.Fatalln(err)
	}
}

func runService(db *gorm.DB, uploadProvider uploadprovider.UploadProvider, jwtSecretToken string) error {
	appContext := component.NewAppContext(db, uploadProvider, jwtSecretToken)

	r := gin.Default()

	// NOTE: Sử dung middleware
	r.Use(middleware.Recover(appContext))

	r.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	v1 := r.Group("/v1")

	v1.POST("/upload", ginupload.Upload(appContext))
	v1.POST("/register", ginuser.Register(appContext))
	v1.POST("/login", ginuser.Login(appContext))
	v1.GET("/profile", middleware.RequiredAuth(appContext), ginuser.GetProfile(appContext))

	restaurants := v1.Group("/restaurants", middleware.RequiredAuth(appContext))
	{
		restaurants.POST("", ginrestaurant.CreateRestaurant(appContext))
		restaurants.GET("/:id", ginrestaurant.GetRestaurantById(appContext))
		restaurants.GET("", ginrestaurant.ListRestaurant(appContext))
		restaurants.PATCH("/:id", ginrestaurant.UpdateRestaurantById(appContext))
		restaurants.DELETE("/:id", ginrestaurant.DeleteRestaurant(appContext))

		restaurants.GET("/:id/liked-users", ginrestaurantlike.ListUserLikeRestaurant(appContext))
	}

	return r.Run()
}
