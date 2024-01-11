package controllers

import (
	models "kredit_plus/models"
	"kredit_plus/structs"
	_ "context"
	"encoding/json"
	_ "fmt"
	"time"

	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
	_ "github.com/leekchan/accounting"
	_ "github.com/shopspring/decimal"
)

type KonsumensController struct {
	beego.Controller
}

func (api *KonsumensController) GetAllKonsumens() {
	o := orm.NewOrm()
	o.Using("default")
	var konsumens []models.Konsumens
	num, err := o.QueryTable("konsumens").All(&konsumens)
	if err == nil && num > 0 {
		api.Data["json"] = konsumens	
	}
	api.ServeJSON()
}

func (api *KonsumensController) GetKonsumensById() {
	o := orm.NewOrm()
	o.Using("default")
	var konsumens []models.Konsumens
	idInt, _ := strconv.Atoi(api.Ctx.Input.Param(":id"))
	num, err := o.QueryTable("konsumens").Filter("id",idInt).All(&konsumens)
	if err == nil && num > 0 {
		api.Data["json"] = konsumens	
	}
	api.ServeJSON()
}

func AllKonsumensCheck(api *KonsumensController) string {
	valid := validation.Validation{}

	frm := api.Ctx.Input.RequestBody
	ul := &structs.GetKonsumen{}
	json.Unmarshal(frm, ul)
	// beego.Debug(ul)
	// Nik := ul.Nik
	// Full_name := ul.Full_name
	// Legal_name := ul.Legal_name
	// Place_birth := ul.Place_birth
	// Date_birth := ul.Date_birth
	// Salary := ul.Salary
	// Foto_ktp := ul.Foto_ktp
	// Foto_selfie := ul.Foto_selfie

	// u := structs.InsertKonsumen{Konsumen : models.Konsumen{Nik : ul.Nik, Full_name: ul.Full_name, Legal_name:ul.Legal_name, Place_birth:ul.Place_birth, Salary:ul.Salary, Foto_ktp:ul.Foto_ktp ,Foto_selfie:ul.Foto_selfie},Date_birth:ul.Date_birth}
	// u.Konsumen = {Nik : Nik, Full_name:Full_name, Legal_name:Legal_name, Place_birth, Salary, Foto_ktp ,Foto_selfie}
	// u.Date_birth = Date_birth
	valid.Required(ul.Nik, "Nik")
	valid.Required(ul.Full_name, "Full_name")
	valid.Required(ul.Legal_name, "Legal_name")
	valid.Required(ul.Place_birth, "Place_birth")
	valid.Required(ul.Date_birth, "Date_birth")
	valid.Required(ul.Salary, "Salary")
	valid.Required(ul.Foto_ktp, "Foto_ktp")
	valid.Required(ul.Foto_selfie, "Foto_selfie")

	if valid.HasErrors() {
		// If there are error messages it means the validation didn't pass
		// Print error message
		for _, err := range valid.Errors {
			return err.Key + err.Message
		}
	}

	return ""
}

func (api *KonsumensController) CreateKonsumens() {
	frm := api.Ctx.Input.RequestBody
	if AllKonsumensCheck(api) != "" {
		api.Data["json"] = AllKonsumensCheck(api)
		api.ServeJSON()
		return
	}

	o := orm.NewOrm()
	o.Using("default")

	ul := &structs.GetKonsumen{}
	json.Unmarshal(frm, ul)
	dateBirth,_ := time.Parse("2006-01-02", ul.Date_birth)
	Qry := models.Konsumens{Konsumen : models.Konsumen{Nik: ul.Nik, Full_name: ul.Full_name, Legal_name: ul.Legal_name, Place_birth: ul.Place_birth, Salary: ul.Salary, Foto_ktp: ul.Foto_ktp, Foto_selfie: ul.Foto_selfie},Date_birth: dateBirth}
	o.Insert(&Qry)
	api.Data["json"] = "Successfully save data"
	api.ServeJSON()
}

func (api *KonsumensController) UpdateKonsumens() {
	frm := api.Ctx.Input.RequestBody
	if AllKonsumensCheck(api) != "" {
		api.Data["json"] = AllKonsumensCheck(api)
		api.ServeJSON()
		return
	}

	o := orm.NewOrm()
	o.Using("default")

	idInt, _ := strconv.Atoi(api.Ctx.Input.Param(":id"))

	ul := &structs.GetKonsumen{}
	json.Unmarshal(frm, ul)
	dateBirth,_ := time.Parse("2006-01-02", ul.Date_birth)
	Qry := models.Konsumens{Id:idInt,Konsumen : models.Konsumen{Nik: ul.Nik, Full_name: ul.Full_name, Legal_name: ul.Legal_name, Place_birth: ul.Place_birth, Salary: ul.Salary, Foto_ktp: ul.Foto_ktp, Foto_selfie: ul.Foto_selfie},Date_birth: dateBirth}
	o.Update(&Qry)
	
	api.Data["json"] = "Successfully update data with id  =  " + api.Ctx.Input.Param(":id")
	
	api.ServeJSON()
}