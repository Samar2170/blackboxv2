package api

import (
	"blackbox-v2/internal/notes"
	"blackbox-v2/internal/userservice"
	"blackbox-v2/pkg/mongoc"
	"blackbox-v2/pkg/mw"
	"blackbox-v2/pkg/response"
	"context"
	"encoding/json"
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

	wrappedMux := mw.LogRequest(mux)
	wrappedMux = mw.TokenMiddleware(wrappedMux)
	mw.LogIt("Server running on " + host + ":8080")
	http.ListenAndServe(":8080", wrappedMux)
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

func signup(w http.ResponseWriter, r *http.Request) {
	var signupRequest SignupRequest
	var errStr string
	err := json.NewDecoder(r.Body).Decode(&signupRequest)
	if err != nil {
		errStr = "Invalid JSON " + err.Error()
		response.BadRequestResponse(w, errStr)
		return
	}
	if signupRequest.Username == "" || signupRequest.Password == "" || signupRequest.Email == "" {
		errStr = "Username, Password and Email are required"
		response.BadRequestResponse(w, errStr)
		return
	}
	userCid, err := userservice.SignupUser(signupRequest.Email, signupRequest.Username, signupRequest.Password)
	if err != nil {
		if userCid == "" {
			errStr = "Error creating user " + err.Error()
			response.InternalServerErrorResponse(w, errStr)
			return
		} else {
			errStr = "User not completed signup " + err.Error()
			response.BadRequestResponse(w, errStr)
			return
		}
	}
	response.SuccessResponse(w, "User created successfully")
	return
}

func login(w http.ResponseWriter, r *http.Request) {
	var loginRequest LoginRequest
	var errStr string
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		errStr = "Invalid JSON " + err.Error()
		response.BadRequestResponse(w, errStr)
		return
	}
	if loginRequest.Username == "" || loginRequest.Password == "" {
		errStr = "Username and Password are required"
		response.BadRequestResponse(w, errStr)
		return
	}
	token, err := userservice.LoginUser(loginRequest.Username, loginRequest.Password)
	if err != nil {
		errStr = "Error logging in " + err.Error()
		response.BadRequestResponse(w, errStr)
		return
	}
	response.JSONResponse(w, map[string]string{
		"token": token,
	})
	return
}

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
