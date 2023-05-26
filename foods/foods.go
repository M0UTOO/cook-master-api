package foods

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

type Food struct {
	IdFood      int     `json:"idfood"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func GetFoods(tokenAPI string) func(c *gin.Context) {
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

		rows, err := db.Query("SELECT * FROM FOODS")
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't query database",
			})
			return
		}
		defer rows.Close()

		var foods []Food

		for rows.Next() {
			var food Food
			err := rows.Scan(&food.IdFood, &food.Name, &food.Description, &food.Price)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "can't scan row",
				})
				return
			}
			foods = append(foods, food)
		}

		c.JSON(200, foods)
	}
}

func GetFoodByID(tokenAPI string) func(c *gin.Context) {
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

		var food Food

		err = db.QueryRow("SELECT * FROM FOODS WHERE Id_FOODS = "+id).Scan(&food.IdFood, &food.Name, &food.Description, &food.Price)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "food not found",
			})
			return
		}

		c.JSON(200, food)
		return

	}
}

func PostFood(tokenAPI string) func(c *gin.Context) {
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

		var food Food
		err = c.BindJSON(&food)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "can't bind json",
			})
			return
		}

		if food.Name == "" || !utils.IsSafeString(food.Name) || len(food.Name) < 0 || len(food.Name) > 100 {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "name can't be empty, contain sql injection or wrong lenght",
			})
			return
		}

		if food.Price < 0 {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "price can't be negative",
			})
			return
		}

		if food.Description == "" || !utils.IsSafeString(food.Description) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "description can't be empty or contain sql injection",
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

		err = db.QueryRow("SELECT Id_FOODS FROM FOODS WHERE name = ?", food.Name).Scan(&id)
		if err == nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "food already exist",
			})
			return
		}

		_, err = db.Exec("INSERT INTO FOODS (name, description, price) VALUES (?, ?, ?)", food.Name, food.Description, food.Price)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't insert into database",
			})
			return
		}

		c.JSON(200, gin.H{
			"error":   false,
			"message": "food added",
		})
	}
}

func DeleteFood(tokenAPI string) func(c *gin.Context) {
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

		_, err = db.Exec("DELETE FROM FOODS WHERE Id_FOODS = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't delete from database",
			})
			return
		}

		c.JSON(200, gin.H{
			"error":   false,
			"message": "food deleted",
		})
	}
}

func UpdateFood(tokenAPI string) func(c *gin.Context) {
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

		var req Food

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

		if len(setClause) == 0 {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "missing field",
			})
			return
		}

		_, err = db.Exec("UPDATE FOODS SET "+strings.Join(setClause, ", ")+" WHERE Id_FOODS = ?", id)
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
			"message": "food updated",
		})
	}
}
