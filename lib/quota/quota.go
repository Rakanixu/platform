package quota

import (
	"golang.org/x/net/context"
)

type Quota interface {
	Checker
	Incrementer
}

type Checker interface {
	// Check returns quota info: srvLabel, icon, rate, resetTimestamp, quota, and if was OK
	Check(ctx context.Context, srvName, uID string) (string, string, int64, int64, int64, bool)
}

type Incrementer interface {
	Increment(ctx context.Context) error
}

var (
	quota Quota
)

func Register(q Quota) {
	quota = q
}

func Check(ctx context.Context, srvName, uID string) (string, string, int64, int64, int64, bool) {
	return quota.Check(ctx, srvName, uID)
}

func Increment(ctx context.Context) error {
	return quota.Increment(ctx)
}
