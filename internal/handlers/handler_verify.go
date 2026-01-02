package handlers

import (
	"net/http"
	"encoding/json"
	"github.com/souls-syntax/Templates/internal/utils"
	"github.com/souls-syntax/Templates/internal/models"
)

func HandlerVerify(w http.ResponseWriter, r *http.Request) {
	var query models.Query
	
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&query); err != nil {
		utils.RespondWithError(w,401,"Incorrect format of JSON query")
	}
	// Will normalize it later till normal return
	response := models.Response{
		Response: "Hey POST seems to be working",
		Status: true,
	}
	utils.RespondWithJson(w,201,&response)
}
