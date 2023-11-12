package fs

import (
	"blackbox-v2/internal/userservice"
	"log"
	"net/http"
	"time"
)

type LoginSession struct {
	UserID    int
	UserCID   string
	SessionID string
	ExpiresAt time.Time
}

func LoginView(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/login.html")
}

func Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == "" || password == "" {
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		return
	}
	session, err := userservice.LoginAppUser(username, password)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.SessionID,
		Expires:  session.ExpiresAt,
		MaxAge:   86400,
		Path:     "/",
		HttpOnly: true,
	})
	w.Write([]byte("Login successful"))
}

func Hello(w http.ResponseWriter, r *http.Request) {
	cookies := r.Cookies()
	log.Println("Cookies:")
	for _, cookie := range cookies {
		log.Println(cookie.Name)
		log.Println(cookie.Value)
	}
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Cookie:")
		log.Println(cookie.Name)
		log.Println(cookie.Value)
	}
	userCid := r.Header.Get("user_cid")
	log.Println("User CID:")
	log.Println(userCid)
	w.Write([]byte("Hello World"))
}
