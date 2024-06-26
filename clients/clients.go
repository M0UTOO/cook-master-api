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
	Email string `json:"email"`
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
	SteetNumber string `json:"streetnumber"`
	PhoneNumber string `json:"phonenumber"`
	Subscription int `json:"subscription"`
	KeepSubscription bool `json:"keepsubscription"`
	IdUsers int `json:"idusers"`
	Language int `json:"language"`
}

type Client struct {
	FidelityPoints int `json:"fidelitypoints"`
	StreetName string `json:"streetname"`
	Country string `json:"country"`
	City string `json:"city"`
	SteetNumber string `json:"streetnumber"`
	PhoneNumber string `json:"phonenumber"`
	Subscription int `json:"subscription"`
	KeepSubscription bool `json:"keepsubscription"`
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

		rows, err := db.Query("SELECT USERS.email, USERS.firstname, USERS.lastname, USERS.profilepicture, USERS.iscreatedat, USERS.lastseen, USERS.isblocked, USERS.Id_LANGUAGES, CLIENTS.Id_CLIENTS, CLIENTS.fidelitypoints, CLIENTS.streetname, CLIENTS.country, CLIENTS.city, CLIENTS.streetnumber, CLIENTS.phonenumber, CLIENTS.subscription, CLIENTS.keepsubscription, CLIENTS.Id_USERS FROM CLIENTS JOIN USERS ON CLIENTS.Id_USERS = USERS.Id_USERS ORDER BY USERS.lastname DESC")
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error": true,
				"message": "user not found",
			})
			return
		}

		var clients []ClientUser

		for rows.Next() {
			var client ClientUser
			err = rows.Scan(&client.Email, &client.FirstName, &client.LastName, &client.ProfilePicture, &client.IsCreatedAt, &client.LastSeen, &client.IsBlocked, &client.Language, &client.IdClient, &client.FidelityPoints, &client.StreetName, &client.Country, &client.City, &client.SteetNumber, &client.PhoneNumber, &client.Subscription, &client.KeepSubscription, &client.IdUsers)
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

		err = db.QueryRow("SELECT USERS.email, USERS.firstname, USERS.lastname, USERS.profilepicture, USERS.iscreatedat, USERS.lastseen, USERS.isblocked, USERS.Id_LANGUAGES, CLIENTS.Id_CLIENTS, CLIENTS.fidelitypoints, CLIENTS.streetname, CLIENTS.country, CLIENTS.city, CLIENTS.streetnumber, CLIENTS.phonenumber, CLIENTS.subscription, CLIENTS.keepsubscription, CLIENTS.Id_USERS FROM CLIENTS JOIN USERS ON CLIENTS.Id_USERS = USERS.Id_USERS WHERE CLIENTS.Id_USERS = " + id).Scan(&client.Email, &client.FirstName, &client.LastName, &client.ProfilePicture, &client.IsCreatedAt, &client.LastSeen, &client.IsBlocked, &client.Language, &client.IdClient, &client.FidelityPoints, &client.StreetName, &client.Country, &client.City, &client.SteetNumber, &client.PhoneNumber, &client.Subscription, &client.KeepSubscription, &client.IdUsers)
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

		type ClientReq struct {
			FidelityPoints int `json:"fidelitypoints"`
			StreetName string `json:"streetname"`
			Country string `json:"country"`
			City string `json:"city"`
			SteetNumber string `json:"streetnumber"`
			PhoneNumber string `json:"phonenumber"`
			Subscription int `json:"subscription"`
			KeepSubscription int `json:"keepsubscription"`
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

		var client ClientReq

		client.FidelityPoints = -1
		client.KeepSubscription = -1

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
		if client.StreetName != "" || !utils.IsSafeString(client.StreetName) || len(client.StreetName) > 100 || len(client.StreetName) < 0 {
			setClause = append(setClause, "StreetName = '"+client.StreetName+"'")
		}
		if client.Country != "" || !utils.IsSafeString(client.Country) || len(client.Country) > 50 || len(client.Country) < 0 {
			setClause = append(setClause, "Country = '"+client.Country+"'")
		}
		if client.City != "" || !utils.IsSafeString(client.City) || len(client.City) > 100 || len(client.City) < 0 {
			setClause = append(setClause, "City = '"+client.City+"'")
		}
		if client.SteetNumber != "" || !utils.IsSafeString(client.SteetNumber) || len(client.SteetNumber) > 10 || len(client.SteetNumber) < 0 {
			setClause = append(setClause, "StreetNumber = '"+client.SteetNumber+"'")
		}
		if client.PhoneNumber != "" || !utils.IsSafeString(client.PhoneNumber) || len(client.PhoneNumber) > 25 || len(client.PhoneNumber) < 0 {
			setClause = append(setClause, "PhoneNumber = '"+client.PhoneNumber+"'")
		}
		if client.KeepSubscription == 0 {
			setClause = append(setClause, "KeepSubscription = false")
		} else if client.KeepSubscription == 1 {
			setClause = append(setClause, "KeepSubscription = true")
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

		iduser := c.Param("iduser")
		if iduser == "" {
			c.JSON(400, gin.H{
				"error": true,
				"message": "iduser can't be empty",
			})
			return
		}

		if !utils.IsSafeString(iduser) {
			c.JSON(400, gin.H{
				"error": true,
				"message": "iduser can't contain sql injection",
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

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		var idclient string

		err = db.QueryRow("SELECT Id_CLIENTS FROM CLIENTS WHERE Id_USERS = " + iduser).Scan(&idclient)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "client not found",
			})
			return
		}

		var idSubscription string

		err = db.QueryRow("SELECT Id_SUBSCRIPTIONS FROM SUBSCRIPTIONS WHERE Id_SUBSCRIPTIONS = " + idsubscription).Scan(&idSubscription)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "subscription not found",
			})
			return
		}

		_, err = db.Exec("INSERT INTO IS_SUBSCRIBED (Id_CLIENTS, Id_SUBSCRIPTIONS, endtime) VALUES (?, ?, DATE_ADD(NOW(), INTERVAL 1 MONTH))", idclient, idsubscription)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot insert is_subscribed",
			})
			return
		}

		_, err = db.Exec("UPDATE CLIENTS SET Subscription = " + idsubscription + " WHERE Id_USERS = " + iduser)
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

		err = db.QueryRow("SELECT Id_CLIENTS FROM WATCHES WHERE Id_LESSONS = '" + idLesson + "' AND Id_CLIENTS = '" + idclient + "'").Scan(&watches)
		if err == nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "client already watched this lesson",
			})
			return
		}

		_, err = db.Exec("INSERT INTO WATCHES (Id_CLIENTS, Id_LESSONS, dateTime) VALUES (?, ?, DEFAULT)", idclient, idlesson)
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

