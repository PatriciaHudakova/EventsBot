package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/yanzay/tbot"
	"strconv"
	"time"
)

//Handles the "/start" command and displays message with button choices
func (a *application) startHandler(request *tbot.Message) {
	buttons := btnStartingChoices()
	a.client.SendMessage(request.Chat.ID, "Hello "+request.Chat.FirstName+", I am your friendly Event Tracker Bot "+
		"and I'll be helping you with keeping track of your "+
		"events! What would you like to do?", tbot.OptInlineKeyboardMarkup(buttons))
}

//Handles any pressed buttons
func (a *application) buttonHandler(pressed *tbot.CallbackQuery) {
	if pressed.Data == "/help" {
		a.helpHandler(pressed.Message)
	} else if pressed.Data == "/new" {
		a.newHandler(pressed.Message)
	} else if pressed.Data == "/newName" {
		a.newNameHandler(pressed.Message)
	} else if pressed.Data == "/newDate" {
		a.newDateHandler(pressed.Message)
	} else if pressed.Data == "/newTime" {
		a.newTimeHandler(pressed.Message)
	} else {
		a.client.SendMessage(pressed.Message.Chat.ID, "Error Occured")
	}
}

//Handles /help command by displaying a list of Bot's functionalities
func (a *application) helpHandler(request *tbot.Message) {
	m := "/new: creates a new event \n/show: shows a current log of all events \n" +
		"/edit: allows you to edit your event \n/help: shows all the things I can do \n" +
		"/deleteAll: erases all events in database \n/today: displays all events happening today"
	a.client.SendMessage(request.Chat.ID, m)
}

//Handler to delete all events
func (a *application) deleteAllHandler(request *tbot.Message) {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres "+
		"password="+getPwd()+" dbname=eventsdb sslmode=disable")

	_, err = db.Exec("DELETE FROM events WHERE chatid = '" + request.Chat.ID + "'")
	if err != nil {
		panic(err)
	} else {
		a.client.SendMessage(request.Chat.ID, "Success: All events have been deleted!")
	}
	defer db.Close()
}

//A reminder for events on the same day
func (a *application) todayHandler(request *tbot.Message) {

	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres "+
		"password="+getPwd()+" dbname=eventsdb sslmode=disable")

	defer db.Close()

	//Only returns entries where the date is today and chatid is identical to request id
	rows, err := db.Query("SELECT name, date, time FROM events WHERE chatid = '" + request.Chat.ID + "' " +
		"AND date = '" + time.Now().Format("2006-01-02") + "' ORDER BY date ASC")

	entries := make([]dbColumns, 0)

	//Loop through the values of rows
	for rows.Next() {
		column := dbColumns{}
		err := rows.Scan(&column.name, &column.date, &column.time)
		if err != nil {
			panic(err)
		}
		entries = append(entries, column)
	}

	//Handle any errors
	if err = rows.Err(); err != nil {
		panic(err)
	}

	//Loop through entries and print all results in a separate message (if any)
	if len(entries) == 0 {
		a.client.SendMessage(request.Chat.ID, "You have no events for today")
	} else {
		a.client.SendMessage(request.Chat.ID, "You have the following events today:")
		for _, i := range entries {
			a.client.SendMessage(request.Chat.ID, i.name+" at "+i.time.Format("15:04"))
		}
	}
}

//Shows all events listed in the table
func (a *application) showEventsHandler(request *tbot.Message) {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres "+
		"password="+getPwd()+" dbname=eventsdb sslmode=disable")

	defer db.Close()

	//check for userID and only display entries that match to a user requesting
	rows, err := db.Query("SELECT name, date, time FROM events WHERE chatid = '" + request.Chat.ID + "' " +
		"ORDER BY date ASC")

	//Placeholder for an array slice
	entries := make([]dbColumns, 0)

	//Loop through the values of rows
	for rows.Next() {
		column := dbColumns{}
		err := rows.Scan(&column.name, &column.date, &column.time)
		if err != nil {
			panic(err)
		}
		entries = append(entries, column)
	}

	//Handle any errors
	if err = rows.Err(); err != nil {
		panic(err)
	}

	//Loop through and print all results (if any) in a separate message
	if len(entries) == 0 {
		a.client.SendMessage(request.Chat.ID, "You have no events scheduled, add some using /new!")
	} else {
		a.client.SendMessage(request.Chat.ID, "There are all your scheduled events:")
		for _, i := range entries {
			a.client.SendMessage(request.Chat.ID, i.name+" on "+i.date.Format("Mon, 02 Jan 2006")+" at "+
				i.time.Format("15:04"))
		}
	}
}

