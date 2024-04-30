package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
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

type Tutor struct {
	ID           string `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName    string `json:"fname,omitempty" bson:"fname,omitempty"`
	LastName     string `json:"lname,omitempty" bson:"lname,omitempty"`
	Subject      string `json:"subject,omitempty" bson:"subject,omitempty"`
	Availability string `json:"availability,omitempty" bson:"availability,omitempty"`
	// Students     []Students `json:"students,omitempty" bson:"students,omitempty"`
	// Image        bytea      `json:"image,omitempty" bson:"image,omitempty"`
	// CheckIn bool   `json:"checkin,omitempty" bson:"checkin,omitempty"`
	// Dates []Date `json:"dates,omitempty" bson:"dates,omitempty"`
}

func getTutors(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting tutors...")
	// if !validtutor(r) {
	// 	http.Error(w, "unauthorized", http.StatusUnauthorized)
	// 	return
	// }
	coll := client.Database(viper.GetString("mongo.db")).Collection("tutors")
	cursor, err := coll.Find(r.Context(), bson.M{})
	if err != nil {
		http.Error(w, "could not find tutors: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())
	tutors := []*Tutor{}
	for cursor.Next(r.Context()) {
		tutor := &Tutor{}
		err := cursor.Decode(tutor)
		if err != nil {
			http.Error(w, "could not decode tutor: "+err.Error(), http.StatusInternalServerError)
			return
		}
		tutors = append(tutors, tutor)
	}

	err = cursor.Err()
	if err != nil {
		http.Error(w, "cursor error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tutors)
	if err != nil {
		http.Error(w, "could not encode tutors: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func getTutor(w http.ResponseWriter, r *http.Request) {
	// if !validtutor(r) {
	// 	http.Error(w, "unauthorized", http.StatusUnauthorized)
	// 	return
	// }
	log.Println("Getting tutor...")
	params := mux.Vars(r)
	tutorID := params["tutorID"]
	if tutorID == "" {
		log.Println("tutorID is required")
		http.Error(w, "tutorID is required", http.StatusBadRequest)
		return
	}

	coll := client.Database(viper.GetString("mongo.db")).Collection("tutors")
	tutor := &Tutor{}
	err := coll.FindOne(r.Context(), bson.M{"_id": tutorID}).Decode(tutor)
	if err != nil {
		log.Printf("could not find tutor: %s\n", err)
		http.Error(w, "could not find tutor: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tutor)
	if err != nil {
		log.Printf("could not encode tutor: %s\n", err)
		http.Error(w, "could not encode tutor: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func createTutor(w http.ResponseWriter, r *http.Request) {
	log.Println("Creating tutor...")
	tutor := &Tutor{}
	err := json.NewDecoder(r.Body).Decode(tutor)
	if err != nil {
		http.Error(w, "could not decode tutor: "+err.Error(), http.StatusBadRequest)
		return
	}
	coll := client.Database(viper.GetString("mongo.db")).Collection("tutors")
	_, err = coll.InsertOne(r.Context(), tutor)
	if err != nil {
		http.Error(w, "could not create tutor: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tutor)
	if err != nil {
		http.Error(w, "could not encode tutor: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func updateTutor(w http.ResponseWriter, r *http.Request) {
	// if !validtutor(r) {
	// 	http.Error(w, "unauthorized", http.StatusUnauthorized)
	// 	return
	// }
	log.Println("Updating tutor...")
	params := mux.Vars(r)
	tutorID := params["tutorID"]
	if tutorID == "" {
		log.Println("tutorID is required")
		http.Error(w, "tutorID is required", http.StatusBadRequest)
		return
	}

	tutor := &Tutor{}
	err := json.NewDecoder(r.Body).Decode(tutor)
	if err != nil {
		http.Error(w, "could not decode tutor: "+err.Error(), http.StatusInternalServerError)
		return
	}

	coll := client.Database(viper.GetString("mongo.db")).Collection("tutors")
	_, err = coll.ReplaceOne(r.Context(), bson.M{"_id": tutorID}, tutor)
	if err != nil {
		http.Error(w, "could not update tutor: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tutor)
	if err != nil {
		http.Error(w, "could not encode tutor: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func deleteTutor(w http.ResponseWriter, r *http.Request) {
	// if !validtutor(r) {
	// 	http.Error(w, "unauthorized", http.StatusUnauthorized)
	// 	return
	// }
	log.Println("Deleting tutor...")
	params := mux.Vars(r)
	tutorID := params["tutorID"]
	if tutorID == "" {
		log.Println("tutorID is required")
		http.Error(w, "tutorID is required", http.StatusBadRequest)
		return
	}

	coll := client.Database(viper.GetString("mongo.db")).Collection("tutors")
	_, err := coll.DeleteOne(r.Context(), bson.M{"_id": tutorID})
	if err != nil {
		http.Error(w, "could not delete tutor: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
