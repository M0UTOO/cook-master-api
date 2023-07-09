package ingredients

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

type Ingredient struct {
	IdIngredient int `json:"idingredient"`
	Name string `json:"name"`
	Allergen string `json:"allergen"`
	IdCookingSpace int `json:"idcookingspace"`
}

func GetIngredients(tokenAPI string) func(c *gin.Context) {
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

		rows, err := db.Query("SELECT * FROM INGREDIENTS")
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "ingredients not found",
			})
			return
		}

		var ingredients []Ingredient

		for rows.Next() {
			var ingredient Ingredient
			err = rows.Scan(&ingredient.IdIngredient, &ingredient.Name, &ingredient.Allergen, &ingredient.IdCookingSpace)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error": true,
					"message": "ingredient not found",
				})
				return
			}
			ingredients = append(ingredients, ingredient)
		}

		c.JSON(200, ingredients)
		return
	}
}

func GetIngredientByID(tokenAPI string) func(c *gin.Context) {
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

		var ingredient Ingredient

		err = db.QueryRow("SELECT * FROM INGREDIENTS WHERE Id_INGREDIENTS = " + id).Scan(&ingredient.IdIngredient, &ingredient.Name, &ingredient.Allergen, &ingredient.IdCookingSpace)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error": true,
				"message": "ingredient not found",
			})
			return
		}

		c.JSON(200, ingredient)
		return

	}
}

func UpdateIngredient(tokenAPI string) func(c *gin.Context) {
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

		var ingredient Ingredient

		err = c.BindJSON(&ingredient)
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
		
		if ingredient.Name != "" {
			if !utils.IsSafeString(ingredient.Name) {
				c.JSON(400, gin.H{
					"error": true,
					"message": "name can't contain sql injection",
				})
				return
			}
			if len(ingredient.Name) < 0 || len(ingredient.Name) > 100 {
				c.JSON(400, gin.H{
					"error": true,
					"message": "name must be between 0 and 100 characters",
				})
				return
			}
			setClause = append(setClause, "Name = '"+ingredient.Name+"'")
		}

		if ingredient.Allergen != "" {
			if !utils.IsSafeString(ingredient.Allergen) {
				c.JSON(400, gin.H{
					"error": true,
					"message": "allergen can't contain sql injection",
				})
				return
			}
			if len(ingredient.Allergen) < 0 || len(ingredient.Allergen) > 50 {
				c.JSON(400, gin.H{
					"error": true,
					"message": "allergen must be between 0 and 50 characters",
				})
				return
			}
			setClause = append(setClause, "Allergen = '"+ingredient.Allergen+"'")
		}

		if ingredient.IdCookingSpace > 0 {
			var idcookingspace int
			err = db.QueryRow("SELECT Id_COOKING_SPACES FROM COOKING_SPACES WHERE Id_COOKING_SPACES = '" + strconv.Itoa(ingredient.IdCookingSpace) + "'").Scan(&idcookingspace)
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "cookingspace not found",
				})
				return
			}
			setClause = append(setClause, "Id_COOKING_SPACES = '"+strconv.Itoa(ingredient.IdCookingSpace)+"'")
		}
		
		if len(setClause) == 0 {
			c.JSON(400, gin.H{
				"error": true,
				"message": "nothing to update",
			})
			return
		}

		_, err = db.Exec("UPDATE INGREDIENTS SET " + strings.Join(setClause, ", ") + " WHERE Id_INGREDIENTS = " + id)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "ingredient not found",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
			"message": "ingredient updated",
		})
		return
	}
}

func PostIngredient(tokenAPI string) func(c *gin.Context) {
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

		var ingredient Ingredient
		c.BindJSON(&ingredient)

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		if ingredient.Name == "" || !utils.IsSafeString(ingredient.Name) || len(ingredient.Name) > 100 || len(ingredient.Name) < 0 {
			c.JSON(400, gin.H{
				"error": true,
				"message": "name can't be empty, above 100 char, below 0 char or contain sql injection",
			})
			return
		}
		if ingredient.Allergen == "" || !utils.IsSafeString(ingredient.Allergen) || len(ingredient.Allergen) > 50 || len(ingredient.Allergen) < 0 {
			c.JSON(400, gin.H{
				"error": true,
				"message": "allergen can't be empty, above 50 char, below 0 char or contain sql injection",
			})
			return
		}

		var idcookingspace int

		err = db.QueryRow("SELECT Id_COOKING_SPACES FROM COOKING_SPACES WHERE Id_COOKING_SPACES = '" + strconv.Itoa(ingredient.IdCookingSpace) + "'").Scan(&idcookingspace)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cookingspace not found",
			})
			return
		}

		rows, err := db.Exec("INSERT INTO INGREDIENTS (Name, Allergen, Id_COOKING_SPACES) VALUES (?, ?, ?)", ingredient.Name, ingredient.Allergen, ingredient.IdCookingSpace)
		fmt.Println(err)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot insert ingredient",
			})
			return
		}

		id, err := rows.LastInsertId()

		c.JSON(200, gin.H{
			"error": false,
			"id": id,
			"message": "ingredient inserted",
		})
	}
}

func DeleteIngredient(tokenAPI string) func(c *gin.Context) {
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

		_, err = db.Exec("DELETE FROM INGREDIENTS WHERE Id_INGREDIENTS = " + id)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "ingredient not found",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
			"message": "ingredient deleted",
		})
	}
}

func GetIngredientsByCookingSpaceID(tokenAPI string) func(c *gin.Context) {
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

		rows, err := db.Query("SELECT * FROM INGREDIENTS WHERE Id_COOKING_SPACES = " + id)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "ingredient not found",
			})
			return
		}
		defer rows.Close()

		var ingredients []Ingredient
		for rows.Next() {
			var ingredient Ingredient
			err := rows.Scan(&ingredient.IdIngredient, &ingredient.Name, &ingredient.Allergen, &ingredient.IdCookingSpace)
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "ingredient not found",
				})
				return
			}
			ingredients = append(ingredients, ingredient)
		}

		c.JSON(200, ingredients)
	}
}

func GetCookingSpaceByIngredientId(tokenAPI string) func(c *gin.Context) {
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

		rows, err := db.Query("SELECT COOKING_SPACES.Id_COOKING_SPACES, COOKING_SPACES.name, COOKING_SPACES.size, COOKING_SPACES.isavailable, COOKING_SPACES.priceperhour, COOKING_SPACES.Id_PREMISES, COOKING_SPACES.Picture FROM COOKING_SPACES JOIN INGREDIENTS ON INGREDIENTS.Id_COOKING_SPACES = COOKING_SPACES.Id_COOKING_SPACES WHERE Id_INGREDIENTS = " + id)
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