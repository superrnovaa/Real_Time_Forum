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

type HomeTemplateData struct {
	Posts      []Post
	UserId     int
	Username   string
	SessionId  string
	ProfileImg string
	LikedPosts []string
}

var (
	user       Account
	Likedposts []string
	HomeData   HomeTemplateData
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/HomePage" {
		http.Error(w, "Page Not Found", 404)
		return
	}

	// Query the user data based on the ID
	user = GetUserData(w, r)
	var request Request
	var All []string
	if r.Method == "POST" {
		if r.Header.Get("Content-Type") == "application/json" {
			err := json.NewDecoder(r.Body).Decode(&request)
			if err != nil {
				// Handle other decoding errors
				http.Error(w, "Failed to decode JSON request", http.StatusBadRequest)
				return
			}
			if request.RequestType == "Data" || request.RequestType == "" {
				ConstructPage(user.Id, w, r, All)
				UpdateUsers([]int{user.Id})
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(HomeData)
			} else {
				categories := strings.Split(request.RequestType, " ")
				ConstructPage(user.Id, w, r, categories)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(HomeData)
			}

		}
	} else {
		ConstructPage(user.Id, w, r, All)
		Template(w)
		
	
	}
}

func GetHomeData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HomeData)
}

func ConstructPage(Id int, w http.ResponseWriter, r *http.Request, categories []string) {
	Posts = nil
	if len(categories) == 0 {
		Posts, err = fetchPostsFromDB(false, "", categories)
	} else {
		Posts, err = fetchPostsFromDB(true, "Categories", categories)
	}

	if err != nil {
		log.Fatal(err)
	}

	// Reverse Posts from new to old
	reverseArray(Posts)
	_, _, Likedposts = GetLikedPosts(Id)

	sessionCookie, _ := r.Cookie("session_id")

	HomeData = HomeTemplateData{
		Posts:      Posts,
		UserId:     Id,
		Username:   user.Username,
		SessionId:  sessionCookie.Value,
		ProfileImg: user.ProfileImg,
		LikedPosts: Likedposts,

	}

}

func fetchPostsFromDB(Filtered bool, Type string, Filter []string) ([]Post, error) {
	Posts = nil
	// Query the database for post data
	var rows *sql.Rows
	if !Filtered {
		rows, err = Postsdb.Query("SELECT post_id, user_id, username,userImg, title, content , image, category, like, dislike, timestamp FROM posts")
		if err != nil {
			return nil, fmt.Errorf("failed to query database: %v", err)
		}
	} else {
		if Type == "ProfileFilter" {
			filterString := strings.Join(Filter, ",")
			rows, err = Postsdb.Query("SELECT post_id, user_id, username, userImg, title, content, image, category, like, dislike, timestamp FROM posts WHERE post_id IN (" + filterString + ")")
			if err != nil {
				return nil, fmt.Errorf("failed to query database: %v", err)
			}
		} else if Type == "Categories" {

			placeholders := make([]string, len(Filter))
			args := make([]interface{}, len(Filter))

			for i, category := range Filter {
				placeholders[i] = "category LIKE ?"
				args[i] = "%" + category + "%"
			}

			splitString := strings.Split(Filter[0], " ")
			f := splitString[0]
			if f == "News" {
				args[0] = "%News%"
			}

			placeholderString := strings.Join(placeholders, " OR ")

			query := fmt.Sprintf(`SELECT post_id, user_id, username, userImg, title, content, image, category, like, dislike, timestamp 
    FROM posts
    WHERE %s`, placeholderString)

			// Execute the query
			rows, err = Postsdb.Query(query, args...)
			if err != nil {
				log.Println("Error querying the database:", err)
			}

		}
	}
	defer rows.Close()
	// Iterate over the rows and populate the Posts slice
	for rows.Next() {
		var post Post
		var category string

		err := rows.Scan(&post.PostId, &post.UserId, &post.UserName, &post.UserImg, &post.Title, &post.Content, &post.Image, &category, &post.Like, &post.Dislike, &post.Time)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}

		// Split the category string into a slice of categories
		post.Category = strings.Split(category, " ")

		// Append the post to the Posts slice
		Posts = append(Posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return Posts, nil
}

func HandleLikeRequest(w http.ResponseWriter, r *http.Request, user Account, request Request) {
	// Get the current like value
	var currentLike int
	err = Postsdb.QueryRow("SELECT "+request.Type+" FROM posts WHERE post_id = ?", request.ID).Scan(&currentLike)
	if err != nil {
		log.Fatal(err)
	}
	if request.Checked {
		// Execute the update statement
		_, err := Postsdb.Exec("UPDATE posts SET "+request.Type+" = ? WHERE post_id = ?", currentLike+1, request.ID)
		if err != nil {
			log.Fatal(err)
		}

		_, err = Postsdb.Exec(`INSERT INTO LikedPosts (user_id, post_id, type) VALUES (?, ?, ?)`, user.Id, request.ID, request.Type)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// Handle unchecked state
		_, err = Postsdb.Exec("UPDATE posts SET "+request.Type+" = ? WHERE post_id = ?", currentLike-1, request.ID)
		if err != nil {
			log.Fatal(err)
		}

		_, err = Postsdb.Exec(`DELETE FROM LikedPosts WHERE user_id = ? AND post_id = ? AND type = ?`, user.Id, request.ID, request.Type)
	}
}

func GetLikedPosts(id int) ([]string, []string, []string) {
	var LikedPosts []string
	var LikedPostsId []string
	var DislikedPostsId []string

	// Query the database for post data
	rows, err := Postsdb.Query("SELECT post_id, type FROM LikedPosts WHERE user_id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate over the rows and populate the Posts slice
	for rows.Next() {
		var PostId int
		var Type string

		err := rows.Scan(&PostId, &Type)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		if Type == "Like" {
			LikedPosts = append(LikedPosts, "Like_"+strconv.Itoa(PostId))
			LikedPostsId = append(LikedPostsId, strconv.Itoa(PostId))

		} else {
			LikedPosts = append(LikedPosts, "Dislike_"+strconv.Itoa(PostId))
			DislikedPostsId = append(DislikedPostsId, strconv.Itoa(PostId))
		}
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return LikedPostsId, DislikedPostsId, LikedPosts
}

// To reverse the posts array
func reverseArray(arr []Post) {
	length := len(arr)
	for i := 0; i < length/2; i++ {
		arr[i], arr[length-i-1] = arr[length-i-1], arr[i]
	}
}
