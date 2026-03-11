package controllers

import (
	_ "context"
	"encoding/json"
	_ "fmt"
	"kredit_plus/models"
	"kredit_plus/structs"

	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
	_ "github.com/leekchan/accounting"
	_ "github.com/shopspring/decimal"
)

type TransaksiDetailsController struct {
	beego.Controller
}

func (api *TransaksiDetailsController) GetAllTransaksiDetails() {
	o := orm.NewOrm()
	o.Using("default")

	var Transaksis []structs.GetAllTransaksi
	sqlTrans := "Select transaksis.*,konsumens.*,transaksis.id_konsumen as konsumen_id, transaksis.id as Id_transaksi FROM konsumens join transaksis on transaksis.id_konsumen = konsumens.id"
	numTrans, errTrans := o.Raw(sqlTrans).QueryRows(&Transaksis)
	sqlTransDeta := "Select transaksi_details.* FROM transaksi_details join transaksis on transaksis.id = transaksi_details.id__transaksi"
	var Transaksi_detail []models.Transaksi_details
	numTransDeta, errTransDeta := o.Raw(sqlTransDeta).QueryRows(&Transaksi_detail)
	Transaksi_details := structs.GetTransaksiDetails{Transaksis: Transaksis, Transaksi_details: Transaksi_detail}
	if errTrans == nil && numTrans > 0 && errTransDeta == nil && numTransDeta > 0 {
		api.Data["json"] = Transaksi_details
	}
	api.ServeJSON()
}

func (api *TransaksiDetailsController) GetTransaksiDetailsByIdKonsumen() {
	o := orm.NewOrm()
	o.Using("default")

	var Transaksis []structs.GetAllTransaksi
	idInt, _ := strconv.Atoi(api.Ctx.Input.Param(":id"))
	sqlTrans := "Select transaksis.*,konsumens.*,transaksis.id_konsumen as konsumen_id, transaksis.id as Id_transaksi FROM konsumens join transaksis on transaksis.id_konsumen = konsumens.id"
	sqlTrans += " where konsumens.id = ?"
	numTrans, errTrans := o.Raw(sqlTrans, idInt).QueryRows(&Transaksis)
	sqlTransDeta := "Select transaksi_details.* FROM transaksi_details join transaksis on transaksis.id = transaksi_details.id__transaksi"
	sqlTransDeta += " where id_konsumen = ?"
	var Transaksi_detail []models.Transaksi_details
	numTransDeta, errTransDeta := o.Raw(sqlTransDeta, idInt).QueryRows(&Transaksi_detail)
	Transaksi_details := structs.GetTransaksiDetails{Transaksis: Transaksis, Transaksi_details: Transaksi_detail}
	if errTrans == nil && numTrans > 0 && errTransDeta == nil && numTransDeta > 0 {
		api.Data["json"] = Transaksi_details
	}
	api.ServeJSON()
}

func AllTransaksiDetailsCheck(api *TransaksiDetailsController) string {
	valid := validation.Validation{}

	frm := api.Ctx.Input.RequestBody
	ul := &models.Transaksi_details{}
	json.Unmarshal(frm, ul)
	valid.Required(ul.Id_Transaksi, "Id_Transaksi")
	valid.Required(ul.Bulan, "Bulan")
	valid.Required(ul.Cicilan, "Cicilan")

	if valid.HasErrors() {
		// If there are error messages it means the validation didn't pass
		// Print error message
		for _, err := range valid.Errors {
			return err.Key + err.Message
		}
	}

	return ""
}

func (api *TransaksiDetailsController) CreateTransaksiDetails() {
	frm := api.Ctx.Input.RequestBody
	if AllTransaksiDetailsCheck(api) != "" {
		api.Data["json"] = AllTransaksiDetailsCheck(api)
		api.ServeJSON()
		return
	}

	o := orm.NewOrm()
	o.Using("default")

	ul := &models.Transaksi_details{}
	json.Unmarshal(frm, ul)
	Qry := models.Transaksi_details{Transaksi_detail: models.Transaksi_detail{Id_Transaksi: ul.Id_Transaksi, Bulan: ul.Bulan, Cicilan: ul.Cicilan}}
	o.Insert(&Qry)
	api.Data["json"] = "Successfully save data"
	api.ServeJSON()
}

func (api *TransaksiDetailsController) UpdateTransaksiDetails() {
	frm := api.Ctx.Input.RequestBody
	if AllTransaksiDetailsCheck(api) != "" {
		api.Data["json"] = AllTransaksiDetailsCheck(api)
		api.ServeJSON()
		return
	}

	o := orm.NewOrm()
	o.Using("default")

	idInt, _ := strconv.Atoi(api.Ctx.Input.Param(":id"))

	ul := &models.Transaksi_details{}
	json.Unmarshal(frm, ul)
	Qry := models.Transaksi_details{Id: idInt, Transaksi_detail: models.Transaksi_detail{Id_Transaksi: ul.Id_Transaksi, Bulan: ul.Bulan, Cicilan: ul.Cicilan}}
	o.Update(&Qry)

	api.Data["json"] = "Successfully update data with id  =  " + api.Ctx.Input.Param(":id")

	api.ServeJSON()
}
