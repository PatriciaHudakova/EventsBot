package main

import (
	"github.com/yanzay/tbot"
)

// Initial choice of buttons upon triggered /start command
func btnStartingChoices() *tbot.InlineKeyboardMarkup {
	btnCreate := tbot.InlineKeyboardButton{
		Text: "Create Event 	\U00002712",
		CallbackData: "/new",
	}

	btnHelp := tbot.InlineKeyboardButton{
		Text: "See what I can do!	\U00002728",
		CallbackData: "/help",
	}
	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]tbot.InlineKeyboardButton{
			{btnCreate, btnHelp},
		},
	}
}

// Button options for event change in /edit command
func btnOptionsChoices() *tbot.InlineKeyboardMarkup {
	btnChangeName := tbot.InlineKeyboardButton{
		Text:         "Rename \U00002712",
		CallbackData: "/newName",
	}

	btnChangeDate := tbot.InlineKeyboardButton{
		Text: "Date 	\U0001F4C5",
		CallbackData: "/newDate",
	}

	btnChangeTime := tbot.InlineKeyboardButton{
		Text: "Time 	\U0001F553",
		CallbackData: "/newTime",
	}

	btnDelete := tbot.InlineKeyboardButton{
		Text:         "Delete \U00002757",
		CallbackData: "/delete",
	}
	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]tbot.InlineKeyboardButton{
			{btnChangeName, btnChangeDate, btnChangeTime, btnDelete},
		},
	}
}
