package lessons

import (
	"cook-master-api/token"
	"cook-master-api/utils"
	"database/sql"
	"fmt"
	"strings"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Lesson struct {
	IdLesson          int    `json:"idlesson"`
	Name              string `json:"name"`
	Content           string `json:"content"`
	Description       string `json:"description"`
	Difficulty        int    `json:"difficulty"`
	GroupDisplayOrder int    `json:"group_display_order"`
	IdLessonGroup     int    `json:"idlessongroup"`
}

func GetLessons(tokenAPI string) func(c *gin.Context) {
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

		rows, err := db.Query("SELECT * FROM LESSONS")
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot get lessons",
			})
			return
		}
		defer rows.Close()

		var lessons []Lesson

		for rows.Next() {
			var lesson Lesson
			err = rows.Scan(&lesson.IdLesson, &lesson.Name, &lesson.Content, &lesson.Description, &lesson.Difficulty, &lesson.GroupDisplayOrder, &lesson.IdLessonGroup)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error":   true,
					"message": "err on scan rows",
				})
				return
			}
			lessons = append(lessons, lesson)
		}

		c.JSON(200, lessons)
	}
}

func GetLessonByID(tokenAPI string) func(c *gin.Context) {
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

		var lesson Lesson

		err = db.QueryRow("SELECT * FROM LESSONS WHERE Id_LESSONS = ?", id).Scan(&lesson.IdLesson, &lesson.Name, &lesson.Content, &lesson.Description, &lesson.Difficulty, &lesson.GroupDisplayOrder, &lesson.IdLessonGroup)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot get lesson",
			})
			return
		}

		c.JSON(200, lesson)
	}
}

func GetLessonsByGroupID(tokenAPI string) func(c *gin.Context) {
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

		rows, err := db.Query("SELECT * FROM LESSONS WHERE Id_LESSONS_GROUPS = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot get lessons",
			})
			return
		}
		defer rows.Close()

		var lessons []Lesson

		for rows.Next() {
			var lesson Lesson
			err = rows.Scan(&lesson.IdLesson, &lesson.Name, &lesson.Content, &lesson.Description, &lesson.Difficulty, &lesson.GroupDisplayOrder, &lesson.IdLessonGroup)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "cannot get lessons",
				})
				return
			}
			lessons = append(lessons, lesson)
		}

		c.JSON(200, lessons)
	}
}

func Postlesson(tokenAPI string) func(c *gin.Context) {
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

		var lesson Lesson
		c.BindJSON(&lesson)

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		if lesson.Name == "" || !utils.IsSafeString(lesson.Name) || len(lesson.Name) < 0 || len(lesson.Name) > 50 {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "error on name field",
			})
			return
		}

		if lesson.Content == "" || !utils.IsSafeString(lesson.Content) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "error on content field",
			})
			return
		}

		if lesson.Description == "" || !utils.IsSafeString(lesson.Description) || len(lesson.Description) < 0 || len(lesson.Description) > 50 {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "error on description field",
			})
			return
		}

		if lesson.Difficulty <= 0 || lesson.Difficulty > 5 {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "error on difficulty field",
			})
			return
		}

		var idContractor int

		err = db.QueryRow("SELECT Id_CONTRACTORS FROM CONTRACTORS WHERE Id_USERS = '" + id + "'").Scan(&idContractor)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot get contractor",
			})
			return
		}

		err = db.QueryRow("SELECT Id_LESSONS_GROUPS FROM LESSONS_GROUPS WHERE Id_LESSONS_GROUPS = 1").Scan(&idContractor)
		if err != nil {
			_, err := db.Exec("INSERT INTO LESSONS_GROUPS (name) VALUES (?)", "default")
			fmt.Println(err)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "cannot insert lesson group",
				})
				return
			}
		}

		result, err := db.Exec("INSERT INTO LESSONS (Name, Content, Description, Difficulty, group_display_order, Id_LESSONS_GROUPS) VALUES (?, ?, ?, ?, ?, ?)", lesson.Name, lesson.Content, lesson.Description, lesson.Difficulty, 0, 1)
		fmt.Println(err)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot insert lesson",
			})
			return
		}

		lastId, err := result.LastInsertId()
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot get last id",
			})
			return
		}

		result, err = db.Exec("INSERT INTO TEACHES (Id_CONTRACTORS, Id_LESSONS) VALUES (?, ?)", idContractor, lastId)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot insert teaches",
			})
			return
		}

		c.JSON(200, gin.H{
			"error":   false,
			"message": "lesson inserted",
		})
	}
}

