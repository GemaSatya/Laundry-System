package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/GemaSatya/LaundrySystem/database"
	"github.com/GemaSatya/LaundrySystem/model"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	// Implementation for user registration
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")

	if !SearchUser(username) {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Create new user
	newUser := model.User{
		Username: username,
		Password: hashedPassword,
		Email:    email,
	}

	err = database.DB.Create(&newUser).Error
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Serve login form
		errorMsg := r.URL.Query().Get("error")
		data := struct {
			Error string
		}{
			Error: errorMsg,
		}

		t, err := template.ParseFiles("auth/login.html")
		if err != nil {
			http.Error(w, "Failed to load template", http.StatusInternalServerError)
			return
		}

		if err := t.Execute(w, data); err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	username := r.FormValue("username")
	password := r.FormValue("password")

	if SearchUser(username) {
		http.Redirect(w, r, "/login?error=User+does+not+exist", http.StatusSeeOther)
		return
	}

	var user model.User
	var login model.Login

	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		http.Redirect(w, r, "/login?error=Invalid+username+or+password", http.StatusSeeOther)
		return
	}

	if !CheckPasswordHash(password, user.Password) {
		http.Redirect(w, r, "/login?error=Invalid+username+or+password", http.StatusSeeOther)
		return
	}

	if !SearchToken(user.ID) {
		http.Redirect(w, r, "/login?error=User+already+logged+in", http.StatusSeeOther)
		return
	}

	sessionToken := GenerateToken(32)
	csrfToken := GenerateToken(32)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(30 * time.Second),
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  time.Now().Add(30 * time.Second),
		HttpOnly: false,
	})

	login = model.Login{
		HashedPassword: user.Password,
		SessionToken:   sessionToken,
		CSRFToken:      csrfToken,
		SessionId:      user.ID,
	}

	err := database.DB.Create(&login).Error
	if err != nil {
		http.Redirect(w, r, "/login?error=Error+creating+login+session", http.StatusSeeOther)
		return
	}

	// Redirect to home page on success
	http.Redirect(w, r, "/frontend", http.StatusSeeOther)
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var session model.Login

	st, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Session token not found", http.StatusUnauthorized)
		return
	}

	if err := database.DB.Where("session_token = ?", st.Value).First(&session).Error; err != nil {
		http.Error(w, "Invalid session token", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: false,
	})

	CleanUpToken(false, session.SessionId)

	if err := json.NewEncoder(w).Encode(map[string]any{
		"msg": "Logout successful",
	}); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	log.Println("Log out successful")

}