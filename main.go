package main

import (
	_ "fmt"
	Connection "kredit_plus/connects"

	Elastic "kredit_plus/elastic"
	_ "kredit_plus/routers"
	"kredit_plus/ssrf"
	"kredit_plus/token"
	"os"
	"runtime"
	"sync"

	"log"
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beego/beego/v2/client/cache"
	"github.com/beego/beego/v2/core/logs"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

var (
	codeCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "http_request_total_code",
			Help: "total request code controller",
		},
	)
)

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of get requests.",
	},
	[]string{"path"},
)

var responseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "response_status",
		Help: "Status of HTTP response",
	},
	[]string{"status"},
)

var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_response_time_seconds",
	Help: "Duration of HTTP requests.",
}, []string{"path"})

func init() { // init instead of int

	Connection.Connects()
	Elastic.Connect()
	prometheus.Register(totalRequests)
	prometheus.Register(responseStatus)
	prometheus.Register(httpDuration)
	prometheus.MustRegister(codeCounter)
}
func main() {
	err := ssrf.Main()
	if err != nil {
		log.Println(err)
		return
	}

	numberOfCores := runtime.NumCPU()
	runtime.GOMAXPROCS(numberOfCores)
	var wg sync.WaitGroup
	for i := 0; i < numberOfCores; i++ {
		wg.Add(1)
		_, err := cache.NewCache("file", `{"CachePath":"./cache","FileSuffix":".cache", "EmbedExpiry": "120"}`)

		beego.InsertFilter("/konsumen/*", beego.BeforeExec, token.ValidateToken)
		beego.InsertFilter("/transaksi/*", beego.BeforeExec, token.ValidateToken)
		beego.InsertFilter("/transaksidetails/*", beego.BeforeExec, token.ValidateToken)
		beego.InsertFilter("/konsumen_mongo_insert/*", beego.BeforeExec, token.ValidateToken)
		beego.InsertFilter("/konsumen_mongo_update/*", beego.BeforeExec, token.ValidateToken)
		beego.InsertFilter("/konsumen_excel_read/*", beego.BeforeExec, token.ValidateToken)

		// Prometheus endpoint
		beego.Handler("/metrics", promhttp.Handler())

		orm.Debug = true

		o := orm.NewOrm()
		o.Using("default")

		if err != nil {
			logs.Error(err)
		}
		log.Println("Env $PORT :", os.Getenv("PORT"))
		if os.Getenv("PORT") != "" {
			port, err := strconv.Atoi(os.Getenv("PORT"))
			if err != nil {
				log.Fatal(err)
				log.Fatal("$PORT must be set")
			}
			log.Println("port : ", port)
			beego.BConfig.Listen.HTTPPort = port
			beego.BConfig.Listen.HTTPSPort = port
		}

		beego.Run()
	}
	wg.Wait()

}
