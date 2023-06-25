package cookingspaces

import (
	"cook-master-api/token"
	"cook-master-api/utils"
	"database/sql"
	"fmt"
	"time"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

const (
	BOOKSSTARTTIME = "07:00:00"
	BOOKSENDTIME   = "20:00:00"
)

type CookingSpace struct {
	IdCookingSpace int     `json:"idCookingSpace"`
	Name           string  `json:"name"`
	Size           int     `json:"size"`
	IsAvailable    bool    `json:"isAvailable"`
	PricePerHour   float64 `json:"pricePerHour"`
	IdPremise      int     `json:"idPremise"`
	Picture 	   string  `json:"picture"`
}

type Books struct {
	IdUser	int `json:"iduser"`
	IdCookingSpace int `json:"idCookingSpace"`
	StartTime string `json:"starttime"`
	EndTime   string `json:"endtime"`
}

func GetCookingSpaces(tokenAPI string) func(c *gin.Context) {
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

		rows, err := db.Query("SELECT * FROM COOKING_SPACES")
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot get cookingspaces",
			})
			return
		}
		defer rows.Close()

		var cookingspaces []CookingSpace

		for rows.Next() {
			var cookingspace CookingSpace
			err = rows.Scan(&cookingspace.IdCookingSpace, &cookingspace.Name, &cookingspace.Size, &cookingspace.IsAvailable, &cookingspace.PricePerHour, &cookingspace.Picture, &cookingspace.IdPremise)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error":   true,
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

		var cookingspace CookingSpace

		err = db.QueryRow("SELECT * FROM COOKING_SPACES WHERE Id_COOKING_SPACES = ?", id).Scan(&cookingspace.IdCookingSpace, &cookingspace.Name, &cookingspace.Size, &cookingspace.IsAvailable, &cookingspace.PricePerHour, &cookingspace.Picture,&cookingspace.IdPremise)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
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

		rows, err := db.Query("SELECT * FROM COOKING_SPACES WHERE Id_PREMISES = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot get cookingspaces",
			})
			return
		}
		defer rows.Close()

		var cookingspaces []CookingSpace

		for rows.Next() {
			var cookingspace CookingSpace
			err = rows.Scan(&cookingspace.IdCookingSpace, &cookingspace.Name, &cookingspace.Size, &cookingspace.IsAvailable, &cookingspace.PricePerHour, &cookingspace.Picture, &cookingspace.IdPremise)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
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

		var cookingspace CookingSpace
		c.BindJSON(&cookingspace)

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		if cookingspace.Name == "" || !utils.IsSafeString(cookingspace.Name) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "name can't be empty or contain sql injection",
			})
			return
		}

		if cookingspace.Size <= 0 {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "size can't be negative or zero",
			})
			return
		}

		if cookingspace.PricePerHour <= 0 {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "price per hour can't be negative or zero",
			})
			return
		}	

		var idPremise int

		err = db.QueryRow("SELECT Id_PREMISES FROM PREMISES WHERE Id_PREMISES = 1").Scan(&idPremise)
		if err != nil {
			_, err := db.Exec("INSERT INTO PREMISES (name, streetNumber, streetName, city, country) VALUES (?, ?, ?, ?, ?)", "default", 0, "default", "default", "default")
			fmt.Println(err)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "cannot insert premise",
				})
				return
			}
		}

		rows, err := db.Exec("INSERT INTO COOKING_SPACES (Name, Size, PricePerHour, picture, Id_PREMISES) VALUES (?, ?, ?, ?, ?)", cookingspace.Name, cookingspace.Size, cookingspace.PricePerHour, cookingspace.Picture, 1)
		fmt.Println(err)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot insert cookingspace",
			})
			return
		}

		id, err := rows.LastInsertId()

		c.JSON(200, gin.H{
			"error":   false,
			"id":      id,
			"message": "cookingspace inserted",
		})
	}
}

