package App

import (
	"encoding/json"
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
	guest.DateCreated = time.Now()

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
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *App) Checkout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	checkout := Checkout{}
	response := ResponseData{}

	err := json.NewDecoder(r.Body).Decode(&checkout)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.Status = false
		response.Message = "Error Checkout Data"
		response.Code = 188

		data, err := json.Marshal(response)
		log.Println(err.Error())

		w.Write(data)
		return
	}

	checkout.ID = vars["id"]

	guest, err := a.FetchData(vars["id"])
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		response.Status = false
		response.Message = "Error Checkout Data"
		response.Code = 187

		data, err := json.Marshal(response)
		log.Println(err.Error())

		w.Write(data)
		return
	}

	guest.Close = checkout.Close

	err = a.Delete(vars["id"])
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		response.Status = false
		response.Message = "Error Checkout Data"
		response.Code = 189

		data, err := json.Marshal(response)
		log.Println(err.Error())

		w.Write(data)
		return
	}

	response.Code = 100
	response.Message = "Success Checkout"
	response.Status = true
	response.Data = guest

	res, _ := json.Marshal(response)

	w.Write(res)

}

func (a *App) Show(w http.ResponseWriter, r *http.Request) {
	log.Println("Call Method Show...")
	w.Header().Add("Content-Type", "application/json")
	response := ResponseData{}
	// var guest Guest
	guest, err := a.getData()
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		response.Status = false
		response.Message = "Error Get Data"
		response.Code = 104

		errEncode, _ := json.Marshal(response)

		w.Write(errEncode)
		return
	}

	w.WriteHeader(http.StatusOK)
	response.Status = true
	response.Message = "Retrive Data"
	response.Code = 199
	response.Data = guest

	res, _ := json.Marshal(response)

	w.Write(res)

}

func (a *App) GetData(w http.ResponseWriter, r *http.Request) {
	log.Println("Call Method GetData...")
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	response := ResponseData{}

	// Check Data if exists
	check, _ := a.Redis.HExists(a.CollectionBook, vars["id"]).Result()
	if !check {
		log.Println("ID Not exist")
		w.WriteHeader(http.StatusNotFound)
		response.Status = false
		response.Message = "ID Not Found"
		response.Code = 104

		errEncode, _ := json.Marshal(response)

		w.Write(errEncode)
		return
	}

	data, err := a.FetchData(vars["id"])
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		response.Status = false
		response.Message = "Error Get Data"
		response.Code = 104

		errEncode, _ := json.Marshal(response)

		w.Write(errEncode)
		return
	}

	response.Code = 100
	response.Message = "Retrive Data"
	response.Status = true
	response.Data = data

	res, _ := json.Marshal(response)

	w.WriteHeader(http.StatusOK)
	w.Write(res)

}
