package clients

import (
	"cook-master-api/token"
	"cook-master-api/utils"
	"database/sql"
	"strconv"
	"strings"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type ClientUser struct {
	Id int `json:"id"`
	Email string `json:"email"`
	Password string `json:"password"`
	FirstName string `json:"firstname"`
	LastName string `json:"lastname"`
	ProfilePicture string `json:"profilepicture"`
	IsCreatedAt string `json:"iscreatedat"`
	LastSeen string `json:"lastseen"`
	IsBlocked string `json:"isblocked"`
	IdClient int `json:"idclient"`
	FidelityPoints int `json:"fidelitypoints"`
	StreetName string `json:"streetname"`
	Country string `json:"country"`
	City string `json:"city"`
	SteetNumber int `json:"streetnumber"`
	PhoneNumber string `json:"phonenumber"`
	Subscription int `json:"subscription"`
	IdUsers int `json:"idusers"`
}

type Client struct {
	FidelityPoints int `json:"fidelitypoints"`
	StreetName string `json:"streetname"`
	Country string `json:"country"`
	City string `json:"city"`
	SteetNumber int `json:"streetnumber"`
	PhoneNumber string `json:"phonenumber"`
	Subscription int `json:"subscription"`
}


func GetClients(tokenAPI string) func(c *gin.Context) {
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

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		rows, err := db.Query("SELECT * FROM CLIENTS JOIN USERS ON CLIENTS.Id_USERS = USERS.Id_USERS ORDER BY USERS.lastname DESC")
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "user not found",
			})
			return
		}

		var clients []ClientUser

		for rows.Next() {
			var client ClientUser
			err = rows.Scan(&client.IdClient, &client.FidelityPoints, &client.StreetName, &client.Country, &client.City, &client.SteetNumber, &client.PhoneNumber, &client.IdUsers, &client.Id, &client.Email, &client.Password, &client.FirstName, &client.LastName, &client.ProfilePicture, &client.IsCreatedAt, &client.LastSeen, &client.IsBlocked)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error": true,
					"message": "client not found",
				})
				return
			}
			clients = append(clients, client)
		}

		c.JSON(200, clients)
		return
	}
}

func GetClientByID(tokenAPI string) func(c *gin.Context) {
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

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		var client ClientUser

		err = db.QueryRow("SELECT * FROM CLIENTS JOIN USERS ON CLIENTS.Id_USERS = USERS.Id_USERS WHERE CLIENTS.Id_USERS = " + id).Scan(&client.IdClient, &client.FidelityPoints, &client.StreetName, &client.Country, &client.City, &client.SteetNumber, &client.PhoneNumber, &client.IdUsers, &client.Id, &client.Email, &client.Password, &client.FirstName, &client.LastName, &client.ProfilePicture, &client.IsCreatedAt, &client.LastSeen, &client.IsBlocked)
		fmt.Println(err)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "client not found",
			})
			return
		}

		c.JSON(200, client)
		return

	}
}

func UpdateClient(tokenAPI string) func(c *gin.Context) {
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

		var client Client

		client.FidelityPoints = -1

		err = c.BindJSON(&client)
		if err != nil {
			c.JSON(400, gin.H{
				"error": true,
				"message": "bad json",
			})
			return
		}

		var setClause []string
		
		if !utils.IsSafeString(client.StreetName) || !utils.IsSafeString(client.Country) || !utils.IsSafeString(client.City) || !utils.IsSafeString(client.PhoneNumber) { 
			c.JSON(400, gin.H{
				"error": true,
				"message": "bad json",
			})
			return
		}

		if client.FidelityPoints >= 0 {
			setClause = append(setClause, "FidelityPoints = "+strconv.Itoa(client.FidelityPoints))
		}
		if client.StreetName != "" {
			setClause = append(setClause, "StreetName = '"+client.StreetName+"'")
		}
		if client.Country != "" {
			setClause = append(setClause, "Country = '"+client.Country+"'")
		}
		if client.City != "" {
			setClause = append(setClause, "City = '"+client.City+"'")
		}
		if client.SteetNumber > 0 {
			setClause = append(setClause, "StreetNumber = '"+strconv.Itoa(client.SteetNumber)+"'")
		}
		if client.PhoneNumber != "" {
			setClause = append(setClause, "PhoneNumber = '"+client.PhoneNumber+"'")
		}

		if len(setClause) == 0 {
			c.JSON(400, gin.H{
				"error": true,
				"message": "nothing to update",
			})
			return
		}

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		_, err = db.Exec("UPDATE CLIENTS SET " + strings.Join(setClause, ", ") + " WHERE Id_USERS = " + id)
		fmt.Println(err)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "client not found",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
			"message": "client updated",
		})
		return
	}
}

