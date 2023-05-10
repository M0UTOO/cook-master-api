package main

import (
	"cook-master-api/token"
	"cook-master-api/users"
	//"cook-master-api/conversations"
	"cook-master-api/subscriptions"
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
	subscription := r.Group("/subscription")
		subscription.GET("/", subscriptions.GetSubscriptions(tokenAPI))
		subscription.POST("/", subscriptions.PostSubscription(tokenAPI))
		subscription.DELETE("/:id", subscriptions.DeleteSubscription(tokenAPI))
		subscription.PATCH("/:id", subscriptions.UpdateSubscription(tokenAPI))
	//conversation := r.Group("/conversations")
		//conversation.POST("/", conversations.PostConversations(tokenAPI))
		//conversation.DELETE("/", conversations.DeleteConversations(tokenAPI))
		//conversation.GET("/all", conversations.GetConversations(tokenAPI))
		//conversation.GET("/:id", conversations.GetConversationByID(tokenAPI))
		//conversation.GET("/user/:id", conversations.GetConversationForUserID(tokenAPI))
		//conversation.POST("/message", conversations.PostMessage(tokenAPI))
	r.Run(":9000")
}

func index(c *gin.Context) {
	c.JSON(200, gin.H{
		"error": true,
		"message": "hello world",
	})
}
