package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimiter struct {
	mu       sync.Mutex
	visitors map[string][]time.Time
	limit    int
	window   time.Duration
}

var limiters = make(map[string]*rateLimiter)

func getLimiter(key string, limit int, window time.Duration) *rateLimiter {
	limiter, ok := limiters[key]
	if !ok {
		limiter = &rateLimiter{
			visitors: make(map[string][]time.Time),
			limit:    limit,
			window:   window,
		}
		limiters[key] = limiter
	}
	return limiter
}

func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	key := "rl_" + time.Now().String()[:10]
	l := getLimiter(key, limit, window)

	return func(c *gin.Context) {
		l.mu.Lock()
		defer l.mu.Unlock()

		ip := c.ClientIP()
		now := time.Now()
		cutoff := now.Add(-l.window)

		timestamps := l.visitors[ip]
		var valid []time.Time
		for _, t := range timestamps {
			if t.After(cutoff) {
				valid = append(valid, t)
			}
		}

		if len(valid) >= l.limit {
			c.JSON(http.StatusTooManyRequests, gin.H{"code": 429, "message": "rate limit exceeded"})
			c.Abort()
			return
		}

		l.visitors[ip] = append(valid, now)
		c.Next()
	}
}
