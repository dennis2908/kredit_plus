package structs

type Env struct {
	Name     string `json:"name"`
	Env_name string `json:"env_name"`
}

type SearchDataStr struct {
	Name string
	App  string
}

type MultipleEnv struct {
	Id   string `json:"id"`
	Env  string `json:"env"`
	Lang string `json:"lang"`
}

type SearchHits struct {
	Hits struct {
		Hits []*struct {
			Source *MultipleEnv `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}
