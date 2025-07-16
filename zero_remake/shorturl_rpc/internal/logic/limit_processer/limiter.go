package limite_processer

import (
	"context"
	"sync"
	"time"
)

var ipLimiters sync.Map // map[string]*IPLimiter

// 简单IP限流器
// 每个IP每秒最多limit次
const ipRateLimit = 1
const ipRateWindow = time.Second * 60 * 30 // 30分钟

type IPLimiter struct {
	lastTime time.Time
	count    int
	lock     sync.Mutex
}

func getIPLimiter(ip string) *IPLimiter {
	limiter, ok := ipLimiters.Load(ip)
	if !ok {
		limiter = &IPLimiter{}
		ipLimiters.Store(ip, limiter)
	}
	return limiter.(*IPLimiter)
}

func AllowIP(ip string) bool {
	limiter := getIPLimiter(ip)
	limiter.lock.Lock()
	defer limiter.lock.Unlock()
	now := time.Now()
	if now.Sub(limiter.lastTime) > ipRateWindow {
		limiter.lastTime = now
		limiter.count = 1
		return true
	}
	if limiter.count < ipRateLimit {
		limiter.count++
		return true
	}
	return false
}

func GetClientIP(ctx context.Context) string {
	if md, ok := ctx.Value("metadata").(map[string]string); ok {
		if ip, ok := md["X-Real-IP"]; ok {
			return ip
		}
		if ip, ok := md["X-Forwarded-For"]; ok {
			return ip
		}
	}
	return "unknown"
}
