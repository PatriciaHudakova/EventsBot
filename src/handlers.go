package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/yanzay/tbot"
	"os"
)

//Handles the "/start" command and displays message with button choices
func (a *application) startHandler(request *tbot.Message) {
	buttons := btnStartingChoices()
	a.client.SendMessage(request.Chat.ID, "Hello, I am your friendly Event Tracker Bot "+
		"and I'll be helping you with keeping track of your "+
		"events! What would you like to do?", tbot.OptInlineKeyboardMarkup(buttons))
}

//Handle pressed buttons created in startHandler()
func (a *application) buttonHandler(pressed *tbot.CallbackQuery) {
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
		"/editEvent <EventName>: allows you to edit your event \n/countdown <EventName>: shows " +
		"how much time between now and specified event \n/options: shows all the things I can do \n" +
		"/deleteAll: erases all events in database"
	a.client.SendMessage(request.Chat.ID, m)
}

//Handles /createEvent command
func (a *application) createEventHandler(request *tbot.Message) {

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
}
