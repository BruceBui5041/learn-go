package subscriber

import (
	"context"
	"learn-go/simple_api/common"
	"learn-go/simple_api/component"
	"learn-go/simple_api/modules/restaurant/restaurantstorage"
	"learn-go/simple_api/modules/restaurantlike/restaurantlikemodel"
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
