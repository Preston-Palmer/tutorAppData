package main

import (
	"context"
	"crypto/rsa"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var privKey *rsa.PrivateKey
var pubKey *rsa.PublicKey

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	viper.SetConfigName("settings")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("could not read config: %s\n", err)
	}

	var err error
	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/api/students", getStudents).Methods("GET")
	router.HandleFunc("/api/students/{id}", getStudent).Methods("GET")
	router.HandleFunc("/api/students", createStudent).Methods("POST")
	router.HandleFunc("/api/students/{id}", updateStudent).Methods("PUT")
	router.HandleFunc("/api/students/{id}", deleteStudent).Methods("DELETE")

	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
