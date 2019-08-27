package servicev1

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/golang/protobuf/jsonpb"
	"github.com/gorilla/mux"
	v1 "github.com/nicolasassi/starWarsApi/pkg/cmd/db/v1"
)

var connection *v1.Mongo

func Serve(dbm *v1.Mongo) {
	r := mux.NewRouter()
	connection = dbm
	r.HandleFunc("/api/v1/planet/{id}", GetByID).Methods("GET")
	r.HandleFunc("/api/v1/planet/name/{name}", GetByName).Methods("GET")
	r.HandleFunc("/api/v1/planet", Add).Methods("POST")
	r.HandleFunc("/api/v1/planet/{id}", Delete).Methods("DELETE")

	var port = ":3000"
	fmt.Println("Server running in port:", port)
	log.Fatal(http.ListenAndServe(port, r))
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}

func getInt(id string) (int, error) {
	i, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func GetByName(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	planet, err := connection.FindByName(params["name"])
	log.Println(params["name"])
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, planet)
}

func GetByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID, err := getInt(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid planet ID")
	}
	planet, err := connection.FindByID(ID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid planet ID")
		return
	}
	respondWithJSON(w, http.StatusOK, planet)
}

func Add(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var planet v1.Planet
	if err := jsonpb.UnmarshalNext(json.NewDecoder(r.Body), &planet); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := connection.AddPlanet(planet); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, planet)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	params := mux.Vars(r)
	ID, err := getInt(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid planet ID")
	}
	if err := connection.DeletePlanet(ID); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
