package service

import (
	"context"
	"time"

	"github.com/souls-syntax/Templates/internal/utils"
	"github.com/souls-syntax/Templates/internal/models"
	"github.com/souls-syntax/Templates/internal/cache"
)


type Verifier struct {
	Cache	*cache.RedisCache
	Bert	*BertClient
}

func NewVerifier(c *cache.RedisCache,b *BertClient) *Verifier {
	return &Verifier{
		Cache: c,
		Bert:  b,
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
	
	var dec models.Decision
	var obs models.Observation

	// Cache hit then
	if hit {

		//Create decision
		dec = models.Decision{
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

	dec, err := v.Bert.GetVerdict(ctx,queryText)
	if err != nil {
		return models.VerifyResponse{}, err
	}


	//-----------------Mock Response------------
	dec.QueryHash = hash
	dec.QueryText = queryText

	obs = models.Observation{
		Source:	"BERT-Python",
		ProcessingTimeMs:	time.Since(start).Milliseconds(),
	}
	//----------------Mock Response--------------
	return BuildResponse(dec, obs), nil
	// Cache miss then the logic
}
