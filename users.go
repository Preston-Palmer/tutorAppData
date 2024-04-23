package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string   `json:"id" bson:"_id"`
	Name     string   `json:"name"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	Scope    []string `json:"scope"`
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting users...")
	if !validUser(r) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	coll := client.Database(viper.GetString("mongo.db")).Collection("users")
	cursor, err := coll.Find(r.Context(), bson.M{})
	if err != nil {
		http.Error(w, "could not find users: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())
	users := []*User{}
	for cursor.Next(r.Context()) {
		user := &User{}
		err := cursor.Decode(user)
		if err != nil {
			http.Error(w, "could not decode user: "+err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	err = cursor.Err()
	if err != nil {
		http.Error(w, "cursor error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, "could not encode users: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	if !validUser(r) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	log.Println("Getting user...")
	params := mux.Vars(r)
	userID := params["userID"]
	if userID == "" {
		log.Println("userID is required")
		http.Error(w, "userID is required", http.StatusBadRequest)
		return
	}

	coll := client.Database(viper.GetString("mongo.db")).Collection("users")
	user := &User{}
	err := coll.FindOne(r.Context(), bson.M{"_id": userID}).Decode(user)
	if err != nil {
		log.Printf("could not find user: %s\n", err)
		http.Error(w, "could not find user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Printf("could not encode user: %s\n", err)
		http.Error(w, "could not encode user: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	if !validUser(r) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	log.Println("Creating user...")
	user := &User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		log.Println("decode error")
		http.Error(w, "could not decode user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("bcrypt error")
		http.Error(w, "could not bcrypt password: "+err.Error(), http.StatusInternalServerError)
		return
	}
	user.Password = string(bcryptPassword)

	coll := client.Database(viper.GetString("mongo.db")).Collection("users")
	_, err = coll.InsertOne(r.Context(), user)
	if err != nil {
		log.Println("insert error")
		http.Error(w, "could not insert user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Println("encode error")
		http.Error(w, "could not encode user: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	if !validUser(r) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	log.Println("Updating user...")
	params := mux.Vars(r)
	userID := params["userID"]
	if userID == "" {
		log.Println("userID is required")
		http.Error(w, "userID is required", http.StatusBadRequest)
		return
	}

	user := &User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		http.Error(w, "could not decode user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	coll := client.Database(viper.GetString("mongo.db")).Collection("users")
	_, err = coll.ReplaceOne(r.Context(), bson.M{"_id": userID}, user)
	if err != nil {
		http.Error(w, "could not update user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "could not encode user: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	if !validUser(r) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	log.Println("Deleting user...")
	params := mux.Vars(r)
	userID := params["userID"]
	if userID == "" {
		log.Println("userID is required")
		http.Error(w, "userID is required", http.StatusBadRequest)
		return
	}

	coll := client.Database(viper.GetString("mongo.db")).Collection("users")
	_, err := coll.DeleteOne(r.Context(), bson.M{"_id": userID})
	if err != nil {
		http.Error(w, "could not delete user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
