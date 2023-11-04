package api

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}
	return t.templates.ExecuteTemplate(w, name, data)
}

func RegisterHtmlPageRoutes(e *echo.Echo) {
	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	e.Renderer = t
  e.Static("/static", "public/static")
	e.GET("/", indexRoute)
	e.GET("/login", loginRoute)
	e.GET("/register", registerRoute)
}

func indexRoute(c echo.Context) error {
  return c.Render(http.StatusOK, "index.html", "")
}

func loginRoute(c echo.Context) error {
  return c.Render(http.StatusOK, "login.html", "")
}

func registerRoute(c echo.Context) error {
  return c.Render(http.StatusOK, "register.html", "")
}
