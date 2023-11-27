package mw

import (
	userservice "blackbox-v2/internal/userservice"
	"blackbox-v2/pkg/response"
	"blackbox-v2/pkg/utils"
	"net/http"
	"strings"
	"time"

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

func ResetCookie(c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = ""
	cookie.Path = "/"
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)
}

//  1. check if cookie exists, if not redirect to login page
//  2. check if cookie is valid, if not redirect to login page
//  3. check if path is exempted, if yes, next
//  4. set user_cid in header

func CookieMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		path := c.Request().URL.Path
		cookie, err := c.Cookie("token")
		if cookie == nil || err != nil || cookie.Value == "" {
			if utils.ArrayContains(CookieExemptedPaths, path) {
				return next(c)
			}
			return c.Redirect(http.StatusFound, "/app/login-view/")
		} else {
			claims, err := userservice.VerifyToken(cookie.Value)
			if err != nil || !claims.IsValid() {
				ResetCookie(c)
				return c.Redirect(http.StatusFound, "/app/login-view/")
			} else {

				// if utils.ArrayContains(CookieExemptedPaths, path) {
				// 	return next(c)
				// }
				c.Request().Header.Set("user_cid", claims.UserCid)
				return next(c)
			}
		}
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
		claims, err := userservice.VerifyToken(token)
		if err != nil {
			response.UnauthorizedResponse(w, "Invalid token")
			return
		}
		r.Header.Set("user_cid", claims.UserCid)
		next.ServeHTTP(w, r)
	})
}
