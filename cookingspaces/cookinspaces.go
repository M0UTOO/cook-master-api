package cookingspaces

import (
	"fmt"
	"database/sql"
	"strconv"
	"strings"
	"cook-master-api/token"
	"cook-master-api/utils"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type CookingSpace struct {
	IdCookingSpace int `json:"idCookingSpace"`
	Name string `json:"name"`
	Size int `json:"size"`
	IsAvailable bool `json:"isAvailable"`
	PricePerHour float64 `json:"pricePerHour"`
	IdPremise int `json:"idPremise"`
}

func GetCookingSpaces(tokenAPI string) func(c *gin.Context) {
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

		rows, err := db.Query("SELECT * FROM COOKING_SPACES")
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot get cookingspaces",
			})
			return
		}
		defer rows.Close()

		var cookingspaces []CookingSpace

		for rows.Next() {
			var cookingspace CookingSpace
			err = rows.Scan(&cookingspace.IdCookingSpace, &cookingspace.Name, &cookingspace.Size, &cookingspace.IsAvailable, &cookingspace.PricePerHour, &cookingspace.IdPremise)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error": true,
					"message": "err on scan rows",
				})
				return
			}
			cookingspaces = append(cookingspaces, cookingspace)
		}

		c.JSON(200, cookingspaces)
	}
}

func GetCookingSpaceByID(tokenAPI string) func(c *gin.Context) {
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

		var cookingspace CookingSpace

		err = db.QueryRow("SELECT * FROM COOKING_SPACES WHERE Id_COOKING_SPACES = ?", id).Scan(&cookingspace.IdCookingSpace, &cookingspace.Name, &cookingspace.Size, &cookingspace.IsAvailable, &cookingspace.PricePerHour, &cookingspace.IdPremise)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot get cookingspace",
			})
			return
		}

		c.JSON(200, cookingspace)
	}
}

func GetCookingSpacesByPremiseID(tokenAPI string) func(c *gin.Context) {
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

		rows, err := db.Query("SELECT * FROM COOKING_SPACES WHERE Id_PREMISES = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot get cookingspaces",
			})
			return
		}
		defer rows.Close()

		var cookingspaces []CookingSpace

		for rows.Next() {
			var cookingspace CookingSpace
			err = rows.Scan(&cookingspace.IdCookingSpace, &cookingspace.Name, &cookingspace.Size, &cookingspace.IsAvailable, &cookingspace.PricePerHour, &cookingspace.IdPremise)
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "cannot get cookingspaces",
				})
				return
			}
			cookingspaces = append(cookingspaces, cookingspace)
		}

		c.JSON(200, cookingspaces)
	}
}

func PostCookingSpace(tokenAPI string) func(c *gin.Context) {
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

		var cookingspace CookingSpace
		c.BindJSON(&cookingspace)

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		if cookingspace.Name == "" || !utils.IsSafeString(cookingspace.Name) {
			c.JSON(400, gin.H{
				"error": true,
				"message": "name can't be empty or contain sql injection",
			})
			return
		}

		if cookingspace.Size <= 0 {
			c.JSON(400, gin.H{
				"error": true,
				"message": "size can't be negative or zero",
			})
			return
		}

		if cookingspace.PricePerHour <= 0 {
			c.JSON(400, gin.H{
				"error": true,
				"message": "price per hour can't be negative or zero",
			})
			return
		}

		var idPremise int

		err = db.QueryRow("SELECT Id_PREMISES FROM PREMISES WHERE Id_PREMISES = 1").Scan(&idPremise)
		if err != nil {
			_, err := db.Exec("INSERT INTO PREMISES (name, streetNumber, streetName, city, country) VALUES (?, ?, ?, ?, ?)", "default", 0 , "default", "default", "default")
			fmt.Println(err)
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "cannot insert premise",
				})
				return
			}
		}

		_, err = db.Exec("INSERT INTO COOKING_SPACES (Name, Size, PricePerHour, Id_PREMISES) VALUES (?, ?, ?, ?)", cookingspace.Name, cookingspace.Size, cookingspace.PricePerHour, 1)
		fmt.Println(err)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot insert cookingspace",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
			"message": "cookingspace inserted",
		})
	}
}

