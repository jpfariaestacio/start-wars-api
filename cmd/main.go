package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	server "github.com/nicolasassi/starWarsApi/pkg/api/v1"
	"github.com/nicolasassi/starWarsApi/pkg/cmd/db"
	mongov1 "github.com/nicolasassi/starWarsApi/pkg/cmd/db/v1"
	"github.com/nicolasassi/starWarsApi/pkg/swapi"
)

var (
	cfg      db.Config
	userResp string
)

func init() {
	f, err := os.Open("cmd/.env")
	if err != nil {
		log.Fatal("error while opening your .env file")
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		rawKeyValuePair := strings.Split(scanner.Text(), "=")
		os.Setenv(rawKeyValuePair[0], rawKeyValuePair[1])
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	flag.StringVar(&cfg.DBHost, "db-host", "localhost", "Database host")
	flag.StringVar(&cfg.DBUser, "db-user", "", "Database user")
	flag.StringVar(&cfg.DBPassword, "db-password", "", "Database password")
	flag.StringVar(&cfg.DBPort, "db-port", "27017", "Database port")
	flag.StringVar(&cfg.DBName, "db-name", "star_wars_db_v1", "Database name")
	flag.BoolVar(&cfg.DBUpdate, "db-update", false, "If true updates values in the MongoDB by comparing with those in https://swapi.co/")
	flag.Parse()
}

func init() {
	if cfg.DBUpdate {
		warningMessageRaw, ok := os.LookupEnv("DB-UPDATE-WARNING-MESSAGE")
		if !ok {
			log.Fatal("error in your .ENV file")
		}
		warningMessage := strings.Split(warningMessageRaw, "%s")
		fmt.Println(warningMessage[0], cfg.DBName, warningMessage[1])
		fmt.Scan(&userResp)
		switch userResp {
		case "y":
			ctx := context.Background()
			defer ctx.Done()
			dbm, err := cfg.Connect(ctx)
			if err != nil {
				log.Fatal(err)
			}
			sc := swapi.NewSwapiClient("planets")
			utils := db.NewUtils(1, dbm, "planets")
			all, err := sc.RetriveAll()
			if err != nil {
				log.Fatal(err)
			}
			var responses []*swapi.SwapiResponse
			for _, obj := range all {
				if resp, ok := obj.(*swapi.SwapiResponse); ok {
					responses = append(responses, resp)
				}
			}
			if err := utils.Update(responses); err != nil {
				log.Fatal(err)
			}
		case "n":
			cfg.DBUpdate = false
			fmt.Println("db-update set to FALSE")
		default:
			fmt.Println("Invalid input " + userResp)
			os.Exit(0)
		}

	}
}

func main() {

	ctx := context.Background()
	defer ctx.Done()
	dbm, err := cfg.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	m := mongov1.NewMongo(ctx, dbm, "planets")
	server.Serve(m)
}
