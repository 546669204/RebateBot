package common

type Msg struct {
	Method string `json:"method"`
	Data   string `json:"data"`
	To     string `json:"to"`
	ID     string `json:"id"`
}
type Resp struct {
	Data   string `json:"data"`
	Status int    `json:"status"`
	ID     string `json:"id"`
}

type ConfigModel struct {
	Services []ServicesModel   `json:"services"`
	Routing  map[string]string `json:"routing"`
}

type ServicesModel struct {
	Name string `json:"name"`
	Run  int    `json:"run"`
}
