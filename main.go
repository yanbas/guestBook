package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	App "guestBook/app"

	"github.com/buger/jsonparser"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

func main() {
	route := mux.NewRouter()

	file, err := os.Open("./config/config.json")
	if err != nil {
		log.Println(err.Error())
	}

	config, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err.Error())
	}

	port, _ := jsonparser.GetString(config, "server", "port")
	redisHost, _ := jsonparser.GetString(config, "redis", "server")
	redisPassword, _ := jsonparser.GetString(config, "redis", "password")
	redisDb, _ := jsonparser.GetString(config, "redis", "db")
	db, _ := strconv.Atoi(redisDb)

	app := App.App{}

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
	route.HandleFunc("/guest", app.Show).Methods("GET")
	route.HandleFunc("/guest/{id}", app.GetData).Methods("GET")

	log.Println("running on port : ", port)
	http.ListenAndServe(port, route)

}
