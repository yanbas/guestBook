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

func (s *App) GetDataById(id string) {

	s.FetchData(id)

}

func (s *App) getData() ([]Guest, error) {
	listKey, err := s.Redis.HKeys(s.CollectionBook).Result()
	guest := []Guest{}
	if err != nil {
		log.Println("Error HKEYS")
		return []Guest{}, errors.New(err.Error())
	}

	for _, v := range listKey {
		result, err := s.FetchData(v)
		if err != nil {
			log.Println("Error Fetch Data, " + err.Error())
		}

		guest = append(guest, result)

	}

	return guest, nil

}

func (s *App) FetchData(id string) (Guest, error) {

	data, _ := s.Redis.HGet(s.CollectionBook, id).Result()

	var guest Guest

	err := json.Unmarshal([]byte(data), &guest)
	if err != nil {
		return guest, errors.New("Error Parse JSON, " + err.Error())
	}

	return guest, nil
}

func (s *App) Delete(id string) error {
	_, err := s.Redis.HDel(s.CollectionBook, id).Result()
	if err != nil {
		return errors.New("Error Deleting Data, " + err.Error())
	}

	return nil
}
