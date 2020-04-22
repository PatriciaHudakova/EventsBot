package main

import (
	"github.com/joho/godotenv" //Loads environment variables from '.env'
	"github.com/yanzay/tbot"   //Go library for Telegram Bot API
	"log"
	"os"
)

//declare a struct called application
type application struct {
	client *tbot.Client
}

//Handle "/start" command
func (a *application) startHandler(request *tbot.Message) {
	m := "Hello, I am your friendly EventsBot and I'll be helping you with keeping track of your " +
		"events! What would you like to do?"
	a.client.SendMessage(request.Chat.ID, m)
}

//assign variables to their types
var (
	app   application
	bot   *tbot.Server
	token string
)

//initialise environment before main() launch
func init() {
	err := godotenv.Load() //assign and declare an error variable
	if err != nil {        //if there is an error loading environment
		log.Fatalln(err)
	}
	token = os.Getenv("TELEGRAM_TOKEN")
}

func main() {
	bot = tbot.New(token)
	app.client = bot.Client()
	bot.HandleMessage("/start", app.startHandler)
	log.Fatal(bot.Start())
}
