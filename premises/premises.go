package premises

import (
	"cook-master-api/token"
	"cook-master-api/utils"
	"database/sql"
	"strconv"
	"fmt"
	"strings"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Premise struct {
	IdPremise int `json:"idPremise"`
	Name string `json:"name"`
	StreetNumber int `json:"streetNumber"`
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

		if premise.Name == "" || !utils.IsSafeString(premise.Name) {
			c.JSON(400, gin.H{
				"error": true,
				"message": "name can't be empty or contain sql injection",
			})
			return
		}

		if premise.StreetNumber < 0 {
			c.JSON(400, gin.H{
				"error": true,
				"message": "street number can't be negative",
			})
			return
		}

		if premise.StreetName == "" || !utils.IsSafeString(premise.StreetName) {
			c.JSON(400, gin.H{
				"error": true,
				"message": "street name can't be empty or contain sql injection",
			})
			return
		}

		if premise.City == "" || !utils.IsSafeString(premise.City) {
			c.JSON(400, gin.H{
				"error": true,
				"message": "city can't be empty or contain sql injection",
			})
			return
		}

		if premise.Country == "" || !utils.IsSafeString(premise.Country) {
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

		_, err = db.Exec("INSERT INTO PREMISES (name, streetnumber, streetname, city, country) VALUES (?, ?, ?, ?, ?)", premise.Name, premise.StreetNumber, premise.StreetName, premise.City, premise.Country)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "can't insert into database",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
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

		if req.StreetNumber > 0 {
			setClause = append(setClause, "streetnumber = '"+strconv.Itoa(req.StreetNumber)+"'")
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
			if len(req.Country) < 0 || len(req.Country) > 50 {
				c.JSON(500, gin.H{
					"error": true,
					"message": "country must be between 0 and 50 characters",
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

