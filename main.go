package main

import (
	"fmt"
	"log"
	"net/http"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Users struct {
	profileId           int
	ID                  int
	FirstName           string
	LastName            string
	number              int
	IsCarOwner          bool
	CarPlateNumber      sql.NullString
	DriverLicenseNumber sql.NullString
}

func main() {
	db, err := sql.Open("mysql",
		"mysql:password@tcp(127.0.0.1:3306)/carpooling_db")
	// handle error
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	GetData(db)

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/", home)
	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello")
}

func GetData(db *sql.DB) {
	results, err := db.Query("SELECT * FROM UserProfiles")
	if err != nil {
		panic(err.Error())
	}
	defer results.Close()

	for results.Next() {
		var u Users
		err = results.Scan(&u.profileId, &u.ID, &u.FirstName, &u.LastName, &u.number, &u.IsCarOwner, &u.CarPlateNumber, &u.DriverLicenseNumber)
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("%d %d %s %s %s %t %s %s\n", u.profileId, u.ID, u.FirstName, u.LastName, u.number, u.IsCarOwner, u.CarPlateNumber, u.DriverLicenseNumber)
	}
}
