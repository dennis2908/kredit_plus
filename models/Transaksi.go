package models

import (
	"time"
	"github.com/astaxie/beego/orm"
)

type Transaksi struct {

	Id_konsumen    int
	No_kontrak      string
	Otr      int
	Admin_fee      int
	Jml_cicilan      int
	Jml_bunga     int
	Nama_aset     string

}

type Transaksis struct {
	Id           int    `orm:"auto;pk;index"`
	Transaksi
	Created_date time.Time `orm:"auto_now;type(datetime)"`
	Updated_date time.Time `orm:"auto_now;type(datetime)"`
}

type TransaksiWIdTrans struct{
	Id_transaksi int
	Transaksi
}

func (a *Transaksis) TableName() string {
	return "transaksis"
}


func init() {
    orm.RegisterModel(new(Transaksis))
}
