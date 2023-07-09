package cookingitems

import (
	"cook-master-api/token"
	"cook-master-api/utils"
	"database/sql"
	"strings"
	"strconv"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type CookingItem struct {
	IdCookingItem int `json:"idcookingitem"`
	Name string `json:"name"`
	Status string `json:"status"`
	IdCookingSpace int `json:"idcookingspace"`
}

func GetCookingItems(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		tokenHeader := c.Request.Header["Token"]
		if tokenHeader == nil{
			c.JSON(498, gin.H{
				"error": true,
				"message": "missing token",
			})
		}

		err := token.CheckAPIToken(tokenAPI, tokenHeader[0], c)
		if err != nil {
			c.JSON(498, gin.H{
				"error": true,
				"message": "wrong token",
			})
			return
		}

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		rows, err := db.Query("SELECT * FROM COOKING_ITEMS")
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cookingitems not found",
			})
			return
		}

		var cookingitems []CookingItem

		for rows.Next() {
			var cookingitem CookingItem
			err = rows.Scan(&cookingitem.IdCookingItem, &cookingitem.Name, &cookingitem.Status, &cookingitem.IdCookingSpace)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error": true,
					"message": "cookingitem not found",
				})
				return
			}
			cookingitems = append(cookingitems, cookingitem)
		}

		c.JSON(200, cookingitems)
		return
	}
}

func GetCookingItemByID(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header["Token"]
		if tokenHeader == nil{
			c.JSON(498, gin.H{
				"error": true,
				"message": "missing token",
			})
		}

		err := token.CheckAPIToken(tokenAPI, tokenHeader[0], c)
		if err != nil {
			c.JSON(498, gin.H{
				"error": true,
				"message": "wrong token",
			})
			return
		}

		id := c.Param("id")
		if id == "" {
			c.JSON(400, gin.H{
				"error": true,
				"message": "id can't be empty",
			})
			return
		}

		if !utils.IsSafeString(id) {
			c.JSON(400, gin.H{
				"error": true,
				"message": "id can't contain sql injection",
			})
			return
		}

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		var cookingitem CookingItem

		err = db.QueryRow("SELECT * FROM COOKING_ITEMS WHERE Id_COOKING_ITEMS = " + id).Scan(&cookingitem.IdCookingItem, &cookingitem.Name, &cookingitem.Status, &cookingitem.IdCookingSpace)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error": true,
				"message": "cookingitem not found",
			})
			return
		}

		c.JSON(200, cookingitem)
		return

	}
}

func UpdateCookingItem(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		tokenHeader := c.Request.Header["Token"]
		if tokenHeader == nil{
			c.JSON(498, gin.H{
				"error": true,
				"message": "missing token",
			})
		}

		err := token.CheckAPIToken(tokenAPI, tokenHeader[0], c)
		if err != nil {
			c.JSON(498, gin.H{
				"error": true,
				"message": "wrong token",
			})
			return
		}

		id := c.Param("id")
		if id == "" {
			c.JSON(400, gin.H{
				"error": true,
				"message": "id can't be empty",
			})
			return
		}

		if !utils.IsSafeString(id) {
			c.JSON(400, gin.H{
				"error": true,
				"message": "id can't contain sql injection",
			})
			return
		}

		var cookingitem CookingItem

		err = c.BindJSON(&cookingitem)
		if err != nil {
			c.JSON(400, gin.H{
				"error": true,
				"message": "bad json",
			})
			return
		}

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		var setClause []string
		
		if cookingitem.Name != "" {
			if !utils.IsSafeString(cookingitem.Name) {
				c.JSON(400, gin.H{
					"error": true,
					"message": "name can't contain sql injection",
				})
				return
			}
			if len(cookingitem.Name) < 0 || len(cookingitem.Name) > 100 {
				c.JSON(400, gin.H{
					"error": true,
					"message": "name must be between 0 and 100 characters",
				})
				return
			}
			setClause = append(setClause, "Name = '"+cookingitem.Name+"'")
		}

		if cookingitem.Status != "" {
			if !utils.IsSafeString(cookingitem.Status) {
				c.JSON(400, gin.H{
					"error": true,
					"message": "status can't contain sql injection",
				})
				return
			}
			if len(cookingitem.Status) < 0 || len(cookingitem.Status) > 50 {
				c.JSON(400, gin.H{
					"error": true,
					"message": "status must be between 0 and 50 characters",
				})
				return
			}
			setClause = append(setClause, "Status = '"+cookingitem.Status+"'")
		}

		if cookingitem.IdCookingSpace > 0 {
			var idcookingspace int
			err = db.QueryRow("SELECT Id_COOKING_SPACES FROM COOKING_SPACES WHERE Id_COOKING_SPACES = '" + strconv.Itoa(cookingitem.IdCookingSpace) + "'").Scan(&idcookingspace)
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "cookingspace not found",
				})
				return
			}
			setClause = append(setClause, "Id_COOKING_SPACES = '"+strconv.Itoa(cookingitem.IdCookingSpace)+"'")
		}
		
		if len(setClause) == 0 {
			c.JSON(400, gin.H{
				"error": true,
				"message": "nothing to update",
			})
			return
		}

		_, err = db.Exec("UPDATE COOKING_ITEMS SET " + strings.Join(setClause, ", ") + " WHERE Id_COOKING_ITEMS = " + id)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cookingitem not found",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
			"message": "cookingitem updated",
		})
		return
	}
}

