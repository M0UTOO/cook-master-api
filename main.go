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
	user.GET("/search/:filter", users.GetUserByFilter(tokenAPI)) // WORKING
	user.PATCH("/:id", users.UpdateUser(tokenAPI))               // WORKING
	user.GET("/all", users.GetUsers(tokenAPI))                   // WORKING
	user.POST("/", users.PostUser(tokenAPI))                     // WORKING
	user.GET("/login", users.LoginUser(tokenAPI)) 				 // WORKING
	user.DELETE("/:id", users.DeleteUser(tokenAPI))              // WORKING
	client := r.Group("/client")
	client.GET("/all", clients.GetClients(tokenAPI))       			// WORKING                                             
	client.GET("/:id", clients.GetClientByID(tokenAPI))             // WORKING                                    
	client.PATCH("/:id", clients.UpdateClient(tokenAPI))            // WORKING                          
	client.PATCH("/subscription/:iduser/:idsubscription", clients.UpdateClientSubscription(tokenAPI))   // WORKING
	client.PATCH("/watch/:idclient/:idlesson", clients.WatchLesson(tokenAPI))  							// WORKING                           
	client.DELETE("/watch/:idclient/:idlesson", clients.UnWatchLesson(tokenAPI))     					// WORKING               
	contractor := r.Group("/contractor")
	contractor.GET("/all", contractors.GetContractors(tokenAPI))    			// WORKING 
	contractor.GET("/:id", contractors.GetContractorByID(tokenAPI))  			// WORKING
	contractor.PATCH("/:id", contractors.UpdateContractor(tokenAPI))  			// WORKING
	contractor.POST("/type", contractors.AddAContractorType(tokenAPI)) 			// WORKING
	contractor.DELETE("/type/:id", contractors.DeleteAContractorType(tokenAPI)) // WORKING
	contractor.GET("/type", contractors.GetContractorTypes(tokenAPI)) 			// WORKING
	manager := r.Group("/manager")
	manager.GET("/all", managers.GetManagers(tokenAPI))      // WORKING
	manager.GET("/:id", managers.GetManagerByID(tokenAPI))   // WORKING
	manager.PATCH("/:id", managers.UpdateManager(tokenAPI))  // WORKING
	subscription := r.Group("/subscription")
	subscription.GET("/all", subscriptions.GetSubscriptions(tokenAPI))      // WORKING
	subscription.GET("/:id", subscriptions.GetSubscriptionByID(tokenAPI))    // WORKING
	subscription.POST("/", subscriptions.PostSubscription(tokenAPI))       // WORKING 
	subscription.DELETE("/:id", subscriptions.DeleteSubscription(tokenAPI))  // WORKING
	subscription.PATCH("/:id", subscriptions.UpdateSubscription(tokenAPI))  // WORKING
	event := r.Group("/event")
	event.GET("/all", events.GetEvents(tokenAPI))              // WORKING                                  
	event.GET("/:id", events.GetEventByID(tokenAPI))                        // WORKING                     
	event.GET("/group/:id", events.GetEventsByGroupID(tokenAPI))            // WORKING                     
	event.POST("/:id", events.PostEvent(tokenAPI))                        // WORKING                       
	event.POST("/group/:id", events.AddEventToAGroup(tokenAPI))                              // WORKING    
	event.DELETE("/group/:id", events.DeleteEventFromAGroup(tokenAPI))                    // WORKING       
	event.PATCH("/:id", events.UpdateEvent(tokenAPI))                                            // WORKING
	event.GET("/animate/:idevent/:iduser", events.AddContractorToAnEvent(tokenAPI))              // WORKING
	event.DELETE("/animate/:idevent/:iduser", events.DeleteContractorFromAnEvent(tokenAPI))      // WORKING
	event.GET("/participate/:idevent/:iduser", events.AddClientToAnEvent(tokenAPI))              // WORKING
	event.DELETE("/participate/:idevent/:iduser", events.DeleteClientFromAnEvent(tokenAPI))      // WORKING
	event.PATCH("/participate/:idevent/:iduser", events.ValidateClientPresence(tokenAPI))        // WORKING
	event.PATCH("/host/:idevent/:idcookingspace", events.AddEventToAnCookingSpace(tokenAPI))     // WORKING
	event.DELETE("/host/:idevent/:idcookingspace", events.DeleteEventToAnCookingSpace(tokenAPI)) // WORKING
	premise := r.Group("/premise")
	premise.GET("/all", premises.GetPremises(tokenAPI))      // WORKING
	premise.GET("/:id", premises.GetPremiseByID(tokenAPI))   // WORKING
	premise.POST("/", premises.PostPremise(tokenAPI))        // WORKING
	premise.DELETE("/:id", premises.DeletePremise(tokenAPI)) // WORKING
	premise.PATCH("/:id", premises.UpdatePremise(tokenAPI))  // WORKING
	cookingSpace := r.Group("/cookingspace")
	cookingSpace.GET("/all", cookingspaces.GetCookingSpaces(tokenAPI))           // WORKING                 
	cookingSpace.GET("/:id", cookingspaces.GetCookingSpaceByID(tokenAPI))      // WORKING                   
	cookingSpace.POST("/", cookingspaces.PostCookingSpace(tokenAPI))       // WORKING                       
	cookingSpace.PATCH("/:id", cookingspaces.UpdateCookingSpace(tokenAPI))                  // WORKING      
	cookingSpace.GET("/premise/:id", cookingspaces.GetCookingSpacesByPremiseID(tokenAPI))       // WORKING  
	cookingSpace.POST("/premise/:id", cookingspaces.AddCookingSpaceToAPremise(tokenAPI))          // WORKING
	cookingSpace.DELETE("/premise/:id", cookingspaces.DeleteCookingSpaceFromAPremise(tokenAPI))   // WORKING
	cookingSpace.PATCH("/books/:idclient/:idcookingspace", cookingspaces.AddABooks(tokenAPI))     // WORKING
	cookingSpace.DELETE("/books/:idclient/:idcookingspace", cookingspaces.DeleteABooks(tokenAPI)) // WORKING
	// TO DO ADD DELETE A COOKING SPACE
	cookingItem := r.Group("/cookingitem")
	cookingItem.GET("/all", cookingitems.GetCookingItems(tokenAPI))                              
	cookingItem.GET("/:id", cookingitems.GetCookingItemByID(tokenAPI))                           
	cookingItem.POST("/", cookingitems.PostCookingItem(tokenAPI))                                
	cookingItem.DELETE("/:id", cookingitems.DeleteCookingItem(tokenAPI))                         
	cookingItem.PATCH("/:id", cookingitems.UpdateCookingItem(tokenAPI))                          
	cookingItem.GET("/cookingspace/:id", cookingitems.GetCookingItemsByCookingSpaceID(tokenAPI)) 
	ingredient := r.Group("/ingredient")
	ingredient.GET("/all", ingredients.GetIngredients(tokenAPI))                              
	ingredient.GET("/:id", ingredients.GetIngredientByID(tokenAPI))                           
	ingredient.POST("/", ingredients.PostIngredient(tokenAPI))                                
	ingredient.DELETE("/:id", ingredients.DeleteIngredient(tokenAPI))                         
	ingredient.PATCH("/:id", ingredients.UpdateIngredient(tokenAPI))                          
	ingredient.GET("/cookingspace/:id", ingredients.GetIngredientsByCookingSpaceID(tokenAPI)) 
	lesson := r.Group("/lesson")
	lesson.GET("/all", lessons.GetLessons(tokenAPI))                      // WORKING
	lesson.GET("/:id", lessons.GetLessonByID(tokenAPI))                   // WORKING
	lesson.GET("/group/:id", lessons.GetLessonsByGroupID(tokenAPI))       // WORKING
	lesson.POST("/:id", lessons.Postlesson(tokenAPI))         // WORKING               
	lesson.POST("/group/:id", lessons.AddLessonToAGroup(tokenAPI))        // WORKING
	lesson.DELETE("/group/:id", lessons.DeleteLessonFromAGroup(tokenAPI)) // WORKING
	lesson.PATCH("/:id", lessons.UpdateLesson(tokenAPI)) // WORKING
	lesson.DELETE("/:id", lessons.DeleteLesson(tokenAPI)) // WORKING
	food := r.Group("/food")
	food.GET("/all", foods.GetFoods(tokenAPI))      
	food.GET("/:id", foods.GetFoodByID(tokenAPI))   
	food.POST("/", foods.PostFood(tokenAPI))        
	food.DELETE("/:id", foods.DeleteFood(tokenAPI)) 
	food.PATCH("/:id", foods.UpdateFood(tokenAPI))  
	shopitem := r.Group("/shopitem")
	shopitem.GET("/all", shopitems.GetShopItems(tokenAPI))      
	shopitem.GET("/:id", shopitems.GetShopItemByID(tokenAPI))   
	shopitem.POST("/", shopitems.PostShopItem(tokenAPI))        
	shopitem.DELETE("/:id", shopitems.DeleteShopItem(tokenAPI)) 
	shopitem.PATCH("/:id", shopitems.UpdateShopItem(tokenAPI))  
	comment := r.Group("/comment")
	comment.GET("/event/:id", comments.GetCommentsByClientID(tokenAPI))   
	comment.GET("/client/:id", comments.GetCommentsByEventID(tokenAPI))   
	comment.POST("/", comments.PostComment(tokenAPI))        
	comment.DELETE("/:id", comments.DeleteComment(tokenAPI)) 
	comment.PATCH("/:id", comments.UpdateComment(tokenAPI))  
	bill := r.Group("/bill")
	bill.GET("/all", bills.GetBills(tokenAPI))      
	bill.GET("/:id", bills.GetBillByID(tokenAPI))   
	bill.GET("/user/:id", bills.GetBillsByUserID(tokenAPI))   
	bill.POST("/", bills.PostBill(tokenAPI))        
	bill.DELETE("/:id", bills.DeleteBill(tokenAPI)) 
	bill.PATCH("/:id", bills.UpdateBill(tokenAPI))  
	order := r.Group("/order")
	order.GET("/all", orders.GetOrders(tokenAPI))      
	order.GET("/:id", orders.GetOrder(tokenAPI))   
	order.GET("/chef/:id", orders.GetOrderByContractor1ID(tokenAPI))   
	order.GET("/deliveryman/:id", orders.GetOrderByContractor2ID(tokenAPI))   
	order.GET("/client/:id", orders.GetOrderByClientID(tokenAPI))   
	order.POST("/", orders.PostOrder(tokenAPI))        
	order.DELETE("/:id", orders.DeleteOrder(tokenAPI)) 
	order.PATCH("/:id", orders.UpdateOrder(tokenAPI))  
	order.PATCH("/item/:iditem/:idorder", orders.AddItemToAnOrder(tokenAPI))  
	order.DELETE("/item/:iditem/:idorder", orders.DeleteItemFromAnOrder(tokenAPI))  
	order.PATCH("/food/:idfood/:idorder", orders.AddFoodToAnOrder(tokenAPI))  
	order.DELETE("/food/:idfood/:idorder", orders.DeleteFoodFromAnOrder(tokenAPI))  
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
