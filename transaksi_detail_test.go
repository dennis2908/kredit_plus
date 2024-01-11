package main_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/astaxie/beego"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("test get transaksidetails", func() {
	Describe("GET /transaksidetails", func() {
		It("response has http code 200", func() {
			request, _ := http.NewRequest("GET", "/transaksidetails", nil)
			response := httptest.NewRecorder()

			beego.BeeApp.Handlers.ServeHTTP(response, request)

			Expect(response.Code).To(Equal(http.StatusOK))
		})
	})
})

var _ = Describe("test get transaksidetails by id", func() {
	Describe("GET /transaksidetails/1", func() {
		It("response has http code 200", func() {
			request, _ := http.NewRequest("GET", "/transaksidetails/1", nil)
			response := httptest.NewRecorder()

			beego.BeeApp.Handlers.ServeHTTP(response, request)

			Expect(response.Code).To(Equal(http.StatusOK))
		})
	})
})

var _ = Describe("test get transaksidetails by id konsumen", func() {
	Describe("GET /transaksidetails/1", func() {
		It("response has http code 200", func() {
			request, _ := http.NewRequest("GET", "/transaksidetails/1", nil)
			response := httptest.NewRecorder()

			beego.BeeApp.Handlers.ServeHTTP(response, request)

			Expect(response.Code).To(Equal(http.StatusOK))
		})
	})
})

var _ = Describe("test post transaksidetails by id", func() {
	Describe("POST /transaksidetails", func() {
		Context("when Full Name is empty", func() {
			It("informs about error", func() {

				url := "http://localhost:9333/transaksidetails"

				var jsonStr = []byte(`{"Id_Transaksi":1}`)
				req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
				req.Header.Set("X-Custom-Header", "myvalue")
				req.Header.Set("Content-Type", "application/json")

				response := httptest.NewRecorder()

				beego.BeeApp.Handlers.ServeHTTP(response, req)

				Expect(response.Code).To(Equal(http.StatusOK))
				beego.Debug(response.Body.String())
				Expect(response.Body.String()).To(ContainSubstring("Bulan"))

			})
		})
	})
})

var _ = Describe("test put /transaksidetails", func() {
	Describe("PUT /transaksidetails", func() {
		Context("when Full Name is empty", func() {
			It("informs about error", func() {

				url := "http://localhost:9333/transaksidetails/1"

				var jsonStr = []byte(`{"Id_Transaksi":1}`)
				req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
				req.Header.Set("X-Custom-Header", "myvalue")
				req.Header.Set("Content-Type", "application/json")

				response := httptest.NewRecorder()

				beego.BeeApp.Handlers.ServeHTTP(response, req)

				Expect(response.Code).To(Equal(http.StatusOK))
				beego.Debug(response.Body.String())
				Expect(response.Body.String()).To(ContainSubstring("Bulan"))

			})
		})
	})
})
