package api

import (
	"fmt"
	"mpt_data/api/auth"
	"mpt_data/api/meeting"
	"mpt_data/api/meeting/absencemeeting"
	"mpt_data/api/person"
	"mpt_data/api/person/absenceperson"
	"mpt_data/api/plan"
	"mpt_data/api/task"
	"mpt_data/helper/config"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func PrepareServer() http.Handler {
	mux := mux.NewRouter()
	registerRoutes(mux)
	return corsHandler().Handler(mux)
}

func registerRoutes(mux *mux.Router) {
	if config.Config.API.UseSwagger {
		fmt.Println("swagger is running")
		initSwagger(mux)
	}

	if config.Config.API.AuthenticationRequired {
		auth.RegisterRoutes(mux)
	}
	meeting.RegisterRoutes(mux)
	task.RegisterRoutes(mux)
	person.RegisterRoutes(mux)
	absencemeeting.RegisterRoutes(mux)
	absenceperson.RegisterRoutes(mux)
	plan.RegisterRoutes(mux)
}

func corsHandler() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions, http.MethodPatch},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		ExposedHeaders: []string{"Authorization"},
	})
}
