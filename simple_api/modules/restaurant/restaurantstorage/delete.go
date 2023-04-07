package restaurantstorage

import (
	"context"
	"learn-go/simple_api/common"
	"learn-go/simple_api/modules/restaurant/restaurantmodel"
)

func (s *sqlStore) SoftDeleteData(ctx context.Context, condition map[string]interface{}) error {
	db := s.db

	for key, value := range condition {
		restaurant := &restaurantmodel.Restaurant{}

		switch key {
		case "id":
			if id, ok := value.(int); ok {
				restaurant.Id = id
			}
		default:
			continue
		}

		db = db.Where(restaurant)
	}

	if err := db.
		Table(restaurantmodel.Restaurant{}.TableName()).
		Updates(map[string]interface{}{"status": 0}).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
