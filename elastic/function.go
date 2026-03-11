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

func InsertDataQueryGeneral(id string, name string, app string, framework string) bool {
	query := fmt.Sprintf(`
	{
		"id": "%s",
		"name": "%s",
		"app": "%s",
		"framework" : "%s",
		"join": {
			"name": "general" 
		}
	}

	`, id, name, app, framework)

	url := "http://127.0.0.1:9200/multiple-env/_doc/" + id + "?refresh"

	req, err := http.NewRequest(http.MethodPut, url,
		bytes.NewBuffer([]byte(query)))
	if err != nil {
		return false
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, errx := client.Do(req)

	if errx != nil {
		return false
	}

	defer resp.Body.Close()

	return true

}

func InsertDataQueryDesc(id1 string, id2 string, env string, lang string) bool {
	query := fmt.Sprintf(`
	{
		"id": "%s",
		"env" : "%s",
		"lang" : "%s",
		"join": {
			"name": "description",
			"parent": "%s"
		}
	}


	`, id2, env, lang, id1)

	url := "http://localhost:9200/multiple-env/_doc/" + id2 + "?routing=" + id1 + "&refresh"

	fmt.Println("1", url)

	req, err := http.Post(url, "application/json",
		bytes.NewBufferString(query))
	if err != nil {
		return false
	}

	_, errx := io.ReadAll(req.Body)

	return errx == nil

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
