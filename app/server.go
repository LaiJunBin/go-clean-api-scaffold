package app

import (
	"net/http"
	"runtime/debug"
	"strconv"
	"time"

	"go-clean-api-scaffold/api"
	sharedtypes "go-clean-api-scaffold/internal/app/shared/types"
	"go-clean-api-scaffold/internal/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Controller interface{}

type ServerImpl struct {
	api.StrictServerStub
}

func SetupHandler(s api.ServerInterface, cfg *config.Config, logger sharedtypes.Logger, m []api.MiddlewareFunc) {
	port := strconv.Itoa(cfg.Server.Port)
	allowedOrigins := cfg.Server.AllowedOrigins

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.CustomRecoveryWithWriter(nil, func(c *gin.Context, err any) {
		logger.Errorf("panic recovered: %v\n%s", err, debug.Stack())
		c.AbortWithStatus(http.StatusInternalServerError)
	}))
	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			for _, allowed := range allowedOrigins {
				if origin == allowed {
					return true
				}
			}
			return false
		},
		MaxAge: 12 * time.Hour,
	}))

	api.RegisterHandlersWithOptions(r, s, api.GinServerOptions{
		BaseURL:     "/api/v1",
		Middlewares: m,
	})

	r.Run(":" + port)
}

func NewServer(ctls []Controller) api.ServerInterface {
	server := &ServerImpl{}
	server.InitControllers(ctls)
	return api.NewStrictHandler(server, nil)
}
