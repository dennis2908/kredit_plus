package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type Transaksi_detail struct {
	Id_Transaksi int
	Bulan        int
	Cicilan      int
}

type Transaksi_details struct {
	Id int `orm:"auto;pk;index"`
	Transaksi_detail
	Created_date time.Time `orm:"auto_now;type(datetime)"`
	Updated_date time.Time `orm:"auto_now;type(datetime)"`
}

func (a *Transaksi_details) TableName() string {
	return "transaksi_details"
}

func init() {
	orm.RegisterModel(new(Transaksi_details))
}
