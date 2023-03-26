package restaurantstorage

import (
	"context"
	"learn-go/simple_api/modules/restaurant/restaurantmodel"
)

func (s *sqlStore) GetRestaurantById(
	ctx context.Context,
	restaurantId int,
	moreKeys ...string,
) (*restaurantmodel.Restaurant, error) {
	var result restaurantmodel.Restaurant
	db := s.db

	db = db.Table(restaurantmodel.Restaurant{}.TableName()).Where("id = ?", restaurantId)

	if err := db.First(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}
