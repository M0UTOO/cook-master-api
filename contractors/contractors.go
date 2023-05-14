package contractors

import (
	"cook-master-api/token"
	"cook-master-api/utils"
	"database/sql"
	"strings"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Contractor struct {
	Presentation string `json:"presentation"`
}

type ContractorUser struct {
	Id int `json:"id"`
	Email string `json:"email"`
	Password string `json:"password"`
	FirstName string `json:"firstname"`
	LastName string `json:"lastname"`
	ProfilePicture string `json:"profilepicture"`
	IdContractor int `json:"idcontractor"`
	Presentation string `json:"presentation"`
	IdUsers int `json:"idusers"`
}

func GetContractors(tokenAPI string) func(c *gin.Context) {
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

		rows, err := db.Query("SELECT * FROM CONTRACTORS JOIN USERS ON CONTRACTORS.Id_USERS = USERS.Id_USERS ORDER BY USERS.lastname DESC")
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "user not found",
			})
			return
		}

		var contractors []ContractorUser

		for rows.Next() {
			var contractor ContractorUser
			err = rows.Scan(&contractor.IdContractor, &contractor.Presentation, &contractor.IdUsers, &contractor.Id, &contractor.Email, &contractor.Password, &contractor.FirstName, &contractor.LastName, &contractor.ProfilePicture)
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "contractor not found",
				})
				return
			}
			contractors = append(contractors, contractor)
		}

		c.JSON(200, contractors)
		return
	}
}

func GetContractorByID(tokenAPI string) func(c *gin.Context) {
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

		var contractor ContractorUser

		err = db.QueryRow("SELECT * FROM CONTRACTORS JOIN USERS ON CONTRACTORS.Id_USERS = USERS.Id_USERS WHERE CONTRACTORS.Id_USERS = " + id).Scan(&contractor.Id, &contractor.Presentation, &contractor.IdContractor, &contractor.Id, &contractor.Email, &contractor.Password, &contractor.FirstName, &contractor.LastName, &contractor.ProfilePicture)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "contractor not found",
			})
			return
		}

		c.JSON(200, contractor)
		return

	}
}

func UpdateContractor(tokenAPI string) func(c *gin.Context) {
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

		var contractor Contractor

		err = c.BindJSON(&contractor)
		if err != nil {
			c.JSON(400, gin.H{
				"error": true,
				"message": "bad json",
			})
			return
		}

		var setClause []string
		
		if !utils.IsSafeString(contractor.Presentation) { 
			c.JSON(400, gin.H{
				"error": true,
				"message": "bad json",
			})
			return
		}

		if contractor.Presentation != "" || !utils.IsSafeString(contractor.Presentation) {
			setClause = append(setClause, "Presentation = '" + contractor.Presentation + "'")
		}

		if len(setClause) == 0 {
			c.JSON(400, gin.H{
				"error": true,
				"message": "nothing to update",
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

		var idcontractor int

		err = db.QueryRow("SELECT Id_CONTRACTORS FROM CONTRACTORS WHERE Id_CONTRACTORS = '" + id + "'").Scan(&idcontractor)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "contractor not found",
			})
			return
		}

		_, err = db.Exec("UPDATE CONTRACTORS SET " + strings.Join(setClause, ", ") + " WHERE Id_CONTRACTORS = " + id)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "contractor not found",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
			"message": "contractor updated",
		})
		return
	}
}