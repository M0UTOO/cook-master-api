package orders

import (
	"cook-master-api/token"
	"cook-master-api/utils"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Order struct {
	IdOrder     int     `json:"idorder"`
	Status	  string  `json:"status"`
	Price 	 float64 `json:"price"`
	DeliveryAddress string `json:"deliveryaddress"`
	IdContractor1 int `json:"idcontractor1"`
	IdContractor2 int `json:"idcontractor2"`
	IdClient int `json:"idclient"`
}

func GetOrders(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		tokenHeader := c.Request.Header["Token"]
		if tokenHeader == nil {
			c.JSON(498, gin.H{
				"error":   true,
				"message": "missing token",
			})
		}

		err := token.CheckAPIToken(tokenAPI, tokenHeader[0], c)
		if err != nil {
			c.JSON(498, gin.H{
				"error":   true,
				"message": "wrong token",
			})
			return
		}

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		rows, err := db.Query("SELECT * FROM ORDERS")
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't query database",
			})
			return
		}

		var orders []Order
		for rows.Next() {
			var order Order
			err := rows.Scan(&order.IdOrder, &order.Status, &order.Price, &order.DeliveryAddress, &order.IdContractor1, &order.IdContractor2, &order.IdClient)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "can't scan rows",
				})
				return
			}
			orders = append(orders, order)
		}

		c.JSON(200, orders)
	}
}

func GetOrder(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		tokenHeader := c.Request.Header["Token"]
		if tokenHeader == nil {
			c.JSON(498, gin.H{
				"error":   true,
				"message": "missing token",
			})
		}

		err := token.CheckAPIToken(tokenAPI, tokenHeader[0], c)
		if err != nil {
			c.JSON(498, gin.H{
				"error":   true,
				"message": "wrong token",
			})
			return
		}

		id := c.Param("id")
		if id == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "missing id",
			})
			return
		}

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		row := db.QueryRow("SELECT * FROM ORDERS WHERE Id_ORDERS = ?", id)

		var order Order
		err = row.Scan(&order.IdOrder, &order.Status, &order.Price, &order.DeliveryAddress, &order.IdContractor1, &order.IdContractor2, &order.IdClient)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "order not found",
			})
			return
		}

		c.JSON(200, order)
	}
}

func PostOrder(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		tokenHeader := c.Request.Header["Token"]
		if tokenHeader == nil {
			c.JSON(498, gin.H{
				"error":   true,
				"message": "missing token",
			})
		}

		err := token.CheckAPIToken(tokenAPI, tokenHeader[0], c)
		if err != nil {
			c.JSON(498, gin.H{
				"error":   true,
				"message": "wrong token",
			})
			return
		}

		var order Order
		err = c.BindJSON(&order)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "can't bind json",
			})
			return
		}

		if order.Status == "" || order.Price < 0 || order.IdContractor1 == 0 || order.IdContractor2 == 0 || order.IdClient == 0 {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "missing parameters",
			})
			return
		}

		if !utils.IsSafeString(order.Status) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "invalid characters",
			})
			return
		}

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		var idcontractor1 int

		err = db.QueryRow("SELECT Id_CONTRACTORS FROM CONTRACTORS WHERE Id_USERS = '" + strconv.Itoa(order.IdContractor1) + "'").Scan(&idcontractor1)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "contractor1 doesn't exist",
			})
			return
		}

		var idcontractor2 int

		err = db.QueryRow("SELECT Id_CONTRACTORS FROM CONTRACTORS WHERE Id_USERS = '" + strconv.Itoa(order.IdContractor2) + "'").Scan(&idcontractor2)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "contractor2 doesn't exist",
			})
			return
		}

		var idclient int

		err = db.QueryRow("SELECT Id_CLIENTS FROM CLIENTS WHERE Id_USERS = '" + strconv.Itoa(order.IdClient) + "'").Scan(&idclient)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "client doesn't exist",
			})
			return
		}

		rows, err := db.Exec("INSERT INTO ORDERS (status, price, deliveryaddress, Id_CONTRACTORS, Id_CONTRACTORS_1, Id_CLIENTS) VALUES (?, ?, DEFAULT, ?, ?, ?)", order.Status, order.Price, idcontractor1, idcontractor2, idclient)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't insert into database",
			})
			return
		}

		id, err := rows.LastInsertId()

		c.JSON(200, gin.H{
			"error": false,
			"id": id,
			"message": "order created",
		})
	}
}

