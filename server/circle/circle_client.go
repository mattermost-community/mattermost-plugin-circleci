package circle

import (
	"context"

	"github.com/darkLord19/circleci-v2/circleci"
)

var (
	client *circleci.APIClient
)

func init() {
	config := circleci.NewConfiguration()
	client = circleci.NewAPIClient(config)
}

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
	user, _, err := client.UserApi.GetCurrentUser(getContext(apiToken))
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetWorkflow returns the info for given workflow id
func GetWorkflow(apiToken string, workflowID string) (*circleci.Workflow, error) {
	wf, _, err := client.WorkflowApi.GetWorkflowById(getContext(apiToken), workflowID)
	if err != nil {
		return nil, err
	}

	return &wf, nil
}

// GetWorkflowJobs returns the info of jobs for given workflow id
func GetWorkflowJobs(apiToken string, workflowID string) (*[]circleci.Job, error) {
	wf, _, err := client.WorkflowApi.ListWorkflowJobs(getContext(apiToken), workflowID)
	if err != nil {
		return nil, err
	}

	return &wf.Items, nil
}
