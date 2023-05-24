package main

import (
	"cook-master-api/token"
	"cook-master-api/users"
	"cook-master-api/clients"
	"cook-master-api/contractors"
	"cook-master-api/managers"
	"cook-master-api/events"
	"cook-master-api/premises"
	"cook-master-api/cookingspaces"
	"cook-master-api/cookingitems"
	"cook-master-api/ingredients"
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
		user.GET("/login", users.LoginUser(tokenAPI)) // WORKING
	client := r.Group("/client")
		client.GET("/all", clients.GetClients(tokenAPI)) // WORKING
		client.GET("/:id", clients.GetClientByID(tokenAPI)) // WORKING
		client.PATCH("/:id", clients.UpdateClient(tokenAPI)) // WORKING
		//client.DELETE("/:id", clients.DeleteClient(tokenAPI)) TO DO AFTER OTHERS TABLES
	contractor := r.Group("/contractor")
		contractor.GET("/all", contractors.GetContractors(tokenAPI)) // WORKING
		contractor.GET("/:id", contractors.GetContractorByID(tokenAPI)) // WORKING
		contractor.PATCH("/:id", contractors.UpdateContractor(tokenAPI)) // WORKING
		//contractor.DELETE("/:id", contractors.DeleteContractor(tokenAPI)) TO DO AFTER OTHERS TABLES
	manager := r.Group("/manager")
		manager.GET("/all", managers.GetManagers(tokenAPI)) // WORKING
		manager.GET("/:id", managers.GetManagerByID(tokenAPI)) // WORKING
		manager.PATCH("/:id", managers.UpdateManager(tokenAPI)) // WORKING
		// manager.DELETE("/:id", managers.DeleteManager(tokenAPI)) // TO DO AFTER OTHERS TABLES
	subscription := r.Group("/subscription")
		subscription.GET("/all", subscriptions.GetSubscriptions(tokenAPI)) // WORKING
		subscription.GET("/:id", subscriptions.GetSubscriptionByID(tokenAPI)) // WORKING
		subscription.POST("/", subscriptions.PostSubscription(tokenAPI)) // WORKING
		subscription.DELETE("/:id", subscriptions.DeleteSubscription(tokenAPI)) // WORKING
		subscription.PATCH("/:id", subscriptions.UpdateSubscription(tokenAPI)) // WORKING
	event := r.Group("/event")
		event.GET("/all", events.GetEvents(tokenAPI)) // WORKING
		event.GET("/:id", events.GetEventByID(tokenAPI)) // WORKING
		event.GET("/group/:id", events.GetEventsByGroupID(tokenAPI)) // WORKING
		event.POST("/:id", events.PostEvent(tokenAPI)) // WORKING
		event.POST("/group/:id", events.AddEventToAGroup(tokenAPI)) // WORKING
		event.DELETE("/group/:id", events.DeleteEventFromAGroup(tokenAPI)) // WORKING
		event.PATCH("/:id", events.UpdateEvent(tokenAPI)) // WORKING
		event.GET("/animate/:idevent/:iduser", events.AddContractorToAnEvent(tokenAPI)) // WORKING
		event.DELETE("/animate/:idevent/:iduser", events.DeleteContractorFromAnEvent(tokenAPI)) // WORKING
		event.GET("/participate/:idevent/:iduser", events.AddClientToAnEvent(tokenAPI)) // WORKING
		event.DELETE("/participate/:idevent/:iduser", events.DeleteClientFromAnEvent(tokenAPI)) // WORKING
		event.PATCH("/participate/:idevent/:iduser", events.ValidateClientPresence(tokenAPI)) // WORKING
	premise := r.Group("/premise")
		premise.GET("/all", premises.GetPremises(tokenAPI)) // WORKING
		premise.GET("/:id", premises.GetPremiseByID(tokenAPI)) // WORKING
		premise.POST("/", premises.PostPremise(tokenAPI)) // WORKING
		premise.DELETE("/:id", premises.DeletePremise(tokenAPI)) // WORKING
		premise.PATCH("/:id", premises.UpdatePremise(tokenAPI)) // WORKING
	cookingSpace := r.Group("/cookingspace")
		cookingSpace.GET("/all", cookingspaces.GetCookingSpaces(tokenAPI)) // WORKING
		cookingSpace.GET("/:id", cookingspaces.GetCookingSpaceByID(tokenAPI)) // WORKING
		cookingSpace.POST("/", cookingspaces.PostCookingSpace(tokenAPI)) // WORKING
		cookingSpace.PATCH("/:id", cookingspaces.UpdateCookingSpace(tokenAPI)) // WORKING
		cookingSpace.GET("/premise/:id", cookingspaces.GetCookingSpacesByPremiseID(tokenAPI)) // WORKING
		cookingSpace.POST("/premise/:id", cookingspaces.AddCookingSpaceToAPremise(tokenAPI)) // WORKING
		cookingSpace.DELETE("/premise/:id", cookingspaces.DeleteCookingSpaceFromAPremise(tokenAPI)) // WORKING
	cookingItem := r.Group("/cookingitem")
		cookingItem.GET("/all", cookingitems.GetCookingItems(tokenAPI)) // WORKING
		cookingItem.GET("/:id", cookingitems.GetCookingItemByID(tokenAPI)) // WORKING
		cookingItem.POST("/", cookingitems.PostCookingItem(tokenAPI)) // WORKING
		cookingItem.DELETE("/:id", cookingitems.DeleteCookingItem(tokenAPI)) // WORKING
		cookingItem.PATCH("/:id", cookingitems.UpdateCookingItem(tokenAPI)) // WORKING
		cookingItem.GET("/cookingspace/:id", cookingitems.GetCookingItemsByCookingSpaceID(tokenAPI)) // WORKING
	ingredient := r.Group("/ingredient")
		ingredient.GET("/all", ingredients.GetIngredients(tokenAPI)) // WORKING
		ingredient.GET("/:id", ingredients.GetIngredientByID(tokenAPI)) // WORKING
		ingredient.POST("/", ingredients.PostIngredient(tokenAPI)) // WORKING
		ingredient.DELETE("/:id", ingredients.DeleteIngredient(tokenAPI)) // WORKING
		ingredient.PATCH("/:id", ingredients.UpdateIngredient(tokenAPI)) // WORKING
		ingredient.GET("/cookingspace/:id", ingredients.GetIngredientsByCookingSpaceID(tokenAPI)) // WORKING
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
