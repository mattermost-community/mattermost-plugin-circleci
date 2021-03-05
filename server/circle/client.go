package circle

import (
	"context"
	"fmt"
	"net/http"

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
	user, resp, err := client.UserApi.GetCurrentUser(getContext(apiToken))
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetWorkflow returns the info for given workflow ID
func GetWorkflow(apiToken string, workflowID string) (*circleci.Workflow, error) {
	wf, resp, err := client.WorkflowApi.GetWorkflowById(getContext(apiToken), workflowID)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return &wf, nil
}

// GetWorkflowJobs returns the info of jobs for given workflow ID
func GetWorkflowJobs(apiToken string, workflowID string) (*[]circleci.Job, error) {
	wf, resp, err := client.WorkflowApi.ListWorkflowJobs(getContext(apiToken), workflowID)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return &wf.Items, nil
}

// GetRecentlyBuiltPipelines get all recently built pipelines in a organization
func GetRecentlyBuiltPipelines(apiToken string, orgSlug string, mine bool) ([]circleci.Pipeline1, error) {
	pl, resp, err := client.PipelineApi.ListPipelines(getContext(apiToken), orgSlug, mine, nil)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return pl.Items, nil
}

// GetAllPipelinesForProject get all pipelines for a given project
func GetAllPipelinesForProject(apiToken string, projectSlug string) ([]circleci.Pipeline1, error) {
	pl, resp, err := client.PipelineApi.ListPipelinesForProject(getContext(apiToken), projectSlug, nil)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return pl.Items, nil
}

// GetAllMyPipelinesForProject get all pipelines triggered by you
func GetAllMyPipelinesForProject(apiToken string, projectSlug string) ([]circleci.Pipeline1, error) {
	pl, resp, err := client.PipelineApi.ListMyPipelines(getContext(apiToken), projectSlug, nil)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return pl.Items, nil
}

// GetWorkflowsByPipeline get all workflows by pipeline ID
func GetWorkflowsByPipeline(apiToken string, pipelineID string) ([]circleci.Workflow1, error) {
	wf, resp, err := client.PipelineApi.ListWorkflowsByPipelineId(getContext(apiToken), pipelineID, nil)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return wf.Items, nil
}

// GetNameByID returns username from user ID
func GetNameByID(apiToken string, id string) (string, error) {
	user, resp, err := client.UserApi.GetUser(getContext(apiToken), id)
	resp.Body.Close()
	return user.Name, err
}

// TriggerPipeline triggers pipeline for given project and given branch/tag
func TriggerPipeline(apiToken string, projectSlug string,
	params circleci.TriggerPipelineParameters) (circleci.PipelineCreation, error) {
	var opts = circleci.PipelineApiTriggerPipelineOpts{}
	opts.Body = optional.NewInterface(params)

	pl, resp, err := client.PipelineApi.TriggerPipeline(getContext(apiToken), projectSlug, &opts)
	resp.Body.Close()
	return pl, err
}

// GetPipelineByNum get info about single pipeline
func GetPipelineByNum(apiToken string, projectSlug string, num string) (circleci.Pipeline, error) {
	pl, resp, err := client.PipelineApi.GetPipelineByNumber(getContext(apiToken), projectSlug, num)
	resp.Body.Close()
	return pl, err
}

// GetPipelineByID get info about single pipeline
func GetPipelineByID(apiToken string, pipelineID string) (circleci.Pipeline, error) {
	pl, resp, err := client.PipelineApi.GetPipelineById(getContext(apiToken), pipelineID)
	resp.Body.Close()
	return pl, err
}

// RerunWorkflow reruns a given workflow
func RerunWorkflow(apiToken string, workflowID string) (circleci.MessageResponse, error) {
	ms, resp, err := client.WorkflowApi.RerunWorkflow(getContext(apiToken), workflowID, nil)
	resp.Body.Close()
	return ms, err
}

// CancelWorkflow cancels a given workflow
func CancelWorkflow(apiToken string, workflowID string) (circleci.MessageResponse, error) {
	ms, resp, err := client.WorkflowApi.CancelWorkflow(getContext(apiToken), workflowID)
	resp.Body.Close()
	return ms, err
}

// GetEnvVarsList returns list of environment variables for given projects
func GetEnvVarsList(apiToken string, projectSlug string) ([]circleci.EnvironmentVariablePair1, error) {
	env, resp, err := client.ProjectApi.ListEnvVars(getContext(apiToken), projectSlug)
	resp.Body.Close()
	return env.Items, err
}

// AddEnvVar add an environment variable for given project
func AddEnvVar(apiToken string, projectSlug string, name string, value string) error {
	opts := new(circleci.ProjectApiCreateEnvVarOpts)
	opts.Body = optional.NewInterface(circleci.EnvironmentVariablePair{Name: name, Value: value})
	_, resp, err := client.ProjectApi.CreateEnvVar(getContext(apiToken), projectSlug, opts)
	resp.Body.Close()
	return err
}

// DelEnvVar delete an environment variable for given project
func DelEnvVar(apiToken string, projectSlug string, name string) error {
	_, resp, err := client.ProjectApi.DeleteEnvVar(getContext(apiToken), projectSlug, name)
	resp.Body.Close()
	return err
}

// GetWorkflowMetrics returns workflow metrics for a project
func GetWorkflowMetrics(apiToken string, projectSlug string) ([]circleci.InlineResponse200Items, error) {
	met, resp, err := client.InsightsApi.GetProjectWorkflowMetrics(getContext(apiToken), projectSlug, nil)
	resp.Body.Close()
	return met.Items, err
}

// GetWorkflowJobsMetrics returns jobs metrics for given workflow
func GetWorkflowJobsMetrics(apiToken string, projectSlug string, workflowName string) ([]circleci.InlineResponse2002Items, error) {
	met, resp, err := client.InsightsApi.GetProjectWorkflowJobMetrics(getContext(apiToken), projectSlug, workflowName, nil)
	resp.Body.Close()
	return met.Items, err
}

// ApproveJob approves the pending job in the workflow
func ApproveJob(apiToken string, approvalRequestID string, workFlowID string) (string, error) {
	response, resp, err := client.WorkflowApi.ApprovePendingApprovalJobById(getContext(apiToken), approvalRequestID, workFlowID)
	resp.Body.Close()
	if err != nil {
		return "", err
	}

	return response.Message, nil
}

// EnvVarExists check if given env var exists
func EnvVarExists(apiToken string, projectSlug string, name string) (circleci.EnvironmentVariablePair, bool, error) {
	val, resp, err := client.ProjectApi.GetEnvVar(getContext(apiToken), projectSlug, name)
	defer resp.Body.Close()
	if err != nil {
		return val, false, err
	}

	switch {
	case resp.StatusCode < 300:
		return val, true, nil

	case resp.StatusCode == http.StatusNotFound:
		return val, false, nil

	default:
		return val, false, fmt.Errorf("error while checking if env var %s exist", name)
	}
}
