package main

import (
	"github.com/yanzay/tbot"
)

//Handles the "/start" command
func (a *application) startHandler(request *tbot.Message) {
	m := "Hello, I am your friendly EventsBot and I'll be helping you with keeping track of your " +
		"events! What would you like to do?"
	a.client.SendMessage(request.Chat.ID, m)
}

func (a *application) optionsHandler(request *tbot.Message) {
	m := "/createEvent: creates a new event \n/showEvents: shows a current log of all events \n" +
		"/editEvent <EventName>: allows you to edit your event \n/countdown <EventName>: shows " +
		"how much time between now and specified event \n/options: shows all the things I can do \n" +
		"/joke: pulls a random (not so great) joke from the internet"
	a.client.SendMessage(request.Chat.ID, m)
}
