package mw

import (
	userservice "blackbox-v2/internal/userservice"
	"blackbox-v2/pkg/response"
	"blackbox-v2/pkg/utils"
	"log"
	"net/http"
	"strings"
)

var ExemptedPaths = []string{
	"/api/v1/signup/",
	"/api/v1/login/",
	"/favicon.ico",
}
var AppExemptedPaths = []string{
	"/app/login-view/",
	"/app/",
	"/favicon.ico",
}

func CookieMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if utils.ArrayContains(AppExemptedPaths, r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}
		if strings.Contains(r.URL.Path, "/api/") {
			next.ServeHTTP(w, r)
			return
		}
		cookie, err := r.Cookie("session_id")
		if err != nil {
			response.UnauthorizedResponse(w, "No session cookie found")
			return
		}
		session, err := userservice.GetSession(cookie.Value)
		if err != nil {
			response.UnauthorizedResponse(w, "Invalid session cookie")
			return
		}
		if session.ExpiresAt.Before(session.ExpiresAt) {
			response.UnauthorizedResponse(w, "Session expired")
			return
		}
		log.Println("User CID:")
		log.Println(session.UserCID)
		r.Header.Set("user_cid", session.UserCID)
		next.ServeHTTP(w, r)
	})
}

func TokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if utils.ArrayContains(ExemptedPaths, r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}
		if strings.Contains(r.URL.Path, "/app/") {
			next.ServeHTTP(w, r)
			return
		}
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.UnauthorizedResponse(w, "No authorization header found")
			return
		}
		token := authHeader[len("Bearer "):]
		user, err := userservice.VerifyToken(token)
		if err != nil {
			response.UnauthorizedResponse(w, "Invalid token")
			return
		}
		r.Header.Set("user_cid", user.UserCID)
		next.ServeHTTP(w, r)
	})
}
