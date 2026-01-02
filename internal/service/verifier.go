package service

import (
	"context"
	"time"
	"errors"
	"log"

	"github.com/souls-syntax/Templates/internal/utils"
	"github.com/souls-syntax/Templates/internal/models"
	"github.com/souls-syntax/Templates/internal/cache"
	"github.com/souls-syntax/Templates/internal/database"
)


type Verifier struct {
	Cache	*cache.RedisCache
	Bert	*BertClient
	DB    *database.Store
}

func NewVerifier(c *cache.RedisCache,b *BertClient, db *database.Store) *Verifier {
	return &Verifier{
		Cache: c,
		Bert:  b,
		DB:		 db,
	}
}


func (v *Verifier) Verify(ctx context.Context, queryText string) (models.VerifyResponse, error) {

	// Start the stopwatch
	start := time.Now()


	// Normalizing and Hashing
	normal := utils.NormalizeQuery(queryText)
	hash := utils.HashQuery(normal)
	
	// Trying the cache out
	cacheVal, hit := v.Cache.Get(ctx, hash)
	
	// Cache hit then
	if hit {

		//Create decision
		dec := models.Decision{
			QueryHash:	hash,
			QueryText:	queryText,
			Verdict:	cacheVal.Verdict,
			Confidence:	cacheVal.Confidence,
			Decider:	cacheVal.Decider,
		}

		//Create observations
		obs := models.Observation{
			Source:	"Redis",
			ProcessingTimeMs:	time.Since(start).Milliseconds(),
		}
		
		resp := BuildResponse(dec,obs)
		return resp, nil
	}

	type raceResult struct {
		dec models.Decision
		src string
		err error
	}

	resultChan := make(chan raceResult, 2)

	go func() {
		dec, err := v.DB.GetDecision(hash)
		if err == nil {
			resultChan <- raceResult{dec: dec, src: "Postgres-History", err: nil}
		} else {
			resultChan <- raceResult{err: errors.New("db miss")}
		}
	}()

	go func() {
		dec, err := v.Bert.GetVerdict(ctx, queryText)
		if err == nil {
			dec.QueryHash = hash
			dec.QueryText = queryText
			resultChan <- raceResult{dec: dec, src: "BERT-Python", err:nil}
		} else {
			resultChan <- raceResult{err: err}
		}
	}()


	var finalDecision models.Decision
	var finalSource string
	
	// such big app first for loop
	for i := 0; i < 2; i++ {
		res := <- resultChan
		if res.err == nil {
			finalDecision = res.dec
			finalSource = res.src
			break
		}
	}


	if finalDecision.Verdict == "" {
		return models.VerifyResponse{}, errors.New("Intelligence failure")
	}
	
	obs := models.Observation{
		Source:	finalSource,
		ProcessingTimeMs:	time.Since(start).Milliseconds(),
	}
	
	go func() {

		bgCtx := context.Background()

		err := v.Cache.Set(bgCtx, hash, finalDecision.ToCacheModel(), 10*time.Minute)
		if err != nil {
			log.Printf("Failed to save to cache Error: %v", err)
		}

		if finalSource != "Postgres-History" && (finalDecision.Verdict == "likely_false" || finalDecision.Verdict == "False") && finalDecision.Confidence > 0.90 {
			v.DB.SaveDecision(finalDecision, finalSource)
			log.Printf("Saved in database")
		}
	}()

	return BuildResponse(finalDecision, obs), nil
}
