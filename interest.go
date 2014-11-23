package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/gorilla/mux"
)

func simpleInterestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	encoder := json.NewEncoder(w)

	var si SimpleInterest

	err := decoder.Decode(&si)
	if err != nil {
		encoder.Encode(Error{"Error decoding json"})
		return
	}

	si.SetSimpleInterest()

	err = encoder.Encode(si)
	if err != nil {
		//TODO: using the same encoder instance. is it ok?
		encoder.Encode(Error{"Error encoding json"})
	}
}

type Error struct {
	Err string `json:"error"`
}

type SimpleInterest struct {
	Principal float64 `json:"principal"`
	Interest  float64 `json:"interest"`
	Period    float64 `json:"period"`

	SimpleInterest float64 `json:"simpleInterest"`
}

func (self *SimpleInterest) SetSimpleInterest() {
	self.SimpleInterest = self.Principal * (self.Interest * 1e-2) * self.Period
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/simpleInterest", simpleInterestHandler).Methods("POST")
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
