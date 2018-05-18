package api

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"encoding/json"
	"net/http"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"
	. "test/config"
	. "test/models"
	"fmt"
	"time"
	"strings"
)

var config = Config{}

var DB_NAME = ""

const (
	NOT_FOUND        = "not found"
	VALIDATION_ERROR = "Validation error"
)

type Repository interface {
	DB(name string) *mgo.Database
	Close()
}

type Api struct {
	Repository Repository
}

func (a *Api) LoginEndPoint(w http.ResponseWriter, r *http.Request) {
	var payload User
	var profile Profile

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	profile, err = Login(payload, a.Repository)
	if err != nil && err.Error() == VALIDATION_ERROR {
		respondWithError(w, http.StatusBadRequest, "VALIDATION_ERROR")
		return
	}
	if err != nil && err.Error() == NOT_FOUND {
		respondWithError(w, http.StatusUnauthorized, "NOT_FOUND")
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "NOT_FOUND")
		return
	}
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{ Name: "token",Value:profile.Token,Expires:expiration}
	http.SetCookie(w, &cookie)

	respondWithJson(w, http.StatusOK, profile)
}

func (a *Api) AllTestEndPoint(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No Auth ")
		return
	}
	authorization := cookie.Value
	if !IsAuthorized(authorization) {
		respondWithError(w, http.StatusUnauthorized, "No Auth ")
		return
	}

	tests, err := TestFindAll(a.Repository)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, tests)
}

func (a *Api) FindTestEndpoint(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No Auth ")
		return
	}
	authorization := cookie.Value
	if !IsAuthorized(authorization) {
		respondWithError(w, http.StatusUnauthorized, "No Auth ")
		return
	}

	params := mux.Vars(r)
	test, err := TestFindById(params["id"], a.Repository)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid test ID")
		return
	}
	respondWithJson(w, http.StatusOK, test)
}

func (a *Api) CreateTestEndPoint(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No Auth ")
		return
	}
	authorization := cookie.Value
	if !IsAuthorized(authorization) {
		respondWithError(w, http.StatusUnauthorized, "No Auth ")
		return
	}

	var test Test
	if err := json.NewDecoder(r.Body).Decode(&test); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	test.ID = bson.NewObjectId()
	if err := InsertTest(test, a.Repository); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, test)
}

func (a *Api) UpdateTestEndPoint(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No Auth ")
		return
	}
	authorization := cookie.Value
	if !IsAuthorized(authorization) {
		respondWithError(w, http.StatusUnauthorized, "No Auth ")
		return
	}
	
	params := mux.Vars(r)

	var test Test
	if err = json.NewDecoder(r.Body).Decode(&test); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	err = UpdateTest(params["id"], test, a.Repository)
	if err != nil && err.Error() == NOT_FOUND {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	if err != nil && err.Error() == VALIDATION_ERROR {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, test)
}

func (a *Api) DeleteTestEndPoint(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No Auth ")
		return
	}
	authorization := cookie.Value
	if !IsAuthorized(authorization) {
		respondWithError(w, http.StatusUnauthorized, "No Auth ")
		return
	}

	params := mux.Vars(r)

	err = DeleteTest(params["id"], a.Repository)
	if err != nil && err.Error() == NOT_FOUND {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusCreated, params["id"] + " deleted")
}

func (a *Api) AllQuestionsEndPoint(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No Auth ")
		return
	}
	authorization := cookie.Value
	if !IsAuthorized(authorization) {
		respondWithError(w, http.StatusUnauthorized, "No Auth ")
		return
	}

	questions, err := QuestionFindAll(a.Repository)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, questions)
}

func (a *Api) FindQuestionsEndpoint(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No Auth ")
		return
	}
	authorization := cookie.Value
	if !IsAuthorized(authorization) {
		respondWithError(w, http.StatusUnauthorized, "No Auth ")
		return
	}

	params := mux.Vars(r)
	question, err := QuestionFindById(params["id"], a.Repository)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid question ID")
		return
	}
	respondWithJson(w, http.StatusOK, question)
}

func (a *Api) CreateQuestionsEndPoint(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No Auth ")
		return
	}
	authorization := cookie.Value
	if !IsAuthorized(authorization) {
		respondWithError(w, http.StatusUnauthorized, "No Auth ")
		return
	}

	var question Question
	if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	question.ID = bson.NewObjectId()
	if err := InsertQuestion(question, a.Repository); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, question)
}

