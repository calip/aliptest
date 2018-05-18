// package main

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"
// 	mgo "gopkg.in/mgo.v2"
// 	"gopkg.in/mgo.v2/bson"

// 	"github.com/gorilla/mux"
// 	. "music/config"
// 	. "music/dao"
// 	. "music/models"
// )

// var config = Config{}
// var dao = MoviesDAO{}

// type Repository struct {
// 	Server   string
// 	Database string
// }
// type Api struct {
// 	Repository Repository
// }

// var db *mgo.Database

// // GET list of movies
// func AllMoviesEndPoint(w http.ResponseWriter, r *http.Request) {
// 	movies, err := dao.FindAll()
// 	if err != nil {
// 		respondWithError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	respondWithJson(w, http.StatusOK, movies)
// }

// // GET a movie by its ID
// func FindMovieEndpoint(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	movie, err := dao.FindById(params["id"])
// 	if err != nil {
// 		respondWithError(w, http.StatusBadRequest, "Invalid Movie ID")
// 		return
// 	}
// 	respondWithJson(w, http.StatusOK, movie)
// }

// // POST a new movie
// func CreateMovieEndPoint(w http.ResponseWriter, r *http.Request) {
// 	defer r.Body.Close()
// 	var movie Movie
// 	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
// 		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
// 		return
// 	}
// 	movie.ID = bson.NewObjectId()
// 	if err := dao.Insert(movie); err != nil {
// 		respondWithError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	respondWithJson(w, http.StatusCreated, movie)
// }

// // PUT update an existing movie
// func UpdateMovieEndPoint(w http.ResponseWriter, r *http.Request) {
// 	defer r.Body.Close()
// 	var movie Movie
// 	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
// 		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
// 		return
// 	}
// 	if err := dao.Update(movie); err != nil {
// 		respondWithError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
// }

// // DELETE an existing movie
// func DeleteMovieEndPoint(w http.ResponseWriter, r *http.Request) {
// 	defer r.Body.Close()
// 	var movie Movie
// 	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
// 		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
// 		return
// 	}
// 	if err := dao.Delete(movie); err != nil {
// 		respondWithError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
// }

// func respondWithError(w http.ResponseWriter, code int, msg string) {
// 	respondWithJson(w, code, map[string]string{"error": msg})
// }

// func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
// 	response, _ := json.Marshal(payload)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(code)
// 	w.Write(response)
// }

// // Parse the configuration file 'config.toml', and establish a connection to DB
// func init() {
// 	config.Read()

// 	dao.Server = config.Server
// 	dao.Database = config.Database
// 	session, err := mgo.Dial(m.Server)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	db = session.DB(m.Database)
// }

// // Define HTTP request routes
// func main() {
// 	app := Api{
// 		Repository: repository,
// 	}
// 	r := mux.NewRouter()
// 	app.ConfigureRoutes(mux)
	
// 	if err := http.ListenAndServe(":9000", r); err != nil {
// 		log.Fatal(err)
// 	}
// }

// func (app *Api) ConfigureRoutes(router *mux.Router) {
// 	r.HandleFunc("/movies", app.AllMoviesEndPoint).Methods("GET")
// 	r.HandleFunc("/movies", app.CreateMovieEndPoint).Methods("POST")
// 	r.HandleFunc("/movies", app.UpdateMovieEndPoint).Methods("PUT")
// 	r.HandleFunc("/movies", app.DeleteMovieEndPoint).Methods("DELETE")
// 	r.HandleFunc("/movies/{id}", app.FindMovieEndpoint).Methods("GET")
// }


package main

import (
	"./api"
)

func main() {
	api.InitServer()
}
