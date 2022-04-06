package main

import (
	"encoding/json"
	b64 "encoding/base64"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
	"os"
	"strconv"
	"net"
)

type Payload struct {
	Fid    string `json:"fid"`
	Src    string `json:"src"`
	Params string `json:"params,omitempty"`
	Lang   string `json:"lang"`
	Worker string `json:"worker"`
}
type Response struct {
	Fid         string        `json:"fid"`
	TimeElapsed time.Duration `json:"time_elapsed"`
	Result      string        `json:"result"`
	Error       string        `json:"error,omitempty"`
}

var Payloads []Payload
var last_ready = time.Now()
var job_received = false

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func updatePoweroff(){
		client := http.Client{
		Timeout: 5 * time.Second,
	}
	fmt.Println("Sending poweroff notification.")
	_, err := client.Get(PROVIDER_IP + "?poweroff=1")

	if (err == nil){
		fmt.Println("Received poweroff acknowledgement.")
		return
	} else {
		updatePoweroff()
	}
}

func poweroff(){
	for {
		if(job_received == false && time.Since(last_ready) >= 15 * time.Second){
			//tell provider that i am turning off
			updatePoweroff()
			
			//turn off
			lol := exec.Command("poweroff")
		
			p_err:= lol.Run()
			if p_err != nil {
				fmt.Println(p_err)
			}

		}
	}
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/showfunctions", returnAllFunctions)
	myRouter.HandleFunc("/showfunctions/{fid}", returnFunction)
	myRouter.HandleFunc("/run", runFunction).Methods("POST")
	// updateHealth()
	// log.Fatal(http.ListenAndServe(":8080", myRouter))

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Error listening")
	}

	// Signal that server is open for business. 
	updateHealth()
	go poweroff()
	if err := http.Serve(l, myRouter); err != nil {
		// handle error
		log.Fatal("Error serving")

	}
	

	
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
	job_received = true
	reqBody, _ := ioutil.ReadAll(r.Body)
	var payload Payload
	res := Response{}
	json.Unmarshal(reqBody, &payload)
	params:= payload.Params
	// src:= payload.Src
	data, _ := b64.StdEncoding.DecodeString(payload.Src)
	s2, _ := strconv.Unquote(string(data))
	// fmt.Println(s2)
	err := os.WriteFile("/root/src.py",[]byte(s2), 0644)
	if err != nil {
		panic(err)
	}
	
	if strings.ToLower(payload.Lang) == "micropython"{
		//fmt.Println(string(out))	
		start := time.Now()
		cmd := exec.Command( "micropython",  "/root/src.py", params)
		out, err := cmd.CombinedOutput()
		t := time.Now()
		fmt.Println(string(out))

		if err != nil {
			res = Response{Fid: payload.Fid, TimeElapsed: t.Sub(start), Result: string(out), Error: err.Error()}
			fmt.Println(err)
		} else {
			res = Response{Fid: payload.Fid, TimeElapsed: t.Sub(start), Result: string(out)}
		}
					// // Testing out rebooting after a function
	// lol := exec.Command("echo", "b", ">", "/proc/sysrq-trigger")

	} else if strings.ToLower(payload.Lang) == "python" {
		start := time.Now()
		cmd := exec.Command("python", "/root/src.py",  params)
		out, err := cmd.CombinedOutput()
		t := time.Now()
		fmt.Println(string(out))
		if err != nil {
			res = Response{Fid: payload.Fid, TimeElapsed: t.Sub(start), Result: string(out), Error: err.Error()}
			fmt.Println(err)
		} else {
			res = Response{Fid: payload.Fid, TimeElapsed: t.Sub(start), Result: string(out)}
		}
	} else {
		res = Response{Fid: payload.Fid, Result: "Language not supported"}
	}

	json.NewEncoder(w).Encode(res)
	fmt.Println("rebooting")

	// lol := exec.Command("sh", "-c", "echo b > /proc/sysrq-trigger")

	lol := exec.Command("reboot")
	
	reboot_err:= lol.Run()
	if reboot_err != nil {
		fmt.Println(err)
	}



	//fmt.Fprintf(w, "%+v", string(reqBody))
}
func returnAllFunctions(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllFunctions")
	json.NewEncoder(w).Encode(Payloads)
}


func updateHealth(){
		client := http.Client{
		Timeout: 5 * time.Second,
	}
	fmt.Println("Sending health response.")
	_, err := client.Get(PROVIDER_IP)
	
	if (err == nil){
		fmt.Println("Received health response.")
		last_ready = time.Now()
		// fmt.Println(resp)
		return
	} else {
		updateHealth()
	}
}

var PROVIDER_IP = "http://192.168.1.203:8083/healthz"

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	Payloads = []Payload{
		{Fid: "helloworld", Src: "print(\"hello\")", Lang: "Python", Params: "{\"A\":\"3\", \"B\":\"5\"}"},
		{Fid: "helloworld2", Src: "print(\"hello\")", Lang: "Python", Params: "{\"A\":\"4\", \"B\":\"6\"}"},
	}

	// go job_timer();
	handleRequests()

}
