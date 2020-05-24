package App

import (
	"errors"
	"strings"
	"time"

	"github.com/go-acme/lego/log"
)

type Guest struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Office         string    `json:"office"`
	Phone          int       `json:"phone"`
	Address        string    `json:"address"`
	IdentityNumber int       `json:"identity_number"`
	Meet           string    `json:"meet"`
	VisitorNumber  int       `json:"visitor_number"`
	OnTime         string    `json:"on_time"`
	Close          string    `json:"close"`
	Concern        string    `json:"concern"`
	DateCreated    time.Time `json:"created_date"`
}

type Checkout struct {
	ID    string `json:"id"`
	Close string `json:"close"`
}

type Office struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Floor       string    `json:"floor"`
	Email       string    `json:"email"`
	DateCreated time.Time `json:"created_date"`
}

type TimeRequest struct {
	time.Time
}

func (t *TimeRequest) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")

	dt, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		log.Println(err.Error())
		return errors.New(err.Error())
	}

	t.Time = dt
	return nil
}

type ResponseData struct {
	Status  bool        `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
