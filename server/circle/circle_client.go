package circle

import (
	"context"

	"github.com/TomTucka/go-circleci/circleci"
)

func getClient() *circleci.APIClient {
	config := circleci.NewConfiguration()
	return circleci.NewAPIClient(config)
}

func getContext(apiToken string) context.Context {
	apiKey := &circleci.APIKey{Key: apiToken}
	parentContext := context.TODO()
	return context.WithValue(parentContext, circleci.ContextAPIKey, *apiKey)
}

// GetCurrentUser returns the current user
func GetCurrentUser(apiToken string) (*circleci.User, error) {
	client := getClient()
	user, _, err := client.UserApi.GetCurrentUser(getContext(apiToken))
	if err != nil {
		return nil, err
	}

	return &user, nil
}