func AddCookingSpaceToAPremise(tokenAPI string) func(c *gin.Context) {
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

		type Premise struct {
			Name string `json:"name"`
		}

		var premise Premise
		c.BindJSON(&premise)

		if !utils.IsSafeString(premise.Name) {
			c.JSON(400, gin.H{
				"error": true,
				"message": "name can't contain sql injection",
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

		var idPremise int

		err = db.QueryRow("SELECT Id_PREMISES FROM PREMISES WHERE name = '" + premise.Name + "'").Scan(&idPremise)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot find Premise",
			})
			return
		}

		var idCookingSpace int

		err = db.QueryRow("SELECT Id_COOKING_SPACES FROM COOKING_SPACES WHERE Id_PREMISES = '" + strconv.Itoa(idPremise) + "' AND Id_COOKING_SPACES = '" + id + "'").Scan(&idCookingSpace)
		if err != nil {
			_, err := db.Exec("UPDATE COOKING_SPACES SET Id_PREMISES = ? WHERE Id_COOKING_SPACES = ?", idPremise, id)
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "cannot update cookingspace",
				})
				return
			}

			c.JSON(200, gin.H{
				"error": false,
				"message": "cookingspace added to Premise",
			})
			return
		}

		c.JSON(500, gin.H{
			"error": false,
			"message": "this cookingspace is already in this premise",
		})
	}
}

func DeleteCookingSpaceFromAPremise(tokenAPI string) func(c *gin.Context) {
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

		_, err = db.Exec("UPDATE COOKING_SPACES SET Id_PREMISES = 1 WHERE Id_COOKING_SPACES = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot update cookingspace",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
			"message": "cookingspace deleted from premise",
		})
	}
}

func UpdateCookingSpace(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		type CookingSpaceReq struct {
			Name string `json:"name"`
			Size int `json:"size"`
			IsAvailable int `json:"isAvailable"`
			PricePerHour float64 `json:"pricePerHour"`
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

		var cookingspace CookingSpaceReq

		cookingspace.IsAvailable = -1

		err = c.BindJSON(&cookingspace)	
		if err != nil {
			c.JSON(400, gin.H{
				"error": true,
				"message": "cannot bind json",
			})
			return
		}

		var setClause []string

		if cookingspace.Name != "" {
			if !utils.IsSafeString(cookingspace.Name) {
				c.JSON(400, gin.H{
					"error": true,
					"message": "name can't contain sql injection",
				})
				return
			}
			if len(cookingspace.Name) < 0 || len(cookingspace.Name) > 50 {
				c.JSON(400, gin.H{
					"error": true,
					"message": "wrong name length",
				})
				return
			}
			setClause = append(setClause, "name = '"+cookingspace.Name+"'")
		}

		if cookingspace.Size > 0 {
			setClause = append(setClause, "size = '"+strconv.Itoa(cookingspace.Size)+"'")
		}
		
		if cookingspace.IsAvailable == 0 {
			setClause = append(setClause, "isAvailable = false")
		} else if cookingspace.IsAvailable == 1 {
			setClause = append(setClause, "isAvailable = true")
		}

		if cookingspace.PricePerHour > 0 {
			setClause = append(setClause, "pricePerHour = '"+strconv.FormatFloat(cookingspace.PricePerHour, 'f', 2, 64)+"'")
		}

		if len(setClause) == 0 {
			c.JSON(400, gin.H{
				"error": true,
				"message": "no field to update",
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

		var idcookingspace string
		
		err = db.QueryRow("SELECT Id_COOKING_SPACES FROM COOKING_SPACES WHERE Id_COOKING_SPACES = '" + id + "'").Scan(&idcookingspace)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cookingspace not found",
			})
			return
		}

		_, err = db.Exec("UPDATE COOKING_SPACES SET "+strings.Join(setClause, ", ")+" WHERE Id_COOKING_SPACES = ?", id)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot update cookingspace",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
			"message": "cookingspace updated",
		})
	}
}