func UpdateOrder(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		tokenHeader := c.Request.Header["Token"]
		if tokenHeader == nil {
			c.JSON(498, gin.H{
				"error":   true,
				"message": "missing token",
			})
		}

		err := token.CheckAPIToken(tokenAPI, tokenHeader[0], c)
		if err != nil {
			c.JSON(498, gin.H{
				"error":   true,
				"message": "wrong token",
			})
			return
		}

		id := c.Param("id")
		if id == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "missing id",
			})
			return
		}

		var order Order

		order.Price = -1

		err = c.BindJSON(&order)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "can't bind json",
			})
			return
		}

		var setClause []string

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		if order.Status != "" {
			if !utils.IsSafeString(order.Status) {
				c.JSON(400, gin.H{
					"error":   true,
					"message": "invalid characters",
				})
				return
			}
			if len(order.Status) < 0 || len(order.Status) > 20 {
				c.JSON(400, gin.H{
					"error":   true,
					"message": "status must be between 0 and 20 characters",
				})
				return
			}
			setClause = append(setClause, "status = '"+order.Status+"'")
		}

		if order.Price >= 0 {
			setClause = append(setClause, "price = '"+strconv.FormatFloat(order.Price, 'f', 2, 64)+"'")
		}

		if order.DeliveryAddress != "" {
			if !utils.IsSafeString(order.DeliveryAddress) {
				c.JSON(400, gin.H{
					"error":   true,
					"message": "invalid characters",
				})
				return
			}
			if len(order.DeliveryAddress) < 0 || len(order.DeliveryAddress) > 50 {
				c.JSON(400, gin.H{
					"error":   true,
					"message": "delivery address must be between 0 and 50 characters",
				})
				return
			}
			setClause = append(setClause, "deliveryaddress = '"+order.DeliveryAddress+"'")
		}

		if order.IdContractor1 > 0 {
			var idcontractor1 int
			err = db.QueryRow("SELECT Id_CONTRACTORS FROM CONTRACTORS WHERE Id_CONTRACTORS = '" + strconv.Itoa(order.IdContractor1) + "'").Scan(&idcontractor1)
			if err != nil {
				c.JSON(400, gin.H{
					"error":   true,
					"message": "contractor1 doesn't exist",
				})
				return
			}
			setClause = append(setClause, "Id_CONTRACTORS = '"+strconv.Itoa(order.IdContractor1)+"'")
		}

		if order.IdContractor2 > 0 {
			var idcontractor2 int
			err = db.QueryRow("SELECT Id_CONTRACTORS FROM CONTRACTORS WHERE Id_CONTRACTORS = '" + strconv.Itoa(order.IdContractor2) + "'").Scan(&idcontractor2)
			if err != nil {
				c.JSON(400, gin.H{
					"error":   true,
					"message": "contractor2 doesn't exist",
				})
				return
			}
			setClause = append(setClause, "ID_CONTRACTORS_1 = '"+strconv.Itoa(order.IdContractor2)+"'")
		}

		if len(setClause) == 0 {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "no field to update",
			})
			return
		}

		var idorder int

		err = db.QueryRow("SELECT Id_ORDERS FROM ORDERS WHERE Id_ORDERS = '" + id + "'").Scan(&idorder)
			if err != nil {
				c.JSON(400, gin.H{
					"error":   true,
					"message": "order doesn't exist",
				})
				return
			}

		_, err = db.Exec("UPDATE ORDERS SET "+strings.Join(setClause, ", ")+" WHERE Id_ORDERS = ?", id)
		fmt.Println(err)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't update database",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
			"message": "order updated",
		})
	}
}

