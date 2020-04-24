package main

import (
	"github.com/joho/godotenv" //Loads environment variables from '.env'
	"github.com/yanzay/tbot"   //Go library for Telegram Bot API
	"log"
	"os"
)

//declare a struct called application
type application struct {
	client *tbot.Client //(pointer) works with the actual value as opposed to a copy
}

var (
	app   application
	bot   *tbot.Server
	token string
)

//initialise environment before main() launches
func init() {
	err := godotenv.Load() //assign and declare an error variable
	if err != nil {        //if there is an error during env launch
		log.Fatalln(err)
	}
	token = os.Getenv("TELEGRAM_TOKEN")
}

func main() {
	bot = tbot.New(token)
	app.client = bot.Client()
	bot.HandleMessage("/start", app.startHandler)
	bot.HandleCallback(app.buttonHandler)
	bot.HandleMessage("/help", app.helpHandler)
	bot.HandleMessage("/createEvent", app.createEventHandler)
	log.Fatal(bot.Start())
}
