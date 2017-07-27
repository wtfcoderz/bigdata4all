package main

import (
	"encoding/json"
	"io/ioutil"
	"flag"
	"fmt"
	"math/rand"
	"github.com/gorilla/mux"
	"log"
	"time"
	"net/http"
	"net/smtp"
	"os"
	"github.com/go-redis/redis"
//	"gopkg.in/gomail.v2"
)

type jsonUserRegister struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type basicResponse struct {
	Result string `json:"result"`
	Msg string `json:"msg"`
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
    letterIdxBits = 6                    // 6 bits to represent a letter index
    letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
    letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func RandStringBytesMaskImprSrc(n int) string {
    b := make([]byte, n)
    // A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
    for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
        if remain == 0 {
            cache, remain = src.Int63(), letterIdxMax
        }
        if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
            b[i] = letterBytes[idx]
            i--
        }
        cache >>= letterIdxBits
        remain--
    }

    return string(b)
}

var port string
var redisDB *redis.Client

func init() {
	flag.StringVar(&port, "port", "80", "give me a port number")
}

func main() {
	flag.Parse()

	// Init DB
	redisDB = redis.NewClient(&redis.Options{
		Addr:		"redis:6379",
		Password:	"",
		DB:		0,
	})

	// Mux Router
	r := mux.NewRouter()
	r.HandleFunc("/user/register", handleUserRegister).Methods("Post")
	r.HandleFunc("/health", health)
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
	var userRegister jsonUserRegister
	_ = json.Unmarshal(bodyb, &userRegister)

	// Verify user
	val2, err := redisDB.Get("user:" + userRegister.Email).Result()
	if err == redis.Nil {
		fmt.Println("user does not exists")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("user already exists", val2)
	}

	token := RandStringBytesMaskImprSrc(64)
	sendMailValidator(userRegister.Email, token)

	// Insert redis
	err = redisDB.Set("user:" + userRegister.Email, token, 0).Err()
	if err != nil {
		panic(err)
	}
	val, err := redisDB.Get("user:" + userRegister.Email).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("user:" + userRegister.Email, val)

	response := basicResponse{
		Result: "success",
		Msg: "You must validate your account, check your emails",
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Fprintln(w, string(jsonResponse))
}

func sendMailValidator(to string, token string) {
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASSWORD")
	from := "contact@bigdata4all.io"
	body := token

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Hello there\n\n" +
		body

	err := smtp.SendMail(os.Getenv("SMTP_SERVER") + ":" + os.Getenv("SMTP_PORT"),
		smtp.PlainAuth("", user, pass, os.Getenv("SMTP_SERVER")),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
}

