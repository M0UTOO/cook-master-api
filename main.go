package main

import (
	"cook-master-api/token"
	"database/sql"
	"github.com/gin-gonic/gin"
	//"io/ioutil"
	_ "github.com/go-sql-driver/mysql"
	//"fmt"
)

func main() {
	tokenAPI := token.GetAPIToken()
	r := gin.Default()
	r.GET("/user/:id", getUserByID(tokenAPI))
	r.POST("/user/email", postUserByEmail(tokenAPI))
	r.Run(":8080")
}

func getUserByID(tokenAPI string) func(c *gin.Context) {
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

		rows, err := db.Query("SELECT id, username FROM users")
		if err != nil {
			c.JSON(498, gin.H{
				"error": true,
				"message": "error on query request to bdd",
			})
			return
		}
		defer rows.Close()

		type User struct {
			Username string `json:"username"`
			Id int `json:"id"`
		}

		var list []User

		for rows.Next() {

			var user User

			err := rows.Scan(&user.Id, &user.Username)
			if err != nil {
				c.AbortWithStatus(500)
				return
			}

			list = append(list, user)

		}

		c.JSON(200, list)
		return
	}
}

func postUserByEmail (tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		type postUserByEmailReq struct {
			Password string `json:"password"`
			Email string `json:"email"`
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

		var req postUserByEmailReq

		err = c.BindJSON(&req)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "can't decode json request",
			})
			return
		}





		c.JSON(200, req)
		return
	}
}

func index(c *gin.Context) {
	c.JSON(200, gin.H{
		"error": true,
		"message": "hello world",
	})
}
