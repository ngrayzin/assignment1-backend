package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"database/sql"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
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
	allowOrigins := handlers.AllowedOrigins([]string{"*"})
	allowMethod := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	allowHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	db, _ = sql.Open("mysql", cfg.FormatDSN())
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/login", login).Methods(http.MethodPost, http.MethodGet)
	router.HandleFunc("/api/v1/signup", signup).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/userProfile/{id}", userProfile).Methods(http.MethodGet, http.MethodPut, http.MethodDelete)

	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", handlers.CORS(allowHeaders, allowMethod, allowOrigins)(router)))
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
		w.WriteHeader(http.StatusForbidden)
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

	fmt.Println("/api/v1/signup")

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
	w.WriteHeader(http.StatusOK)
}

func userProfile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if _, ok := params["id"]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No ID")
	}
	id, _ := strconv.Atoi(params["id"])

	switch r.Method {
	case http.MethodGet:
		fmt.Printf("/api/v1/userProfile/%d\n", id)
		results, err := db.Query("SELECT UserID, Email, FirstName, LastName, MobileNumber, DriverLicenseNumber, CarPlateNumber FROM Users WHERE UserID = ?;", id)
		if err != nil {
			panic(err.Error())
		}
		defer results.Close()
		user := User{}
		for results.Next() {
			err = results.Scan(&user.UserID, &user.Email, &user.FirstName, &user.LastName, &user.Number, &user.DriverLicenseNumber, &user.CarPlateNumber)
			if err != nil {
				panic(err.Error())
			}
		}
		userJSON, err := json.Marshal(user)
		if err != nil {
			panic(err.Error())
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(userJSON)
	case http.MethodPut:
		var updateFields map[string]interface{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&updateFields); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Invalid request body")
			return
		}

		fmt.Printf("/api/v1/userProfile/%d\n", id)

		var setClauses []string
		var values []interface{}

		for key, value := range updateFields {
			if key == "IsCarOwner" {
				setClauses = append(setClauses, fmt.Sprintf("%s = ?", key))
				// Check if the value is a string representation of a boolean
				if strValue, ok := value.(string); ok {
					// Convert string representation to boolean
					boolValue, err := strconv.ParseBool(strValue)
					if err != nil {
						boolValue = false // Default value set to false
					}
					values = append(values, boolValue)
				} else if boolValue, ok := value.(bool); ok {
					values = append(values, boolValue)
				} else {
					values = append(values, value)
				}
			} else {
				setClauses = append(setClauses, fmt.Sprintf("%s = ?", key))
				values = append(values, value)
			}
		}

		query := fmt.Sprintf(`
			UPDATE Users
			SET %s
			WHERE UserID = ?;
		`, strings.Join(setClauses, ", "))

		values = append(values, id)
		rows, err := db.Query(query, values...)
		if err != nil {
			panic(err.Error())
		}
		defer rows.Close()
		fmt.Printf("User with id %d updated\n", id)
		fmt.Fprintf(w, "User data updated successfully\n")
		w.WriteHeader(http.StatusOK)
	case http.MethodDelete:
		results, err := db.Exec("DELETE FROM Users WHERE UserID = ? AND AccountCreationDate < DATE_SUB(NOW(),INTERVAL 1 YEAR);", id)
		if err != nil {
			panic(err.Error())
		}

		RowsEffected, err := results.RowsAffected()
		if err != nil {
			panic(err.Error())
		}

		if RowsEffected > 0 {
			w.WriteHeader(http.StatusOK)
			fmt.Printf("User with id %d deleted\n", id)
			fmt.Fprintf(w, "deleted user with id %d\n", id)
		} else {
			w.WriteHeader(http.StatusConflict)
			fmt.Println("Account cannot be deleted (1yr retention policy)")
			fmt.Fprintln(w, "Account cannot be deleted (1yr retention policy)")
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Error")
	}

}
