package subscriber

import (
    "golang.org/x/net/context"
    save "github.com/kazoup/platform/lib/protomsg/saveremotefiles"
    "log"
    "reflect"
    "fmt"
)

type AgentServiceTaskHandler struct{}

// Logs the message sent to this subscriber which contains file details
func(h *AgentServiceTaskHandler) LogAgentSaveMessage(ctx context.Context, msg *save.SaveRemoteFilesMessage) error {
    values := reflect.ValueOf(msg).Elem()

    for i := 0; i < values.NumField(); i++ {
        log.Println(fmt.Sprintf("%s: %v", values.Type().Field(i).Name,values.Field(i).Interface()))
    }

    return nil
}
