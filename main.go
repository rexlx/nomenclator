package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type Response struct {
	Data string `json:"data"`
}

type Service struct {
	Adjectives []string
	Nouns      []string
}

//go:embed data
var data embed.FS

func main() {
	// init our little library
	svr := Service{
		Adjectives: []string{},
		Nouns:      []string{},
	}

	// populate the service with words
	svr.instantiate()

	// handle incoming requests
	http.HandleFunc("/", svr.handler)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

}

// instantiate populates the word slices with the embeded files
func (svr *Service) instantiate() {
	adj, _ := data.ReadFile("data/adj.gof")
	noun, _ := data.ReadFile("data/noun.gof")
	svr.Adjectives = strings.Split(string(adj), " ")
	svr.Nouns = strings.Split(string(noun), " ")
}

// name creates a random name
func (svr *Service) name(seed int64) string {
	if seed < 1 {
		seed = time.Now().Unix()
	}
	rand.Seed(seed)

	return fmt.Sprintf("%s-%s",
		svr.Adjectives[rand.Intn(len(svr.Adjectives))],
		svr.Nouns[rand.Intn(len(svr.Nouns))])

}

// handler handles http requests
func (svr *Service) handler(w http.ResponseWriter, r *http.Request) {
	var out Response

	out.Data = svr.name(time.Now().Unix())
	res, err := json.Marshal(out)
	if err != nil {
		log.Fatalln(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(res)
	if err != nil {
		log.Fatalln("failed to write the response", err)
	}
}