func UpdateClientSubscription(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		type Subscription struct {
			EndTime string `json:"endTime"`
		}

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

		idclient := c.Param("idclient")
		if idclient == "" {
			c.JSON(400, gin.H{
				"error": true,
				"message": "idclient can't be empty",
			})
			return
		}

		if !utils.IsSafeString(idclient) {
			c.JSON(400, gin.H{
				"error": true,
				"message": "idclient can't contain sql injection",
			})
			return
		}

		idsubscription := c.Param("idsubscription")
		if idsubscription == "" {
			c.JSON(400, gin.H{
				"error": true,
				"message": "idsubscription can't be empty",
			})
			return
		}

		if !utils.IsSafeString(idsubscription) {
			c.JSON(400, gin.H{
				"error": true,
				"message": "idsubscription can't contain sql injection",
			})
			return
		}

		var subscription Subscription
		err = c.BindJSON(&subscription)
		if err != nil {
			c.JSON(400, gin.H{
				"error": true,
				"message": "bad json",
			})
			return
		}

		if subscription.EndTime == "" || !utils.IsSafeString(subscription.EndTime) {
			c.JSON(400, gin.H{
				"error": true,
				"message": "bad json",
			})
			return
		}

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		err = db.QueryRow("SELECT Id_USERS FROM CLIENTS WHERE Id_USERS = " + idclient).Scan(&idclient)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "client not found",
			})
			return
		}

		err = db.QueryRow("SELECT Id_SUBSCRIPTIONS FROM SUBSCRIPTIONS WHERE Id_SUBSCRIPTIONS = " + idsubscription).Scan(&idsubscription)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "subscription not found",
			})
			return
		}

		_, err = db.Exec("INSERT INTO IS_SUBSCRIBED (Id_CLIENTS, Id_SUBSCRIPTIONS, endtime) VALUES (?, ?, ?)", idclient, idsubscription, subscription.EndTime)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "canno insert is_subscribed",
			})
			return
		}

		_, err = db.Exec("UPDATE CLIENTS SET Subscription = " + idsubscription + " WHERE Id_USERS = " + idclient)
		fmt.Println(err)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot update client subscription",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
			"message": "subscription update",
		})
		return
	}
}

func WatchLesson(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header["Token"]

		if len(tokenHeader) == 0 {
			c.JSON(400, gin.H{
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

		idClient := c.Param("idclient")
		if idClient == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id client can't be empty",
			})
			return
		}

		if !utils.IsSafeString(idClient) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id client can't contain sql injection",
			})
			return
		}

		idLesson := c.Param("idlesson")
		if idLesson == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "idlesson can't be empty",
			})
			return
		}

		if !utils.IsSafeString(idLesson) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "idlesson can't contain sql injection",
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

		var idclient string

		err = db.QueryRow("SELECT Id_CLIENTS FROM CLIENTS WHERE Id_USERS = '" + idClient + "'").Scan(&idclient)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "client not found",
			})
			return
		}

		var idlesson string

		err = db.QueryRow("SELECT Id_LESSONS FROM LESSONS WHERE Id_LESSONS = '" + idLesson + "'").Scan(&idlesson)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "lesson not found",
			})
			return
		}

		var watches string

		err = db.QueryRow("SELECT Id_CLIENTS FROM WATCHES WHERE Id_LESSONS = '" + idLesson + "' AND Id_CLIENTS = '" + idClient + "'").Scan(&watches)
		if err == nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "client already watched this lesson",
			})
			return
		}

		_, err = db.Exec("INSERT INTO WATCHES (Id_CLIENTS, Id_LESSONS) VALUES (?, ?)", idLesson, idclient)
		fmt.Println(err)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot insert watches",
			})
			return
		}

		c.JSON(200, gin.H{
			"error":   false,
			"message": "lesson added to watched",
		})
	}
}

func UnwatchLesson(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header["Token"]

		if len(tokenHeader) == 0 {
			c.JSON(400, gin.H{
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

		idClient := c.Param("idclient")
		if idClient == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id event can't be empty",
			})
			return
		}

		if !utils.IsSafeString(idClient) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id event can't contain sql injection",
			})
			return
		}

		idLesson := c.Param("idlesson")
		if idLesson == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "idlesson can't be empty",
			})
			return
		}

		if !utils.IsSafeString(idLesson) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "idlesson can't contain sql injection",
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

		var idclient string

		err = db.QueryRow("SELECT Id_CLIENTS FROM CLIENTS WHERE Id_USERS = '" + idClient + "'").Scan(&idclient)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "client not found",
			})
			return
		}

		var idlesson string

		err = db.QueryRow("SELECT Id_LESSONS FROM LESSONS WHERE Id_LESSONS = '" + idLesson + "'").Scan(&idlesson)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "lesson not found",
			})
			return
		}

		var watches string

		err = db.QueryRow("SELECT Id_CLIENTS FROM WATCHES WHERE Id_LESSONS = '" + idLesson + "' AND Id_CLIENTS = '" + idclient + "'").Scan(&watches)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "client didn't watch this lesson",
			})
			return
		}

		_, err = db.Exec("DELETE FROM WATCHES WHERE Id_CLIENTS = '" + idClient + "' AND Id_LESSONS = '" + idLesson + "'")
		fmt.Println(err)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot delete watches",
			})
			return
		}

		c.JSON(200, gin.H{
			"error":   false,
			"message": "lesson deleted from watched",
		})
	}
}
