package events

import (
	"fmt"
	"database/sql"
	"strconv"
	"strings"
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
	IsClosed bool `json:"isclosed"`
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
			err = rows.Scan(&event.IdEvent, &event.Name, &event.Type, &event.EndTime, &event.IsClosed, &event.StartTime, &event.IsInternal, &event.IsPrivate, &event.GroupDisplayOrder, &event.DefaultPicture, &event.IdGroups)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"error": true,
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

		err = db.QueryRow("SELECT * FROM EVENTS WHERE Id_EVENTS = ?", id).Scan(&event.IdEvent, &event.Name, &event.Type, &event.EndTime, &event.IsClosed, &event.StartTime, &event.IsInternal, &event.IsPrivate, &event.GroupDisplayOrder, &event.DefaultPicture, &event.IdGroups)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
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
			err = rows.Scan(&event.IdEvent, &event.Name, &event.Type, &event.EndTime, &event.IsClosed, &event.StartTime, &event.IsInternal, &event.IsPrivate, &event.GroupDisplayOrder, &event.DefaultPicture, &event.IdGroups)
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
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

func DeleteEventFromAGroup(tokenAPI string) func(c *gin.Context) {
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

		_, err = db.Exec("UPDATE EVENTS SET Id_GROUPS = 1, group_display_order = 0 WHERE Id_EVENTS = ?", id)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot update event",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
			"message": "event deleted from group",
		})
	}
}

func UpdateEvent(tokenAPI string) func(c *gin.Context) {
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

		err = c.BindJSON(&event)	
		if err != nil {
			c.JSON(400, gin.H{
				"error": true,
				"message": "cannot bind json",
			})
			return
		}

		var setClause []string

		if event.Name != "" {
			if !utils.IsSafeString(event.Name) {
				c.JSON(400, gin.H{
					"error": true,
					"message": "name can't contain sql injection",
				})
				return
			}
			if len(event.Name) < 0 || len(event.Name) > 100 {
				c.JSON(400, gin.H{
					"error": true,
					"message": "wrong name length",
				})
				return
			}
			setClause = append(setClause, "name = '"+event.Name+"'")
		}

		if event.Type != "" {
			if !utils.IsSafeString(event.Type) {
				c.JSON(400, gin.H{
					"error": true,
					"message": "type can't contain sql injection",
				})
				return
			}
			if len(event.Type) < 0 || len(event.Type) > 50 {
				c.JSON(400, gin.H{
					"error": true,
					"message": "wrong type length",
				})
				return
			}
			setClause = append(setClause, "type = '"+event.Type+"'")
		}

		if event.EndTime != "" {
			if !utils.IsSafeString(event.EndTime) {
				c.JSON(400, gin.H{
					"error": true,
					"message": "end time can't contain sql injection",
				})
				return
			}
			setClause = append(setClause, "endtime = '"+event.EndTime+"'")
		}

		if event.IsClosed == false || event.IsClosed == true {
			setClause = append(setClause, "isclosed = "+strconv.FormatBool(event.IsClosed))
		}

		if event.StartTime != "" {
			if !utils.IsSafeString(event.StartTime) {
				c.JSON(400, gin.H{
					"error": true,
					"message": "start time can't contain sql injection",
				})
				return
			}
			setClause = append(setClause, "starttime = '"+event.StartTime+"'")
		}

		if event.IsInternal == false || event.IsInternal == true {
			setClause = append(setClause, "isinternal = "+strconv.FormatBool(event.IsInternal))
		}

		if event.IsPrivate == false || event.IsPrivate == true {
			setClause = append(setClause, "isprivate = "+strconv.FormatBool(event.IsPrivate))
		}

		if event.DefaultPicture != "" {
			if !utils.IsSafeString(event.DefaultPicture) {
				c.JSON(400, gin.H{
					"error": true,
					"message": "default picture can't contain sql injection",
				})
				return
			}
			if len(event.DefaultPicture) < 0 || len(event.DefaultPicture) > 255 {
				c.JSON(400, gin.H{
					"error": true,
					"message": "wrong default picture length",
				})
				return
			}
			setClause = append(setClause, "default_picture = '"+event.DefaultPicture+"'")
		}

		if len(setClause) == 0 {
			c.JSON(400, gin.H{
				"error": true,
				"message": "no field to update",
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

		var idevent string
		
		err = db.QueryRow("SELECT Id_EVENTS FROM EVENTS WHERE Id_EVENTS = '" + id + "'").Scan(&idevent)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "event not found",
			})
			return
		}

		_, err = db.Exec("UPDATE EVENTS SET "+strings.Join(setClause, ", ")+" WHERE Id_EVENTS = ?", id)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot update event",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
			"message": "event updated",
		})
	}
}

