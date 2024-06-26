package users

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"cook-master-api/token"
	"cook-master-api/utils"
	"regexp"
	"fmt"
	"strings"
	"strconv"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id int `json:"id"`
	Email string `json:"email"`
	Password string `json:"password"`
	FirstName string `json:"firstname"`
	LastName string `json:"lastname"`
	ProfilePicture string `json:"profilepicture"`
	IsCreatedAt string `json:"iscreatedat"`
	LastSeen string `json:"lastseen"`
	IsBlocked string `json:"isblocked"`
	Language int `json:"language"`
}

type UserNoPassword struct {
	Id int `json:"id"`
	Email string `json:"email"`
	FirstName string `json:"firstname"`
	LastName string `json:"lastname"`
	ProfilePicture string `json:"profilepicture"`
	IsCreatedAt string `json:"iscreatedat"`
	LastSeen string `json:"lastseen"`
	IsBlocked string `json:"isblocked"`
	Language int `json:"language"`
}

func GetUserByID(tokenAPI string) func(c *gin.Context) {
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

		var user UserNoPassword

		err = db.QueryRow("SELECT USERS.Id_USERS, USERS.email, USERS.firstname, USERS.lastname, USERS.profilepicture, USERS.iscreatedat, USERS.lastseen, USERS.isblocked, USERS.Id_LANGUAGES FROM USERS WHERE Id_USERS=" + id).Scan(&user.Id, &user.Email, &user.FirstName, &user.LastName, &user.ProfilePicture, &user.IsCreatedAt, &user.LastSeen, &user.IsBlocked, &user.Language)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error": true,
				"message": "user not found",
			})
			return
		}

		c.JSON(200, user)
		return
	}
}

