package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	_ "fmt"
	"kredit_plus/elastic"
	models "kredit_plus/models"
	"time"

	"strconv"

	"github.com/astaxie/beego"
	_ "github.com/beego/beego/v2/core/logs"
	_ "github.com/leekchan/accounting"
	_ "github.com/shopspring/decimal"
)

type ElasticController struct {
	beego.Controller
}

func (api *ElasticController) ElasticInsert() {
	frm := api.Ctx.Input.RequestBody
	ul := &models.Env{}
	json.Unmarshal(frm, ul)

	data, err := elastic.SearchData(ul.Name)

	if err != nil {

		api.Ctx.ResponseWriter.WriteHeader(400)
		api.Ctx.ResponseWriter.Write([]byte("error fetch data"))
		return
	}

	fmt.Println(111, len(data))

	if len(data) > 0 {
		api.Ctx.ResponseWriter.WriteHeader(503)
		api.Ctx.ResponseWriter.Write([]byte("data exists"))
		return
	}
	id := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	Qry := models.Env{Name: ul.Name, Env_name: ul.Env_name, Id: id}
	dataJSON, _ := json.Marshal(Qry)
	js := string(dataJSON)
	ctx := context.Background()
	esclient, err := elastic.Connect()
	if err != nil {
		fmt.Println("Error initializing : ", err)
		panic("Client fail ")
	}

	_, err_ind := esclient.Index().Id(id).Index("env").
		BodyJson(js).
		Do(ctx)

	if err_ind != nil {
		panic(err_ind)
	}

	api.Data["json"] = "Successfully save data"
	api.ServeJSON()
}

func (api *ElasticController) InsertData() {
	frm := api.Ctx.Input.RequestBody
	ul := &models.InsertEnv{}
	json.Unmarshal(frm, ul)

	_, err := elastic.SearchDataQuery(ul.Name, ul.App)

	if err == nil {

		api.Ctx.ResponseWriter.WriteHeader(400)
		api.Ctx.ResponseWriter.Write([]byte("data exists"))
		return
	}

	id := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

	if elastic.InsertDataQueryGeneral(id, ul.Name, ul.App, ul.Framework) {
		id2 := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
		if elastic.InsertDataQueryDesc(id, id2, ul.Env, ul.Lang) {
			api.Data["json"] = "Successfully save data"
			api.ServeJSON()
			return
		}
	}

	api.Ctx.ResponseWriter.WriteHeader(503)
	api.Ctx.ResponseWriter.Write([]byte("data not saved"))

}

func (api *ElasticController) ElasticSearch() {

	keyword := api.GetString("search")

	data, err := elastic.SearchData(keyword)

	if err != nil {

		api.Ctx.ResponseWriter.WriteHeader(400)
		api.Ctx.ResponseWriter.Write([]byte("error search data"))
		return
	}

	if len(data) == 0 {
		api.Ctx.ResponseWriter.WriteHeader(503)
		api.Ctx.ResponseWriter.Write([]byte("data not found"))
		return
	}

	api.Ctx.Output.SetStatus(200)
	api.Data["json"] = data
	api.ServeJSON()

	// if err != nil {
	// 	fmt.Println("Fetching datas fail: ", err)
	// } else {
	// 	for _, s := range datas {
	// 		fmt.Printf("Datas found Name: %s, Env Name: %d \n", s.Name, s.Env_name)
	// 	}
	// }

}

func (api *ElasticController) SearchData() {

	t, err := elastic.SearchDataQuery(api.GetString("name"), api.GetString("app"))

	if err != nil {
		api.Ctx.ResponseWriter.WriteHeader(503)
		api.Ctx.ResponseWriter.Write([]byte("data not found"))
		return
	}

	api.Data["json"] = t.Hits.Hits[0].Source

	api.ServeJSON()

}

func (api *ElasticController) ElasticDelete() {

	id := api.Ctx.Input.Param(":id")

	fmt.Println(12212, id)

	esclient, _ := elastic.Connect()

	_, err := esclient.Delete().Index("env").Id(id).Refresh("true").Do(context.Background())

	if err != nil {

		api.Ctx.ResponseWriter.WriteHeader(400)
		api.Ctx.ResponseWriter.Write([]byte("error delete data"))
		return
	}

	api.Ctx.Output.SetStatus(200)
	api.Data["json"] = "data had been delete"
	api.ServeJSON()

	// if err != nil {
	// 	fmt.Println("Fetching datas fail: ", err)
	// } else {
	// 	for _, s := range datas {
	// 		fmt.Printf("Datas found Name: %s, Env Name: %d \n", s.Name, s.Env_name)
	// 	}
	// }

}

func (api *ElasticController) ElasticGetAllData() {

	data, err := elastic.SearchAllData()

	if err != nil {

		api.Ctx.ResponseWriter.WriteHeader(400)
		api.Ctx.ResponseWriter.Write([]byte("error search data"))
		return
	}

	if len(data) == 0 {
		api.Ctx.ResponseWriter.WriteHeader(503)
		api.Ctx.ResponseWriter.Write([]byte("data not found"))
		return
	}

	api.Ctx.Output.SetStatus(200)
	api.Data["json"] = data
	api.ServeJSON()

	// if err != nil {
	// 	fmt.Println("Fetching datas fail: ", err)
	// } else {
	// 	for _, s := range datas {
	// 		fmt.Printf("Datas found Name: %s, Env Name: %d \n", s.Name, s.Env_name)
	// 	}
	// }

}

func (api *ElasticController) ElasticGetDataById() {
	id := api.Ctx.Input.Param(":id")
	data, err := elastic.SearchDataById(id)

	if err != nil {

		api.Ctx.ResponseWriter.WriteHeader(400)
		api.Ctx.ResponseWriter.Write([]byte("error search data"))
		return
	}

	if len(data) == 0 {
		api.Ctx.ResponseWriter.WriteHeader(503)
		api.Ctx.ResponseWriter.Write([]byte("data not found"))
		return
	}

	api.Ctx.Output.SetStatus(200)
	api.Data["json"] = data
	api.ServeJSON()

}
