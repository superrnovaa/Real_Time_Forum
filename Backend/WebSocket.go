package forum

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var userSockets = make(map[int]*websocket.Conn)

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "No active session", http.StatusUnauthorized)
			return
		}

		userID := getUserIDFromSessionID(cookie.Value)

		sender, receiver, err := getChatInfoFromURL(r)
		if err != nil {
			http.Error(w, "Page Not Found", 404)
			return
		}
		n1, err1 := accountInfo(sender, receiver, userID)

		if err1 != nil || n1 == false {
			http.Error(w, "Page Not Found", 404)
		} else {
			Template(w)
		}
	}
}

func accountInfo(sender, receiver, userID int) (bool, error) {
	var senderCount, receiverCount int

	// Check if sender account exists
	query := "SELECT COUNT(*) FROM accounts WHERE id = ?"
	err := Accountsdb.QueryRow(query, sender).Scan(&senderCount)
	if err != nil {
		return false, err
	}
	if senderCount == 0 {
		return false, fmt.Errorf("sender account not found")
	}

	// Check if receiver account exists
	query = "SELECT COUNT(*) FROM accounts WHERE id = ?"
	err = Accountsdb.QueryRow(query, receiver).Scan(&receiverCount)
	if err != nil {
		return false, err
	}
	if receiverCount == 0 {
		return false, fmt.Errorf("receiver account not found")
	}

	// Check if the receiver is the same as the user ID or if the sender is not the user ID
	if receiver == userID || sender != userID {
		return false, nil
	}

	return true, nil
}

func getChatInfoFromURL(r *http.Request) (int, int, error) {
	// Parse the query parameters
	query := r.URL.Query()

	// Extract the sender and receiver IDs
	senderIDStr := query.Get("Sender")
	receiverIDStr := query.Get("Receiver")
	// Check if the sender or receiver ID is missing
	if senderIDStr == "" || receiverIDStr == "" {
		return 0, 0, fmt.Errorf("missing sender or receiver ID in the URL")
	}

	// Convert the IDs to integers
	senderID, err := strconv.Atoi(senderIDStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid sender ID: %v", err)
	}

	receiverID, err := strconv.Atoi(receiverIDStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid receiver ID: %v", err)
	}

	return senderID, receiverID, nil
}

