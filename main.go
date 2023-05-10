package main

import (
	"cook-master-api/token"
	"cook-master-api/users"
	"cook-master-api/conversations"
	"github.com/gin-gonic/gin"
)

func main() {
	tokenAPI := token.GetAPIToken()
	r := gin.Default()
	r.GET("/", index)
	user := r.Group("/user")
		user.GET("/:id", users.GetUserByID(tokenAPI))
		user.GET("/search/:filter", users.GetUserByFilter(tokenAPI))
		user.PATCH("/", users.UpdateUser(tokenAPI))
		user.DELETE("/", users.DeleteUser(tokenAPI))
		user.GET("/all", users.GetUsers(tokenAPI))
		user.POST("/", users.PostUser(tokenAPI))
	//conversation := r.Group("/conversations")
		//conversation.POST("/", conversations.PostConversations(tokenAPI))
	r.Run(":9000")
}

func index(c *gin.Context) {
	c.JSON(200, gin.H{
		"error": true,
		"message": "hello world",
	})
}
