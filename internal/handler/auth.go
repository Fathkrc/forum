package handler

import "net/http"

const sessionCookieName = "session"

type authData struct {
	Error    string
	Email    string
	Username string
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		renderTemplate(w, "login", authData{})
	case http.MethodPost:
		email := r.FormValue("email")
		password := r.FormValue("password")

		session, err := h.svc.Login(r.Context(), email, password)
		if err != nil {
			renderTemplate(w, "login", authData{Error: "Invalid email or password", Email: email})
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     sessionCookieName,
			Value:    session.Token,
			Path:     "/",
			Expires:  session.ExpiresAt,
			HttpOnly: true,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		renderTemplate(w, "register", authData{})
	case http.MethodPost:
		email := r.FormValue("email")
		username := r.FormValue("username")
		password := r.FormValue("password")
		confirm := r.FormValue("confirmPassword")

		err := h.svc.Register(r.Context(), email, username, password, confirm)
		if err != nil {
			renderTemplate(w, "register", authData{Error: err.Error(), Email: email, Username: username})
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	_ = h.svc.Logout(r.Context(), cookie.Value)

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // expires it immediately
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
