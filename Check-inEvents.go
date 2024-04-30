package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
)

type Clock struct {
	ID      string    `json:"id,omitempty" bson:"_id,omitempty"`
	Date    time.Time `json:"date,omitempty" bson:"date,omitempty"`
	OutDate time.Time `json:"outdate,omitempty" bson:"outdate,omitempty"`
	FName   string    `json:"fname,omitempty" bson:"fname,omitempty"`
	LName   string    `json:"lname,omitempty" bson:"lname,omitempty"`
	Subject string    `json:"subject,omitempty" bson:"subject,omitempty"`
	Notes   string    `json:"notes,omitempty" bson:"notes,omitempty"`
	ClockIn bool      `json:"clockin,omitempty" bson:"clockin,omitempty"`
}

func getClocks(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting clocks...")
	// if !validclock(r) {
	// 	http.Error(w, "unauthorized", http.StatusUnauthorized)
	// 	return
	// }
	coll := client.Database(viper.GetString("mongo.db")).Collection("clocks")
	cursor, err := coll.Find(r.Context(), bson.M{})
	if err != nil {
		http.Error(w, "could not find clocks: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())
	clocks := []*Clock{}
	for cursor.Next(r.Context()) {
		clock := &Clock{}
		err := cursor.Decode(clock)
		if err != nil {
			http.Error(w, "could not decode clock: "+err.Error(), http.StatusInternalServerError)
			return
		}
		clocks = append(clocks, clock)
	}

	err = cursor.Err()
	if err != nil {
		http.Error(w, "cursor error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(clocks)
	if err != nil {
		http.Error(w, "could not encode clocks: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func getClock(w http.ResponseWriter, r *http.Request) {
	// if !validClock(r) {
	// 	http.Error(w, "unauthorized", http.StatusUnauthorized)
	// 	return
	// }
	log.Println("Getting clock...")
	params := mux.Vars(r)
	clockID := params["clockID"]
	if clockID == "" {
		log.Println("clockID is required")
		http.Error(w, "clockID is required", http.StatusBadRequest)
		return
	}

	coll := client.Database(viper.GetString("mongo.db")).Collection("clocks")
	clock := &Clock{}
	err := coll.FindOne(r.Context(), bson.M{"_id": clockID}).Decode(clock)
	if err != nil {
		log.Printf("could not find clock: %s\n", err)
		http.Error(w, "could not find clock: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(clock)
	if err != nil {
		log.Printf("could not encode clock: %s\n", err)
		http.Error(w, "could not encode clock: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func createClock(w http.ResponseWriter, r *http.Request) {
	log.Println("Creating clock...")
	clock := &Clock{}
	err := json.NewDecoder(r.Body).Decode(clock)
	if err != nil {
		http.Error(w, "could not decode clock: "+err.Error(), http.StatusBadRequest)
		return
	}
	coll := client.Database(viper.GetString("mongo.db")).Collection("clocks")
	_, err = coll.InsertOne(r.Context(), clock)
	if err != nil {
		http.Error(w, "could not create clock: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(clock)
	if err != nil {
		http.Error(w, "could not encode clock: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func updateClock(w http.ResponseWriter, r *http.Request) {
	// if !validclock(r) {
	// 	http.Error(w, "unauthorized", http.StatusUnauthorized)
	// 	return
	// }
	log.Println("Updating clock...")
	params := mux.Vars(r)
	clockID := params["clockID"]
	if clockID == "" {
		log.Println("clockID is required")
		http.Error(w, "clockID is required", http.StatusBadRequest)
		return
	}

	clock := &Clock{}
	err := json.NewDecoder(r.Body).Decode(clock)
	if err != nil {
		http.Error(w, "could not decode clock: "+err.Error(), http.StatusInternalServerError)
		return
	}

	coll := client.Database(viper.GetString("mongo.db")).Collection("clocks")
	_, err = coll.ReplaceOne(r.Context(), bson.M{"_id": clockID}, clock)
	if err != nil {
		http.Error(w, "could not update clock: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(clock)
	if err != nil {
		http.Error(w, "could not encode clock: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func deleteClock(w http.ResponseWriter, r *http.Request) {
	// if !validclock(r) {
	// 	http.Error(w, "unauthorized", http.StatusUnauthorized)
	// 	return
	// }
	log.Println("Deleting clock...")
	params := mux.Vars(r)
	clockID := params["clockID"]
	if clockID == "" {
		log.Println("clockID is required")
		http.Error(w, "clockID is required", http.StatusBadRequest)
		return
	}

	coll := client.Database(viper.GetString("mongo.db")).Collection("clocks")
	_, err := coll.DeleteOne(r.Context(), bson.M{"_id": clockID})
	if err != nil {
		http.Error(w, "could not delete clock: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
