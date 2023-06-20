package restaurantmodel

import (
	"errors"
	"learn-go/simple_api/common"
	"strings"
)

const Entity = "Restaurant"

type Restaurant struct {
	common.SQLModel `json:",inline"`   // NOTE: json:",inline": to make it spread properties into Restaurant, not create a new SQLModel key
	Name            string             `json:"name" gorm:"column:name;"`
	UserId          int                `json:"-" gorm:"column:owner_id"`
	Addr            string             `json:"address" gorm:"column:addr;"`
	Logo            *common.Image      `json:"logo" gorm:"column:logo;"`
	Cover           *common.Images     `json:"cover" gorm:"column:cover;"`
	LikeCount       int                `json:"like_count" gorm:"-"`
	User            *common.SimpleUser `json:"user" gorm:"preload:false;"`
}

func (Restaurant) TableName() string {
	return "restaurants"
}

type RestaurantUpdate struct {
	Name *string       `json:"name" gorm:"column:name;"`
	Addr *string       `json:"address" gorm:"column:addr;"`
	Logo *common.Image `json:"logo" gorm:"column:logo;"`
	// NOTE: Nếu sử dụng []Image thì sẽ gặp lỗi nên phải tạo alias Images type trong common để sử dụng
	Cover  *common.Images `json:"cover" gorm:"column:cover;"`
	UserId int            `json:"-" gorm:"column:owner_id"`
}

func (RestaurantUpdate) TableName() string {
	return Restaurant{}.TableName()
}

type RestaurantCreate struct {
	common.SQLModel `json:",inline"` // NOTE: json:",inline": to make it spread properties into Restaurant, not create a new SQLModel key
	Name            string           `json:"name" gorm:"column:name;"`
	Addr            string           `json:"address" gorm:"column:addr;"`
	Logo            *common.Image    `json:"logo" gorm:"column:logo;"`
	Cover           *common.Images   `json:"cover" gorm:"column:cover;"`
	UserId          int              `json:"-" gorm:"column:owner_id"`
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

// This to parse DB id into custome UID before send back to client
func (data *Restaurant) Mask(isAdminOrOwner bool) {
	data.GenUID(common.DbTypeRestaurant)

	if user := data.User; user != nil {
		data.User.Mask(false)
	}
}

// This to parse DB id into custome UID before send back to client
func (data *RestaurantCreate) Mask(isAdminOrOwner bool) {
	data.GenUID(common.DbTypeRestaurant)
}
