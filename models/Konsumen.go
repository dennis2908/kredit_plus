package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type Konsumen struct {
	Email       string
	Nik         string
	Full_name   string
	Legal_name  string
	Place_birth string
	Salary      int
	Foto_ktp    string
	Foto_selfie string
}

type KonsumenId struct {
	Id int
}

type Konsumens struct {
	Id int `orm:"auto;pk;index"`
	Konsumen
	Date_birth   time.Time `orm:"type(date)"`
	Created_date time.Time `orm:"auto_now;type(datetime)"`
	Updated_date time.Time `orm:"auto_now;type(datetime)"`
}

type KonsumensWIdKons struct {
	Konsumen_id int
	Konsumen
}

func (a *Konsumens) TableName() string {
	return "konsumens"
}

func init() {
	orm.RegisterModel(new(Konsumens))
}
