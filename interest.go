package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"os"
	"runtime"

	"github.com/gorilla/mux"
)

func interestHandler(si SetInterest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		decoder := json.NewDecoder(r.Body)
		encoder := json.NewEncoder(w)

		err := decoder.Decode(&si)
		if err != nil {
			encoder.Encode(Error{"Error decoding json"})
			return
		}

		si.SetInterest()

		err = encoder.Encode(si)
		if err != nil {
			//TODO: using the same encoder instance. is it ok?
			encoder.Encode(Error{"Error encoding json"})
		}
	}
}

var simpleInterestHandler http.HandlerFunc = interestHandler(&SimpleInterest{})
var compoundInterestHandler http.HandlerFunc = interestHandler(&CompoundInterest{})

type Error struct {
	Err string `json:"error"`
}

type SetInterest interface {
	SetInterest()
}

type SimpleInterestRequest struct {
	Principal float64 `json:"principal"`
	Interest  float64 `json:"interest"`
	Period    float64 `json:"period"`
}

type SimpleInterest struct {
	*SimpleInterestRequest

	SimpleInterest float64 `json:"simpleInterest"`
}

func (self *SimpleInterest) SetInterest() {
	self.SimpleInterest = self.Principal * (self.Interest * 1e-2) * self.Period
}

type CompoundInterestRequest struct {
	*SimpleInterestRequest
	Frequency float64 `json:"frequency"`
}

type CompoundInterest struct {
	*CompoundInterestRequest

	CompoundInterest float64 `json:"compoundInterest"`
}

func (self *CompoundInterest) SetInterest() {
	if self.Frequency == 0 {
		return
	}
	s := self.Principal * math.Pow(1+self.Interest*1e-2/self.Frequency, self.Frequency*self.Period)
	self.CompoundInterest = s - self.Principal
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/simpleInterest", simpleInterestHandler).Methods("POST")
	r.HandleFunc("/compoundInterest", compoundInterestHandler).Methods("POST")
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
