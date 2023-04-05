package galactus_client

type AuthInfo struct {
}

// TODO: requirements:
// 		Login: username + password + multiaddr -> token + all relevant albums + all relevant multiaddrs
// 		Sync: token + multiaddr -> all relevant albums + all relevant multiaddrs
//		Will there ever be a case where we add an album WITHOUT updating Galactus?
//		Can we add users to an album after creation?

// Endpoint used to log authenticate user through Galactus
// The same endpoint can be used for new and returning users
// If a new user uses this endpoint to sign up, a new document will be created in the User collection in MongoDB
// TODO: can this fetch the albums the user has access to, like gc.Sync?
func (gc *GalactusClient) Login() *AuthInfo {
	return nil
}

// Fetches all the albums the user has access to and the multiaddrs of all the peers who also have access to each album
func (gc *GalactusClient) Sync() *AuthInfo {
	return nil
}

// Creates an album and saves it to MongoDB. Must also send a list of authorized users with payload
func (gc *GalactusClient) AddAlbum() *AuthInfo {
	return nil
}

// Updates the authenticated users array of an album
// Token required and must be creator of album to edit
func (gc *GalactusClient) UpdateAlbum() *AuthInfo {
	return nil
}

// Retrieves list of all users (username)
// Token required
func (gc *GalactusClient) GetUsers() *AuthInfo {
	return nil
}

// Delete an album
// Token required
func (gc *GalactusClient) DeleteAlbum() *AuthInfo {
	return nil
}
