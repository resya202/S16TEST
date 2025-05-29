package api

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"gorm.io/gorm"

	"github.com/resya202/S16TEST/internal/model"
)

// rate limiter store per-IP
var (
	visitors = make(map[string]*rate.Limiter)
	mu       sync.Mutex
)

// getVisitor returns the rate.Limiter for a given IP
func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	lim, exists := visitors[ip]
	if !exists {
		// default: 5 req/sec with burst 10
		lim = rate.NewLimiter(rate.Limit(5), 10)
		visitors[ip] = lim
	}
	return lim
}

// RateLimitMiddleware enforces a simple per-IP rate limit
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getVisitor(ip)
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			return
		}
		c.Next()
	}
}

// RegisterRoutes mounts the API endpoints onto Gin
func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	// apply rate limiting globally
	r.Use(RateLimitMiddleware())

	v1 := r.Group("/api/v1/validators")
	v1.GET("/:valAddr/delegations/hourly", getHourly(db))
	v1.GET("/:valAddr/delegations/daily", getDaily(db))
	v1.GET("/:valAddr/delegator/:delAddr/history", getDelegatorHistory(db))
}

func getHourly(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		val := c.Param("valAddr")
		var rows []model.HourlyDelegation
		db.Where("validator_addr = ?", val).Find(&rows)
		c.JSON(http.StatusOK, rows)
	}
}

func getDaily(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		val := c.Param("valAddr")
		type DailyResponse struct {
			ValidatorAddr    string `json:"validator_addr"`
			DelegatorAddr    string `json:"delegator_addr"`
			Date             string `json:"date"`
			TotalChangeUAtom int64  `json:"total_change_uatom"`
		}
		var rows []DailyResponse
		db.Table("delegation_hourly").
			Select("validator_addr, delegator_addr, DATE(timestamp) as date, SUM(change_uatom) as total_change_uatom").
			Where("validator_addr = ?", val).
			Group("validator_addr, delegator_addr, DATE(timestamp)").
			Scan(&rows)
		c.JSON(http.StatusOK, rows)
	}
}

func getDelegatorHistory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		val := c.Param("valAddr")
		del := c.Param("delAddr")
		var rows []model.HourlyDelegation
		db.Where("validator_addr = ? AND delegator_addr = ?", val, del).
			Order("timestamp").
			Find(&rows)
		c.JSON(http.StatusOK, rows)
	}
}
