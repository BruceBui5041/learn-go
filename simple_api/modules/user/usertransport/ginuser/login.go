package ginuser

import (
	"learn-go/simple_api/common"
	"learn-go/simple_api/component"
	"learn-go/simple_api/component/hasher"
	"learn-go/simple_api/component/tokenprovider/jwt"
	"learn-go/simple_api/modules/user/userbiz"
	"learn-go/simple_api/modules/user/usermodel"
	"learn-go/simple_api/modules/user/userstorage"
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