func DeleteOrder(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		tokenHeader := c.Request.Header["Token"]
		if tokenHeader == nil {
			c.JSON(498, gin.H{
				"error":   true,
				"message": "missing token",
			})
		}

		err := token.CheckAPIToken(tokenAPI, tokenHeader[0], c)
		if err != nil {
			c.JSON(498, gin.H{
				"error":   true,
				"message": "wrong token",
			})
			return
		}

		id := c.Param("id")
		if id == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "missing id",
			})
			return
		}

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		_, err = db.Exec("DELETE FROM CONTAINS_ITEM WHERE Id_ORDERS = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't delete from database",
			})
			return
		}

		_, err = db.Exec("DELETE FROM CONTAINS_FOOD WHERE Id_ORDERS = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't delete from database",
			})
			return
		}

		_, err = db.Exec("DELETE FROM ORDERS WHERE Id_ORDERS = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't delete from database",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
			"message": "order deleted",
		})
	}
}

func GetOrderByContractor1ID(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		tokenHeader := c.Request.Header["Token"]
		if tokenHeader == nil {
			c.JSON(498, gin.H{
				"error":   true,
				"message": "missing token",
			})
		}

		err := token.CheckAPIToken(tokenAPI, tokenHeader[0], c)
		if err != nil {
			c.JSON(498, gin.H{
				"error":   true,
				"message": "wrong token",
			})
			return
		}
		id := c.Param("id")
		if id == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "missing id",
			})
			return
		}
		var orders []Order
		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()
		rows, err := db.Query("SELECT * FROM ORDERS WHERE Id_CONTRACTORS = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't get orders from database",
			})
			return
		}
		defer rows.Close()
		for rows.Next() {
			var order Order
			err := rows.Scan(&order.IdOrder, &order.Status, &order.Price, &order.DeliveryAddress, &order.IdContractor1, &order.IdContractor2, &order.IdClient)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "chef don't have orders",
				})
				return
			}
			orders = append(orders, order)
		}
		c.JSON(200, orders)
	}
}

func GetOrderByContractor2ID(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		tokenHeader := c.Request.Header["Token"]
		if tokenHeader == nil {
			c.JSON(498, gin.H{
				"error":   true,
				"message": "missing token",
			})
		}

		err := token.CheckAPIToken(tokenAPI, tokenHeader[0], c)
		if err != nil {
			c.JSON(498, gin.H{
				"error":   true,
				"message": "wrong token",
			})
			return
		}
		id := c.Param("id")
		if id == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "missing id",
			})
			return
		}
		var orders []Order
		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()
		rows, err := db.Query("SELECT * FROM ORDERS WHERE Id_CONTRACTORS_1 = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't get orders from database",
			})
			return
		}
		defer rows.Close()
		for rows.Next() {
			var order Order
			err := rows.Scan(&order.IdOrder, &order.Status, &order.Price, &order.DeliveryAddress, &order.IdContractor1, &order.IdContractor2, &order.IdClient)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "delivery man don't have orders",
				})
				return
			}
			orders = append(orders, order)
		}
		c.JSON(200, orders)
	}
}

func GetOrderByStatus(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		tokenHeader := c.Request.Header["Token"]
		if tokenHeader == nil {
			c.JSON(498, gin.H{
				"error":   true,
				"message": "missing token",
			})
		}

		err := token.CheckAPIToken(tokenAPI, tokenHeader[0], c)
		if err != nil {
			c.JSON(498, gin.H{
				"error":   true,
				"message": "wrong token",
			})
			return
		}
		status := c.Param("status")
		if status == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "missing status",
			})
			return
		}
		var orders []Order
		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()
		rows, err := db.Query("SELECT * FROM ORDERS WHERE Status = ?", status)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't get orders from database",
			})
			return
		}
		defer rows.Close()
		for rows.Next() {
			var order Order
			err := rows.Scan(&order.IdOrder, &order.Status, &order.Price, &order.DeliveryAddress, &order.IdContractor1, &order.IdContractor2, &order.IdClient)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "no order with this status",
				})
				return
			}
			orders = append(orders, order)
		}
		c.JSON(200, orders)
	}
}