func (a *Api) UpdateQuestionsEndPoint(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No Auth ")
		return
	}
	authorization := cookie.Value
	if !IsAuthorized(authorization) {
		respondWithError(w, http.StatusUnauthorized, "No Auth ")
		return
	}

	params := mux.Vars(r)

	var question Question
	if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	err = UpdateQuestion(params["id"], question, a.Repository)
	if err != nil && err.Error() == NOT_FOUND {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	if err != nil && err.Error() == VALIDATION_ERROR {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, question)
}

func (a *Api) DeleteQuestionsEndPoint(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No Auth ")
		return
	}
	authorization := cookie.Value
	if !IsAuthorized(authorization) {
		respondWithError(w, http.StatusUnauthorized, "No Auth ")
		return
	}

	params := mux.Vars(r)

	err = DeleteQuestion(params["id"], a.Repository)
	if err != nil && err.Error() == NOT_FOUND {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusCreated, params["id"] + " deleted")
}

func (a *Api) AllUsersEndPoint(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No Auth ")
		return
	}
	authorization := cookie.Value
	if !IsAuthorized(authorization) {
		respondWithError(w, http.StatusUnauthorized, "No Auth ")
		return
	}

	users, err := UserFindAll(a.Repository)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, users)
}

func (a *Api) FindUserEndpoint(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No Auth ")
		return
	}
	authorization := cookie.Value
	if !IsAuthorized(authorization) {
		respondWithError(w, http.StatusUnauthorized, "No Auth ")
		return
	}

	params := mux.Vars(r)
	user, err := UserFindById(params["id"], a.Repository)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid question ID")
		return
	}
	respondWithJson(w, http.StatusOK, user)
}


func IsAuthorized(authorization string) bool {
	if len(authorization) == 0 {
		return false
	}
	if authorized := ValidateToken(authorization); !authorized {
		return false
	}

	return true
}

func ValidateToken(hash string) bool {
	hash = strings.Replace(hash, "Bearer ", "", -1)
	token, err := jwt.Parse(hash, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		hmacSecret := []byte("agVrqfQK9CiQZceW")
		return hmacSecret, nil
	})

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true
	}

	log.Errorf("API - ValidateToken - Error: %s", err)
	return false
}


func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (app *Api) ConfigureRoutes(router *mux.Router) {
	/* test end point */
	router.HandleFunc("/users", app.AllUsersEndPoint).Methods("get")
	router.HandleFunc("/users/{id}", app.FindUserEndpoint).Methods("GET")
	/* test end point */
	router.HandleFunc("/tests", app.AllTestEndPoint).Methods("get")
	router.HandleFunc("/tests/{id}", app.FindTestEndpoint).Methods("GET")
	router.HandleFunc("/tests", app.CreateTestEndPoint).Methods("post")
	router.HandleFunc("/tests/{id}", app.UpdateTestEndPoint).Methods("put")
	router.HandleFunc("/tests/{id}", app.DeleteTestEndPoint).Methods("delete")
	/* questions end point */
	router.HandleFunc("/questions", app.AllQuestionsEndPoint).Methods("get")
	router.HandleFunc("/questions/{id}", app.FindQuestionsEndpoint).Methods("GET")
	router.HandleFunc("/questions", app.CreateQuestionsEndPoint).Methods("post")
	router.HandleFunc("/questions/{id}", app.UpdateQuestionsEndPoint).Methods("put")
	router.HandleFunc("/questions/{id}", app.DeleteQuestionsEndPoint).Methods("delete")
	router.HandleFunc("/login", app.LoginEndPoint).Methods("post")
}

func InitServer() {
	config.Read()
	DB_NAME = config.Database
	repository := DatabseInit(config)
	defer repository.Close()

	app := Api{
		Repository: repository,
	}

	mux := mux.NewRouter()
	app.ConfigureRoutes(mux)

	server := negroni.New(negroni.NewRecovery())
	server.UseHandler(mux)

	serverAddr := ":" + config.Port
	server.Run(serverAddr)
}

func DatabseInit(config Config) Repository {
	session, err := mgo.Dial(config.Server)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	return session
}
