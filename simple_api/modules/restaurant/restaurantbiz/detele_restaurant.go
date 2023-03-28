package restaurantbiz

import (
	"context"
	"errors"
	"learn-go/simple_api/modules/restaurant/restaurantmodel"
)

type DeleteRestaurantStore interface {
	GetRestaurant(ctx context.Context, condition map[string]interface{}, moreKeys ...string) (*restaurantmodel.Restaurant, error)
	SoftDeleteData(ctx context.Context, condition map[string]interface{}) error
}

type deleteRestaurantBiz struct {
	store DeleteRestaurantStore
}

func NewDeleteRestaurantBiz(store DeleteRestaurantStore) *deleteRestaurantBiz {
	return &deleteRestaurantBiz{store: store}
}

func (biz *deleteRestaurantBiz) SoftDeleteRestaurant(
	ctx context.Context,
	restaurantId int,
	moreKeys ...string,
) error {
	var condition = map[string]interface{}{"id": restaurantId}
	oldData, err := biz.store.GetRestaurant(
		ctx,
		condition,
		moreKeys...,
	)

	if err != nil {
		return err
	}

	if oldData.Status == 0 {
		return errors.New("data deleted")
	}

	if err := biz.store.SoftDeleteData(ctx, condition); err != nil {
		return err
	}

	return nil
}
