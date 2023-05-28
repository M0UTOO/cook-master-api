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

		row := db.QueryRow("SELECT * FROM ORDERS WHERE idorder = ?", id)

		var order Order
		err = row.Scan(&order.IdOrder, &order.Status, &order.Price, &order.DeliveryAddress, &order.IdContractor1, &order.IdContractor2, &order.IdClient)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't scan row",
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

		if order.Status == "" || order.Price < 0 || order.DeliveryAddress == "" || order.IdContractor1 == 0 || order.IdContractor2 == 0 || order.IdClient == 0 {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "missing parameters",
			})
			return
		}

		if !utils.IsSafeString(order.Status) || !utils.IsSafeString(order.DeliveryAddress) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "invalid characters",
			})
			return
		}

		var id int

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		err = db.QueryRow("SELECT Id_USERS FROM CONTRACTORS WHERE Id_USERS = '" + strconv.Itoa(order.IdContractor1) + "'").Scan(&id)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "contractor1 doesn't exist",
			})
			return
		}

		err = db.QueryRow("SELECT Id_USERS FROM CONTRACTORS WHERE Id_USERS = '" + strconv.Itoa(order.IdContractor2) + "'").Scan(&id)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "contractor2 doesn't exist",
			})
			return
		}

		err = db.QueryRow("SELECT Id_USERS FROM CLIENTS WHERE Id_USERS = '" + strconv.Itoa(order.IdClient) + "'").Scan(&id)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "client doesn't exist",
			})
			return
		}

		_, err = db.Exec("INSERT INTO ORDERS (status, price, deliveryaddress, idcontractor1, idcontractor2, idclient) VALUES (?, ?, ?, ?, ?, ?)", order.Status, order.Price, order.DeliveryAddress, order.IdContractor1, order.IdContractor2, order.IdClient)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't insert into database",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
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
		err = c.BindJSON(&order)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "can't bind json",
			})
			return
		}

		order.Price = -1

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

		if order.Price > -1 {
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
			var idInt int
			err = db.QueryRow("SELECT Id_USERS FROM CONTRACTORS WHERE Id_USERS = '" + strconv.Itoa(order.IdContractor1) + "'").Scan(&idInt)
			if err != nil {
				c.JSON(400, gin.H{
					"error":   true,
					"message": "contractor1 doesn't exist",
				})
				return
			}
			setClause = append(setClause, "idcontractor1 = '"+strconv.Itoa(order.IdContractor1)+"'")
		}

		if order.IdContractor2 > 0 {
			var idInt int
			err = db.QueryRow("SELECT Id_USERS FROM CONTRACTORS WHERE Id_USERS = '" + strconv.Itoa(order.IdContractor2) + "'").Scan(&idInt)
			if err != nil {
				c.JSON(400, gin.H{
					"error":   true,
					"message": "contractor2 doesn't exist",
				})
				return
			}
			setClause = append(setClause, "idcontractor2 = '"+strconv.Itoa(order.IdContractor2)+"'")
		}

		if len(setClause) == 0 {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "no field to update",
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
		rows, err := db.Query("SELECT * FROM ORDERS WHERE Id_CONTRACTOR1 = ?", id)
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
			err := rows.Scan(&order.IdOrder, &order.Status, &order.Price, &order.DeliveryAddress, &order.IdContractor1, &order.IdContractor2)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "can't get orders from database",
				})
				return
			}
			orders = append(orders, order)
		}
		c.JSON(200, gin.H{
			"error": false,
			"orders": orders,
		})
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
		rows, err := db.Query("SELECT * FROM ORDERS WHERE Id_CONTRACTOR2 = ?", id)
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
			err := rows.Scan(&order.IdOrder, &order.Status, &order.Price, &order.DeliveryAddress, &order.IdContractor1, &order.IdContractor2)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "can't get orders from database",
				})
				return
			}
			orders = append(orders, order)
		}
		c.JSON(200, gin.H{
			"error": false,
			"orders": orders,
		})
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
			err := rows.Scan(&order.IdOrder, &order.Status, &order.Price, &order.DeliveryAddress, &order.IdContractor1, &order.IdContractor2)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "can't get orders from database",
				})
				return
			}
			orders = append(orders, order)
		}
		c.JSON(200, gin.H{
			"error": false,
			"orders": orders,
		})
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
		rows, err := db.Query("SELECT * FROM ORDERS WHERE Id_CLIENT = ?", id)
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
			err := rows.Scan(&order.IdOrder, &order.Status, &order.Price, &order.DeliveryAddress, &order.IdContractor1, &order.IdContractor2)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "can't get orders from database",
				})
				return
			}
			orders = append(orders, order)
		}
		c.JSON(200, gin.H{
			"error": false,
			"orders": orders,
		})
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