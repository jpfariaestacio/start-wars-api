package db

import (
	"context"

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
	GetTimeOnMovies(id int) (int, error)
}

func NewUtils(dbv uint, db *mongo.Database, collectionName string) Utils {
	ctx := context.Background()
	switch dbv {
	case version1:
		return v1.NewMongo(ctx, db, collectionName)
	}
	return nil
}
