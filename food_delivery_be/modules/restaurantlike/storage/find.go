package restaurantlikestorage

import (
	"context"
	"learn-go/food_delivery_be/modules/restaurantlike/restaurantlikemodel"
)

func (s *sqlStore) Find(ctx context.Context, condition map[string]interface{}) (*restaurantlikemodel.Like, error) {
	db := s.db

	var like restaurantlikemodel.Like

	if err := db.Table(restaurantlikemodel.Like{}.TableName()).First(&like, condition).Error; err != nil {
		return nil, err
	}

	return &like, nil
}
