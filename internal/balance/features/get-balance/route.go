package getbalance

import (
	"net/http"
	"strconv"
	"time"

	"github.com/andre/code-styles-golang/pkg/cqrs"
	"github.com/gin-gonic/gin"
)

// GET Balance godoc
// @Summary Return balance by user
// @Schemes
// @Description Return a balance by user
// @Tags balance
// @Accept json
// @Produce json
// @Param user_id path int true "user_id"
// @Param snapshot_date query string false "snapshot_date" Format(date-time)
// @Success 200 {object} Model
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /balance/user/{user_id} [get]
func GetBalance(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	snapshotDateStr := ctx.Query("snapshot_date")
	var snapshotDate time.Time
	if snapshotDateStr != "" {
		snapshotDate, err = time.Parse(time.RFC3339, snapshotDateStr)
		if err != nil {
			ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid snapshot_date format, use RFC3339"})
			return
		}
	} else {
		snapshotDate = time.Now()
	}

	query := &Query{
		UserID:       userID,
		SnapshotDate: snapshotDate,
	}

	response, err := cqrs.Request[*Query, *Model](ctx.Request.Context(), query)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if response == nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "no balance found"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, response)
}
