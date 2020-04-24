package main

import (
	"github.com/yanzay/tbot"
)

// Initial choice of buttons upon triggered /start command
func btnStartingChoices() *tbot.InlineKeyboardMarkup {
	btnCreate := tbot.InlineKeyboardButton{
		Text:         "Create Event",
		CallbackData: "/createEvent",
	}

	btnHelp := tbot.InlineKeyboardButton{
		Text:         "See what I can do!",
		CallbackData: "/help",
	}
	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]tbot.InlineKeyboardButton{
			{btnCreate, btnHelp},
		},
	}
}
