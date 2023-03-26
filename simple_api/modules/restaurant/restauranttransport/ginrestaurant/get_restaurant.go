package ginrestaurant

import (
	"learn-go/simple_api/common"
	"learn-go/simple_api/component"
	"learn-go/simple_api/modules/restaurant/restaurantbiz"
	"learn-go/simple_api/modules/restaurant/restaurantstorage"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetRestaurantById(appContext component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		restaurantId, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "id must be a integer number",
			})
			return
		}

		store := restaurantstorage.NewSQLStore(appContext.GetMainDBConnection())
		biz := restaurantbiz.NewGetRestaurantBiz(store)

		result, err := biz.GetRestaurantById(c.Request.Context(), restaurantId)

		if err != nil {
			statusCode := http.StatusUnauthorized

			if err.Error() == "record not found" {
				statusCode = http.StatusNotFound
			}

			c.JSON(statusCode, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))
	}
}
