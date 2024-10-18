package database

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

func Connect(dbName string) error {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		panic("You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	err := mgm.SetDefaultConfig(nil, dbName, options.Client().ApplyURI(uri))

	if err != nil {
		return err
	}
	return nil
}
