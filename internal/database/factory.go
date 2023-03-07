package database

import (
	"github.com/akhmettolegen/proxy/internal/database/drivers"
	"github.com/akhmettolegen/proxy/internal/database/drivers/mongo"
)

func New(conf drivers.DataStoreConfig) (drivers.DataStore, error) {
	if conf.DataStoreName == "mongo" {
		return mongo.New(conf)
	}

	return nil, ErrDatastoreNotImplemented
}
