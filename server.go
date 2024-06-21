package main

import (
	"log"
	"myapp/database"
	"myapp/routes"

	"github.com/joho/godotenv"
)

func init() {
    var err error
    err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	err = database.InitDataBase()
    if err != nil {
        log.Fatal(err)
    }
    log.Println("DataBase Connected Successfuly.")
	err = routes.InitRoutes()
    if err != nil {
        log.Fatal(err)
    }
}

func main() {
	routes.Run()
}
