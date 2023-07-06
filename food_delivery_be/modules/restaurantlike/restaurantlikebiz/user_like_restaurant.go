package restaurantlikebiz

import (
	"context"
	"learn-go/food_delivery_be/common"
	"learn-go/food_delivery_be/modules/restaurantlike/restaurantlikemodel"
	"learn-go/food_delivery_be/pubsub"
)

type UsersLikeRestaurantStore interface {
	Create(ctx context.Context, data *restaurantlikemodel.Like) error
	Find(ctx context.Context, condition map[string]interface{}) (*restaurantlikemodel.Like, error)
}

// type IncreaseLikeCountStore interface {
// 	IncreaseLikeCount(ctx context.Context, id int) error
// }

type userLikeRestaurantBiz struct {
	store UsersLikeRestaurantStore
	// increaseStore IncreaseLikeCountStore
	pubsub pubsub.Pubsub
}

func NewUserLikeRestaurantBiz(
	store UsersLikeRestaurantStore,
	// increaseStore IncreaseLikeCountStore,
	pubsub pubsub.Pubsub,
) *userLikeRestaurantBiz {
	return &userLikeRestaurantBiz{
		store: store,
		// increaseStore: increaseStore,
		pubsub: pubsub,
	}
}

func (biz *userLikeRestaurantBiz) UserLikeRestaurant(
	ctx context.Context,
	data *restaurantlikemodel.Like,
) error {
	// NOTE: nên find coi coi user đã like restaurant chưa trước khi create
	// vì tìm theo 2 khoá chính user_id và restaurant_id nên sẽ rất nhanh => nên làm
	if like, err := biz.store.Find(ctx, map[string]interface{}{"restaurant_id": data.RestaurantId, "user_id": data.UserId}); like != nil {
		return restaurantlikemodel.ErrCannotLikeRestaurant(err)
	}

	err := biz.store.Create(ctx, data)

	if err != nil {
		return restaurantlikemodel.ErrCannotLikeRestaurant(err)
	}

	biz.pubsub.Publish(ctx, common.TopicUserLikeRestaurant, pubsub.NewMessage(data))

	// Solution 2
	// This is consider as side effect so no need to handle error
	// go func() {
	// 	defer common.AppRecover()
	// 	job := asyncjob.NewJob(func(ctx context.Context) error {
	// 		return biz.increaseStore.IncreaseLikeCount(ctx, data.RestaurantId)
	// 	})

	// 	_ = asyncjob.NewGroup(true, job).Run(ctx)
	// }()

	// Solution 1
	// This is consider as side effect so no need to handle error
	// go func() {
	// 	defer common.AppRecover()
	// 	_ = biz.increaseStore.IncreaseLikeCount(ctx, data.RestaurantId)
	// }()

	return nil
}
