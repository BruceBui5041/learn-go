package main

import (
	"context"
	"errors"
	"fmt"
	"learn-go/food_delivery_be/component"
	"learn-go/food_delivery_be/component/tokenprovider/jwt"
	"learn-go/food_delivery_be/component/uploadprovider"
	"learn-go/food_delivery_be/middleware"
	"learn-go/food_delivery_be/modules/restaurant/restauranttransport/ginrestaurant"
	"learn-go/food_delivery_be/modules/restaurantlike/restaurantliketransport/ginrestaurantlike"
	"learn-go/food_delivery_be/modules/upload/uploadtransport/ginupload"
	"learn-go/food_delivery_be/modules/user/userstorage"
	"learn-go/food_delivery_be/modules/user/usertransport/ginuser"
	"learn-go/food_delivery_be/pubsub/pblocal"
	"learn-go/food_delivery_be/subscriber"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
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
	appContext := component.NewAppContext(db, uploadProvider, jwtSecretToken, pblocal.NewPubSub())

	// subscriber.SetUp(appContext)
	if err := subscriber.NewEngine(appContext).Start(); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	// NOTE: Sử dung middleware
	r.Use(middleware.Recover(appContext))

	// server demo.html static file to test socketio
	r.StaticFile("/demo/", "./demo/demo.html")

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
		restaurants.POST("/:id/like", ginrestaurantlike.UserLikeRestaurant(appContext))
		restaurants.DELETE("/:id/unlike", ginrestaurantlike.UserUnlikeRestaurant(appContext))
	}

	startSocketIOServer(r, appContext)

	return r.Run()
}

func startSocketIOServer(engine *gin.Engine, appCtx component.AppContext) {
	server, _ := socketio.NewServer(&engineio.Options{
		// ép kiểu về websocket luôn vi có thể nó sẽ tạo ra transport theo Long-polling nếu client k support
		Transports: []transport.Transport{websocket.Default},
	})

	server.OnConnect("/", func(c socketio.Conn) error {
		fmt.Println("Connected: ", c.ID(), " Ip:", c.RemoteAddr())

		return nil
	})

	server.OnError("/", func(c socketio.Conn, err error) {
		fmt.Println("Connect error: ", err)
	})

	server.OnDisconnect("/", func(c socketio.Conn, reason string) {
		fmt.Println("Disconnected: ", reason)
	})

	server.OnEvent("/", "authenticate", func(c socketio.Conn, token string) {
		db := appCtx.GetMainDBConnection()
		store := userstorage.NewSQLStore(db)
		tokenProvider := jwt.NewTokenJWTProvider(appCtx.SecretKey())

		user_payload, err := tokenProvider.Validate(token)

		if err != nil {
			c.Emit("authentication failed", err.Error())
			c.Close()
			return
		}

		user, err := store.FindUser(context.Background(), map[string]interface{}{"id": user_payload.UserId})
		if err != nil {
			c.Emit("authentication failed", err.Error())
			c.Close()
			return
		}

		if user.Status == 0 {
			c.Emit("authentication failed", errors.New("you has been banned or deleted"))
			c.Close()
			return
		}

		user.Mask(false)

		c.Emit("authenticated", user)
	})

	server.OnEvent("/", "test", func(c socketio.Conn, message string) {
		log.Println(message)
	})

	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	server.OnEvent("/", "notice", func(c socketio.Conn, p Person) {
		fmt.Println("Server recieved notice: ", p.Name, p.Age)

		p.Age = 33

		c.Emit("notice", p)
	})

	go server.Serve()

	// bản chất của socket.io nó phải từ http đi lên
	// nên phải có 1 req đi lên để xin quyền upgrade
	engine.GET("/socket.io/*any", gin.WrapH(server))
	engine.POST("/socket.io/*any", gin.WrapH(server))
}
