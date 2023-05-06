package main

import (
	"cook-master-api/token"

	"github.com/gin-gonic/gin"
)

func main() {
	tokenAPI := token.GetAPIToken()
	r := gin.Default()
	r.GET("/user/:id", getUserByID(tokenAPI))
	r.Run(":8080")
}

func getUserByID(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header["Token"]
		err := token.CheckAPIToken(tokenAPI, tokenHeader[0], c)
		if err != nil {
			c.JSON(401, err)
			return
		}
		c.JSON(200, "lala")
		return
	}
}

func index(c *gin.Context) {
	c.String(404, "Not Found")
}
