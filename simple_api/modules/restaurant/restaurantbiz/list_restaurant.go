package restaurantbiz

import (
	"context"
	"learn-go/simple_api/common"
	"learn-go/simple_api/modules/restaurant/restaurantmodel"
)

type ListRestaurantRepo interface {
	ListRestaurant(
		ctx context.Context,
		filter *restaurantmodel.Filter,
		paging *common.Paging,
		moreKeys ...string,
	) ([]restaurantmodel.Restaurant, error)
}

type listRestaurantBiz struct {
	retaurantRepo ListRestaurantRepo
}

// /	NOTE: *listRestaurantBiz | Return về con trỏ là để tối ưu cho khỏi copy rồi return về 1 copy object của listRestaurantBiz
func NewListRestaurantBiz(repo ListRestaurantRepo) *listRestaurantBiz {
	return &listRestaurantBiz{retaurantRepo: repo}
}

func (biz *listRestaurantBiz) ListRestaurant(
	ctx context.Context,
	filter *restaurantmodel.Filter,
	paging *common.Paging,
	moreKeys ...string,
) ([]restaurantmodel.Restaurant, error) {
	// NOTE: "User" is the key in the Restaurant struct not the table name of User struct
	restaurants, err := biz.retaurantRepo.ListRestaurant(ctx, filter, paging)

	if err != nil {
		return nil, common.ErrCannotListEntity(restaurantmodel.Entity, err)
	}

	return restaurants, nil
}