func PostUser(tokenAPI string) func(c *gin.Context) {
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

		typeOfUser := c.Request.Header["Type"]
		if typeOfUser == nil{
			c.JSON(498, gin.H{
				"error": true,
				"message": "missing type",
			})
		}

		if typeOfUser[0] != "Client" && typeOfUser[0] != "Contractor" && typeOfUser[0] != "Manager" {
			c.JSON(498, gin.H{
				"error": true,
				"message": "wrong type",
			})
			return
		}

		type UserRequest struct {
			Email string `json:"email"`
			Password string `json:"password"`
			FirstName string `json:"firstname"`
			LastName string `json:"lastname"`
			Country string `json:"country"`
			Subscription int `json:"subscription"`
			Language int `json:"language"`
			Presentation string `json:"presentation"`
			ContractStart string `json:"contractstart"`
			ContractEnd string `json:"contractend"`
			IdContractorType int `json:"idcontractortype"`
			IsItemManager bool `json:"isitemmanager"`
			IsClientManager bool `json:"isclientmanager"`
			IsContractorManager bool `json:"iscontractormanager"`
			IsSuperAdmin bool `json:"issuperadmin"`
		}

		var req UserRequest
		
		err = c.BindJSON(&req)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "can't decode json request",
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

		if !utils.IsSafeString(req.Email) || !utils.IsSafeString(req.Password) || !utils.IsSafeString(req.FirstName) || !utils.IsSafeString(req.LastName) {
			c.JSON(400, gin.H{
				"error": true,
				"message": "field can't contain sql injection",
			})
			return
		}

		if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" {
			c.JSON(400, gin.H{
				"error": true,
				"message": "missing field",
			})
			return
		}

		match, _ := regexp.MatchString("^[a-zA-Z0-9._-]+@[a-zA-Z0-9._-]+\\.[a-z]{2,6}$", req.Email)
		if !match {
			c.JSON(400, gin.H{
				"error": true,
				"message": "wrong email format",
			})
			return
		}

		if len(req.Password) < 0 || len(req.Password) > 255 {
			c.JSON(400, gin.H{
				"error": true,
				"message": "wrong password length",
			})
			return
		}

		if len(req.FirstName) < 0 || len(req.FirstName) > 50 {
			c.JSON(400, gin.H{
				"error": true,
				"message": "wrong firstname length",
			})
			return
		}

		if len(req.LastName) < 0 || len(req.LastName) > 50 {
			c.JSON(400, gin.H{
				"error": true,
				"message": "wrong lastname length",
			})
			return
		}

		if req.Language <= 0 {
			c.JSON(400, gin.H{
				"error": true,
				"message": "wrong language id",
			})
			return
		}

		var idlanguage int

		err = db.QueryRow("SELECT Id_LANGUAGES FROM LANGUAGES WHERE Id_LANGUAGES = '" + strconv.Itoa(req.Language) + "'").Scan(&idlanguage)
		if err != nil {
			c.JSON(400, gin.H{
				"error": true,
				"message": "wrong language id",
			})
			return
		}

		if typeOfUser[0] == "Client" {
			if req.Country == "" || req.Subscription <= 0 {
				c.JSON(400, gin.H{
					"error": true,
					"message": "missing field",
				})
				return
			}

			if len(req.Country) < 0 || len(req.Country) > 100 {
				c.JSON(400, gin.H{
					"error": true,
					"message": "wrong country length",
				})
				return
			}

			if !utils.IsSafeString(req.Country) {
				c.JSON(400, gin.H{
					"error": true,
					"message": "field can't contain sql injection",
				})
				return
			}

			var id int

			err = db.QueryRow("SELECT Id_SUBSCRIPTIONS FROM SUBSCRIPTIONS WHERE Id_SUBSCRIPTIONS = '" + strconv.Itoa(req.Subscription) + "'").Scan(&id)
			if err != nil {
				c.JSON(400, gin.H{
					"error": true,
					"message": "wrong subscription id",
				})
				return
			}

		} else if typeOfUser[0] == "Manager" {
			if req.IsItemManager == false && req.IsClientManager == false && req.IsContractorManager == false && req.IsSuperAdmin == false {
				c.JSON(400, gin.H{
					"error": true,
					"message": "manager must have at least one role",
				})
				return
			}
		} else if typeOfUser[0] == "Contractor" {
			if req.Presentation == "" || req.ContractStart == "" || req.ContractEnd == "" {
				c.JSON(400, gin.H{
					"error": true,
					"message": "missing field",
				})
				return
			}

			if !utils.IsSafeString(req.Presentation) || !utils.IsSafeString(req.ContractStart) || !utils.IsSafeString(req.ContractEnd) {
				c.JSON(400, gin.H{
					"error": true,
					"message": "field can't contain sql injection",
				})
				return
			}

			var idcontractortype int

			err = db.QueryRow("SELECT Id_CONTRACTOR_TYPES FROM CONTRACTOR_TYPES WHERE Id_CONTRACTOR_TYPES = '" + strconv.Itoa(req.IdContractorType) + "'").Scan(&idcontractortype)
			if err != nil {
				c.JSON(400, gin.H{
					"error": true,
					"message": "wrong contractortype id",
				})
				return
			}


		}

		result, err := db.Exec("INSERT INTO USERS VALUES(NULL, '" + req.Email + "', '" + req.Password + "', '" + req.FirstName + "', '" + req.LastName + "', DEFAULT, DEFAULT, DEFAULT, DEFAULT, '" + strconv.Itoa(req.Language) + "')")
		fmt.Println(err)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot insert user in bdd",
			})
			return
		}

		lastId, err := result.LastInsertId()
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot get user id",
			})
			return
		}

		if typeOfUser[0] == "Client" {
			rows, err := db.Exec("INSERT INTO CLIENTS VALUES(NULL, DEFAULT, '', '" + req.Country + "', '', '', '', '" + strconv.Itoa(req.Subscription) + "', DEFAULT, '" + strconv.FormatInt(lastId, 10) + "')")
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error": true,
					"message": "error on query request to bdd",
				})

				_, err = db.Exec("DELETE FROM USERS WHERE Id_USERS=" + strconv.FormatInt(lastId, 10))
				if err != nil {
					c.JSON(500, gin.H{
						"error": true,
						"message": "cannot delete user",
					})
					return
				}
				return
			}

			lastIdClient, err := rows.LastInsertId()
				if err != nil {
					c.JSON(500, gin.H{
						"error": true,
						"message": "cannot get client id",
					})
					return
				}

			_, err = db.Query("INSERT INTO IS_SUBSCRIBED VALUES(?, ?, DATE_ADD(NOW(), INTERVAL 1 MONTH))", strconv.FormatInt(lastIdClient, 10), strconv.Itoa(req.Subscription))
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "cannot insert client subscription",
				})
				return
			}

			c.JSON(200, gin.H{
				"error": false,
				"iduser": lastId,
				"idclient": lastIdClient,
				"message": "client created",
			})
			return

		} else if typeOfUser[0] == "Manager" {
			rows, err := db.Exec("INSERT INTO MANAGERS VALUES(NULL, ?, ?, ?, ?, ?)", req.IsItemManager, req.IsClientManager, req.IsContractorManager, req.IsSuperAdmin, strconv.FormatInt(lastId, 10))
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "error on query request to bdd",
				})

				_, err = db.Exec("DELETE FROM USERS WHERE Id_USERS=" + strconv.FormatInt(lastId, 10))
				if err != nil {
					c.JSON(500, gin.H{
						"error": true,
						"message": "cannot delete user",
					})
					return
				}
				return
			}

			lastIdManager, err := rows.LastInsertId()

			c.JSON(200, gin.H{
				"error": false,
				"iduser": lastId,
				"idmanager": lastIdManager,
				"message": "manager created",
			})
			return

		} else if typeOfUser[0] == "Contractor" {
			rows, err := db.Exec("INSERT INTO CONTRACTORS VALUES(NULL, '" + req.Presentation + "', '" + req.ContractStart + "', '" + req.ContractEnd + "', '" + strconv.Itoa(req.IdContractorType) + "', '" + strconv.FormatInt(lastId, 10) + "')")
			fmt.Println(err)
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "error on query request to bdd",
				})
				
				_, err = db.Exec("DELETE FROM USERS WHERE Id_USERS=" + strconv.FormatInt(lastId, 10))
				if err != nil {
					c.JSON(500, gin.H{
						"error": true,
						"message": "cannot delete user",
					})
					return
				}
				return
			}

			lastIdContractor, err := rows.LastInsertId()

			c.JSON(200, gin.H{
				"error": false,
				"iduser": lastId,
				"idcontractor": lastIdContractor,
				"message": "contractor created",
			})
			return
		}
	}
}

