package main

import (
	"database/sql"
	"fmt"
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
	return c.Render(http.StatusOK, "index", message)
}

func handlePost(c echo.Context) error {
	// Retrieve the file from the form
	file, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
	}
	db.Exec("INSERT INTO videos (user, video_data) VALUES (?,?);", "emuslu", src)

	defer src.Close()

	/* Create a destination file
	dst, err := os.Create("/home/emuslu/new.mp4")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer dst.Close()

	// Copy the file content to the destination file
	if _, err := io.Copy(dst, src); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}*/

	return c.String(http.StatusOK, "File uploaded successfully")
}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/css", "css")

	e.Renderer = newTemplate()
	// Routes
	e.GET("/", handleGet)
	e.POST("/upload", handlePost)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
