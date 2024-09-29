package forum

import (
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"strconv"
	"strings"
)

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var account SignupData

		err := json.NewDecoder(r.Body).Decode(&account)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var ID int
		errorMessage := ""

		exists, err := AccountExists(account.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ageMessage := AgeValidation(account.Age)

		if !isValidEmail(account.Email) {
			errorMessage = "Invalid Email Format"
		} else if ageMessage != ""{
			errorMessage = ageMessage
		} else if !isValidString(account.Email) || !isValidString(account.Username) || !isValidString(account.Password) {
			errorMessage = "Unallowed Charachters"
		} else if exists {
			errorMessage = "Email already exists! Go to the Login page"
			http.Redirect(w, r, "/", http.StatusFound)
		} else if account.ConfirmPassword != account.Password {
			errorMessage = "Password confirmation does not match the account password."
		} else {
			password := encrytion(account.Password)
			ID, errorMessage, err = AddAccount(account.Email, account.Username, account.FirstName, account.LastName, password, account.Age, account.Gender)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		
		if errorMessage == "" {
			w.WriteHeader(http.StatusOK)
			session := createSession(ID)
			// Set the session cookie in the response
			http.SetCookie(w, &http.Cookie{
				Name:     "session_id",
				Value:    session.ID,
				HttpOnly: true,
				Path:     "/",
				Expires:  session.ExpiresAt,
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
		} else {
			http.Error(w, errorMessage, http.StatusInternalServerError)
		}
	} else {
		Template(w)
	}
}

func AgeValidation(age string) string {
	Age, err := strconv.Atoi(age) 
	if err !=nil {
		return "Please Enter Valid Age"
	} else if Age < 18 {
		return "You have to be 18 or older to Sign Up"
	}
	return ""
}

func AddAccount(email, username, firstname, lastname, password, age, gender string) (int, string, error) {
	errorMessage := ""
	insertQuery := "INSERT INTO accounts (Email, Username, FirstName, LastName, Password, Age, Gender) VALUES (?, ?, ?, ?, ?, ?, ?)"
	result, err := Accountsdb.Exec(insertQuery, email, username, firstname, lastname, password, age, gender)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: accounts.Username") {
			errorMessage = "Username is already taken. Please choose a different username."
			return 0, errorMessage, nil
		} else {
			errorMessage = "An error occurred while creating the account. Please try again later."
			return 0, errorMessage, nil
		}
	}
	Id, err := result.LastInsertId()
	if err != nil {
		return 0, errorMessage, err
	}
	return int(Id), errorMessage, nil
}