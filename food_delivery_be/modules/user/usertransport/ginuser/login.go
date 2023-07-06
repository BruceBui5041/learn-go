package ginuser

import (
	"learn-go/food_delivery_be/common"
	"learn-go/food_delivery_be/component"
	"learn-go/food_delivery_be/component/hasher"
	"learn-go/food_delivery_be/component/tokenprovider/jwt"
	"learn-go/food_delivery_be/modules/user/userbiz"
	"learn-go/food_delivery_be/modules/user/usermodel"
	"learn-go/food_delivery_be/modules/user/userstorage"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(appCtx component.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var loginUserData usermodel.UserLogin

		if err := ctx.ShouldBind(&loginUserData); err != nil {
			panic(err)
		}

		db := appCtx.GetMainDBConnection()
		tokenProvider := jwt.NewTokenJWTProvider(appCtx.SecretKey())

		md5 := hasher.NewMD5Hash()

		userStore := userstorage.NewSQLStore(db)
		loginbiz := userbiz.NewLoginBusiness(userStore, tokenProvider, md5, 60*60*24*30)

		account, err := loginbiz.Login(ctx, &loginUserData)

		if err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(account))
	}
}
