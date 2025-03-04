package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sekthor/puzzleinvite/internal/repo"
	"github.com/sekthor/puzzleinvite/internal/server/assets"
)

type Server struct {
	repo repo.Repo
}

func (s *Server) Run() error {

	s.repo = repo.NewInMemoryRepo()

	router := gin.Default()

	router.HTMLRender = renderer()
	router.StaticFS("/assets", http.FS(assets.Assets))

	router.GET("/", s.HomeHandler)
	router.GET("/quiz/:id", s.QuizHandler)
	router.GET("/new", s.NewQuizFormHandler)
	router.POST("/new", s.NewQuizHandler)

	return router.Run()
}
