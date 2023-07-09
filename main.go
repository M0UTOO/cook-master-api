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
	"cook-master-api/utils"
	"cook-master-api/languages"
	"cook-master-api/messages"
	"cook-master-api/subscriptions"
	"github.com/gin-gonic/gin"
	// "github.com/unrolled/secure"
)

func main() {
	tokenAPI := token.GetAPIToken()
	//gin.SetMode(gin.ReleaseMode)
	// secureFunc := func() gin.HandlerFunc {
	// 	return func(c *gin.Context) {
	// 		secureMiddleware := secure.New(secure.Options{
	// 			SSLRedirect: true,
	// 			SSLHost:     "api.becomeacookmaster.live:9000",
	// 		})
	// 		err := secureMiddleware.Process(c.Writer, c.Request)

	// 		// If there was an error, do not continue.
	// 		if err != nil {
	// 			return
	// 		}

	// 		c.Next()
	// 	}
	// }()
	r := gin.Default()
	r.Use(utils.CorsMiddleware())
	// r.Use(secureFunc)
	r.GET("/", index)
	user := r.Group("/user")
	user.GET("/:id", users.GetUserByID(tokenAPI))                										// WORKING
	user.GET("/search/:filter", users.GetUserByFilter(tokenAPI)) 										// WORKING
	user.PATCH("/:id", users.UpdateUser(tokenAPI))               										// WORKING
	user.GET("/all", users.GetUsers(tokenAPI))                   										// WORKING
	user.POST("/", users.PostUser(tokenAPI))                     										// WORKING
	user.POST("/login", users.LoginUser(tokenAPI)) 				 										// WORKING
	user.DELETE("/:id", users.DeleteUser(tokenAPI))              										// WORKING
	user.POST("/password", users.GetPasswordByEmail(tokenAPI)) 											// WORKING
	client := r.Group("/client")
	client.GET("/all", clients.GetClients(tokenAPI))       												// WORKING                                             
	client.GET("/:id", clients.GetClientByID(tokenAPI))             									// WORKING                                    
	client.PATCH("/:id", clients.UpdateClient(tokenAPI))            									// WORKING                          
	client.PATCH("/subscription/:iduser/:idsubscription", clients.UpdateClientSubscription(tokenAPI))   // WORKING
	client.PATCH("/watch/:idclient/:idlesson", clients.WatchLesson(tokenAPI))  							// WORKING                           
	client.DELETE("/watch/:idclient/:idlesson", clients.UnWatchLesson(tokenAPI))     					// WORKING     
	client.GET("/subscription", clients.GetAllSubscription(tokenAPI)) 									// WORKING          
	client.GET("/participation", clients.GetAverageClientParticipationByMonth(tokenAPI)) 				// WORKING
	client.GET("/money", clients.GetAverageMoneySpentByClient(tokenAPI)) 								// WORKING
	client.GET("/country", clients.GetClientsByCountry(tokenAPI)) 										// WORKING
	client.GET("/top5", clients.GetTop5Client(tokenAPI)) 												// WORKING
	contractor := r.Group("/contractor")
	contractor.GET("/all", contractors.GetContractors(tokenAPI))    									// WORKING 
	contractor.GET("/:id", contractors.GetContractorByID(tokenAPI))  									// WORKING
	contractor.PATCH("/:id", contractors.UpdateContractor(tokenAPI))  									// WORKING
	contractor.POST("/type", contractors.AddAContractorType(tokenAPI)) 									// WORKING
	contractor.DELETE("/type/:id", contractors.DeleteAContractorType(tokenAPI)) 						// WORKING
	contractor.GET("/type", contractors.GetContractorTypes(tokenAPI)) 									// WORKING
	contractor.GET("/role/:search", contractors.GetContractorByRoles(tokenAPI))							// WORKING
	contractor.GET("/type/:id", contractors.GetCreatorRoleById(tokenAPI))							    // WORKING
	manager := r.Group("/manager")
	manager.GET("/all", managers.GetManagers(tokenAPI))      											// WORKING
	manager.GET("/:id", managers.GetManagerByID(tokenAPI))   											// WORKING
	manager.PATCH("/:id", managers.UpdateManager(tokenAPI))  											// WORKING
	subscription := r.Group("/subscription")
	subscription.GET("/all", subscriptions.GetSubscriptions(tokenAPI))      							// WORKING
	subscription.GET("/:id", subscriptions.GetSubscriptionByID(tokenAPI))    							// WORKING
	subscription.POST("/", subscriptions.PostSubscription(tokenAPI))       								// WORKING 
	subscription.DELETE("/:id", subscriptions.DeleteSubscription(tokenAPI))  							// WORKING
	subscription.PATCH("/:id", subscriptions.UpdateSubscription(tokenAPI))  							// WORKING
	event := r.Group("/event")
	event.GET("/all", events.GetEvents(tokenAPI))            			 								// WORKING                                  
	event.GET("/:id", events.GetEventByID(tokenAPI))                        							// WORKING                     
	event.GET("/group/:id", events.GetEventsByGroupID(tokenAPI))            							// WORKING                     
	event.POST("/:id", events.PostEvent(tokenAPI))                        								// WORKING                       
	event.POST("/group/:id", events.AddEventToAGroup(tokenAPI))                              			// WORKING    
	event.DELETE("/group/:id", events.DeleteEventFromAGroup(tokenAPI))                    				// WORKING       
	event.PATCH("/:id", events.UpdateEvent(tokenAPI))                                            		// WORKING
	event.GET("/animate/:idevent/:iduser", events.AddContractorToAnEvent(tokenAPI))              		// WORKING
	event.DELETE("/animate/:idevent/:iduser", events.DeleteContractorFromAnEvent(tokenAPI))      		// WORKING
	event.GET("/animate/:idevent", events.GetContractorsByEventID(tokenAPI))                     		// WORKING
	event.GET("/organize/:idevent", events.GetManagersByEventID(tokenAPI))                       		// WORKING
	event.GET("/participate/:idevent", events.GetClientsByEventID(tokenAPI))                      		// WORKING
	event.GET("/groups/:idevent", events.GetGroupsByEventID(tokenAPI))                             		// WORKING
	event.GET("/host/:idevent", events.GetCookingSpacesByEventID(tokenAPI))                       		// WORKING
	event.POST("/participate/:idevent/:iduser", events.AddClientToAnEvent(tokenAPI))              		// WORKING
	event.DELETE("/participate/:idevent/:iduser", events.DeleteClientFromAnEvent(tokenAPI))      		// WORKING
	event.PATCH("/participation/:idevent/:iduser", events.ValidateClientPresence(tokenAPI))        		// WORKING
	event.DELETE("/participation/:idevent/:iduser", events.UnvalidateClientPresence(tokenAPI))        		// WORKING
	event.PATCH("/host/:idevent/:idcookingspace", events.AddEventToAnCookingSpace(tokenAPI))     		// WORKING
	event.DELETE("/host/:idevent/:idcookingspace", events.DeleteEventToAnCookingSpace(tokenAPI)) 		// WORKING
	event.GET("/group/all", events.GetGroupEvents(tokenAPI)) 											// WORKING
	event.GET("/formation/:iduser", events.GetAllFormationsByUserID(tokenAPI))							// WORKING
	event.GET("/month", events.GetEventsByMonth(tokenAPI)) // WORKING
	event.GET("/months", events.GetEventsByMonthInAYear(tokenAPI)) // WORKING
	event.GET("/type", events.GetEventsByType(tokenAPI)) // WORKING
	event.GET("/week", events.GetEventsByDayOfTheWeek(tokenAPI)) // WORKING
	event.GET("/top5", events.GetTop5Events(tokenAPI)) // WORKING
	event.GET("/formation", events.GetFormationsDone(tokenAPI)) // WORKING
	event.POST("/search/:search", events.SearchForEvents(tokenAPI)) // WORKING
	event.GET("/rate/:id", events.GetRateByEventID(tokenAPI)) // WORKING
	event.GET("/comment/:id", events.GetEventComments(tokenAPI)) // WORKING
	event.GET("/coming/:id", events.GetComingEventByClientId(tokenAPI)) // WORKING
	event.GET("/past/:id", events.GetPastEventByClientIdfunc(tokenAPI)) // WORKING
	event.DELETE("/:id", events.DeleteEvent(tokenAPI)) // WORKING
	event.GET("/group/search/:search", events.SearchForEventsGroups(tokenAPI)) // WORKING
	event.GET("/animate/get/:id", events.GetEventsByUserId(tokenAPI)) // WORKING
	event.GET("/ispresent/:idevent/:idclient", events.GetClientParticipationToEvent(tokenAPI)) // WORKING
	event.GET("/formation/get/:id", events.GetFomationsForUser(tokenAPI)) // WORKING
	event.GET("/group/get/:id", events.GetGroupByGroupId(tokenAPI)) // WORKING
	premise := r.Group("/premise")
	premise.GET("/all", premises.GetPremises(tokenAPI))      // WORKING
	premise.GET("/:id", premises.GetPremiseByID(tokenAPI))   // WORKING
	premise.POST("/", premises.PostPremise(tokenAPI))        // WORKING
	premise.DELETE("/:id", premises.DeletePremise(tokenAPI)) // WORKING
	premise.PATCH("/:id", premises.UpdatePremise(tokenAPI))  // WORKING
	premise.GET("/books", premises.GetBooksByPremises(tokenAPI)) // WORKING
	premise.GET("/cookingspace/:id", premises.GetPremiseByCookingSpace(tokenAPI)) // WORKING
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
	cookingSpace.GET("/books/:id", cookingspaces.GetBooksByCookingSpaceID(tokenAPI)) // WORKING
	cookingSpace.GET("/books/all", cookingspaces.GetCookingSpacesBooks(tokenAPI)) // WORKING
	cookingSpace.DELETE("/:id", cookingspaces.DeleteCookingSpace(tokenAPI)) // WORKING
	cookingSpace.GET("/top5", cookingspaces.GetTop5CookingSpaces(tokenAPI)) // WORKING
	cookingSpace.GET("/event/:id", cookingspaces.GetEventsByCookingSpaceId(tokenAPI)) // WORKING
	cookingSpace.GET("/books/client/:id", cookingspaces.GetBooksByUserId(tokenAPI)) // WORKING
	cookingSpace.GET("/search/:search", cookingspaces.SearchForCookingSpaces(tokenAPI)) // WORKING
	cookingItem := r.Group("/cookingitem")
	cookingItem.GET("/all", cookingitems.GetCookingItems(tokenAPI))       // WORKING                       
	cookingItem.GET("/:id", cookingitems.GetCookingItemByID(tokenAPI))   // WORKING                        
	cookingItem.POST("/", cookingitems.PostCookingItem(tokenAPI))       // WORKING                         
	cookingItem.DELETE("/:id", cookingitems.DeleteCookingItem(tokenAPI))           // WORKING              
	cookingItem.PATCH("/:id", cookingitems.UpdateCookingItem(tokenAPI))          // WORKING                
	cookingItem.GET("/cookingspace/:id", cookingitems.GetCookingItemsByCookingSpaceID(tokenAPI)) // WORKING
	cookingItem.GET("/incookingspace/:id", cookingitems.GetCookingSpaceByCookingItemId(tokenAPI)) // WORKING
	ingredient := r.Group("/ingredient")
	ingredient.GET("/all", ingredients.GetIngredients(tokenAPI))                  // WORKING            
	ingredient.GET("/:id", ingredients.GetIngredientByID(tokenAPI))                         // WORKING   
	ingredient.POST("/", ingredients.PostIngredient(tokenAPI))                           // WORKING     
	ingredient.DELETE("/:id", ingredients.DeleteIngredient(tokenAPI))               // WORKING          
	ingredient.PATCH("/:id", ingredients.UpdateIngredient(tokenAPI))                  // WORKING        
	ingredient.GET("/cookingspace/:id", ingredients.GetIngredientsByCookingSpaceID(tokenAPI)) // WORKING
	ingredient.GET("/incookingspace/:id", ingredients.GetCookingSpaceByIngredientId(tokenAPI)) // WORKING
	lesson := r.Group("/lesson")
	lesson.GET("/all", lessons.GetLessons(tokenAPI))                      // WORKING
	lesson.GET("/:id", lessons.GetLessonByID(tokenAPI))                   // WORKING
	lesson.GET("/group/:id", lessons.GetLessonsByGroupID(tokenAPI))       // WORKING
	lesson.POST("/:id", lessons.Postlesson(tokenAPI))         // WORKING               
	lesson.POST("/group/:id", lessons.AddLessonToAGroup(tokenAPI))        // WORKING
	lesson.DELETE("/group/:id", lessons.DeleteLessonFromAGroup(tokenAPI)) // WORKING
	lesson.PATCH("/:id", lessons.UpdateLesson(tokenAPI)) // WORKING
	lesson.DELETE("/:id", lessons.DeleteLesson(tokenAPI)) // WORKING
	lesson.GET("/group/all", lessons.GetGroupLessons(tokenAPI)) // WORKING
	lesson.GET("/user/:id", lessons.GetUserIdByLessonId(tokenAPI)) // WORKING
	lesson.GET("/difficulty", lessons.GetLessonsWatchedByDifficulty(tokenAPI)) // WORKING
	lesson.GET("/group/get/:id", lessons.GetGroupByGroupId(tokenAPI)) // WORKING
	lesson.DELETE("/group/delete/:id", lessons.DeleteLessonGroup(tokenAPI)) // WORKING
	lesson.POST("/group/post", lessons.CreateLessonGroup(tokenAPI)) // WORKING
	lesson.GET("/suggested", lessons.GetSuggestedLessons(tokenAPI)) // WORKING
	lesson.DELETE("/views/:id", lessons.UpdateLessonViews(tokenAPI)) // WORKING
	lesson.GET("/views/:id", lessons.GetAllClientViews(tokenAPI)) // WORKING
	lesson.GET("/watch/:iduser/:idlesson", lessons.IsLessonWatched(tokenAPI)) // WORKING
	lesson.GET("/search/:search", lessons.SearchForLessons(tokenAPI)) // WORKING
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
	comment.GET("/:id", comments.GetCommentByCommentID(tokenAPI))   // WORKING
	comment.GET("/event/:id", comments.GetCommentsByEventID(tokenAPI))   // WORKING
	comment.GET("/client/:id", comments.GetCommentsByUserID(tokenAPI))   // WORKING
	comment.POST("/", comments.PostComment(tokenAPI))        // WORKING
	comment.DELETE("/:id", comments.DeleteComment(tokenAPI)) // WORKING
	comment.PATCH("/:id", comments.UpdateComment(tokenAPI))  // WORKING
	bill := r.Group("/bill")
	bill.GET("/all", bills.GetBills(tokenAPI))      // WORKING
	bill.GET("/:id", bills.GetBillByID(tokenAPI))   // WORKING
	bill.GET("/user/:id", bills.GetBillsByUserID(tokenAPI))   // WORKING
	bill.POST("/", bills.PostBill(tokenAPI))        // WORKING
	bill.DELETE("/:id", bills.DeleteBill(tokenAPI)) // WORKING
	bill.PATCH("/:id", bills.UpdateBill(tokenAPI))  // WORKING
	order := r.Group("/order")
	order.GET("/all", orders.GetOrders(tokenAPI))     // WORKING 
	order.GET("/:id", orders.GetOrder(tokenAPI))   // WORKING
	order.GET("/chef/:id", orders.GetOrderByContractor1ID(tokenAPI))   // WORKING
	order.GET("/deliveryman/:id", orders.GetOrderByContractor2ID(tokenAPI))   // WORKING
	order.GET("/client/:id", orders.GetOrderByClientID(tokenAPI))   // WORKING
	order.GET("/status/:status", orders.GetOrderByStatus(tokenAPI))   // WORKING
	order.POST("/", orders.PostOrder(tokenAPI))        // WORKING
	order.DELETE("/:id", orders.DeleteOrder(tokenAPI)) // WORKING
	order.PATCH("/:id", orders.UpdateOrder(tokenAPI))  // WORKING
	order.PATCH("/item/:iditem/:idorder", orders.AddItemToAnOrder(tokenAPI))  // WORKING
	order.DELETE("/item/:iditem/:idorder", orders.DeleteItemFromAnOrder(tokenAPI))  // WORKING
	order.PATCH("/food/:idfood/:idorder", orders.AddFoodToAnOrder(tokenAPI))  // WORKING
	order.DELETE("/food/:idfood/:idorder", orders.DeleteFoodFromAnOrder(tokenAPI))  // WORKING
	order.GET("/item/:id", orders.GetItemsByOrderID(tokenAPI))  // WORKING
	order.GET("/food/:id", orders.GetFoodsByOrderID(tokenAPI))  // WORKING
	order.GET("/month", orders.GetOrdersByMonth(tokenAPI))  // WORKING
	order.GET("/top5item", orders.GetTop5Items(tokenAPI))  // WORKING
	order.GET("/top5food", orders.GetTop5Food(tokenAPI))  // WORKING
	order.GET("/random", orders.GetRandomFoods(tokenAPI))  // WORKING
	language := r.Group("/language")
	language.GET("/all", languages.GetLanguages(tokenAPI))      // WORKING
	language.GET("/:id", languages.GetLanguageByID(tokenAPI))   // WORKING
	language.POST("/", languages.PostLanguage(tokenAPI))        // WORKING
	language.DELETE("/:id", languages.DeleteLanguage(tokenAPI)) // WORKING
	language.PATCH("/:id", languages.UpdateLanguage(tokenAPI))  // WORKING
	message := r.Group("/message")
	message.GET("/:idsender/:idreceiver", messages.GetMessageForIdSenderAndIdReceiver(tokenAPI))
	message.POST("/:idsender/:idreceiver", messages.PostMessage(tokenAPI))
	message.GET("/chief/:id", messages.GetAllClientForAContractorUserId(tokenAPI))
	r.Run(":9000")
	//r.RunTLS(":9000", "/home/debian/.ssh/api_certificate.pem", "/home/debian/.ssh/api_private_key.pem")
}

func index(c *gin.Context) {
	c.JSON(200, gin.H{
		"error":   true,
		"message": "hello world",
	})
}
