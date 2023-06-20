package ginrestaurantlike

import (
	"learn-go/simple_api/common"
	"learn-go/simple_api/component"
	"learn-go/simple_api/modules/restaurantlike/restaurantlikebiz"
	"learn-go/simple_api/modules/restaurantlike/restaurantlikemodel"
	restaurantlikestorage "learn-go/simple_api/modules/restaurantlike/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /v1/restaurants/:id/liked-users

func ListUserLikeRestaurant(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, err := common.FromBase58(c.Param("id"))

		if err != nil {
			panic(err)
		}

		// var filter restaurantlikemodel.Filter

		// if err := c.ShouldBind(&filter); err != nil {
		// 	c.JSON(http.StatusUnauthorized, gin.H{
		// 		"error": err.Error(),
		// 	})
		// 	return
		// }

		filter := restaurantlikemodel.Filter{
			RestaurantId: int(uid.GetLocalID()),
		}

		var paging common.Paging

		if err := c.ShouldBind(&paging); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		paging.Fulfill()

		restaurantLikeStore := restaurantlikestorage.NewSQLStore(appCtx.GetMainDBConnection())
		biz := restaurantlikebiz.NewListUserLikeRestaurantBiz(restaurantLikeStore)

		result, err := biz.ListUserLikeRestaurant(c.Request.Context(), &filter, &paging)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})

			return
		}

		for i := range result {
			result[i].Mask(false)

			// NOTE: Ở đây không thể định nghĩa NextCursor chỉ trong storage mới có thể định nghĩa vì đang order by create_at field
			// if i == len(result)-1 {
			// 	paging.NextCursor = result[i].FakeId.String()
			// }
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(result, paging, filter))
	}
}
