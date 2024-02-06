package model

import "time"

type Wallet struct {
	Username string  `form:"username,omitempty" gorm:"not null"`
	Balance  float64 `form:"balance,omitempty" gorm:"default :0"`
}
type Cart struct {
	Username string  `form:"username,omitempty" gorm:"not null"`
	Gid      int     `form:"gid,omitempty" gorm:"not null"`
	Gname    string  `form:"gname,omitempty" gorm:"not null"`
	Price    float64 `form:"price,omitempty" gorm:"not null"`
	Count    int     `form:"count,omitempty" gorm:"not null"`
}

type Order struct {
	Oid       int       `json:"oid,omitempty" gorm:"primary_key"`
	Username  string    `json:"username" gorm:"not null"`
	Gname     string    `json:"gname,omitempty" gorm:"not null"`
	Price     float64   `json:"price,omitempty" gorm:"not null"`
	Count     int       `json:"count,omitempty" gorm:"not null"`
	OrderTime time.Time `json:"order_time" gorm:"not null"`
}
