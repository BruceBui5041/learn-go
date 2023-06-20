package restaurantbiz

import (
	"context"
	"learn-go/simple_api/common"
	"learn-go/simple_api/modules/restaurant/restaurantmodel"
	"log"
)

type ListRestaurantStore interface {
	ListDataByCondition(
		ctx context.Context,
		conditions map[string]interface{},
		filter *restaurantmodel.Filter,
		paging *common.Paging,
		moreKeys ...string,
	) ([]restaurantmodel.Restaurant, error)
}

type LikeStore interface {
	GetRestaurantLikes(ctx context.Context, ids []int) (map[int]int, error)
}

type listRestaurantBiz struct {
	store     ListRestaurantStore
	likeStore LikeStore
}

// /	NOTE: *listRestaurantBiz | Return về con trỏ là để tối ưu cho khỏi copy rồi return về 1 copy object của listRestaurantBiz
func NewListRestaurantBiz(store ListRestaurantStore, likeStore LikeStore) *listRestaurantBiz {
	return &listRestaurantBiz{store: store, likeStore: likeStore}
}

func (biz *listRestaurantBiz) ListRestaurant(
	ctx context.Context,
	filter *restaurantmodel.Filter,
	paging *common.Paging,
	moreKeys ...string,
) ([]restaurantmodel.Restaurant, error) {
	// NOTE: "User" is the key in the Restaurant struct not the table name of User struct
	restaurants, err := biz.store.ListDataByCondition(ctx, nil, filter, paging, "User")

	if err != nil {
		return nil, common.ErrCannotListEntity(restaurantmodel.Entity, err)
	}

	ids := make([]int, len(restaurants))

	for i := range restaurants {
		ids[i] = restaurants[i].Id
	}

	mapRestaurantAndLikes, err := biz.likeStore.GetRestaurantLikes(ctx, ids)

	// Should not return error when just cannot get likes of a restaurant because it not important
	if err != nil {
		log.Println("cannot get restaurant likes: ", err)
	}

	if value := mapRestaurantAndLikes; value != nil {
		for i, restaurant := range restaurants {
			restaurants[i].LikeCount = mapRestaurantAndLikes[restaurant.Id]
		}
	}

	return restaurants, nil
}
