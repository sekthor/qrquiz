package server

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
}

func (s *Server) Run() error {
	router := gin.Default()

	router.HTMLRender = renderer()

	router.GET("/", s.HomeHandler)
	router.GET("/quiz/:id", s.QuizHandler)
	router.GET("/new", s.NewQuizFormHandler)
	router.POST("/new", s.NewQuizHandler)

	return router.Run()
}
