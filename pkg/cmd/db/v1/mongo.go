package v1

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	"github.com/nicolasassi/starWarsApi/pkg/swapi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Handler abstracts the methods to interact with a mongo collection
type Handler interface {
	AddPlanet(interface{}) error
	ListPlanets() ([]Planet, error)
	FindByName(name string) (Planet, error)
	FindById(id int) (Planet, error)
	RemovePlanet(id int) error
}

// NewHandler returns a Handler to abstract the methods to interact with the mongo collection
func NewHandler() *Handler {
	return new(Handler)
}

// Mongo control the behavior of the Mongo collection and interaction
type Mongo struct {
	Ctx         *context.Context
	CancelSinal chan struct{}
	Cancel      context.CancelFunc
	Collection  *mongo.Collection
}

// NewMongo creates a new Mongo intance of the version 1
func NewMongo(ctx context.Context, db *mongo.Database, collectionName string) *Mongo {
	ctx, cancel := context.WithCancel(ctx)
	cancelSinal := make(chan struct{}, 1)
	return &Mongo{&ctx, cancelSinal, cancel, db.Collection(collectionName)}
}

// Update retrives all planets from swapi/planets and update its values in the given mongo collection
func (m *Mongo) Update(responses []swapi.SwapiResponse) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Make sure it's called to release resources even if no errors
	errorChan := make(chan error)
	doneChan := make(chan struct{})
	go func() {
		defer ctx.Done()
		for _, planets := range responses {
			for _, planet := range planets.GetResults() {
				go func(swapiPlanet *swapi.SwapiPlanetResponse) {
					planet, err := swapiToPlanetProto(*swapiPlanet)
					if err != nil {
						errorChan <- err
					}
					if err := m.AddPlanet(planet); err != nil {
						errorChan <- err
					}
				}(planet)
			}
		}
	}()
	for {
		select {
		case err := <-errorChan:
			return err
		case <-doneChan:
			return nil
		}
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

// Generate creates a new database by using the given database name and default data
// TODO
func (m *Mongo) Generate() error {
	return nil
}

// DeletePlanet searches if a planet exists in the mongodb collection by its id and if so
// deletes it
func (m *Mongo) DeletePlanet(id int) error {
	var planet Planet
	if err := m.Collection.FindOneAndDelete(*m.Ctx, bson.M{"id": id}).Decode(&planet); err != nil {
		return err
	}
	return nil
}

// AddPlanet checks if the planet data already is a document in the mongodb
// by checking the planet name
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

// GetTimeOnMovies requests the data from swapi/planets and parses the times on movies of
// the given planet returing zero if the planet is not part of the swapi/planets dataset
func (m *Mongo) GetTimeOnMovies(id int) (int, error) {
	client := swapi.NewSwapiClient("planets")
	resp, err := client.RetriveItem(id)
	if err != nil {
		return 0, err
	}
	switch resp.(type) {
	case *swapi.SwapiPlanetResponse:
		return len(resp.(*swapi.SwapiPlanetResponse).GetFilms()), nil
	default:
		return 0, errors.New("error converting the response from swapi to proto SwapiPlanetResponse")
	}

	return 0, nil
}

// FindByID returns the planet information by searching it by id in the mongodb database
func (m *Mongo) FindByID(id int) (*Planet, error) {
	planet := new(Planet)
	if err := m.Collection.FindOne(*m.Ctx, bson.M{"id": id}).Decode(&planet); err != nil {
		return nil, err
	}
	return planet, nil
}

// FindByName returns the planet information by searching it by name in the mongodb database
func (m *Mongo) FindByName(name string) (*Planet, error) {
	planet := new(Planet)
	if err := m.Collection.FindOne(*m.Ctx, bson.M{"name": name}).Decode(&planet); err != nil {
		return nil, err
	}
	return planet, nil
}

// WaitToCancelContext just to remember after
// TODO will became a select statmemnt
func (m *Mongo) WaitToCancelContext() {
	<-m.CancelSinal
	m.Cancel()
}
