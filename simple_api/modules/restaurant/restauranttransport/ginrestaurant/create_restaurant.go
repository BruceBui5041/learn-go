package ginrestaurant

import (
	"learn-go/simple_api/common"
	"learn-go/simple_api/component"
	"learn-go/simple_api/modules/restaurant/restaurantbiz"
	"learn-go/simple_api/modules/restaurant/restaurantmodel"
	"learn-go/simple_api/modules/restaurant/restaurantstorage"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateRestaurant(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data restaurantmodel.RestaurantCreate

		if err := c.ShouldBind(&data); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		store := restaurantstorage.NewSQLStore(appCtx.GetMainDBConnection())
		biz := restaurantbiz.NewCreateRestaurantBiz(store)

		if err := biz.CreateRestaurant(c.Request.Context(), &data); err != nil {
			panic(common.ErrDB(err))
		}

		c.JSON(http.StatusCreated, common.SimpleSuccessResponse(data))
	}
}
