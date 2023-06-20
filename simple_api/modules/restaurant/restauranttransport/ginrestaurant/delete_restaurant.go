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
