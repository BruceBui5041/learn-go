package subscriber

import (
	"context"
	"learn-go/food_delivery_be/common"
	"learn-go/food_delivery_be/component"
	"learn-go/food_delivery_be/modules/restaurant/restaurantstorage"
	"learn-go/food_delivery_be/modules/restaurantlike/restaurantlikemodel"
)

func IncreaseLikeCountAfterUserLikedRestaurant(appCtx component.AppContext, ctx context.Context) {
	c, _ := appCtx.GetPubSub().Subscribe(ctx, common.TopicUserLikeRestaurant)

	store := restaurantstorage.NewSQLStore(appCtx.GetMainDBConnection())

	go func() {
		defer common.AppRecover()
		for {
			msg := <-c
			var like = msg.Data().(*restaurantlikemodel.Like)
			store.IncreaseLikeCount(ctx, like.RestaurantId)
		}
	}()
}
