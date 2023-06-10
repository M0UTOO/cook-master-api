package managers

import (
	"cook-master-api/token"
	"cook-master-api/utils"
	"database/sql"
	"strings"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Manager struct {
	IdManager int `json:"idmanager"`
	IsItemManager bool `json:"isitemmanager"`
	IsClientManager bool `json:"isclientmanager"`
	IsContractorManager bool `json:"iscontractormanager"`
	IsSuperAdmin bool `json:"issuperadmin"`
}

type ManagerUser struct {
	Email string `json:"email"`
	FirstName string `json:"firstname"`
	LastName string `json:"lastname"`
	ProfilePicture string `json:"profilepicture"`
	IsCreatedAt string `json:"iscreatedat"`
	LastSeen string `json:"lastseen"`
	IsBlocked string `json:"isblocked"`
	IdManager int `json:"idmanager"`
	IsItemManager bool `json:"isitemmanager"`
	IsClientManager bool `json:"isclientmanager"`
	IsContractorManager bool `json:"iscontractormanager"`
	IsSuperAdmin bool `json:"issuperadmin"`
	IdUsers int `json:"idusers"`
	Language int `json:"language"`
}

func GetManagers(tokenAPI string) func(c *gin.Context) {
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

		rows, err := db.Query("SELECT USERS.email, USERS.firstname, USERS.lastname, USERS.profilepicture, USERS.iscreatedat, USERS.lastseen, USERS.isblocked, USERS.Id_LANGUAGES, MANAGERS.Id_MANAGERS, MANAGERS.isItemManager, MANAGERS.isClientManager, MANAGERS.isContractorManager, MANAGERS.isSuperAdmin, MANAGERS.Id_USERS FROM MANAGERS JOIN USERS ON MANAGERS.Id_USERS = USERS.Id_USERS ORDER BY USERS.lastname DESC")
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "user not found",
			})
			return
		}

		var managers []ManagerUser

		for rows.Next() {
			var manager ManagerUser
			err = rows.Scan(&manager.Email, &manager.FirstName, &manager.LastName, &manager.ProfilePicture, &manager.IsCreatedAt, &manager.LastSeen, &manager.IsBlocked, &manager.Language, &manager.IdManager, &manager.IsItemManager, &manager.IsClientManager, &manager.IsContractorManager, &manager.IsSuperAdmin, &manager.IdUsers)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error": true,
					"message": "manager not found",
				})
				return
			}
			managers = append(managers, manager)
		}

		c.JSON(200, managers)
		return
	}
}

func GetManagerByID(tokenAPI string) func(c *gin.Context) {
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

		var manager ManagerUser

		err = db.QueryRow("SELECT USERS.email, USERS.firstname, USERS.lastname, USERS.profilepicture, USERS.iscreatedat, USERS.lastseen, USERS.isblocked, USERS.Id_LANGUAGES, MANAGERS.Id_MANAGERS, MANAGERS.isItemManager, MANAGERS.isClientManager, MANAGERS.isContractorManager, MANAGERS.isSuperAdmin, MANAGERS.Id_USERS FROM MANAGERS JOIN USERS ON MANAGERS.Id_USERS = USERS.Id_USERS WHERE MANAGERS.Id_USERS = " + id).Scan(&manager.Email, &manager.FirstName, &manager.LastName, &manager.ProfilePicture, &manager.IsCreatedAt, &manager.LastSeen, &manager.IsBlocked, &manager.Language, &manager.IdManager, &manager.IsItemManager, &manager.IsClientManager, &manager.IsContractorManager, &manager.IsSuperAdmin, &manager.IdUsers)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "manager not found",
			})
			return
		}

		c.JSON(200, manager)
		return

	}
}

func UpdateManager(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		type Manager struct {
			IsItemManager int `json:"isitemmanager"`
			IsClientManager int `json:"isclientmanager"`
			IsContractorManager int `json:"iscontractormanager"`
			IsSuperAdmin int `json:"issuperadmin"`
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

		var manager Manager

		manager.IsItemManager = -1
		manager.IsClientManager = -1
		manager.IsContractorManager = -1
		manager.IsSuperAdmin = -1

		err = c.BindJSON(&manager)
		if err != nil {
			c.JSON(400, gin.H{
				"error": true,
				"message": "bad json",
			})
			return
		}

		var setClause []string

		if manager.IsItemManager == 0 {
			setClause = append(setClause, "isitemmanager = false")
		} else if manager.IsItemManager == 1 {
			setClause = append(setClause, "isitemmanager = true")
		}

		if manager.IsClientManager == 0 {
			setClause = append(setClause, "isclientmanager = false")
		} else if manager.IsClientManager == 1 {
			setClause = append(setClause, "isclientmanager = true")
		}

		if manager.IsContractorManager == 0 {
			setClause = append(setClause, "iscontractormanager = false")
		} else if manager.IsContractorManager == 1 {
			setClause = append(setClause, "iscontractormanager = true")
		}

		if manager.IsSuperAdmin == 0 {
			setClause = append(setClause, "issuperadmin = false")
		} else if manager.IsSuperAdmin == 1 {
			setClause = append(setClause, "issuperadmin = true")
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

		var idmanager int

		err = db.QueryRow("SELECT Id_MANAGERS FROM MANAGERS WHERE Id_USERS = '" + id + "'").Scan(&idmanager)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "manager not found",
			})
			return
		}

		_, err = db.Exec("UPDATE MANAGERS SET " + strings.Join(setClause, ", ") + " WHERE Id_USERS = " + id)
		fmt.Println(err)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot update manager",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
			"message": "manager updated",
		})
		return
	}
}