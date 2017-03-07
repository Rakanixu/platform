package quota

import (
	"fmt"
	"github.com/kazoup/platform/lib/globals"
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

// GetQuota returns quota info: rate, resetTimestamp, and quota
func GetQuota(srvName, uID string) (int64, int64, int64) {

	hq := int64(globals.SRV_LIMIT_DICTIONARY.M[srvName]["handler"])
	sq := int64(globals.SRV_LIMIT_DICTIONARY.M[srvName]["subscriber"])

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

	return rate1 + rate2, ts, hq + sq
}
