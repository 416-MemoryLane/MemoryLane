package galactus_client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type GalactusClient struct {
	API       string
	Username  string
	Password  string
	Multiaddr string
	AuthToken string
	l         *log.Logger
}

func NewGalactusClient(api, un, pw, ma string, l *log.Logger) *GalactusClient {
	return &GalactusClient{api, un, pw, ma, "", l}
}

// TODO: requirements:
// 		Sync: token + multiaddr -> all relevant albums + all relevant multiaddrs

// Endpoint used to log authenticate user through Galactus
// The same endpoint can be used for new and returning users
// If a new user uses this endpoint to sign up, a new document will be created in the User collection in MongoDB
func (gc *GalactusClient) Login() (*LoginResponse, error) {
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`{
    "username": "%s",
    "password": "%s",
    "multiaddr": "%s"
	}`, gc.Username, gc.Password, gc.Multiaddr))

	client := &http.Client{}
	url := fmt.Sprintf("%s/login", gc.API)
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, fmt.Errorf("error creating Galactus request to login: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing HTTP request to Galactus while logging in: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading HTTP response from Galactus while logging in: %w", err)
	}

	// Unmarshal JSON response into struct
	var loginResp LoginResponse
	err = json.Unmarshal(body, &loginResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response from Galactus while logging in: %w", err)
	}

	return &loginResp, nil
}

// Fetches all the albums the user has access to and the multiaddrs of all the peers who also have access to each album
func (gc *GalactusClient) Sync() *SyncResponse {
	return nil
}
