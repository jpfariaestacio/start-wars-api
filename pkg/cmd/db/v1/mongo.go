package v1

import (
	"context"
	"encoding/json"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	"github.com/nicolasassi/starWarsApi/pkg/swapi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler interface {
	AddPlanet(interface{}) error
	ListPlanets() ([]Planet, error)
	FindByName(name string) (Planet, error)
	FindById(id int) (Planet, error)
	RemovePlanet(id int) error
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

func (m *Mongo) Update(planets swapi.SwapiResponse) error {
	for _, planet := range planets.GetResults() {
		go func(swapiPlanet *swapi.SwapiPlanetResponse) {
			planet, err := swapiToPlanetProto(*swapiPlanet)
			if err != nil {

			}
			if err := m.AddPlanet(planet); err != nil {

			}
		}(planet)
		return nil
	}
	return nil
}

func swapiToPlanetProto(swapiPlanet swapi.SwapiPlanetResponse) (Planet, error) {
	planet := Planet{}
	b, err := json.Marshal(swapiPlanet)
	if err != nil {
		return planet, err
	}
	jsonpb.UnmarshalString(string(b), &planet)
	return planet, nil
}

func (m *Mongo) Generate() error {
	return nil
}

func (m *Mongo) DeletePlanet(id int) error {
	var planet Planet
	if err := m.Collection.FindOneAndDelete(*m.Ctx, bson.M{"id": id}).Decode(&planet); err != nil {
		return err
	}
	return nil
}

func (m *Mongo) AddPlanet(planet Planet) error {
	p, err := m.FindByName(planet.GetName())
	if err != nil {
		return err
	}
	timeOfMovies, err := m.GetTimeOnMovies(int(planet.GetId()))
	if err != nil {
		return err
	}
	planet.TimesOnMovies = int32(timeOfMovies)
	planet.AddedAt = ptypes.TimestampNow()
	if p == new(Planet) {
		m.Collection.InsertOne(*m.Ctx, planet)
		return nil
	}
	if err := m.Collection.FindOneAndUpdate(*m.Ctx, bson.M{"name": p.GetName()}, planet).Decode(&p); err != nil {
		return err
	}
	return nil
}

func (m *Mongo) GetTimeOnMovies(id int) (int, error) {
	SwapiPlanetResponse, err := swapi.GetPlanet(id)
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

func (m *Mongo) FindByName(name string) (*Planet, error) {
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
