package main_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/astaxie/beego"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("test get transaksi", func() {
	Describe("GET /transaksi", func() {
		It("response has http code 200", func() {
			request, _ := http.NewRequest("GET", "/transaksi", nil)
			response := httptest.NewRecorder()

			beego.BeeApp.Handlers.ServeHTTP(response, request)

			Expect(response.Code).To(Equal(http.StatusOK))
		})
	})
})

var _ = Describe("test get transaksi by id", func() {
	Describe("GET /transaksi/1", func() {
		It("response has http code 200", func() {
			request, _ := http.NewRequest("GET", "/transaksi/1", nil)
			response := httptest.NewRecorder()

			beego.BeeApp.Handlers.ServeHTTP(response, request)

			Expect(response.Code).To(Equal(http.StatusOK))
		})
	})
})

var _ = Describe("test get transaksi by id konsumen", func() {
	Describe("GET /transaksi/konsumen/1", func() {
		It("response has http code 200", func() {
			request, _ := http.NewRequest("GET", "/transaksi/konsumen/1", nil)
			response := httptest.NewRecorder()

			beego.BeeApp.Handlers.ServeHTTP(response, request)

			Expect(response.Code).To(Equal(http.StatusOK))
		})
	})
})

var _ = Describe("test post transaksi by id", func() {
	Describe("POST /transaksi", func() {
		Context("when Full Name is empty", func() {
			It("informs about error", func() {

				url := "http://localhost:9333/transaksi"

				var jsonStr = []byte(`{"Id_konsumen":1}`)
				req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
				req.Header.Set("X-Custom-Header", "myvalue")
				req.Header.Set("Content-Type", "application/json")

				response := httptest.NewRecorder()

				beego.BeeApp.Handlers.ServeHTTP(response, req)

				Expect(response.Code).To(Equal(http.StatusOK))
				beego.Debug(response.Body.String())
				Expect(response.Body.String()).To(ContainSubstring("No_kontrak"))

			})
		})
	})
})

var _ = Describe("test put /transaksi", func() {
	Describe("PUT /transaksi", func() {
		Context("when Full Name is empty", func() {
			It("informs about error", func() {

				url := "http://localhost:9333/transaksi/1"

				var jsonStr = []byte(`{"Id_konsumen":1}`)
				req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
				req.Header.Set("X-Custom-Header", "myvalue")
				req.Header.Set("Content-Type", "application/json")

				response := httptest.NewRecorder()

				beego.BeeApp.Handlers.ServeHTTP(response, req)

				Expect(response.Code).To(Equal(http.StatusOK))
				beego.Debug(response.Body.String())
				Expect(response.Body.String()).To(ContainSubstring("No_kontrak"))

			})
		})
	})
})
