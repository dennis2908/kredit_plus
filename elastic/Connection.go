package elastic

import (
	"fmt"

	elastic "github.com/olivere/elastic/v7"
)

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
