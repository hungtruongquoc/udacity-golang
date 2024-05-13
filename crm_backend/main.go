package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Customer struct to hold the details of a customer.
type Customer struct {
	ID        int
	Name      string
	Role      string
	Email     string
	Phone     string
	Contacted bool
}

// Slice of Customer to act as a database.
var customers = []Customer{
	{ID: 1, Name: "Alice Johnson", Role: "Manager", Email: "alice.johnson@example.com", Phone: "555-0101", Contacted: false},
	{ID: 2, Name: "Bob Smith", Role: "Developer", Email: "bob.smith@example.com", Phone: "555-0102", Contacted: true},
	{ID: 3, Name: "Carol Taylor", Role: "Analyst", Email: "carol.taylor@example.com", Phone: "555-0103", Contacted: false},
}

var dictionary = map[string]string{
	"Go":     "A programming language created by Google.",
	"Gopher": "A software engineer who builds with Go.",
	"Golang": "Another name for Go.",
}

func getDictionary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(dictionary)
}

func getCustomers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customers)
}

func createACustomer(w http.ResponseWriter, r *http.Request) {
	var newCustomer Customer
	if err := json.NewDecoder(r.Body).Decode(&newCustomer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newCustomer.ID = len(customers) + 1 // simplistic approach to assign a new ID
	customers = append(customers, newCustomer)
	json.NewEncoder(w).Encode(newCustomer)
}

func getCustomerByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	for _, customer := range customers {
		if customer.ID == id {
			json.NewEncoder(w).Encode(customer)
			return
		}
	}

	http.NotFound(w, r)
}

func deleteCustomerByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	// Find and delete the customer from the slice
	for i, customer := range customers {
		if customer.ID == id {
			// Delete the customer by slicing out their entry
			customers = append(customers[:i], customers[i+1:]...)
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "Deleted successfully")
			return
		}
	}

	// If no customer found, return an error
	http.NotFound(w, r)
}

// patchCustomer updates specified fields of a customer.
func patchCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	// Find the customer
	for i, customer := range customers {
		if customer.ID == id {
			// Read request body
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body", http.StatusInternalServerError)
				return
			}
			defer r.Body.Close()

			var update Customer
			// Parse JSON data
			if err := json.Unmarshal(body, &update); err != nil {
				http.Error(w, "Error parsing JSON data", http.StatusBadRequest)
				return
			}

			// Update the customer's fields if provided in the JSON
			if update.Name != "" {
				customer.Name = update.Name
			}
			if update.Role != "" {
				customer.Role = update.Role
			}
			if update.Email != "" {
				customer.Email = update.Email
				fmt.Println(customer.Email)
			}
			if update.Phone != "" {
				customer.Phone = update.Phone
			}
			if update.Contacted {
				customer.Contacted = update.Contacted
			}

			customers[i] = customer // Make sure to update the slice's element

			json.NewEncoder(w).Encode(customer)
			return
		}
	}

	http.NotFound(w, r)
}

func main() {
	// Create a new router
	router := mux.NewRouter()

	router.HandleFunc("/customers/{id:[0-9]+}", getCustomerByID).Methods("GET")
	router.HandleFunc("/customers/{id:[0-9]+}", deleteCustomerByID).Methods("DELETE")
	router.HandleFunc("/customers/{id:[0-9]+}", patchCustomer).Methods("PATCH")
	router.HandleFunc("/customers", getCustomers).Methods("GET")
	router.HandleFunc("/customers", createACustomer).Methods("POST")
	// Serve static files under the 'static' directory only if no other route matches.
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("static/"))))

	fmt.Println("Server is starting on port 3000...")
	http.ListenAndServe(":3000", router)
}
