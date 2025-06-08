package main

import (
	"kairon/cmd/api/infrastructure/datastore"
	"kairon/cmd/api/infrastructure/router"
	"kairon/config"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("API :: setup")

	config.ReadConf()

	db, err := datastore.NewDBConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fbConfig := firebase.Config{
		ProjectID: config.C.ProjectID,
	}
	app, err := firebase.NewApp(db.Ctx, &fbConfig)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	authClient, err := app.Auth(db.Ctx)
	if err != nil {
		log.Fatalf("error getting Firebase client: %v\n", err)
	}

	server := router.NewServer(db, authClient)
	server.Run(config.C.Server.Address)
}
