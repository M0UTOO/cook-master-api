package bills

import (
	"cook-master-api/token"
	"cook-master-api/utils"
	"database/sql"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Bill struct {
	IdBill      int     `json:"idbill"`
	Name 	  string  `json:"name"`
	Type 	  string  `json:"type"`
	CreatedAt string  `json:"createdat"`
	IdUser 	  int     `json:"iduser"`
}

func GetBills(tokenAPI string) func(c *gin.Context) {
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

		rows, err := db.Query("SELECT * FROM BILLS")
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't query database",
			})
			return
		}

		var bills []Bill
		for rows.Next() {
			var bill Bill
			err = rows.Scan(&bill.IdBill, &bill.Name, &bill.Type, &bill.CreatedAt, &bill.IdUser)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "can't scan database",
				})
				return
			}
			bills = append(bills, bill)
		}

		c.JSON(200, bills)
	}
}

func GetBillsByUserID(tokenAPI string) func(c *gin.Context) {
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

		iduser := c.Param("iduser")

		rows, err := db.Query("SELECT * FROM BILLS WHERE ID_USERS = ?", iduser)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't query database",
			})
			return
		}

		var bills []Bill
		for rows.Next() {
			var bill Bill
			err = rows.Scan(&bill.IdBill, &bill.Name, &bill.Type, &bill.CreatedAt, &bill.IdUser)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "can't scan database",
				})
				return
			}
			bills = append(bills, bill)
		}

		c.JSON(200, bills)
	}
}

func GetBillByID(tokenAPI string) func(c *gin.Context) {
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

		rows, err := db.Query("SELECT * FROM BILLS WHERE ID_BILLS = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't query database",
			})
			return
		}

		var bill Bill
		for rows.Next() {
			err = rows.Scan(&bill.IdBill, &bill.Name, &bill.Type, &bill.CreatedAt, &bill.IdUser)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "can't scan database",
				})
				return
			}
		}

		c.JSON(200, bill)
	}
}

func PostBill(tokenAPI string) func(c *gin.Context) {
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

		var bill Bill
		err = c.BindJSON(&bill)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "can't bind json",
			})
			return
		}

		if bill.Name == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "missing name",
			})
			return
		}

		if !utils.IsSafeString(bill.Name) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "name is not safe",
			})
			return
		}

		if len(bill.Name) > 255 || len(bill.Name) < 0 {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "name is too long or too short",
			})
			return
		}

		if bill.Type == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "missing type",
			})
			return
		}

		if !utils.IsSafeString(bill.Type) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "type is not safe",
			})
			return
		}

		if len(bill.Type) > 50 || len(bill.Type) < 0 {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "type is too long or too short",
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

		err = db.QueryRow("SELECT Id_USERS FROM USERS WHERE Id_USERS = '" + strconv.Itoa(bill.IdUser) + "'").Scan(&id)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "client doesn't exist",
			})
			return
		}

		var idbill int

		err = db.QueryRow("SELECT Id_BILLS FROM BILLS WHERE name = ?", bill.Name).Scan(&idbill)
		if err == nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "bill already exist",
			})
			return
		}

		_, err = db.Exec("INSERT INTO BILLS (name, type, createdAt, Id_USERS) VALUES (?, ?, DEFAULT, ?)", bill.Name, bill.Type, bill.CreatedAt, bill.IdUser)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't insert into database",
			})
			return
		}

		c.JSON(200, gin.H{
			"error":   false,
			"message": "bill added",
		})
	}
}

func UpdateBill(tokenAPI string) func(c *gin.Context) {
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

		var bill Bill
		err = c.BindJSON(&bill)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "can't bind json",
			})
			return
		}

		if bill.Name == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "missing name",
			})
			return
		}

		if !utils.IsSafeString(bill.Name) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "name is not safe",
			})
			return
		}

		if len(bill.Name) > 255 || len(bill.Name) < 0 {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "name is too long or too short",
			})
			return
		}

		if bill.Type == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "missing type",
			})
			return
		}

		if !utils.IsSafeString(bill.Type) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "type is not safe",
			})
			return
		}

		if len(bill.Type) > 50 || len(bill.Type) < 0 {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "type is too long or too short",
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

		err = db.QueryRow("SELECT Id_USERS FROM USERS WHERE Id_USERS = '" + strconv.Itoa(bill.IdUser) + "'").Scan(&id)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "client doesn't exist",
			})
			return
		}

		var idbill int

		err = db.QueryRow("SELECT Id_BILLS FROM BILLS WHERE name = ?", bill.Name).Scan(&idbill)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "bill doesn't exist",
			})
			return
		}

		_, err = db.Exec("UPDATE BILLS SET name = ?, type = ?, Id_USERS = ?, WHERE Id_BILLS = ?", bill.Name, bill.Type, bill.IdUser, idbill)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't update database",
			})
			return
		}

		c.JSON(200, gin.H{
			"error":   false,
			"message": "bill updated",
		})
	}
}

func DeleteBill(tokenAPI string) func(c *gin.Context) {
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
				"message": "missing id",
			})
			return
		}

		if !utils.IsSafeString(id) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id is not safe",
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

		var idbill int

		err = db.QueryRow("SELECT Id_BILLS FROM BILLS WHERE Id_BIILS = ?", id).Scan(&idbill)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "bill doesn't exist",
			})
			return
		}

		_, err = db.Exec("DELETE FROM BILLS WHERE Id_BILLS = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't delete bill",
			})
			return
		}

		c.JSON(200, gin.H{
			"error":   false,
			"message": "bill deleted",
		})
	}
}