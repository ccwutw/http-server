package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

// Create an unexported global variable to hold the database connection pool.
var db *sql.DB

type insertpair struct {
	RecKey       string `json:"key"`
	RecValue     string `json:"value"`
	RecTimestamp int    `json:"timestamp"`
}

type readpair struct {
	RecKey       string `json:"key"`
	RecTimestamp int    `json:"timestamp"`
}

func handleGet(w http.ResponseWriter, req *http.Request) {
	fmt.Println("handleGet")
	//fmt.Fprintf(w, "in get method:\n")
	//Decode JSON
	data := json.NewDecoder(req.Body)
	var kvmap readpair
	err := data.Decode(&kvmap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("handleGet Error")
	}
	fmt.Fprintf(w, "Get:%+v \n", kvmap)

	//Read DB
	var value string

	row := db.QueryRow("select value from mytable where key = ? AND timestamp <= ? order by timestamp DESC limit 1", kvmap.RecKey, kvmap.RecTimestamp)
	err = row.Scan(&value)

	if err != nil {
		if err == sql.ErrNoRows {
			// Handle the case of no rows found
			return
		}
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Get:%s \n", value)
}

func handlePut(w http.ResponseWriter, req *http.Request) {
	fmt.Println("handlePut")
	//fmt.Fprintf(w, "in put method:\n")
	//Decode JSON
	data := json.NewDecoder(req.Body)
	var kvmap insertpair
	err := data.Decode(&kvmap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "Put:%+v \n", kvmap)

	//write DB
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("insert into mytable(key, value, timestamp) values(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(kvmap.RecKey, kvmap.RecValue, kvmap.RecTimestamp)
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

}

func main() {

	//Database
	os.Remove("./kvmap.db")
	var err error
	// Make sure not to shadow your global - just assign with = - don't initialise a new variable and assign with :=
	db, err = sql.Open("sqlite3", "./kvmap.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table IF NOT EXISTS mytable (id integer not null primary key, key text, value text, timestamp datetime DEFAULT(STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')));
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	//Server
	fmt.Println("Server Start !")
	router := mux.NewRouter()
	router.HandleFunc("/", handlePut).Methods("PUT")
	router.HandleFunc("/", handleGet).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))

}
