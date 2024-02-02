package main

import (
	"encoding/json"
	"fmt"
	"log"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Movie struct {
	ID       int       `json:"id,string"`
	Name     string    `json:"name"`
	Genre    string    `json:"genre"`
	Director *Director `json:"director"`
}

type Director struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

var movies []Movie

func main() {
	movie1 := Movie{
		ID:    1,
		Name:  "Dabaang",
		Genre: "Adventure",
		Director: &Director{
			FirstName: "Arbaaz",
			LastName:  "Khan",
		},
	}
	movie2 := Movie{
		ID:    2,
		Name:  "MS Dhone: The untold story",
		Genre: "Sports",
		Director: &Director{
			FirstName: "Neeraj",
			LastName:  "Pandey",
		},
	}
	movie3 := Movie{
		ID:    3,
		Name:  "Deewar",
		Genre: "Adventure",
		Director: &Director{
			FirstName: "Yash",
			LastName:  "Chopra",
		},
	}
	movies = append(movies, movie1, movie2, movie3)

	router := mux.NewRouter()

	router.HandleFunc("/", HomePage)
	router.HandleFunc("/movies", GetMovies).Methods("GET")
	router.HandleFunc("/movies/{id}", GetMovie).Methods("GET")
	router.HandleFunc("/movies/{id}", DeleteMovie).Methods("DELETE")
	router.HandleFunc("/movies", CreateMovie).Methods("POST")
	router.HandleFunc("/movies/{id}", UpdateMovie).Methods("PUT")
	http.ListenAndServe(":8080", router)
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the home page"))
}

func GetMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	jsonBytes, err := json.Marshal(movies)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func GetMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	params := mux.Vars(r)
	res := params["id"]
	id, err := strconv.Atoi(res)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	for _, movie := range movies {
		if movie.ID == id {
			jsonBytes, err := json.Marshal(movie)
			if err != nil {
				log.Fatalf("%v", err)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write(jsonBytes)
			break
		}
	}
	fmt.Printf("SUCCESS!")
}

func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	params := mux.Vars(r)
	res := params["id"]
	id, err := strconv.Atoi(res)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	for index, movie := range movies {
		if movie.ID == id {
			movies = append(movies[:index], movies[index+1:]...)
			jsonBytes, err := json.Marshal(movie)
			if err != nil {
				log.Fatalf("%v", err)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write(jsonBytes)
			break
		}
	}
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var newMovie Movie

	for _, movie := range movies {
		if movie.ID == newMovie.ID {
			w.Write([]byte("Movie ID already exists"))
			return
		}
	}

	if err := json.NewDecoder(r.Body).Decode(&newMovie); err != nil {
		log.Fatalf("%v", err)
		return
	}

	movies = append(movies, newMovie)
	json.NewEncoder(w).Encode(movies)
	w.WriteHeader(http.StatusOK)
}

func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	res := params["id"]
	id, _ := strconv.Atoi(res)

	// check if movie exist already
	for index, movie := range movies {
		if movie.ID == id {
			movies = append(movies[:index], movies[index+1:]...)
			var newMovie Movie
			json.NewDecoder(r.Body).Decode(&newMovie)
			movies = append(movies, newMovie)
			json.NewEncoder(w).Encode(newMovie)
			return
		}
	}

	fmt.Fprintf(w,"Movie of id %d does not exist",id)
}
