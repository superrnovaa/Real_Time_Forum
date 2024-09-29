package forum

import (
	"net/http"
	"time"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type Session struct {
	ID        string
	UserID    int
	ExpiresAt time.Time
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
		// Get the session ID from the cookie
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "No active session", http.StatusUnauthorized)
			return
		}
		// Revoke the user's session
		revokeSession(cookie.Value)
	
		// Remove the session cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    "",
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			Expires:  time.Now(),
		}) 

		UpdateOnlineUsers()
		http.Redirect(w, r, "/", http.StatusFound)
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
			// Get the session ID from the cookie
			sessionCookie, err := r.Cookie("session_id")
			if err != nil {
				// Redirect to the login page if the session ID cookie is missing
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
			// Session cookie found, retrieve the session ID
			if !isSessionExisting(sessionCookie.Value) {
				LogoutHandler(w,r)
			} else {
				
				// Call the next handler if the session ID is valid
				next.ServeHTTP(w, r)
		}
	}
}

var sessions = make(map[string]*Session)

func createSession(userID int) *Session {
	session := &Session{
		ID:        uuid.New().String(),
		UserID:    userID,
		ExpiresAt: time.Now().Add(24 * time.Hour), // Set expiration to 24 hours
	}
	sessions[session.ID] = session

	go func() {
		waitForSocketAndUpdate(userID) 
	  }()

	return session
}

func waitForSocketAndUpdate(userID int) {
	for {
	  _, ok := userSockets[userID]
	  if ok {
		UpdateOnlineUsers()
		return
	  }
	  time.Sleep(100 * time.Millisecond)
	}
  }

func getSession(sessionID string) *Session {
	session, ok := sessions[sessionID]
	if !ok || session.ExpiresAt.Before(time.Now()) {
		return nil
	}
	return session
}

func revokeSession(sessionID string) {
	delete(sessions, sessionID)

}

func isSessionExisting(sessionID string) bool {
	_, ok := sessions[sessionID]
	if(ok) {return true}else{return false}
	
}

func CleanupExpiredSessions() {
	for sessionID, session := range sessions {
		if session.ExpiresAt.Before(time.Now()) {
			delete(sessions, sessionID)
		}
	}
}

func revokeUserSessions(userID int) {
	for _, session := range sessions {
		if session.UserID == userID {
			delete(sessions, session.ID)
		}
	}
}

func getUserIDFromSessionID(sessionID string) (int) {
	session, ok := sessions[sessionID]
	if !ok || session.ExpiresAt.Before(time.Now()) {
		return 0
	}
	return session.UserID
}
