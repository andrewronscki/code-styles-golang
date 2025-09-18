package routes

import (
	_ "github.com/andre/code-styles-golang/api"
	balanceroutes "github.com/andre/code-styles-golang/internal/balance/routes"
	"github.com/andre/code-styles-golang/pkg/config"
	"github.com/andre/code-styles-golang/pkg/datadog/env"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
)

func Router(ddEnvs *env.DatadogEnvironment) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	env := config.Env.GetString("GOLANG_ENV")

	if env != "production" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	if ddEnvs.DATADOG_ENABLED {
		r.Use(gintrace.Middleware(config.Env.GetString("DD_SERVICE")))
	}

	api := r.Group("/api")
	{
		balanceroutes.Map(api.Group("/balance"))
	}

	return r
}
