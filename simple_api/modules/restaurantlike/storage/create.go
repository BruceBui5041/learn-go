package restaurantlikestorage

import (
	"context"
	"learn-go/simple_api/modules/restaurantlike/restaurantlikemodel"
)

func (s *sqlStore) Create(ctx context.Context, data *restaurantlikemodel.Like) error {
	db := s.db

	if err := db.Create(data).Error; err != nil {
		return err
	}

	return nil
}
