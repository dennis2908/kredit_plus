package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"kredit_plus/models"
	"kredit_plus/structs"
	"net/http"

	elastic "github.com/olivere/elastic/v7"
)

func SearchData(keyword string) ([]models.Env, error) {
	esclient, _ := Connect()
	searchSource := elastic.NewSearchSource()
	searchSource.Query(elastic.NewMatchQuery("Name", keyword).Operator("and"))

	/* this block will basically print out the es query */
	queryStr, err1 := searchSource.Source()
	_, err2 := json.Marshal(queryStr)

	if err1 != nil {
		return []models.Env{}, err1
	}
	if err2 != nil {
		return []models.Env{}, err2
	}

	searchService := esclient.Search().Index("env").SearchSource(searchSource)

	searchResult, err := searchService.Do(context.Background())
	if err != nil {
		fmt.Println("Error", err)
		return []models.Env{}, err
	}

	var datas []models.Env

	for _, hit := range searchResult.Hits.Hits {
		var data models.Env
		err := json.Unmarshal(hit.Source, &data)
		if err != nil {
			fmt.Println("[Getting Datas][Unmarshal] Err=", err)
		}

		datas = append(datas, data)
	}

	return datas, nil
}

func SearchDataQuery(name string, app string) (structs.SearchHits, error) {
	var t structs.SearchHits
	query := fmt.Sprintf(`
	{
		"query": {
			"has_parent": {
			"parent_type": "general",
			"query": {
				"bool": {
				"must": [
					{"match_phrase": {"name": "%s"}},
					{"match_phrase": {"app": "%s"}}
				]
				}
			},
			"inner_hits": {}    
			}
		}
	}
	`, name, app)

	req, err := http.Post("http://localhost:9200/multiple-env/_search", "application/json",
		bytes.NewBufferString(query))
	if err != nil {
		return t, err
	}

	responseBody, err := io.ReadAll(req.Body)

	if err != nil {
		return t, err
	}

	json.Unmarshal(responseBody, &t)

	if len(t.Hits.Hits) == 0 {
		return t, fmt.Errorf("Error")
	}

	return t, nil

}

func SearchDataById(keyword string) ([]models.Env, error) {
	esclient, _ := Connect()
	searchSource := elastic.NewSearchSource()
	searchSource.Query(elastic.NewMatchQuery("Id", keyword).Operator("and"))

	/* this block will basically print out the es query */
	queryStr, err1 := searchSource.Source()
	_, err2 := json.Marshal(queryStr)

	if err1 != nil {
		return []models.Env{}, err1
	}
	if err2 != nil {
		return []models.Env{}, err2
	}

	searchService := esclient.Search().Index("env").SearchSource(searchSource)

	searchResult, err := searchService.Do(context.Background())
	if err != nil {
		fmt.Println("Error", err)
		return []models.Env{}, err
	}

	var datas []models.Env

	for _, hit := range searchResult.Hits.Hits {
		var data models.Env
		err := json.Unmarshal(hit.Source, &data)
		if err != nil {
			fmt.Println("[Getting Datas][Unmarshal] Err=", err)
		}

		datas = append(datas, data)
	}

	return datas, nil
}

func SearchAllData() ([]models.Env, error) {
	esclient, _ := Connect()

	/* this block will basically print out the es query */

	searchService := esclient.Search().Index("env")

	searchResult, err := searchService.Do(context.Background())
	if err != nil {
		fmt.Println("Error", err)
		return []models.Env{}, err
	}

	var datas []models.Env

	for _, hit := range searchResult.Hits.Hits {
		var data models.Env
		err := json.Unmarshal(hit.Source, &data)
		if err != nil {
			fmt.Println("[Getting Datas][Unmarshal] Err=", err)
		}

		datas = append(datas, data)
	}

	return datas, nil
}
