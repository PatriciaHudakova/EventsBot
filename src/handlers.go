package main

import (
	"github.com/yanzay/tbot"
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

//Handles /createEvent
func (a *application) createEventHandler(request *tbot.Message) {
	return //TBD
}
