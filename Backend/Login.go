package forum

import (
	"database/sql"
	"encoding/json"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		Template(w)
		return
	}

	if r.Method == "POST" {
		var account LoginData

		err := json.NewDecoder(r.Body).Decode(&account)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		errorMessage := ""

		// check if account Exists
		exists, err := AccountExists(account.Username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if account.Username == "" || account.Password == "" {
			errorMessage = "Please Fill All The Boxes"
		} else if !exists {
			errorMessage = "This Account Does Not Exist. Please Try Again"
		} else {
			var password string
			ID, err := GetAccountID(account.Username)
			if ID == 0 {
				errorMessage = "This Account Does Not Exist. Please Try Again"
			} else {
				er := Accountsdb.QueryRow("SELECT Password FROM accounts WHERE id = ?", ID).Scan(&password)
				if !decryption(account.Password, password) {
					errorMessage = "worng password"
					if er != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				} else {
					w.WriteHeader(http.StatusOK)
					session := createSession(ID)
					// Set the session cookie in the response
					http.SetCookie(w, &http.Cookie{
						Name:     "session_id",
						Value:    session.ID,
						HttpOnly: true,
						Path:    "/",
						Expires: session.ExpiresAt,
					})
					// Prepare the response data
					responseData := struct {
						SessionID string `json:"sessionId"`
					}{
						SessionID: session.ID,
					}

					// Write the response as JSON
					w.Header().Set("Content-Type", "application/json")
					err = json.NewEncoder(w).Encode(responseData)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}

					http.Redirect(w, r, "/HomePage", http.StatusFound)
				}
			}
		}

		if errorMessage != "" {
			http.Error(w, errorMessage, http.StatusInternalServerError)
			return
		}

	} else {
		Template(w)
	}
}

func Error(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Content-Type", "application/json")
	//http.Error(w, "Page Not Found", 404)

    errorResponse := ErrorResponse{
        Code:    http.StatusNotFound,
        Message: "Page Not Found",
    }

    err := json.NewEncoder(w).Encode(errorResponse)
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
	//Template(w)
}


func encrytion(password string) string {
	// Generate a bcrypt hash of the password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	// Convert the hash to a string and store it in the database
	hashedPassword := string(hash)

	return hashedPassword

}

func decryption(providedPassword string, hashedPassword string) bool {
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(providedPassword))
	if err == nil {
		return true
	} else if err == bcrypt.ErrMismatchedHashAndPassword {
		return false
	} else {
		log.Fatal(err)
		return false
	}
}

func AccountExists(user string) (bool, error) {
	query := "SELECT COUNT(*) FROM accounts WHERE Email = ? OR Username = ?"
	row := Accountsdb.QueryRow(query, user, user)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func SessionID(ID int, sessionID string, expt time.Time) error {
	updateQuery := "UPDATE accounts SET SessionID = ?, Expiration = ? WHERE id = ?"
	_, err := Accountsdb.Exec(updateQuery, sessionID, expt, ID)
	if err != nil {
		return err
	}
	return nil
}

func GetAccountID(user string) (int, error) {
	// Prepare the SQL statement
	query := "SELECT id FROM accounts WHERE Email = ? OR Username = ?"
	stmt, err := Accountsdb.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	// Execute the query and get the account ID
	var accountID int
	err = stmt.QueryRow(user, user).Scan(&accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Account does not exist
			return 0, nil
		}
		return 0, err
	}

	return accountID, nil
}
