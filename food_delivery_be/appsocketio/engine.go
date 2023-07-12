package appsocketio

import (
	"context"
	"sync"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

type RealtimeEngine interface {
	UserSockets(userId int) []AppSocket
	EmmitToRoom(room string, key string, data interface{}) error
	EmmitToUser(userId int, key string, data interface{}) error
	Run(ctx context.Context, engine *gin.Engine)
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
