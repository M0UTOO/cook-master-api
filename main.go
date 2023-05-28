package main

import (
	"cook-master-api/clients"
	"cook-master-api/contractors"
	"cook-master-api/cookingitems"
	"cook-master-api/cookingspaces"
	"cook-master-api/events"
	"cook-master-api/foods"
	"cook-master-api/ingredients"
	"cook-master-api/lessons"
	"cook-master-api/managers"
	"cook-master-api/premises"
	"cook-master-api/shopitems"
	"cook-master-api/comments"
	"cook-master-api/orders"
	"cook-master-api/bills"
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
	user.GET("/:id", users.GetUserByID(tokenAPI))                // WORKING
	user.GET("/search/:filter", users.GetUserByFilter(tokenAPI)) // WORKING BUT MAYBE NOT USEFUL
	user.PATCH("/:id", users.UpdateUser(tokenAPI))               // WORKING
	user.GET("/all", users.GetUsers(tokenAPI))                   // WORKING BUT MAYBE NOT USEFUL
	user.POST("/", users.PostUser(tokenAPI))                     // WORKING
	user.GET("/login", users.LoginUser(tokenAPI))                // WORKING
	client := r.Group("/client")
	client.GET("/all", clients.GetClients(tokenAPI))                                                    // WORKING
	client.GET("/:id", clients.GetClientByID(tokenAPI))                                                 // WORKING
	client.PATCH("/:id", clients.UpdateClient(tokenAPI))                                                // WORKING
	client.PATCH("/subscription/:idclient/:idsubscription", clients.UpdateClientSubscription(tokenAPI)) // WORKING
	client.GET("/watch/:idclient/:idlesson", clients.WatchLesson(tokenAPI))                             // MUST BE TESTED
	client.DELETE("/watch/:idclient/:idlesson", clients.UnwatchLesson(tokenAPI))                        // MUST BE TESTED
	//client.DELETE("/:id", clients.DeleteClient(tokenAPI)) TO DO AFTER OTHERS TABLES
	contractor := r.Group("/contractor")
	contractor.GET("/all", contractors.GetContractors(tokenAPI))     // WORKING
	contractor.GET("/:id", contractors.GetContractorByID(tokenAPI))  // WORKING
	contractor.PATCH("/:id", contractors.UpdateContractor(tokenAPI)) // WORKING
	contractor.POST("/type", contractors.AddAContractorType(tokenAPI)) // MUST BE TESTED
	contractor.DELETE("/type/:id", contractors.DeleteAContractorType(tokenAPI)) // MUST BE TESTED
	contractor.GET("/type", contractors.GetContractorTypes(tokenAPI)) // MUST BE TESTED
	//contractor.DELETE("/:id", contractors.DeleteContractor(tokenAPI)) TO DO AFTER OTHERS TABLES
	manager := r.Group("/manager")
	manager.GET("/all", managers.GetManagers(tokenAPI))     // WORKING
	manager.GET("/:id", managers.GetManagerByID(tokenAPI))  // WORKING
	manager.PATCH("/:id", managers.UpdateManager(tokenAPI)) // WORKING
	// manager.DELETE("/:id", managers.DeleteManager(tokenAPI)) // TO DO AFTER OTHERS TABLES
	subscription := r.Group("/subscription")
	subscription.GET("/all", subscriptions.GetSubscriptions(tokenAPI))      // WORKING
	subscription.GET("/:id", subscriptions.GetSubscriptionByID(tokenAPI))   // WORKING
	subscription.POST("/", subscriptions.PostSubscription(tokenAPI))        // WORKING
	subscription.DELETE("/:id", subscriptions.DeleteSubscription(tokenAPI)) // WORKING
	subscription.PATCH("/:id", subscriptions.UpdateSubscription(tokenAPI))  // WORKING
	event := r.Group("/event")
	event.GET("/all", events.GetEvents(tokenAPI))                                                // WORKING
	event.GET("/:id", events.GetEventByID(tokenAPI))                                             // WORKING
	event.GET("/group/:id", events.GetEventsByGroupID(tokenAPI))                                 // WORKING
	event.POST("/:id", events.PostEvent(tokenAPI))                                               // WORKING
	event.POST("/group/:id", events.AddEventToAGroup(tokenAPI))                                  // WORKING
	event.DELETE("/group/:id", events.DeleteEventFromAGroup(tokenAPI))                           // WORKING
	event.PATCH("/:id", events.UpdateEvent(tokenAPI))                                            // WORKING
	event.GET("/animate/:idevent/:iduser", events.AddContractorToAnEvent(tokenAPI))              // WORKING
	event.DELETE("/animate/:idevent/:iduser", events.DeleteContractorFromAnEvent(tokenAPI))      // WORKING
	event.GET("/participate/:idevent/:iduser", events.AddClientToAnEvent(tokenAPI))              // WORKING
	event.DELETE("/participate/:idevent/:iduser", events.DeleteClientFromAnEvent(tokenAPI))      // WORKING
	event.PATCH("/participate/:idevent/:iduser", events.ValidateClientPresence(tokenAPI))        // WORKING
	event.PATCH("/host/:idevent/:idcookingspace", events.AddEventToAnCookingSpace(tokenAPI))     // MUST BE TESTED
	event.DELETE("/host/:idevent/:idcookingspace", events.DeleteEventToAnCookingSpace(tokenAPI)) // MUST BE TESTED
	premise := r.Group("/premise")
	premise.GET("/all", premises.GetPremises(tokenAPI))      // WORKING
	premise.GET("/:id", premises.GetPremiseByID(tokenAPI))   // WORKING
	premise.POST("/", premises.PostPremise(tokenAPI))        // WORKING
	premise.DELETE("/:id", premises.DeletePremise(tokenAPI)) // WORKING
	premise.PATCH("/:id", premises.UpdatePremise(tokenAPI))  // WORKING
	cookingSpace := r.Group("/cookingspace")
	cookingSpace.GET("/all", cookingspaces.GetCookingSpaces(tokenAPI))                            // WORKING
	cookingSpace.GET("/:id", cookingspaces.GetCookingSpaceByID(tokenAPI))                         // WORKING
	cookingSpace.POST("/", cookingspaces.PostCookingSpace(tokenAPI))                              // WORKING
	cookingSpace.PATCH("/:id", cookingspaces.UpdateCookingSpace(tokenAPI))                        // WORKING
	cookingSpace.GET("/premise/:id", cookingspaces.GetCookingSpacesByPremiseID(tokenAPI))         // WORKING
	cookingSpace.POST("/premise/:id", cookingspaces.AddCookingSpaceToAPremise(tokenAPI))          // WORKING
	cookingSpace.DELETE("/premise/:id", cookingspaces.DeleteCookingSpaceFromAPremise(tokenAPI))   // WORKING
	cookingSpace.PATCH("/books/:idclient/:idcookingspace", cookingspaces.AddABooks(tokenAPI))     // MUST BE TESTED
	cookingSpace.DELETE("/books/:idclient/:idcookingspace", cookingspaces.DeleteABooks(tokenAPI)) // MUST BE TESTED
	cookingItem := r.Group("/cookingitem")
	cookingItem.GET("/all", cookingitems.GetCookingItems(tokenAPI))                              // WORKING
	cookingItem.GET("/:id", cookingitems.GetCookingItemByID(tokenAPI))                           // WORKING
	cookingItem.POST("/", cookingitems.PostCookingItem(tokenAPI))                                // WORKING
	cookingItem.DELETE("/:id", cookingitems.DeleteCookingItem(tokenAPI))                         // WORKING
	cookingItem.PATCH("/:id", cookingitems.UpdateCookingItem(tokenAPI))                          // WORKING
	cookingItem.GET("/cookingspace/:id", cookingitems.GetCookingItemsByCookingSpaceID(tokenAPI)) // WORKING
	ingredient := r.Group("/ingredient")
	ingredient.GET("/all", ingredients.GetIngredients(tokenAPI))                              // WORKING
	ingredient.GET("/:id", ingredients.GetIngredientByID(tokenAPI))                           // WORKING
	ingredient.POST("/", ingredients.PostIngredient(tokenAPI))                                // WORKING
	ingredient.DELETE("/:id", ingredients.DeleteIngredient(tokenAPI))                         // WORKING
	ingredient.PATCH("/:id", ingredients.UpdateIngredient(tokenAPI))                          // WORKING
	ingredient.GET("/cookingspace/:id", ingredients.GetIngredientsByCookingSpaceID(tokenAPI)) // WORKING
	lesson := r.Group("/lesson")
	lesson.GET("/all", lessons.GetLessons(tokenAPI))                      // MUST BE TESTED
	lesson.GET("/:id", lessons.GetLessonByID(tokenAPI))                   // MUST BE TESTED
	lesson.GET("/group/:id", lessons.GetLessonsByGroupID(tokenAPI))       // MUST BE TESTED
	lesson.POST("/", lessons.Postlesson(tokenAPI))                        // MUST BE TESTED
	lesson.POST("/group/:id", lessons.AddLessonToAGroup(tokenAPI))        // MUST BE TESTED
	lesson.DELETE("/group/:id", lessons.DeleteLessonFromAGroup(tokenAPI)) // MUST BE TESTED
	food := r.Group("/food")
	food.GET("/all", foods.GetFoods(tokenAPI))      // WORKING
	food.GET("/:id", foods.GetFoodByID(tokenAPI))   // WORKING
	food.POST("/", foods.PostFood(tokenAPI))        // WORKING
	food.DELETE("/:id", foods.DeleteFood(tokenAPI)) // WORKING
	food.PATCH("/:id", foods.UpdateFood(tokenAPI))  // WORKING
	shopitem := r.Group("/shopitem")
	shopitem.GET("/all", shopitems.GetShopItems(tokenAPI))      // WORKING
	shopitem.GET("/:id", shopitems.GetShopItemByID(tokenAPI))   // WORKING
	shopitem.POST("/", shopitems.PostShopItem(tokenAPI))        // WORKING
	shopitem.DELETE("/:id", shopitems.DeleteShopItem(tokenAPI)) // WORKING
	shopitem.PATCH("/:id", shopitems.UpdateShopItem(tokenAPI))  // WORKING
	comment := r.Group("/comment")
	comment.GET("/event/:id", comments.GetCommentsByClientID(tokenAPI))   // MUST BE TESTED
	comment.GET("/client/:id", comments.GetCommentsByEventID(tokenAPI))   // MUST BE TESTED
	comment.POST("/", comments.PostComment(tokenAPI))        // MUST BE TESTED
	comment.DELETE("/:id", comments.DeleteComment(tokenAPI)) // MUST BE TESTED
	comment.PATCH("/:id", comments.UpdateComment(tokenAPI))  // MUST BE TESTED
	bill := r.Group("/bill")
	bill.GET("/all", bills.GetBills(tokenAPI))      // MUST BE TESTED
	bill.GET("/:id", bills.GetBillByID(tokenAPI))   // MUST BE TESTED
	bill.GET("/user/:id", bills.GetBillsByUserID(tokenAPI))   // MUST BE TESTED
	bill.POST("/", bills.PostBill(tokenAPI))        // MUST BE TESTED
	bill.DELETE("/:id", bills.DeleteBill(tokenAPI)) // MUST BE TESTED
	bill.PATCH("/:id", bills.UpdateBill(tokenAPI))  // MUST BE TESTED
	order := r.Group("/order")
	order.GET("/all", orders.GetOrders(tokenAPI))      // MUST BE TESTED
	order.GET("/:id", orders.GetOrder(tokenAPI))   // MUST BE TESTED
	order.GET("/chef/:id", orders.GetOrderByContractor1ID(tokenAPI))   // MUST BE TESTED
	order.GET("/deliveryman/:id", orders.GetOrderByContractor2ID(tokenAPI))   // MUST BE TESTED
	order.GET("/client/:id", orders.GetOrderByClientID(tokenAPI))   // MUST BE TESTED
	order.POST("/", orders.PostOrder(tokenAPI))        // MUST BE TESTED
	order.DELETE("/:id", orders.DeleteOrder(tokenAPI)) // MUST BE TESTED
	order.PATCH("/:id", orders.UpdateOrder(tokenAPI))  // MUST BE TESTED
	order.PATCH("/item/:iditem/:idorder", orders.AddItemToAnOrder(tokenAPI))  // MUST BE TESTED
	order.DELETE("/item/:iditem/:idorder", orders.DeleteItemFromAnOrder(tokenAPI))  // MUST BE TESTED
	order.PATCH("/food/:idfood/:idorder", orders.AddFoodToAnOrder(tokenAPI))  // MUST BE TESTED
	order.DELETE("/food/:idfood/:idorder", orders.DeleteFoodFromAnOrder(tokenAPI))  // MUST BE TESTED
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
		"error":   true,
		"message": "hello world",
	})
}
