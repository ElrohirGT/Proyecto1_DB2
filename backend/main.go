package main

import (
	"context"
	"net/http"

	"github.com/ElrohirGT/Proyecto1_DB2/api"
	mw "github.com/ElrohirGT/Proyecto1_DB2/api/middlewares"
	"github.com/ElrohirGT/Proyecto1_DB2/config"
	"github.com/ElrohirGT/Proyecto1_DB2/db_client"
	"github.com/ElrohirGT/Proyecto1_DB2/utils"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func main() {

	// Loading env variables from .env file
	godotenv.Load()

	config := config.LoadConfig()

	// Configure Logger
	utils.ConfigureLogger()

	// Database Client
	dbClient, closeDB, err := db_client.NewDriver(&config.DatabaseConfig)

	if err != nil {
		log.Panic().Err(err).Msg("Failed to create Neo4j driver: %v")
	}
	defer closeDB(context.Background())

	if err != nil {
		log.Fatal().
			Str("message", "Could not initialize DB Client").
			Err(err)
	}

	// App and Services Configuration
	app := api.NewApi(dbClient)

	// Routes
	r := chi.NewRouter()

	r.Use(mw.Logging)
	r.Use(mw.CreateCors(config.CorsConfig))

	r.Get("/health", app.CheckHealthHandler)
	
	r.Group(func(r chi.Router) {

		// Nodes
		r.Post("/node", app.CreateNodeHandler)
		r.Get("/node", app.ReadNodeHandler)
		r.Put("/node", app.UpdateNodeHandler)
		r.Delete("/node", app.DeleteNodeHandler)

		// Relations
		r.Post("/relation", app.CreateUserHandler)
		r.Get("/relation", app.ReadRelationHandler)
		r.Put("/relation", app.UpdateRelationHandler)
		r.Delete("/relation", app.CreateUserHandler)

	})

	// Start server
	log.Printf("Starting server on port %s", config.APIPort)
	log.Fatal().Err(http.ListenAndServe(":"+config.APIPort, r))
}
