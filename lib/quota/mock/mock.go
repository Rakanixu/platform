package mock

import (
	"github.com/kazoup/platform/lib/quota"
	"golang.org/x/net/context"
)

type Mock struct{}

func init() {
	quota.Register(new(Mock))
}

// Check returns quota info: srvLabel, icon, rate, resetTimestamp, quota, and if was OK
func (m *Mock) Check(ctx context.Context, srvName, uID string) (string, string, int64, int64, int64, bool) {
	return "", "", 0, 0, 0, false
}

//TODO: implement
func (m *Mock) Increment(ctx context.Context) error {
	return nil
}
