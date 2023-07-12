package appsocketio

import (
	"context"
	"errors"
	"fmt"
	"learn-go/food_delivery_be/component"
	"learn-go/food_delivery_be/component/tokenprovider/jwt"
	"learn-go/food_delivery_be/modules/user/userstorage"
	"sync"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
)

type RealtimeEngine interface {
	UserSockets(userId int) []AppSocket // Note: this is because 1 user can use more than 1 device
	EmmitToRoom(room string, key string, data interface{}) error
	EmmitToUser(userId int, key string, data interface{}) error
	Run(ctx component.AppContext, ginEngine *gin.Engine) error
}

type rtEngine struct {
	server  *socketio.Server
	storage map[int][]AppSocket
	locker  *sync.RWMutex
}

func NewEngine(
	server *socketio.Server,
	storage map[int][]AppSocket,
	locker *sync.RWMutex,
) *rtEngine {
	return &rtEngine{server: server, storage: storage, locker: locker}
}

func (engine *rtEngine) saveAppSocket(userId int, appSpcket AppSocket) {
	engine.locker.Lock()
	defer engine.locker.Unlock()

	if sockets, ok := engine.storage[userId]; ok {
		engine.storage[userId] = append(sockets, appSpcket)
	} else {
		engine.storage[userId] = []AppSocket{appSpcket}
	}
}

func (engine *rtEngine) getAppSocket(userId int) []AppSocket {
	engine.locker.Lock()

	defer engine.locker.Unlock()
	return engine.storage[userId]
}

func (engine *rtEngine) removeAppSocket(userId int, appSpcket AppSocket) {
	engine.locker.Lock()

	defer engine.locker.Unlock()
	if v, ok := engine.storage[userId]; ok {
		for i := range v {
			if v[i] == appSpcket {
				engine.storage[userId] = append(v[:i], v[i+1:]...)
				break
			}
		}
	}
}

func (engine *rtEngine) UserSockets(userId int) []AppSocket {
	var sockets []AppSocket

	if skts, ok := engine.storage[userId]; ok {
		return skts
	}

	return sockets
}

func (engine *rtEngine) EmmitToRoom(room string, key string, data interface{}) error {
	engine.server.BroadcastToRoom("/", room, key, data)
	return nil
}

func (engine *rtEngine) EmmitToUser(userId int, key string, data interface{}) error {
	sockets := engine.storage[userId]

	for _, s := range sockets {
		s.Emit(key, data)
	}

	return nil
}

func (engine *rtEngine) Run(appCtx component.AppContext, ginEngine *gin.Engine) error {
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

		appSkt := NewAppSocket(c, user)
		engine.saveAppSocket(user.Id, appSkt)

		user.Mask(false)

		c.Emit("authenticated", user)
	})

	go server.Serve()

	// bản chất của socket.io nó phải từ http đi lên
	// nên phải có 1 req đi lên để xin quyền upgrade
	ginEngine.GET("/socket.io/*any", gin.WrapH(server))
	ginEngine.POST("/socket.io/*any", gin.WrapH(server))

	return nil
}
