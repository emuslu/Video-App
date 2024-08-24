package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

type Message struct {
	Message string
}

func handleGet(c echo.Context) error {
	message := Message{Message: "Hello son"}
	return c.Render(http.StatusOK, "video", message)
}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Renderer = newTemplate()
	// Routes
	e.GET("/", handleGet)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
