package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/cockroachdb/cmux"
	pb "github.com/contd/links-rpc/links"
	"google.golang.org/grpc"
)

const (
	port   = 5051
	dbpath = "saved.sqlite"
)

func main() {
	// formatted address to listen on
	addr := fmt.Sprintf(":%d", port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	// create the cmux object that will multiplex 2 protocols on same port
	m := cmux.New(l)
	// match gRPC requests, otherwise regular HTTP requests
	grpcL := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpL := m.Match(cmux.Any())

	// create the go-grpc example greeter server
	s := NewServer(dbpath)
	grpcS := grpc.NewServer()
	pb.RegisterLinksServer(grpcS, s)

	// create the regular HTTP requests muxer
	h := http.NewServeMux()
	h.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})
	h.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// The "/" pattern matches everything not matched by previous handlers
		fmt.Fprintf(w, "Links Service Web Interface\n")
	})
	h.HandleFunc("/links", func(w http.ResponseWriter, r *http.Request) {
		//links, err := getLinks(s.DB)
		rows, err := s.DB.Query("SELECT id, url, category, created_on, done FROM links")
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		defer rows.Close()
		links := []link{}

		for rows.Next() {
			var l link
			if err := rows.Scan(&l.ID, &l.Url, &l.Category, &l.Created, &l.Done); err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
			links = append(links, l)
		}
		respondWithJSON(w, http.StatusOK, links)
	})
	httpS := &http.Server{
		Handler: h,
	}

	// collect on this channel the exits of each protocol's .Serve() call
	eps := make(chan error, 2)

	// start the listeners for each protocol
	go func() { eps <- grpcS.Serve(grpcL) }()
	go func() { eps <- httpS.Serve(httpL) }()

	log.Println("listening and serving (multiplexed) on", addr)
	err = m.Serve()

	// the rest of the code handles exit errors of the muxes

	var failed bool
	if err != nil {
		log.Printf("cmux serve error: %v\n", err)
		failed = true
	}
	var i int
	for err := range eps {
		if err != nil {
			log.Printf("protocol serve error: %v", err)
			failed = true
		}
		i++
		if i == cap(eps) {
			close(eps)
			break
		}
	}
	if failed {
		os.Exit(1)
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
