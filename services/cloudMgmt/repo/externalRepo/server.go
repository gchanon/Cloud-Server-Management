package externalrepo

type ResponseGetAllServer struct {
	IsFound    bool                    `json:"is_found"`
	ServerList []SkuByServerIdResponse `json:"server_list"`
}

type SkuByServerIdResponse struct {
	ServerID  string `json:"server_id"`
	SKU       string `json:"sku"`
	IsPowerOn bool   `json:"is_power_on"`
}

type ResponseAddServer struct {
	Success bool `json:"success"`
}

type ResponseControlPower struct {
	Success bool   `json:"success"`
	State   string `json:"state"`
}
