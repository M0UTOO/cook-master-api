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
	Id int `json:"id"`
	Email string `json:"email"`
	Password string `json:"password"`
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

		rows, err := db.Query("SELECT * FROM MANAGERS JOIN USERS ON MANAGERS.Id_USERS = USERS.Id_USERS ORDER BY USERS.lastname DESC")
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
			err = rows.Scan(&manager.IdManager, &manager.IsItemManager, &manager.IsClientManager, &manager.IsContractorManager, &manager.IsSuperAdmin, &manager.IdUsers, &manager.Id, &manager.Email, &manager.Password, &manager.FirstName, &manager.LastName, &manager.ProfilePicture, &manager.IsCreatedAt, &manager.LastSeen, &manager.IsBlocked)
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

		err = db.QueryRow("SELECT * FROM MANAGERS JOIN USERS ON MANAGERS.Id_USERS = USERS.Id_USERS WHERE MANAGERS.Id_USERS = " + id).Scan(&manager.IdManager, &manager.IsItemManager, &manager.IsClientManager, &manager.IsContractorManager, &manager.IsSuperAdmin, &manager.IdUsers, &manager.Id, &manager.Email, &manager.Password, &manager.FirstName, &manager.LastName, &manager.ProfilePicture, &manager.IsCreatedAt, &manager.LastSeen, &manager.IsBlocked)
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

		err = c.BindJSON(&manager)
		if err != nil {
			c.JSON(400, gin.H{
				"error": true,
				"message": "bad json",
			})
			return
		}

		var setClause []string

		if manager.IsItemManager == true {
			setClause = append(setClause, "isItemManager = 1")
		}
		if manager.IsItemManager == false {
			setClause = append(setClause, "isItemManager = 0")
		}
		if manager.IsClientManager == true {
			setClause = append(setClause, "isClientManager = 1")
		}
		if manager.IsClientManager == false {
			setClause = append(setClause, "isClientManager = 0")
		}
		if manager.IsContractorManager == true {
			setClause = append(setClause, "isContractorManager = 1")
		}
		if manager.IsContractorManager == false {
			setClause = append(setClause, "isContractorManager = 0")
		}
		if manager.IsSuperAdmin == true {
			setClause = append(setClause, "isSuperAdmin = 1")
		}
		if manager.IsSuperAdmin == false {
			setClause = append(setClause, "isSuperAdmin = 0")
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

		err = db.QueryRow("SELECT Id_MANAGERS FROM MANAGERS WHERE Id_MANAGERS = '" + id + "'").Scan(&idmanager)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "manager not found",
			})
			return
		}

		_, err = db.Exec("UPDATE MANAGERS SET " + strings.Join(setClause, ", ") + " WHERE Id_MANAGERS = " + id)
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