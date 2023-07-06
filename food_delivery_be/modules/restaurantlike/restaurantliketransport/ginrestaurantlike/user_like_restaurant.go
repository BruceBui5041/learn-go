package ginrestaurantlike

import (
	"learn-go/food_delivery_be/common"
	"learn-go/food_delivery_be/component"
	"learn-go/food_delivery_be/modules/restaurantlike/restaurantlikebiz"
	"learn-go/food_delivery_be/modules/restaurantlike/restaurantlikemodel"
	restaurantlikestorage "learn-go/food_delivery_be/modules/restaurantlike/storage"
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
		biz := restaurantlikebiz.NewUserLikeRestaurantBiz(restaurantLikeStore, appCtx.GetPubSub())

		err = biz.UserLikeRestaurant(c.Request.Context(), &restaurantLike)

		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
