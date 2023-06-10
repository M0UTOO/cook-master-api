package languages

import (
	"cook-master-api/token"
	"cook-master-api/utils"
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Language struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func GetLanguages(tokenAPI string) func(c *gin.Context) {
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

		rows, err := db.Query("SELECT * FROM LANGUAGES")
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't query database",
			})
			return
		}

		var languages []Language
		for rows.Next() {
			var language Language
			err := rows.Scan(&language.Id, &language.Name)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "can't scan rows",
				})
				return
			}
			languages = append(languages, language)
		}

		c.JSON(200, languages)
	}
}

func GetLanguageByID(tokenAPI string) func(c *gin.Context) {
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

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		rows, err := db.Query("SELECT * FROM LANGUAGES WHERE Id_LANGUAGES = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't query database",
			})
			return
		}

		var language Language
		for rows.Next() {
			err := rows.Scan(&language.Id, &language.Name)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "can't scan rows",
				})
				return
			}
		}

		c.JSON(200, language)
	}
}

func PostLanguage(tokenAPI string) func(c *gin.Context) {
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

		var language Language
		err = c.BindJSON(&language)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't bind json",
			})
			return
		}

		if language.Name == "" {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "missing name",
			})
			return
		}

		if !utils.IsSafeString(language.Name) {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "invalid name",
			})
			return
		}

		if len(language.Name) > 100 {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "name too long",
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

		rows, err := db.Exec("INSERT INTO LANGUAGES (Name) VALUES (?)", language.Name)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't insert into database",
			})
			return
		}

		id, err := rows.LastInsertId()

		c.JSON(200, gin.H{
			"error": false,
			"id": id,
			"message": "language added",
		})
	}
}

func UpdateLanguage(tokenAPI string) func(c *gin.Context) {
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

		var language Language
		err = c.BindJSON(&language)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't bind json",
			})
			return
		}

		if language.Name == "" {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "missing name",
			})
			return
		}

		if !utils.IsSafeString(language.Name) {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "invalid name",
			})
			return
		}

		if len(language.Name) > 100 {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "name too long",
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

		var idlanguage int

		err = db.QueryRow("SELECT Id_LANGUAGES FROM LANGUAGES WHERE Id_LANGUAGES = ?", id).Scan(&idlanguage)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "wrong id",
			})
			return
		}

		_, err = db.Exec("UPDATE LANGUAGES SET Name = ? WHERE Id_LANGUAGES = ?", language.Name, id)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't update database",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
			"message": "language updated",
		})
	}
}

func DeleteLanguage(tokenAPI string) func(c *gin.Context) {
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

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		var idlanguage int

		err = db.QueryRow("SELECT Id_LANGUAGES FROM LANGUAGES WHERE Id_LANGUAGES = ?", id).Scan(&idlanguage)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "wrong id",
			})
			return
		}

		_, err = db.Exec("DELETE FROM LANGUAGES WHERE Id_LANGUAGES = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't delete from database",
			})
			return
		}

		c.JSON(200, gin.H{
			"error":   false,
			"message": "language deleted",
		})
	}
}
