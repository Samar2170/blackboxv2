package fs

import (
	"io"
	"net/http"
	"text/template"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func StartEchoServer() {
	e := echo.New()
	subUrl := e.Group("/app")
	t := &Template{
		templates: template.Must(template.ParseGlob("fs/templates/*.html")),
	}
	e.Renderer = t
	subUrl.GET("/hello", Hello)
	e.Logger.Fatal(e.Start(":8080"))
}

func Hello(c echo.Context) error {
	return c.Render(http.StatusOK, "hello", "World")
}
