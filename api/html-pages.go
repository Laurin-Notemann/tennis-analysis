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

func RegisterHtmlPageRoutes(e *echo.Echo, middleware Middleware) {
	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	e.Renderer = t
  e.Static("/static", "public/static")
	e.GET("/", indexRoute)
	e.GET("/login", loginRoute)
	e.GET("/register", registerRoute)
	e.GET("/create-player", createPlayerRoute)
	e.GET("/players", playersRoute)
  e.GET("/edit-player/:id", editPlayerRoute)
	e.GET("/create-team", createTeamRoute)
	e.GET("/teams", teamsRoute)
  e.GET("/edit-team/:id", editTeamRoute)
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

func playersRoute(c echo.Context) error {
  return c.Render(http.StatusOK, "players.html", "")
}

func createPlayerRoute(c echo.Context) error {
  return c.Render(http.StatusOK, "create-player.html", "")
}

func editPlayerRoute(c echo.Context) error {
  return c.Render(http.StatusOK, "edit-player.html", "")
}

func teamsRoute(c echo.Context) error {
  return c.Render(http.StatusOK, "teams.html", "")
}

func createTeamRoute(c echo.Context) error {
  return c.Render(http.StatusOK, "create-team.html", "")
}

func editTeamRoute(c echo.Context) error {
  return c.Render(http.StatusOK, "edit-team.html", "")
}
