package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/nicolasassi/starWarsApi/pkg/cmd/db"
	"github.com/nicolasassi/starWarsApi/pkg/swapi"
)

var cfg db.Config

func init() {
	flag.StringVar(&cfg.DBHost, "db-host", "localhost", "Database host")
	flag.StringVar(&cfg.DBUser, "db-user", "", "Database user")
	flag.StringVar(&cfg.DBPassword, "db-password", "", "Database password")
	flag.StringVar(&cfg.DBPort, "db-port", "27017", "Database port")
	flag.StringVar(&cfg.DBName, "db-name", "starWarsDB", "Database name")
	flag.BoolVar(&cfg.DBGenerate, "db-generate", false, "// If true initialize the MongoDB database with default data and schema")
	flag.BoolVar(&cfg.DBUpdate, "db-update", false, "If true updates values in the MongoDB by comparing with those in https://swapi.co/")
	flag.Parse()
}

func init() {
	var userResp string
	if cfg.DBGenerate {
		fmt.Println(fmt.Sprintf("db-generate is set to TRUE. If there is a database identified as %s in your server it will be dumped and a new database with the same name and default values will be created. Are you sure? [y/n]", cfg.DBName))
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
	if cfg.DBUpdate {
		fmt.Println(fmt.Sprintf("db-update is set to TRUE. It might overwrite values set by your users in %s. Are you sure? [y/n]", cfg.DBName))
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
	tp, err := swapi.GetTotalPages()
	if err != nil {
		log.Fatal(err)
	}
	allValues := swapi.RetiveAllPlanets(tp)
	for _, v := range allValues {
		for _, v1 := range v.GetResults() {
			fmt.Println(v1)
		}
	}
}
