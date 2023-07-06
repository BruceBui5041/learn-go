package restaurantrepo

import (
	"context"
	"learn-go/food_delivery_be/common"
	"learn-go/food_delivery_be/modules/restaurant/restaurantmodel"
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

type listRestaurantRepo struct {
	store     ListRestaurantStore
	likeStore LikeStore
}

func NewRestaurantRepo(store ListRestaurantStore, likeStore LikeStore) *listRestaurantRepo {
	return &listRestaurantRepo{store: store, likeStore: likeStore}
}

func (restaurantRepo *listRestaurantRepo) ListRestaurant(
	ctx context.Context,
	filter *restaurantmodel.Filter,
	paging *common.Paging,
	moreKeys ...string,
) ([]restaurantmodel.Restaurant, error) {
	// NOTE: "User" is the key in the Restaurant struct not the table name of User struct
	restaurants, err := restaurantRepo.store.ListDataByCondition(ctx, nil, filter, paging, "User")

	if err != nil {
		return nil, common.ErrCannotListEntity(restaurantmodel.Entity, err)
	}

	// ids := make([]int, len(restaurants))

	// for i := range restaurants {
	// 	ids[i] = restaurants[i].Id
	// }

	// mapRestaurantAndLikes, err := restaurantRepo.likeStore.GetRestaurantLikes(ctx, ids)

	// // Should not return error when just cannot get likes of a restaurant because it not important
	// if err != nil {
	// 	log.Println("cannot get restaurant likes: ", err)
	// }

	// if value := mapRestaurantAndLikes; value != nil {
	// 	for i, restaurant := range restaurants {
	// 		restaurants[i].LikeCount = mapRestaurantAndLikes[restaurant.Id]
	// 	}
	// }

	return restaurants, nil
}
