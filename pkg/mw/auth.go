package mw

import (
	userservice "blackbox-v2/internal/userservice"
	"blackbox-v2/pkg/response"
	"blackbox-v2/pkg/utils"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

var ExemptedPaths = []string{
	"/api/v1/signup/",
	"/api/v1/login/",
}
var CookieExemptedPaths = []string{
	"/app/login-view/",
	"/app/login",
}

func CookieMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		path := c.Request().URL.Path
		if utils.ArrayContains(CookieExemptedPaths, path) {
			return next(c)
		}
		cookie, err := c.Cookie("token")
		if err != nil {
			return c.HTML(http.StatusUnauthorized, "Unauthorized")
		}
		user, err := userservice.VerifyToken(cookie.Value)
		if err != nil {
			return c.HTML(http.StatusUnauthorized, "Unauthorized")
		}
		c.Request().Header.Set("user_cid", user.UserCID)
		return next(c)
	}
}

func TokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if utils.ArrayContains(ExemptedPaths, r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}
		if r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}
		if strings.Contains(r.URL.Path, "/app") {
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
