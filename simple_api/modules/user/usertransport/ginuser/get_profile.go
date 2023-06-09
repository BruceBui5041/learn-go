package ginuser

import (
	"learn-go/simple_api/common"
	"learn-go/simple_api/component"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProfile(appCtx component.AppContext) func(*gin.Context) {
	return func(ctx *gin.Context) {
		data := ctx.MustGet(common.CurrentUser).(common.Requester)

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(data))
	}
}
