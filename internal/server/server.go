package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sekthor/qrquiz/internal/config"
	"github.com/sekthor/qrquiz/internal/repo"
	"github.com/sekthor/qrquiz/internal/server/assets"
)

type Server struct {
	repo repo.Repo
}

func (s *Server) Run(config *config.Config) error {

	switch config.Database {
	case "sqlite":
		s.repo = repo.NewSqliteRepo()
	default:
		s.repo = repo.NewInMemoryRepo()
	}

	router := gin.Default()

	router.HTMLRender = renderer()
	router.StaticFS("/assets", http.FS(assets.Assets))

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

	return router.Run(config.Listen)
}
