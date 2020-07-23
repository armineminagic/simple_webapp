package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	_ "fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// CmdArgs holds values about command line flags
type CmdArgs struct {
	DbName   string
	User     string
	Password string
}

// cmdArgs for db connection info
var cmdArgs CmdArgs

// Student holds information about stud
type Student struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Indexnum uint   `json:"indexnum"`
	IDnumber uint64 `json:"id"`
}

func ifError(e error) {
	if e != nil {
		log.Fatal(e.Error())
	}
}

// DbConn creates DB and table if not exists
func DbConn(ca CmdArgs) (db *sql.DB) {

	db, err := sql.Open("mysql", ca.User+":"+ca.Password+"@tcp(127.0.0.1:3306)/"+ca.DbName)

	ifError(err)

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `students` (`index` INTEGER NOT NULL, `idstud` INTEGER NOT NULL, `name` VARCHAR(30) not null, `surname` VARCHAR(30) not null);")

	if err != nil {
		log.Fatal(err)
	}

	return db
}

// Students handler for homepage
func getStudents(w http.ResponseWriter, r *http.Request) {
	db := DbConn(cmdArgs)
	defer db.Close()
	rows, err := db.Query("Select * from students")
	if err != nil {
		log.Fatal(err.Error())
	}
	stud := Student{}
	students := []Student{}
	for rows.Next() {
		err = rows.Scan(&stud.Indexnum, &stud.IDnumber, &stud.Name, &stud.Surname)
		if err != nil {
			log.Fatal(err.Error())
		}
		students = append(students, stud)
	}
	studentsJSON, err := json.Marshal(students)
	ifError(err)
	w.Write(studentsJSON)
}

func searchStudent(w http.ResponseWriter, r *http.Request) {

	searchValue := mux.Vars(r)

	db := DbConn(cmdArgs)
	defer db.Close()

	row, err := db.Query("Select * from students where `index`= ?", searchValue["indexnum"])

	ifError(err)
	student := Student{}
	students := []Student{}
	for row.Next() {
		err = row.Scan(&student.Indexnum, &student.IDnumber, &student.Name, &student.Surname)
		ifError(err)
		students = append(students, student)
	}
	studentsJSON, err := json.Marshal(students)
	ifError(err)
	w.Write(studentsJSON)

}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		log.Println(r.Method)
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func main() {

	flag.StringVar(&cmdArgs.User, "u", "root", "database connection user")
	flag.StringVar(&cmdArgs.Password, "p", "root", "database connection password")
	flag.StringVar(&cmdArgs.DbName, "db", "faculty", "database name")
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/", getStudents)
	r.HandleFunc("/search/{indexnum}", searchStudent)
	r.Use(middleware)

	log.Fatal(http.ListenAndServe(":8080", r))
}
