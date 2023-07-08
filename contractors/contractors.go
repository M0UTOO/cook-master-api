package contractors

import (
	"cook-master-api/token"
	"cook-master-api/utils"
	"database/sql"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Contractor struct {
	Presentation  string `json:"presentation"`
	ContractStart string `json:"contractstart"`
	ContractEnd   string `json:"contractend"`
	ContractorType          string `json:"contractortype"`
}

type ContractorUser struct {
	Email          string `json:"email"`
	FirstName      string `json:"firstname"`
	LastName       string `json:"lastname"`
	ProfilePicture string `json:"profilepicture"`
	IsCreatedAt    string `json:"iscreatedat"`
	LastSeen       string `json:"lastseen"`
	IsBlocked      string `json:"isblocked"`
	IdContractor   int    `json:"idcontractor"`
	Presentation   string `json:"presentation"`
	ContractStart  string `json:"contractstart"`
	ContractEnd    string `json:"contractend"`
	ContractorType           int `json:"contractortype"`
	IdUsers        int    `json:"idusers"`
	Language 	 int `json:"language"`
}

type ContractorType struct {
	IdContractorType int    `json:"idcontractortype"`
	Name 		   string `json:"name"`
}

func GetContractors(tokenAPI string) func(c *gin.Context) {
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

		rows, err := db.Query("SELECT USERS.email, USERS.firstname, USERS.lastname, USERS.profilepicture, USERS.iscreatedat, USERS.lastseen, USERS.isblocked, USERS.Id_LANGUAGES, CONTRACTORS.Id_CONTRACTORS, CONTRACTORS.presentation, CONTRACTORS.contractStart, CONTRACTORS.contractEnd, CONTRACTORS.Id_CONTRACTOR_TYPES, CONTRACTORS.Id_USERS FROM CONTRACTORS JOIN USERS ON CONTRACTORS.Id_USERS = USERS.Id_USERS ORDER BY USERS.lastname DESC")
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "user not found",
			})
			return
		}

		var contractors []ContractorUser

		for rows.Next() {
			var contractor ContractorUser
			err = rows.Scan(&contractor.Email, &contractor.FirstName, &contractor.LastName, &contractor.ProfilePicture, &contractor.IsCreatedAt, &contractor.LastSeen, &contractor.IsBlocked, &contractor.Language, &contractor.IdContractor, &contractor.Presentation, &contractor.ContractStart, &contractor.ContractEnd, &contractor.ContractorType, &contractor.IdUsers)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error":   true,
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

		var contractor ContractorUser

		err = db.QueryRow("SELECT USERS.email, USERS.firstname, USERS.lastname, USERS.profilepicture, USERS.iscreatedat, USERS.lastseen, USERS.isblocked, USERS.Id_LANGUAGES, CONTRACTORS.Id_CONTRACTORS, CONTRACTORS.presentation, CONTRACTORS.contractStart, CONTRACTORS.contractEnd, CONTRACTORS.Id_CONTRACTOR_TYPES, CONTRACTORS.Id_USERS FROM CONTRACTORS JOIN USERS ON CONTRACTORS.Id_USERS = USERS.Id_USERS WHERE CONTRACTORS.Id_USERS = "+id).Scan(&contractor.Email, &contractor.FirstName, &contractor.LastName, &contractor.ProfilePicture, &contractor.IsCreatedAt, &contractor.LastSeen, &contractor.IsBlocked, &contractor.Language, &contractor.IdContractor, &contractor.Presentation, &contractor.ContractStart, &contractor.ContractEnd, &contractor.ContractorType, &contractor.IdUsers)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
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

		var contractor Contractor

		err = c.BindJSON(&contractor)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "bad json",
			})
			return
		}

		var setClause []string

		if !utils.IsSafeString(contractor.Presentation) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "bad json",
			})
			return
		}

		if contractor.Presentation != "" || !utils.IsSafeString(contractor.Presentation) {
			setClause = append(setClause, "Presentation = '"+contractor.Presentation+"'")
		}
		if contractor.ContractStart != "" || !utils.IsSafeString(contractor.ContractStart) {
			setClause = append(setClause, "ContractStart = '"+contractor.ContractStart+"'")
		}
		if contractor.ContractEnd != "" || !utils.IsSafeString(contractor.ContractEnd) {
			setClause = append(setClause, "ContractEnd = '"+contractor.ContractEnd+"'")
		}
		if contractor.ContractorType != "" || !utils.IsSafeString(contractor.ContractorType) {
			setClause = append(setClause, "ContractorType = '"+contractor.ContractorType+"'")
		}

		if len(setClause) == 0 {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "nothing to update",
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

		var idcontractor int

		err = db.QueryRow("SELECT Id_CONTRACTORS FROM CONTRACTORS WHERE Id_USERS = '" + id + "'").Scan(&idcontractor)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "contractor not found",
			})
			return
		}

		_, err = db.Exec("UPDATE CONTRACTORS SET " + strings.Join(setClause, ", ") + " WHERE Id_USERS = '" + id + "'")
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot update contractor",
			})
			return
		}

		c.JSON(200, gin.H{
			"error":   false,
			"message": "contractor updated",
		})
		return
	}
}

