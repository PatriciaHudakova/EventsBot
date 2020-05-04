package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv" //Loads environment variables from '.env'
	_ "github.com/lib/pq"
	"github.com/yanzay/tbot" //Go library for Telegram Bot API
	"log"
	"os"
	"time"
)

//Declare a struct called application
type application struct {
	client *tbot.Client //(pointer) works with the actual value as opposed to a copy
}

//Visible user columns
type dbColumns struct {
	name string
	date time.Time
	time time.Time
}

//Derived Table from row_number() used for indexing in /edit and /delete
type derivedTable struct {
	name   string
	date   time.Time
	time   time.Time
	rownum string
}

var (
	app         application
	bot         *tbot.Server
	token       string
	eventName   string
	eventDate   string
	eventTime   string
	eventChatId string
	db          *sql.DB
	port        = os.Getenv("PORT")
	publicURL   = os.Getenv("PUBLIC_URL")
)

//Initialise environment before main() launches
func init() {
	err := godotenv.Load() //assign and declare an error variable
	if err != nil {        //if there is an error during env launch
		log.Fatalln(err)
	}
	token = os.Getenv("TELEGRAM_BOT_TOKEN")

	checkDBConnection()
}

//Main entry for the program
func main() {
	bot = tbot.New(token, tbot.WithWebhook("PUBLIC_URL", ":"+os.Getenv("PORT")))
	app.client = bot.Client()

	//All handler-related code
	bot.HandleMessage("/start", app.startHandler)
	bot.HandleCallback(app.buttonHandler)
	bot.HandleMessage("/help", app.helpHandler)
	bot.HandleMessage("/new", app.newHandler)
	bot.HandleMessage("/deleteAll", app.deleteAllHandler)
	bot.HandleMessage("/show", app.showEventsHandler)
	bot.HandleMessage("/today", app.todayHandler)
	bot.HandleMessage("/edit", app.editHandler)
	log.Fatal(bot.Start())
}

//Opens and verifies a connection to our database, handles any errors
func checkDBConnection() {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres "+
		"password="+getPwd()+" dbname=eventsdb sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("Error: Cannot connect to database")
	}

	fmt.Println("Successfully connected!")
}

//Factory function to return env var
func getPwd() string {
	pwd := os.Getenv("POSTGRES_PASSWORD")
	return pwd
}
