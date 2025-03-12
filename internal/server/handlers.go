package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sekthor/qrquiz/internal/domain"
	"github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
)

func (s *Server) HomeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Title": "QR Quiz",
	})
}

func (s *Server) QuizHandler(c *gin.Context) {

	id := c.Param("id")

	quiz, err := s.repo.GetQuiz(c.Request.Context(), id)

	if err != nil {
		logrus.WithContext(c.Request.Context()).
			WithError(err).
			Debug("could not get quiz from db")
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
		logrus.WithContext(c.Request.Context()).
			WithError(err).
			Debug("could not unmarshal quiz")
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not unmarshal quiz"})
		return
	}

	quiz, err := domain.NewQuiz(req.Title, req.Secret, req.Questions)
	if err != nil {
		logrus.WithContext(c.Request.Context()).
			WithError(err).
			Debug("could not create new quiz")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.repo.Save(c.Request.Context(), quiz); err != nil {
		logrus.WithContext(c.Request.Context()).
			WithError(err).
			Debug("could not save new quiz")
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not save new quiz"})
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
	pageStr := c.Param("page")
	if pageStr == "" {
		pageStr = "1"
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		logrus.WithContext(c.Request.Context()).
			WithError(err).
			Debug("could not convert :page path param to integer")
		c.Status(http.StatusBadRequest)
		return
	}

	if page < 1 {
		page = 1
	}

	quiz, _ := s.repo.List(c.Request.Context(), page, 100)

	c.HTML(http.StatusOK, "list.html", gin.H{
		"Title":    "Quiz List",
		"Quizlist": quiz,
	})
}

func (s *Server) QrHandler(c *gin.Context) {
	data := c.Query("q")

	png, err := qrcode.Encode(data, qrcode.Medium, 256)
	if err != nil {
		logrus.WithContext(c.Request.Context()).
			WithError(err).
			Debug("could not encode query to qr code")
		c.Status(http.StatusBadRequest)
		return
	}

	c.Data(http.StatusOK, "image/png", png)
}

func (s *Server) ImprintHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "imprint.html", gin.H{
		"Title":   "Imprint",
		"Version": tag,
		"Commit":  commit,
		"Date":    date,
		"Contact": s.config.Contact,
	})
}
