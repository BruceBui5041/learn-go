package ginrestaurantlike

import (
	"learn-go/simple_api/common"
	"learn-go/simple_api/component"
	"learn-go/simple_api/modules/restaurant/restaurantstorage"
	"learn-go/simple_api/modules/restaurantlike/restaurantlikebiz"
	"learn-go/simple_api/modules/restaurantlike/restaurantlikemodel"
	restaurantlikestorage "learn-go/simple_api/modules/restaurantlike/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

// POST /v1/restaurants/:id/like

func UserLikeRestaurant(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, err := common.FromBase58(c.Param("id"))

		if err != nil {
			panic(err)
		}

		requester := c.MustGet(common.CurrentUser).(common.Requester)

		restaurantLike := restaurantlikemodel.Like{
			RestaurantId: int(uid.GetLocalID()),
			UserId:       requester.GetUserId(),
		}

		restaurantLikeStore := restaurantlikestorage.NewSQLStore(appCtx.GetMainDBConnection())
		restaurantStore := restaurantstorage.NewSQLStore(appCtx.GetMainDBConnection())
		biz := restaurantlikebiz.NewUserLikeRestaurantBiz(restaurantLikeStore, restaurantStore)

		err = biz.UserLikeRestaurant(c.Request.Context(), &restaurantLike)

		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
