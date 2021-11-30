package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

type Payload struct {
	Fid string `json:"fid"`
	Src string `json:"src"`
	Params string `json:"params"`
	Lang string `json:"lang"`
}
var Payloads []Payload

func homePage(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/showfunctions", returnAllFunctions)
	myRouter.HandleFunc("/showfunctions/{fid}", returnFunction)
	myRouter.HandleFunc("/run", runFunction).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", myRouter))
}
func returnFunction(w http.ResponseWriter, r *http.Request){
	fmt.Println("Endpoint Hit: returnFunction")
	vars := mux.Vars(r)
	key := vars["fid"]
	for _, payload := range Payloads{
		if payload.Fid == key {
			json.NewEncoder(w).Encode(payload)
		}
	}
}
func runFunction(w http.ResponseWriter, r *http.Request){
	reqBody, _ := ioutil.ReadAll(r.Body)
	var payload Payload
	json.Unmarshal(reqBody, &payload)
	json.NewEncoder(w).Encode(payload.Params)
	//fmt.Fprintf(w, "%+v", string(reqBody))
}
func returnAllFunctions(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllFunctions")
	json.NewEncoder(w).Encode(Payloads)
}
func main()  {
	fmt.Println("Rest API v2.0 - Mux Routers")
	Payloads = []Payload{
		{Fid: "helloworld", Src: "print(\"hello\")", Lang: "Python", Params: "Yanni"},
		{Fid: "helloworld2", Src: "print(\"hello\")", Lang: "Python", Params: "Yanni"},
	}
	handleRequests()
}