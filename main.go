package main

import (
	"cook-master-api/token"
	"cook-master-api/users"
	"cook-master-api/clients"
	"cook-master-api/contractors"
	"cook-master-api/managers"
	//"cook-master-api/conversations"
	"cook-master-api/subscriptions"
	"github.com/gin-gonic/gin"
)

func main() {
	tokenAPI := token.GetAPIToken()
	r := gin.Default()
	r.GET("/", index)
	user := r.Group("/user")
		user.GET("/:id", users.GetUserByID(tokenAPI)) // WORKING
		user.GET("/search/:filter", users.GetUserByFilter(tokenAPI)) // WORKING BUT MAYBE NOT USEFUL
		user.PATCH("/:id", users.UpdateUser(tokenAPI)) // WORKING
		user.GET("/all", users.GetUsers(tokenAPI)) // WORKING BUT MAYBE NOT USEFUL
		user.POST("/", users.PostUser(tokenAPI)) // WORKING
	client := r.Group("/client")
		client.GET("/all", clients.GetClients(tokenAPI)) // WORKING
		client.GET("/:id", clients.GetClientByID(tokenAPI)) // WORKING
		client.PATCH("/:id", clients.UpdateClient(tokenAPI)) // WORKING
		client.GET("/login" , clients.LoginClient(tokenAPI)) // WORKING
		//client.DELETE("/:id", clients.DeleteClient(tokenAPI)) TO DO AFTER OTHERS TABLES
	contractor := r.Group("/contractor")
		contractor.GET("/all", contractors.GetContractors(tokenAPI)) // WORKING
		contractor.GET("/:id", contractors.GetContractorByID(tokenAPI)) // WORKING
		contractor.PATCH("/:id", contractors.UpdateContractor(tokenAPI)) // WORKING
		contractor.GET("/login" , contractors.LoginContractor(tokenAPI)) // WORKING
		//contractor.DELETE("/:id", contractors.DeleteContractor(tokenAPI)) TO DO AFTER OTHERS TABLES
	manager := r.Group("/manager")
		manager.GET("/all", managers.GetManagers(tokenAPI)) // WORKING
		manager.GET("/:id", managers.GetManagerByID(tokenAPI)) // WORKING
		manager.PATCH("/:id", managers.UpdateManager(tokenAPI)) // WORKING
		manager.GET("/login" , managers.LoginManager(tokenAPI)) // WORKING
		// manager.DELETE("/:id", managers.DeleteManager(tokenAPI)) // TO DO AFTER OTHERS TABLES
	subscription := r.Group("/subscription")
		subscription.GET("/", subscriptions.GetSubscriptions(tokenAPI)) // WORKING
		subscription.POST("/", subscriptions.PostSubscription(tokenAPI)) // WORKING
		subscription.DELETE("/:id", subscriptions.DeleteSubscription(tokenAPI)) // WORKING
		subscription.PATCH("/:id", subscriptions.UpdateSubscription(tokenAPI)) // WORKING
	//conversation := r.Group("/conversations") TO DO AFTER OTHERS TABLES AND RE WORK ON THE MDC
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
