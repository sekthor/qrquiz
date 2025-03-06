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

	router.GET("/", s.HomeHandler)
	router.GET("/quiz/:id", s.QuizHandler)
	router.GET("/new", s.NewQuizFormHandler)
	router.GET("/new/question", s.NewQuestionFormHandler)
	router.GET("/new/review", s.NewQuizReviewFormHandler)
	router.POST("/new", s.NewQuizHandler)

	return router.Run(config.Listen)
}
