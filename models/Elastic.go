package models

type Env struct {
	Id       string
	Name     string
	Env_name string
}

type InsertEnv struct {
	Id        string
	Name      string
	Env       string
	App       string
	Framework string
	Lang      string
}

type MultipleEnv struct {
	Id   string
	Env  string
	Lang string
}
