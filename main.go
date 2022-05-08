package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/imdario/mergo"
)

type Movie struct {
	Id       string  `json:"id"`
	Title    string  `json:"title"`
	Desc     string  `json:"desc"`
	Duration float64 `json:"duration"`
}

var Movies []Movie

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)
	router.HandleFunc("/movies/delete/{id}", deleteSingleMovie).Methods("DELETE")
	router.HandleFunc("/movies/{id}", updateSingleMovie).Methods("PUT")
	router.HandleFunc("/movies", createSingleMovie).Methods("POST")
	router.HandleFunc("/movies", returnAllMovies)
	router.HandleFunc("/movies/{id}", returnSingleMovie)

	log.Fatal(http.ListenAndServe(":10000", router))
}

// Return all movies in Movies array
func returnAllMovies(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	json.NewEncoder(w).Encode(Movies)
}

// Return single Movie from Movies array with given id
func returnSingleMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	var movieToOutput Movie
	for _, movie := range Movies {
		if movie.Id == key {
			movieToOutput = movie
		}
	}
	outputHandler(movieToOutput, key, w)

}

// Create single movie and add to Movies array
func createSingleMovie(w http.ResponseWriter, r *http.Request) {

	movie := unmarshalJson(r)
	Movies = append(Movies, movie)

	json.NewEncoder(w).Encode(movie)
}

func deleteSingleMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	var deletedMovie Movie

	for i, movie := range Movies {
		if movie.Id == key {
			deletedMovie = Movies[i]
			Movies[i] = Movies[len(Movies)-1]
			Movies = Movies[:len(Movies)-1]
		}
	}
	outputHandler(deletedMovie, key, w)

}

func updateSingleMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	newMovie := unmarshalJson(r)
	var movieToOutput Movie
	for i, movie := range Movies {
		if movie.Id == key {
			mergo.Merge(&newMovie, Movies[i])
			Movies[i] = newMovie
			movieToOutput = Movies[i]
		}
	}
	outputHandler(movieToOutput, key, w)

}

func unmarshalJson(r *http.Request) Movie {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var movie Movie
	json.Unmarshal(reqBody, &movie)
	return movie
}

func outputHandler(movieToOutput Movie, key string, w http.ResponseWriter) {
	if (Movie{} == movieToOutput) {
		json.NewEncoder(w).Encode("Movie with id " + key + " does not exist")
	} else {
		json.NewEncoder(w).Encode(movieToOutput)
	}
}

func main() {

	Movies = []Movie{
		Movie{Id: "1", Title: "Title1", Desc: "Desc1", Duration: 90.64},
		Movie{Id: "2", Title: "Title2", Desc: "Desc2", Duration: 150},
	}

	handleRequests()

}