func AddCookingSpaceToAPremise(tokenAPI string) func(c *gin.Context) {
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

		type Premise struct {
			Name string `json:"name"`
		}

		var premise Premise
		c.BindJSON(&premise)

		if !utils.IsSafeString(premise.Name) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "name can't contain sql injection",
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

		var idPremise int

		err = db.QueryRow("SELECT Id_PREMISES FROM PREMISES WHERE name = '" + premise.Name + "'").Scan(&idPremise)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot find Premise",
			})
			return
		}

		var idCookingSpace int

		err = db.QueryRow("SELECT Id_COOKING_SPACES FROM COOKING_SPACES WHERE Id_COOKING_SPACES= '" + id + "'").Scan(&idCookingSpace)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot find cookingspace",
			})
			return
		}

		err = db.QueryRow("SELECT Id_COOKING_SPACES FROM COOKING_SPACES WHERE Id_PREMISES = '" + strconv.Itoa(idPremise) + "' AND Id_COOKING_SPACES = '" + id + "'").Scan(&idCookingSpace)
		if err != nil {
			_, err := db.Exec("UPDATE COOKING_SPACES SET Id_PREMISES = ? WHERE Id_COOKING_SPACES = ?", idPremise, id)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "cannot update cookingspace",
				})
				return
			}

			c.JSON(200, gin.H{
				"error":   false,
				"message": "cookingspace added to Premise",
			})
			return
		}

		c.JSON(500, gin.H{
			"error":   false,
			"message": "this cookingspace is already in this premise",
		})
	}
}

func DeleteCookingSpaceFromAPremise(tokenAPI string) func(c *gin.Context) {
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

		_, err = db.Exec("UPDATE COOKING_SPACES SET Id_PREMISES = 1 WHERE Id_COOKING_SPACES = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot update cookingspace",
			})
			return
		}

		c.JSON(200, gin.H{
			"error":   false,
			"message": "cookingspace deleted from premise",
		})
	}
}

func UpdateCookingSpace(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		type CookingSpaceReq struct {
			Name         string  `json:"name"`
			Size         int     `json:"size"`
			IsAvailable  int     `json:"isAvailable"`
			PricePerHour float64 `json:"pricePerHour"`
			Picture 	 string  `json:"picture"`
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

		var cookingspace CookingSpaceReq

		cookingspace.IsAvailable = -1

		err = c.BindJSON(&cookingspace)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "cannot bind json",
			})
			return
		}

		var setClause []string

		if cookingspace.Name != "" {
			if !utils.IsSafeString(cookingspace.Name) {
				c.JSON(400, gin.H{
					"error":   true,
					"message": "name can't contain sql injection",
				})
				return
			}
			if len(cookingspace.Name) < 0 || len(cookingspace.Name) > 50 {
				c.JSON(400, gin.H{
					"error":   true,
					"message": "wrong name length",
				})
				return
			}
			setClause = append(setClause, "name = '"+cookingspace.Name+"'")
		}

		if cookingspace.Picture != "" {
			if !utils.IsSafeString(cookingspace.Picture) {
				c.JSON(400, gin.H{
					"error":   true,
					"message": "picture can't contain sql injection",
				})
				return
			}
			setClause = append(setClause, "picture = '"+cookingspace.Picture+"'")
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
				"error":   true,
				"message": "no field to update",
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

		var idcookingspace string

		err = db.QueryRow("SELECT Id_COOKING_SPACES FROM COOKING_SPACES WHERE Id_COOKING_SPACES = '" + id + "'").Scan(&idcookingspace)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cookingspace not found",
			})
			return
		}

		_, err = db.Exec("UPDATE COOKING_SPACES SET "+strings.Join(setClause, ", ")+" WHERE Id_COOKING_SPACES = ?", id)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot update cookingspace",
			})
			return
		}

		c.JSON(200, gin.H{
			"error":   false,
			"id" : idcookingspace,
			"message": "cookingspace updated",
		})
	}
}

