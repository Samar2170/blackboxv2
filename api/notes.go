package api

import (
	"blackbox-v2/internal/notes"
	"blackbox-v2/pkg/response"
	"net/http"

	"github.com/gorilla/mux"
)

func uploadNotes(w http.ResponseWriter, r *http.Request) {
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
	err = notes.SaveFile(file, fileHeader, userCid)
	if err != nil {
		response.InternalServerErrorResponse(w, "Error saving file  "+err.Error())
		return
	}
	response.SuccessResponse(w, "File uploaded successfully")
	return
}

func listNotes(w http.ResponseWriter, r *http.Request) {
	userCid := r.Header.Get("user_cid")
	if userCid == "" {
		response.UnauthorizedResponse(w, "Invalid token")
		return
	}
	noteMds, err := notes.GetNoteMetaDataByUser(userCid)
	if err != nil {
		response.InternalServerErrorResponse(w, "Error fetching notes "+err.Error())
		return
	}
	response.JSONResponse(w, noteMds)
}

func getNote(w http.ResponseWriter, r *http.Request) {
	userCid := r.Header.Get("user_cid")
	if userCid == "" {
		response.UnauthorizedResponse(w, "Invalid token")
		return
	}
	vars := mux.Vars(r)
	noteId := vars["id"]
	if noteId == "" {
		response.BadRequestResponse(w, "Invalid note id")
		return
	}
	note, err := notes.GetNoteByNoteID(noteId)
	if err != nil {
		response.InternalServerErrorResponse(w, "Error fetching note "+err.Error())
		return
	}
	response.JSONResponse(w, note)
}