func WsConn(w http.ResponseWriter, r *http.Request) {
	// Return true to allow connection from any origin
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// Upgrade the HTTP connection to  a WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "No active session", http.StatusUnauthorized)
		return
	}

	userID := getUserIDFromSessionID(cookie.Value)
	if userID != 0 {
		userSockets[userID] = conn
		UpdateUsers([]int{userID})
	}

	defer func() {
		deleteConnection(conn)
		conn.Close()
	}()
	log.Println("Client successfully connected")

	for {
		var wsRequest WSRequest
		err := conn.ReadJSON(&wsRequest)
		if err != nil {
			log.Println(err)
			return
		}
		currentTime := time.Now()
		if wsRequest.Type == "send_message" {
			var content SendMessageContent
			contentBytes, err := json.Marshal(wsRequest.Content)
			if err != nil {
				log.Println("Error marshaling content to JSON:", err)
				// Additional error handling or response
				return
			}
			err = json.Unmarshal(contentBytes, &content)
			if err != nil {
				log.Println("failed to unmarshal wsRequest.Content:", err)
				return
			}
			Id, err := AddChat(content.Text, content.Sender, content.Receiver, currentTime)
			if err != nil {
				log.Println(err)
			}

			response := WSRequest{
				Type: "send_message",
				Content: struct {
					ID        int64     `json:"ID"`
					Receiver  int       `json:"Receiver"`
					Sender    int       `json:"Sender"`
					Text      string    `json:"Text"`
					Timestamp time.Time `json:"Time"`
				}{
					ID:        Id,
					Receiver:  content.Receiver,
					Sender:    content.Sender,
					Text:      content.Text,
					Timestamp: currentTime,
				},
			}

			// Write the JSON-encoded chat messages to the WebSocket
			err = conn.WriteJSON(response)
			if err != nil {
				log.Println("Error writing to WebSocket:", err)
				return
			}

			recipientConn, ok := userSockets[content.Receiver]
			if ok {
				err = recipientConn.WriteJSON(response)
				if err != nil {
					log.Println("Error writing to WebSocket:", err)
					return
				}
			}
			UpdateUsers([]int{content.Sender, content.Receiver})
		} else if wsRequest.Type == "notification" {
			var content NotificationContent
			contentBytes, err := json.Marshal(wsRequest.Content)
			if err != nil {
				log.Println("Error marshaling content to JSON:", err)
				return
			}
			err = json.Unmarshal(contentBytes, &content)
			if err != nil {
				log.Println("Error unmarshaling JSON:", err)
				log.Println("JSON data:", string(contentBytes))
				return
			}
			UpdateNotification(content.Sender, content.Receiver)
			UpdateUsers([]int{content.Sender, content.Receiver})
		} else if wsRequest.Type == "load_messages" {
			var content LoadMessagesContent
			contentBytes, err := json.Marshal(wsRequest.Content)
			if err != nil {
				log.Println("Error marshaling content to JSON:", err)
				return
			}
			err = json.Unmarshal(contentBytes, &content)
			if err != nil {
				log.Println("Error unmarshaling JSON:", err)
				log.Println("JSON data:", string(contentBytes))
				return
			}
			chatMessages, err := LoadMessages(content.Sender, content.Receiver, content.Scroll)
			if err != nil {
				log.Println(err)
			}
			// Convert the []Com slice to []interface{}

			var messages []interface{}
			for _, msg := range chatMessages {
				messages = append(messages, msg)
			}

			var rInfo RInfo

			err = Accountsdb.QueryRow("SELECT Username, UserImg FROM accounts WHERE id = ?", content.Receiver).Scan(&rInfo.RUname, &rInfo.RPImg)
			if err != nil {
				log.Fatal(err)
			}
			var response WSResponse
			response.Type = "load_messages"
			response.RInfo = rInfo
			response.Data = messages

			// Write the JSON-encoded chat messages to the WebSocket
			err = conn.WriteJSON(response)
			if err != nil {
				log.Println("Error writing to WebSocket:", err)
				return
			}
		} else if wsRequest.Type == "typing" {
			var content TypingEvent
			contentBytes, err := json.Marshal(wsRequest.Content)
			if err != nil {
				log.Println("Error marshaling content to JSON:", err)
				return
			}
			err = json.Unmarshal(contentBytes, &content)
			if err != nil {
				log.Println("failed to unmarshal wsRequest.Content:", err)
				return
			}
			recipientConn, ok := userSockets[content.Receiver]
			if ok {
				response := WSRequest{
					Type: "typing",
					Content: struct {
						Receiver int    `json:"rceiver"`
						Sender   int    `json:"sender"`
						Type     string `json:"type"`
					}{
						Receiver: content.Receiver,
						Sender:   content.Sender,
						Type:     content.Type,
					},
				}
				err = recipientConn.WriteJSON(response)
				if err != nil {
					log.Println("Error writing to WebSocket:", err)
					return
				}
			}
		}

	}
}

func deleteConnection(conn *websocket.Conn) {
	for userID, userConn := range userSockets {
		if userConn == conn {
			delete(userSockets, userID)
			return
		}
	}
}

func UpdateUsers(Users []int) {
	var activeUsers []int
	for id := range userSockets {
		activeUsers = append(activeUsers, id)
	}

	for _, id := range Users {
		OtherUsers := GetOtherUsersData(id, activeUsers)
		socket, ok := userSockets[id]
		fmt.Println()
		if ok {
			err := socket.WriteJSON(map[string]interface{}{
				"type":    "activeUsers",
				"Content": OtherUsers,
			})
			if err != nil {
				// Handle the error, e.g., log it or close the connection
				log.Printf("Update Users: Error sending active users to user: %v", err)
			}
		}
	}
}

