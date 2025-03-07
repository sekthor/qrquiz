package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sekthor/qrquiz/internal/domain"
)

func (s *Server) HomeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Title": "QR Quiz",
	})
}

func (s *Server) QuizHandler(c *gin.Context) {

	id := c.Param("id")

	quiz, err := s.repo.GetQuiz(id)

	if err != nil {
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"Title": "Quiz",
		})
		return
	}

	c.HTML(http.StatusOK, "quiz.html", gin.H{
		"Title": quiz.Title,
		"Quiz":  quiz,
	})
}

func (s *Server) NewQuizHandler(c *gin.Context) {
	var req struct {
		Title     string            `form:"title" json:"title"`
		Secret    string            `form:"secret" json:"secret"`
		Questions []domain.Question `json:"questions"`
	}

	if err := c.Bind(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	quiz, err := domain.NewQuiz(req.Title, req.Secret, req.Questions)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := s.repo.Save(quiz); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, &quiz)
}

func (s *Server) NewQuizFormHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "form.html", gin.H{
		"Title": "Create a Quiz",
	})
}

func (s *Server) NewQuestionFormHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "question.html", gin.H{
		"Title": "Add Questions",
	})
}

func (s *Server) NewQuizReviewFormHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "review.html", gin.H{
		"Title": "Review Quiz",
	})
}

func (s *Server) QuizlistHandler(c *gin.Context) {

	quiz, _ := s.repo.List()

	c.HTML(http.StatusOK, "list.html", gin.H{
		"Title":    "Quiz List",
		"Quizlist": quiz,
	})
}
