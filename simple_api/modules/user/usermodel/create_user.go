package usermodel

import "learn-go/simple_api/common"

type CreateUser struct {
	common.SQLModel `json:",inline"`
	Email           string        `json:"email" gorm:"column:email;"`
	Password        string        `json:"password" gorm:"column:password;"`
	LastName        string        `json:"last_name" gorm:"column:last_name;"`
	FirstName       string        `json:"first_name" gorm:"column:first_name;"`
	Role            string        `json:"-" gorm:"column:role;"`
	Salt            string        `json:"-" gorm:"column:salt;"`
	Avatar          *common.Image `json:"avatar,omitempty" gorm:"column:avatar;type:json"`
}

func (u *CreateUser) Mask(isAdmin bool) {
	u.GenUID(common.DbTypeUser)
	return
}

func (CreateUser) TableName() string {
	return User{}.TableName()
}