func UpdateOnlineUsers() {
	var activeUsers []int
	for id := range userSockets {
		activeUsers = append(activeUsers, id)
	}

	for id, socket := range userSockets {
		OtherUsers := GetOtherUsersData(id, activeUsers)
		// Send the active users data to each connected socket
		err := socket.WriteJSON(map[string]interface{}{
			"type":    "activeUsers",
			"Content": OtherUsers,
		})
		if err != nil {
			log.Printf("Update Online Users:Error sending active users to user: %v", err)
		}
	}
}

func GetOtherUsersData(UserId int, activeUsers []int) []User {
	var Friends []User
	var Stranger []User

	// Query the database for post data
	rows, err := Accountsdb.Query("SELECT id, Username, UserImg FROM accounts WHERE id != ?", UserId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate over the rows and populate the Posts slice
	for rows.Next() {
		var User User
		err := rows.Scan(&User.UserId, &User.Username, &User.ProfileImg)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		User.MessageID, User.LastMessage = GetLastMessage(UserId, User.UserId)
		if IsUserActive(User.UserId, activeUsers) {
			User.Status = "Online"
		} else {
			User.Status = "Offline"
		}
		User.Notifications = GetNumOfNotifications(User.UserId, UserId)

		if User.LastMessage == "" {
			Stranger = append(Stranger, User)
		} else {
			Friends = append(Friends, User)
		}
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	Users := CombineUsers(Friends, Stranger)

	return Users
}

func GetLastMessage(senderId int, receiverId int) (int, string) {
	var id int
	var message string

	err := Accountsdb.QueryRow(`
    SELECT id, text FROM chat
    WHERE (Sender = ? AND Receiver = ?) OR (Sender = ? AND Receiver = ?)
    ORDER BY id DESC
    LIMIT 1`,
		senderId, receiverId, receiverId, senderId,
	).Scan(&id, &message)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ""
		} else {
			log.Fatal(err)
		}
	}
	if len(message) > 15 {
		message = message[:20] + "..."
	}

	return id, message
}

func IsUserActive(userId int, activeUsers []int) bool {
	for _, id := range activeUsers {
		if id == userId {
			return true
		}
	}

	return false
}

func CombineUsers(friends []User, strangers []User) []User {
	// Sort friends by descending MessageID
	sort.Slice(friends, func(i, j int) bool {
		return friends[i].MessageID > friends[j].MessageID
	})

	// Sort strangers alphabetically by Username
	sort.Slice(strangers, func(i, j int) bool {
		return strings.ToLower(strangers[i].Username) < strings.ToLower(strangers[j].Username)
	})

	// Combine slices
	allUsers := make([]User, len(friends)+len(strangers))
	copy(allUsers[:len(friends)], friends)
	copy(allUsers[len(friends):], strangers)

	return allUsers
}

func LoadMessages(sender, receiver, scroll int) ([]Com, error) {
	var chatMessages []Com
	CHAT := `
        SELECT *
        FROM (
          SELECT * FROM chat WHERE (Sender = ? AND Receiver = ?) OR (Receiver = ? AND Sender = ?)
          ORDER BY id DESC
          LIMIT ?
        ) sub
        ORDER BY id ASC;
    `

	rows, err := Accountsdb.Query(CHAT, sender, receiver, sender, receiver, scroll+10)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var chat Com
		if err := rows.Scan(&chat.Id, &chat.Sender, &chat.Receiver, &chat.Text, &chat.Time); err != nil {
			return nil, err
		}
		chatMessages = append(chatMessages, chat)
	}

	return chatMessages, nil
}

func AddChat(text string, sender int, receiver int, time time.Time) (int64, error) {
	insertQuery := "INSERT INTO chat (Sender, Receiver, Text, Timestamp) VALUES (?, ?, ?, ?)"

	result, err := Accountsdb.Exec(insertQuery, sender, receiver, text, time)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
