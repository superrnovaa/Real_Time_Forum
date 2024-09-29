package forum

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type Payload struct {
	Message      string `json:"message"`
	PostId       int    `json:"postId"`
	Username     string `json:"username"`
	ProfileImage string `json:"profileImage"`
}

type CommentResponse struct {
	CommentId int64 `json:"CommentId"`
}

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed" , http.StatusMethodNotAllowed)
		return
	}

	var payload Payload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid request payload" , http.StatusBadRequest)
		return
	}
	

	message := payload.Message
	postId := payload.PostId
	username := payload.Username
	profileImage := payload.ProfileImage

	// Create a new comment
	comment := Comment{
		UserName:     username,
		Text:         message,
		ProfileImage: profileImage,
	}
	// Insert the comment into the database
	result, err := Commentsdb.Exec("INSERT INTO comments (user_id, userName, post_id, text, ProfileImage) VALUES (?,?, ?, ?,?)",
		user.Id, comment.UserName, postId, comment.Text, profileImage)
	if err != nil {
		log.Fatal(err)
	}

	CommentId, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	// Create a CommentResponse object
	response := CommentResponse{
		CommentId: CommentId,
	}

	// Set the response content type and encode the response as JSON
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	err = enc.Encode(response)
	if err != nil {
		http.Error(w, "Error encoding response" , http.StatusInternalServerError)
		return
	}

}

func CommentLikeHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w,"Method not allowed" , http.StatusMethodNotAllowed)
		return
	}

	var request Request
	//var user Account
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w,"Invalid request payload" , http.StatusBadRequest)
		return
	}

	var user Account
	// Query the user data based on the ID
	user = GetUserData(w, r)
	var currentLike int
	err = Commentsdb.QueryRow("SELECT C"+request.Type+" FROM comments WHERE comment_id = ?", request.ID).Scan(&currentLike)
	if err != nil {
		log.Fatal(err)
	}
	if request.Checked {
		// Execute the update statement
		_, err := Commentsdb.Exec("UPDATE comments SET C"+request.Type+" = ? WHERE comment_id = ?", currentLike+1, request.ID)
		if err != nil {
			log.Fatal(err)
		}

		_, err = Commentsdb.Exec(`INSERT INTO Likedcomments (user_id, comment_id, post_id, type) VALUES (?, ?, ?, ?)`, user.Id, request.ID, request.PostId, request.Type)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// Handle unchecked state
		_, err = Commentsdb.Exec("UPDATE comments SET C"+request.Type+" = ? WHERE comment_id = ?", currentLike-1, request.ID)
		if err != nil {
			log.Fatal(err)
		}

		_, err = Commentsdb.Exec(`DELETE FROM Likedcomments WHERE user_id = ? AND comment_id = ? AND type = ?`, user.Id, request.ID, request.Type)
		if err != nil {
			log.Fatal(err)
		}

	}
}

func GetLikedComments(UserId, PostId int) []string {
	var LikedComments []string
	var commentIds []int

	rows, err := Commentsdb.Query("SELECT comment_id FROM Likedcomments WHERE post_id = ? AND user_id = ?", PostId, UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return LikedComments
		} else {
			log.Fatal(err)
		}
	}

	defer rows.Close()
	for rows.Next() {
		var id int

		if err := rows.Scan(&id); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		commentIds = append(commentIds, id)
	}

	for _, id := range commentIds {
		var Type string
		err := Commentsdb.QueryRow("SELECT type FROM Likedcomments WHERE comment_id = ? AND user_id = ?", id, UserId).Scan(&Type)
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			} else {
				log.Fatal("No rows for liked comment")
			}
		}
		if Type == "Like" {
			LikedComments = append(LikedComments, "CLike_"+strconv.Itoa(id))
		} else {
			LikedComments = append(LikedComments, "CDislike_"+strconv.Itoa(id))
		}
	}

	return LikedComments
}

func GetCreatedComments(UserId int, PostId int) []string {
	var CommentsId []string

	// Query the database for post data
	rows, err := Commentsdb.Query("SELECT comment_id FROM comments WHERE user_id = ? and post_id = ?", UserId, PostId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate over the rows and populate the Posts slice
	for rows.Next() {
		var CommentId int

		err := rows.Scan(&CommentId)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		CommentsId = append(CommentsId, strconv.Itoa(CommentId))
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return CommentsId
}

func HandleDeleteCommentRequest(request Request) {
	CommentId := request.ID
	_, err = Commentsdb.Exec("DELETE FROM comments WHERE comment_id = ?", CommentId)
	if err != nil {
		log.Fatal(err)
	}
	_, err = Commentsdb.Exec("DELETE FROM Likedcomments WHERE comment_id = ?", CommentId)
	if err != nil {
		log.Fatal(err)
	}
}

func HandleUpdateCommentRequest(request Request) {
	CommentId := request.ID
	TextContent := request.Text
	_, err = Commentsdb.Exec("UPDATE comments SET text = ? WHERE comment_id = ?", TextContent, CommentId) // Updated query
	if err != nil {
		log.Fatal(err)
	}
}
