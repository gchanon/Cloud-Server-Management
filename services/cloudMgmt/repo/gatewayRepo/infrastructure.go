package gatewayrepo

type ResponseGWGetAllInfra struct {
	SkuData []SkuBodyFromGWGetAllInfra `json:"skus"`
}

type SkuBodyFromGWGetAllInfra struct {
	ID           string  `json:"id"`
	Sku          string  `json:"sku"`
	Type         string  `json:"type"`
	Name         string  `json:"name"`
	CPU          int     `json:"cpu"`
	RAM          int     `json:"ram"`
	Disk         int     `json:"disk"`
	PriceHourly  float64 `json:"price_hourly"`
	PriceMonthly float64 `json:"price_monthly"`
}

type RequestAddInfra struct {
	Sku string `json:"sku"`
}

type ResponseAddInfra struct {
	Success bool   `json:"success"`
	Id      string `json:"id"`
}

type ResponseGetRegistedInfra struct {
	Resources []RegistedInfra `json:"resources"`
}

type RegistedInfra struct {
	ID        string `json:"id"`
	SKU       string `json:"sku"`
	Status    string `json:"status"`
	IP        string `json:"ip"`
	CreatedAt string `json:"created_at"`
}

type RequestControlInfraPower struct {
	Action string `json:"action"`
}

type ResponseControlInfraPower struct {
	Status string `json:"status"`
	State  string `json:"state"`
}
