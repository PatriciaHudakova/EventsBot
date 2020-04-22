package main

import "github.com/yanzay/tbot"

//Handles the "/start" command
func (a *application) startHandler(request *tbot.Message) {
	m := "Hello, I am your friendly EventsBot and I'll be helping you with keeping track of your " +
		"events! What would you like to do?"
	a.client.SendMessage(request.Chat.ID, m)
}