func UpdateUser(tokenAPI string) func(c *gin.Context) {
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

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		var req User

		err = c.BindJSON(&req)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "can't decode json request",
			})
			return
		}

		var setClause []string

		if req.FirstName != "" {
			if utils.IsSafeString(req.FirstName) == false {
				c.JSON(400, gin.H{
					"error": true,
					"message": "wrong firstname format",
				})
				return
			}
			if len(req.FirstName) < 0 || len(req.FirstName) > 50 {
				c.JSON(400, gin.H{
					"error": true,
					"message": "wrong firstname length",
				})
				return
			}
			setClause = append(setClause, fmt.Sprintf("firstname = '%s'", req.FirstName))
		}
		if req.LastName != "" {
			if utils.IsSafeString(req.LastName) == false {
				c.JSON(400, gin.H{
					"error": true,
					"message": "wrong lastname format",
				})
				return
			}
			if len(req.LastName) < 0 || len(req.LastName) > 50 {
				c.JSON(400, gin.H{
					"error": true,
					"message": "wrong lastname length",
				})
				return
			}
			setClause = append(setClause, fmt.Sprintf("lastname = '%s'", req.LastName))
		}
		if req.Email != "" {
			if utils.IsSafeString(req.Email) == false {
				c.JSON(400, gin.H{
					"error": true,
					"message": "wrong email format",
				})
				return
			}
			if len(req.Email) < 0 || len(req.Email) > 100 {
				c.JSON(400, gin.H{
					"error": true,
					"message": "wrong email length",
				})
				return
			}
			match, _ := regexp.MatchString("^[a-zA-Z0-9._-]+@[a-zA-Z0-9._-]+\\.[a-z]{2,6}$", req.Email)
			if !match {
				c.JSON(400, gin.H{
					"error": true,
					"message": "wrong email format",
				})
				return
			}
			setClause = append(setClause, fmt.Sprintf("email = '%s'", req.Email))
		}
		if req.Password != "" {
			if utils.IsSafeString(req.Password) == false {
				c.JSON(400, gin.H{
					"error": true,
					"message": "wrong password format",
				})
				return
			}
			if len(req.Password) < 0 || len(req.Password) > 255 {
				c.JSON(400, gin.H{
					"error": true,
					"message": "wrong password length",
				})
				return
			}
			setClause = append(setClause, fmt.Sprintf("password = '%s'", req.Password))
		}
		if req.ProfilePicture != "" {
			if utils.IsSafeString(req.ProfilePicture) == false {
				c.JSON(400, gin.H{
					"error": true,
					"message": "wrong profilepicture format",
				})
				return
			}
			if len(req.ProfilePicture) < 0 || len(req.ProfilePicture) > 100 {
				c.JSON(400, gin.H{
					"error": true,
					"message": "wrong profilepicture length",
				})
				return
			}
			setClause = append(setClause, fmt.Sprintf("profilepicture = '%s'", req.ProfilePicture))
		}

		if req.IsCreatedAt != "" {
			if utils.IsSafeString(req.IsCreatedAt) == false {
				c.JSON(400, gin.H{
					"error": true,
					"message": "wrong iscreatedat format",
				})
				return
			}
			setClause = append(setClause, fmt.Sprintf("iscreatedat = '%s'", req.IsCreatedAt))
		}

		if req.LastSeen != "" {
			if utils.IsSafeString(req.LastSeen) == false {
				c.JSON(400, gin.H{
					"error": true,
					"message": "wrong lastseen format",
				})
				return
			}
			setClause = append(setClause, fmt.Sprintf("lastseen = '%s'", req.LastSeen))
		}

		if req.IsBlocked != "" {
			if utils.IsSafeString(req.IsBlocked) == false {
				c.JSON(400, gin.H{
					"error": true,
					"message": "wrong isblocked format",
				})
				return
			}
			setClause = append(setClause, fmt.Sprintf("isblocked = '%s'", req.IsBlocked))
		}

		if req.Language != 0 {
			var idlanguage int

			err = db.QueryRow("SELECT Id_LANGUAGES FROM LANGUAGES WHERE Id_LANGUAGES = '" + strconv.Itoa(req.Language) + "'").Scan(&idlanguage)
			if err != nil {
				c.JSON(400, gin.H{
					"error": true,
					"message": "wrong language id",
				})
				return
			}
			setClause = append(setClause, fmt.Sprintf("id_languages = '%s'", strconv.Itoa(req.Language)))
		}

		if len(setClause) == 0 {
			c.JSON(400, gin.H{
				"error": true,
				"message": "no field to update",
			})
			return
		}

		var user User

		err = db.QueryRow("SELECT * FROM USERS WHERE Id_USERS=" + id).Scan(&user.Id, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.ProfilePicture, &user.IsCreatedAt, &user.LastSeen, &user.IsBlocked, &user.Language)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error": true,
				"message": "user not found",
			})
			return
		}

		_, err = db.Exec("UPDATE USERS SET " + strings.Join(setClause, ", ") + " WHERE Id_USERS=" + id)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot update user",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
			"message": "user updated",
		})
		return
	}
}

