package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Config sets the credentials to access the MongoDB server
type Config struct {

	// DBHost is host of database
	DBHost string
	// DBUser is username to connect to database
	DBUser string
	// DBPassword password to connect to database
	DBPassword string
	// DBPort is the port to accessed
	DBPort string
	// DBName is the name of the database to be accessed
	DBName string
	// If true initialize the MongoDB database with default data and schema
	DBGenerate bool
	// If true updates values in the MongoDB by comparing with those in https://swapi.co/
	DBUpdate bool
}

// Connect uses the given config and context to connect to the MongoDB server
func (cfg Config) Connect(ctx context.Context) (*mongo.Database, error) {
	uri := cfg.makeURI()

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("todo: couldn't connect to mongo: %v", err)
	}
	err = client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("todo: mongo client couldn't connect with background context: %v", err)
	}

	starWarsDB := client.Database(cfg.DBName)

	return starWarsDB, nil
}

// makeURI returns the URI for connecting to MongoDB using the given config
func (cfg Config) makeURI() string {
	if cfg.DBUser != "" && cfg.DBPassword != "" {
		return fmt.Sprintf("mongodb://%s:%s@%s/%s",
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBHost,
			cfg.DBPort,
		)
	}
	return fmt.Sprintf("mongodb://%s:%s",
		cfg.DBHost,
		cfg.DBPort,
	)
}
