package api

import (
	"blackbox-v2/pkg/mongoc"
	"blackbox-v2/pkg/mw"
	"blackbox-v2/pkg/response"
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

var host = "localhost"

func RunServer() {

	mux := mux.NewRouter()

	v1 := mux.PathPrefix("/api/v1").Subrouter()
	// v1.HandleFunc("/", hello).Methods("GET")
	signupHandler := http.HandlerFunc(signup)
	v1.Handle("/signup/", signupHandler).Methods("POST")
	loginHandler := http.HandlerFunc(login)
	v1.Handle("/login/", loginHandler).Methods("POST")

	uploadNotesHandler := http.HandlerFunc(uploadNotes)
	v1.Handle("/upload/notes/", uploadNotesHandler).Methods("POST")
	listNotesHandler := http.HandlerFunc(listNotes)
	v1.Handle("/list/notes/", listNotesHandler).Methods("GET")
	getNoteHandler := http.HandlerFunc(getNote)
	v1.Handle("/get/note/{id}", getNoteHandler).Methods("GET")

	uploadFileHandler := http.HandlerFunc(UploadFile)
	v1.Handle("/upload/file/", uploadFileHandler).Methods("POST")
	listFileHandler := http.HandlerFunc(listFiles)
	v1.Handle("/list/file/", listFileHandler).Methods("GET")
	getFileHandler := http.HandlerFunc(getFile)
	v1.Handle("/get/file/{id}", getFileHandler).Methods("GET")

	healthCheckHandler := http.HandlerFunc(healthCheck)
	v1.Handle("/health/", healthCheckHandler).Methods("GET")

	wrappedMux := mw.LogRequest(mux)
	wrappedMux = mw.TokenMiddleware(wrappedMux)
	mw.LogIt("Server running on " + host + ":8081")
	http.ListenAndServe(":8081", wrappedMux)
}

func Shutdown() {
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	err := mongoc.MongoClient.Ping(context.Background(), nil)
	if err != nil {
		response.InternalServerErrorResponse(w, "Mongo setup failing")
		return
	}
	response.SuccessResponse(w, "Looks good")
	return
}
