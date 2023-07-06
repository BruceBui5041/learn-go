package usermodel

import "learn-go/food_delivery_be/common"

const EntityName = "User"

type User struct {
	common.SQLModel `json:",inline"`
	Email           string        `json:"email" gorm:"column:email;"`
	Password        string        `json:"-" gorm:"column:password;"`
	LastName        string        `json:"last_name" gorm:"column:last_name;"`
	FirstName       string        `json:"first_name" gorm:"column:first_name;"`
	Role            string        `json:"role" gorm:"column:role;"`
	Phone           string        `json:"phone" gorm:"column:phone;"`
	Salt            string        `json:"-" gorm:"column:salt;"`
	Avatar          *common.Image `json:"avatar,omitempty" gorm:"column:avatar;type:json"`
}

func (user *User) GetUserId() int {
	return user.Id
}

func (user *User) GetEmail() string {
	return user.Email
}

func (user *User) GetRole() string {
	return user.Role
}

func (User) TableName() string {
	return "users"
}

func (u *User) Mask(isAdmin bool) {
	u.GenUID(common.DbTypeUser)
	return
}
