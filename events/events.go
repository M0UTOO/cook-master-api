package events

import (
	"fmt"
	"database/sql"
	"strconv"
	"cook-master-api/token"
	"cook-master-api/utils"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Event struct {
	IdEvent int `json:"idevent"`
	Name string `json:"name"`
	Type string `json:"type"`
	EndTime string `json:"endtime"`
	StartTime string `json:"starttime"`
	IsInternal bool `json:"isinternal"`
	IsPrivate bool `json:"isprivate"`
	GroupDisplayOrder int `json:"groupdisplayorder"`
	DefaultPicture string `json:"defaultpicture"`
	IdGroups int `json:"idgroups"`
}

func GetEvents(tokenAPI string) func(c *gin.Context) {
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

		rows, err := db.Query("SELECT * FROM EVENTS")
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot get events",
			})
			return
		}
		defer rows.Close()

		var events []Event

		for rows.Next() {
			var event Event
			err = rows.Scan(&event.IdEvent, &event.Name, &event.Type, &event.EndTime, &event.StartTime, &event.IsInternal, &event.IsPrivate, &event.GroupDisplayOrder, &event.DefaultPicture, &event.IdGroups)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error": true,
					"message": "cannot get events",
				})
				return
			}
			events = append(events, event)
		}

		c.JSON(200, gin.H{
			"error": false,
			"events": events,
		})
	}
}

func GetEventByID(tokenAPI string) func(c *gin.Context) {
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

		var event Event

		err = db.QueryRow("SELECT * FROM EVENTS WHERE Id_EVENTS = ?", id).Scan(&event.IdEvent, &event.Name, &event.Type, &event.EndTime, &event.StartTime, &event.IsInternal, &event.IsPrivate, &event.GroupDisplayOrder, &event.DefaultPicture, &event.IdGroups)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot get event",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
			"event": event,
		})
	}
}

func GetEventsByGroupID(tokenAPI string) func(c *gin.Context) {
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

		rows, err := db.Query("SELECT * FROM EVENTS WHERE Id_GROUPS = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot get events",
			})
			return
		}
		defer rows.Close()

		var events []Event

		for rows.Next() {
			var event Event
			err = rows.Scan(&event.IdEvent, &event.Name, &event.Type, &event.EndTime, &event.StartTime, &event.IsInternal, &event.IsPrivate, &event.GroupDisplayOrder, &event.DefaultPicture, &event.IdGroups)
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "cannot get events",
				})
				return
			}
			events = append(events, event)
		}

		c.JSON(200, gin.H{
			"error": false,
			"events": events,
		})
	}
}

func PostEvent(tokenAPI string) func(c *gin.Context) {
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

		var event Event
		c.BindJSON(&event)

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		if !utils.IsSafeString(event.Name) || !utils.IsSafeString(event.Type) { 
			c.JSON(400, gin.H{
				"error": true,
				"message": "name or type can't contain sql injection",
			})
			return
		}

		if event.Name == "" || event.Type == "" || event.EndTime == "" || event.StartTime == "" {
			c.JSON(400, gin.H{
				"error": true,
				"message": "name or type or endtime or starttime can't be empty",
			})
			return
		}

		var idManager int

		err = db.QueryRow("SELECT Id_MANAGERS FROM MANAGERS WHERE Id_USERS = '" + id + "'").Scan(&idManager)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot get manager",
			})
			return
		}

		err = db.QueryRow("SELECT Id_GROUPS FROM GROUPS WHERE Id_GROUPS = 1").Scan(&idManager)
		if err != nil {
			_, err := db.Exec("INSERT INTO GROUPS (name) VALUES (?)", "default")
			fmt.Println(err)
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "cannot insert group",
				})
				return
			}
		}

		result, err := db.Exec("INSERT INTO EVENTS (Name, Type, EndTime, StartTime, isInternal, isPrivate, group_display_order, DefaultPicture, Id_GROUPS) VALUES (?, ?, ?, ?, ?, ?, ?, DEFAULT, ?)", event.Name, event.Type, event.EndTime, event.StartTime, event.IsInternal, event.IsPrivate, 0, 1)
		fmt.Println(err)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot insert event",
			})
			return
		}

		lastId, err := result.LastInsertId()
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot get last id",
			})
			return
		}

		result, err = db.Exec("INSERT INTO ORGANIZES (Id_MANAGERS, Id_EVENTS) VALUES (?, ?)", idManager, lastId)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot insert organizes",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
			"message": "event inserted",
		})
	}
}

func AddEventToAGroup(tokenAPI string) func(c *gin.Context) {
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

		type Group struct {
			Name string `json:"name"`
			GroupDisplayOrder int `json:"group_display_order"`
		}

		var group Group
		c.BindJSON(&group)

		if !utils.IsSafeString(group.Name) {
			c.JSON(400, gin.H{
				"error": true,
				"message": "name can't contain sql injection",
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

		var idGroup int

		err = db.QueryRow("SELECT Id_GROUPS FROM GROUPS WHERE name = '" + group.Name + "'").Scan(&idGroup)
		if err != nil {
			result, err := db.Exec("INSERT INTO GROUPS (name) VALUES (?)", group.Name)
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "cannot insert group",
				})
				return
			}

			lastId, err := result.LastInsertId()
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "cannot get last id",
				})
				return
			}

			result, err = db.Exec("UPDATE EVENTS SET Id_GROUPS = ?, group_display_order = ? WHERE Id_EVENTS = ?", lastId, group.GroupDisplayOrder, id)
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "cannot update event",
				})
				return
			}

			c.JSON(200, gin.H{
				"error": false,
				"message": "event added to group",
			})
			return
		}

		var idEvent int

		err = db.QueryRow("SELECT Id_EVENTS FROM EVENTS WHERE Id_GROUPS = '" + strconv.Itoa(idGroup) + "' AND group_display_order = '" + strconv.Itoa(group.GroupDisplayOrder) + "'").Scan(&idEvent)
		if err != nil {
			_, err := db.Exec("UPDATE EVENTS SET Id_GROUPS = ?, group_display_order = ? WHERE Id_EVENTS = ?", idGroup, group.GroupDisplayOrder, id)
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "cannot update event",
				})
				return
			}

			c.JSON(200, gin.H{
				"error": false,
				"message": "event added to group",
			})
			return
		}

		c.JSON(500, gin.H{
			"error": false,
			"message": "an event is already in this group with this display order",
		})
	}
}