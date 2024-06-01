package server

import (
	"auth/common"
	"auth/database"
	auth "auth/handler"
	"auth/logger"
	"auth/mailer"
	"auth/recaptcha"
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Server struct {
	config     *Config
	router     *gin.Engine
	httpServer *http.Server
}

func NewServer(config *Config, db *database.Database, mailer mailer.Mailer, accessLogger *zap.Logger) *Server {
	router := gin.New()
	router.LoadHTMLGlob(config.HtmlTemplates)
	handleRecovery := func(c *gin.Context, err interface{}) {
		common.ErrorResponse(c, http.StatusInternalServerError, err.(string), nil)
		c.Abort()
	}
	bodySizeMiddleware := func(c *gin.Context) {
		var w http.ResponseWriter = c.Writer
		c.Request.Body = http.MaxBytesReader(w, c.Request.Body, config.MaxBodyBytes)

		c.Next()
	}
	router.Use(gin.CustomRecovery(handleRecovery))
	router.Use(bodySizeMiddleware)
	router.Use(logger.MiddlewareFunc(accessLogger))
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, ResponseType, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	s := &http.Server{
		Addr:           config.BindAddress,
		Handler:        router,
		ReadTimeout:    time.Duration(config.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(config.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	authGroup := router.Group("/auth")
	recaptchaHandler := recaptcha.New(config.Recaptcha)
	authHandler := auth.New(db, mailer, recaptchaHandler, config.WebServer)
	authHandler.RegisterHandlers(authGroup)

	return &Server{
		config:     config,
		router:     router,
		httpServer: s,
	}
}

func (s *Server) ListenAndServer() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown() error {
	return s.httpServer.Shutdown(context.Background())
}
