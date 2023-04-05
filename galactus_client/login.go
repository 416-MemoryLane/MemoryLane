package galactus_client

type LoginResponse struct {
	Username string `json:"username"`
	Token    string `json:"token"`
	Message  string `json:"message"`
}