func AddABooks(tokenAPI string) func(c *gin.Context) {
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

		idClient := c.Param("idclient")
		if idClient == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "idclient can't be empty",
			})
			return
		}

		if !utils.IsSafeString(idClient) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "idclient can't contain sql injection",
			})
			return
		}

		idcooking := c.Param("idcookingspace")
		if idcooking == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "idcooking can't be empty",
			})
			return
		}

		if !utils.IsSafeString(idcooking) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "idcooking can't contain sql injection",
			})
			return
		}

		var books Books
		err = c.BindJSON(&books)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "cannot bind json",
			})
			return
		}

		if books.StartTime == "" || !utils.IsSafeString(books.StartTime) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "startime can't be empty or contain sql injection",
			})
			return
		}
		
		startTime, err := time.Parse("2006-01-02 15:04:05", books.StartTime)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "invalid StartTime format",
			})
			return
		}

		referenceTime, _ := time.Parse("2006-01-02 15:04:05", "2006-01-02 06:59:59")

		startTimeOnly := startTime.Format("15:04:05")
		referenceTimeOnly := referenceTime.Format("15:04:05")

		if startTimeOnly <= referenceTimeOnly {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "StartTime should be after 7am",
			})
			return
		}

		endTime, err := time.Parse("2006-01-02 15:04:05", books.EndTime)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "invalid EndTime format",
			})
			return
		}

		referenceTime, _ = time.Parse("2006-01-02 15:04:05", "2006-01-02 20:00:01")

		endTimeOnly := endTime.Format("15:04:05")
		referenceTimeOnly = referenceTime.Format("15:04:05")

		if endTimeOnly >= referenceTimeOnly {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "EndTime should be before 8pm",
			})
			return
		}

		if !startTime.Truncate(24*time.Hour).Equal(endTime.Truncate(24*time.Hour)) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "StartTime and EndTime must be on the same day",
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

		var existingReservationID int
		err = db.QueryRow("SELECT Id_BOOKS FROM BOOKS WHERE startTime <= ? AND endTime >= ?", endTime.Format("2006-01-02 15:04:05"), startTime.Format("2006-01-02 15:04:05")).Scan(&existingReservationID)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("reservation time range is available")
			} else {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "failed to query existing reservations",
				})
				return
			}
		} else {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "reservation time range is not available",
			})
			return
		}

		if books.EndTime == "" || !utils.IsSafeString(books.EndTime) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "endtime can't be empty or contain sql injection",
			})
			return
		}

		var idclient string

		err = db.QueryRow("SELECT Id_CLIENTS FROM CLIENTS WHERE Id_USERS = '" + idClient + "'").Scan(&idclient)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "client not found",
			})
			return
		}

		var idcookingspace string

		err = db.QueryRow("SELECT Id_COOKING_SPACES FROM COOKING_SPACES WHERE Id_COOKING_SPACES = '" + idcooking + "'").Scan(&idcookingspace)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cookingspace not found",
			})
			return
		}

		_, err = db.Exec("INSERT INTO BOOKS (Id_COOKING_SPACES, Starttime, Endtime, Id_CLIENTS) VALUES (?, ?, ?, ?)", idcookingspace, books.StartTime, books.EndTime, idclient)
		fmt.Println(err)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot insert books",
			})
			return
		}

		c.JSON(200, gin.H{
			"error":   false,
			"message": "books added",
		})
	}
}

func DeleteABooks(tokenAPI string) func(c *gin.Context) {
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

		idClient := c.Param("idclient")
		if idClient == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "idclient can't be empty",
			})
			return
		}

		if !utils.IsSafeString(idClient) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "idclient can't contain sql injection",
			})
			return
		}

		idcooking := c.Param("idcookingspace")
		if idcooking == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "idcooking can't be empty",
			})
			return
		}

		if !utils.IsSafeString(idcooking) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "idcooking can't contain sql injection",
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

		var idclient string

		err = db.QueryRow("SELECT Id_CLIENTS FROM CLIENTS WHERE Id_USERS = '" + idClient + "'").Scan(&idclient)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "client not found",
			})
			return
		}

		var idcookingspace string

		err = db.QueryRow("SELECT Id_COOKING_SPACES FROM COOKING_SPACES WHERE Id_COOKING_SPACES = '" + idcooking + "'").Scan(&idcookingspace)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cookingspace not found",
			})
			return
		}

		var idbooks string

		err = db.QueryRow("SELECT Id_CLIENTS FROM BOOKS WHERE Id_COOKING_SPACES = '" + idcookingspace + "' AND Id_CLIENTS = '" + idclient + "'").Scan(&idbooks)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "books not found",
			})
			return
		}

		_, err = db.Exec("DELETE FROM BOOKS WHERE Id_CLIENTS = '" + idclient + "' AND Id_COOKING_SPACES = '" + idcookingspace + "'")
		fmt.Println(err)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot delete books",
			})
			return
		}

		c.JSON(200, gin.H{
			"error":   false,
			"message": "books delete",
		})
	}
}

