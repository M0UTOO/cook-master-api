package premises

import (
	"cook-master-api/token"
	"cook-master-api/utils"
	"database/sql"
	"fmt"
	"strings"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Premise struct {
	IdPremise int `json:"idPremise"`
	Name string `json:"name"`
	StreetNumber string `json:"streetNumber"`
	StreetName string `json:"streetName"`
	City string `json:"city"`
	Country string `json:"country"`
}

func GetPremises(tokenAPI string) func(c *gin.Context) {
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

		rows, err := db.Query("SELECT * FROM PREMISES")
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "can't query database",
			})
			return
		}
		defer rows.Close()

		var premises []Premise

		for rows.Next() {
			var premise Premise
			err := rows.Scan(&premise.IdPremise, &premise.Name, &premise.StreetNumber, &premise.StreetName, &premise.City, &premise.Country)
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "can't scan row",
				})
				return
			}
			premises = append(premises, premise)
		}

		c.JSON(200, premises)
	}
}

func GetPremiseByID(tokenAPI string) func(c *gin.Context) {
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

		var premise Premise

		err = db.QueryRow("SELECT * FROM PREMISES WHERE Id_PREMISES = " + id).Scan(&premise.IdPremise, &premise.Name, &premise.StreetNumber, &premise.StreetName, &premise.City, &premise.Country)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "premise not found",
			})
			return
		}

		c.JSON(200, premise)
		return

	}
}

func PostPremise(tokenAPI string) func(c *gin.Context) {
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

		var premise Premise
		err = c.BindJSON(&premise)
		if err != nil {
			c.JSON(400, gin.H{
				"error": true,
				"message": "can't bind json",
			})
			return
		}

		if premise.Name == "" || !utils.IsSafeString(premise.Name) || len(premise.Name) > 100 || len(premise.Name) < 0 {
			c.JSON(400, gin.H{
				"error": true,
				"message": "name can't be empty, contain sql injection or have wrong length",
			})
			return
		}

		if premise.StreetName == "" || !utils.IsSafeString(premise.StreetName) || len(premise.StreetName) > 100 || len(premise.StreetName) < 0 {
			c.JSON(400, gin.H{
				"error": true,
				"message": "street name can't be empty, contain sql injection or have wrong length",
			})
			return
		}

		if premise.StreetNumber == "" || !utils.IsSafeString(premise.StreetNumber) || len(premise.StreetNumber) > 10 || len(premise.StreetNumber) < 0 {
			c.JSON(400, gin.H{
				"error": true,
				"message": "street name can't be empty, contain sql injection or have wrong length",
			})
			return
		}

		if premise.City == "" || !utils.IsSafeString(premise.City) || len(premise.City) > 100 || len(premise.City) < 0 {
			c.JSON(400, gin.H{
				"error": true,
				"message": "city can't be empty, contain sql injection or have wrong length",
			})
			return
		}

		if premise.Country == "" || !utils.IsSafeString(premise.Country) || len(premise.Country) > 100 || len(premise.Country) < 0 {
			c.JSON(400, gin.H{
				"error": true,
				"message": "country can't be empty or contain sql injection",
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

		var id int

		err = db.QueryRow("SELECT Id_PREMISES FROM PREMISES WHERE name = ?", premise.Name).Scan(&id)
		if err == nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "premise already exist",
			})
			return
		}

		rows, err := db.Exec("INSERT INTO PREMISES (name, streetnumber, streetname, city, country) VALUES (?, ?, ?, ?, ?)", premise.Name, premise.StreetNumber, premise.StreetName, premise.City, premise.Country)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "can't insert into database",
			})
			return
		}

		idpremise, err := rows.LastInsertId()

		c.JSON(200, gin.H{
			"error": false,
			"id": idpremise,
			"message": "premise added",
		})
	}
}

func DeletePremise(tokenAPI string) func(c *gin.Context) {
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
			c.JSON(500, gin.H{
				"error": true,
				"message": "missing id",
			})
			return
		}

		if !utils.IsSafeString(id) {
			c.JSON(500, gin.H{
				"error": true,
				"message": "unsafe string",
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

		_, err = db.Exec("DELETE FROM PREMISES WHERE Id_PREMISES = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "can't delete from database",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
			"message": "premise deleted",
		})
	}
}

