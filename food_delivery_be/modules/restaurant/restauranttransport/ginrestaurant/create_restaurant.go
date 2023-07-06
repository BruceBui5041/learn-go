package ginrestaurant

import (
	"learn-go/food_delivery_be/common"
	"learn-go/food_delivery_be/component"
	"learn-go/food_delivery_be/modules/restaurant/restaurantbiz"
	"learn-go/food_delivery_be/modules/restaurant/restaurantmodel"
	"learn-go/food_delivery_be/modules/restaurant/restaurantstorage"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateRestaurant(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data restaurantmodel.RestaurantCreate

		if err := c.ShouldBind(&data); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		requester := c.MustGet(common.CurrentUser).(common.Requester)

		data.UserId = requester.GetUserId()

		store := restaurantstorage.NewSQLStore(appCtx.GetMainDBConnection())
		biz := restaurantbiz.NewCreateRestaurantBiz(store)

		if err := biz.CreateRestaurant(c.Request.Context(), &data); err != nil {
			panic(common.ErrDB(err))
		}

		data.Mask(false)

		c.JSON(http.StatusCreated, common.SimpleSuccessResponse(data.FakeId.String()))
	}
}