func DeleteCookingSpace(tokenAPI string) func(c *gin.Context) {
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

		var idcookingspace int
		var picture string

		err = db.QueryRow("SELECT Id_COOKING_SPACES, picture FROM COOKING_SPACES WHERE Id_COOKING_SPACES = ?", id).Scan(&idcookingspace, &picture)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cookingspace doesn't exist",
			})
			return
		}

		_, err = db.Exec("DELETE FROM IS_HOSTED WHERE Id_COOKING_SPACES = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't delete cookingspace",
			})
			return
		}

		_, err = db.Exec("DELETE FROM BOOKS WHERE Id_COOKING_SPACES = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't delete cookingspace",
			})
			return
		}

		_, err = db.Exec("UPDATE COOKING_ITEMS SET Id_COOKING_SPACES = 1 WHERE Id_COOKING_SPACES = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot update cookingitems",
			})
			return
		}

		_, err = db.Exec("UPDATE INGREDIENTS SET Id_COOKING_SPACES = 1 WHERE Id_COOKING_SPACES = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot update ingredients",
			})
			return
		}

		_, err = db.Exec("DELETE FROM COOKING_SPACES WHERE Id_COOKING_SPACES = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't delete cookingspace",
			})
			return
		}

		c.JSON(200, gin.H{
			"error":   false,
			"message": "cookingspace deleted",
			"picture": picture,
		})
	}
}

func GetBooksByCookingSpaceID(tokenAPI string) func(c *gin.Context) {
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
			c.JSON(400, gin.H{
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
				"message": "cannot connect to database",
			})
			return
		}
		defer db.Close()
		rows, err := db.Query("SELECT BOOKS.startTime, BOOKS.endTime, BOOKS.Id_COOKING_SPACES, CLIENTS.Id_USERS FROM BOOKS JOIN CLIENTS ON CLIENTS.Id_CLIENTS = BOOKS.Id_CLIENTS WHERE BOOKS.Id_COOKING_SPACES = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot get books",
			})
			return
		}
		defer rows.Close()
		var books []Books
		for rows.Next() {
			var book Books
			err = rows.Scan(&book.StartTime, &book.EndTime, &book.IdCookingSpace, &book.IdUser)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "cannot get books",
				})
				return
			}
			books = append(books, book)
		}
		c.JSON(200, books)
	}
}

func GetCookingSpacesBooks(tokenAPI string) func(c *gin.Context) {
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
			c.JSON(400, gin.H{
				"error":   true,
				"message": "wrong token",
			})
			return
		}
		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot connect to database",
			})
			return
		}
		defer db.Close()
		rows, err := db.Query("SELECT BOOKS.startTime, BOOKS.endTime, BOOKS.Id_COOKING_SPACES, CLIENTS.Id_USERS FROM BOOKS JOIN CLIENTS ON CLIENTS.Id_CLIENTS = BOOKS.Id_CLIENTS")
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot get books",
			})
			return
		}
		defer rows.Close()
		var books []Books
		for rows.Next() {
			var book Books
			err = rows.Scan(&book.StartTime, &book.EndTime, &book.IdCookingSpace, &book.IdUser)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "cannot get books",
				})
				return
			}
			books = append(books, book)
		}
		c.JSON(200, books)
	}
}

func GetTop5CookingSpaces(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header["Token"]

		type Top5Event struct {
			IdCookingSpace int `json:"idcookingspace"`
			Name string `json:"name"`
			Count        int    `json:"count"`
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

		rows, err := db.Query("SELECT cs.Id_COOKING_SPACES, cs.name AS cookingSpace, COUNT(b.Id_BOOKS) AS bookCount FROM COOKING_SPACES cs INNER JOIN BOOKS b ON b.Id_COOKING_SPACES = cs.Id_COOKING_SPACES GROUP BY cs.name ORDER BY bookCount DESC LIMIT 5;")
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "month not found",
			})
			return
		}

		var top5events []Top5Event

		for rows.Next() {
			var top5event Top5Event
			err = rows.Scan(&top5event.IdCookingSpace, &top5event.Name, &top5event.Count)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error":   true,
					"message": "top 5 not found",
				})
				return
			}
			top5events = append(top5events, top5event)
		}

		c.JSON(200, top5events)
		return
	}
}