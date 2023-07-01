package restaurantlikebiz

import (
	"context"
	"learn-go/simple_api/common"
	"learn-go/simple_api/modules/restaurantlike/restaurantlikemodel"
)

type UsersLikeRestaurantStore interface {
	Create(ctx context.Context, data *restaurantlikemodel.Like) error
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

	err := biz.store.Create(ctx, data)

	if err != nil {
		return restaurantlikemodel.ErrCannotLikeRestaurant(err)
	}

	// This is consider as side effect so no need to handle error
	go func() {
		defer common.AppRecover()
		_ = biz.increaseStore.IncreaseLikeCount(ctx, data.RestaurantId)
	}()

	return nil
}
