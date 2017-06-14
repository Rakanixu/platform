package subscriber

import (
	"fmt"
	save "github.com/kazoup/platform/lib/protomsg/saveremotefiles"
	"golang.org/x/net/context"
	"log"
)

type AgentServiceTaskHandler struct{}

// Logs the message sent to this subscriber which contains file details
func (h *AgentServiceTaskHandler) LogAgentSaveMessage(ctx context.Context, msg *save.SaveRemoteFilesMessage) error {
	log.Println(fmt.Sprintf("%s: %s", "UserId", msg.UserId))
	log.Println(fmt.Sprintf("%s: %s", "Index", msg.Index))
	log.Println(fmt.Sprintf("%s: %s", "Data", msg.Data))

	return nil
}
