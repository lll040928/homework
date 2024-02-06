package model

type UserInfo struct {
	Uid      int    `form:"uid,omitempty" gorm:"primary_key"`
	Username string `form:"username,omitempty" gorm:"not null,index:idx_username,unique"`
	Password string `form:"password,omitempty" gorm:"not null"`
	Phone    string `form:"phone,omitempty" gorm:"not null"`
	Address  string `form:"address,omitempty" gorm:"default:'未填写'"`
	Role     int    `form:"role,omitempty" gorm:"default:0"`
}
