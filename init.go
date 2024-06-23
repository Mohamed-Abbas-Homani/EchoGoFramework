package main

import (
	"github.com/joho/godotenv"
	"log"
	"myapp/cache"
	"myapp/database"
	"myapp/routes"
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
	log.Println("DataBase Connected Successfully.")
	err = cache.InitCache()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Cache Connected Successfully.")
	err = routes.InitRoutes()
	if err != nil {
		log.Fatal(err)
	}
}