func AddLessonToAGroup(tokenAPI string) func(c *gin.Context) {
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

		type Group struct {
			Name              string `json:"name"`
			GroupDisplayOrder int    `json:"group_display_order"`
		}

		var group Group
		c.BindJSON(&group)

		if group.Name != "" {
			if !utils.IsSafeString(group.Name) {
				c.JSON(400, gin.H{
					"error":   true,
					"message": "name can't contain sql injection",
				})
				return
			}
			if len(group.Name) < 0 || len(group.Name) > 100 {
				c.JSON(400, gin.H{
					"error":   true,
					"message": "wrong name lenght",
				})
			}
		}

		if group.GroupDisplayOrder <= 0 {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "group_display_order can't be empty or negative",
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

		var idGroup int

		err = db.QueryRow("SELECT Id_LESSONS_GROUPS FROM LESSONS_GROUPS WHERE name = '" + group.Name + "'").Scan(&idGroup)
		if err != nil {
			result, err := db.Exec("INSERT INTO LESSONS_GROUPS (name) VALUES (?)", group.Name)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error":   true,
					"message": "cannot insert group",
				})
				return
			}

			lastId, err := result.LastInsertId()
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "cannot get last id",
				})
				return
			}

			result, err = db.Exec("UPDATE LESSONS SET Id_LESSONS_GROUPS = ?, group_display_order = ? WHERE Id_LESSONS = ?", lastId, group.GroupDisplayOrder, id)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error":   true,
					"message": "cannot update lesson",
				})
				return
			}

			c.JSON(200, gin.H{
				"error":   false,
				"message": "lesson added to group",
			})
			return
		}

		var idLesson int

		err = db.QueryRow("SELECT Id_LESSONS FROM LESSONS WHERE Id_LESSONS_GROUPS = '" + strconv.Itoa(idGroup) + "' AND group_display_order = '" + strconv.Itoa(group.GroupDisplayOrder) + "'").Scan(&idLesson)
		if err != nil {
			_, err := db.Exec("UPDATE LESSONS SET Id_LESSONS_GROUPS = ?, group_display_order = ? WHERE Id_LESSONS = ?", idGroup, group.GroupDisplayOrder, id)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "cannot update lesson",
				})
				return
			}

			c.JSON(200, gin.H{
				"error":   false,
				"message": "lesson added to group",
			})
			return
		}

		c.JSON(500, gin.H{
			"error":   false,
			"message": "an lesson is already in this group with this display order",
		})
	}
}

func DeleteLessonFromAGroup(tokenAPI string) func(c *gin.Context) {
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

		_, err = db.Exec("UPDATE LESSONS SET Id_LESSONS_GROUPS = 1, group_display_order = 0 WHERE Id_LESSONS = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot update lesson",
			})
			return
		}

		c.JSON(200, gin.H{
			"error":   false,
			"message": "lesson deleted from group",
		})
	}
}

func DeleteLesson(tokenAPI string) func(c *gin.Context) {
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

		_, err = db.Exec("DELETE FROM TEACHES WHERE Id_LESSONS = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot delete teaches",
			})
			return
		}

		_, err = db.Exec("DELETE FROM WATCHES WHERE Id_LESSONS = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot delete watches",
			})
			return
		}

		_, err = db.Exec("DELETE FROM LESSONS WHERE Id_LESSONS = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot delete lesson",
			})
			return
		}

		c.JSON(200, gin.H{
			"error":   false,
			"message": "lesson deleted",
		})
	}
}

func UpdateLesson(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		var setClause []string

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

		var lesson Lesson
		c.BindJSON(&lesson)

		if lesson.Name != "" {
			if !utils.IsSafeString(lesson.Name) {
				c.JSON(400, gin.H{
					"error":   true,
					"message": "name can't contain sql injection",
				})
				return
			}
			if len(lesson.Name) < 0 || len(lesson.Name) > 100 {
				c.JSON(400, gin.H{
					"error":   true,
					"message": "wrong name lenght",
				})
			}
			setClause = append(setClause, "name = '"+lesson.Name+"'")
		}

		if lesson.Content != "" {
			if !utils.IsSafeString(lesson.Content) {
				c.JSON(400, gin.H{
					"error":   true,
					"message": "content can't contain sql injection",
				})
				return
			}
			setClause = append(setClause, "content = '"+lesson.Content+"'")
		}

		if lesson.Description != "" {
			if !utils.IsSafeString(lesson.Description) {
				c.JSON(400, gin.H{
					"error":   true,
					"message": "description can't contain sql injection",
				})
				return
			}
			if len(lesson.Description) < 0 || len(lesson.Description) > 50 {
				c.JSON(400, gin.H{
					"error":   true,
					"message": "wrong description lenght",
				})
			}
			setClause = append(setClause, "description = '"+lesson.Description+"'")
		}

		if lesson.Difficulty != 0 {
			if lesson.Difficulty <= 0 || lesson.Difficulty > 5 {
				c.JSON(400, gin.H{
					"error":   true,
					"message": "error on difficulty field",
				})
				return
			}
			setClause = append(setClause, "difficulty = '"+strconv.Itoa(lesson.Difficulty)+"'")
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

		_, err = db.Exec("UPDATE LESSONS SET " + strings.Join(setClause, ", ") + " WHERE Id_LESSONS = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot update lesson",
			})
			return
		}

		c.JSON(200, gin.H{
			"error":   false,
			"message": "lesson updated",
		})
	}
}
