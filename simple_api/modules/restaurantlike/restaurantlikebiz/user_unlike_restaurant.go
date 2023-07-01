package restaurantlikebiz

import (
	"context"
	"learn-go/simple_api/common"
	"learn-go/simple_api/modules/restaurantlike/restaurantlikemodel"
)

type UserUnlikeRestaurantStore interface {
	Delete(ctx context.Context, userId, restaurantId int) error
}

type DecreaseLikeCountStore interface {
	DecreaseLikeCount(ctx context.Context, id int) error
}

type userUnlikeRestaurantBiz struct {
	store             UserUnlikeRestaurantStore
	decreaseLikeStore DecreaseLikeCountStore
}

func NewUserUnlikeRestaurantBiz(store UserUnlikeRestaurantStore, decreaseLikeStore DecreaseLikeCountStore) *userUnlikeRestaurantBiz {
	return &userUnlikeRestaurantBiz{store: store, decreaseLikeStore: decreaseLikeStore}
}

func (biz *userUnlikeRestaurantBiz) UserUnlikeRestaurant(
	ctx context.Context,
	userId, restaurantId int,
) error {
	// NOTE: nên find coi coi user đã like restaurant chưa trước khi create
	// vì tìm theo 2 khoá chính user_id và restaurant_id nên sẽ rất nhanh => nên làm

	err := biz.store.Delete(ctx, userId, restaurantId)

	if err != nil {
		return restaurantlikemodel.ErrCannotUnlikeRestaurant(err)
	}

	// This is consider as side effect so no need to handle error
	go func() {
		defer common.AppRecover()
		_ = biz.decreaseLikeStore.DecreaseLikeCount(ctx, restaurantId)
	}()

	return nil
}
