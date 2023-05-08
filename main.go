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
	r.GET("/users")
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

		id := c.Param("id")
		if id == "" {
			c.JSON(400, gin.H{
				"error": true,
				"message": "id can't be empty",
			})
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

		type User struct {
			Id int `json:"id"`
			Username string `json:"username"`
			Password string `json:"password"`
		}

		var user User

		err = db.QueryRow("SELECT * FROM users WHERE id=" + id).Scan(&user.Id, &user.Username, &user.Password)
		if err != nil {
			c.JSON(498, gin.H{
				"error": err,
				"message": "error on query request to bdd",
			})
			return
		}

		c.JSON(200, user)
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
