package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

// NEED TO FOLLOW USERS.GO for Student.ID instead of _id
type Clock struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	Date      time.Time `json:"date,omitempty" bson:"date,omitempty"`
	Tutor     string    `json:"tutor,omitempty" bson:"tutor,omitempty"`
	Subject   string    `json:"subject,omitempty" bson:"subject,omitempty"`
	Notes     string    `json:"notes,omitempty" bson:"notes,omitempty"`
	Completed bool      `json:"completed,omitempty" bson:"completed,omitempty"`
}

func getClocks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var clocks []Clock
	collection := client.Database("clock").Collection("clocks")
	cur, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var clock Clock
		err := cur.Decode(&clock)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		clocks = append(clocks, clock)
	}
	if err := cur.Err(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(clocks)
}

func getClock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var clock Clock
	collection := client.Database("clock").Collection("clocks")
	params := mux.Vars(r)
	cur := collection.FindOne(context.Background(), bson.M{"_id": params["id"]})
	err := cur.Decode(&clock)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(clock)
}

func createClock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var clock Clock
	collection := client.Database("clock").Collection("clocks")
	err := json.NewDecoder(r.Body).Decode(&clock)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	clock.ID = uuid.New().String()
	_, err = collection.InsertOne(context.Background(), clock)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(clock)
}

func updateClock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var clock Clock
	collection := client.Database("clock").Collection("clocks")
	params := mux.Vars(r)
	err := json.NewDecoder(r.Body).Decode(&clock)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": params["id"]}, bson.M{"$set": clock})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(clock)
}

func deleteClock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	collection := client.Database("clock").Collection("clocks")
	params := mux.Vars(r)
	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": params["id"]})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
