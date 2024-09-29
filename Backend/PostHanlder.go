package forum

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	user = GetUserData(w, r)
	if r.Method == "GET" {
		Template(w)
	} else {
		var request Request
		// Check if the Content-Type is application/json
		if r.Header.Get("Content-Type") == "application/json" {
			err := json.NewDecoder(r.Body).Decode(&request)
			if err != nil {
				// Handle other decoding errors
				http.Error(w, "Failed to decode JSON request", http.StatusBadRequest)
				return
			}
			switch request.RequestType {
			case "Data":
				HandleDataRequest(w, r, user, request)
			case "like":
				HandleLikeRequest(w, r, user, request)
			case "delete":
				HandleDeleteRequest(request)
			case "DeleteComment":
				HandleDeleteCommentRequest(request)
			case "UpdateComment":
				HandleUpdateCommentRequest(request)
			default:
				http.Error(w, "Invalid request type", http.StatusBadRequest)
				return
			}
		}
	}
}

func HandleDataRequest(w http.ResponseWriter, r *http.Request, user Account, request Request) {
	PostData.Username = user.Username
	PostData.ProfileImg = user.ProfileImg
	postId, err := getPostIDFromURL(r)
	if err != nil {
		http.Error(w, "Page Not Found", 404)
		return
	}
	err = PostInf(postId)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Page Not Found", 404)
			return
		} else {
			log.Fatal(err)
		}
	}
	isLiked(user.Id, postId)
	isCreated(user.Id, postId)
	PostData.LikedComments = GetLikedComments(user.Id, postId)
	PostData.CreatedComments = GetCreatedComments(user.Id, postId)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(PostData)
}

func GetPostData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(PostData)
}

func PostInf(postId int) (err error) {
	var category string
	query := `SELECT post_id, user_id, username, userImg, title, content, image, category, like, dislike, timestamp 
            FROM posts 
            WHERE post_id = $1`

	err = Postsdb.QueryRow(query, postId).Scan(
		&post.PostId,
		&post.UserId,
		&post.UserName,
		&post.UserImg,
		&post.Title,
		&post.Content,
		&post.Image,
		&category,
		&post.Like,
		&post.Dislike,
		&post.Time,
	)
	if err != nil {
		return err
	}

	post.Category = strings.Split(category, " ")
	var commentArray []Comment
	var comment Comment
	crow, err := Commentsdb.Query("SELECT comment_id, userName, text, ProfileImage, timestamp, CLike, CDislike FROM comments WHERE post_id = ? ", post.PostId)
	if err != nil {
		fmt.Errorf("failed to query database: %v", err)
	}
	defer crow.Close()
	for crow.Next() {
		err := crow.Scan(&comment.CommentId, &comment.UserName, &comment.Text, &comment.ProfileImage, &comment.Time, &comment.CLike, &comment.CDislike)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		commentArray = append(commentArray, comment)
	}
	post.Comments = commentArray
	PostData.Post = post
	return nil
}

func getPostIDFromURL(r *http.Request) (int, error) {
	path := r.URL.Path
	// Split path on "/"
	segments := strings.Split(path, "/")
	if len(segments) > 3 {
		return 0, err
	}
	// Last segment is the ID
	idString := segments[len(segments)-1]
	// Convert ID to int
	id, err := strconv.Atoi(idString)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func isLiked(userId int, postId int) {
	var Type string

	row := Postsdb.QueryRow("SELECT type FROM LikedPosts WHERE user_id=? AND post_id=?", userId, postId)

	if err = row.Scan(&Type); err != nil {
		if err == sql.ErrNoRows {
			PostData.IsLiked = false
			PostData.IsDisliked = false
		} else {
			log.Fatal(err)
		}
	} else {
		if Type == "Like" {
			PostData.IsLiked = true
			PostData.IsDisliked = false
		} else if Type == "Dislike" {
			PostData.IsLiked = false
			PostData.IsDisliked = true
		}
	}
}

func isCreated(userId int, postId int) {
	var UserID int
	row := Postsdb.QueryRow("SELECT user_id From posts WHERE post_id=?", postId)

	if err = row.Scan(&UserID); err != nil {
		log.Fatal(err)
	} else {
		if UserID == userId {
			PostData.IsCreated = true
		} else {
			PostData.IsCreated = false
		}
	}
}

func HandleDeleteRequest(request Request) {
	PostId := request.ID
	_, err = Postsdb.Exec("DELETE FROM posts WHERE post_id = ?", PostId)
	if err != nil {
		log.Fatal(err)
	}
	_, err = Postsdb.Exec("DELETE FROM LikedPosts WHERE post_id = ?", PostId)
	if err != nil {
		log.Fatal(err)
	}
}
