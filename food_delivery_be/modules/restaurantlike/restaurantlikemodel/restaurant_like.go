package restaurantlikemodel

import (
	"learn-go/food_delivery_be/common"
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

func (like Like) GetRestaurantId() int {
	return like.RestaurantId
}

func (like Like) GetUserId() int {
	return like.UserId
}

func ErrCannotLikeRestaurant(err error) *common.AppError {
	return common.NewCustomError(
		err,
		"Cannot like restaurant",
		"ErrCannotLikeRestaurant",
	)
}

func ErrCannotUnlikeRestaurant(err error) *common.AppError {
	return common.NewCustomError(
		err,
		"Cannot unlike restaurant",
		"ErrCannotUnlikeRestaurant",
	)
}
