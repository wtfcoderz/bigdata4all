package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type jsonUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type basicResponse struct {
	Result string `json:"result"`
	Msg    string `json:"msg"`
}

var port string
var secretKey string
var redisDB *redis.Client

func init() {
	flag.StringVar(&port, "port", "80", "give me a port number")
}

func main() {
	flag.Parse()

	// Init DB
	redisDB = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	// Secret
	secretKey = string(securecookie.GenerateRandomKey(32))

	// Mux Router
	r := mux.NewRouter()
	r.HandleFunc("/user/register", handleUserRegister).Methods("Post")
	r.HandleFunc("/user/auth", handleUserAuth).Methods("Post")
	r.HandleFunc("/health", health)
	r.HandleFunc("/data", handleData).Methods("Post")
	r.HandleFunc("/", health)
	http.Handle("/", r)

	// Start Http
	fmt.Println("Starting up on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func health(w http.ResponseWriter, req *http.Request) {
	hostname, _ := os.Hostname()
	fmt.Fprintln(w, "Hostname:", hostname)
}

func handleUserRegister(w http.ResponseWriter, req *http.Request) {
	bodyb, _ := ioutil.ReadAll(req.Body)
	var user jsonUser
	_ = json.Unmarshal(bodyb, &user)

	var response basicResponse

	// Verify user
	_, err := redisDB.Get("user:" + user.Email).Result()
	if err == redis.Nil {
		// User not present in DB
		token := randStringBytesMaskImprSrc(64)
		// Insert redis
		err = redisDB.Set("user:"+user.Email, token, 0).Err()
		if err != nil {
			panic(err)
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		err = redisDB.Set("user:"+user.Email+":password", hash, 0).Err()
		if err != nil {
			panic(err)
		}
		// Send email
		sendMailValidator(user.Email, token)
		// Set response
		response = basicResponse{
			Result: "success",
			Msg:    "You must validate your account, check your emails",
		}
	} else if err != nil {
		panic(err)
	} else {
		// User already present in DB
		response = basicResponse{
			Result: "failure",
			Msg:    "User already exists in database",
		}
	}

	// Send response
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Fprintln(w, string(jsonResponse))
}

func handleData(w http.ResponseWriter, req *http.Request) {
	var response basicResponse
	// Here we just verify that the user is authenticated
	reqToken := req.Header.Get("X-Bigdata4all-Token")
	token, err := jwt.Parse(reqToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err == nil && token.Valid {
		response = basicResponse{
			Result: "success",
			Msg:    "Valid token",
		}
	} else {
		response = basicResponse{
			Result: "failure",
			Msg:    "Invalid token",
		}
	}
	// Print response
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Fprintln(w, string(jsonResponse))
}

func handleUserAuth(w http.ResponseWriter, req *http.Request) {
	bodyb, _ := ioutil.ReadAll(req.Body)
	var user jsonUser
	_ = json.Unmarshal(bodyb, &user)
	var response basicResponse

	hash, err := redisDB.Get("user:" + user.Email + ":password").Result()
	if err == redis.Nil {
		// user doesn't exists
		response = basicResponse{
			Result: "failure",
			Msg:    "Invalid credentials",
		}
	} else if err != nil {
		panic(err)
	} else {
		// user exists, check password
		if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(user.Password)); err != nil {
			// Bad password
			response = basicResponse{
				Result: "failure",
				Msg:    "Invalid credentials",
			}
		} else {
			// Good password
			// Create JWT token
			token := jwt.New(jwt.GetSigningMethod("HS256"))
			claims := make(jwt.MapClaims)
			claims["user"] = user.Email
			claims["exp"] = time.Now().Add(time.Minute * 3600).Unix()
			token.Claims = claims
			tokenString, err := token.SignedString([]byte(secretKey))
			if err != nil {
				panic(err)
			}
			w.Header().Set("X-Bigdata4all-Token", tokenString)
			response = basicResponse{
				Result: "sucess",
				Msg:    "Feel free to use token",
			}
		}
	}

	// Send response
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Fprintln(w, string(jsonResponse))
}
