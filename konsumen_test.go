package main_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/astaxie/beego"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("test get konsumen", func() {
	Describe("GET /konsumen", func() {
		It("response has http code 200", func() {
			request, _ := http.NewRequest("GET", "/konsumen", nil)
			response := httptest.NewRecorder()

			beego.BeeApp.Handlers.ServeHTTP(response, request)

			Expect(response.Code).To(Equal(http.StatusOK))
		})
	})
})

var _ = Describe("test get konsumen by id", func() {
	Describe("GET /konsumen/1", func() {
		It("response has http code 200", func() {
			request, _ := http.NewRequest("GET", "/konsumen/1", nil)
			response := httptest.NewRecorder()

			beego.BeeApp.Handlers.ServeHTTP(response, request)

			Expect(response.Code).To(Equal(http.StatusOK))
		})
	})
})

var _ = Describe("test post konsumen by id", func() {
	Describe("POST /konsumen", func() {
		Context("when Full Name is empty", func() {
			It("informs about error", func() {

				url := "http://localhost:9333/konsumen"

				var jsonStr = []byte(`{"nik":"1919919119"}`)
				req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
				req.Header.Set("X-Custom-Header", "myvalue")
				req.Header.Set("Content-Type", "application/json")

				response := httptest.NewRecorder()

				beego.BeeApp.Handlers.ServeHTTP(response, req)

				Expect(response.Code).To(Equal(http.StatusOK))
				beego.Debug(response.Body.String())
				Expect(response.Body.String()).To(ContainSubstring("Full_name"))

			})
		})
	})
})

var _ = Describe("test put /konsumen", func() {
	Describe("PUT /konsumen", func() {
		Context("when Full Name is empty", func() {
			It("informs about error", func() {

				url := "http://localhost:9333/konsumen/1"

				var jsonStr = []byte(`{"nik":"1919919119"}`)
				req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
				req.Header.Set("X-Custom-Header", "myvalue")
				req.Header.Set("Content-Type", "application/json")

				response := httptest.NewRecorder()

				beego.BeeApp.Handlers.ServeHTTP(response, req)

				Expect(response.Code).To(Equal(http.StatusOK))
				beego.Debug(response.Body.String())
				Expect(response.Body.String()).To(ContainSubstring("Full_name"))

			})
		})
	})
})
