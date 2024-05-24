package api

import (
	"github.com/gorilla/mux"

	_ "mpt_data/docs" // Import the Swag-generated docs package

	httpSwagger "github.com/swaggo/http-swagger"
)

func initSwagger(mux *mux.Router) {
	mux.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // The url pointing to API definition
		httpSwagger.Layout(httpSwagger.BaseLayout),
	))
}
