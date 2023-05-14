package events

type Event struct {
	IdEvent int `json:"idevent"`
	Type string `json:"type"`
	EndTime string `json:"endtime"`
	StartTime string `json:"starttime"`
	isInternal bool `json:"isinternal"`
	isPrivate bool `json:"isprivate"`
	groupDisplayOrder int `json:"groupdisplayorder"`
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
			err = rows.Scan(&event.IdEvent, &event.Type, &event.EndTime, &event.StartTime, &event.isInternal, &event.isPrivate, &event.groupDisplayOrder, &event.IdGroups)
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

		err = db.QueryRow("SELECT * FROM EVENTS WHERE Id_EVENTS = ?", id).Scan(&event.IdEvent, &event.Type, &event.EndTime, &event.StartTime, &event.isInternal, &event.isPrivate, &event.groupDisplayOrder, &event.IdGroups)
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
			err = rows.Scan(&event.IdEvent, &event.Type, &event.EndTime, &event.StartTime, &event.isInternal, &event.isPrivate, &event.groupDisplayOrder, &event.IdGroups)
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

		var idGroups int

		err = db.QueryRow("SELECT Id_GROUPS FROM GROUPS WHERE Id_GROUPS = ?", event.IdGroups).Scan(&idGroups)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "group does not exist",
			})
			return
		}

		err = db.QueryRow("SELECT Id_EVENTS FROM EVENTS WHERE Id_GROUPS = ? AND group_display_order = ?", event.IdGroups, event.groupDisplayOrder).Scan(&idGroups)
		if err == nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "event already exist at this display order",
			})
			return
		}

		stmt, err := db.Prepare("INSERT INTO EVENTS (Type, EndTime, StartTime, isInternal, isPrivate, GroupDisplayOrder, Id_GROUPS) VALUES (?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot prepare request",
			})
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(event.Type, event.EndTime, event.StartTime, event.isInternal, event.isPrivate, event.groupDisplayOrder, event.IdGroups)
		if err != nil {
			c.JSON(500, gin.H{
				"error": true,
				"message": "cannot insert event",
			})
			return
		}

		c.JSON(200, gin.H{
			"error": false,
			"message": "event inserted",
		})
	}
}