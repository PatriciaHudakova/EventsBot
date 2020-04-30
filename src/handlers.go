package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/yanzay/tbot"
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
		"/edit <EventName>: allows you to edit your event \n/help: shows all the things I can do \n" +
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

//Sends a reminder for events on the same day
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

	//Loop through and print all results in a separate message (if any)
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
	rows, err := db.Query("SELECT name, date, time FROM events WHERE chatid = '" + request.Chat.ID + "' ORDER BY date ASC")

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

//Handles /edit command
func (a *application) editHandler(request *tbot.Message) {
	a.client.SendMessage(request.Chat.ID, "Okay, what is the event called?")
	bot.HandleMessage("[a-zA-Z]", app.editEnterNameHandler)
}

//Searched through events table to check if entry exists
func (a *application) editEnterNameHandler(request *tbot.Message) {
	searchEvent := tbot.Message{Text: request.Text}.Text

	db, _ := sql.Open("postgres", "host=localhost port=5432 user=postgres "+
		"password="+getPwd()+" dbname=eventsdb sslmode=disable")

	defer db.Close()

	rows, err := db.Query("SELECT name, date, time FROM events WHERE chatid = '" + request.Chat.ID + "' AND " +
		"name = '" + searchEvent + "' ORDER BY date ASC")

	if err != nil {
		panic(err)
	}

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

	//if there are no entries matching user input
	if len(entries) != 1 {
		a.client.SendMessage(request.Chat.ID, "I'm sorry, there are no events of this name. Send /show to check")
	} else {
		for _, i := range entries {
			eventName = i.name
			buttons := btnOptionsChoices()
			a.client.SendMessage(request.Chat.ID, "Got it, what would you like to change?", tbot.OptInlineKeyboardMarkup(buttons))
		}
	}
}

func (a *application) newNameHandler(request *tbot.Message) {
	a.client.SendMessage(request.Chat.ID, "Enter a new name:")
	bot.HandleMessage("[0-9]", app.newNameDBHandler)
}

func (a *application) newNameDBHandler(request *tbot.Message) {
	newEventName := tbot.Message{Text: request.Text}.Text

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

func (a *application) newDateHandler(request *tbot.Message) {
	a.client.SendMessage(request.Chat.ID, "Enter a new Date (YYYY:MM:DD):")
	bot.HandleMessage("\\d{4}-\\d{2}-\\d{2}", app.newDateDBHandler)
}

func (a *application) newDateDBHandler(request *tbot.Message) {
	newEventDate := tbot.Message{Text: request.Text}.Text

	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres "+
		"password="+getPwd()+" dbname=eventsdb sslmode=disable")

	_, err = db.Exec("UPDATE events SET date='" + newEventDate + "' WHERE name='" + eventName + "'")

	if err != nil {
		a.client.SendMessage(request.Chat.ID, "An error occurred, please try again.")
	} else {
		a.client.SendMessage(request.Chat.ID, "All done! Run /show to see your changes!")
	}

	defer db.Close()
}

func (a *application) newTimeHandler(request *tbot.Message) {
	a.client.SendMessage(request.Chat.ID, "Enter a new Time (YYYY:MM:DD:")
	bot.HandleMessage("^([0-9]|0[0-9]|1[0-9]|2[0-3]):([0-9]|[0-5][0-9])$", app.newTimeDBHandler)
}

func (a *application) newTimeDBHandler(request *tbot.Message) {
	newEventTime := tbot.Message{Text: request.Text}.Text

	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres "+
		"password="+getPwd()+" dbname=eventsdb sslmode=disable")

	_, err = db.Exec("UPDATE events SET time='" + newEventTime + "' WHERE name='" + eventName + "'")

	if err != nil {
		a.client.SendMessage(request.Chat.ID, "An error occurred, please try again.")
	} else {
		a.client.SendMessage(request.Chat.ID, "All done! Run /show to see your changes!")
	}

	defer db.Close()
}
