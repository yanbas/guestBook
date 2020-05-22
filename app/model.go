package App

import (
	"errors"

	"github.com/go-acme/lego/v3/log"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/square/go-jose.v2/json"
)

func (s *App) Insert(param *Guest) error {
	var data = make(map[string]interface{})

	raw, err := json.Marshal(param)
	if err != nil {
		log.Println(err.Error())
	}

	data[uuid.NewV4().String()] = raw

	_, err = s.Redis.HMSet(s.CollectionBook, data).Result()

	if err != nil {
		return errors.New("Error Set Data : " + err.Error())
	}

	log.Println("done")

	return nil
}
