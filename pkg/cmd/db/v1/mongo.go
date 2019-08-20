package v1

import (
	"github.com/nicolasassi/starWarsApi/pkg/swapi"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler interface {
	AddPlanet()
	ListPlanets()
	FindByName()
	FindById()
	RemovePlanet()
}

func NewHandler() *Handler {
	return new(Handler)
}

type Mongo struct {
	DB *mongo.Collection
}

func NewMongo(db *mongo.Database, collectionName string) *Mongo {
	return &Mongo{db.Collection(collectionName)}
}

func (m *Mongo) Update(planets swapi.SwapiPlanetResponse) error {
	return nil
}

func (m *Mongo) Generate() error {
	return nil
}

func (m *Mongo) AddPlanet() {
}
