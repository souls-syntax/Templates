package handlers

import (
	"net/http"
	"time"
	"encoding/json"
	"context"
	"github.com/souls-syntax/Templates/internal/utils"
	"github.com/souls-syntax/Templates/internal/models"
	"github.com/souls-syntax/Templates/internal/service"
)

func (c *ApiConfig) HandlerVerify(w http.ResponseWriter, r *http.Request) {
	var req models.Query
	
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		utils.RespondWithError(w,401,"Incorrect format of JSON query")
		return
	}
	
	// Normalizing and Hashing
	normal := utils.NormalizeQuery(req.Query)
	hash := utils.HashQuery(normal)
	
	// Start the stopwatch
	start := time.Now()

	// Trying the cache out
	ctx := r.Context()
	cacheVal, hit := c.Cache.Get(ctx, hash)
	
	if hit {

		//Create decision
		dec := models.Decision{
			QueryHash:	hash,
			QueryText:	req.Query,
			Verdict:	cacheVal.Verdict,
			Confidence:	cacheVal.Confidence,
			Decider:	cacheVal.Decider,
		}

		//Create observations
		obs := models.Observations{
			Source:	"Redis",
			ProcessingTimeMs:	time.Since(start).Milliseconds()
		}
		
		resp := service.BuildResponse(dec,obs)

		utils.RespondWithJson(w, 200, resp)
		return
	}


	res := models.Response{
		Response: "Hey POST seems to be working",
		Status: true,
	}


	utils.RespondWithJson(w,201,&res)
}
