package fs

import (
	"blackbox-v2/internal/fileservice"
	"blackbox-v2/internal/notes"
	"net/http"
	"path"
	"text/template"

	"github.com/labstack/echo/v4"
)

func uploadNotesView(c echo.Context) error {
	fp := path.Join("fs", "templates", "upload-notes.html")
	tmpl, err := template.ParseFiles("fs/templates/base.html", fp)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if err := tmpl.Execute(c.Response().Writer, nil); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func uploadFilesView(c echo.Context) error {
	fp := path.Join("fs", "templates", "upload-files.html")
	tmpl, err := template.ParseFiles("fs/templates/base.html", fp)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if err := tmpl.Execute(c.Response().Writer, nil); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func uploadNotes(c echo.Context) error {
	userCid := c.Request().Header.Get("user_cid")
	if userCid == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	}
	file, fileHeader, err := c.Request().FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid file")
	}
	err = notes.SaveFile(file, fileHeader, userCid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error saving file  "+err.Error())
	}
	return c.String(http.StatusOK, "File uploaded successfully")
}

func uploadFiles(c echo.Context) error {
	userCid := c.Request().Header.Get("user_cid")
	if userCid == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	}
	file, fileHeader, err := c.Request().FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid file")
	}
	err = fileservice.SaveFile(file, fileHeader, userCid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error saving file  "+err.Error())
	}
	return c.String(http.StatusOK, "File uploaded successfully")
}
func listNotes(c echo.Context) error {
	userCid := c.Request().Header.Get("user_cid")
	if userCid == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	}
	noteMds, err := notes.GetNoteMetaDataByUser(userCid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error fetching notes "+err.Error())
	}
	fp := path.Join("fs", "templates", "list-notes.html")
	tmpl, err := template.ParseFiles(fp)
	if err = tmpl.Execute(c.Response().Writer, map[string]interface{}{
		"notes": noteMds,
	}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func listFiles(c echo.Context) error {
	userCid := c.Request().Header.Get("user_cid")
	if userCid == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	}
	fileMds, err := fileservice.GetFilesByUser(userCid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error fetching files "+err.Error())
	}
	fp := path.Join("fs", "templates", "list-files.html")
	tmpl, err := template.ParseFiles(fp)
	if err = tmpl.Execute(c.Response().Writer, map[string]interface{}{
		"files": fileMds,
	}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return nil
}