func GetUserByFilter(tokenAPI string) func(c *gin.Context) {
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

		filter := c.Param("filter")
		if filter == "" {
			c.JSON(400, gin.H{
				"error": true,
				"message": "id can't be empty",
			})
		}

		if !utils.IsSafeString(filter) {
			c.JSON(400, gin.H{
				"error": true,
				"message": "wrong filter format",
			})
			return
		}
		
		query := fmt.Sprintf("SELECT USERS.Id_USERS, USERS.email, USERS.firstname, USERS.lastname, USERS.profilepicture, USERS.iscreatedat, USERS.lastseen, USERS.isblocked, USERS.Id_LANGUAGES FROM USERS WHERE lastname LIKE '%%%s%%' OR firstname LIKE '%%%s%%' OR email LIKE '%%%s%%'", filter, filter, filter)

		rows, err := db.Query(query)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "error on query request to bdd",
			})
			return
		}
		defer rows.Close()

		var users []UserNoPassword

		for rows.Next() {
			var user UserNoPassword
			err = rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.ProfilePicture, &user.IsCreatedAt, &user.LastSeen, &user.IsBlocked, &user.Language)
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "error on scan rows",
				})
				return
			}
			users = append(users, user)
		}

		if len(users) == 0 {
			c.JSON(404, gin.H{
				"error": true,
				"message": "no user found",
			})
			return
		}

		c.JSON(200, users)
		return
	}
}

