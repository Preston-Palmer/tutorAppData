package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TodoClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	FName    string `json:"fname"`
	LName    string `json:"lname"`
	Scope    string `json:"scope"`
	jwt.RegisteredClaims
}

func getLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting login...")
	lr := &LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(lr)
	if err != nil {
		log.Printf("could not decode login request: %s\n", err)
		http.Error(w, "could not decode login request: "+err.Error(), http.StatusBadRequest)
		return
	}

	coll := client.Database(viper.GetString("mongo.db")).Collection("users")
	user := &User{}
	err = coll.FindOne(r.Context(), bson.M{"username": lr.Username}).Decode(user)
	if err != nil {
		log.Printf("could not find user: %s\n", err)
		http.Error(w, "could not find user: "+err.Error(), http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(lr.Password))
	if err != nil {
		log.Printf("could not compare passwords: %s\n", err)
		http.Error(w, "could not compare passwords: "+err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := createToken(user)
	if err != nil {
		log.Printf("could not create token: %s\n", err)
		http.Error(w, "could not create token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{"token": token})
	if err != nil {
		log.Printf("could not encode token: %s\n", err)
		http.Error(w, "could not encode token: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func createToken(user *User) (string, error) {
	claims := &TodoClaims{
		ID:       user.ID,
		Username: user.Username,
		FName:    user.FName,
		LName:    user.LName,
		Scope:    user.Scope,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "tutorDB",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	ss, err := token.SignedString(privKey)
	if err != nil {
		log.Printf("could not sign token: %s\n", err)
		return "", err
	}
	return ss, nil
}

func validUser(r *http.Request) bool {
	claims, err := getTokenClaims(r)
	if err != nil {
		log.Printf("could not get token claims: %s\n", err)
		return false
	}
	if claims == nil {
		return false
	}
	return true
}

func getTokenClaims(r *http.Request) (*TodoClaims, error) {
	authorization := r.Header.Get("Authorization")
	if !strings.HasPrefix(authorization, "Bearer ") {
		return nil, errors.New("invalid authorization header")
	}
	tokenString := strings.TrimPrefix(authorization, "Bearer ")
	claims := &TodoClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodRSA)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return pubKey, nil
	})
	if err != nil {
		return nil, errors.New("invalid token: " + err.Error())
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
