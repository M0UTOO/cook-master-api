package messages

import (
	"cook-master-api/token"
	"cook-master-api/utils"
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Message struct {
	IdMessage      int     `json:"idmessage"`
	Content	  string  `json:"content"`
	CreatedAt 	  string  `json:"createdat"`
	IdSender int  `json:"idsender"`
	IdReceiver 	  int     `json:"idreceiver"`
}

func GetMessageForIdSenderAndIdReceiver(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		tokenHeader := c.Request.Header["Token"]
		if tokenHeader == nil {
			c.JSON(498, gin.H{
				"error":   true,
				"message": "missing token",
			})
		}

		err := token.CheckAPIToken(tokenAPI, tokenHeader[0], c)
		if err != nil {
			c.JSON(498, gin.H{
				"error":   true,
				"message": "wrong token",
			})
			return
		}

		idsender := c.Param("idsender")
		if idsender == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "idsender can't be empty",
			})
			return
		}
		if !utils.IsSafeString(idsender) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "idsender can't contain sql injection",
			})
			return
		}

		idreceiver := c.Param("idreceiver")
		if idreceiver == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "idreceiver can't be empty",
			})
			return
		}
		if !utils.IsSafeString(idreceiver) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "idreceiver can't contain sql injection",
			})
			return
		}

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		rows, err := db.Query("SELECT * FROM MESSAGES WHERE Id_USERS = '" + idsender + "' AND Id_USERS_1 = '" + idreceiver + "' OR (Id_USERS = '" + idreceiver + "' AND Id_USERS_1 = '" + idsender + "') ORDER BY createdAt ASC")
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't query database",
			})
			return
		}

		var messages []Message
		for rows.Next() {
			var message Message
			err = rows.Scan(&message.IdMessage, &message.Content, &message.CreatedAt, &message.IdSender, &message.IdReceiver)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "can't scan database",
				})
				return
			}
			messages = append(messages, message)
		}

		c.JSON(200, messages)
	}
}

func PostMessage(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		tokenHeader := c.Request.Header["Token"]
		if tokenHeader == nil {
			c.JSON(498, gin.H{
				"error":   true,
				"message": "missing token",
			})
		}

		err := token.CheckAPIToken(tokenAPI, tokenHeader[0], c)
		if err != nil {
			c.JSON(498, gin.H{
				"error":   true,
				"message": "wrong token",
			})
			return
		}

		idsender := c.Param("idsender")
		if idsender == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "idsender can't be empty",
			})
			return
		}
		if !utils.IsSafeString(idsender) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "idsender can't contain sql injection",
			})
			return
		}

		idreceiver := c.Param("idreceiver")
		if idreceiver == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "idreceiver can't be empty",
			})
			return
		}
		if !utils.IsSafeString(idreceiver) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "idreceiver can't contain sql injection",
			})
			return
		}

		var message Message
		err = c.BindJSON(&message)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "can't bind json",
			})
			return
		}

		if message.Content == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "missing content",
			})
			return
		}

		if !utils.IsSafeString(message.Content) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "content is not safe",
			})
			return
		}

		if len(message.Content) < 0 {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "content is too long or too short",
			})
			return
		}

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		var usersender int

		err = db.QueryRow("SELECT Id_USERS FROM USERS WHERE Id_USERS = '" + idsender + "'").Scan(&usersender)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "user sender doesn't exist",
			})
			return
		}

		var userreceiver int

		err = db.QueryRow("SELECT Id_USERS FROM USERS WHERE Id_USERS = '" + idreceiver + "'").Scan(&userreceiver)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "user receiver doesn't exist",
			})
			return
		}

		rows, err := db.Exec("INSERT INTO MESSAGES (content, createdAt, Id_USERS, Id_USERS_1) VALUES (?, DEFAULT, ?, ?)", message.Content, idsender, idreceiver)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't insert into database",
			})
			return
		}

		messageid, err := rows.LastInsertId()

		c.JSON(200, gin.H{
			"error":   false,
			"idmessage":  messageid,
			"message": "message added",
		})
	}
}

func GetAllClientForAContractorUserId(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		type User struct {
			IdUser      int     `json:"iduser"`
			FirstName	  string  `json:"firstname"`
			LastName 	  string  `json:"lastname"`
			ProfilePicture 	  string     `json:"profilepicture"`
		}

		tokenHeader := c.Request.Header["Token"]
		if tokenHeader == nil {
			c.JSON(498, gin.H{
				"error":   true,
				"message": "missing token",
			})
		}

		err := token.CheckAPIToken(tokenAPI, tokenHeader[0], c)
		if err != nil {
			c.JSON(498, gin.H{
				"error":   true,
				"message": "wrong token",
			})
			return
		}

		id := c.Param("id")
		if id == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id can't be empty",
			})
			return
		}
		if !utils.IsSafeString(id) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id can't contain sql injection",
			})
			return
		}

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		var idcontractor int

		err = db.QueryRow("SELECT Id_CONTRACTORS FROM CONTRACTORS JOIN CONTRACTOR_TYPES ON CONTRACTOR_TYPES.Id_CONTRACTOR_TYPES = CONTRACTORS.Id_CONTRACTOR_TYPES WHERE Id_USERS = '" + id + "' AND CONTRACTOR_TYPES.name = 'chief'").Scan(&idcontractor)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "user doesn't exist or is not a chief",
			})
			return
		}

		rows, err := db.Query("SELECT DISTINCT USERS.Id_USERS, USERS.firstName, USERS.lastname, USERS.profilePicture FROM USERS JOIN CLIENTS ON USERS.Id_USERS = CLIENTS.Id_USERS JOIN MESSAGES ON USERS.Id_USERS = MESSAGES.Id_USERS_1 WHERE MESSAGES.Id_USERS = '" + id + "' OR Id_USERS_1 = '" + id + "'")
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't query database",
			})
			return
		}

		var users []User
		for rows.Next() {
			var user User
			err = rows.Scan(&user.IdUser, &user.FirstName, &user.LastName, &user.ProfilePicture)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "can't scan database",
				})
				return
			}
			users = append(users, user)
		}

		c.JSON(200, users)
	}
}