func AddContractorToAnEvent(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header["Token"]

		if len(tokenHeader) == 0 {
			c.JSON(400, gin.H{
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

		idEvent := c.Param("idevent")
		if idEvent == "" {
			c.JSON(400, gin.H{
				"error": true,
				"message": "id event can't be empty",
			})
			return
		}

		if !utils.IsSafeString(idEvent) {
			c.JSON(400, gin.H{
				"error": true,
				"message": "id event can't contain sql injection",
			})
			return
		}

		iduser := c.Param("iduser")
		if iduser == "" {
			c.JSON(400, gin.H{
				"error": true,
				"message": "id contractor can't be empty",
			})
			return
		}

		if !utils.IsSafeString(iduser) {
			c.JSON(400, gin.H{
				"error": true,
				"message": "id contractor can't contain sql injection",
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

		var idevent string
		
		err = db.QueryRow("SELECT Id_EVENTS FROM EVENTS WHERE Id_EVENTS = '" + idEvent + "'").Scan(&idevent)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "event not found",
			})
			return
		}

		var idcontractor string
		
		err = db.QueryRow("SELECT Id_CONTRACTORS FROM CONTRACTORS WHERE Id_USERS = '" + iduser + "'").Scan(&idcontractor)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "contractor not found",
			})
			return
		}

		var idanimates string

		err = db.QueryRow("SELECT Id_EVENTS FROM ANIMATES WHERE Id_CONTRACTORS = '" + idcontractor + "' AND Id_EVENTS = '" + idevent + "'").Scan(&idanimates)
		if err == nil {
			c.JSON(400, gin.H{
				"error": true,
				"message": "contractor already added to event",
			})
			return
		}

		_, err = db.Exec("INSERT INTO ANIMATES (Id_CONTRACTORS, Id_EVENTS) VALUES (?, ?)", idcontractor, idevent)
		fmt.Println(err)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot insert animates",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
			"message": "contractor added to event",
		})
	}
}

func DeleteContractorFromAnEvent(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header["Token"]

		if len(tokenHeader) == 0 {
			c.JSON(400, gin.H{
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

		idEvent := c.Param("idevent")
		if idEvent == "" {
			c.JSON(400, gin.H{
				"error": true,
				"message": "id event can't be empty",
			})
			return
		}

		if !utils.IsSafeString(idEvent) {
			c.JSON(400, gin.H{
				"error": true,
				"message": "id event can't contain sql injection",
			})
			return
		}

		iduser := c.Param("iduser")
		if iduser == "" {
			c.JSON(400, gin.H{
				"error": true,
				"message": "id contractor can't be empty",
			})
			return
		}

		if !utils.IsSafeString(iduser) {
			c.JSON(400, gin.H{
				"error": true,
				"message": "id contractor can't contain sql injection",
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

		var idevent string
		
		err = db.QueryRow("SELECT Id_EVENTS FROM EVENTS WHERE Id_EVENTS = '" + idEvent + "'").Scan(&idevent)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "event not found",
			})
			return
		}

		var idcontractor string
		
		err = db.QueryRow("SELECT Id_CONTRACTORS FROM CONTRACTORS WHERE Id_USERS = '" + iduser + "'").Scan(&idcontractor)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "contractor not found",
			})
			return
		}

		var idanimates string

		err = db.QueryRow("SELECT Id_EVENTS FROM ANIMATES WHERE Id_CONTRACTORS = '" + idcontractor + "' AND Id_EVENTS = '" + idevent + "'").Scan(&idanimates)
		if err != nil {
			c.JSON(400, gin.H{
				"error": true,
				"message": "contractor not added to event",
			})
			return
		}

		_, err = db.Exec("DELETE FROM ANIMATES WHERE Id_CONTRACTORS = '" + idcontractor + "' AND Id_EVENTS = '" + idevent + "'")
		fmt.Println(err)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot delete animates",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
			"message": "contractor removed from event",
		})
	}
}

func AddClientToAnEvent(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header["Token"]

		if len(tokenHeader) == 0 {
			c.JSON(400, gin.H{
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

		idEvent := c.Param("idevent")
		if idEvent == "" {
			c.JSON(400, gin.H{
				"error": true,
				"message": "id event can't be empty",
			})
			return
		}

		if !utils.IsSafeString(idEvent) {
			c.JSON(400, gin.H{
				"error": true,
				"message": "id event can't contain sql injection",
			})
			return
		}

		iduser := c.Param("iduser")
		if iduser == "" {
			c.JSON(400, gin.H{
				"error": true,
				"message": "id client can't be empty",
			})
			return
		}

		if !utils.IsSafeString(iduser) {
			c.JSON(400, gin.H{
				"error": true,
				"message": "id client can't contain sql injection",
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

		var idevent string
		
		err = db.QueryRow("SELECT Id_EVENTS FROM EVENTS WHERE Id_EVENTS = '" + idEvent + "'").Scan(&idevent)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "event not found",
			})
			return
		}

		var idclient string
		
		err = db.QueryRow("SELECT Id_CLIENTS FROM CLIENTS WHERE Id_USERS = '" + iduser + "'").Scan(&idclient)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "client not found",
			})
			return
		}

		var idclients string

		err = db.QueryRow("SELECT Id_EVENTS FROM PARTICIPATES WHERE Id_CLIENTS = '" + idclient + "' AND Id_EVENTS = '" + idevent + "'").Scan(&idclients)
		if err == nil {
			c.JSON(400, gin.H{
				"error": true,
				"message": "client already added to event",
			})
			return
		}

		_, err = db.Exec("INSERT INTO PARTICIPATES (Id_CLIENTS, Id_EVENTS, isPresent) VALUES ('" + idclient + "', '" + idevent + "', DEFAULT)")
		fmt.Println(err)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot add participates",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
			"message": "client added to event",
		})
	}
}

