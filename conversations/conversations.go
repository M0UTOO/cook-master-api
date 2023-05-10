package conversations

import (
	"cook-master-api/token"
	"cook-master-api/utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Conversations struct {
	UserId int `json:"conversationid"`
}

type Talks struct {
	UserID int `json:"userid"`
	ConversationID int `json:"conversationid"`
}

type Messages struct {
	MessageID int `json:"messageid"`
	ConversationID int `json:"conversationid"`
	Content string `json:"content"`
	IsFromUser1 bool `json:"isfromuser1"`
}
	
func PostConversations(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		tokenHeader := c.Request.Header["Token"]
		if tokenHeader == nil{
			c.JSON(498, gin.H{
				"error": true,
				"message": "missing token",
			})
		}

		err := token.CheckAPIToken(tokenAPI, tokenHeader[0], c)
		if err != nil {
			c.JSON(498, gin.H{
				"error": true,
				"message": "wrong token",
			})
			return
		}

		type ConversationsReq struct {
			UserId int `json:"userid"`
		}

		var req ConversationsReq

		err = c.BindJSON(&req)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "can't decode json request",
			})
			return
		}

		db, err := sql.Open("mysql", "admin:Respons11@tcp(localhost:3306)/cookmaster")
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		result, err := db.Exec("INSERT INTO CONVERSATIONS VALUES (NULL)")
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot insert conversation",
			})
			return
		}

		conversationId, err := result.LastInsertId()
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot get conversation id",
			})
			return
		}

		talk := Talks{
			UserID:         req.UserId,
			ConversationID: int(conversationId),
		}

		_, err = db.Exec("INSERT INTO TALKS (Id_USERS, Id_CONVERSATIONS) VALUES (?, ?)", talk.UserID, talk.ConversationID)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot insert talk",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
			"message": "conversation created",
		})
		return
	}
}

func DeleteConversations(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		tokenHeader := c.Request.Header["Token"]
		if tokenHeader == nil{
			c.JSON(498, gin.H{
				"error": true,
				"message": "missing token",
			})
		}

		err := token.CheckAPIToken(tokenAPI, tokenHeader[0], c)
		if err != nil {
			c.JSON(498, gin.H{
				"error": true,
				"message": "wrong token",
			})
			return
		}

		var req Conversations

		err = c.BindJSON(&req)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "can't decode json request",
			})
			return
		}

		db, err := sql.Open("mysql", "admin:Respons11@tcp(localhost:3306)/cookmaster")
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		_, err = db.Exec("DELETE FROM TALKS WHERE Id_CONVERSATIONS = ?", req.UserId)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot delete talk",
			})
			return
		}

		_, err = db.Exec("DELETE FROM CONVERSATIONS WHERE Id_CONVERSATIONS = ?", req.UserId)
		fmt.Println(err)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot delete conversation",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
			"message": "conversation deleted",
		})
		return
	}
}

func GetConversations(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		tokenHeader := c.Request.Header["Token"]
		if tokenHeader == nil{
			c.JSON(498, gin.H{
				"error": true,
				"message": "missing token",
			})
		}

		err := token.CheckAPIToken(tokenAPI, tokenHeader[0], c)
		if err != nil {
			c.JSON(498, gin.H{
				"error": true,
				"message": "wrong token",
			})
			return
		}

		db, err := sql.Open("mysql", "admin:Respons11@tcp(localhost:3306)/cookmaster")
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		rows, err := db.Query("SELECT * FROM CONVERSATIONS")
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot get conversations",
			})
			return
		}

		var conversations []Conversations

		for rows.Next() {
			var conversation Conversations
			err := rows.Scan(&conversation.UserId)
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "cannot scan conversations",
				})
				return
			}
			conversations = append(conversations, conversation)
		}

		c.JSON(200, conversations)
		return
	}
}

