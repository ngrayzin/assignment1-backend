package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"database/sql"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Response struct {
	Message string `json:"message"`
}

type User struct {
	UserID              int            `json:"userID"`
	Email               string         `json:"email"`
	FirstName           string         `json:"firstName"`
	LastName            string         `json:"lastName"`
	Number              int            `json:"number"`
	IsCarOwner          bool           `json:"isCarOwner"`
	CarPlateNumber      sql.NullString `json:"carPlateNumber"`
	DriverLicenseNumber sql.NullString `json:"driverLicenseNumber"`
	Password            string         `json:"password"`
	AccountCreation     string         `json:"accountCreationDate"`
	LastUpdated         string         `json:"lastUpdated"`
}

type Trips struct {
	TripID             int
	ownerUserID        int
	pickupLoc          string
	altPickupLoc       sql.NullString
	startTravelTime    string
	destinationAddress string
	availableSeats     int
	isActive           bool
	createdAt          string
	lastUpdated        string
}

type TripEnrollment struct {
	EnrolmentID    int
	TripID         int
	PassengerID    int
	EnrollmentTime string
}

var db *sql.DB

var cfg = mysql.Config{
	User:      "user",
	Passwd:    "password",
	Net:       "tcp",
	Addr:      "localhost:3306",
	DBName:    "carpooling_db",
	ParseTime: true,
}

func main() {
	db, _ = sql.Open("mysql", cfg.FormatDSN())
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/login", login).Methods(http.MethodPost, http.MethodGet)
	router.HandleFunc("/api/v1/signup", signup).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/userProfile/{id}", userProfile).Methods(http.MethodGet, http.MethodPut, http.MethodDelete)
	router.HandleFunc("/api/v1/trips", trips).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/trips/{id}", trips).Methods(http.MethodPut)
	router.HandleFunc("/api/v1/myEnrolments/{id}", myEnrolments).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/publishTrip", publishTrip).Methods(http.MethodGet, http.MethodPut)
	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}

func login(w http.ResponseWriter, r *http.Request) {
	type LoginRequest struct {
		Email string `json:"email"`
		Pwd   string `json:"pwd"`
	}
	decoder := json.NewDecoder(r.Body)
	var loginRequest LoginRequest
	err := decoder.Decode(&loginRequest)
	if err != nil {
		// Handle JSON decoding error
	}

	email := loginRequest.Email
	pwd := loginRequest.Pwd

	if email == "" && pwd == "" {
		fmt.Println("Invalid params")
		return
	}

	fmt.Println("/api/v1/login")

	results, err := db.Query("SELECT * FROM Users WHERE Email = ? AND Password = ?", email, pwd)
	if err != nil {
		panic(err.Error())
	}
	defer results.Close()

	user := User{}
	userFound := false
	for results.Next() {
		userFound = true
		err = results.Scan(&user.UserID, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Number, &user.IsCarOwner, &user.DriverLicenseNumber, &user.CarPlateNumber, &user.AccountCreation, &user.LastUpdated)
		if err != nil {
			panic(err.Error())
		}
	}

	if userFound {
		userJSON, err := json.Marshal(user)
		if err != nil {
			panic(err.Error())
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Println("Logged in :D")
		w.Write(userJSON)
	} else {
		fmt.Println("Invalid login credentials")
		fmt.Fprintf(w, "Invalid login credentials")
	}
}

func signup(w http.ResponseWriter, r *http.Request) {
	type newUser struct {
		Email     string `json:"email"`
		Pwd       string `json:"pwd"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Number    int    `json:"number"`
	}
	decoder := json.NewDecoder(r.Body)
	var user newUser
	err := decoder.Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if user.Email == "" || user.Pwd == "" || user.FirstName == "" || user.LastName == "" || len(strconv.Itoa(user.Number)) != 8 {
		fmt.Println("Invalid params")
		http.Error(w, "Invalid parameters", http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO Users (Email, Password, FirstName, LastName, MobileNumber) VALUES (?, ?, ?, ?, ?)",
		user.Email, user.Pwd, user.FirstName, user.LastName, user.Number)
	if err != nil {
		panic(err.Error())
	}

	id, err := result.LastInsertId()
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("User created with ID:", id)
	fmt.Fprintf(w, "User created with ID: %d", id)
}

func userProfile(w http.ResponseWriter, r *http.Request) {
	querystringmap := r.URL.Query()
	userId := querystringmap["userId"]
	fmt.Println(userId)
	switch r.Method {
	case http.MethodPut:
	case http.MethodDelete:
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Error")
	}

}

func trips(w http.ResponseWriter, r *http.Request) {
	querystringmap := r.URL.Query()
	userId := querystringmap["userId"]
	fmt.Println(userId)
	switch r.Method {
	case http.MethodGet:
	case http.MethodPut:
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Error")
	}
}

func myEnrolments(w http.ResponseWriter, r *http.Request) {
	querystringmap := r.URL.Query()
	userId := querystringmap["userId"]
	fmt.Println(userId)

}

func publishTrip(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
	case http.MethodPut:
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Error")
	}
}
