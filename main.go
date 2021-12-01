package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

type Payload struct {
	Fid    string `json:"fid"`
	Src    string `json:"src"`
	Params string `json:"params"`
	Lang   string `json:"lang"`
}

var Payloads []Payload

func homePage(w http.ResponseWriter, r *http.Request) {
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
func returnFunction(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnFunction")
	vars := mux.Vars(r)
	key := vars["fid"]
	for _, payload := range Payloads {
		if payload.Fid == key {
			json.NewEncoder(w).Encode(payload)
		}
	}
}
func runFunction(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var payload Payload
	json.Unmarshal(reqBody, &payload)
	json.NewEncoder(w).Encode(payload.Params)
	params, _ := json.Marshal(payload.Params)
	src, _ := json.Marshal(payload.Src)
	cmd := exec.Command("python", "run_script.py", string(src), string(params))
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
	//fmt.Fprintf(w, "%+v", string(reqBody))
}
func returnAllFunctions(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllFunctions")
	json.NewEncoder(w).Encode(Payloads)
}
func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	Payloads = []Payload{
		{Fid: "helloworld", Src: "print(\"hello\")", Lang: "Python", Params: "{\"A\":\"3\", \"B\":\"5\"}"},
		{Fid: "helloworld2", Src: "print(\"hello\")", Lang: "Python", Params: "{\"A\":\"4\", \"B\":\"6\"}"},
	}
	handleRequests()
}