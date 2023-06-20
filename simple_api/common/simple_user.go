package common

// NOTE: This SimpleUser is to show only public information of user
type SimpleUser struct {
	SQLModel  `json:",inline"`
	LastName  string `json:"last_name" gorm:"column:last_name;"`
	FirstName string `json:"first_name" gorm:"column:first_name;"`
	Role      string `json:"role" gorm:"column:role;"`
	Avatar    *Image `json:"avatar,omitempty" gorm:"column:avatar;type:json"`
}

func (u *SimpleUser) Mask(isAdmin bool) {
	u.GenUID(DbTypeUser)
	return
}

func (SimpleUser) TableName() string {
	return "users"
}
