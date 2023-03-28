package restaurantbiz

import (
	"context"
	"errors"
	"learn-go/simple_api/modules/restaurant/restaurantmodel"
)

type UpdateRestaurantStore interface {
	GetRestaurant(ctx context.Context, restaurant map[string]interface{}, moreKeys ...string) (*restaurantmodel.Restaurant, error)
	UpdateData(ctx context.Context, id int, data *restaurantmodel.RestaurantUpdate) error
}

type updateRestaurantBiz struct {
	store UpdateRestaurantStore
}

func NewUpdateRestaurantBiz(store UpdateRestaurantStore) *updateRestaurantBiz {
	return &updateRestaurantBiz{store: store}
}

func (biz *updateRestaurantBiz) UpdateRestaurant(
	ctx context.Context,
	restaurantId int,
	data *restaurantmodel.RestaurantUpdate,
	moreKeys ...string,
) error {

	oldData, err := biz.store.GetRestaurant(
		ctx,
		map[string]interface{}{"id": restaurantId},
		moreKeys...,
	)

	if err != nil {
		return err
	}

	if oldData.Status == 0 {
		return errors.New("data deleted")
	}

	if err := biz.store.UpdateData(ctx, restaurantId, data); err != nil {
		return err
	}

	return nil
}
