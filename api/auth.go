package api

import (
	"blackbox-v2/internal/userservice"
	"blackbox-v2/pkg/response"
	"encoding/json"
	"net/http"
)

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
