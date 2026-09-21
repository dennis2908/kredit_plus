package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"

	_ "github.com/lib/pq"

	_ "fmt"
	_ "kredit_plus/routers"

	_ "github.com/lib/pq"

	elastic "github.com/olivere/elastic/v7"
)

type ElasticServer int64

type Args struct{}

func main() {
	elasticServer := new(ElasticServer)
	rpc.Register(elasticServer)

	rpc.HandleHTTP()
	// Start listening for the requests on port 1234
	listener, err := net.Listen("tcp", "0.0.0.0:1234")
	if err != nil {
		log.Fatal("Listener error: ", err)
	}
	http.Serve(listener, nil)
}

func Connect() (*elastic.Client, error) {

	client, err := elastic.NewClient(elastic.SetURL("http://127.0.0.1:9200"),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false))

	if err != nil {
		fmt.Println("Some error", err.Error())
		panic("Failed to initialize elastic-search client")
	}

	fmt.Println("ES initialized...")

	return client, err

}
