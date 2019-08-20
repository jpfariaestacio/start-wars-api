package db

import (
	v1 "github.com/nicolasassi/starWarsApi/pkg/cmd/db/v1"
	"github.com/nicolasassi/starWarsApi/pkg/swapi"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	version1 = iota + 1
)

type Utils interface {
	Generate() error
	Update(planets swapi.SwapiPlanetResponse) error
}

func NewUtils(dbv uint, db *mongo.Database, collectionName string) Utils {
	switch dbv {
	case version1:
		return v1.NewMongo(db, collectionName)
	}
	return nil
}
