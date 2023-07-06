package uploadstorage

import (
	"context"
	"learn-go/food_delivery_be/common"
)

func (store *sqlStore) CreateImage(c context.Context, data *common.Image) error {
	db := store.db

	if err := db.Table(data.TableName()).Create(data).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
