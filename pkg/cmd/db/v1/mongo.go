package mongov1

import (
	"context"
	"encoding/json"
	"errors"
	fmt "fmt"
	"log"
	"strings"
	"sync"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	"github.com/nicolasassi/starWarsApi/pkg/swapi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var once sync.Once

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
func (m *Mongo) Update(responses []*swapi.SwapiResponse) error {
	errorChan := make(chan error)
	var wg sync.WaitGroup
	doneChan := make(chan struct{})
	go func() {
		for _, planets := range responses {
			for _, planet := range planets.GetResults() {
				wg.Add(1)
				once.Do(func() {
					go func() {
						wg.Wait()
						doneChan <- struct{}{}
					}()
				})
				go func(swapiPlanet *swapi.SwapiPlanetResponse) {
					defer wg.Done()
					planet, err := swapiToPlanetProto(*swapiPlanet)
					if err != nil {
						errorChan <- err
					}
					if err := m.updatePlanet(planet); err != nil {
						errorChan <- err
					}
					log.Println(planet.GetName())
				}(planet)
			}
		}
	}()

	select {
	case err := <-errorChan:
		return err
	case <-doneChan:
		return nil
	}

}

func swapiToPlanetProto(swapiPlanet swapi.SwapiPlanetResponse) (Planet, error) {
	planet := Planet{}
	intID, err := swapiPlanet.GetIDFromURL(swapiPlanet.GetUrl())
	if err != nil {
		return planet, err
	}
	swapiPlanet.Id = intID
	b, err := json.Marshal(swapiPlanet)
	if err != nil {
		return planet, err
	}
	jsonpb.UnmarshalString(string(b), &planet)
	planet.Terrain = strings.Split(strings.TrimSpace(swapiPlanet.GetTerrain()), ",")
	return planet, nil
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
	_, err := m.FindByName(planet.GetName())
	if err == nil {
		return errors.New("the planet already exist in the database")
	}
	timeOfMovies, err := m.GetTimeOnMovies(int(planet.GetId()))
	if err != nil {
		return err
	}
	planet.TimesOnMovies = int32(timeOfMovies)
	planet.AddedAt = ptypes.TimestampNow()
	if _, err := m.Collection.InsertOne(*m.Ctx, planet); err != nil {
		return err
	}
	return nil
}

func (m *Mongo) updatePlanet(planet Planet) error {
	p, err := m.FindByName(planet.GetName())
	if err != nil {
		return fmt.Errorf("could not update the planet %s", planet.GetName())
	}
	timeOfMovies, err := m.GetTimeOnMovies(int(planet.GetId()))
	if err != nil {
		return err
	}
	planet.TimesOnMovies = int32(timeOfMovies)
	planet.AddedAt = ptypes.TimestampNow()
	m.Collection.FindOneAndUpdate(*m.Ctx, bson.M{"name": p.GetName()}, bson.M{"$set": planet})
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
