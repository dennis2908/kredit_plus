package controllers

import (
	_ "fmt"
	"kredit_plus/structs"
	"sync"
)

var (
	wg sync.WaitGroup
)

type RetAPIController struct{}

type DepRetAPI interface {
	RetError(a *ElasticController, header int, str string)
	RetAPIElasticGet(a *ElasticController, reply structs.SearchHits)
	RetAPIElasticGetAll(a *ElasticController, reply structs.SearchHits)
}

type RetAPI struct {
	depRetAPI DepRetAPI
}

func (api RetAPIController) RetAPIElasticGetAll(a *ElasticController, reply structs.SearchHits) {

	a.Data["json"] = reply.Hits.Hits
	wg.Add(1)
	go a.Ctx.Output.SetStatus(200)
	wg.Done()
	go a.ServeJSON()
	wg.Wait()

}

func (api RetAPIController) RetAPIElasticGet(a *ElasticController, reply structs.SearchHits) {

	a.Data["json"] = reply.Hits.Hits[0].Source
	wg.Add(1)
	go a.Ctx.Output.SetStatus(200)
	wg.Done()
	go a.ServeJSON()

}

func (api RetAPIController) RetError(a *ElasticController, header int, str string) {

	go a.Ctx.ResponseWriter.WriteHeader(header)
	go a.Ctx.ResponseWriter.Write([]byte(str))

}
