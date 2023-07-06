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

func DeleteRestaurant(appContext component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		restaurantUID, err := common.FromBase58(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "id must be a integer number",
			})
			return
		}

		store := restaurantstorage.NewSQLStore(appContext.GetMainDBConnection())
		biz := restaurantbiz.NewDeleteRestaurantBiz(store)

		if err := biz.SoftDeleteRestaurant(c.Request.Context(), int(restaurantUID.GetLocalID())); err != nil {

			if err.Error() == "record not found" {
				panic(common.ErrEntityNotFound(restaurantmodel.Entity, err))
			}

			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
