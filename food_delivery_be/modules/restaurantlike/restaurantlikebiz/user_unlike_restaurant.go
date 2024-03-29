package restaurantlikebiz

import (
	"context"
	"learn-go/food_delivery_be/common"
	"learn-go/food_delivery_be/modules/restaurantlike/restaurantlikemodel"
	"learn-go/food_delivery_be/pubsub"
)

type UserUnlikeRestaurantStore interface {
	Delete(ctx context.Context, userId, restaurantId int) error
	Find(ctx context.Context, condition map[string]interface{}) (*restaurantlikemodel.Like, error)
}

type DecreaseLikeCountStore interface {
	DecreaseLikeCount(ctx context.Context, id int) error
}

type userUnlikeRestaurantBiz struct {
	store UserUnlikeRestaurantStore
	// decreaseLikeStore DecreaseLikeCountStore
	pubsub pubsub.Pubsub
}

func NewUserUnlikeRestaurantBiz(
	store UserUnlikeRestaurantStore,
	// decreaseLikeStore DecreaseLikeCountStore,
	pubsub pubsub.Pubsub,
) *userUnlikeRestaurantBiz {
	return &userUnlikeRestaurantBiz{
		store: store,
		// decreaseLikeStore: decreaseLikeStore,
		pubsub: pubsub,
	}
}

func (biz *userUnlikeRestaurantBiz) UserUnlikeRestaurant(
	ctx context.Context,
	userId, restaurantId int,
) error {
	// NOTE: nên find coi coi user đã like restaurant chưa trước khi create
	// vì tìm theo 2 khoá chính user_id và restaurant_id nên sẽ rất nhanh => nên làm
	if like, err := biz.store.Find(ctx, map[string]interface{}{"restaurant_id": restaurantId, "user_id": userId}); like == nil {
		return restaurantlikemodel.ErrCannotUnlikeRestaurant(err)
	}

	err := biz.store.Delete(ctx, userId, restaurantId)

	if err != nil {
		return restaurantlikemodel.ErrCannotUnlikeRestaurant(err)
	}

	biz.pubsub.Publish(ctx, common.TopicUserDislikeRestaurant, pubsub.NewMessage(restaurantId))

	// This is consider as side effect so no need to handle error
	// go func() {
	// 	defer common.AppRecover()
	// 	job := asyncjob.NewJob(func(ctx context.Context) error {
	// 		return biz.decreaseLikeStore.DecreaseLikeCount(ctx, restaurantId)
	// 	})

	// 	_ = asyncjob.NewGroup(true, job).Run(ctx)
	// }()

	// This is consider as side effect so no need to handle error
	// go func() {
	// 	defer common.AppRecover()
	// 	_ = biz.decreaseLikeStore.DecreaseLikeCount(ctx, restaurantId)
	// }()

	return nil
}
