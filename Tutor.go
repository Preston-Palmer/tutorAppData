package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

// type Date struct {
// 	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
// 	Date      time.Time `json:"date,omitempty" bson:"date,omitempty"`
// 	Tutor     string    `json:"tutor,omitempty" bson:"tutor,omitempty"`
// 	Subject   string    `json:"subject,omitempty" bson:"subject,omitempty"`
// 	Notes     string    `json:"notes,omitempty" bson:"notes,omitempty"`
// 	Completed bool      `json:"completed,omitempty" bson:"completed,omitempty"`
// }

type Tutors struct {
	ID           string     `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName    string     `json:"firstname,omitempty" bson:"firstname,omitempty"`
	LastName     string     `json:"lastname,omitempty" bson:"lastname,omitempty"`
	Subjects     string     `json:"subjects,omitempty" bson:"subjects,omitempty"`
	Availability string     `json:"availability,omitempty" bson:"availability,omitempty"`
	Students     []Students `json:"students,omitempty" bson:"students,omitempty"`
	Image        bytea      `json:"image,omitempty" bson:"image,omitempty"`
	CheckIn      bool       `json:"checkin,omitempty" bson:"checkin,omitempty"`
	Dates        []Date     `json:"dates,omitempty" bson:"dates,omitempty"`
}

func gettutors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var tutors []Tutors
	collection := client.Database("tutor").Collection("tutors")
	cur, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var tutor Tutors
		err := cur.Decode(&tutor)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		tutors = append(tutors, tutor)
	}
	if err := cur.Err(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tutors)
}

func gettutor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var tutor Tutors
	collection := client.Database("tutor").Collection("tutors")
	params := mux.Vars(r)
	cur := collection.FindOne(context.Background(), bson.M{"_id": params["id"]})
	err := cur.Decode(&tutor)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tutor)
}

func createtutor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var tutor Tutors
	collection := client.Database("tutor").Collection("tutors")
	err := json.NewDecoder(r.Body).Decode(&tutor)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tutor.ID = uuid.New().String()
	_, err = collection.InsertOne(context.Background(), tutor)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tutor)
}

func updatetutor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var tutor Tutors
	collection := client.Database("tutor").Collection("tutors")
	params := mux.Vars(r)
	err := json.NewDecoder(r.Body).Decode(&tutor)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": params["id"]}, bson.M{"$set": tutor})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tutor)
}

func deletetutor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	collection := client.Database("tutor").Collection("tutors")
	params := mux.Vars(r)
	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": params["id"]})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