func GetUsers(tokenAPI string) func(c *gin.Context) {
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

		rows, err := db.Query("SELECT USERS.Id_USERS, USERS.email, USERS.firstname, USERS.lastname, USERS.profilepicture, USERS.iscreatedat, USERS.lastseen, USERS.isblocked, USERS.Id_LANGUAGES FROM USERS ORDER BY lastname")
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "users not found",
			})
			return
		}

		var users []UserNoPassword

		for rows.Next() {
			var user UserNoPassword
			err = rows.Scan(&user.Id, &user.Email, &user.FirstName, &user.LastName, &user.ProfilePicture, &user.IsCreatedAt, &user.LastSeen, &user.IsBlocked, &user.Language)
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "err of scan rows",
				})
				return
			}
			users = append(users, user)
		}

		c.JSON(200, users)
		return
	}
}


func LoginUser(tokenAPI string) func(c *gin.Context) {
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

		type Login struct {
			Email string `json:"email"`
			Password string `json:"password"`
		}

		var login Login

		err = c.BindJSON(&login)
		if err != nil {
			c.JSON(400, gin.H{
				"error": true,
				"message": "bad json",
			})
			return
		}

		if !utils.IsSafeString(login.Email) || !utils.IsSafeString(login.Password) {
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

		var id int

		err = db.QueryRow("SELECT Id_USERS FROM USERS WHERE email = '" + login.Email + "'").Scan(&id)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "email not found",
			})
			return
		}

		var isblocked string

		err = db.QueryRow("SELECT Id_USERS, isBlocked FROM USERS WHERE Email = '" + login.Email + "' AND Password = '" + login.Password + "'").Scan(&id, &isblocked)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "wrong password",
			})
			return
		}

		var idEntity int
		var idSub int

		err = db.QueryRow("SELECT Id_CLIENTS, subscription FROM CLIENTS WHERE Id_USERS = (SELECT Id_USERS FROM USERS WHERE Email = '" + login.Email + "' AND Password = '" + login.Password + "')").Scan(&idEntity, &idSub)
		if err == nil {
			c.JSON(200, gin.H{
				"error": false,
				"role": "client",
				"subscription": idSub,
				"id": id,
				"isblocked": isblocked,
			})
			return
		}

		err = db.QueryRow("SELECT Id_MANAGERS FROM MANAGERS WHERE Id_USERS = (SELECT Id_USERS FROM USERS WHERE Email = '" + login.Email + "' AND Password = '" + login.Password + "')").Scan(&idEntity)
		if err == nil {
			c.JSON(200, gin.H{
				"error": false,
				"role": "manager",
				"id": id,
				"isblocked": isblocked,
			})
			return
		}

		err = db.QueryRow("SELECT Id_CONTRACTORS FROM CONTRACTORS WHERE Id_USERS = (SELECT Id_USERS FROM USERS WHERE Email = '" + login.Email + "' AND Password = '" + login.Password + "')").Scan(&idEntity)
		if err == nil {
			c.JSON(200, gin.H{
				"error": false,
				"role": "contractor",
				"id": id,
				"isblocked": isblocked,
			})
			return
		}

		c.JSON(500, gin.H{
			"error": false,
			"message": "user as no role",
		})
		return
	}
}

