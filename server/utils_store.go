package main

const (
	KVStoreSuffix = "_circleci_token"
)

// Return false if no token is saved for this user
func (p *Plugin) getTokenFromKVStore(userID string) (string, bool) {
	raw, appErr := p.API.KVGet(userID + KVStoreSuffix)
	if appErr != nil {
		p.API.LogError("Unable to reach KVStore", "KVStore error", appErr)
		return "", false
	}

	if raw == nil {
		return "", false
	}

	userToken := string(raw)
	return userToken, true
}

// Return false if the token has not been saved
func (p *Plugin) storeTokenInKVStore(userID string, circleciToken string) bool {
	appErr := p.API.KVSet(userID+KVStoreSuffix, []byte(circleciToken))
	if appErr != nil {
		p.API.LogError("Unable to write in KVStore", "KVStore error", appErr)
		return false
	}

	return true
}

// Return false if the token has not been deleted
func (p *Plugin) deleteTokenFromKVStore(userID string) bool {
	if appErr := p.API.KVDelete(userID + KVStoreSuffix); appErr != nil {
		p.API.LogError("Unable to delete from KVStore", "KVStore error", appErr)
		return false
	}

	return true
}
