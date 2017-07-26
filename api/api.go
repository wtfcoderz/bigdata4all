package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

var port string

func init() {
	flag.StringVar(&port, "port", "80", "give me a port number")
}

func main() {
	flag.Parse()
	r := mux.NewRouter()
	r.HandleFunc("/health", health)
	r.HandleFunc("/", health)
	http.Handle("/", r)
	fmt.Println("Starting up on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func health(w http.ResponseWriter, req *http.Request) {
	hostname, _ := os.Hostname()
	fmt.Fprintln(w, "Hostname:", hostname)
}
