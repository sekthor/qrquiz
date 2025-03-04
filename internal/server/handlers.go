package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) HomeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Title": "QR Quiz",
	})
}

func (s *Server) QuizHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "quiz.html", gin.H{
		"Title": "Quiz",
	})
}

func (s *Server) NewQuizHandler(c *gin.Context) {
	var req struct {
		Title string `form:"title" json:"title"`
	}

	if err := c.Bind(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, gin.H{"title": req.Title})

}

func (s *Server) NewQuizFormHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "form.html", gin.H{
		"Title": "Quiz",
	})
}
