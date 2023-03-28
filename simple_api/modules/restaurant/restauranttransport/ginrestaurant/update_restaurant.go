package ginrestaurant

import (
	"learn-go/simple_api/common"
	"learn-go/simple_api/component"
	"learn-go/simple_api/modules/restaurant/restaurantbiz"
	"learn-go/simple_api/modules/restaurant/restaurantmodel"
	"learn-go/simple_api/modules/restaurant/restaurantstorage"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UpdateRestaurantById(appContext component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		restaurantId, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "id must be a integer number",
			})
			return
		}

		var data restaurantmodel.RestaurantUpdate

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
		}

		store := restaurantstorage.NewSQLStore(appContext.GetMainDBConnection())
		biz := restaurantbiz.NewUpdateRestaurantBiz(store)

		if err := biz.UpdateRestaurant(c.Request.Context(), restaurantId, &data); err != nil {

			statusCode := http.StatusUnauthorized

			if err.Error() == "record not found" {
				statusCode = http.StatusNotFound
			}

			c.JSON(statusCode, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
