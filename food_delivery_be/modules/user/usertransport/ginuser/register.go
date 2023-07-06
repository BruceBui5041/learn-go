package ginuser

import (
	"learn-go/food_delivery_be/common"
	"learn-go/food_delivery_be/component"
	"learn-go/food_delivery_be/component/hasher"
	"learn-go/food_delivery_be/modules/user/userbiz"
	"learn-go/food_delivery_be/modules/user/usermodel"
	"learn-go/food_delivery_be/modules/user/userstorage"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(appCtx component.AppContext) func(*gin.Context) {
	return func(ctx *gin.Context) {
		db := appCtx.GetMainDBConnection()
		var data usermodel.CreateUser

		if err := ctx.ShouldBind(&data); err != nil {
			panic(err)
		}

		store := userstorage.NewSQLStore(db)
		md5 := hasher.NewMD5Hash()
		business := userbiz.NewRegisterBusiness(store, md5)

		if err := business.Register(ctx, &data); err != nil {
			panic(err)
		}

		data.Mask(false)

		ctx.JSON(http.StatusCreated, common.SimpleSuccessResponse(data.FakeId.String()))

	}
}