func GetOrderByClientID(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		tokenHeader := c.Request.Header["Token"]
		if tokenHeader == nil {
			c.JSON(498, gin.H{
				"error":   true,
				"message": "missing token",
			})
		}

		err := token.CheckAPIToken(tokenAPI, tokenHeader[0], c)
		if err != nil {
			c.JSON(498, gin.H{
				"error":   true,
				"message": "wrong token",
			})
			return
		}
		id := c.Param("id")
		if id == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "missing id",
			})
			return
		}
		var orders []Order
		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()
		rows, err := db.Query("SELECT * FROM ORDERS WHERE Id_CLIENTS = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't get orders from database",
			})
			return
		}
		defer rows.Close()
		for rows.Next() {
			var order Order
			err := rows.Scan(&order.IdOrder, &order.Status, &order.Price, &order.DeliveryAddress, &order.IdContractor1, &order.IdContractor2, &order.IdClient)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "client don't have orders",
				})
				return
			}
			orders = append(orders, order)
		}
		c.JSON(200, orders)
	}
}

func AddItemToAnOrder(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		tokenHeader := c.Request.Header["Token"]
		if tokenHeader == nil {
			c.JSON(498, gin.H{
				"error":   true,
				"message": "missing token",
			})
		}

		err := token.CheckAPIToken(tokenAPI, tokenHeader[0], c)
		if err != nil {
			c.JSON(498, gin.H{
				"error":   true,
				"message": "wrong token",
			})
			return
		}

		idOrder := c.Param("idorder")
		if idOrder == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "missing idOrder",
			})
			return
		}

		idItem := c.Param("iditem")
		if idItem == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "missing idItem",
			})
			return
		}

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		var idInt int

		err = db.QueryRow("SELECT Id_SHOP_ITEMS FROM SHOP_ITEMS WHERE Id_SHOP_ITEMS = '" + idItem + "'").Scan(&idInt)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "item doesn't exist",
			})
			return
		}

		err = db.QueryRow("SELECT Id_ORDERS FROM ORDERS WHERE Id_ORDERS = '" + idOrder + "'").Scan(&idInt)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "order doesn't exist",
			})
			return
		}

		_, err = db.Exec("INSERT INTO CONTAINS_ITEM (Id_ORDERS, Id_SHOP_ITEMS) VALUES (?, ?)", idOrder, idItem)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't insert item into database",
			})
			return
		}
		c.JSON(200, gin.H{
			"error": false,
			"message": "item added to order",
		})
	}
}

func AddFoodToAnOrder(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		tokenHeader := c.Request.Header["Token"]
		if tokenHeader == nil {
			c.JSON(498, gin.H{
				"error":   true,
				"message": "missing token",
			})
		}

		err := token.CheckAPIToken(tokenAPI, tokenHeader[0], c)
		if err != nil {
			c.JSON(498, gin.H{
				"error":   true,
				"message": "wrong token",
			})
			return
		}

		idOrder := c.Param("idorder")
		if idOrder == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "missing idOrder",
			})
			return
		}

		idFood := c.Param("idfood")
		if idFood == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "missing idFood",
			})
			return
		}

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		var idInt int

		err = db.QueryRow("SELECT Id_FOODS FROM FOODS WHERE Id_FOODS = '" + idFood + "'").Scan(&idInt)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "food doesn't exist",
			})
			return
		}

		err = db.QueryRow("SELECT Id_ORDERS FROM ORDERS WHERE Id_ORDERS = '" + idOrder + "'").Scan(&idInt)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "order doesn't exist",
			})
			return
		}

		_, err = db.Exec("INSERT INTO CONTAINS_FOOD (Id_ORDERS, Id_FOODS) VALUES (?, ?)", idOrder, idFood)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't insert food into database",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
			"message": "food added to order",
		})
	}
}

