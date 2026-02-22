package config

type AppConfig struct {
	ServiceName        string
	ServicePort        string
	JWTSecret          string
	JWTExpireHour      int
	CookieDomain       string
	AllowedOrigin      string
	InfraAPIBaseDomain string
	InfraGetAllPath    string
	InfraInsertPath    string
}

func LoadConfig() *AppConfig {
	return &AppConfig{
		ServiceName:        "Cloud Management Service",
		ServicePort:        "8080",
		JWTSecret:          "golf_jwt_krub",
		JWTExpireHour:      4,
		CookieDomain:       "localhost",
		AllowedOrigin:      "http://localhost:3000", // responding to task 2 req.
		InfraAPIBaseDomain: "http://localhost:8081", // responding to task 2 req.
		InfraGetAllPath:    "/v1/skus",              // responding to task 2 req.
		InfraInsertPath:    "/v1/resources",         // responding to task 2 req.
	}
}
