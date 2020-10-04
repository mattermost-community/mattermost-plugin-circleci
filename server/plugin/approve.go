package plugin

import (
	"fmt"
	"net/http"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/nathanaelhoun/mattermost-plugin-circleci/server/circle"
)

func (p *Plugin) httpHandleApprove(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-Id")
	circleciToken, exists := p.Store.GetTokenForUser(userID, p.getConfiguration().EncryptionKey)
	if !exists {
		http.NotFound(w, r)
	}
	requestData := model.PostActionIntegrationRequestFromJson(r.Body)
	if requestData == nil {
		p.API.LogError("Empty request data")
		p.sendEphemeralResponse(&model.CommandArgs{}, "Cannot approve the workflow from mattermost. Please go [here](http://app.circleci.com)")
		return
	}

	workFlowID := fmt.Sprintf("%v", requestData.Context["WorkflowID"])
	jobs, err := circle.GetWorkflowJobs(circleciToken, fmt.Sprintf("%v", requestData.Context["WorkflowID"]))

	if err != nil {
		p.API.LogError("Error occurred while getting workflow jobs", err)
		// TODO: replace with actual workflow URL to approve in circle as a fallback
		p.sendEphemeralResponse(&model.CommandArgs{}, "Cannot approve the workflow from mattermost. Please go [here](http://app.circleci.com)")
		return
	}
	var approvalRequestID string
	for _, job := range *jobs {
		if job.ApprovalRequestId != "" {
			fmt.Println(fmt.Sprintf("Job with Approval request Id %s"), job.Id)
			approvalRequestID = fmt.Sprintf("%v", job.ApprovalRequestId)
		}
	}
	_, err = circle.ApproveJob(circleciToken, approvalRequestID, workFlowID)

	if err != nil {
		p.API.LogError("Error occurred while approving", err)
		// TODO: replace with actual workflow URL to approve in circle as a fallback
		p.sendEphemeralResponse(&model.CommandArgs{}, "Cannot approve the workflow from mattermost. Please go [here](http://app.circleci.com)")
	} else {
		p.sendEphemeralResponse(&model.CommandArgs{}, "Successfully approved :+1:")
	}


}
