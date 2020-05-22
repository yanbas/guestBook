package App

import "time"

type Guest struct {
	ID             string    `json:"id`
	Name           string    `json:"name`
	Office         string    `json:"office`
	Phone          int       `json:"phone`
	Address        string    `json:"address`
	IdentityNumber int       `json:"identity_number`
	Meet           string    `json:"meet`
	VisitorNumber  int       `json:"visitor_number`
	OnTime         time.Time `json:"on_time`
	Close          time.Time `json:"close`
	Concern        string    `json:"concern`
	DateCreated    time.Time `json:"created_date`
}

type Office struct {
	ID          string    `json:"id`
	Name        string    `json:"name`
	Floor       string    `json:"floor`
	Email       string    `json:"email`
	DateCreated time.Time `json:"created_date`
}