func DeleteClientFromAnEvent(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header["Token"]

		if len(tokenHeader) == 0 {
			c.JSON(400, gin.H{
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

		idEvent := c.Param("idevent")
		if idEvent == "" {
			c.JSON(400, gin.H{
				"error": true,
				"message": "id event can't be empty",
			})
			return
		}

		if !utils.IsSafeString(idEvent) {
			c.JSON(400, gin.H{
				"error": true,
				"message": "id event can't contain sql injection",
			})
			return
		}

		iduser := c.Param("iduser")
		if iduser == "" {
			c.JSON(400, gin.H{
				"error": true,
				"message": "id client can't be empty",
			})
			return
		}

		if !utils.IsSafeString(iduser) {
			c.JSON(400, gin.H{
				"error": true,
				"message": "id client can't contain sql injection",
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

		var idevent string
		
		err = db.QueryRow("SELECT Id_EVENTS FROM EVENTS WHERE Id_EVENTS = '" + idEvent + "'").Scan(&idevent)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "event not found",
			})
			return
		}

		var idclient string
		
		err = db.QueryRow("SELECT Id_CLIENTS FROM CLIENTS WHERE Id_USERS = '" + iduser + "'").Scan(&idclient)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "client not found",
			})
			return
		}

		var idclients string

		err = db.QueryRow("SELECT Id_EVENTS FROM PARTICIPATES WHERE Id_CLIENTS = '" + idclient + "' AND Id_EVENTS = '" + idevent + "'").Scan(&idclients)
		if err != nil {
			c.JSON(400, gin.H{
				"error": true,
				"message": "client not added to event",
			})
			return
		}

		_, err = db.Exec("DELETE FROM PARTICIPATES WHERE Id_CLIENTS = '" + idclient + "' AND Id_EVENTS = '" + idevent + "'")
		fmt.Println(err)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot remove participates",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
			"message": "client removed from event",
		})
	}
}

func ValidateClientPresence(tokenAPI string) func(c *gin.Context) {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header["Token"]

		if len(tokenHeader) == 0 {
			c.JSON(400, gin.H{
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

		idEvent := c.Param("idevent")
		if idEvent == "" {
			c.JSON(400, gin.H{
				"error": true,
				"message": "id event can't be empty",
			})
			return
		}

		if !utils.IsSafeString(idEvent) {
			c.JSON(400, gin.H{
				"error": true,
				"message": "id event can't contain sql injection",
			})
			return
		}

		iduser := c.Param("iduser")
		if iduser == "" {
			c.JSON(400, gin.H{
				"error": true,
				"message": "id client can't be empty",
			})
			return
		}

		if !utils.IsSafeString(iduser) {
			c.JSON(400, gin.H{
				"error": true,
				"message": "id client can't contain sql injection",
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

		var idevent string
		
		err = db.QueryRow("SELECT Id_EVENTS FROM EVENTS WHERE Id_EVENTS = '" + idEvent + "'").Scan(&idevent)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "event not found",
			})
			return
		}

		var idclient string
		
		err = db.QueryRow("SELECT Id_CLIENTS FROM CLIENTS WHERE Id_USERS = '" + iduser + "'").Scan(&idclient)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "client not found",
			})
			return
		}

		var ispresent bool

		err = db.QueryRow("SELECT IsPresent FROM PARTICIPATES WHERE Id_CLIENTS = '" + idclient + "' AND Id_EVENTS = '" + idevent + "'").Scan(&ispresent)
		if err != nil {
			c.JSON(400, gin.H{
				"error": true,
				"message": "client not added to event",
			})
			return
		}

		if ispresent {
			c.JSON(400, gin.H{
				"error": true,
				"message": "client already validated",
			})
			return
		}

		_, err = db.Exec("UPDATE PARTICIPATES SET IsPresent = ? WHERE Id_CLIENTS = '" + idclient + "' AND Id_EVENTS = '" + idevent + "'", true)
			if err != nil {
				c.JSON(500, gin.H{
					"error": true,
					"message": "cannot update participates",
				})
				return
			}

		c.JSON(200, gin.H{
			"error": false,
			"message": "client presence validated",
		})
	}
}
