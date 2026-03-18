package ratelimiter

import (
	"net/http"
	"time"

	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	clients map[string]*client
	mu      sync.Mutex
	r       rate.Limit
	burst   int
	ttl     time.Duration
}

func NewRateLimiter(r rate.Limit, burst int, ttl time.Duration) *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]*client),
		r:       r,
		burst:   burst,
		ttl:     ttl,
	}

	go rl.cleanupExpiredClients()

	return rl
}

func (rl *RateLimiter) getClient(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	c, exists := rl.clients[ip]
	if !exists {
		limiter := rate.NewLimiter(rl.r, rl.burst)
		rl.clients[ip] = &client{limiter, time.Now()}
		return limiter
	}

	c.lastSeen = time.Now()
	return c.limiter
}

func (rl *RateLimiter) cleanupExpiredClients() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		for ip, c := range rl.clients {
			if time.Since(c.lastSeen) > rl.ttl {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := rl.getClient(ip)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			return
		}

		c.Next()
	}
}
