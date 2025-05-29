package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/resya202/S16TEST/internal/model"
)

// RegisterRoutes mounts the API endpoints onto Gin
func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	v1 := r.Group("/api/v1/validators")
	v1.GET("/:valAddr/delegations/hourly", getHourly(db))
	v1.GET("/:valAddr/delegations/daily", getDaily(db))
}

func getHourly(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		val := c.Param("valAddr")
		var rows []model.HourlyDelegation
		db.Where("validator_addr = ?", val).Find(&rows)
		c.JSON(http.StatusOK, rows)
	}
}

// getDaily aggregates hourly snapshots into daily totals
func getDaily(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		val := c.Param("valAddr")
		// DailyResponse defines the aggregated structure
		type DailyResponse struct {
			ValidatorAddr    string `json:"validator_addr"`
			DelegatorAddr    string `json:"delegator_addr"`
			Date             string `json:"date"`
			TotalChangeUAtom int64  `json:"total_change_uatom"`
		}
		var rows []DailyResponse

		// Aggregate hourly change by date
		db.Table("delegation_hourly").
			Select("validator_addr, delegator_addr, DATE(timestamp) as date, SUM(change_uatom) as total_change_uatom").
			Where("validator_addr = ?", val).
			Group("validator_addr, delegator_addr, DATE(timestamp)").
			Scan(&rows)

		c.JSON(http.StatusOK, rows)
	}
}
