package restaurantlikestorage

import (
	"context"
	"learn-go/food_delivery_be/modules/restaurantlike/restaurantlikemodel"
)

func (s *sqlStore) Create(ctx context.Context, data *restaurantlikemodel.Like) error {
	db := s.db

	if err := db.Create(data).Error; err != nil {
		return err
	}

	return nil
}
