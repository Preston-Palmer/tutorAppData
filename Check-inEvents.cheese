package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

type Date struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	Date      time.Time `json:"date,omitempty" bson:"date,omitempty"`
	Tutor     string    `json:"tutor,omitempty" bson:"tutor,omitempty"`
	Subject   string    `json:"subject,omitempty" bson:"subject,omitempty"`
	Notes     string    `json:"notes,omitempty" bson:"notes,omitempty"`
	Completed bool      `json:"completed,omitempty" bson:"completed,omitempty"`
}

func getDates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var dates []Date
	collection := client.Database("date").Collection("dates")
	cur, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var date Date
		err := cur.Decode(&date)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		dates = append(dates, date)
	}
	if err := cur.Err(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(dates)
}

func getDate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var date Date
	collection := client.Database("date").Collection("dates")
	params := mux.Vars(r)
	cur := collection.FindOne(context.Background(), bson.M{"_id": params["id"]})
	err := cur.Decode(&date)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(date)
}