func UnWatchLesson(tokenAPI string) func(c *gin.Context) {
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

func GetAllSubscription(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header["Token"]

		type SubscriptionCount struct {
			Subscription string `json:"subscription"`
			Count        int    `json:"count"`
		}

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

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		rows, err := db.Query("SELECT SUBSCRIPTIONS.name AS name, COUNT(*) AS count FROM CLIENTS JOIN SUBSCRIPTIONS ON CLIENTS.subscription = SUBSCRIPTIONS.Id_SUBSCRIPTIONS GROUP BY subscription")
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "subscription not found",
			})
			return
		}

		var subscriptions []SubscriptionCount

		for rows.Next() {
			var subscription SubscriptionCount
			err = rows.Scan(&subscription.Subscription, &subscription.Count)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error":   true,
					"message": "subscription not found",
				})
				return
			}
			subscriptions = append(subscriptions, subscription)
		}

		c.JSON(200, subscriptions)
		return
	}
}

func GetAverageClientParticipationByMonth(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header["Token"]

		type Participation struct {
			Month string `json:"month"`
			Count float64    `json:"count"`
		}

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

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		rows, err := db.Query("SELECT DATE_FORMAT(endTime, '%Y-%m') AS month, AVG(CAST(isPresent AS INT)) AS average_participation FROM EVENTS JOIN PARTICIPATES ON EVENTS.Id_EVENTS = PARTICIPATES.Id_EVENTS GROUP BY month ORDER BY month")
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "subscription not found",
			})
			return
		}

		var participations []Participation

		for rows.Next() {
			var participation Participation
			err = rows.Scan(&participation.Month, &participation.Count)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error":   true,
					"message": "subscription not found",
				})
				return
			}
			participations = append(participations, participation)
		}

		c.JSON(200, participations)
		return
	}
}

