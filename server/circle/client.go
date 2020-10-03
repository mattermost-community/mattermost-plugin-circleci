package circle

import (
	"context"

	"github.com/antihax/optional"
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
	return user.Name, err
}

// TriggerPipeline triggers pipeline for given project and given branch/tag
func TriggerPipeline(apiToken string, projectSlug string,
	params circleci.TriggerPipelineParameters) (circleci.PipelineCreation, error) {
	var opts = circleci.PipelineApiTriggerPipelineOpts{}
	opts.Body = optional.NewInterface(params)
	pl, _, err := client.PipelineApi.TriggerPipeline(getContext(apiToken), projectSlug, &opts)
	return pl, err
}

// GetPipelineByNum get info about single pipeline
func GetPipelineByNum(apiToken string, projectSlug string, num string) (circleci.Pipeline, error) {
	pl, _, err := client.PipelineApi.GetPipelineByNumber(getContext(apiToken), projectSlug, num)
	return pl, err
}

// GetPipelineByID get info about single pipeline
func GetPipelineByID(apiToken string, pipelineID string) (circleci.Pipeline, error) {
	pl, _, err := client.PipelineApi.GetPipelineById(getContext(apiToken), pipelineID)
	return pl, err
}

// RerunWorkflow reruns a given workflow
func RerunWorkflow(apiToken string, workflowID string) (circleci.MessageResponse, error) {
	ms, _, err := client.WorkflowApi.RerunWorkflow(getContext(apiToken), workflowID, nil)
	return ms, err
}

// CancelWorkflow reruns a given workflow
func CancelWorkflow(apiToken string, workflowID string) (circleci.MessageResponse, error) {
	ms, _, err := client.WorkflowApi.CancelWorkflow(getContext(apiToken), workflowID)
	return ms, err
}

// GetEnvVarsList returns list of environment variables for given projects
func GetEnvVarsList(apiToken string, projectSlug string) ([]circleci.EnvironmentVariablePair1, error) {
	env, _, err := client.ProjectApi.ListEnvVars(getContext(apiToken), projectSlug)
	return env.Items, err
}

// AddEnvVar returns list of environment variables for given projects
func AddEnvVar(apiToken string, projectSlug string, name string, value string) error {
	opts := new(circleci.ProjectApiCreateEnvVarOpts)
	opts.Body = optional.NewInterface(circleci.EnvironmentVariablePair{Name: name, Value: value})
	_, _, err := client.ProjectApi.CreateEnvVar(getContext(apiToken), projectSlug, opts)
	return err
}

// DelEnvVar returns list of environment variables for given projects
func DelEnvVar(apiToken string, projectSlug string, name string) error {
	_, _, err := client.ProjectApi.DeleteEnvVar(getContext(apiToken), projectSlug, name)
	return err
}

// GetWorkflowMetrics retuirns workflow metrics for a project
func GetWorkflowMetrics(apiToken string, projectSlug string) ([]circleci.InlineResponse200Items, error) {
	met, _, err := client.InsightsApi.GetProjectWorkflowMetrics(getContext(apiToken), projectSlug, nil)
	return met.Items, err
}

// GetWorkflowJobsMetrics returns jobs metrics for given workflow
func GetWorkflowJobsMetrics(apiToken string, projectSlug string, workflowName string) ([]circleci.InlineResponse2002Items, error) {
	met, _, err := client.InsightsApi.GetProjectWorkflowJobMetrics(getContext(apiToken), projectSlug, workflowName, nil)
	return met.Items, err
}
// ApproveJob approves the pending job in the workflow
func ApproveJob(apiToken string, approvalRequestID string, workFlowID string) (string, error) {
	response, _, err := client.WorkflowApi.ApprovePendingApprovalJobById(getContext(apiToken), approvalRequestID, workFlowID)
	if err != nil {
		return "", err
	}

	return response.Message, nil
}
