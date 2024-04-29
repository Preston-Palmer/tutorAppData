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
	w.Header().Set("Content-Type", "application/json")

	var students []Student
	collection := client.Database("student").Collection("students")
	cur, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var student Student
		err := cur.Decode(&student)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		students = append(students, student)
	}
	if err := cur.Err(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(students)
}

func getStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var student Student
	collection := client.Database("student").Collection("students")
	params := mux.Vars(r)
	cur := collection.FindOne(context.Background(), bson.M{"_id": params["id"]})
	err := cur.Decode(&student)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(student)
}

func createStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var student Student
	collection := client.Database("student").Collection("students")
	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	student.ID = uuid.New().String()
	_, err = collection.InsertOne(context.Background(), student)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(student)
}

func updateStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var student Student
	collection := client.Database("student").Collection("students")
	params := mux.Vars(r)
	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": params["id"]}, bson.M{"$set": student})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(student)
}

func deleteStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	collection := client.Database("student").Collection("students")
	params := mux.Vars(r)
	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": params["id"]})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
