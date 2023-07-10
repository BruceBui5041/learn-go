package subscriber

import (
	"context"
	"learn-go/food_delivery_be/component"
	"learn-go/food_delivery_be/modules/restaurant/restaurantstorage"
	"learn-go/food_delivery_be/pubsub"
)

func RunDecreaseLikeCountAfterUserUnlikedRestaurant(appCtx component.AppContext) consumerJob {
	return consumerJob{
		Title: "Decrease like count after user unliked restaurant",
		Hld: func(ctx context.Context, message *pubsub.Message) error {
			store := restaurantstorage.NewSQLStore(appCtx.GetMainDBConnection())
			var restaurantId = message.Data().(int)
			return store.DecreaseLikeCount(ctx, restaurantId)
		},
	}
}
