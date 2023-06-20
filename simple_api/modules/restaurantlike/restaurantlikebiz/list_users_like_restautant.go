package restaurantlikebiz

import (
	"context"
	"learn-go/simple_api/common"
	"learn-go/simple_api/modules/restaurantlike/restaurantlikemodel"
)

type ListUsersLikeRestaurantStore interface {
	GetUsersLikeRestaurant(
		ctx context.Context,
		conditions map[string]interface{},
		filter *restaurantlikemodel.Filter,
		paging *common.Paging,
		moreKeys ...string,
	) ([]common.SimpleUser, error)
}

type listUserLikeRestaurantBiz struct {
	store ListUsersLikeRestaurantStore
}

func NewListUserLikeRestaurantBiz(store ListUsersLikeRestaurantStore) *listUserLikeRestaurantBiz {
	return &listUserLikeRestaurantBiz{store: store}
}

func (biz *listUserLikeRestaurantBiz) ListUserLikeRestaurant(
	ctx context.Context,
	filter *restaurantlikemodel.Filter,
	paging *common.Paging,
	moreKeys ...string,
) ([]common.SimpleUser, error) {
	users, err := biz.store.GetUsersLikeRestaurant(ctx, nil, filter, paging)

	if err != nil {
		return nil, common.ErrCannotListEntity(restaurantlikemodel.EntityName, err)
	}

	return users, nil
}
