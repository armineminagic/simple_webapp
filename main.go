package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	_ "fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// CmdArgs holds values about command line flags
type CmdArgs struct {
	DbName string
	User string
	Password string
}

// cmdArgs for db connection info
var cmdArgs CmdArgs

// Phone holds information about phone
type Phone struct {
	DeviceModel 	string	`json:"model"`
	SerialNumber 	uint	`json:"serial"`
	Storage			uint	`json:"storage"`
	Color 			string	`json:"string"`	
}

func ifError(e error){
	if e != nil {
		panic(e.Error())
	}
}

// DbConn creates DB and table if not exists
func DbConn(ca CmdArgs) (db *sql.DB) {

	db, err := sql.Open("mysql", ca.User + ":" + ca.Password +"@tcp(127.0.0.1:3306)/?charset=utf8&autocommit=true")
	ifError(err)
	
	_, err = db.Exec(`create database if not exists `+ ca.DbName + `;`)	
	ifError(err)
	
	_, err = db.Exec(`use ` + ca.DbName + `;`)
	ifError(err)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS phone(
		serial_num int not null primary key,
		dev_model varchar(100) not null,
		storage int not null,
		color varchar(20) not null
		);`)
	ifError(err)
	
	return db
}


func IndexHandler(w http.ResponseWriter, r *http.Request){
	db := DbConn(cmdArgs)
	defer db.Close()
	rows, err := db.Query("Select * from phone")
	if err != nil { 
		panic(err.Error())
	}
	phone := Phone{}
	phones := []Phone{}
	for rows.Next() {
		err = rows.Scan(&phone.SerialNumber,&phone.DeviceModel, &phone.Storage, &phone.Color)
		if err != nil {
			panic(err.Error())
		}
		phones = append(phones, phone)
	}
	phonesJSON, err := json.Marshal(phones)
	ifError(err)
	fmt.Println(string(phonesJSON))
	w.Write(phonesJSON)
}

// func AddPhoneHandler(w http.ResponseWriter, r *http.Request){
// 	db := DbConn(cmdArgs)
// 	if r.Method == "POST" {

// 	}
// }

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Cotrol-Allow-Methods", "OPTIONS, GET, POST, PUT, PATCH, DELETE")
		log.Println(r.Method)
		log.Println(r.RequestURI)
		next.ServeHTTP(w,r)
	})
}

func main() {
	
	flag.StringVar(&cmdArgs.User, "u", "root", "database connection user")
	flag.StringVar(&cmdArgs.Password, "p", "root", "database connection password")
	flag.StringVar(&cmdArgs.DbName, "db", "phones", "database connection password")
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)
	r.Use(middleware)

	log.Fatal(http.ListenAndServe(":8080", r))

}