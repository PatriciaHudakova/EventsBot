package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/yanzay/tbot"
)

//Handles the "/start" command and displays message with button choices
func (a *application) startHandler(request *tbot.Message) {
	buttons := btnStartingChoices()
	a.client.SendMessage(request.Chat.ID, "Hello "+request.Chat.FirstName+", I am your friendly Event Tracker Bot "+
		"and I'll be helping you with keeping track of your "+
		"events! What would you like to do?", tbot.OptInlineKeyboardMarkup(buttons))
}

//Handle pressed buttons created in startHandler()
func (a *application) startButtonHandler(pressed *tbot.CallbackQuery) {
	if pressed.Data == "/help" {
		a.helpHandler(pressed.Message)
	} else if pressed.Data == "/new" {
		a.newHandler(pressed.Message)
	} else {
		a.client.SendMessage(pressed.Message.Chat.ID, "Error Occured. Type /help for options.")
	}
}

//Handles /help command by displaying a list of Bot's functionalities
func (a *application) helpHandler(request *tbot.Message) {
	m := "/new: creates a new event \n/show: shows a current log of all events \n" +
		"/editEvent <EventName>: allows you to edit your event \n/help: shows all the things I can do \n" +
		"/deleteAll: erases all events in database"
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

//Shows all events listed in the table
func (a *application) showEventsHandler(request *tbot.Message) {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres "+
		"password="+getPwd()+" dbname=eventsdb sslmode=disable")

	if err != nil {
		panic(err)
	}
	defer db.Close()

	//check for userID and only display entries that match to a user requesting
	rows, err := db.Query("SELECT name, date, time FROM events WHERE chatid = '" + request.Chat.ID + "'")

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

	//Loop through and print all results in a separate message
	for _, i := range entries {
		a.client.SendMessage(request.Chat.ID, i.name+" on "+i.date.Format("Mon, 02 Jan 2006")+" at "+
			i.time.Format("15:04"))
	}
}

//Handles /createEvent command using a chain of handlers creating a database entry based on user input
func (a *application) newHandler(request *tbot.Message) {
	eventChatId = request.Chat.ID
	a.client.SendMessage(request.Chat.ID, "Great! What would you like to call your event?")
	bot.HandleMessage("[a-zA-Z]", app.eventNameHandler)
}

//Logs event name and asks for the Date
func (a *application) eventNameHandler(request *tbot.Message) {
	eventName = tbot.Message{Text: request.Text}.Text
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
	a.client.SendMessage(request.Chat.ID, "Let me make a note for you ^.^")

	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres "+
		"password="+getPwd()+" dbname=eventsdb sslmode=disable")

	_, err = db.Exec("INSERT INTO events(name, date, time, chatid) VALUES " +
		"('" + eventName + "', '" + eventDate + "', '" + eventTime + "', '" + eventChatId + "');")
	if err != nil {
		panic(err)
	} else {
		a.client.SendMessage(request.Chat.ID, " Great, all done! Send /show to see your new event!")
	}
	db.Close()
}
