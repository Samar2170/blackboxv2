package api

import (
	"blackbox-v2/internal/fileservice"
	"blackbox-v2/pkg/response"
	"bytes"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func UploadFile(w http.ResponseWriter, r *http.Request) {
	userCid := r.Header.Get("user_cid")
	if userCid == "" {
		response.UnauthorizedResponse(w, "Invalid token")
		return
	}
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		response.BadRequestResponse(w, "Invalid file")
		return
	}
	err = fileservice.SaveFile(file, fileHeader, userCid)
	if err != nil {
		response.InternalServerErrorResponse(w, "Error saving file  "+err.Error())
		return
	}
	response.SuccessResponse(w, "File uploaded successfully")
	return

}

func listFiles(w http.ResponseWriter, r *http.Request) {
	userCid := r.Header.Get("user_cid")
	if userCid == "" {
		response.UnauthorizedResponse(w, "Invalid token")
		return
	}
	noteMds, err := fileservice.GetFilesByUser(userCid)
	if err != nil {
		response.InternalServerErrorResponse(w, "Error fetching notes "+err.Error())
		return
	}
	response.JSONResponse(w, noteMds)

}
func getFile(w http.ResponseWriter, r *http.Request) {
	userCid := r.Header.Get("user_cid")
	if userCid == "" {
		response.UnauthorizedResponse(w, "Invalid token")
		return
	}
	vars := mux.Vars(r)
	fileId := vars["id"]
	if fileId == "" {
		response.BadRequestResponse(w, "Invalid note id")
		return
	}
	data, fmd, err := fileservice.GetFileByID(fileId)
	if err != nil {
		response.InternalServerErrorResponse(w, "Error fetching file "+err.Error())
		return
	}
	if err != nil {
		fmt.Fprint(w, err)
	}
	http.ServeContent(w, r, fmd.FilePath, fmd.CreatedAt, bytes.NewReader(data))
}
