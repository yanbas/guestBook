package App

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"log"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

type App struct {
	Redis          *redis.Client
	Loc            *time.Location
	CollectionBook string
}

func (a *App) Create(w http.ResponseWriter, r *http.Request) {
	log.Println("Call Method Create...")

	// err := json.NewDecoder(r.Body).Decode(&guest)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// loc, _ := time.LoadLocation("Asia/Jakarta")
	// datee := "2010-01-23 11:44:20"
	// dt, _ := time.Parse("2006-01-02 15:04:05", datee)
	// dtstr2 := dt.Format("Jan 2 '06 at 15:04")
	// guest.Tanggal = dt.In(loc)

	// log.Println(guest)

	var guest Guest

	err := json.NewDecoder(r.Body).Decode(&guest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.Insert(&guest)
	if err != nil {
		log.Println(err.Error)
	}
}

func (a *App) Show(w http.ResponseWriter, r *http.Request) {
	log.Println("Call Method Show...")

}

func (a *App) GetData(w http.ResponseWriter, r *http.Request) {
	log.Println("Call Method GetData...")

	vars := mux.Vars(r)

	fmt.Println(vars["id"])

}
