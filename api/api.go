package api

import (
	"blackbox-v2/internal/notes"
	"blackbox-v2/internal/userservice"
	"blackbox-v2/pkg/mw"
	"blackbox-v2/pkg/response"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func RunServer() {
	mux := mux.NewRouter()

	v1 := mux.PathPrefix("api/v1").Subrouter()
	// v1.HandleFunc("/", hello).Methods("GET")
	signupHandler := http.HandlerFunc(signup)
	v1.Handle("/signup/", mw.TokenMiddleware(signupHandler)).Methods("POST")
	loginHandler := http.HandlerFunc(login)
	v1.Handle("/login/", mw.TokenMiddleware(loginHandler)).Methods("POST")
	uploadNotesHandler := http.HandlerFunc(uploadNotes)
	v1.Handle("/upload/notes/", mw.TokenMiddleware(uploadNotesHandler)).Methods("POST")

	http.ListenAndServe(":8080", mux)
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
	}
	if signupRequest.Username == "" || signupRequest.Password == "" || signupRequest.Email == "" {
		errStr = "Username, Password and Email are required"
		response.BadRequestResponse(w, errStr)
	}
	userCid, err := userservice.SignupUser(signupRequest.Email, signupRequest.Username, signupRequest.Password)
	if err != nil {
		if userCid == "" {
			errStr = "Error creating user " + err.Error()
			response.InternalServerErrorResponse(w, errStr)
		} else {
			errStr = "User not completed signup " + err.Error()
			response.BadRequestResponse(w, errStr)
		}
	}
	response.SuccessResponse(w, "User created successfully")
}

func login(w http.ResponseWriter, r *http.Request) {
	var loginRequest LoginRequest
	var errStr string
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		errStr = "Invalid JSON " + err.Error()
		response.BadRequestResponse(w, errStr)
	}
	if loginRequest.Username == "" || loginRequest.Password == "" {
		errStr = "Username and Password are required"
		response.BadRequestResponse(w, errStr)
	}
	token, err := userservice.LoginUser(loginRequest.Username, loginRequest.Password)
	if err != nil {
		errStr = "Error logging in " + err.Error()
		response.BadRequestResponse(w, errStr)
	}
	response.JSONResponse(w, map[string]string{
		"token": token,
	})
}

func uploadNotes(w http.ResponseWriter, r *http.Request) {
	userCid := r.Header.Get("user_cid")
	if userCid == "" {
		response.UnauthorizedResponse(w, "Invalid token")
	}
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		response.BadRequestResponse(w, "Invalid file")
	}
	err = notes.SaveFile(file, fileHeader, userCid)
	if err != nil {
		response.InternalServerErrorResponse(w, "Error saving file")
	}
	response.SuccessResponse(w, "File uploaded successfully")
}

func Shutdown() {
}
