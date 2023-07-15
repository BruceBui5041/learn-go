package wsuser

import (
	"learn-go/food_delivery_be/common"
	"learn-go/food_delivery_be/component"
	"log"

	socketio "github.com/googollee/go-socket.io"
)

type LocationData struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lg"`
}

func OnUserUpdateLocation(appCtx component.AppContext, requester common.Requester) func(c socketio.Conn, location LocationData) {
	return func(c socketio.Conn, location LocationData) {
		log.Print("User ", requester.GetUserId(), "upload location", location)
	}
}
