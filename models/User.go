package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type User struct {
	Id           int       `orm:"auto;pk" json:"id"`
	Email        string    `orm:"size(255);unique" json:"email"`
	PasswordHash string    `orm:"size(255)" json:"-"`
	FullName     string    `orm:"size(255)" json:"full_name"`
	IsActive     bool      `orm:"default(true)" json:"is_active"`
	CreatedAt    time.Time `orm:"auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt    time.Time `orm:"auto_now;type(datetime)" json:"updated_at"`
}

func (u *User) TableName() string {
	return "users"
}

func init() {
	orm.RegisterModel(new(User))
}
