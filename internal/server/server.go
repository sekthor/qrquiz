package server

import (
	"github.com/gin-gonic/gin"
	"github.com/sekthor/puzzleinvite/internal/repo"
)

type Server struct {
	repo repo.Repo
}

func (s *Server) Run() error {

	s.repo = repo.NewInMemoryRepo()

	router := gin.Default()

	router.HTMLRender = renderer()

	router.GET("/", s.HomeHandler)
	router.GET("/quiz/:id", s.QuizHandler)
	router.GET("/new", s.NewQuizFormHandler)
	router.POST("/new", s.NewQuizHandler)

	return router.Run()
}
