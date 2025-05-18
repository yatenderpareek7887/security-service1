package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter holds the rate limiters for each client IP
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

// NewRateLimiter creates a new rate limiter with the specified rate (requests per second) and burst
func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     r,
		burst:    b,
	}
}

// GetLimiter returns or creates a rate limiter for the given IP
func (rl *RateLimiter) GetLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[ip] = limiter
		// Clean up old limiters periodically
		go rl.cleanup(ip)
	}
	return limiter
}

// cleanup removes a limiter after a period of inactivity
func (rl *RateLimiter) cleanup(ip string) {
	time.Sleep(1 * time.Hour)
	rl.mu.Lock()
	delete(rl.limiters, ip)
	rl.mu.Unlock()
}

// RateLimitMiddleware creates Gin middleware for rate limiting
func RateLimitMiddleware(r rate.Limit, b int) gin.HandlerFunc {
	rl := NewRateLimiter(r, b)
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := rl.GetLimiter(ip)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
