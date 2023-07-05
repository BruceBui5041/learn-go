package restaurantlikebiz

import (
	"context"
	"learn-go/simple_api/common"
	"learn-go/simple_api/component/asyncjob"
	"learn-go/simple_api/modules/restaurantlike/restaurantlikemodel"
)

type UsersLikeRestaurantStore interface {
	Create(ctx context.Context, data *restaurantlikemodel.Like) error
	Find(ctx context.Context, condition map[string]interface{}) (*restaurantlikemodel.Like, error)
}

type IncreaseLikeCountStore interface {
	IncreaseLikeCount(ctx context.Context, id int) error
}

type userLikeRestaurantBiz struct {
	store         UsersLikeRestaurantStore
	increaseStore IncreaseLikeCountStore
}

func NewUserLikeRestaurantBiz(
	store UsersLikeRestaurantStore,
	increaseStore IncreaseLikeCountStore,
) *userLikeRestaurantBiz {
	return &userLikeRestaurantBiz{store: store, increaseStore: increaseStore}
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

	// This is consider as side effect so no need to handle error
	go func() {
		defer common.AppRecover()
		job := asyncjob.NewJob(func(ctx context.Context) error {
			return biz.increaseStore.IncreaseLikeCount(ctx, data.RestaurantId)
		})

		_ = asyncjob.NewGroup(true, job).Run(ctx)
	}()

	// This is consider as side effect so no need to handle error
	// go func() {
	// 	defer common.AppRecover()
	// 	_ = biz.increaseStore.IncreaseLikeCount(ctx, data.RestaurantId)
	// }()

	return nil
}
