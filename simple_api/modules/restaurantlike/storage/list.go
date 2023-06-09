package restaurantlikestorage

import (
	"context"
	"learn-go/simple_api/common"
	"learn-go/simple_api/modules/restaurantlike/restaurantlikemodel"
)

func (s sqlStore) GetRestaurantLikes(ctx context.Context, ids []int) (map[int]int, error) {
	result := make(map[int]int)

	type sqlData struct {
		RestaurantId int `gorm:"column:restaurant_id;"`
		LikeCount    int `gorm:"column:count;"`
	}

	var listLikes []sqlData

	if err := s.db.Table(restaurantlikemodel.Like{}.TableName()).
		Select("restaurant_id, count(restaurant_id) as count").
		Where("restaurant_id in (?)", ids).
		Group("restaurant_id").Find(&listLikes).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	for _, item := range listLikes {
		result[item.RestaurantId] = item.LikeCount
	}

	return result, nil
}
