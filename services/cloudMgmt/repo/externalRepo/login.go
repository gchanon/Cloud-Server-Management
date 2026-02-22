package externalrepo

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	LoginStatus string `json:"login_status"`
}
