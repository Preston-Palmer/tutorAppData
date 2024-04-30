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

type Date struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	Date      time.Time `json:"date,omitempty" bson:"date,omitempty"`
	Tutor     string    `json:"tutor,omitempty" bson:"tutor,omitempty"`
	Student   string    `json:"student,omitempty" bson:"student,omitempty"`
	Subject   string    `json:"subject,omitempty" bson:"subject,omitempty"`
	Notes     string    `json:"notes,omitempty" bson:"notes,omitempty"`
	Completed bool      `json:"completed,omitempty" bson:"completed,omitempty"`
}

func getDates(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting dates...")
	// if !validdate(r) {
	// 	http.Error(w, "unauthorized", http.StatusUnauthorized)
	// 	return
	// }
	coll := client.Database(viper.GetString("mongo.db")).Collection("dates")
	cursor, err := coll.Find(r.Context(), bson.M{})
	if err != nil {
		http.Error(w, "could not find dates: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())
	dates := []*Date{}
	for cursor.Next(r.Context()) {
		date := &Date{}
		err := cursor.Decode(date)
		if err != nil {
			http.Error(w, "could not decode date: "+err.Error(), http.StatusInternalServerError)
			return
		}
		dates = append(dates, date)
	}

	err = cursor.Err()
	if err != nil {
		http.Error(w, "cursor error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(dates)
	if err != nil {
		http.Error(w, "could not encode dates: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func getDate(w http.ResponseWriter, r *http.Request) {
	// if !validdate(r) {
	// 	http.Error(w, "unauthorized", http.StatusUnauthorized)
	// 	return
	// }
	log.Println("Getting date...")
	params := mux.Vars(r)
	dateID := params["dateID"]
	if dateID == "" {
		log.Println("dateID is required")
		http.Error(w, "dateID is required", http.StatusBadRequest)
		return
	}

	coll := client.Database(viper.GetString("mongo.db")).Collection("dates")
	date := &Date{}
	err := coll.FindOne(r.Context(), bson.M{"_id": dateID}).Decode(date)
	if err != nil {
		log.Printf("could not find date: %s\n", err)
		http.Error(w, "could not find date: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(date)
	if err != nil {
		log.Printf("could not encode date: %s\n", err)
		http.Error(w, "could not encode date: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func createDate(w http.ResponseWriter, r *http.Request) {
	log.Println("Creating date...")
	date := &Date{}
	err := json.NewDecoder(r.Body).Decode(date)
	if err != nil {
		http.Error(w, "could not decode date: "+err.Error(), http.StatusBadRequest)
		return
	}
	coll := client.Database(viper.GetString("mongo.db")).Collection("dates")
	_, err = coll.InsertOne(r.Context(), date)
	if err != nil {
		http.Error(w, "could not create date: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(date)
	if err != nil {
		http.Error(w, "could not encode date: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func updateDate(w http.ResponseWriter, r *http.Request) {
	// if !validdate(r) {
	// 	http.Error(w, "unauthorized", http.StatusUnauthorized)
	// 	return
	// }
	log.Println("Updating date...")
	params := mux.Vars(r)
	dateID := params["dateID"]
	if dateID == "" {
		log.Println("dateID is required")
		http.Error(w, "dateID is required", http.StatusBadRequest)
		return
	}

	date := &Date{}
	err := json.NewDecoder(r.Body).Decode(date)
	if err != nil {
		http.Error(w, "could not decode date: "+err.Error(), http.StatusInternalServerError)
		return
	}

	coll := client.Database(viper.GetString("mongo.db")).Collection("dates")
	_, err = coll.ReplaceOne(r.Context(), bson.M{"_id": dateID}, date)
	if err != nil {
		http.Error(w, "could not update date: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(date)
	if err != nil {
		http.Error(w, "could not encode date: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func deleteDate(w http.ResponseWriter, r *http.Request) {
	// if !validdate(r) {
	// 	http.Error(w, "unauthorized", http.StatusUnauthorized)
	// 	return
	// }
	log.Println("Deleting date...")
	params := mux.Vars(r)
	dateID := params["dateID"]
	if dateID == "" {
		log.Println("dateID is required")
		http.Error(w, "dateID is required", http.StatusBadRequest)
		return
	}

	coll := client.Database(viper.GetString("mongo.db")).Collection("dates")
	_, err := coll.DeleteOne(r.Context(), bson.M{"_id": dateID})
	if err != nil {
		http.Error(w, "could not delete date: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
