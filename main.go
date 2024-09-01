package main

import (
	"database/sql"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
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

	db, err := sql.Open("sqlite3", "db/database.db")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	//Create a destination file
	dst, err := os.Create("/home/emuslu/" + file.Filename)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Copy the file content to the destination file
	if _, err := io.Copy(dst, src); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer src.Close()
	defer dst.Close()

	db_res, err := db.Exec("INSERT INTO videos (user, video_path) VALUES (?,?);", "emuslu", dst.Name())
	db.Close()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	db_res.LastInsertId()

	message := Message{Message: dst.Name()}
	return c.Render(http.StatusOK, "video", message)
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