//Handles /createEvent through a chain of handlers resulting in a database entry based on user input
func (a *application) newHandler(request *tbot.Message) {
	eventChatId = request.Chat.ID
	a.client.SendMessage(request.Chat.ID, "Great! What would you like to call your event?: use format n<eventName>")
	bot.HandleMessage("[n][a-zA-Z]", app.eventNameHandler)

}

//Logs event name and asks for the Date
func (a *application) eventNameHandler(request *tbot.Message) {
	eventNameRAW := tbot.Message{Text: request.Text}.Text
	eventName = eventNameRAW[1:]
	a.client.SendMessage(request.Chat.ID, "Awesome, when will it happen? Format: YYYY-MM-DD")
	bot.HandleMessage("\\d{4}-\\d{2}-\\d{2}", app.eventDateHandler)
}

//Logs date and asks for a time input
func (a *application) eventDateHandler(request *tbot.Message) {
	eventDate = tbot.Message{Text: request.Text}.Text
	a.client.SendMessage(request.Chat.ID, "Perfect, what time? 24h format: HH:MM")
	bot.HandleMessage("^([0-9]|0[0-9]|1[0-9]|2[0-3]):([0-9]|[0-5][0-9])$", app.eventDBHandler)
}

// Logs time input and creates a database query based on user input and executes
func (a *application) eventDBHandler(request *tbot.Message) {
	eventTime = tbot.Message{Text: request.Text}.Text

	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres "+
		"password="+getPwd()+" dbname=eventsdb sslmode=disable")

	_, err = db.Exec("INSERT INTO events(name, date, time, chatid) VALUES " +
		"('" + eventName + "', '" + eventDate + "', '" + eventTime + "', '" + eventChatId + "');")
	if err != nil {
		a.client.SendMessage(request.Chat.ID, "I'm sorry, something went wrong adding to the database. "+
			"Please stick to valid date/ time and specified format.")
	} else {
		a.client.SendMessage(request.Chat.ID, " Great, all done! Send /show to see your new event!")
	}
	db.Close()
}

//Handles /edit command
func (a *application) editHandler(request *tbot.Message) {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres "+
		"password="+getPwd()+" dbname=eventsdb sslmode=disable")

	defer db.Close()

	//check for userID and only display Entries that match to a user requesting, creating a new column counting rows
	rows, err := db.Query("SELECT * FROM (SELECT name, date, time, ROW_NUMBER() OVER (ORDER BY date ASC)" +
		" FROM events WHERE chatid = '" + request.Chat.ID + "') AS derivedTable;")

	if err != nil {
		panic(err)
	}

	//Placeholder for an array slice
	Entries := make([]derivedTable, 0)

	//Loop through the values of rows
	for rows.Next() {
		column := derivedTable{}
		err := rows.Scan(&column.name, &column.date, &column.time, &column.rownum)
		if err != nil {
			panic(err)
		}
		Entries = append(Entries, column)
	}

	//Handle any errors
	if err = rows.Err(); err != nil {
		panic(err)
	}

	//Loop through and print all results (if any) in a separate message
	if len(Entries) == 0 {
		a.client.SendMessage(request.Chat.ID, "You have no events scheduled, add some using /new!")
	} else {
		a.client.SendMessage(request.Chat.ID, "Enter the index of the event you wish to edit using format: e<number>")
		for _, i := range Entries {
			a.client.SendMessage(request.Chat.ID, i.rownum+") "+i.name+" on "+i.date.Format("Mon, 02 Jan 2006")+" at "+
				i.time.Format("15:04"))
		}
		bot.HandleMessage("[e]\\d{1}", app.editEnterNameHandler)
	}
}