func GetAverageMoneySpentByClient(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header["Token"]

		type Participation struct {
			Month string `json:"month"`
			Count float64    `json:"count"`
		}

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

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		rows, err := db.Query("SELECT COALESCE(DATE_FORMAT(createdAt, '%Y-%m'), '2023-07') AS month, COALESCE(AVG(total_price), 0) AS average_money_spent FROM (SELECT c.Id_CLIENTS, o.createdAt, SUM(o.price) + COALESCE(SUM(sub.price), 0) + COALESCE(SUM(cs.PricePerHour), 0) AS total_price FROM CLIENTS c LEFT JOIN ORDERS o ON c.Id_CLIENTS = o.Id_CLIENTS LEFT JOIN BOOKS b ON c.Id_CLIENTS = b.Id_CLIENTS LEFT JOIN IS_SUBSCRIBED s ON c.Id_CLIENTS = s.Id_CLIENTS LEFT JOIN SUBSCRIPTIONS sub ON s.Id_SUBSCRIPTIONS = sub.Id_SUBSCRIPTIONS LEFT JOIN COOKING_SPACES cs ON b.Id_COOKING_SPACES = cs.Id_COOKING_SPACES WHERE c.keepSubscription = TRUE GROUP BY c.Id_CLIENTS, o.createdAt ) AS t GROUP BY month ORDER BY month")
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "subscription not found",
			})
			return
		}

		var participations []Participation

		for rows.Next() {
			var participation Participation
			err = rows.Scan(&participation.Month, &participation.Count)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error":   true,
					"message": "subscription not found",
				})
				return
			}
			participations = append(participations, participation)
		}

		c.JSON(200, participations)
		return
	}
}

func GetClientsByCountry(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header["Token"]

		type Participation struct {
			Country string `json:"country"`
			Count float64    `json:"count"`
		}

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

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		rows, err := db.Query("SELECT country, COUNT(*) AS client_count FROM CLIENTS GROUP BY country ORDER BY client_count DESC;")
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "subscription not found",
			})
			return
		}

		var participations []Participation

		for rows.Next() {
			var participation Participation
			err = rows.Scan(&participation.Country, &participation.Count)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error":   true,
					"message": "subscription not found",
				})
				return
			}
			participations = append(participations, participation)
		}

		c.JSON(200, participations)
		return
	}
}

func GetTop5Client(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header["Token"]

		type Top5Client struct {
			IdClient int `json:"idclient"`
			Name string `json:"name"`
			Count        int    `json:"count"`
		}

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

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		rows, err := db.Query("SELECT c.Id_CLIENTS, u.firstname, COUNT(*) AS participationCount FROM CLIENTS c LEFT JOIN PARTICIPATES p ON c.Id_CLIENTS = p.Id_CLIENTS LEFT JOIN USERS u ON c.Id_USERS = u.Id_USERS LEFT JOIN COMMENTS cm ON c.Id_CLIENTS = cm.Id_CLIENTS LEFT JOIN IS_SUBSCRIBED s ON c.Id_CLIENTS = s.Id_CLIENTS LEFT JOIN ORDERS o ON c.Id_CLIENTS = o.Id_CLIENTS LEFT JOIN FORMATIONS f ON c.Id_CLIENTS = f.Id_CLIENTS LEFT JOIN BOOKS b ON c.Id_CLIENTS = b.Id_CLIENTS GROUP BY c.Id_CLIENTS ORDER BY participationCount DESC LIMIT 30;")
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "client not found",
			})
			return
		}

		var top5clients []Top5Client

		for rows.Next() {
			var top5client Top5Client
			err = rows.Scan(&top5client.IdClient, &top5client.Name, &top5client.Count)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error":   true,
					"message": "top 5 not found",
				})
				return
			}
			top5clients = append(top5clients, top5client)
		}

		c.JSON(200, top5clients)
		return
	}
}