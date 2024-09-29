package forum

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/CreatePost" {
		http.Error(w, "Page Not Found" , 404)
		return
	}

	var user Account
	user = GetUserData(w, r)
	if r.Method == "POST" {
		// Parse the form data
		err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

		// Get the uploaded image file
		file, handler, err := r.FormFile("file")
		if err != nil && err != http.ErrMissingFile {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var imgPath string
		if file != nil {
			defer file.Close()
			// Save the image file
			imgPath = "Posts/" + handler.Filename
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

		// Get the Title and Content value
		title := r.FormValue("Title")
		content := r.FormValue("Content")

		// Get the selected checkboxes
		checkboxes := r.Form["checkbox"]
		// Now you can use the retrieved id and username variables as needed
		InsertPost(user.Id, user.Username, user.ProfileImg, title, content, imgPath, strings.Join(checkboxes, ","), 0, 0)
		//http.Redirect(w, r, "/HomePage", http.StatusFound)
	} else {
		HomeData = HomeTemplateData{
			Username:        user.Username,
			ProfileImg:      user.ProfileImg,
		}
		Template(w)
	}
}

func InsertPost(userID int, username string, userImg string, title string, content string, image string, category string, like int, dislike int) {
	query := `
        INSERT INTO posts (user_id, username,userImg,title,content, image,category, like, dislike)
        VALUES (?, ?, ?, ?, ?,?,?,?,?)
    `

	_, err := Postsdb.Exec(query, userID, username, userImg, title, content, image, category, like, dislike)
	if err != nil {
		log.Fatal(err)
	}
}
