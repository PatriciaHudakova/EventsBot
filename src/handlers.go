package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/yanzay/tbot"
	"os"
	"time"
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
	} else if pressed.Data == "/createEvent" {
		a.createEventHandler(pressed.Message)
	} else {
		a.client.SendMessage(pressed.Message.Chat.ID, "Error Occured. Type /help for options.")
	}
}

//Handles /help command by displaying a list of Bot's functionalities
func (a *application) helpHandler(request *tbot.Message) {
	m := "/createEvent: creates a new event \n/showEvents: shows a current log of all events \n" +
		"/editEvent <EventName>: allows you to edit your event \n/options: shows all the things I can do \n" +
		"/deleteAll: erases all events in database"
	a.client.SendMessage(request.Chat.ID, m)
}

//Handler to delete all events
func (a *application) deleteAllHandler(request *tbot.Message) {
	pwd := os.Getenv("POSTGRES_PASSWORD")

	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres "+
		"password="+pwd+" dbname=eventsdb sslmode=disable")

	_, err = db.Exec("TRUNCATE events")
	if err != nil {
		panic(err)
	} else {
		a.client.SendMessage(request.Chat.ID, "Success: All events have been deleted!")
	}
	defer db.Close()
}

//Shows all events listen in the table
func (a *application) showEventsHandler(request *tbot.Message) {
	pwd := os.Getenv("POSTGRES_PASSWORD")

	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres "+
		"password="+pwd+" dbname=eventsdb sslmode=disable")

	rows, err := db.Query("SELECT * FROM events")
	if err != nil {
		panic(err)
	}
	defer db.Close()

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
		a.client.SendMessage(request.Chat.ID, i.name+" "+i.time.Format(time.RFC1123))
	}
}

//Handles /createEvent command
func (a *application) createEventHandler(request *tbot.Message) {
	a.client.SendMessage(request.Chat.ID, "Great! What would you like to call your event?")
	bot.HandleMessage(".", app.eventNameHandler)
}

func (a *application) eventNameHandler(request *tbot.Message) {
	eventName = tbot.Message{Text: request.Text}.Text
	a.client.SendMessage(request.Chat.ID, "Awesome, when will it happen?")
	bot.HandleMessage("\\d{4}-\\d{2}-\\d{2}", app.eventDateHandler)
}

func (a *application) eventDateHandler(request *tbot.Message) {
	eventDate = tbot.Message{Text: request.Text}.Text
	a.client.SendMessage(request.Chat.ID, "Perfect, what time?")
	bot.HandleMessage("^([0-9]|0[0-9]|1[0-9]|2[0-3]):([0-9]|[0-5][0-9])$", app.eventTimeHandler)
}

func (a *application) eventTimeHandler(request *tbot.Message) {
	eventTime = tbot.Message{Text: request.Text}.Text
	a.client.SendMessage(request.Chat.ID, "Let me make a note for you ^.^")
	bot.HandleMessage(".", app.eventDBHandler)
}

func (a *application) eventDBHandler(request *tbot.Message) {
	pwd := os.Getenv("POSTGRES_PASSWORD")

	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres "+
		"password="+pwd+" dbname=eventsdb sslmode=disable")

	_, err = db.Exec("INSERT INTO eventsdb VALUES (" + eventName + ", " + eventDate + ", " + eventTime + ")")
	if err != nil {
		panic(err)
	} else {
		a.client.SendMessage(request.Chat.ID, "Success: Event Created!")
	}

	a.client.SendMessage(request.Chat.ID, "Great, All Done!")
}