func DeleteUser(tokenAPI string) func(c *gin.Context) {
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

		typeOfUser := c.Request.Header["Type"]
		if typeOfUser == nil{
			c.JSON(498, gin.H{
				"error": true,
				"message": "missing type",
			})
		}

		if typeOfUser[0] != "Client" && typeOfUser[0] != "Contractor" && typeOfUser[0] != "Manager" {
			c.JSON(498, gin.H{
				"error": true,
				"message": "wrong type",
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

		if typeOfUser[0] == "Client" {
			var idClient int

			err = db.QueryRow("SELECT Id_CLIENTS FROM CLIENTS WHERE Id_USERS = '" + id + "'").Scan(&idClient)
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "client not found",
				})
				return
			}

			_, err = db.Exec("DELETE FROM IS_SUBSCRIBED WHERE Id_CLIENTS=" + strconv.Itoa(idClient))
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "cannot delete client subscription",
				})
				return
			}

			_, err = db.Exec("DELETE FROM WATCHES WHERE Id_CLIENTS=" + strconv.Itoa(idClient))
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "cannot delete client watch",
				})
				return
			}

			_, err = db.Exec("DELETE FROM BOOKS WHERE Id_CLIENTS=" + strconv.Itoa(idClient))
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "cannot delete client watch",
				})
				return
			}
			_, err = db.Exec("DELETE FROM PARTICIPATES WHERE Id_CLIENTS=" + strconv.Itoa(idClient))
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "cannot delete client books",
				})
				return
			}

			_, err = db.Exec("DELETE FROM COMMENTS WHERE Id_CLIENTS=" + strconv.Itoa(idClient))
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "cannot delete client comments",
				})
				return
			}

			rows, err := db.Query("SELECT Id_ORDERS FROM ORDERS WHERE Id_CLIENTS=" + strconv.Itoa(idClient))
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "cannot get client orders",
				})
				return
			}

			var idOrder int

			for rows.Next() {
				err = rows.Scan(&idOrder)
				if err != nil {
					c.JSON(500, gin.H{
						"error": true,
						"message": "cannot scan client orders",
					})
					return
				}

				_, err = db.Exec("DELETE FROM CONTAINS_ITEM WHERE Id_ORDERS=" + strconv.Itoa(idOrder))
				if err != nil {
					fmt.Println(err)
					c.JSON(500, gin.H{
						"error": true,
						"message": "cannot delete client orders items",
					})
					return
				}

				_, err = db.Exec("DELETE FROM CONTAINS_FOOD WHERE Id_ORDERS=" + strconv.Itoa(idOrder))
				if err != nil {
					c.JSON(500, gin.H{
						"error": true,
						"message": "cannot delete client orders food",
					})
					return
				}
			}

			_, err = db.Exec("DELETE FROM ORDERS WHERE Id_CLIENTS=" + strconv.Itoa(idClient))
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error": true,
					"message": "cannot delete client orders",
				})
				return
			}

			_, err = db.Exec("DELETE FROM FORMATIONS WHERE Id_CLIENTS=" + strconv.Itoa(idClient))
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error": true,
					"message": "cannot delete client formations",
				})
				return
			}

			_, err = db.Exec("DELETE FROM CLIENTS WHERE Id_CLIENTS=" + strconv.Itoa(idClient))
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error": true,
					"message": "cannot delete client",
				})
				return
			}
		} else if typeOfUser[0] == "Manager" {
			var idManager int

			err = db.QueryRow("SELECT Id_MANAGERS FROM MANAGERS WHERE Id_USERS = '" + id + "'").Scan(&idManager)
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "manager not found",
				})
				return
			}

			_, err = db.Exec("UPDATE ORGANIZES SET Id_MANAGERS = 1 WHERE Id_MANAGERS =" + strconv.Itoa(idManager))
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "cannot update organizes",
				})
				return
			}

			_, err = db.Exec("DELETE FROM MANAGERS WHERE Id_MANAGERS=" + strconv.Itoa(idManager))
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error": true,
					"message": "cannot delete manager",
				})
				return
			}

		} else if typeOfUser[0] == "Contractor" {

			var idContractor int

			err = db.QueryRow("SELECT Id_CONTRACTORS FROM CONTRACTORS WHERE Id_USERS = '" + id + "'").Scan(&idContractor)
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "contractor not found",
				})
				return
			}

			rows, err := db.Query("SELECT Id_ORDERS FROM ORDERS WHERE Id_CONTRACTORS =" + strconv.Itoa(idContractor))
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "cannot get client orders",
				})
				return
			}

			var idOrder int

			for rows.Next() {
				err = rows.Scan(&idOrder)
				if err != nil {
					c.JSON(500, gin.H{
						"error": true,
						"message": "cannot scan client orders",
					})
					return
				}

				_, err = db.Exec("DELETE FROM CONTAINS_ITEM WHERE Id_ORDERS=" + strconv.Itoa(idOrder))
				if err != nil {
					fmt.Println(err)
					c.JSON(500, gin.H{
						"error": true,
						"message": "cannot delete client orders items",
					})
					return
				}

				_, err = db.Exec("DELETE FROM CONTAINS_FOOD WHERE Id_ORDERS=" + strconv.Itoa(idOrder))
				if err != nil {
					c.JSON(500, gin.H{
						"error": true,
						"message": "cannot delete client orders food",
					})
					return
				}
			}

			_, err = db.Exec("DELETE FROM ORDERS WHERE Id_CONTRACTORS=" + strconv.Itoa(idContractor))
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error": true,
					"message": "cannot delete client orders",
				})
				return
			}

			rows, err = db.Query("SELECT Id_ORDERS FROM ORDERS WHERE Id_CONTRACTORS_1=" + strconv.Itoa(idContractor))
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "cannot get client orders",
				})
				return
			}

			for rows.Next() {
				err = rows.Scan(&idOrder)
				if err != nil {
					c.JSON(500, gin.H{
						"error": true,
						"message": "cannot scan client orders",
					})
					return
				}

				_, err = db.Exec("DELETE FROM CONTAINS_ITEM WHERE Id_ORDERS=" + strconv.Itoa(idOrder))
				if err != nil {
					fmt.Println(err)
					c.JSON(500, gin.H{
						"error": true,
						"message": "cannot delete client orders items",
					})
					return
				}

				_, err = db.Exec("DELETE FROM CONTAINS_FOOD WHERE Id_ORDERS=" + strconv.Itoa(idOrder))
				if err != nil {
					c.JSON(500, gin.H{
						"error": true,
						"message": "cannot delete client orders food",
					})
					return
				}
			}

			_, err = db.Exec("DELETE FROM ORDERS WHERE Id_CONTRACTORS_1=" + strconv.Itoa(idContractor))
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error": true,
					"message": "cannot delete client orders",
				})
				return
			}

			_, err = db.Exec("UPDATE ANIMATES SET Id_CONTRACTORS = 1 WHERE Id_CONTRACTORS =" + strconv.Itoa(idContractor))
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "cannot update animates",
				})
				return
			}

			_, err = db.Exec("UPDATE TEACHES SET Id_CONTRACTORS = 1 WHERE Id_CONTRACTORS =" + strconv.Itoa(idContractor))
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "cannot update teaches",
				})
				return
			}

			_, err = db.Exec("UPDATE ORDERS SET Id_CONTRACTORS = 1 WHERE Id_CONTRACTORS =" + strconv.Itoa(idContractor))
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "cannot update orders",
				})
				return
			}

			_, err = db.Exec("UPDATE ORDERS SET Id_CONTRACTORS = 1 WHERE Id_CONTRACTORS_1 =" + strconv.Itoa(idContractor))
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "cannot update orders",
				})
				return
			}

			_, err = db.Exec("DELETE FROM CONTRACTORS WHERE Id_CONTRACTORS=" + strconv.Itoa(idContractor))
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error": true,
					"message": "cannot delete contractor",
				})
				return
			}
		}

		_, err = db.Exec("DELETE FROM BILLS WHERE Id_USERS=" + id)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot delete user",
			})
			return
		}

		_, err = db.Exec("DELETE FROM MESSAGES WHERE Id_USERS=" + id)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot delete user",
			})
			return
		}

		_, err = db.Exec("DELETE FROM MESSAGES WHERE Id_USERS_1=" + id)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot delete user",
			})
			return
		}

		_, err = db.Exec("DELETE FROM USERS WHERE Id_USERS=" + id)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot delete user",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
			"id": id,
			"message": "user deleted",
		})
		return
	}
}

func GetPasswordByEmail(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		type Email struct {
			Email string `json:"email"`
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
		
		var email Email

		err = c.BindJSON(&email)
		if err != nil {
			c.JSON(400, gin.H{
				"error": true,
				"message": "bad json",
			})
			return
		}

		if !utils.IsSafeString(email.Email) {
			c.JSON(400, gin.H{
				"error": true,
				"message": "email can't contain sql injection",
			})
			return
		}

		match, _ := regexp.MatchString("^[a-zA-Z0-9._-]+@[a-zA-Z0-9._-]+\\.[a-z]{2,6}$", email.Email)
		if !match {
			c.JSON(400, gin.H{
				"error": true,
				"message": "wrong email format",
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

		var password string

		err = db.QueryRow("SELECT Password FROM USERS WHERE Email = '" + email.Email + "'").Scan(&password)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "email not found",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
			"password": password,
		})
		return
	}
}