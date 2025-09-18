package routes

import (
	"github.com/gin-gonic/gin"
	getbalance "github.com/andre/code-styles-golang/internal/balance/features/get-balance"
)

func Map(g *gin.RouterGroup) {
	g.GET("/user/:user_id", getbalance.GetBalance)

}
