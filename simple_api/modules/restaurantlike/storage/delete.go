package restaurantlikestorage

import (
	"context"
	"learn-go/simple_api/modules/restaurantlike/restaurantlikemodel"
)

func (s *sqlStore) Delete(ctx context.Context, userId, restaurantId int) error {
	db := s.db

	if err := db.Table(restaurantlikemodel.Like{}.TableName()).
		Where("user_id = ? and restaurant_id = ?", userId, restaurantId).
		Delete(nil).
		Error; err != nil {
		return err
	}

	return nil
}
