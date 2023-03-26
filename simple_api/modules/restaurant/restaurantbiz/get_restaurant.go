package restaurantbiz

import (
	"context"
	"learn-go/simple_api/modules/restaurant/restaurantmodel"
)

type GetRestaurantStore interface {
	GetRestaurantById(ctx context.Context, restaurantId int, moreKeys ...string) (*restaurantmodel.Restaurant, error)
}

type getRestaurantBiz struct {
	store GetRestaurantStore
}

func NewGetRestaurantBiz(store GetRestaurantStore) *getRestaurantBiz {
	return &getRestaurantBiz{store: store}
}

func (biz *getRestaurantBiz) GetRestaurantById(
	ctx context.Context,
	restaurantId int,
	moreKeys ...string,
) (*restaurantmodel.Restaurant, error) {
	result, err := biz.store.GetRestaurantById(ctx, restaurantId, moreKeys...)

	return result, err
}
