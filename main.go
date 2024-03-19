package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// User struct to represent a user
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var adminUsers []User

func main() {
	loadAdmins()

	router := mux.NewRouter()

	// CORS ayarları
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// Admin login endpoint
	router.HandleFunc("/admin/login", adminLogin).Methods("POST")

	// CORS ayarlarını kullanarak HTTP sunucusunu başlat
	fmt.Println("Server is running on port 8000")
	http.ListenAndServe(":8000", handlers.CORS(headersOk, originsOk, methodsOk)(router))
}

// Function to load admins from JSON file
func loadAdmins() {
	fileContent, err := ioutil.ReadFile("admins.json")
	if err != nil {
		fmt.Println("Error reading admins.json:", err)
		return
	}

	err = json.Unmarshal(fileContent, &adminUsers)
	if err != nil {
		fmt.Println("Error parsing admins.json:", err)
		return
	}
}

// Function to authenticate admin users
func authenticateAdmin(username, password string) bool {
	for _, admin := range adminUsers {
		if admin.Username == username && admin.Password == password {
			return true
		}
	}
	return false
}

// Function to simulate admin login
func adminLogin(w http.ResponseWriter, r *http.Request) {
	// Parse Basic Auth headers
	username, password, ok := r.BasicAuth()
	if !ok {
		http.Error(w, "Unauthorized2", http.StatusUnauthorized)
		return
	}

	// Authenticate admin
	if authenticateAdmin(username, password) {
		// Admin login successful
		fmt.Println("Admin logged in")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Admin logged in successfully")
	} else {
		// Admin login failed
		fmt.Println("Admin login failed")
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Unauthorized1")
	}
}
