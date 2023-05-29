package shopitems

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

type ShopItem struct {
	IdShopItem  int     `json:"idshopitem"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Reward      string  `json:"reward"`
	Picture	 string  `json:"picture"`
}

func GetShopItems(tokenAPI string) func(c *gin.Context) {
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

		rows, err := db.Query("SELECT * FROM SHOP_ITEMS")
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't query database",
			})
			return
		}
		defer rows.Close()

		var shopitems []ShopItem

		for rows.Next() {
			var shopitem ShopItem
			err := rows.Scan(&shopitem.IdShopItem, &shopitem.Name, &shopitem.Description, &shopitem.Price, &shopitem.Stock, &shopitem.Reward, &shopitem.Picture)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "can't scan row",
				})
				return
			}
			shopitems = append(shopitems, shopitem)
		}

		c.JSON(200, shopitems)
	}
}

func GetShopItemByID(tokenAPI string) func(c *gin.Context) {
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
				"message": "id can't be empty",
			})
			return
		}

		if !utils.IsSafeString(id) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id can't contain sql injection",
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

		var shopitem ShopItem

		err = db.QueryRow("SELECT * FROM SHOP_ITEMS WHERE Id_SHOP_ITEMS = "+id).Scan(&shopitem.IdShopItem, &shopitem.Name, &shopitem.Description, &shopitem.Price, &shopitem.Stock, &shopitem.Reward, &shopitem.Picture)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "shopitem not found",
			})
			return
		}

		c.JSON(200, shopitem)
		return

	}
}

func PostShopItem(tokenAPI string) func(c *gin.Context) {
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

		var shopitem ShopItem
		err = c.BindJSON(&shopitem)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "can't bind json",
			})
			return
		}

		if shopitem.Name == "" || !utils.IsSafeString(shopitem.Name) || len(shopitem.Name) < 0 || len(shopitem.Name) > 100 {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "name can't be empty, contain sql injection or wrong lenght",
			})
			return
		}

		if shopitem.Price < 0 {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "price can't be negative",
			})
			return
		}

		if shopitem.Description == "" || !utils.IsSafeString(shopitem.Description) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "description can't be empty or contain sql injection",
			})
			return
		}

		if shopitem.Stock < 0 {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "stock can't be negative",
			})
			return
		}

		if shopitem.Reward == "" || !utils.IsSafeString(shopitem.Reward) || len(shopitem.Reward) < 0 || len(shopitem.Reward) > 50 {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "reward can't be empty, contain sql injection or wrong lenght",
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

		var id int

		err = db.QueryRow("SELECT Id_SHOP_ITEMS FROM SHOP_ITEMS WHERE name = ?", shopitem.Name).Scan(&id)
		if err == nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "shopitem already exist",
			})
			return
		}

		_, err = db.Exec("INSERT INTO SHOP_ITEMS (name, description, price, stock, reward, picture) VALUES (?, ?, ?, ?, ?, DEFAULT)", shopitem.Name, shopitem.Description, shopitem.Price, shopitem.Stock, shopitem.Reward)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't insert into database",
			})
			return
		}

		c.JSON(200, gin.H{
			"error":   false,
			"message": "shopitem added",
		})
	}
}

func DeleteShopItem(tokenAPI string) func(c *gin.Context) {
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
			c.JSON(500, gin.H{
				"error":   true,
				"message": "missing id",
			})
			return
		}

		if !utils.IsSafeString(id) {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "unsafe string",
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

		_, err = db.Exec("DELETE FROM SHOP_ITEMS WHERE Id_SHOP_ITEMS = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't delete from database",
			})
			return
		}

		c.JSON(200, gin.H{
			"error":   false,
			"message": "shopitem deleted",
		})
	}
}

func UpdateShopItem(tokenAPI string) func(c *gin.Context) {
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

		id := c.Param("id")

		if id == "" {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "missing id",
			})
			return
		}

		if !utils.IsSafeString(id) {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "unsafe string",
			})
			return
		}

		var req ShopItem

		err = c.BindJSON(&req)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't decode json request",
			})
			return
		}

		var setClause []string

		if req.Name != "" {
			if !utils.IsSafeString(req.Name) {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "unsafe string",
				})
				return
			}
			if len(req.Name) < 0 || len(req.Name) > 100 {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "name must be between 0 and 100 characters",
				})
				return
			}
			setClause = append(setClause, "name = '"+req.Name+"'")
		}

		if req.Description != "" {
			if !utils.IsSafeString(req.Description) {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "unsafe string",
				})
				return
			}
			setClause = append(setClause, "description = '"+req.Description+"'")
		}

		if req.Price > 0 {
			setClause = append(setClause, "price = '"+strconv.FormatFloat(req.Price, 'f', 2, 64)+"'")
		}

		if req.Stock > 0 {
			setClause = append(setClause, "stock = '"+strconv.Itoa(req.Stock)+"'")
		}

		if req.Reward != "" {
			if !utils.IsSafeString(req.Reward) {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "unsafe string",
				})
				return
			}
			if len(req.Reward) < 0 || len(req.Reward) > 50 {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "reward must be between 0 and 50 characters",
				})
				return
			}
			setClause = append(setClause, "reward = '"+req.Reward+"'")
		}

		if req.Picture != "" {
			if !utils.IsSafeString(req.Picture) {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "unsafe string",
				})
				return
			}
			if len(req.Picture) < 0 || len(req.Picture) > 255 {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "picture must be between 0 and 255 characters",
				})
				return
			}
			setClause = append(setClause, "picture = '"+req.Picture+"'")
		}

		if len(setClause) == 0 {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "missing field",
			})
			return
		}

		_, err = db.Exec("UPDATE SHOP_ITEMS SET "+strings.Join(setClause, ", ")+" WHERE Id_SHOP_ITEMS = ?", id)
		fmt.Println(err)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't update database",
			})
			return
		}

		c.JSON(200, gin.H{
			"error":   false,
			"message": "shopitem updated",
		})
	}
}
