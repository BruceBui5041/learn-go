package userstorage

import (
	"context"
	"learn-go/food_delivery_be/common"
	"learn-go/food_delivery_be/modules/user/usermodel"
)

func (s *sqlStore) CreateUser(ctx context.Context, data *usermodel.CreateUser) error {
	db := s.db.Begin()

	if err := db.Table(data.TableName()).Create(data).Error; err != nil {
		db.Rollback()
		return common.ErrDB(err)
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return common.ErrDB(err)
	}

	return nil
}
