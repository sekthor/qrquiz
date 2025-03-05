package server

import (
	"embed"

	"github.com/gin-contrib/multitemplate"
)

//go:embed templates/*.html
var templates embed.FS

func renderer() multitemplate.Renderer {
	renderer := multitemplate.NewRenderer()
	renderer.AddFromFS("index.html", templates, "templates/base.html", "templates/index.html")
	renderer.AddFromFS("quiz.html", templates, "templates/base.html", "templates/quiz.html")
	renderer.AddFromFS("form.html", templates, "templates/base.html", "templates/form.html")
	renderer.AddFromFS("question.html", templates, "templates/base.html", "templates/question.html")
	renderer.AddFromFS("review.html", templates, "templates/base.html", "templates/review.html")
	renderer.AddFromFS("404.html", templates, "templates/base.html", "templates/404.html")
	return renderer
}
