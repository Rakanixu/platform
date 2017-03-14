package quota

import (
	"fmt"
	"github.com/kazoup/platform/lib/globals"
	"golang.org/x/net/context"
	timerate "golang.org/x/time/rate"
	"gopkg.in/go-redis/rate.v5"
	"gopkg.in/redis.v5"
	"time"
)

var limiter *rate.Limiter

func init() {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server1": "redis:6379",
		},
	})
	fallbackLimiter := timerate.NewLimiter(timerate.Every(time.Second), 1000)
	limiter = rate.NewLimiter(ring, fallbackLimiter)
}

// GetQuota returns quota info: srvLabel, icon, rate, resetTimestamp, quota, and if was OK
func GetQuota(ctx context.Context, srvName, uID string) (string, string, int64, int64, int64, bool) {
	r, err := globals.ParseRolesFromContext(ctx)
	if err != nil {
		return "", "", 0, 0, 0, false
	}

	var product string
	for _, v := range r {
		switch v {
		case globals.PRODUCT_TYPE_PERSONAL, globals.PRODUCT_TYPE_TEAM, globals.PRODUCT_TYPE_ENTERPRISE:
			product = v
		}
	}

	if globals.PRODUCT_QUOTAS.M[product] == nil {
		return "", "", 0, 0, 0, false
	}

	if globals.PRODUCT_QUOTAS.M[product][srvName] == nil {
		return "", "", 0, 0, 0, false
	}

	hq := int64(globals.PRODUCT_QUOTAS.M[product][srvName]["handler"].(int))
	sq := int64(globals.PRODUCT_QUOTAS.M[product][srvName]["subscriber"].(int))
	sl := globals.PRODUCT_QUOTAS.M[product][srvName]["label"].(string)
	i := globals.PRODUCT_QUOTAS.M[product][srvName]["icon"].(string)

	// Get quota for srv handler
	rate1, resetTimestamp1, _ := limiter.AllowN(fmt.Sprintf("%s-handler-%s", srvName, uID), hq, globals.QUOTA_TIME_LIMITER, 0)
	// Get quota for srv subscriber
	rate2, resetTimestamp2, _ := limiter.AllowN(fmt.Sprintf("%s-subs-%s", srvName, uID), sq, globals.QUOTA_TIME_LIMITER, 0)

	// Merge info
	var ts int64
	if resetTimestamp1 > resetTimestamp2 {
		ts = resetTimestamp1
	} else {
		ts = resetTimestamp2
	}

	return sl, i, rate1 + rate2, ts, hq + sq, true
}
