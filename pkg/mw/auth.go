package mw

import (
	userservice "blackbox-v2/internal/userservice"
	"blackbox-v2/pkg/response"
	"blackbox-v2/pkg/utils"
	"net/http"
)

var ExemptedPaths = []string{
	"/api/v1/signup/",
	"/api/v1/login/",
}

func TokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if utils.ArrayContains(ExemptedPaths, r.URL.Path) {
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
