package forum

import (
	"html/template"
	"io"
	"log"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func Template(w http.ResponseWriter) {
	tmpl, err := template.ParseFiles("Pages/forum.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}
	if err = tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}

func isValidEmail(email string) bool {
	// Regular expression pattern for email validation
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Create a regular expression object
	r := regexp.MustCompile(pattern)

	// Match the email against the pattern
	return r.MatchString(email)
}

func isValidString(str string) bool {
	match, _ := regexp.MatchString("^[a-zA-Z0-9!@#$%^&*()-_=+,.?/:;{}|~<>]+$", str)
	return match
}

func generateSessionID() string {
	// Generate a unique session ID using UUID
	sessionID := uuid.New().String()
	return sessionID
}

func GetUserData(w http.ResponseWriter, r *http.Request) Account {
	var user Account
	sessionCookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		// No session cookie found, redirect to login page
		http.Redirect(w, r, "/", http.StatusFound)
	} else if err != nil {
		// Handle other cookie-related errors
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	//Get id
	user.Id = getUserIDFromSessionID(sessionCookie.Value)
	//err = Accountsdb.QueryRow("SELECT id FROM accounts WHERE SessionID = ?", sessionCookie.Value).Scan(&user.Id)
	// Query the user data based on the ID
	err = Accountsdb.QueryRow("SELECT Email, Username, UserImg FROM accounts WHERE id = ?", user.Id).
		Scan(&user.Email, &user.Username, &user.ProfileImg)
	if err != nil {
		log.Fatal(err)
	}
	return user
}


func getIDBySessionID(sessionid string) (int, error) {
	// Prepare the SQL query
	query := "SELECT id FROM accounts WHERE SessionID = ?"

	// Execute the query and retrieve the result
	row := Accountsdb.QueryRow(query, sessionid)

	// Scan the result into a string variable
	ID := 0
	err := row.Scan(&ID)
	if err != nil {
		return 0, err
	}

	return ID, nil
}

func ProfileImageHandler(w http.ResponseWriter, r *http.Request) {
	user := GetUserData(w, r)
	// Parse the multipart form data
	err := r.ParseMultipartForm(32 << 20) // Set the maximum file size
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}
	file, handler, err := r.FormFile("PImg") // Get the uploaded file
	if err != nil {
		http.Error(w, "Failed to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	var imgPath string
	if file != nil {
		// Save the image file
		imgPath = "ProfileImages/" + handler.Filename
		distinationFile, err := os.OpenFile(imgPath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer distinationFile.Close()
		io.Copy(distinationFile, file)
	} else {
		imgPath = ""
	}
	_, err = Accountsdb.Exec(`UPDATE accounts SET UserImg = ? WHERE id= ? `, handler.Filename, user.Id)
	if err != nil {
		log.Fatal(err)
	}

	_, err = Commentsdb.Exec(`UPDATE comments SET ProfileImage = ? WHERE user_id= ? `, handler.Filename, user.Id)
	if err != nil {
		log.Fatal(err)
	}
	_, err = Postsdb.Exec(`UPDATE posts SET userImg = ? WHERE user_id= ? `, handler.Filename, user.Id)
	if err != nil {
		log.Fatal(err)
	}
	// Send a response back to the client
	fmt.Fprintf(w, "File uploaded successfully.")

}
