package subscriptions

import (
	"cook-master-api/token"
	"cook-master-api/utils"
	"database/sql"
	"strconv"
	"fmt"
	"strings"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Subscription struct {
	IdSubscription int `json:"idsubscription"`
	Name string `json:"name"`
	Price float64 `json:"price"`
	MaxLessonAccess int `json:"maxlessonaccess"`
	Picture string `json:"picture"`
	AllowRoomBooking bool `json:"allowroombooking"`
	AllowShopReduction bool `json:"allowshoreduction"`
	AllowChat bool `json:"allowchat"`
	Description string `json:"description"`
}

func GetSubscriptions(tokenAPI string) func(c *gin.Context) {
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

		rows, err := db.Query("SELECT * FROM SUBSCRIPTIONS")
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "can't query database",
			})
			return
		}
		defer rows.Close()

		var subscriptions []Subscription

		for rows.Next() {
			var subscription Subscription
			err := rows.Scan(&subscription.IdSubscription, &subscription.Name, &subscription.Price, &subscription.MaxLessonAccess, &subscription.Picture, &subscription.AllowRoomBooking, &subscription.AllowShopReduction, &subscription.AllowChat, &subscription.Description)
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "can't scan row",
				})
				return
			}
			subscriptions = append(subscriptions, subscription)
		}

		c.JSON(200, subscriptions)
	}
}

func GetSubscriptionByID(tokenAPI string) func(c *gin.Context) {
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

		var subscription Subscription

		err = db.QueryRow("SELECT * FROM SUBSCRIPTIONS WHERE Id_SUBSCRIPTIONS = " + id).Scan(&subscription.IdSubscription, &subscription.Name, &subscription.Price, &subscription.MaxLessonAccess, &subscription.Picture, &subscription.AllowRoomBooking, &subscription.AllowShopReduction, &subscription.AllowChat, &subscription.Description)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error": true,
				"message": "subscription not found",
			})
			return
		}

		c.JSON(200, subscription)
		return

	}
}

func PostSubscription(tokenAPI string) func(c *gin.Context) {
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

		var subscription Subscription

		err = c.BindJSON(&subscription)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "can't bind json",
			})
			return
		}

		if !utils.IsSafeString(subscription.Name) || !utils.IsSafeString(subscription.Description) {
			c.JSON(500, gin.H{
				"error": true,
				"message": "unsafe string",
			})
			return
		}

		if subscription.Name == "" || subscription.Price < 0 || subscription.MaxLessonAccess <= 0 || subscription.Description == "" {
			c.JSON(500, gin.H{
				"error": true,
				"message": "missing field",
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

		var id int

		err = db.QueryRow("SELECT Id_SUBSCRIPTIONS FROM SUBSCRIPTIONS WHERE name = ?", subscription.Name).Scan(&id)
		if err == nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "subscription already exist",
			})
			return
		}

		result, err := db.Exec("INSERT INTO SUBSCRIPTIONS (name, price, max_lesson_access, picture, description, allowRoomBooking, allowShopReduction, allowChat) VALUES (?, ?, ?, DEFAULT, ?, ?, ?, ?)", subscription.Name, subscription.Price, subscription.MaxLessonAccess, subscription.Description, subscription.AllowRoomBooking, subscription.AllowShopReduction, subscription.AllowChat)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error": true,
				"message": "can't insert into database",
			})
			return
		}

		lastId, err := result.LastInsertId()
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "can't get last id",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
			"message": "subscription added",
			"id": lastId,
		})
	}
}

func DeleteSubscription(tokenAPI string) func(c *gin.Context) {
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
			c.JSON(500, gin.H{
				"error": true,
				"message": "missing id",
			})
			return
		}

		if !utils.IsSafeString(id) {
			c.JSON(500, gin.H{
				"error": true,
				"message": "unsafe string",
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

		_, err = db.Exec("DELETE FROM IS_SUBSCRIBED WHERE Id_SUBSCRIPTIONS = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "can't delete from database",
			})
			return
		}

		_, err = db.Exec("DELETE FROM SUBSCRIPTIONS WHERE Id_SUBSCRIPTIONS = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "can't delete from database",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
			"message": "subscription deleted",
		})
	}
}

func UpdateSubscription(tokenAPI string) func(c *gin.Context) {
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

		id := c.Param("id")

		if id == "" {
			c.JSON(500, gin.H{
				"error": true,
				"message": "missing id",
			})
			return
		}

		if !utils.IsSafeString(id) {
			c.JSON(500, gin.H{
				"error": true,
				"message": "unsafe string",
			})
			return
		}

		var req Subscription

		req.Price = -1

		err = c.BindJSON(&req)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "can't decode json request",
			})
			return
		}

		var setClause []string

		if req.Name != "" {
			if !utils.IsSafeString(req.Name) {
				c.JSON(500, gin.H{
					"error": true,
					"message": "unsafe string",
				})
				return
			}
			setClause = append(setClause, "name = '"+req.Name+"'")
		}

		if req.Price >= 0 {
			setClause = append(setClause, "price = "+strconv.FormatFloat(req.Price, 'f', 2, 64))
		}

		if req.MaxLessonAccess > 0 {
			setClause = append(setClause, "max_lesson_access = "+strconv.Itoa(req.MaxLessonAccess))
		}

		if req.Picture != "" {
			if !utils.IsSafeString(req.Picture) {
				c.JSON(500, gin.H{
					"error": true,
					"message": "unsafe string",
				})
				return
			}
			setClause = append(setClause, "picture = '"+req.Picture+"'")
		}

		if req.Description != "" {
			if !utils.IsSafeString(req.Description) {
				c.JSON(500, gin.H{
					"error": true,
					"message": "unsafe string",
				})
				return
			}
			setClause = append(setClause, "description = '"+req.Description+"'")
		}

		if len(setClause) == 0 {
			c.JSON(500, gin.H{
				"error": true,
				"message": "missing field",
			})
			return
		}

		_, err = db.Exec("UPDATE SUBSCRIPTIONS SET " + strings.Join(setClause, ", ") + " WHERE Id_SUBSCRIPTIONS = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "can't update database",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
			"message": "subscription updated",
		})
	}
}

