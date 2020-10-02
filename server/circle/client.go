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

// GetRecentlyBuiltPipelines get all recently built pipelines in a org
func GetRecentlyBuiltPipelines(apiToken string, orgSlug string, mine bool) ([]circleci.Pipeline1, error) {
	pl, _, err := client.PipelineApi.ListPipelines(getContext(apiToken), orgSlug, mine, nil)
	if err != nil {
		return nil, err
	}

	return pl.Items, nil
}

// GetAllPipelinesForProject get all pipelines for a given project
func GetAllPipelinesForProject(apiToken string, projectSlug string) ([]circleci.Pipeline1, error) {
	pl, _, err := client.PipelineApi.ListPipelinesForProject(getContext(apiToken), projectSlug, nil)
	if err != nil {
		return nil, err
	}

	return pl.Items, nil
}

// GetAllMyPipelinesForProject get all pipelines triggered by you
func GetAllMyPipelinesForProject(apiToken string, projectSlug string) ([]circleci.Pipeline1, error) {
	pl, _, err := client.PipelineApi.ListMyPipelines(getContext(apiToken), projectSlug, nil)
	if err != nil {
		return nil, err
	}

	return pl.Items, nil
}

// GetWorkflowsByPipeline get all workflows by pipeline ID
func GetWorkflowsByPipeline(apiToken string, pipelineID string) ([]circleci.Workflow1, error) {
	wf, _, err := client.PipelineApi.ListWorkflowsByPipelineId(getContext(apiToken), pipelineID, nil)
	if err != nil {
		return nil, err
	}

	return wf.Items, nil
}

// GetNameByID returns username from user id
func GetNameByID(apiToken string, id string) (string, error) {
	user, _, err := client.UserApi.GetUser(getContext(apiToken), id)
	if err != nil {
		return "", err
	}

	return user.Name, nil
}