func GetConversationByID(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		tokenHeader := c.Request.Header["Token"]
		if tokenHeader == nil{
			c.JSON(498, gin.H{
				"error": true,
				"message": "missing token",
			})
		}

		err := token.CheckAPIToken(tokenAPI, tokenHeader[0], c)
		if err != nil {
			c.JSON(498, gin.H{
				"error": true,
				"message": "wrong token",
			})
			return
		}

		db, err := sql.Open("mysql", "admin:Respons11@tcp(localhost:3306)/cookmaster")
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		id := c.Param("id")
		if id == "" {
			c.JSON(400, gin.H{
				"error": true,
				"message": "id can't be empty",
			})
			return
		}

		if !utils.IsSafeString(id) {
			c.JSON(400, gin.H{
				"error": true,
				"message": "id can't contain sql injection",
			})
			return
		}

		rows, err := db.Query("SELECT * FROM MESSAGES WHERE Id_CONVERSATIONS = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot get messages",
			})
			return
		}

		var messages []Messages

		for rows.Next() {
			var message Messages
			err := rows.Scan(&message.MessageID, &message.ConversationID, &message.Content, &message.IsFromUser1)
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "cannot scan messages",
				})
				return
			}
			messages = append(messages, message)
		}

		c.JSON(200, messages)
		return
	}
}

func GetConversationForUserID(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		tokenHeader := c.Request.Header["Token"]
		if tokenHeader == nil{
			c.JSON(498, gin.H{
				"error": true,
				"message": "missing token",
			})
		}

		err := token.CheckAPIToken(tokenAPI, tokenHeader[0], c)
		if err != nil {
			c.JSON(498, gin.H{
				"error": true,
				"message": "wrong token",
			})
			return
		}

		db, err := sql.Open("mysql", "admin:Respons11@tcp(localhost:3306)/cookmaster")
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		id := c.Param("id")
		if id == "" {
			c.JSON(400, gin.H{
				"error": true,
				"message": "id can't be empty",
			})
			return
		}

		if !utils.IsSafeString(id) {
			c.JSON(400, gin.H{
				"error": true,
				"message": "id can't contain sql injection",
			})
			return
		}

		rows, err := db.Query("SELECT C.* FROM CONVERSATIONS C INNER JOIN TALKS T ON C.Id_CONVERSATIONS = T.Id_CONVERSATIONS WHERE T.Id_USERS = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot get conversations",
			})
			return
		}

		var conversations []Conversations

		for rows.Next() {
			var conversation Conversations
			err := rows.Scan(&conversation.UserId)
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "cannot scan conversations",
				})
				return
			}
			conversations = append(conversations, conversation)
		}

		c.JSON(200, conversations)
		return
	}
}

func PostMessage(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		tokenHeader := c.Request.Header["Token"]
		if tokenHeader == nil{
			c.JSON(498, gin.H{
				"error": true,
				"message": "missing token",
			})
		}

		err := token.CheckAPIToken(tokenAPI, tokenHeader[0], c)
		if err != nil {
			c.JSON(498, gin.H{
				"error": true,
				"message": "wrong token",
			})
			return
		}

		db, err := sql.Open("mysql", "admin:Respons11@tcp(localhost:3306)/cookmaster")
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		var message Messages

		err = c.BindJSON(&message)
		if err != nil {
			c.JSON(400, gin.H{
				"error": true,
				"message": "cannot bind message",
			})
			return
		}

		_, err = db.Exec("SELECT * FROM CONVERSATIONS WHERE Id_CONVERSATIONS = ?", message.ConversationID)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "wrong conversation id",
			})
			return
		}

		if message.Content == "" {
			c.JSON(400, gin.H{
				"error": true,
				"message": "content can't be empty",
			})
			return
		}

		if !utils.IsSafeString(message.Content) {
			c.JSON(400, gin.H{
				"error": true,
				"message": "content can't contain sql injection",
			})
			return
		}

		_, err = db.Exec("INSERT INTO MESSAGES (Id_CONVERSATIONS, Content, IsFromUser1) VALUES (?, ?, ?)", message.ConversationID, message.Content, message.IsFromUser1)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot insert message",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
			"message": "message inserted",
		})
		return
	}
}


		



