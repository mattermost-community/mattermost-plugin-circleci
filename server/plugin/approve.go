package plugin

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (p *Plugin) httpHandleApprove(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-Id")
	circleciToken, exists := p.Store.GetTokenForUser(userID, p.getConfiguration().EncryptionKey)
	if !exists {
		http.NotFound(w, r)
	}
	vars := mux.Vars(r)
	fmt.Println("===============")
	fmt.Println(fmt.Sprintf("Workflowid %s", vars["workflowID"]))
	fmt.Println(fmt.Sprintf("jobID %s", circleciToken))
	fmt.Println("===============")

}
