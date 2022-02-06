package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// Database to share and download graphs
var db *sql.DB

// Not great security but if someone wants to waste their time messing stuff up then can make it more secure
const dbUrl = "postgres://vdsjgptjfnkwzg:e1ced5319036f0ab82cf7607947ada6a39fe7e34bbef358f1b0900a691c924e8@ec2-52-45-83-163.compute-1.amazonaws.com:5432/ddrck8fch7vb48"

type GraphDB struct {
	Name  string
	Graph string
	Date  time.Time
}
type GraphsDB []*GraphDB

// Initialize the database
func InitDB() {
	var err error
	db, err = sql.Open("postgres", dbUrl)
	if HandleError(err) {
		return
	}
	err = db.Ping()
	HandleError(err)
}

// Upload a graph to the database
func UploadGraph(name, data string) {
	b, _ := time.Now().UTC().MarshalText()
	cmd := fmt.Sprintf("INSERT INTO Graphs(Name, Data, Time) VALUES ('%v', '%v', '%v');", name, data, string(b))
	_, err := db.Exec(cmd)
	HandleError(err)
}

// Get the graphs from the database
func GetGraphs() GraphsDB {
	cmd := "SELECT * FROM Graphs"
	rows, err := db.Query(cmd)
	if HandleError(err) {
		return nil
	}
	defer rows.Close()
	theGraphsDB := make(GraphsDB, 0, 10)
	for rows.Next() {
		var name, data, dstring string
		var id int
		err := rows.Scan(&id, &name, &data, &dstring)
		if HandleError(err) {
			break
		}
		date, err := time.Parse(time.RFC3339, dstring)
		if HandleError(err) {
			break
		}
		if time.Since(date).Hours() >= 168 { // if a week has passed since a graph was published remove it
			RemoveGraph(id)
			break
		}
		theGraphsDB = append(theGraphsDB, &GraphDB{name, data, date})
	}
	return theGraphsDB

}

// Remove a graph given id
func RemoveGraph(id int) {
	cmd := fmt.Sprintf("DELETE FROM Graphs WHERE Id=%v", id)
	_, err := db.Exec(cmd)
	HandleError(err)
}