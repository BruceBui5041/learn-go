package restaurantbiz

import (
	"context"
	"errors"
	"learn-go/simple_api/common"
	"learn-go/simple_api/modules/restaurant/restaurantmodel"
)

type GetRestaurantStore interface {
	GetRestaurant(ctx context.Context, restaurant map[string]interface{}, moreKeys ...string) (*restaurantmodel.Restaurant, error)
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
	condition := map[string]interface{}{
		"id": restaurantId,
	}

	result, err := biz.store.GetRestaurant(ctx, condition, moreKeys...)

	if err != nil {
		if err == common.RecordNotFound {
			return nil, err
		}

		return nil, err
	}

	if result.Status == 0 {
		return nil, errors.New("data deleted")
	}

	return result, err
}