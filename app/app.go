package App

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	Service "guestBook/service"

	log "github.com/sirupsen/logrus"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

type App struct {
	Redis          *redis.Client
	Loc            *time.Location
	CollectionBook string
	Ctx            context.Context
}

var mail = Service.MailConfig{}

func (a *App) Create(w http.ResponseWriter, r *http.Request) {
	log.Info("Call Method Create...")

	var guest Guest

	req, _ := ioutil.ReadAll(r.Body)
	log.Debug(string(req))

	err := json.Unmarshal(req, &guest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error("Error Unmarshal Request, ", err.Error())
		return
	}
	guest.DateCreated = time.Now()

	response := ResponseData{}

	err = a.Insert(&guest)
	if err != nil {
		log.Info("Error Insert Data, ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		response.Status = false
		response.Message = "Error Push Data"
		response.Code = 199

		data, err := json.Marshal(response)
		if err != nil {
			log.Error("Error Marshal the Request, ", err.Error())
		}

		w.Write(data)
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Info("Sending Mail")
	mail.Ctx = a.Ctx
	mail.Send("in")
	log.Info("Create Success")
}

func (a *App) Checkout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Add("Content-Type", "application/json")
	checkout := Checkout{}
	response := ResponseData{}

	log.Info("Decode Request Checkout")
	err := json.NewDecoder(r.Body).Decode(&checkout)
	if err != nil {
		log.Error("Error Decode Json, ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		response.Status = false
		response.Message = "Error Checkout Data"
		response.Code = 188

		data, err := json.Marshal(response)
		if err != nil {
			log.Info("Error Marshal Response, ", err.Error())
		}

		w.Write(data)
		return
	}

	checkout.ID = vars["id"]
	log.Info("ID : ", vars["id"])
	log.Info("Fetch Checkout")
	guest, err := a.FetchData(vars["id"])
	if err != nil {
		log.Info("Error Fetch Data, ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		response.Status = false
		response.Message = "Error Checkout Data"
		response.Code = 187

		data, err := json.Marshal(response)
		if err != nil {
			log.Error("Error Marshal Request, ", err.Error())
		}

		w.Write(data)
		return
	}

	guest.Close = checkout.Close

	log.Info("Checkout Time: ", checkout.Close)
	err = a.Delete(vars["id"])
	if err != nil {
		log.Info("Error Deleting Data, ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		response.Status = false
		response.Message = "Error Checkout Data"
		response.Code = 189

		data, err := json.Marshal(response)
		if err != nil {
			log.Info("Error Marshal Request, ", err.Error())
		}

		w.Write(data)
		return
	}

	log.Info("Insert Data Checkout")
	err = a.Insert(&guest)
	if err != nil {
		log.Info("Error Insert Data, ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		response.Status = false
		response.Message = "Error Put Data"
		response.Code = 176

		data, err := json.Marshal(response)
		if err != nil {
			log.Error("Error Marshal Request, ", err.Error())
		}

		w.Write(data)
		return
	}

	response.Code = 100
	response.Message = "Success Checkout"
	response.Status = true
	response.Data = guest

	res, _ := json.Marshal(response)
	log.Info("Success Checkout: ", string(res))

	w.Write(res)
	log.Info("Sending Mail")
	mail.Ctx = a.Ctx
	mail.Send("out")
	log.Info("Checkout Success")
}

func (a *App) Show(w http.ResponseWriter, r *http.Request) {
	log.Info("Call Method Show...")
	w.Header().Add("Content-Type", "application/json")

	// r.Header.Get("correlation_id")
	response := ResponseData{}
	// var guest Guest
	guest, err := a.getData()
	if err != nil {
		log.Error("Error Get Data: ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		response.Status = false
		response.Message = "Error Get Data"
		response.Code = 104

		errEncode, err := json.Marshal(response)
		if err != nil {
			log.Error("Error Marshal Request, ", err.Error())
		}

		w.Write(errEncode)
		return
	}

	w.WriteHeader(http.StatusOK)
	response.Status = true
	response.Message = "Retrive Data"
	response.Code = 199
	response.Data = guest

	res, _ := json.Marshal(response)
	log.Info("Success Calling Method")
	w.Write(res)

}

func (a *App) GetData(w http.ResponseWriter, r *http.Request) {
	log.Info("Call Method GetData...")
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	response := ResponseData{}

	log.Info("ID : ", vars["id"])

	// Check Data if exists
	check, _ := a.Redis.HExists(a.CollectionBook, vars["id"]).Result()
	if !check {
		log.Error("ID Not exist")
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
		log.Error("Error Fetch Data: ", err.Error())
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
	log.Info("Success Call Method")
}
