package restaurantstorage

import (
	"context"
	"learn-go/simple_api/common"
	"learn-go/simple_api/modules/restaurant/restaurantmodel"

	"gorm.io/gorm"
)

func (s *sqlStore) GetRestaurant(
	ctx context.Context,
	condition map[string]interface{},
	moreKeys ...string,
) (*restaurantmodel.Restaurant, error) {
	var result restaurantmodel.Restaurant
	db := s.db

	// Loop through the condition map
	for key, value := range condition {
		// Initialize an empty restaurant model
		restaurant := &restaurantmodel.Restaurant{}

		// Use a switch statement to set the corresponding field of the restaurant model
		switch key {
		case "id":
			if id, ok := value.(int); ok {
				restaurant.Id = id
			}
		case "name":
			if name, ok := value.(string); ok {
				restaurant.Name = name
			}
		case "address":
			if address, ok := value.(string); ok {
				restaurant.Addr = address
			}
		default:
			// If the key is not supported, continue to the next iteration
			continue
		}

		// Update the query with the restaurant model
		db = db.Where(restaurant)
	}

	if err := db.First(&result).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.RecordNotFound
		}
		return nil, err
	}

	return &result, nil
}