//Searched through events table to check if entry exists
func (a *application) editEnterNameHandler(request *tbot.Message) {
	buttons := btnOptionsChoices()
	searchIndexRAW := tbot.Message{Text: request.Text}.Text //raw user input
	searchIndex := searchIndexRAW[1:]                       //lose the first char to obtain valid input
	var searchIndexInt int
	searchIndexInt, _ = strconv.Atoi(searchIndex) //convert string to int
	postgresOffset := searchIndexInt - 1          //calculate offset to get intended event
	var postgresOffsetSTR string
	postgresOffsetSTR = strconv.Itoa(postgresOffset) //int to string convert to use in the query

	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres "+
		"password="+getPwd()+" dbname=eventsdb sslmode=disable")

	//check for userID and only display Entries that match to a user requesting, creating a new column counting rows
	rows, err := db.Query("SELECT * FROM (SELECT name, date, time, ROW_NUMBER() OVER (ORDER BY date ASC) " +
		" FROM events WHERE chatid = '" + request.Chat.ID + "') AS derivedTable;")
	defer db.Close()
	if err != nil {
		panic(err)
	}

	//Placeholder for an array slice
	entries := make([]derivedTable, 0)

	//Loop through the values of rows
	for rows.Next() {
		column := derivedTable{}
		err := rows.Scan(&column.name, &column.date, &column.time, &column.rownum)
		if err != nil {
			panic(err)
		}
		entries = append(entries, column)
	}

	//Handle any errors
	if err = rows.Err(); err != nil {
		panic(err)
	}

	//Check if entry exists
	if len(entries) == 0 || len(entries) < searchIndexInt {
		a.client.SendMessage(request.Chat.ID, "Entry doesn't exist")
	} else {
		db, _ := sql.Open("postgres", "host=localhost port=5432 user=postgres "+
			"password="+getPwd()+" dbname=eventsdb sslmode=disable")

		rows2, _ := db.Query("SELECT name FROM events WHERE chatid = '" + request.Chat.ID + "' ORDER BY date " +
			"ASC LIMIT 1 OFFSET " + postgresOffsetSTR)

		defer db.Close()

		//since rownum column isn't accessible, use temporary index to obtain event name and log for future reference
		for rows2.Next() {
			err := rows2.Scan(&eventName)
			if err != nil {
				panic(err)
			}
		}
		a.client.SendMessage(request.Chat.ID, "Got it, what would you like to edit?", tbot.OptInlineKeyboardMarkup(buttons))
	}

}

//Handle new name input
func (a *application) newNameHandler(request *tbot.Message) {
	a.client.SendMessage(request.Chat.ID, "Enter a new name using format: e<newName>")
	bot.HandleMessage("[e][a-zA-Z]", app.newNameDBHandler)
}

//Extract valid input and parse into the query
func (a *application) newNameDBHandler(request *tbot.Message) {
	newEventNameRAW := tbot.Message{Text: request.Text}.Text
	newEventName := newEventNameRAW[1:]

	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres "+
		"password="+getPwd()+" dbname=eventsdb sslmode=disable")

	_, err = db.Exec("UPDATE events SET name='" + newEventName + "' WHERE name='" + eventName + "'")

	if err != nil {
		a.client.SendMessage(request.Chat.ID, "An error occurred, please try again.")
	} else {
		a.client.SendMessage(request.Chat.ID, "All done! Run /show to see your changes!")
	}

	defer db.Close()
}

//Handle new date input
func (a *application) newDateHandler(request *tbot.Message) {
	a.client.SendMessage(request.Chat.ID, "Enter a new Date (YYYY:MM:DD):")
	bot.HandleMessage("[e]\\d{4}-\\d{2}-\\d{2}", app.newDateDBHandler)
}

//Extract valid date from raw input and insert into query
func (a *application) newDateDBHandler(request *tbot.Message) {
	newEventDateRAW := tbot.Message{Text: request.Text}.Text
	newEventDate := newEventDateRAW[1:]

	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres "+
		"password="+getPwd()+" dbname=eventsdb sslmode=disable")

	_, err = db.Exec("UPDATE events SET date='" + newEventDate + "' WHERE name='" + eventName + "'")

	if err != nil {
		a.client.SendMessage(request.Chat.ID, "Please enter a valid date.")
	} else {
		a.client.SendMessage(request.Chat.ID, "All done! Run /show to see your changes!")
	}

	defer db.Close()
}

//Handle new tim input
func (a *application) newTimeHandler(request *tbot.Message) {
	a.client.SendMessage(request.Chat.ID, "Enter a new time: HH:MM")
	bot.HandleMessage("[t]([0-1]?[0-9]|2[0-3]):[0-5][0-9]", app.newTimeDBHandler)
}

//Extract valid time from raw input and insert into query
func (a *application) newTimeDBHandler(request *tbot.Message) {
	newEventTimeRAW := tbot.Message{Text: request.Text}.Text
	newEventTime := newEventTimeRAW[1:]

	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres "+
		"password="+getPwd()+" dbname=eventsdb sslmode=disable")

	_, err = db.Exec("UPDATE events SET time='" + newEventTime + "' WHERE name='" + eventName + "'")

	if err != nil {
		a.client.SendMessage(request.Chat.ID, "Please enter a valid time.")
	} else {
		a.client.SendMessage(request.Chat.ID, "All done! Run /show to see your changes!")
	}

	defer db.Close()
}
