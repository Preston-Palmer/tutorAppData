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

// type Date struct {
// 	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
// 	Date      time.Time `json:"date,omitempty" bson:"date,omitempty"`
// 	Tutor     string    `json:"tutor,omitempty" bson:"tutor,omitempty"`
// 	Subject   string    `json:"subject,omitempty" bson:"subject,omitempty"`
// 	Notes     string    `json:"notes,omitempty" bson:"notes,omitempty"`
// 	Completed bool      `json:"completed,omitempty" bson:"completed,omitempty"`
// }

type Student struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName string    `json:"firstname,omitempty" bson:"firstname,omitempty"`
	LastName  string    `json:"lastname,omitempty" bson:"lastname,omitempty"`
	TimeIn    time.Time `json:"timein,omitempty" bson:"timein,omitempty"`
	TimeOut   time.Time `json:"timeout,omitempty" bson:"timeout,omitempty"`
	CheckIn   bool      `json:"checkin,omitempty" bson:"checkin,omitempty"`
	// Dates     []Date    `json:"dates,omitempty" bson:"dates,omitempty"`
}

func getStudents(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting students...")
	// if !validstudent(r) {
	// 	http.Error(w, "unauthorized", http.StatusUnauthorized)
	// 	return
	// }
	coll := client.Database(viper.GetString("mongo.db")).Collection("students")
	cursor, err := coll.Find(r.Context(), bson.M{})
	if err != nil {
		http.Error(w, "could not find students: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())
	students := []*Student{}
	for cursor.Next(r.Context()) {
		student := &Student{}
		err := cursor.Decode(student)
		if err != nil {
			http.Error(w, "could not decode student: "+err.Error(), http.StatusInternalServerError)
			return
		}
		students = append(students, student)
	}

	err = cursor.Err()
	if err != nil {
		http.Error(w, "cursor error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(students)
	if err != nil {
		http.Error(w, "could not encode students: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func getStudent(w http.ResponseWriter, r *http.Request) {
	// if !validstudent(r) {
	// 	http.Error(w, "unauthorized", http.StatusUnauthorized)
	// 	return
	// }
	log.Println("Getting student...")
	params := mux.Vars(r)
	studentID := params["studentID"]
	if studentID == "" {
		log.Println("studentID is required")
		http.Error(w, "studentID is required", http.StatusBadRequest)
		return
	}

	coll := client.Database(viper.GetString("mongo.db")).Collection("students")
	student := &Student{}
	err := coll.FindOne(r.Context(), bson.M{"_id": studentID}).Decode(student)
	if err != nil {
		log.Printf("could not find student: %s\n", err)
		http.Error(w, "could not find student: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(student)
	if err != nil {
		log.Printf("could not encode student: %s\n", err)
		http.Error(w, "could not encode student: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func createStudent(w http.ResponseWriter, r *http.Request) {
	log.Println("Creating student...")
	student := &Student{}
	err := json.NewDecoder(r.Body).Decode(student)
	if err != nil {
		http.Error(w, "could not decode student: "+err.Error(), http.StatusBadRequest)
		return
	}
	coll := client.Database(viper.GetString("mongo.db")).Collection("students")
	_, err = coll.InsertOne(r.Context(), student)
	if err != nil {
		http.Error(w, "could not create student: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(student)
	if err != nil {
		http.Error(w, "could not encode student: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func updateStudent(w http.ResponseWriter, r *http.Request) {
	// if !validstudent(r) {
	// 	http.Error(w, "unauthorized", http.StatusUnauthorized)
	// 	return
	// }
	log.Println("Updating student...")
	params := mux.Vars(r)
	studentID := params["studentID"]
	if studentID == "" {
		log.Println("studentID is required")
		http.Error(w, "studentID is required", http.StatusBadRequest)
		return
	}

	student := &Student{}
	err := json.NewDecoder(r.Body).Decode(student)
	if err != nil {
		http.Error(w, "could not decode student: "+err.Error(), http.StatusInternalServerError)
		return
	}

	coll := client.Database(viper.GetString("mongo.db")).Collection("students")
	_, err = coll.ReplaceOne(r.Context(), bson.M{"_id": studentID}, student)
	if err != nil {
		http.Error(w, "could not update student: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(student)
	if err != nil {
		http.Error(w, "could not encode student: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func deleteStudent(w http.ResponseWriter, r *http.Request) {
	// if !validstudent(r) {
	// 	http.Error(w, "unauthorized", http.StatusUnauthorized)
	// 	return
	// }
	log.Println("Deleting student...")
	params := mux.Vars(r)
	studentID := params["studentID"]
	if studentID == "" {
		log.Println("studentID is required")
		http.Error(w, "studentID is required", http.StatusBadRequest)
		return
	}

	coll := client.Database(viper.GetString("mongo.db")).Collection("students")
	_, err := coll.DeleteOne(r.Context(), bson.M{"_id": studentID})
	if err != nil {
		http.Error(w, "could not delete student: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
