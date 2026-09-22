package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type User struct {
	Id           int       `orm:"auto;pk;index" json:"id"`
	Username     string    `orm:"unique;size(100)" json:"username"`
	PasswordHash string    `orm:"size(255)" json:"-"`
	Role         string    `orm:"size(50)" json:"role"`
	CreatedAt    time.Time `orm:"auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt    time.Time `orm:"auto_now;type(datetime)" json:"updated_at"`
}

func (user *User) TableName() string {
	return "users"
}

func init() {
	orm.RegisterModel(new(User))
}
