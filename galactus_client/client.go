package galactus_client

import "log"

type GalactusClient struct {
	Username  string
	Password  string
	Multiaddr string
	AuthToken string
	l         *log.Logger
}

func NewGalactusClient(un, pw, ma string, l *log.Logger) *GalactusClient {
	return &GalactusClient{un, pw, ma, "", l}
}

// TODO: requirements:
// 		Login: username + password + multiaddr -> token
// 		Sync: token + multiaddr -> all relevant albums + all relevant multiaddrs

// Endpoint used to log authenticate user through Galactus
// The same endpoint can be used for new and returning users
// If a new user uses this endpoint to sign up, a new document will be created in the User collection in MongoDB
func (gc *GalactusClient) Login() *LoginResponse {
	return nil
}

// Fetches all the albums the user has access to and the multiaddrs of all the peers who also have access to each album
func (gc *GalactusClient) Sync() *SyncResponse {
	return nil
}
