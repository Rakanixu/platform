package mock

import (
	"fmt"
	"github.com/kazoup/platform/lib/quota"
	"github.com/micro/go-micro/metadata"
	"golang.org/x/net/context"
)

type Mock struct{}

func init() {
	quota.Register(new(Mock))
}

// Check returns quota info: srvLabel, icon, rate, resetTimestamp, quota, and if was OK
func (m *Mock) Check(ctx context.Context, srvName, uID string) (string, string, int64, int64, int64, bool) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		panic(fmt.Sprintf("Mock: metadata.FromContext"))
	}

	if len(md["Quota-Exceeded"]) == 0 {
		panic(fmt.Sprintf("Mock: quota_exceeded required"))
	}

	if md["Quota-Exceeded"] == "true" {
		// Rate: 10, Quota: 0
		return "", "", 10, 0, 0, true
	}

	if md["Quota-Exceeded"] == "false" {
		// Rate: 10, Quota: 11
		return "", "", 10, 0, 11, true
	}

	return "", "", 0, 0, 0, false
}

//TODO: implement
func (m *Mock) Increment(ctx context.Context) error {
	return nil
}
