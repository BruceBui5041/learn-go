package restaurantlikemodel

import (
	"learn-go/simple_api/common"
	"time"
)

const EntityName = "RestaurantLike"

type Like struct {
	RestaurantId int                `json:"restaurant_id" gorm:"column:restaurant_id"`
	UserId       int                `json:"user_id" gorm:"column:user_id"`
	CreatedAt    *time.Time         `json:"created_at" gorm:"column:created_at"`
	User         *common.SimpleUser `json:"user" gorm:"preload:false"`
}

func (like Like) TableName() string {
	return "restaurant_likes"
}
