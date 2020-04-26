package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv" //Loads environment variables from '.env'
	_ "github.com/lib/pq"
	"github.com/yanzay/tbot" //Go library for Telegram Bot API
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

//Initialise environment before main() launches
func init() {
	err := godotenv.Load() //assign and declare an error variable
	if err != nil {        //if there is an error during env launch
		log.Fatalln(err)
	}
	token = os.Getenv("TELEGRAM_TOKEN")

	//Call a method to create and connection to our database
	openDBConnection()
}

func main() {
	bot = tbot.New(token)
	app.client = bot.Client()

	//All handler-related code
	bot.HandleMessage("/start", app.startHandler)
	bot.HandleCallback(app.buttonHandler)
	bot.HandleMessage("/help", app.helpHandler)
	bot.HandleMessage("/createEvent", app.createEventHandler)
	bot.HandleMessage("/deleteAll", app.deleteAllHandler)
	log.Fatal(bot.Start())
}

//Opens and verifies a connection to our database, handles any errors
func openDBConnection() {
	pwd := os.Getenv("POSTGRES_PASSWORD")

	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres "+
		"password="+pwd+" dbname=eventsdb sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
}
