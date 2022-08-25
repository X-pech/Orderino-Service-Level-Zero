package cacheController

import (
	"wbservice/config"
	"wbservice/controllers/dbController"
)

type CacheController struct {
	data map[string]string
	dbController.DBController
}

func New(psqlConfig config.PostgresConfig) (CacheController, error) {
	result := new(CacheController)
	result.data = make(map[string]string)
	var err error
	result.DBController, err = dbController.New(psqlConfig)
	if err != nil {
		return *result, err
	}
	err = result.LoadFromDB()
	if err != nil {
		return *result, err
	}
	return *result, nil
}

func (cc *CacheController) Push(order_uid string, object string) error {
	cc.data[order_uid] = object
	err := cc.LogToDB(order_uid, object)
	return err
}

func (cc *CacheController) Erase(order_uid string) error {
	delete(cc.data, order_uid)
	_, err := cc.DBController.RemoveKey(order_uid)
	return err
}

func (cc *CacheController) Get(order_uid string) string {
	return cc.data[order_uid]
}

func (cc *CacheController) LoadFromDB() error {
	offset := 0
	limit := 100
	for {
		rows, next, err := cc.LoadBatch(limit, offset)
		if err != nil {
			return err
		}
		offset += limit
		for i := range *rows {
			cc.data[(*rows)[i][0]] = (*rows)[i][1]
		}

		if !next {
			break
		}
	}
	return nil
}
