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
	Port       string
	Adjectives []string
	Nouns      []string
}

//go:embed data
var data embed.FS

func main() {
	svr := newServer(os.Getenv("PORT"))

	// handle incoming requests
	http.HandleFunc("/", svr.handler)

	// Determine port for HTTP service.

	// Start HTTP server.
	log.Printf("listening on port %s", svr.Port)
	if err := http.ListenAndServe(":"+svr.Port, nil); err != nil {
		log.Fatal(err)
	}

}

func newServer(port string) *Service {
	if port == "" {
		port = "8080"
	}
	svr := Service{
		Adjectives: []string{},
		Nouns:      []string{},
		Port:       port,
	}
	svr.instantiate()
	return &svr
}

// instantiate populates the word slices with the embeded files
func (svr *Service) instantiate() {
	adj, _ := data.ReadFile("data/adj.txt")
	noun, _ := data.ReadFile("data/noun.txt")
	svr.Adjectives = strings.Split(string(adj), " ")
	svr.Nouns = strings.Split(string(noun), " ")
}

// name creates a random name
func (svr *Service) name(seed int64) string {
	if seed < 1 {
		seed = time.Now().UnixNano() / int64(time.Millisecond)
	}
	r := rand.New(rand.NewSource(seed))

	return fmt.Sprintf("%s-%s",
		svr.Adjectives[r.Intn(len(svr.Adjectives))],
		svr.Nouns[r.Intn(len(svr.Nouns))])
}

// handler handles http requests
func (svr *Service) handler(w http.ResponseWriter, r *http.Request) {
	var out Response

	out.Data = svr.name(time.Now().UnixMicro())
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
