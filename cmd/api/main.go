package main

import (
	"os"
	"net/http"
	"log"
	"github.com/joho/godotenv"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/souls-syntax/Templates/internal/handlers"
	"github.com/redis/go-redis/v9"
	"github.com/souls-syntax/Templates/internal/cache"
	"github.com/souls-syntax/Templates/internal/service"
)

func main() {
	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	redisUrl := os.Getenv("REDIS_URL")
	bertUrl := os.Getenv("BERT_URL")

	if redisUrl == "" {
		log.Fatal("Redis url not set")
	}
	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		log.Fatal(err)
	}
	
	client := redis.NewClient(opt)
	myCache := cache.NewRedisCache(client)
	
	bertClient := service.NewBertClient(bertUrl)

	myVerifier := service.NewVerifier(myCache, bertClient)

	apiCfg := &handlers.ApiConfig{
		Verifier: myVerifier,
	}

	router := chi.NewRouter()
	
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:				[]string{"https://*","http://*"},
		AllowedMethods:				[]string{"GET","POST","PUT","DELETE","OPTIONS"},
		AllowedHeaders:				[]string{"*"},
		ExposedHeaders:				[]string{"Link"},
		AllowCredentials:			false,
		MaxAge:								300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/healthz",handlers.HandlerReadiness)
	v1Router.Get("/err",handlers.HandlerErr)

	v1Router.Post("/verify",apiCfg.HandlerVerify)

	router.Mount("/v1",v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:		 ":"+portString,
	}
	log.Printf("Server Starting on port %v", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
