package controllers

import (
	_ "context"
	"encoding/json"
	_ "fmt"
	models "kredit_plus/models"
	"kredit_plus/structs"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
	_ "github.com/leekchan/accounting"
	_ "github.com/shopspring/decimal"
)

type TransaksisController struct {
	beego.Controller
}

type merged struct {
	models.Konsumens
	models.Transaksis
}

func (api *TransaksisController) GetAllTransaksis() {
	o := orm.NewOrm()
	o.Using("default")

	var Transaksis []structs.GetAllTransaksi
	sql := "Select transaksis.*,konsumens.*,transaksis.id_konsumen as konsumen_id, transaksis.id as Id_transaksi FROM konsumens join transaksis on transaksis.id_konsumen = konsumens.id"
	num, err := o.Raw(sql).QueryRows(&Transaksis)
	if err == nil && num > 0 {
		api.Data["json"] = Transaksis
	}
	api.ServeJSON()
}

func (api *TransaksisController) GetTransaksiByIdKons() {
	o := orm.NewOrm()
	o.Using("default")

	var Transaksis []structs.GetAllTransaksi
	idInt, _ := strconv.Atoi(api.Ctx.Input.Param(":id"))
	sql := "Select transaksis.*,konsumens.*,transaksis.id_konsumen as konsumen_id, transaksis.id as Id_transaksi FROM konsumens join transaksis on transaksis.id_konsumen = konsumens.id"
	sql += " where konsumens.id = ?"
	num, err := o.Raw(sql, idInt).QueryRows(&Transaksis)
	if err == nil && num > 0 {
		api.Data["json"] = Transaksis
	}
	api.ServeJSON()
}

func (api *TransaksisController) GetTransaksisById() {
	o := orm.NewOrm()
	o.Using("default")
	var Transaksis []models.Transaksis
	idInt, _ := strconv.Atoi(api.Ctx.Input.Param(":id"))
	num, err := o.QueryTable("Transaksis").Filter("id", idInt).All(&Transaksis)
	if err == nil && num > 0 {
		api.Data["json"] = Transaksis
	}
	api.ServeJSON()
}

func AllTransaksisCheck(api *TransaksisController) string {
	valid := validation.Validation{}

	frm := api.Ctx.Input.RequestBody
	ul := &models.Transaksis{}
	json.Unmarshal(frm, ul)

	valid.Required(ul.Id_konsumen, "Id_konsumen")
	valid.Required(ul.No_kontrak, "No_kontrak")
	valid.Required(ul.Otr, "Otr")
	valid.Required(ul.Admin_fee, "Admin_fee")
	valid.Required(ul.Jml_cicilan, "Jml_cicilan")
	valid.Required(ul.Jml_bunga, "Jml_bunga")
	valid.Required(ul.Nama_aset, "Nama_aset")

	if valid.HasErrors() {
		// If there are error messages it means the validation didn't pass
		// Print error message
		for _, err := range valid.Errors {
			return err.Key + err.Message
		}
	}

	return ""
}

func (api *TransaksisController) CreateTransaksis() {
	frm := api.Ctx.Input.RequestBody
	if AllTransaksisCheck(api) != "" {
		api.Data["json"] = AllTransaksisCheck(api)
		api.ServeJSON()
		return
	}

	o := orm.NewOrm()
	o.Using("default")

	ul := &models.Transaksis{}
	json.Unmarshal(frm, ul)
	Qry := models.Transaksis{Transaksi: models.Transaksi{Id_konsumen: ul.Id_konsumen, No_kontrak: ul.No_kontrak, Otr: ul.Otr, Admin_fee: ul.Admin_fee, Jml_cicilan: ul.Jml_cicilan, Jml_bunga: ul.Jml_bunga, Nama_aset: ul.Nama_aset}}
	o.Insert(&Qry)
	api.Data["json"] = "Successfully save data"
	api.ServeJSON()
}

func (api *TransaksisController) UpdateTransaksis() {
	frm := api.Ctx.Input.RequestBody
	if AllTransaksisCheck(api) != "" {
		api.Data["json"] = AllTransaksisCheck(api)
		api.ServeJSON()
		return
	}

	o := orm.NewOrm()
	o.Using("default")

	idInt, _ := strconv.Atoi(api.Ctx.Input.Param(":id"))

	ul := &models.Transaksis{}
	json.Unmarshal(frm, ul)
	Qry := models.Transaksis{Id: idInt, Transaksi: models.Transaksi{Id_konsumen: ul.Id_konsumen, No_kontrak: ul.No_kontrak, Otr: ul.Otr, Admin_fee: ul.Admin_fee, Jml_cicilan: ul.Jml_cicilan, Jml_bunga: ul.Jml_bunga, Nama_aset: ul.Nama_aset}}
	o.Update(&Qry)
	api.Data["json"] = "Successfully update data with id  =  " + api.Ctx.Input.Param(":id")

	api.ServeJSON()
}
