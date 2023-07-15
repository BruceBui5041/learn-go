package subscriber

import (
	"context"
	"learn-go/food_delivery_be/appsocketio"
	"learn-go/food_delivery_be/component"
	"learn-go/food_delivery_be/modules/restaurant/restaurantstorage"
	"learn-go/food_delivery_be/pubsub"
)

type HasRestaurantId interface {
	GetRestaurantId() int
	GetUserId() int
}

// func IncreaseLikeCountAfterUserLikedRestaurant(appCtx component.AppContext, ctx context.Context) {
// 	c, _ := appCtx.GetPubSub().Subscribe(ctx, common.TopicUserLikeRestaurant)
// 	store := restaurantstorage.NewSQLStore(appCtx.GetMainDBConnection())
// 	go func() {
// 		defer common.AppRecover()
// 		for {
// 			msg := <-c
// 			var like = msg.Data().(HasRestaurantId)
// 			store.IncreaseLikeCount(ctx, like.GetRestaurantId())
// 		}
// 	}()
// }

// func RunIncreaseLikeCountAfterUserLikedRestaurant(appCtx component.AppContext) func(ctx context.Context, message *pubsub.Message,) error {
// 	store := restaurantstorage.NewSQLStore(appCtx.GetMainDBConnection())

// 	return func(ctx context.Context, message *pubsub.Message) error {
// 		var like = message.Data().(HasRestaurantId)
// 		return store.IncreaseLikeCount(ctx, like.GetRestaurantId())
// 	}
// }

func RunIncreaseLikeCountAfterUserLikedRestaurant(appCtx component.AppContext) consumerJob {
	return consumerJob{
		Title: "Increase like count after user liked restaurant",
		Hld: func(ctx context.Context, message *pubsub.Message) error {
			store := restaurantstorage.NewSQLStore(appCtx.GetMainDBConnection())
			var like = message.Data().(HasRestaurantId)

			return store.IncreaseLikeCount(ctx, like.GetRestaurantId())
		},
	}
}

func EmitRealtimeAfterUserLikedRestaurant(appCtx component.AppContext, rtEngine appsocketio.RealtimeEngine) consumerJob {
	return consumerJob{
		Title: "Emit realtime after user liked restaurant",
		Hld: func(ctx context.Context, message *pubsub.Message) error {
			var like = message.Data().(HasRestaurantId)
			userId := like.GetUserId()
			return rtEngine.EmmitToUser(userId, string(message.Channel()), like)
		},
	}
}