func DeleteItemFromAnOrder(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		tokenHeader := c.Request.Header["Token"]
		if tokenHeader == nil {
			c.JSON(498, gin.H{
				"error":   true,
				"message": "missing token",
			})
		}

		err := token.CheckAPIToken(tokenAPI, tokenHeader[0], c)
		if err != nil {
			c.JSON(498, gin.H{
				"error":   true,
				"message": "wrong token",
			})
			return
		}

		idOrder := c.Param("idorder")
		if idOrder == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "missing idOrder",
			})
			return
		}

		idItem := c.Param("iditem")
		if idItem == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "missing idItem",
			})
			return
		}

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		var idInt int

		err = db.QueryRow("SELECT Id_SHOP_ITEMS FROM SHOP_ITEMS WHERE Id_SHOP_ITEMS = '" + idItem + "'").Scan(&idInt)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "item doesn't exist",
			})
			return
		}

		err = db.QueryRow("SELECT Id_ORDERS FROM ORDERS WHERE Id_ORDERS = '" + idOrder + "'").Scan(&idInt)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "order doesn't exist",
			})
			return
		}

		_, err = db.Exec("DELETE FROM CONTAINS_ITEM WHERE Id_ORDERS = '" + idOrder + "' AND Id_SHOP_ITEMS = '" + idItem + "'")
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't delete item from database",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
			"message": "item deleted from order",
		})
	}
}

func DeleteFoodFromAnOrder(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		tokenHeader := c.Request.Header["Token"]
		if tokenHeader == nil {
			c.JSON(498, gin.H{
				"error":   true,
				"message": "missing token",
			})
		}

		err := token.CheckAPIToken(tokenAPI, tokenHeader[0], c)
		if err != nil {
			c.JSON(498, gin.H{
				"error":   true,
				"message": "wrong token",
			})
			return
		}

		idOrder := c.Param("idorder")
		if idOrder == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "missing idOrder",
			})
			return
		}

		idFood := c.Param("idfood")
		if idFood == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "missing idFood",
			})
			return
		}

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		var idInt int

		err = db.QueryRow("SELECT Id_FOODS FROM FOODS WHERE Id_FOODS = '" + idFood + "'").Scan(&idInt)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "food doesn't exist",
			})
			return
		}

		err = db.QueryRow("SELECT Id_ORDERS FROM ORDERS WHERE Id_ORDERS = '" + idOrder + "'").Scan(&idInt)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "order doesn't exist",
			})
			return
		}

		_, err = db.Exec("DELETE FROM CONTAINS_FOOD WHERE Id_ORDERS = '" + idOrder + "' AND Id_FOODS = '" + idFood + "'")
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't delete food from database",
			})
			return
		}
		
		c.JSON(200, gin.H{
			"error": false,
			"message": "food deleted from order",
		})
	}
}

func GetItemsByOrderID(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		type Item struct {
			IdItem int `json:"iditem"`
			Name string `json:"name"`
			Price float64 `json:"price"`
		}

		id := c.Param("id")
		if id == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "missing id",
			})
			return
		}
		var items []Item
		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot connect to database",
			})
			return
		}
		defer db.Close()
		rows, err := db.Query("SELECT Id_SHOP_ITEMS, name, price FROM SHOP_ITEMS WHERE Id_SHOP_ITEMS IN (SELECT Id_SHOP_ITEMS FROM CONTAINS_ITEM WHERE Id_ORDERS = ?)", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't get items from database",
			})
			return
		}
		defer rows.Close()
		for rows.Next() {
			var item Item
			err := rows.Scan(&item.IdItem, &item.Name, &item.Price)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "no item in this order",
				})
				return
			}
			items = append(items, item)
		}
		c.JSON(200, items)
	}
}

func GetFoodsByOrderID(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		type Food struct {
			IdFood int `json:"idfood"`
			Name string `json:"name"`
			Price float64 `json:"price"`
		}

		id := c.Param("id")
		if id == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "missing id",
			})
			return
		}
		var foods []Food
		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot connect to database",
			})
			return
		}
		defer db.Close()
		rows, err := db.Query("SELECT Id_FOODS, name, price FROM FOODS WHERE Id_FOODS IN (SELECT Id_FOODS FROM CONTAINS_FOOD WHERE Id_ORDERS = ?)", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't get foods from database",
			})
			return
		}
		defer rows.Close()
		for rows.Next() {
			var food Food
			err := rows.Scan(&food.IdFood, &food.Name, &food.Price)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "no food in this order",
				})
				return
			}
			foods = append(foods, food)
		}
		c.JSON(200, foods)
	}
}