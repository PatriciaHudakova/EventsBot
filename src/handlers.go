package main

import (
	"github.com/yanzay/tbot"
)

//Handles the "/start" command and displays choices
func (a *application) startHandler(m *tbot.Message) {
	btns := btnStartingChoices()
	a.client.SendMessage(m.Chat.ID, "Hello, I am your friendly Event Tracker Bot and I'll be helping you with keeping track of your "+
		"events! What would you like to do?", tbot.OptInlineKeyboardMarkup(btns))
}

//Handles /help command by displaying a list of Bot's functionalities
func (a *application) helpHandler(request *tbot.Message) {
	m := "/createEvent: creates a new event \n/showEvents: shows a current log of all events \n" +
		"/editEvent <EventName>: allows you to edit your event \n/countdown <EventName>: shows " +
		"how much time between now and specified event \n/options: shows all the things I can do \n" +
		"/deleteAll: erases all events in database"
	a.client.SendMessage(request.Chat.ID, m)
}
