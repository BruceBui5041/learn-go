package ginrestaurant

import (
	"learn-go/food_delivery_be/common"
	"learn-go/food_delivery_be/component"
	"learn-go/food_delivery_be/modules/restaurant/restaurantbiz"
	"learn-go/food_delivery_be/modules/restaurant/restaurantmodel"
	"learn-go/food_delivery_be/modules/restaurant/restaurantrepo"
	"learn-go/food_delivery_be/modules/restaurant/restaurantstorage"
	restaurantlikestorage "learn-go/food_delivery_be/modules/restaurantlike/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListRestaurant(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filter restaurantmodel.Filter

		if err := c.ShouldBind(&filter); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		var paging common.Paging

		if err := c.ShouldBind(&paging); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		paging.Fulfill()

		restaurantStore := restaurantstorage.NewSQLStore(appCtx.GetMainDBConnection())
		restaurantLikeStore := restaurantlikestorage.NewSQLStore(appCtx.GetMainDBConnection())
		repo := restaurantrepo.NewRestaurantRepo(restaurantStore, restaurantLikeStore)
		biz := restaurantbiz.NewListRestaurantBiz(repo)

		result, err := biz.ListRestaurant(c.Request.Context(), &filter, &paging)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})

			return
		}

		for i := range result {
			result[i].Mask(false)

			if i == len(result)-1 {
				paging.NextCursor = result[i].FakeId.String()
			}
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(result, paging, filter))
	}
}
