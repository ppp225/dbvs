package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/kataras/muxie"
	"github.com/ppp225/dbvs/postgres/db"
	"github.com/ppp225/envp"
)

var database db.Database

type response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func main() {
	// load env
	envp.LoadEnvFromEnvFiles("")

	// setup db connection
	dbUser, dbPassword, dbName :=
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB")
	var err error
	database, err = db.Initialize(dbUser, dbPassword, dbName)
	if err != nil {
		log.Fatalf("db init failed with msg: %v", err)
	}
	defer database.Conn.Close()

	// setup mux
	mux := muxie.NewMux()
	mux.PathCorrection = true

	// setup handlers
	v1 := mux.Of("/v1")
	v1.Handle("/items/:id", muxie.Methods().
		HandleFunc(http.MethodGet, getItem))
	v1.Handle("/items", muxie.Methods().
		HandleFunc(http.MethodGet, get1000Items).
		HandleFunc(http.MethodPost, postItem))

	// start server
	server := &http.Server{Addr: ":8090", Handler: mux}
	go func() {
		log.Printf("server exiting with msg: %v", server.ListenAndServe())
	}()
	log.Print("server started")
	defer stop(server)

	// handle interrupts
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	log.Printf("server received signal: %s", fmt.Sprint(<-s))
}

func getItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(muxie.GetParam(w, "id"))
	if err != nil {
		respondWithStatus(w, 400, "invalid id")
		return
	}
	item, err := database.GetItemById(id)
	if err != nil {
		switch err {
		case db.ErrNotFound:
			respondWithStatus(w, 404, err.Error())
		default:
			respondWithStatus(w, 500, err.Error())
		}
		return
	}
	respond(w, item)
}

func get1000Items(w http.ResponseWriter, r *http.Request) {
	items, err := database.Get1000Items()
	if err != nil {
		respondWithStatus(w, 500, err.Error())
		return
	}
	respond(w, items)
}

func postItem(w http.ResponseWriter, r *http.Request) {
	item := &db.Item{}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithStatus(w, 400, "invalid input body: "+err.Error())
		return
	}
	if err := json.Unmarshal(bodyBytes, &item); err != nil {
		respondWithStatus(w, 400, "invalid input body: "+err.Error())
		return
	}
	if err := database.AddItem(item); err != nil {
		respondWithStatus(w, 500, err.Error())
		return
	}
	respond(w, item)
}

func respond(w http.ResponseWriter, data interface{}) {
	resp := response{}
	writeDefaultHeaders(w)
	resp.Code = http.StatusOK
	resp.Msg = http.StatusText(resp.Code)
	resp.Data = data

	marshalledResponse, err := json.Marshal(resp)
	if err != nil {
		panic(err) // should not happen rly
	}
	w.Write([]byte(marshalledResponse))
}

func respondWithStatus(w http.ResponseWriter, code int, msg string) {
	resp := response{}
	writeDefaultHeaders(w)
	w.WriteHeader(code)
	resp.Code = code
	resp.Msg = msg

	marshalledResponse, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	w.Write([]byte(marshalledResponse))
}

func writeDefaultHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", w.Header().Get("Origin"))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
}

func stop(server *http.Server) {
	log.Print("server stop initiated...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
		os.Exit(1)
	}
	log.Print("server exited cleanly")
}
