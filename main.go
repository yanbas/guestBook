package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	App "guestBook/app"

	log "github.com/sirupsen/logrus"

	"github.com/buger/jsonparser"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

func main() {
	route := mux.NewRouter()
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)

	file, err := os.Open("./config/config.json")
	if err != nil {
		log.Error("Error Open File, ", err.Error())
	}

	config, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error("Error Read File, ", err.Error())
	}

	//email initialize
	sender, _ := jsonparser.GetString(config, "smtp", "user")
	password, _ := jsonparser.GetString(config, "smtp", "password")
	portMail, _ := jsonparser.GetString(config, "smtp", "port")
	host, _ := jsonparser.GetString(config, "smtp", "host")
	portMailer, _ := strconv.Atoi(portMail)

	// Set Context For Mail
	contextParent := context.Background()
	ctxMail := context.WithValue(contextParent, "mailHost", host)
	ctxMail = context.WithValue(ctxMail, "mailSender", sender)
	ctxMail = context.WithValue(ctxMail, "mailPassword", password)
	ctxMail = context.WithValue(ctxMail, "mailSender", sender)
	ctxMail = context.WithValue(ctxMail, "mailPort", portMailer)

	port, _ := jsonparser.GetString(config, "server", "port")
	redisHost, _ := jsonparser.GetString(config, "redis", "server")
	redisPassword, _ := jsonparser.GetString(config, "redis", "password")
	redisDb, _ := jsonparser.GetString(config, "redis", "db")
	db, _ := strconv.Atoi(redisDb)

	app := App.App{}

	app.Ctx = ctxMail

	connRedis := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPassword,
		DB:       db,
	})

	timezone, _ := time.LoadLocation("Asia/Jakarta")
	collection, _ := jsonparser.GetString(config, "collection", "guestBook")

	app.Loc = timezone
	app.Redis = connRedis
	app.CollectionBook = collection

	route.HandleFunc("/guest", app.Create).Methods("POST")
	route.HandleFunc("/guest/checkout/{id}", app.Checkout).Methods("POST")
	route.HandleFunc("/guest", app.Show).Methods("GET")
	route.HandleFunc("/guest/{id}", app.GetData).Methods("GET")

	log.Info("Running On Port : ", port)
	http.ListenAndServe(port, route)

}