func AddAContractorType(tokenAPI string) func(c *gin.Context) {
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

		var contractorType ContractorType

		err = c.BindJSON(&contractorType)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "bad json",
			})
			return
		}

		if contractorType.Name == "" || !utils.IsSafeString(contractorType.Name) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "name can't be empty",
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

		_, err = db.Exec("INSERT INTO CONTRACTOR_TYPES (Name) VALUES ('" + contractorType.Name + "')")
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot add contractor type",
			})
			return
		}

		c.JSON(200, gin.H{
			"error":   false,
			"message": "contractor type added",
		})
		return
	}
}

func DeleteAContractorType(tokenAPI string) func(c *gin.Context) {
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

		_, err = db.Exec("UPDATE CONTRACTORS SET Id_CONTRACTOR_TYPES = 0 WHERE Id_CONTRACTOR_TYPES = " + id)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "contractor not found",
			})
			return
		}

		_, err = db.Exec("DELETE FROM CONTRACTOR_TYPES WHERE Id_CONTRACTOR_TYPES = " + id)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot delete contractor type",
			})
			return
		}

		c.JSON(200, gin.H{
			"error":   false,
			"message": "contractor type deleted",
		})
		return
	}
}

func GetContractorTypes(tokenAPI string) func(c *gin.Context) {
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

		rows, err := db.Query("SELECT Id_CONTRACTOR_TYPES, Name FROM CONTRACTOR_TYPES")
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot get contractor types",
			})
			return
		}
		defer rows.Close()

		var contractorTypes []ContractorType

		for rows.Next() {
			var contractorType ContractorType
			err = rows.Scan(&contractorType.IdContractorType, &contractorType.Name)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "cannot get contractor types",
				})
				return
			}
			contractorTypes = append(contractorTypes, contractorType)
		}

		c.JSON(200, contractorTypes)
		return
	}
}

func GetContractorByRoles(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		type ContractorRole struct {
			IdContractor int `json:"idcontractor"`
			Firstname string `json:"firstname"`
			Lastname string `json:"lastname"`
			ProfilePicture string `json:"profilepicture"`
			Name string `json:"name"`
		}

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

		search := c.Param("search")
		if search == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "search can't be empty",
			})
			return
		}

		if !utils.IsSafeString(search) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "search can't contain sql injection",
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

		rows, err := db.Query("SELECT CONTRACTORS.Id_CONTRACTORS, USERS.firstname, USERS.lastname, USERS.profilePicture, CONTRACTOR_TYPES.name FROM CONTRACTORS JOIN CONTRACTOR_TYPES ON CONTRACTORS.Id_CONTRACTOR_TYPES = CONTRACTOR_TYPES.Id_CONTRACTOR_TYPES JOIN USERS ON CONTRACTORS.Id_USERS = USERS.Id_USERS WHERE CONTRACTOR_TYPES.name = '" + search + "'")
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot get contractor types",
			})
			return
		}
		defer rows.Close()

		var contractorroles []ContractorRole

		for rows.Next() {
			var contractorrole ContractorRole
			err = rows.Scan(&contractorrole.IdContractor, &contractorrole.Firstname, &contractorrole.Lastname, &contractorrole.ProfilePicture, &contractorrole.Name)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "cannot get contractor role",
				})
				return
			}
			contractorroles = append(contractorroles, contractorrole)
		}

		c.JSON(200, contractorroles)
		return
	}
}

func GetCreatorRoleById(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		type ContractorRole struct {
			Name string `json:"name"`
		}

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

		rows, err := db.Query("SELECT CONTRACTOR_TYPES.name FROM CONTRACTOR_TYPES WHERE Id_CONTRACTOR_TYPES = " + id)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot get contractor types",
			})
			return
		}
		defer rows.Close()

		var contractorroles []ContractorRole

		for rows.Next() {
			var contractorrole ContractorRole
			err = rows.Scan(&contractorrole.Name)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "cannot get contractor role",
				})
				return
			}
			contractorroles = append(contractorroles, contractorrole)
		}

		c.JSON(200, contractorroles)
		return
	}
}