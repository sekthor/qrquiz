package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sekthor/qrquiz/internal/config"
	"github.com/sekthor/qrquiz/internal/repo"
	"github.com/sekthor/qrquiz/internal/server/assets"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Server struct {
	config *config.Config
	repo   repo.Repo
	tracer trace.Tracer
}

func (s *Server) Run(config *config.Config) error {

	s.config = config

	switch config.Database {
	case "sqlite":
		s.repo = repo.NewSqliteRepo()
	default:
		s.repo = repo.NewInMemoryRepo()
	}

	s.tracer = otel.Tracer("server")

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(otelgin.Middleware("qrquiz"))
	router.Use(GinLogger())

	router.HTMLRender = renderer()

	// serve static assets with optional configurable cache policy
	asset := router.Group("/assets")
	if config.StaticCacheMaxAge != 0 {
		asset.Use(func(c *gin.Context) {
			c.Header("Cache-Control", fmt.Sprintf("max-age=%d", config.StaticCacheMaxAge))
		})
	}
	asset.StaticFS("", http.FS(assets.Assets))

	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"Title": "Not Found",
		})
	})

	router.GET("/", s.HomeHandler)
	router.GET("/quiz/:id", s.QuizHandler)
	router.GET("/new", s.NewQuizFormHandler)
	router.GET("/new/question", s.NewQuestionFormHandler)
	router.GET("/new/review", s.NewQuizReviewFormHandler)
	router.POST("/new", s.NewQuizHandler)
	router.GET("/list", s.QuizlistHandler)
	router.GET("/list/:page", s.QuizlistHandler)
	router.GET("/qr", s.QrHandler)
	router.GET("/imprint", s.ImprintHandler)

	// TODO: make interval configureable
	logrus.Infof("deleting expired quizzes every %d minutes", 15)
	go func() {
		for {
			func() {
				ctx, span := s.tracer.Start(context.Background(), "deleteExpired")
				defer span.End()
				s.repo.DeleteExpired(ctx)
				time.Sleep(time.Minute * 15)
			}()
		}
	}()

	logrus.Infof("starting server, listening on '%s'", config.Listen)
	return router.Run(config.Listen)
}
