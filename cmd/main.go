package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nicolasassi/starWarsApi/pkg/cmd/db"
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
	flag.BoolVar(&cfg.DBGenerate, "db-generate", false, "// If true initialize the MongoDB database with default data and schema")
	flag.BoolVar(&cfg.DBUpdate, "db-update", false, "If true updates values in the MongoDB by comparing with those in https://swapi.co/")
	flag.Parse()
}

func init() {
	if cfg.DBGenerate {
		warningMessageRaw, ok := os.LookupEnv("DB-GENERATE-WARNING-MESSAGE")
		if !ok {
			log.Fatal("error in your .ENV file")
		}
		warningMessage := strings.Split(warningMessageRaw, "%s")
		fmt.Println(warningMessage[0], cfg.DBName, warningMessage[1])
		fmt.Scan(&userResp)
		switch userResp {
		case "y":
		case "n":
			cfg.DBGenerate = false
			fmt.Println("db-generate set to FALSE")
		default:
			fmt.Println("Invalid input " + userResp)
			os.Exit(0)
		}
	}
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
	// ctx := context.Background()
	// db, err := cfg.Connect(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	sc := swapi.NewSwapiClient("planets")
	log.Println(sc.RetriveAll())
	// allValues := swapi.RetriveAllPlanets(tp)
	// for _, v := range allValues {
	// 	for _, v1 := range v.GetResults() {
	// 		fmt.Println(v1)
	// 	}
	// }
}
