package database

import (
	"log"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() {
	err := mgm.SetDefaultConfig(nil, "studynest", options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("MongoDB connect aagilla", err)
	}

	log.Println("MongoDB connect aagide")
}
