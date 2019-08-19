package db

import (
	"github.com/nicolasassi/starWarsApi/pkg/swapi"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	v1 = iota + 1
)

type Handler interface {
	Generate() error
	Update(planets swapi.SwapiPlanetResponse, collectionName string) error
}

type V1 struct {
	DB *mongo.Database
}

func NewHandler(dbtype uint) Handler {
	switch dbtype {
	case v1:
		return new(V1)
	}
	return nil
}

func (v *V1) Update(planets swapi.SwapiPlanetResponse, collectionName string) error {
	v.DB.Collection(collectionName)
}

func (v *V1) Generate() error {

}
