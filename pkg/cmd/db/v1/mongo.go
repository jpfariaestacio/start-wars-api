package v1

import (
	"context"

	"github.com/nicolasassi/starWarsApi/pkg/swapi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler interface {
	AddPlanet(interface{}) error
	ListPlanets() ([]Planet, error)
	FindByName(name string) (Planet, error)
	FindById(id int) (Planet, error)
	RemovePlanet(interface{}) error
}

func NewHandler() *Handler {
	return new(Handler)
}

type Mongo struct {
	Ctx         *context.Context
	CancelSinal chan struct{}
	Cancel      context.CancelFunc
	Collection  *mongo.Collection
}

func NewMongo(ctx context.Context, db *mongo.Database, collectionName string) *Mongo {
	ctx, cancel := context.WithCancel(ctx)
	cancelSinal := make(chan struct{}, 1)
	return &Mongo{&ctx, cancelSinal, cancel, db.Collection(collectionName)}
}

func (m *Mongo) Update(planets swapi.SwapiPlanetResponse) error {
	return nil
}

func (m *Mongo) Generate() error {
	return nil
}

func (m *Mongo) AddPlanet(interface{}) error {
	// TODO make middleware to send the write type
	timeOfMovies, err := GetTimeOnMovies()
	return nil
}

func GetTimeOnMovies(id int) (int, error) {
	SwapiPlanetResponse, err := swapi.GetPlanet()
	if err != nil {
		return 0, err
	}
	return len(SwapiPlanetResponse.GetFilms()), nil
}

func (m *Mongo) FindById(id int) (*Planet, error) {
	planet := new(Planet)
	if err := m.Collection.FindOne(*m.Ctx, bson.M{"id": id}).Decode(&planet); err != nil {
		return nil, err
	}
	return planet, nil
}

func (m *Mongo) FindByName(name int) (*Planet, error) {
	planet := new(Planet)
	if err := m.Collection.FindOne(*m.Ctx, bson.M{"name": name}).Decode(&planet); err != nil {
		return nil, err
	}
	return planet, nil
}

func (m *Mongo) WaitToCancelContext() {
	<-m.CancelSinal
	m.Cancel()
}
