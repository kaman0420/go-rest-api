package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Create a Dummy database
type event struct {
	ID        string  `json: "ID"`
	Latitude  float64 `json: "Latitude"`
	Longitude float64 `json: "Longitude"`
}

type allEvent []event

var events = allEvent{
	{
		ID:        "1",
		Latitude:  -33.906561,
		Longitude: 151.192165,
	},
}

// Set up the HTTP server: homepage
func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Thinxtra")
}

// Create an event
func createEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent event

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event latitude and longitude only in order to update")
	}
	json.Unmarshal(reqBody, &newEvent)
	events = append(events, newEvent)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newEvent)
	fmt.Println("created!")
}

// Get one event
func getOneEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	for _, singleEvent := range events {
		if singleEvent.ID == eventID {
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}

// Get all events
func getAllEvents(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(events)
}

// Update an event
func updateEvent(w http.ResponseWriter, r *http.Request) {
	// declare event id
	eventID := mux.Vars(r)["id"]

	// create an updated event
	var updatedEvent event
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event latitude and longitude only in order to update")
	}

	// decode
	json.Unmarshal(reqBody, &updatedEvent)

	// filter out a specific event from the events slice
	for i, singleEvent := range events {
		// when the id resembles the one of an event in the slice
		if singleEvent.ID == eventID {
			// update the values of Latitude & Longitude
			singleEvent.Latitude = updatedEvent.Latitude
			singleEvent.Longitude = updatedEvent.Longitude
			// update the value of the struct in the events slice
			events = append(events[:i], singleEvent)
			// encode
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}

// Delete an event
func deleteEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	for i, singleEvent := range events {
		if singleEvent.ID == eventID {
			// Remove: slice = append(slice[:i], slice[i+1:]...)
			events = append(events[:i], events[i+1:]...)
			fmt.Fprintf(w, "The event with ID %v has been deleted successfully", eventID)
		}
	}
}

func main() {
	// Set up the HTTP server using Gorilla Mux
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/event", createEvent).Methods("POST")
	router.HandleFunc("/events", getAllEvents).Methods("GET")
	router.HandleFunc("/events/{id}", getOneEvent).Methods("GET")
	router.HandleFunc("/events/{id}", updateEvent).Methods("PATCH")
	router.HandleFunc("/events/{id}", deleteEvent).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// {
// 	"Id": "2",
// 	"Latitude": -33.891321,
// 	"Longitude": 151.198857
// }

// {
// 	"Latitude": -34.426632,
// 	"Longitude": 150.888373
// }
