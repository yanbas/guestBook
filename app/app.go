package App

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	var guest Guest

	req, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(req, &guest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := ResponseData{}

	err = a.Insert(&guest)
	if err != nil {
		log.Println(err.Error)
		w.WriteHeader(http.StatusInternalServerError)
		response.Status = false
		response.Message = "Error Push Data"
		response.Code = 199

		data, err := json.Marshal(response)
		log.Println(err.Error())

		w.Write(data)
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *App) Show(w http.ResponseWriter, r *http.Request) {
	log.Println("Call Method Show...")

}

func (a *App) GetData(w http.ResponseWriter, r *http.Request) {
	log.Println("Call Method GetData...")

	vars := mux.Vars(r)

	fmt.Println(vars["id"])

}
