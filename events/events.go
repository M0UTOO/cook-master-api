package events

import (
	"cook-master-api/token"
	"cook-master-api/utils"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Event struct {
	IdEvent           int    `json:"idevent"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	Type              string `json:"type"`
	EndTime           string `json:"endtime"`
	IsClosed          bool   `json:"isclosed"`
	StartTime         string `json:"starttime"`
	IsInternal        bool   `json:"isinternal"`
	IsPrivate         bool   `json:"isprivate"`
	GroupDisplayOrder int    `json:"groupdisplayorder"`
	DefaultPicture    string `json:"defaultpicture"`
	IdEventGroups          int    `json:"ideventgroups"`
}

type EventGroup struct {
	IdEvent           int    `json:"idevent"`
	Name              string `json:"name"`
}

func GetEvents(tokenAPI string) func(c *gin.Context) {
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

		rows, err := db.Query("SELECT * FROM EVENTS")
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot get events",
			})
			return
		}
		defer rows.Close()

		var events []Event

		for rows.Next() {
			var event Event
			err = rows.Scan(&event.IdEvent, &event.Name, &event.Description, &event.Type, &event.EndTime, &event.IsClosed, &event.StartTime, &event.IsInternal, &event.IsPrivate, &event.GroupDisplayOrder, &event.DefaultPicture, &event.IdEventGroups)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error":   true,
					"message": "err on scan rows",
				})
				return
			}
			events = append(events, event)
		}

		c.JSON(200, events)
	}
}

func GetGroupEvents(tokenAPI string) func(c *gin.Context) {
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

		rows, err := db.Query("SELECT * FROM EVENTS_GROUPS")
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot get events groups",
			})
			return
		}
		defer rows.Close()

		var events []EventGroup

		for rows.Next() {
			var event EventGroup
			err = rows.Scan(&event.IdEvent, &event.Name)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error":   true,
					"message": "err on scan rows",
				})
				return
			}
			events = append(events, event)
		}

		c.JSON(200, events)
	}
}

func GetEventByID(tokenAPI string) func(c *gin.Context) {
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

		var event Event

		err = db.QueryRow("SELECT * FROM EVENTS WHERE Id_EVENTS = ?", id).Scan(&event.IdEvent, &event.Name, &event.Description, &event.Type, &event.EndTime, &event.IsClosed, &event.StartTime, &event.IsInternal, &event.IsPrivate, &event.GroupDisplayOrder, &event.DefaultPicture, &event.IdEventGroups)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot get event",
			})
			return
		}

		c.JSON(200, event)
	}
}

func GetEventsByGroupID(tokenAPI string) func(c *gin.Context) {
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

		rows, err := db.Query("SELECT * FROM EVENTS WHERE Id_EVENTS_GROUPS = ?", id)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot get events",
			})
			return
		}
		defer rows.Close()

		var events []Event

		for rows.Next() {
			var event Event
			err = rows.Scan(&event.IdEvent, &event.Name, &event.Description, &event.Type, &event.EndTime, &event.IsClosed, &event.StartTime, &event.IsInternal, &event.IsPrivate, &event.GroupDisplayOrder, &event.DefaultPicture, &event.IdEventGroups)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error":   true,
					"message": "cannot get events",
				})
				return
			}
			events = append(events, event)
		}

		c.JSON(200, events)
	}
}

func PostEvent(tokenAPI string) func(c *gin.Context) {
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

		var event Event
		c.BindJSON(&event)

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot connect to bdd",
			})
			return
		}
		defer db.Close()

		if !utils.IsSafeString(event.Name) || !utils.IsSafeString(event.Type) || !utils.IsSafeString(event.Description) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "name or type or description can't contain sql injection",
			})
			return
		}

		if event.Name == "" || event.Type == "" || event.EndTime == "" || event.StartTime == "" || event.Description == ""{
			c.JSON(400, gin.H{
				"error":   true,
				"message": "name or type or endtime or starttime or description can't be empty",
			})
			return
		}

		var idManager int

		err = db.QueryRow("SELECT Id_MANAGERS FROM MANAGERS WHERE Id_USERS = '" + id + "'").Scan(&idManager)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot get manager",
			})
			return
		}

		err = db.QueryRow("SELECT Id_EVENTS_GROUPS FROM EVENTS_GROUPS WHERE Id_EVENTS_GROUPS = 1").Scan(&idManager)
		if err != nil {
			_, err := db.Exec("INSERT INTO EVENTS_GROUPS (name) VALUES (?)", "default")
			fmt.Println(err)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "cannot insert group",
				})
				return
			}
		}

		result, err := db.Exec("INSERT INTO EVENTS (Name, Description, Type, EndTime, StartTime, isInternal, isPrivate, group_display_order, DefaultPicture, Id_EVENTS_GROUPS) VALUES (?, ?, ?, ?, ?, ?, ?, ?, DEFAULT, ?)", event.Name, event.Description, event.Type, event.EndTime, event.StartTime, event.IsInternal, event.IsPrivate, 0, 1)
		fmt.Println(err)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot insert event",
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

		result, err = db.Exec("INSERT INTO ORGANIZES (Id_MANAGERS, Id_EVENTS) VALUES (?, ?)", idManager, lastId)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot insert organizes",
			})
			return
		}

		c.JSON(200, gin.H{
			"error":   false,
			"id":      lastId,
			"message": "event inserted",
		})
	}
}

func AddEventToAGroup(tokenAPI string) func(c *gin.Context) {
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

		fmt.Println(group)

		if group.Name != "" {
			if !utils.IsSafeString(group.Name) {
				c.JSON(400, gin.H{
					"error":   true,
					"message": "name can't contain sql injection",
				})
				return
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

		err = db.QueryRow("SELECT Id_EVENTS_GROUPS FROM EVENTS_GROUPS WHERE name = '" + group.Name + "'").Scan(&idGroup)
		if err != nil {
			result, err := db.Exec("INSERT INTO EVENTS_GROUPS (name) VALUES (?)", group.Name)
			if err != nil {
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

			result, err = db.Exec("UPDATE EVENTS SET Id_EVENTS_GROUPS = ?, group_display_order = ? WHERE Id_EVENTS = ?", lastId, group.GroupDisplayOrder, id)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "cannot update event",
				})
				return
			}

			c.JSON(200, gin.H{
				"error":   false,
				"message": "event added to group",
			})
			return
		}

		var idEvent int

		err = db.QueryRow("SELECT Id_EVENTS FROM EVENTS WHERE Id_EVENTS_GROUPS = '" + strconv.Itoa(idGroup) + "' AND group_display_order = '" + strconv.Itoa(group.GroupDisplayOrder) + "'").Scan(&idEvent)
		if err != nil {
			_, err := db.Exec("UPDATE EVENTS SET Id_EVENTS_GROUPS = ?, group_display_order = ? WHERE Id_EVENTS = ?", idGroup, group.GroupDisplayOrder, id)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "cannot update event",
				})
				return
			}

			c.JSON(200, gin.H{
				"error":   false,
				"message": "event added to group",
			})
			return
		}

		c.JSON(500, gin.H{
			"error":   false,
			"message": "an event is already in this group with this display order",
		})
	}
}

func DeleteEventFromAGroup(tokenAPI string) func(c *gin.Context) {
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

		_, err = db.Exec("UPDATE EVENTS SET Id_EVENTS_GROUPS = 1, group_display_order = 0 WHERE Id_EVENTS = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot update event",
			})
			return
		}

		c.JSON(200, gin.H{
			"error":   false,
			"message": "event deleted from group",
		})
	}
}

func UpdateEvent(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		type EventReq struct {
			Name           string `json:"name"`
			Description	string `json:"description"`
			Type           string `json:"type"`
			EndTime        string `json:"endtime"`
			IsClosed       int    `json:"isclosed"`
			StartTime      string `json:"starttime"`
			IsInternal     int    `json:"isinternal"`
			IsPrivate      int    `json:"isprivate"`
			DefaultPicture string `json:"defaultpicture"`
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

		var event EventReq

		event.IsClosed = -1
		event.IsInternal = -1
		event.IsPrivate = -1

		err = c.BindJSON(&event)
		if err != nil {
			fmt.Println(err)
			c.JSON(400, gin.H{
				"error":   true,
				"message": "cannot bind json",
			})
			return
		}

		var setClause []string

		if event.Name != "" {
			if !utils.IsSafeString(event.Name) {
				c.JSON(400, gin.H{
					"error":   true,
					"message": "name can't contain sql injection",
				})
				return
			}
			if len(event.Name) < 0 || len(event.Name) > 100 {
				c.JSON(400, gin.H{
					"error":   true,
					"message": "wrong name length",
				})
				return
			}
			setClause = append(setClause, "name = '"+event.Name+"'")
		}

		if event.Description != "" {
			if !utils.IsSafeString(event.Description) {
				c.JSON(400, gin.H{
					"error":   true,
					"message": "description can't contain sql injection",
				})
				return
			}
			if len(event.Description) < 0 {
				c.JSON(400, gin.H{
					"error":   true,
					"message": "wrong description length",
				})
				return
			}
			setClause = append(setClause, "description = '"+event.Description+"'")
		}

		if event.Type != "" {
			if !utils.IsSafeString(event.Type) {
				c.JSON(400, gin.H{
					"error":   true,
					"message": "type can't contain sql injection",
				})
				return
			}
			if len(event.Type) < 0 || len(event.Type) > 50 {
				c.JSON(400, gin.H{
					"error":   true,
					"message": "wrong type length",
				})
				return
			}
			setClause = append(setClause, "type = '"+event.Type+"'")
		}

		if event.EndTime != "" {
			if !utils.IsSafeString(event.EndTime) {
				c.JSON(400, gin.H{
					"error":   true,
					"message": "end time can't contain sql injection",
				})
				return
			}
			setClause = append(setClause, "endtime = '"+event.EndTime+"'")
		}

		if event.IsClosed == 0 {
			setClause = append(setClause, "isclosed = false")
		} else if event.IsClosed == 1 {
			setClause = append(setClause, "isclosed = true")
		}

		if event.StartTime != "" {
			if !utils.IsSafeString(event.StartTime) {
				c.JSON(400, gin.H{
					"error":   true,
					"message": "start time can't contain sql injection",
				})
				return
			}
			setClause = append(setClause, "starttime = '"+event.StartTime+"'")
		}

		if event.IsInternal == 0 {
			setClause = append(setClause, "isinternal = false")
		} else if event.IsInternal == 1 {
			setClause = append(setClause, "isinternal = true")
		}

		if event.IsPrivate == 0 {
			setClause = append(setClause, "isprivate = false")
		} else if event.IsPrivate == 1 {
			setClause = append(setClause, "isprivate = true")
		}

		if event.DefaultPicture != "" {
			if !utils.IsSafeString(event.DefaultPicture) {
				c.JSON(400, gin.H{
					"error":   true,
					"message": "default picture can't contain sql injection",
				})
				return
			}
			if len(event.DefaultPicture) < 0 || len(event.DefaultPicture) > 255 {
				c.JSON(400, gin.H{
					"error":   true,
					"message": "wrong default picture length",
				})
				return
			}
			setClause = append(setClause, "defaultPicture = '"+event.DefaultPicture+"'")
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

		var idevent string

		err = db.QueryRow("SELECT Id_EVENTS FROM EVENTS WHERE Id_EVENTS = '" + id + "'").Scan(&idevent)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "event not found",
			})
			return
		}

		_, err = db.Exec("UPDATE EVENTS SET "+strings.Join(setClause, ", ")+" WHERE Id_EVENTS = ?", id)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot update event",
			})
			return
		}

		c.JSON(200, gin.H{
			"error":   false,
			"message": "event updated",
		})
	}
}

func AddContractorToAnEvent(tokenAPI string) func(c *gin.Context) {
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

		idEvent := c.Param("idevent")
		if idEvent == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id event can't be empty",
			})
			return
		}

		if !utils.IsSafeString(idEvent) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id event can't contain sql injection",
			})
			return
		}

		iduser := c.Param("iduser")
		if iduser == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id contractor can't be empty",
			})
			return
		}

		if !utils.IsSafeString(iduser) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id contractor can't contain sql injection",
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

		var idevent string

		err = db.QueryRow("SELECT Id_EVENTS FROM EVENTS WHERE Id_EVENTS = '" + idEvent + "'").Scan(&idevent)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "event not found",
			})
			return
		}

		var idcontractor string

		err = db.QueryRow("SELECT Id_CONTRACTORS FROM CONTRACTORS WHERE Id_USERS = '" + iduser + "'").Scan(&idcontractor)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "contractor not found",
			})
			return
		}

		var idanimates string

		err = db.QueryRow("SELECT Id_EVENTS FROM ANIMATES WHERE Id_CONTRACTORS = '" + idcontractor + "' AND Id_EVENTS = '" + idevent + "'").Scan(&idanimates)
		if err == nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "contractor already added to event",
			})
			return
		}

		_, err = db.Exec("INSERT INTO ANIMATES (Id_CONTRACTORS, Id_EVENTS) VALUES (?, ?)", idcontractor, idevent)
		fmt.Println(err)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot insert animates",
			})
			return
		}

		c.JSON(200, gin.H{
			"error":   false,
			"message": "contractor added to event",
		})
	}
}

func DeleteContractorFromAnEvent(tokenAPI string) func(c *gin.Context) {
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

		idEvent := c.Param("idevent")
		if idEvent == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id event can't be empty",
			})
			return
		}

		if !utils.IsSafeString(idEvent) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id event can't contain sql injection",
			})
			return
		}

		iduser := c.Param("iduser")
		if iduser == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id contractor can't be empty",
			})
			return
		}

		if !utils.IsSafeString(iduser) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id contractor can't contain sql injection",
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

		var idevent string

		err = db.QueryRow("SELECT Id_EVENTS FROM EVENTS WHERE Id_EVENTS = '" + idEvent + "'").Scan(&idevent)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "event not found",
			})
			return
		}

		var idcontractor string

		err = db.QueryRow("SELECT Id_CONTRACTORS FROM CONTRACTORS WHERE Id_USERS = '" + iduser + "'").Scan(&idcontractor)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "contractor not found",
			})
			return
		}

		var idanimates string

		err = db.QueryRow("SELECT Id_EVENTS FROM ANIMATES WHERE Id_CONTRACTORS = '" + idcontractor + "' AND Id_EVENTS = '" + idevent + "'").Scan(&idanimates)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "contractor not added to event",
			})
			return
		}

		_, err = db.Exec("DELETE FROM ANIMATES WHERE Id_CONTRACTORS = '" + idcontractor + "' AND Id_EVENTS = '" + idevent + "'")
		fmt.Println(err)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot delete animates",
			})
			return
		}

		c.JSON(200, gin.H{
			"error":   false,
			"message": "contractor removed from event",
		})
	}
}

func GetManagersByEventID(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		type ManagerReq struct {
			IdUser string `json:"iduser"`
			Email  string `json:"email"`
			IdManager string `json:"idmanager"`
		}

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

		idevent := c.Param("idevent")
		if idevent == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id event can't be empty",
			})
			return
		}

		if !utils.IsSafeString(idevent) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id event can't contain sql injection",
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

		var idmanager string

		err = db.QueryRow("SELECT Id_MANAGERS FROM ORGANIZES WHERE Id_EVENTS = '" + idevent + "'").Scan(&idmanager)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "event not found",
			})
			return
		}

		var iduser string

		err = db.QueryRow("SELECT Id_USERS FROM MANAGERS WHERE Id_MANAGERS = '" + idmanager + "'").Scan(&iduser)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "event not found",
			})
			return
		}

		var managers []ManagerReq
		var manager ManagerReq

		manager.IdManager = idmanager

		err = db.QueryRow("SELECT Id_USERS, email FROM USERS WHERE Id_USERS = '" + iduser + "'").Scan(&manager.IdUser, &manager.Email)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "event not found",
			})
			return
		}

		managers = append(managers, manager)

		c.JSON(200, managers)
	}
}

func GetGroupsByEventID(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		type GroupReq struct {
			IdGroup string `json:"idgroup"`
			Name 		 string `json:"name"`
		}

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

		idevent := c.Param("idevent")
		if idevent == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id event can't be empty",
			})
			return
		}

		if !utils.IsSafeString(idevent) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id event can't contain sql injection",
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

		var idgroup string

		err = db.QueryRow("SELECT Id_EVENTS_GROUPS FROM EVENTS WHERE Id_EVENTS = '" + idevent + "'").Scan(&idgroup)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "event not found",
			})
			return
		}

		var groups []GroupReq
		var group GroupReq

		err = db.QueryRow("SELECT Id_EVENTS_GROUPS, Name FROM EVENTS_GROUPS WHERE Id_EVENTS_GROUPS = '" + idgroup + "'").Scan(&group.IdGroup, &group.Name)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cookingspace not found",
			})
			return
		}

		groups = append(groups, group)

		c.JSON(200, groups)
	}
}

func GetClientsByEventID(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		type ClientReq struct {
			IdUser string `json:"iduser"`
			Email  string `json:"email"`
			FirstName  string `json:"firstname"`
			LastName string `json:"lastname"`
			IdClient string `json:"idclient"`
			IsPresent int `json:"ispresent"`
		}

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

		idevent := c.Param("idevent")
		if idevent == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id event can't be empty",
			})
			return
		}

		if !utils.IsSafeString(idevent) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id event can't contain sql injection",
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

		var clients []ClientReq
		var client ClientReq

		rows, err := db.Query("SELECT Id_CLIENTS, IsPresent FROM PARTICIPATES WHERE Id_EVENTS = '" + idevent + "'")
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "event not found",
			})
			return
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&client.IdClient, &client.IsPresent)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "error while scanning",
				})
				return
			}

			err = db.QueryRow("SELECT Id_USERS FROM CLIENTS WHERE Id_CLIENTS = '" + client.IdClient + "'").Scan(&client.IdUser)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "client not found",
				})
				return
			}

			err = db.QueryRow("SELECT email, firstname, lastname FROM USERS WHERE Id_USERS = '" + client.IdUser + "'").Scan(&client.Email, &client.FirstName, &client.LastName)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "user not found",
				})
				return
			}

			clients = append(clients, client)
		}

		c.JSON(200, clients)
	}
}


func GetCookingSpacesByEventID(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		type CookingSpaceReq struct {
			IdCookingSpace string `json:"idcookingspace"`
			Name 		 string `json:"name"`
			IdPremise string `json:"idpremise"`
		}

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

		idevent := c.Param("idevent")
		if idevent == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id event can't be empty",
			})
			return
		}

		if !utils.IsSafeString(idevent) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id event can't contain sql injection",
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

		err = db.QueryRow("SELECT Id_COOKING_SPACES FROM IS_HOSTED WHERE Id_EVENTS = '" + idevent + "'").Scan(&idcookingspace)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "event not found",
			})
			return
		}

		var cookingspaces []CookingSpaceReq
		var cookingspace CookingSpaceReq

		err = db.QueryRow("SELECT Name, Id_PREMISES FROM COOKING_SPACES WHERE Id_COOKING_SPACES = '" + idcookingspace + "'").Scan(&cookingspace.Name, &cookingspace.IdPremise)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cookingspace not found",
			})
			return
		}

		cookingspace.IdCookingSpace = idcookingspace

		cookingspaces = append(cookingspaces, cookingspace)

		c.JSON(200, cookingspaces)
	}
}


func GetContractorsByEventID(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header["Token"]

		type ContractorReq struct {
			IdContractor string `json:"idcontractor"`
			Email        string `json:"email"`
			IdUser	   string `json:"iduser"`
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

		idEvent := c.Param("idevent")
		if idEvent == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id event can't be empty",
			})
			return
		}

		if !utils.IsSafeString(idEvent) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id event can't contain sql injection",
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

		rows, err := db.Query("SELECT Id_CONTRACTORS FROM ANIMATES WHERE Id_EVENTS = '" + idEvent + "'")
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "contractors not found",
			})
			return
		}
		defer rows.Close()

		var contractors []ContractorReq	

		for rows.Next() {
			var contractor ContractorReq
			var idcontractor string
			err = rows.Scan(&idcontractor)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "cannot scan contractors",
				})
				return
			}

			var iduser string

			err = db.QueryRow("SELECT Id_USERS FROM CONTRACTORS WHERE Id_CONTRACTORS = '" + idcontractor + "'").Scan(&iduser)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "event not found",
				})
				return
			}

			err = db.QueryRow("SELECT Id_USERS, email FROM USERS WHERE Id_USERS = '" + iduser + "'").Scan(&contractor.IdUser, &contractor.Email)
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "event not found",
				})
				return
			}

			contractor.IdContractor = idcontractor

			contractors = append(contractors, contractor)
		}

		c.JSON(200, contractors)
	}
}

func AddClientToAnEvent(tokenAPI string) func(c *gin.Context) {
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

		idEvent := c.Param("idevent")
		if idEvent == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id event can't be empty",
			})
			return
		}

		if !utils.IsSafeString(idEvent) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id event can't contain sql injection",
			})
			return
		}

		iduser := c.Param("iduser")
		if iduser == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id client can't be empty",
			})
			return
		}

		if !utils.IsSafeString(iduser) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id client can't contain sql injection",
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

		var idevent string

		err = db.QueryRow("SELECT Id_EVENTS FROM EVENTS WHERE Id_EVENTS = '" + idEvent + "'").Scan(&idevent)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "event not found",
			})
			return
		}

		var idclient string

		err = db.QueryRow("SELECT Id_CLIENTS FROM CLIENTS WHERE Id_USERS = '" + iduser + "'").Scan(&idclient)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "client not found",
			})
			return
		}

		var idclients string

		err = db.QueryRow("SELECT Id_EVENTS FROM PARTICIPATES WHERE Id_CLIENTS = '" + idclient + "' AND Id_EVENTS = '" + idevent + "'").Scan(&idclients)
		fmt.Println(idclient, idevent, idclients)
		if err == nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "client already added to event",
			})
			return
		}

		_, err = db.Exec("INSERT INTO PARTICIPATES (Id_CLIENTS, Id_EVENTS, isPresent) VALUES ('" + idclient + "', '" + idevent + "', DEFAULT)")
		fmt.Println(err)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot add participates",
			})
			return
		}

		c.JSON(200, gin.H{
			"error":   false,
			"message": "client added to event",
		})
	}
}

func DeleteClientFromAnEvent(tokenAPI string) func(c *gin.Context) {
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

		idEvent := c.Param("idevent")
		if idEvent == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id event can't be empty",
			})
			return
		}

		if !utils.IsSafeString(idEvent) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id event can't contain sql injection",
			})
			return
		}

		iduser := c.Param("iduser")
		if iduser == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id client can't be empty",
			})
			return
		}

		if !utils.IsSafeString(iduser) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id client can't contain sql injection",
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

		var idevent string

		err = db.QueryRow("SELECT Id_EVENTS FROM EVENTS WHERE Id_EVENTS = '" + idEvent + "'").Scan(&idevent)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "event not found",
			})
			return
		}

		var idclient string

		err = db.QueryRow("SELECT Id_CLIENTS FROM CLIENTS WHERE Id_USERS = '" + iduser + "'").Scan(&idclient)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "client not found",
			})
			return
		}

		var idclients string

		err = db.QueryRow("SELECT Id_EVENTS FROM PARTICIPATES WHERE Id_CLIENTS = '" + idclient + "' AND Id_EVENTS = '" + idevent + "'").Scan(&idclients)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "client not added to event",
			})
			return
		}

		_, err = db.Exec("DELETE FROM PARTICIPATES WHERE Id_CLIENTS = '" + idclient + "' AND Id_EVENTS = '" + idevent + "'")
		fmt.Println(err)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot remove participates",
			})
			return
		}

		c.JSON(200, gin.H{
			"error":   false,
			"message": "client removed from event",
		})
	}
}

func ValidateClientPresence(tokenAPI string) func(c *gin.Context) {
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

		idEvent := c.Param("idevent")
		if idEvent == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id event can't be empty",
			})
			return
		}

		if !utils.IsSafeString(idEvent) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id event can't contain sql injection",
			})
			return
		}

		iduser := c.Param("iduser")
		if iduser == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id client can't be empty",
			})
			return
		}

		if !utils.IsSafeString(iduser) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id client can't contain sql injection",
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

		var idevent string

		err = db.QueryRow("SELECT Id_EVENTS FROM EVENTS WHERE Id_EVENTS = '" + idEvent + "'").Scan(&idevent)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "event not found",
			})
			return
		}

		var idclient string

		err = db.QueryRow("SELECT Id_CLIENTS FROM CLIENTS WHERE Id_USERS = '" + iduser + "'").Scan(&idclient)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "client not found",
			})
			return
		}

		var ispresent bool

		err = db.QueryRow("SELECT IsPresent FROM PARTICIPATES WHERE Id_CLIENTS = '" + idclient + "' AND Id_EVENTS = '" + idevent + "'").Scan(&ispresent)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "client not added to event",
			})
			return
		}

		if ispresent {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "client already validated",
			})
			return
		}

		_, err = db.Exec("UPDATE PARTICIPATES SET IsPresent = ? WHERE Id_CLIENTS = '"+idclient+"' AND Id_EVENTS = '"+idevent+"'", true)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot update participates",
			})
			return
		}

		var idgroup string

		err = db.QueryRow("SELECT Id_EVENTS_GROUPS FROM EVENTS WHERE Id_EVENTS = '" + idEvent + "'").Scan(&idgroup)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "client not added to event",
			})
			return
		}

		if idgroup != "1" && idgroup != "0"{

			var isformation int

			err = db.QueryRow("SELECT Id_CLIENTS FROM FORMATIONS WHERE Id_CLIENTS = '" + idclient + "' AND Id_EVENTS_GROUPS = '" + idgroup + "'").Scan(&isformation)
			if err != nil {
				fmt.Println(err)
				_, err = db.Exec("INSERT INTO FORMATIONS (Id_CLIENTS, Id_EVENTS_GROUPS, counter) VALUES ('" + idclient + "', '" + idgroup + "', DEFAULT)")
				fmt.Println(err)
				if err != nil {
					c.JSON(500, gin.H{
						"error":   true,
						"message": "cannot add formations",
					})
					return
				}

			} else {
				_, err = db.Exec("UPDATE FORMATIONS SET counter = counter + 1 WHERE Id_CLIENTS = '"+idclient+"' AND Id_EVENTS_GROUPS = '"+idgroup+"'")
				if err != nil {
					c.JSON(500, gin.H{
						"error":   true,
						"message": "cannot update formations",
					})
					return
				}
			}
		}
			
		c.JSON(200, gin.H{
			"error":   false,
			"message": "client presence validated",
		})
	}
}

func GetAllFormationsByUserID(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header["Token"]

		if len(tokenHeader) == 0 {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "missing token",
			})

			return
		}

		err := token.CheckAPIToken(tokenAPI, tokenHeader[0], c)
		if err != nil {
			c.JSON(498, gin.H{
				"error":   true,
				"message": "wrong token",
			})

			return
		}

		iduser := c.Param("iduser")
		if iduser == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id user can't be empty",
			})

			return
		}

		if !utils.IsSafeString(iduser) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id user can't contain sql injection",
			})

			return
		}

		db, err := sql.Open("mysql", token.DbLogins)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot connect to bdd",
			})

			return
		}
		defer db.Close()

		var idclient string

		err = db.QueryRow("SELECT Id_CLIENTS FROM CLIENTS WHERE Id_USERS = '" + iduser + "'").Scan(&idclient)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "client not found",
			})

			return
		}

		rows, err := db.Query("SELECT Id_EVENTS_GROUPS, counter FROM FORMATIONS WHERE Id_CLIENTS = '" + idclient + "'")
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "formations not found",
			})

			return
		}
		defer rows.Close()

		var formations []string

		for rows.Next() {
			var idgroup string
			var counter int
			err = rows.Scan(&idgroup, &counter)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error":   true,
					"message": "cannot scan formations",
				})

				return
			}

			var name string

			err = db.QueryRow("SELECT Name FROM EVENTS_GROUPS WHERE Id_EVENTS_GROUPS = '" + idgroup + "'").Scan(&name)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error":   true,
					"message": "event not found",
				})

				return
			}

			rows, err := db.Query("SELECT Id_EVENTS FROM EVENTS WHERE Id_EVENTS_GROUPS = '" + idgroup + "'")
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error":   true,
					"message": "event not found",
				})

				return
			}
			defer rows.Close()

			var counterEvents int

			for rows.Next() {
				counterEvents += 1
			}

			formations = append(formations, name + " : " + strconv.Itoa(counter) + "/" + strconv.Itoa(counterEvents))
		}

		c.JSON(200, formations)
	}
}

func AddEventToAnCookingSpace(tokenAPI string) func(c *gin.Context) {
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

		idEvent := c.Param("idevent")
		if idEvent == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id event can't be empty",
			})
			return
		}

		if !utils.IsSafeString(idEvent) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id event can't contain sql injection",
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

		var idevent string

		err = db.QueryRow("SELECT Id_EVENTS FROM EVENTS WHERE Id_EVENTS = '" + idEvent + "'").Scan(&idevent)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "event not found",
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

		var idishosted string

		err = db.QueryRow("SELECT Id_EVENTS FROM IS_HOSTED WHERE Id_COOKING_SPACES = '" + idcookingspace + "' AND Id_EVENTS = '" + idevent + "'").Scan(&idishosted)
		if err == nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "cookingspace already added to event",
			})
			return
		}

		_, err = db.Exec("INSERT INTO IS_HOSTED (Id_COOKING_SPACES, Id_EVENTS) VALUES (?, ?)", idcookingspace, idevent)
		fmt.Println(err)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot insert ishosted",
			})
			return
		}

		c.JSON(200, gin.H{
			"error":   false,
			"message": "cookingspace added to event",
		})
	}
}

func DeleteEventToAnCookingSpace(tokenAPI string) func(c *gin.Context) {
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

		idEvent := c.Param("idevent")
		if idEvent == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id event can't be empty",
			})
			return
		}

		if !utils.IsSafeString(idEvent) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id event can't contain sql injection",
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

		var idevent string

		err = db.QueryRow("SELECT Id_EVENTS FROM EVENTS WHERE Id_EVENTS = '" + idEvent + "'").Scan(&idevent)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "event not found",
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

		var idishosted string

		err = db.QueryRow("SELECT Id_EVENTS FROM IS_HOSTED WHERE Id_COOKING_SPACES = '" + idcookingspace + "' AND Id_EVENTS = '" + idevent + "'").Scan(&idishosted)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "event not added to the cookingspace",
			})
			return
		}

		_, err = db.Exec("DELETE FROM IS_HOSTED WHERE Id_EVENTS = '" + idEvent + "' AND Id_COOKING_SPACES = '" + idcooking + "'")
		fmt.Println(err)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot delete ishosted",
			})
			return
		}

		c.JSON(200, gin.H{
			"error":   false,
			"message": "cookingspace delete from event",
		})
	}
}

func GetEventsByMonth(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {

		type EventReq struct {
			Month  string `json:"month"`
			Counter int `json:"counter"`
		}

		tokenHeader := c.Request.Header["Token"]

		if len(tokenHeader) == 0 {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "missing token",
			})

			return
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
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot connect to bdd",
			})

			return
		}
		defer db.Close()

		rows, err := db.Query("SELECT DATE_FORMAT(endTime, '%m/%Y') AS month, COUNT(Id_EVENTS) AS counter FROM EVENTS GROUP BY DATE_FORMAT(endTime, '%m/%Y') ORDER BY endTime")
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "events not found",
			})

			return
		}
		defer rows.Close()

		var events []EventReq

		for rows.Next() {
			var event EventReq
			err = rows.Scan(&event.Month, &event.Counter)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error":   true,
					"message": "cannot scan events",
				})

				return
			}

			events = append(events, event)
		}

		c.JSON(200, events)
	}
}

func GetEventsByType(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header["Token"]

		type EventType struct {
			Type string `json:"type"`
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

		rows, err := db.Query("SELECT EVENTS.type AS name, COUNT(*) AS count FROM EVENTS GROUP BY EVENTS.type")
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "types not found",
			})
			return
		}

		var eventtypes []EventType

		for rows.Next() {
			var eventtype EventType
			err = rows.Scan(&eventtype.Type, &eventtype.Count)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error":   true,
					"message": "type not found",
				})
				return
			}
			eventtypes = append(eventtypes, eventtype)
		}

		c.JSON(200, eventtypes)
		return
	}
}

func GetEventsByDayOfTheWeek(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header["Token"]

		type EventWeek struct {
			WeekDay string `json:"weekday"`
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

		rows, err := db.Query("SELECT DAYNAME(endTime) AS dayOfWeek, COUNT(*) AS eventCount FROM EVENTS WHERE endTime IS NOT NULL GROUP BY dayOfWeek ORDER BY FIELD(dayOfWeek, 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday');")
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "days not found",
			})
			return
		}

		var eventweeks []EventWeek

		for rows.Next() {
			var eventweek EventWeek
			err = rows.Scan(&eventweek.WeekDay, &eventweek.Count)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error":   true,
					"message": "type not found",
				})
				return
			}
			eventweeks = append(eventweeks, eventweek)
		}

		c.JSON(200, eventweeks)
		return
	}
}

func GetEventsByMonthInAYear(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header["Token"]

		type EventMonth struct {
			Month string `json:"month"`
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

		rows, err := db.Query("SELECT DATE_FORMAT(endTime, '%M') AS month, COUNT(*) AS eventCount FROM EVENTS WHERE endTime IS NOT NULL GROUP BY month ORDER BY FIELD(month, 'January', 'February', 'March', 'April', 'May', 'June', 'July', 'August', 'September', 'October', 'November', 'December');")
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "month not found",
			})
			return
		}

		var eventmonths []EventMonth

		for rows.Next() {
			var eventmonth EventMonth
			err = rows.Scan(&eventmonth.Month, &eventmonth.Count)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error":   true,
					"message": "type not found",
				})
				return
			}
			eventmonths = append(eventmonths, eventmonth)
		}

		c.JSON(200, eventmonths)
		return
	}
}

func GetTop5Events(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header["Token"]

		type Top5Event struct {
			IdEvent int `json:"idevent"`
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

		rows, err := db.Query("SELECT e.Id_EVENTS, e.name, COUNT(*) AS participationCount FROM EVENTS e JOIN PARTICIPATES p ON e.Id_EVENTS = p.Id_EVENTS WHERE p.isPresent = TRUE GROUP BY e.Id_EVENTS, e.name ORDER BY participationCount DESC LIMIT 5;")
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
			err = rows.Scan(&top5event.IdEvent, &top5event.Name, &top5event.Count)
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

func GetFormationsDone(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header["Token"]

		type Formation struct {
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

		rows, err := db.Query("SELECT EG.name AS eventGroup, COUNT(F.Id_CLIENTS) AS formationCount FROM EVENTS_GROUPS EG LEFT JOIN FORMATIONS F ON EG.Id_EVENTS_GROUPS = F.Id_EVENTS_GROUPS LEFT JOIN EVENTS E ON EG.Id_EVENTS_GROUPS = E.Id_EVENTS_GROUPS GROUP BY EG.Id_EVENTS_GROUPS, EG.name HAVING COUNT(F.Id_CLIENTS) = COUNT(E.Id_EVENTS) ORDER BY EG.name;")
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "month not found",
			})
			return
		}

		var formations []Formation

		for rows.Next() {
			var formation Formation
			err = rows.Scan(&formation.Name, &formation.Count)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error":   true,
					"message": "formation not found",
				})
				return
			}
			formations = append(formations, formation)
		}

		c.JSON(200, formations)
		return
	}
}

func SearchForEvents(tokenAPI string) func(c *gin.Context) {
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

		rows, err := db.Query("SELECT * FROM EVENTS WHERE name LIKE '%" + search + "%'")
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "month not found",
			})
			return
		}

		var events []Event

		for rows.Next() {
			var event Event
			err = rows.Scan(&event.IdEvent, &event.Name, &event.Description, &event.Type, &event.EndTime, &event.IsClosed, &event.StartTime, &event.IsInternal, &event.IsPrivate, &event.GroupDisplayOrder, &event.DefaultPicture, &event.IdEventGroups)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error":   true,
					"message": "events not found",
				})
				return
			}
			events = append(events, event)
		}

		c.JSON(200, events)
		return
	}
}

func GetRateByEventID(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header["Token"]

		type Rate struct {
			Grade float64 `json:"grade"`
			Id string `json:"id"`
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

		var rate Rate

		err = db.QueryRow("SELECT AVG(grade) AS average_grade FROM EVENTS JOIN COMMENTS ON EVENTS.Id_EVENTS = COMMENTS.Id_EVENTS WHERE EVENTS.Id_EVENTS = '" + id + "'").Scan(&rate.Grade)
		if err != nil {
			rate.Grade = 0
		}

		rate.Id = id

		c.JSON(200, rate)
		return
	}
}

func GetEventComments(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header["Token"]

		type Comment struct {
			IdComment int `json:"idcomment"`
			IdUser int `json:"iduser"`
			Grade float64 `json:"grade"`
			Content string `json:"content"`
			Picture string `json:"picture"`
			Firstname string `json:"firstname"`
			Lastname string `json:"lastname"`
			ProfilePicture string `json:"profilepicture"`
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

		rows, err := db.Query("SELECT COMMENTS.Id_COMMENTS, COMMENTS.grade, COMMENTS.picture, COMMENTS.content, USERS.firstname, USERS.lastname, USERS.profilePicture, USERS.Id_USERS FROM COMMENTS JOIN CLIENTS ON CLIENTS.Id_CLIENTS = COMMENTS.Id_CLIENTS JOIN USERS ON USERS.Id_USERS = CLIENTS.Id_USERS WHERE COMMENTS.Id_EVENTS = '" + id + "'")
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "month not found",
			})
			return
		}

		var comments []Comment

		for rows.Next() {
			var comment Comment
			err = rows.Scan(&comment.IdComment, &comment.Grade, &comment.Picture, &comment.Content, &comment.Firstname, &comment.Lastname, &comment.ProfilePicture, &comment.IdUser)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error":   true,
					"message": "comment not found",
				})
				return
			}
			comments = append(comments, comment)
		}

		c.JSON(200, comments)
		return
	}
}

func GetComingEventByClientId(tokenAPI string) func (c *gin.Context) {
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

		var idClient string

		err = db.QueryRow("SELECT Id_CLIENTS FROM CLIENTS WHERE Id_USERS = '" + id + "'").Scan(&idClient)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "client not found",
			})
			return
		}

		rows, err := db.Query("SELECT e.Id_EVENTS, e.name, e.description, e.type, e.endTime, e.isClosed, e.startTime, e.isInternal, e.isPrivate, e.group_display_order, e.defaultPicture, e.Id_EVENTS_GROUPS FROM EVENTS e JOIN PARTICIPATES p ON e.Id_EVENTS = p.Id_EVENTS WHERE p.Id_CLIENTS = " + idClient + " AND e.endTime > NOW();")
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "events not found",
			})
			return
		}

		var events []Event

		for rows.Next() {
			var event Event
			err = rows.Scan(&event.IdEvent, &event.Name, &event.Description, &event.Type, &event.EndTime, &event.IsClosed, &event.StartTime, &event.IsInternal, &event.IsPrivate, &event.GroupDisplayOrder, &event.DefaultPicture, &event.IdEventGroups)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error":   true,
					"message": "events not found",
				})
				return
			}
			events = append(events, event)
		}

		c.JSON(200, events)
		return
	}
}

func GetPastEventByClientIdfunc(tokenAPI string) func (c *gin.Context) {
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

		var idClient string

		err = db.QueryRow("SELECT Id_CLIENTS FROM CLIENTS WHERE Id_USERS = '" + id + "'").Scan(&idClient)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "client not found",
			})
			return
		}

		rows, err := db.Query("SELECT e.Id_EVENTS, e.name, e.description, e.type, e.endTime, e.isClosed, e.startTime, e.isInternal, e.isPrivate, e.group_display_order, e.defaultPicture, e.Id_EVENTS_GROUPS FROM EVENTS e JOIN PARTICIPATES p ON e.Id_EVENTS = p.Id_EVENTS WHERE p.Id_CLIENTS = " + idClient + " AND e.endTime < NOW() AND p.IsPresent = 1;")
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "events not found",
			})
			return
		}

		var events []Event

		for rows.Next() {
			var event Event
			err = rows.Scan(&event.IdEvent, &event.Name, &event.Description, &event.Type, &event.EndTime, &event.IsClosed, &event.StartTime, &event.IsInternal, &event.IsPrivate, &event.GroupDisplayOrder, &event.DefaultPicture, &event.IdEventGroups)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error":   true,
					"message": "events not found",
				})
				return
			}
			events = append(events, event)
		}

		c.JSON(200, events)
		return
	}
}

func DeleteEvent(tokenAPI string) func(c *gin.Context) {
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

		id := c.Param("id")
		if id == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id event can't be empty",
			})
			return
		}

		if !utils.IsSafeString(id) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id event can't contain sql injection",
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

		var idevent string

		err = db.QueryRow("SELECT Id_EVENTS FROM EVENTS WHERE Id_EVENTS = '" + id + "'").Scan(&idevent)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "event not found",
			})
			return
		}

		var animate string

		err = db.QueryRow("SELECT Id_EVENTS FROM ANIMATES WHERE Id_EVENTS = '" + id + "'").Scan(&animate)
		if err == nil {
			_, err = db.Exec("DELETE FROM ANIMATES WHERE Id_EVENTS = '" + id + "'")
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "cannot delete animates",
				})
				return
			}
		}

		var ishosted string

		err = db.QueryRow("SELECT Id_EVENTS FROM IS_HOSTED WHERE Id_EVENTS = '" + id + "'").Scan(&ishosted)
		if err == nil {
			_, err = db.Exec("DELETE FROM IS_HOSTED WHERE Id_EVENTS = '" + id + "'")
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "cannot delete is_hosted",
				})
				return
			}
		}

		var organize string

		err = db.QueryRow("SELECT Id_EVENTS FROM ORGANIZES WHERE Id_EVENTS = '" + id + "'").Scan(&organize)
		if err == nil {
			_, err = db.Exec("DELETE FROM ORGANIZES WHERE Id_EVENTS = '" + id + "'")
			if err != nil {
				c.JSON(500, gin.H{
					"error":   true,
					"message": "cannot delete organizes",
				})
				return
			}
		}

		rows, err := db.Query("SELECT Id_EVENTS FROM PARTICIPATES WHERE Id_EVENTS = '" + id + "'")
		if err == nil {
			for rows.Next() {
				var participate string
				err = rows.Scan(&participate)
				if err != nil {
					c.JSON(500, gin.H{
						"error":   true,
						"message": "err on scan rows",
					})
					return
				}

				_, err = db.Exec("DELETE FROM PARTICIPATES WHERE Id_EVENTS = '" + id + "'")
				if err != nil {
					c.JSON(500, gin.H{
						"error":   true,
						"message": "cannot delete participate",
					})
					return
				}
			}
		}

		rows, err = db.Query("SELECT Id_EVENTS FROM COMMENTS WHERE Id_EVENTS = '" + id + "'")
		if err == nil {
			for rows.Next() {
				var comment string
				err = rows.Scan(&comment)
				if err != nil {
					c.JSON(500, gin.H{
						"error":   true,
						"message": "err on scan rows",
					})
					return
				}

				_, err = db.Exec("DELETE FROM COMMENTS WHERE Id_EVENTS = '" + id + "'")
				if err != nil {
					c.JSON(500, gin.H{
						"error":   true,
						"message": "cannot delete comment",
					})
					return
				}
			}
		}

		_, err = db.Exec("DELETE FROM EVENTS WHERE Id_EVENTS = '" + id + "'")
		fmt.Println(err)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot delete ishosted",
			})
			return
		}

		c.JSON(200, gin.H{
			"error":   false,
			"message": "event deleted",
		})
	}
}

func SearchForEventsGroups(tokenAPI string) func(c *gin.Context) {
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

		var eventGroup EventGroup

		err = db.QueryRow("SELECT * FROM EVENTS_GROUPS WHERE name = '" + search + "';").Scan(&eventGroup.IdEvent, &eventGroup.Name)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "groups not found",
			})
			return
		}

		c.JSON(200, eventGroup)
		return
	}
}

func GetEventsByUserId(tokenAPI string) func(c *gin.Context) {
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

		var idcontractor string

		err = db.QueryRow("SELECT Id_CONTRACTORS FROM CONTRACTORS WHERE Id_USERS = '" + id + "';").Scan(&idcontractor)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "contractor not found",
			})
			return
		}

		rows, err := db.Query("SELECT EVENTS.Id_EVENTS, EVENTS.Name, EVENTS.Description, EVENTS.Type, EVENTS.Endtime, EVENTS.IsClosed, EVENTS.Starttime, EVENTS.IsInternal, EVENTS.IsPrivate, EVENTS.Group_Display_Order, EVENTS.defaultPicture, EVENTS.Id_EVENTS_GROUPS FROM EVENTS JOIN ANIMATES ON EVENTS.Id_EVENTS = ANIMATES.Id_EVENTS WHERE Id_CONTRACTORS = '" + idcontractor + "';")
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "events not found",
			})
			return
		}

		var events []Event

		for rows.Next() {
			var event Event
			err = rows.Scan(&event.IdEvent, &event.Name, &event.Description, &event.Type, &event.EndTime, &event.IsClosed, &event.StartTime, &event.IsInternal, &event.IsPrivate, &event.GroupDisplayOrder, &event.DefaultPicture, &event.IdEventGroups)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error":   true,
					"message": "events not found",
				})
				return
			}
			events = append(events, event)
		}

		c.JSON(200, events)
		return
	}
}

func GetClientParticipationToEvent(tokenAPI string) func(c *gin.Context) {
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

		idclient := c.Param("idclient")
		if idclient == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "idclient can't be empty",
			})
			return
		}

		if !utils.IsSafeString(idclient) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "idclient can't contain sql injection",
			})
			return
		}

		idevent := c.Param("idevent")
		if idevent == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "idevent can't be empty",
			})
			return
		}

		if !utils.IsSafeString(idevent) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "idevent can't contain sql injection",
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

		var verifclient string

		err = db.QueryRow("SELECT Id_CLIENTS FROM CLIENTS WHERE Id_CLIENTS = '" + idclient + "';").Scan(&verifclient)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "client not found",
			})
			return
		}

		var verifevent string

		err = db.QueryRow("SELECT Id_EVENTS FROM EVENTS WHERE Id_EVENTS = '" + idevent + "';").Scan(&verifevent)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "event not found",
			})
			return
		}

		type IsPresent struct {
			IsPresent int `json:"ispresent"`
		}

		var isPresent IsPresent

		err = db.QueryRow("SELECT isPresent FROM PARTICIPATES WHERE Id_CLIENTS = '" + idclient + "' AND Id_EVENTS = '" + idevent + "';").Scan(&isPresent.IsPresent)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "participation not found",
			})
			return
		}

		c.JSON(200, isPresent)
		return
	}
}

func UnvalidateClientPresence(tokenAPI string) func(c *gin.Context) {
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

		idEvent := c.Param("idevent")
		if idEvent == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id event can't be empty",
			})
			return
		}

		if !utils.IsSafeString(idEvent) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id event can't contain sql injection",
			})
			return
		}

		iduser := c.Param("iduser")
		if iduser == "" {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id client can't be empty",
			})
			return
		}

		if !utils.IsSafeString(iduser) {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "id client can't contain sql injection",
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

		var idevent string

		err = db.QueryRow("SELECT Id_EVENTS FROM EVENTS WHERE Id_EVENTS = '" + idEvent + "'").Scan(&idevent)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "event not found",
			})
			return
		}

		var idclient string

		err = db.QueryRow("SELECT Id_CLIENTS FROM CLIENTS WHERE Id_USERS = '" + iduser + "'").Scan(&idclient)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "client not found",
			})
			return
		}

		var ispresent bool

		err = db.QueryRow("SELECT IsPresent FROM PARTICIPATES WHERE Id_CLIENTS = '" + idclient + "' AND Id_EVENTS = '" + idevent + "'").Scan(&ispresent)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "client not added to event",
			})
			return
		}

		if ispresent == false{
			c.JSON(400, gin.H{
				"error":   true,
				"message": "client not validated",
			})
			return
		}

		_, err = db.Exec("UPDATE PARTICIPATES SET IsPresent = ? WHERE Id_CLIENTS = '"+idclient+"' AND Id_EVENTS = '"+idevent+"'", false)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": "cannot update participates",
			})
			return
		}

		var idgroup string

		err = db.QueryRow("SELECT Id_EVENTS_GROUPS FROM EVENTS WHERE Id_EVENTS = '" + idEvent + "'").Scan(&idgroup)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   true,
				"message": "client not added to event",
			})
			return
		}

		if idgroup != "1" && idgroup != "0"{

			var isformation int

			err = db.QueryRow("SELECT Id_CLIENTS FROM FORMATIONS WHERE Id_CLIENTS = '" + idclient + "' AND Id_EVENTS_GROUPS = '" + idgroup + "'").Scan(&isformation)
			if err != nil {
				c.JSON(400, gin.H{
					"error":   true,
					"message": "client not added to event",
				})
				return

			} else {
				_, err = db.Exec("UPDATE FORMATIONS SET counter = counter - 1 WHERE Id_CLIENTS = '"+idclient+"' AND Id_EVENTS_GROUPS = '"+idgroup+"'")
				if err != nil {
					c.JSON(500, gin.H{
						"error":   true,
						"message": "cannot update formations",
					})
					return
				}
			}
		}
			
		c.JSON(200, gin.H{
			"error":   false,
			"message": "client presence deleted",
		})
	}
}

func GetFomationsForUser(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header["Token"]

		type FormationReq struct {
			IdEventGroup string `json:"ideventgroup"`
			NameEventGroup string `json:"nameeventgroup"`
			EventCount int `json:"eventcount"`
			Count int `json:"count"`
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

		var verifclient string

		err = db.QueryRow("SELECT Id_CLIENTS FROM CLIENTS WHERE Id_USERS = '" + id + "';").Scan(&verifclient)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "client not found",
			})
			return
		}

		rows, err := db.Query("SELECT F.Id_EVENTS_GROUPS, EG.name AS event_group_name, COUNT(E.Id_EVENTS) AS event_count, F.counter AS formations_counter FROM EVENTS_GROUPS EG JOIN EVENTS E ON E.Id_EVENTS_GROUPS = EG.Id_EVENTS_GROUPS JOIN FORMATIONS F ON F.Id_EVENTS_GROUPS = EG.Id_EVENTS_GROUPS WHERE F.Id_CLIENTS = '" + verifclient + "' GROUP BY EG.Id_EVENTS_GROUPS;")
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "formations not found",
			})
			return
		}

		var formations []FormationReq

		for rows.Next() {
			var formation FormationReq
			err = rows.Scan(&formation.IdEventGroup, &formation.NameEventGroup, &formation.EventCount, &formation.Count)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error":   true,
					"message": "formations not found",
				})
				return
			}
			formations = append(formations, formation)
		}

		c.JSON(200, formations)
		return
	}
}

func GetGroupByGroupId(tokenAPI string) func(c *gin.Context) {
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

		var eventGroup EventGroup

		err = db.QueryRow("SELECT * FROM EVENTS_GROUPS WHERE Id_EVENTS_GROUPS = '" + id + "';").Scan(&eventGroup.IdEvent, &eventGroup.Name)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "group not found",
			})
			return
		}

		c.JSON(200, eventGroup)
		return
	}
}
