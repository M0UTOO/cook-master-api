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
	IsBlocked bool `json:"isblocked"`
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
			err = rows.Scan(&client.IdClient, &client.FidelityPoints, &client.StreetName, &client.Country, &client.City, &client.SteetNumber, &client.PhoneNumber, &client.Subscription, &client.IdUsers, &client.Id, &client.Email, &client.Password, &client.FirstName, &client.LastName, &client.ProfilePicture, &client.IsCreatedAt, &client.LastSeen, &client.IsBlocked)
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

		err = db.QueryRow("SELECT * FROM CLIENTS JOIN USERS ON CLIENTS.Id_USERS = USERS.Id_USERS WHERE CLIENTS.Id_USERS = " + id).Scan(&client.IdClient, &client.FidelityPoints, &client.StreetName, &client.Country, &client.City, &client.SteetNumber, &client.PhoneNumber, &client.Subscription, &client.IdUsers, &client.Id, &client.Email, &client.Password, &client.FirstName, &client.LastName, &client.ProfilePicture, &client.IsCreatedAt, &client.LastSeen, &client.IsBlocked)
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
		if client.SteetNumber >= 0 {
			setClause = append(setClause, "StreetNumber = '"+strconv.Itoa(client.SteetNumber)+"'")
		}
		if client.PhoneNumber != "" {
			setClause = append(setClause, "PhoneNumber = '"+client.PhoneNumber+"'")
		}
		if client.Subscription >= 0 {
			setClause = append(setClause, "Id_SUBSCRIPTIONS = '"+strconv.Itoa(client.Subscription)+"'")
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

		_, err = db.Exec("UPDATE CLIENTS SET " + strings.Join(setClause, ", ") + " WHERE Id_CLIENTS = " + id)
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