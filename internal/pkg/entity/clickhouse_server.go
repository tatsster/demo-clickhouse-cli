package entity

type ClickHouseServer struct {
	Id       string   `json:"id"`
	OrgId    string   `json:"org_id"`
	Host     string   `json:"host"`
	Port     string   `json:"port"`
	Cluster  string   `json:"cluster"`
	Username string   `json:"username"`
	Password string   `json:"-"`
	Shards   []string `json:"shards"`
}

type ClickHouseShard struct {
	Name    string `json:"name"`
	Index   int    `json:"index"`
	Replica string `json:"replica"`
}
