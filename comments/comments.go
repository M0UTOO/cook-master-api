package comments

import (
	"cook-master-api/token"
	"cook-master-api/utils"
	"database/sql"
	"strconv"
	"strings"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Comment struct {
	IdComment      int     `json:"idcomment"`
	Grade 	  float64  `json:"grade"`
	Content string  `json:"content"`
	Picture	   string  `json:"picture"`
	IdClient 	 int     `json:"idclient"`
	IdEvent int `json:"idevent"`
}

func GetCommentsByClientID(tokenAPI string) func(c *gin.Context) {
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

		rows, err := db.Query("SELECT * FROM COMMENTS WHERE ID_CLIENTS = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't query database",
			})
			return
		}

		var comments []Comment
		for rows.Next() {
			var comment Comment
			err = rows.Scan(&comment.IdComment, &comment.Grade, &comment.Content, &comment.Picture, &comment.IdClient, &comment.IdEvent)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "can't scan rows",
				})
				return
			}
			comments = append(comments, comment)
		}

		c.JSON(200, comments)
	}
}

func GetCommentsByEventID(tokenAPI string) func(c *gin.Context) {
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

		rows, err := db.Query("SELECT * FROM COMMENTS WHERE ID_EVENTS = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't query database",
			})
			return
		}

		var comments []Comment
		for rows.Next() {
			var comment Comment
			err = rows.Scan(&comment.IdComment, &comment.Grade, &comment.Content, &comment.Picture, &comment.IdClient, &comment.IdEvent)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error":   true,
					"message": "can't scan rows",
				})
				return
			}
			comments = append(comments, comment)
		}

		c.JSON(200, comments)
	}
}

func PostComment(tokenAPI string) func(c *gin.Context) {
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

		var comment Comment
		err = c.BindJSON(&comment)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "can't bind json",
			})
			return
		}

		if comment.Grade < 1 || comment.Grade > 5 {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "grade must be between 1 and 5",
			})
			return
		}

		if comment.Content == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "content can't be empty",
			})
			return
		}

		if !utils.IsSafeString(comment.Content) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "content can't contain sql injection",
			})
			return
		}

		if !utils.IsSafeString(comment.Picture) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "picture can't contain sql injection",
			})
			return
		}

		if len(comment.Picture) > 255 || len(comment.Picture) < 0 {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "picture can't be more than 255 characters or less than 0",
			})
			return
		}

		var id int

		err = db.QueryRow("SELECT Id_USERS FROM CLIENTS WHERE Id_USERS = '" + strconv.Itoa(comment.IdClient) + "'").Scan(&id)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "client doesn't exist",
			})
			return
		}

		err = db.QueryRow("SELECT Id_EVENTS FROM EVENTS WHERE Id_EVENTS = '" + strconv.Itoa(comment.IdEvent) + "'").Scan(&id)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "event doesn't exist",
			})
			return
		}

		stmt, err := db.Prepare("INSERT INTO COMMENTS (Grade, Content, Picture, Id_CLIENTS, Id_EVENTS) VALUES (?, ?, ?, ?, ?)")
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't prepare statement",
			})
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(comment.Grade, comment.Content, comment.Picture, comment.IdClient, comment.IdEvent)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't execute statement",
			})
			return
		}

		c.JSON(200, gin.H{
			"error":   false,
			"message": "comment added",
		})
	}
}

func DeleteComment(tokenAPI string) func(c *gin.Context) {
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

		var idComment int

		err = db.QueryRow("SELECT Id_COMMENTS FROM COMMENTS WHERE Id_COMMENTS = '" + id + "'").Scan(&idComment)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "comment doesn't exist",
			})
			return
		}

		stmt, err := db.Prepare("DELETE FROM COMMENTS WHERE Id_COMMENTS = ?")
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't prepare statement",
			})
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(idComment)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't execute statement",
			})
			return
		}

		c.JSON(200, gin.H{
			"error":   false,
			"message": "comment deleted",
		})
	}
}

func UpdateComment(tokenAPI string) func(c *gin.Context) {
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

		var idComment int

		err = db.QueryRow("SELECT Id_COMMENTS FROM COMMENTS WHERE Id_COMMENTS = '" + id + "'").Scan(&idComment)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "comment doesn't exist",
			})
			return
		}

		var comment Comment

		err = c.BindJSON(&comment)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "can't bind json",
			})
			return
		}

		var setClause []string

		if comment.Grade != 0 {
			if comment.Grade < 1 || comment.Grade > 5 {
				c.JSON(400, gin.H{
					"error":   true,
					"message": "grade must be between 1 and 5",
				})
				return
			}
			setClause = append(setClause, "Grade = '"+strconv.FormatFloat(comment.Grade, 'f', 2, 64)+"'")
		}

		if comment.Content != "" {
			if !utils.IsSafeString(comment.Content) {
				c.JSON(400, gin.H{
					"error":   true,
					"message": "content can't contain sql injection",
				})
				return
			}
			setClause = append(setClause, "Content = '"+comment.Content+"'")
		}

		if comment.Picture != "" {
			if !utils.IsSafeString(comment.Picture) {
				c.JSON(400, gin.H{
					"error":   true,
					"message": "picture can't contain sql injection",
				})
				return
			}

			if len(comment.Picture) > 255 || len(comment.Picture) < 0 {
				c.JSON(400, gin.H{
					"error":   true,
					"message": "picture can't be more than 255 characters or less than 0",
				})
				return
			}
			setClause = append(setClause, "Picture = '"+comment.Picture+"'")
		}

		if len(setClause) == 0 {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "no value to update",
			})
			return
		}

		stmt, err := db.Prepare("UPDATE COMMENTS SET " + strings.Join(setClause, ", ") + " WHERE Id_COMMENTS = ?")
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't prepare statement",
			})
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(idComment)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "can't execute statement",
			})
			return
		}

		c.JSON(200, gin.H{
			"error":   false,
			"message": "comment updated",
		})
	}
}
