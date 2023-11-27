package fs

import (
	"blackbox-v2/internal/userservice"
	"blackbox-v2/pkg/mw"
	"io"
	"log"
	"net/http"
	"path"
	"text/template"
	"time"

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
	e.Use(mw.CookieMiddleware)

	subUrl.GET("/hello/", hello)
	subUrl.GET("/login-view/", loginView)
	subUrl.POST("/login", login)

	subUrl.GET("/upload-notes-view/", uploadNotesView)
	subUrl.GET("/upload-files-view/", uploadFilesView)

	subUrl.POST("/upload/notes", uploadNotes)
	subUrl.POST("/upload/files", uploadFiles)

	subUrl.GET("/list/notes/", listNotes)
	subUrl.GET("/list/files/", listFiles)

	e.Logger.Fatal(e.Start(":8080"))
}

func hello(c echo.Context) error {
	UserCid := c.Request().Header.Get("user_cid")
	user, err := userservice.GetUserByCID(UserCid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	tmpl := template.Must(template.ParseFiles("fs/templates/base.html", "fs/templates/hello.html"))
	if tmpl == nil {
		log.Println("tmpl is nil")
	}
	v := map[string]string{
		"name": user.Username,
	}
	if err := tmpl.Execute(c.Response().Writer, v); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func loginView(c echo.Context) error {
	userCid := c.Request().Header.Get("user_cid")
	if userCid != "" {
		return c.Redirect(http.StatusMovedPermanently, "/app/hello/")
	}
	fp := path.Join("fs", "templates", "login.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if err := tmpl.Execute(c.Response().Writer, nil); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	if username == "" && password == "" {
		return c.Render(http.StatusOK, "login", "Please enter username and password")
	}
	token, err := userservice.LoginUser(username, password)
	if err != nil {
		return c.Render(http.StatusOK, "login", err.Error())
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(24 * 10 * time.Hour),
		Path:     "/",
		HttpOnly: true,
	}
	c.SetCookie(cookie)
	return c.Redirect(http.StatusMovedPermanently, "/app/hello/")
}
