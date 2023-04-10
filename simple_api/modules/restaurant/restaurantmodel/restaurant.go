package restaurantmodel

import (
	"errors"
	"learn-go/simple_api/common"
	"strings"
)

const Entity = "Restaurant"

type Restaurant struct {
	common.SQLModel `json:",inline"` // NOTE: json:",inline": to make it spread properties into Restaurant, not create a new SQLModel key
	Name            string           `json:"name" gorm:"column:name;"`
	Addr            string           `json:"address" gorm:"column:addr;"`
}

func (Restaurant) TableName() string {
	return "restaurants"
}

type RestaurantUpdate struct {
	Name *string `json:"name" gorm:"column:name;"`
	Addr *string `json:"address" gorm:"column:addr;"`
}

func (RestaurantUpdate) TableName() string {
	return Restaurant{}.TableName()
}

type RestaurantCreate struct {
	Id   int    `json:"id" gorm:"column:id;"`
	Name string `json:"name" gorm:"column:name;"`
	Addr string `json:"address" gorm:"column:addr;"`
}

func (RestaurantCreate) TableName() string {
	return Restaurant{}.TableName()
}

func (res *RestaurantCreate) Validate() error {
	res.Name = strings.TrimSpace(res.Name)

	if len(res.Name) == 0 {
		return errors.New("restaurant name cannot be blank")
	}

	return nil
}