func UpdatePremise(tokenAPI string) func(c *gin.Context) {
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

		id := c.Param("id")

		if id == "" {
			c.JSON(500, gin.H{
				"error": true,
				"message": "missing id",
			})
			return
		}

		if !utils.IsSafeString(id) {
			c.JSON(500, gin.H{
				"error": true,
				"message": "unsafe string",
			})
			return
		}

		var req Premise

		err = c.BindJSON(&req)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "can't decode json request",
			})
			return
		}

		var setClause []string

		if req.Name != "" {
			if !utils.IsSafeString(req.Name) {
				c.JSON(500, gin.H{
					"error": true,
					"message": "unsafe string",
				})
				return
			}
			if len(req.Name) < 0 || len(req.Name) > 100 {
				c.JSON(500, gin.H{
					"error": true,
					"message": "name must be between 0 and 100 characters",
				})
				return
			}
			setClause = append(setClause, "name = '"+req.Name+"'")
		}

		if req.StreetNumber != "" {
			if !utils.IsSafeString(req.StreetNumber) {
				c.JSON(500, gin.H{
					"error": true,
					"message": "unsafe string",
				})
				return
			}
			if len(req.StreetNumber) < 0 || len(req.StreetNumber) > 100 {
				c.JSON(500, gin.H{
					"error": true,
					"message": "streetnumber must be between 0 and 100 characters",
				})
				return
			}
			setClause = append(setClause, "streetnumber = '"+req.StreetNumber+"'")
		}

		if req.StreetName != "" {
			if !utils.IsSafeString(req.StreetName) {
				c.JSON(500, gin.H{
					"error": true,
					"message": "unsafe string",
				})
				return
			}
			if len(req.StreetName) < 0 || len(req.StreetName) > 100 {
				c.JSON(500, gin.H{
					"error": true,
					"message": "street name must be between 0 and 100 characters",
				})
				return
			}
			setClause = append(setClause, "streetname = '"+req.StreetName+"'")
		}

		if req.City != "" {
			if !utils.IsSafeString(req.City) {
				c.JSON(500, gin.H{
					"error": true,
					"message": "unsafe string",
				})
				return
			}
			if len(req.City) < 0 || len(req.City) > 100 {
				c.JSON(500, gin.H{
					"error": true,
					"message": "city must be between 0 and 100 characters",
				})
				return
			}
			setClause = append(setClause, "city = '"+req.City+"'")
		}

		if req.Country != "" {
			if !utils.IsSafeString(req.Country) {
				c.JSON(500, gin.H{
					"error": true,
					"message": "unsafe string",
				})
				return
			}
			if len(req.Country) < 0 || len(req.Country) > 100 {
				c.JSON(500, gin.H{
					"error": true,
					"message": "country must be between 0 and 100 characters",
				})
				return
			}
			setClause = append(setClause, "country = '"+req.Country+"'")
		}

		if len(setClause) == 0 {
			c.JSON(500, gin.H{
				"error": true,
				"message": "missing field",
			})
			return
		}

		_, err = db.Exec("UPDATE PREMISES SET " + strings.Join(setClause, ", ") + " WHERE Id_PREMISES = ?", id)
		fmt.Println(err)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "can't update database",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
			"message": "premise updated",
		})
	}
}

func GetBooksByPremises(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header["Token"]

		type BooksCount struct {
			Name string `json:"name"`
			Count int    `json:"count"`
		}

		if len(tokenHeader) == 0 {
			c.JSON(400, gin.H{
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

		rows, err := db.Query("SELECT p.name AS premises_name, COUNT(b.Id_BOOKS) AS book_count FROM BOOKS b JOIN COOKING_SPACES cs ON b.Id_COOKING_SPACES = cs.Id_COOKING_SPACES JOIN PREMISES p ON cs.Id_PREMISES = p.Id_PREMISES GROUP BY p.Id_PREMISES, p.name ORDER BY p.Id_PREMISES ASC;")
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "premises not found",
			})
			return
		}

		var bookscounts []BooksCount

		for rows.Next() {
			var bookcount BooksCount
			err = rows.Scan(&bookcount.Name, &bookcount.Count)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error":   true,
					"message": "books not found",
				})
				return
			}
			bookscounts = append(bookscounts, bookcount)
		}

		c.JSON(200, bookscounts)
		return
	}
}

func GetPremiseByCookingSpace(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header["Token"]

		if len(tokenHeader) == 0 {
			c.JSON(400, gin.H{
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

		rows, err := db.Query("SELECT * FROM PREMISES WHERE Id_PREMISES = (SELECT Id_PREMISES FROM COOKING_SPACES WHERE Id_COOKING_SPACES = ?)", id)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "premises not found",
			})
			return
		}

		var premises []Premise

		for rows.Next() {
			var premise Premise
			err = rows.Scan(&premise.IdPremise, &premise.Name, &premise.StreetNumber, &premise.StreetName, &premise.City, &premise.Country)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error":   true,
					"message": "premises not found",
				})
				return
			}
			premises = append(premises, premise)
		}

		c.JSON(200, premises)
		return
	}
}
