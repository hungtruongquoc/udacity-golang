package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestEndpoints(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/v1/customers/{id:[0-9]+}", getCustomerByID).Methods("GET")
	router.HandleFunc("/v1/customers/{id:[0-9]+}", deleteCustomerByID).Methods("DELETE")
	router.HandleFunc("/v1/customers/{id:[0-9]+}", patchCustomer).Methods("PATCH")
	router.HandleFunc("/v1/customers", getCustomers).Methods("GET")
	router.HandleFunc("/v1/customers", createACustomer).Methods("POST")

	// Test getCustomers
	t.Run("TestGetCustomers", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/v1/customers", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var actual []Customer
		err = json.Unmarshal(rr.Body.Bytes(), &actual)
		if err != nil {
			t.Fatal(err)
		}

		expected := customers

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
		}
	})

	t.Run("TestAddCustomers", func(t *testing.T) {
		// Create a new customer JSON
		newCustomer := Customer{
			Name:      "Dave Lee",
			Role:      "Designer",
			Email:     "dave.lee@example.com",
			Phone:     "555-0104",
			Contacted: false,
		}
		jsonData, err := json.Marshal(newCustomer)
		if err != nil {
			t.Fatalf("Error marshalling json: %v", err)
		}

		// Create a request to pass to our handler.
		request, err := http.NewRequest("POST", "/v1/customers", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatal(err)
		}
		request.Header.Set("Content-Type", "application/json")

		// Create a ResponseRecorder to record the response.
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)

		// Check the status code is what we expect.
		if status := recorder.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		var addedCustomer Customer
		if err := json.Unmarshal(recorder.Body.Bytes(), &addedCustomer); err != nil {
			t.Fatalf("Could not parse json: %v", err)
		}

		// Verify that the customer returned in response matches what we sent.
		if addedCustomer.Name != newCustomer.Name || addedCustomer.Email != newCustomer.Email {
			t.Errorf("handler returned unexpected body: got %v, want %v", addedCustomer, newCustomer)
		}

		// Optionally, check that the customer was added to the slice.
		if len(customers) == 0 || customers[len(customers)-1].Name != newCustomer.Name {
			t.Errorf("Customer was not added successfully to the database")
		}
	})

	// Test getCustomerByID
	t.Run("TestGetCustomerByID", func(t *testing.T) {
		validReq, err := http.NewRequest("GET", fmt.Sprintf("/v1/customers/%d", 1), nil)
		if err != nil {
			t.Fatal(err)
		}

		invalidReq, err := http.NewRequest("GET", "/v1/customers/9999", nil)
		if err != nil {
			t.Fatal(err)
		}

		testCases := []struct {
			name     string
			request  *http.Request
			expected int
			valid    bool // This is used to indicate if the test expects a valid customer or not
		}{
			{"Valid ID", validReq, http.StatusOK, true},
			{"Invalid ID", invalidReq, http.StatusNotFound, false},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				rr := httptest.NewRecorder()
				router.ServeHTTP(rr, tc.request)

				if status := rr.Code; status != tc.expected {
					t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expected)
				}

				if tc.valid {
					var customer Customer
					err := json.Unmarshal(rr.Body.Bytes(), &customer)
					if err != nil {
						t.Fatal(err)
					}
					expected := customers[0]
					if customer != expected {
						t.Errorf("handler returned unexpected body: got %v want %v", customer, expected)
					}
				}
			})
		}
	})

	t.Run("TestDeleteCustomerByID", func(t *testing.T) {
		// Create a request to delete a customer.
		request, err := http.NewRequest("DELETE", "/v1/customers/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)

		// Check the status code is what we expect.
		if status := recorder.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		// Optionally, check that the customer was actually deleted.
		found := false
		for _, customer := range customers {
			if customer.ID == 1 {
				found = true
				break
			}
		}
		if found {
			t.Errorf("Customer was not deleted successfully")
		}
	})

	t.Run("TestDeleteCustomerByInvalidID", func(t *testing.T) {
		// Create a request to delete a customer.
		request, err := http.NewRequest("DELETE", "/v1/customers/999999", nil)
		if err != nil {
			t.Fatal(err)
		}

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)

		// Check the status code is what we expect.
		if status := recorder.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("TestUpdateCustomerByInvalidID", func(t *testing.T) {
		// Data for updating the customer
		updateData := map[string]interface{}{
			"Email": "updated.email@example.com",
		}
		updateBody, _ := json.Marshal(updateData)

		request, err := http.NewRequest("PATCH", "/v1/customers/999999", bytes.NewBuffer(updateBody))
		if err != nil {
			t.Fatal(err)
		}
		request.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)

		// Check the status code is what we expect.
		if status := recorder.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	})

	t.Run("TestUpdateCustomerByValidID", func(t *testing.T) {
		// Data for updating the customer
		updateData := map[string]interface{}{
			"Email": "updated.email@example.com",
		}
		updateBody, _ := json.Marshal(updateData)

		request, err := http.NewRequest("PATCH", "/v1/customers/3", bytes.NewBuffer(updateBody))
		if err != nil {
			t.Fatal(err)
		}
		request.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)

		// Check the status code is what we expect.
		if status := recorder.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		// Check that the email was updated correctly.
		var updatedCustomer Customer
		json.Unmarshal(recorder.Body.Bytes(), &updatedCustomer)
		if updatedCustomer.Email != "updated.email@example.com" {
			t.Errorf("Email was not updated correctly: got %v want %v", updatedCustomer.Email, "updated.email@example.com")
		}
	})
}