func PostCookingItem(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		tokenHeader := c.Request.Header["Token"]
		if tokenHeader == nil{
			c.JSON(498, gin.H{
				"error": true,
				"message": "missing token",
			})
		}

		err := token.CheckAPIToken(tokenAPI, tokenHeader[0], c)
		if err != nil {
			c.JSON(498, gin.H{
				"error": true,
				"message": "wrong token",
			})
			return
		}

		var cookingitem CookingItem
		c.BindJSON(&cookingitem)

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		if cookingitem.Name == "" || !utils.IsSafeString(cookingitem.Name) || len(cookingitem.Name) > 100 || len(cookingitem.Name) < 0 {
			c.JSON(400, gin.H{
				"error": true,
				"message": "name can't be empty, above 100 char, below 0 char or contain sql injection",
			})
			return
		}
		if cookingitem.Status == "" || !utils.IsSafeString(cookingitem.Status) || len(cookingitem.Status) > 50 || len(cookingitem.Status) < 0 {
			c.JSON(400, gin.H{
				"error": true,
				"message": "status can't be empty, above 50 char, below 0 char or contain sql injection",
			})
			return
		}

		var idcookingspace int

		err = db.QueryRow("SELECT Id_COOKING_SPACES FROM COOKING_SPACES WHERE Id_COOKING_SPACES = '" + strconv.Itoa(cookingitem.IdCookingSpace) + "'").Scan(&idcookingspace)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cookingspace not found",
			})
			return
		}

		rows, err := db.Exec("INSERT INTO COOKING_ITEMS (Name, Status, Id_COOKING_SPACES) VALUES (?, ?, ?)", cookingitem.Name, cookingitem.Status, cookingitem.IdCookingSpace)
		fmt.Println(err)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot insert cookingitem",
			})
			return
		}

		id, err := rows.LastInsertId()

		c.JSON(200, gin.H{
			"error": false,
			"id": id,
			"message": "cookingitem inserted",
		})
	}
}

func DeleteCookingItem(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		tokenHeader := c.Request.Header["Token"]
		if tokenHeader == nil{
			c.JSON(498, gin.H{
				"error": true,
				"message": "missing token",
			})
		}

		err := token.CheckAPIToken(tokenAPI, tokenHeader[0], c)
		if err != nil {
			c.JSON(498, gin.H{
				"error": true,
				"message": "wrong token",
			})
			return
		}

		id := c.Param("id")
		if id == "" {
			c.JSON(400, gin.H{
				"error": true,
				"message": "missing id",
			})
			return
		}

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		_, err = db.Exec("DELETE FROM COOKING_ITEMS WHERE Id_COOKING_ITEMS = " + id)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cookingitem not found",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
			"message": "cookingitem deleted",
		})
	}
}

func GetCookingItemsByCookingSpaceID(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		tokenHeader := c.Request.Header["Token"]
		if tokenHeader == nil{
			c.JSON(498, gin.H{
				"error": true,
				"message": "missing token",
			})
		}

		err := token.CheckAPIToken(tokenAPI, tokenHeader[0], c)
		if err != nil {
			c.JSON(498, gin.H{
				"error": true,
				"message": "wrong token",
			})
			return
		}

		id := c.Param("id")
		if id == "" {
			c.JSON(400, gin.H{
				"error": true,
				"message": "missing id",
			})
			return
		}

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		rows, err := db.Query("SELECT * FROM COOKING_ITEMS WHERE Id_COOKING_SPACES = " + id)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cookingitem not found",
			})
			return
		}
		defer rows.Close()

		var cookingitems []CookingItem
		for rows.Next() {
			var cookingitem CookingItem
			err := rows.Scan(&cookingitem.IdCookingItem, &cookingitem.Name, &cookingitem.Status, &cookingitem.IdCookingSpace)
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "cookingitem not found",
				})
				return
			}
			cookingitems = append(cookingitems, cookingitem)
		}

		c.JSON(200, cookingitems)
	}
}

func GetCookingSpaceByCookingItemId(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		type CookingSpace struct {
			IdCookingSpace int     `json:"idCookingSpace"`
			Name           string  `json:"name"`
			Size           int     `json:"size"`
			IsAvailable    bool    `json:"isAvailable"`
			PricePerHour   float64 `json:"pricePerHour"`
			IdPremise      int     `json:"idPremise"`
			Picture 	   string  `json:"picture"`
		}

		tokenHeader := c.Request.Header["Token"]
		if tokenHeader == nil{
			c.JSON(498, gin.H{
				"error": true,
				"message": "missing token",
			})
		}

		err := token.CheckAPIToken(tokenAPI, tokenHeader[0], c)
		if err != nil {
			c.JSON(498, gin.H{
				"error": true,
				"message": "wrong token",
			})
			return
		}

		id := c.Param("id")
		if id == "" {
			c.JSON(400, gin.H{
				"error": true,
				"message": "missing id",
			})
			return
		}

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		rows, err := db.Query("SELECT COOKING_SPACES.Id_COOKING_SPACES, COOKING_SPACES.name, COOKING_SPACES.size, COOKING_SPACES.isavailable, COOKING_SPACES.priceperhour, COOKING_SPACES.Id_PREMISES, COOKING_SPACES.Picture FROM COOKING_SPACES JOIN COOKING_ITEMS ON COOKING_ITEMS.Id_COOKING_SPACES = COOKING_SPACES.Id_COOKING_SPACES WHERE Id_COOKING_ITEMS = " + id)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error": true,
				"message": "cookingspace not found",
			})
			return
		}
		defer rows.Close()

		var cookingspaces []CookingSpace
		for rows.Next() {
			var cookingspace CookingSpace
			err := rows.Scan(&cookingspace.IdCookingSpace, &cookingspace.Name, &cookingspace.Size, &cookingspace.IsAvailable, &cookingspace.PricePerHour, &cookingspace.IdPremise, &cookingspace.Picture)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error": true,
					"message": "cookingspace not found",
				})
				return
			}
			cookingspaces = append(cookingspaces, cookingspace)
		}

		c.JSON(200, cookingspaces)
	}
}