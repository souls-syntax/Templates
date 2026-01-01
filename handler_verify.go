package main

import (
	"net/http"
	"encoding/json"
)

type Query struct {
	Query			string	`json:"query"`
	UserID		string	`json:"userid"`
}

type Response struct {
	Response		string
	Status			bool
}

func handlerVerify(w http.ResponseWriter, r *http.Request) {
	var query Query
	
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&query); err != nil {
		respondWithError(w,401,"Incorrect format of JSON query")
	}
	// Will normalize it later till normal return
	response := Response{
		Response: "Hey POST seems to be working",
		Status: true,
	}
	respondWithJson(w,201,&response)
}
