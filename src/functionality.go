package main

import (
	"github.com/yanzay/tbot"
)

// Initial choice of buttons upon triggered /start command
func btnStartingChoices() *tbot.InlineKeyboardMarkup {
	btnCreate := tbot.InlineKeyboardButton{
		Text:         "Create Event",
		CallbackData: "/new",
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

// Button options for event change in /edit command
func btnOptionsChoices() *tbot.InlineKeyboardMarkup {
	btnChangeName := tbot.InlineKeyboardButton{
		Text:         "Rename",
		CallbackData: "/newName",
	}

	btnChangeDate := tbot.InlineKeyboardButton{
		Text:         "Date",
		CallbackData: "/newDate",
	}

	btnChangeTime := tbot.InlineKeyboardButton{
		Text:         "Time",
		CallbackData: "/newTime",
	}

	btnDelete := tbot.InlineKeyboardButton{
		Text:         "Delete",
		CallbackData: "/delete",
	}
	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]tbot.InlineKeyboardButton{
			{btnChangeName, btnChangeDate, btnChangeTime, btnDelete},
		},
	}
}
