package api_test

import (
	"bytes"
	"encoding/json"
	"mpt_data/api/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

type RequestData struct {
	Data   interface{}
	Route  string
	Method string
	Router func(http.ResponseWriter, *http.Request)
	Path   string
}

func DoRequest(t *testing.T, reqData RequestData) *httptest.ResponseRecorder {
	payload, err := json.Marshal(reqData.Data)
	if err != nil {
		t.Fatal("error test prep:", err)
		return nil
	}

	req, err := http.NewRequest(reqData.Method, reqData.Route, bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	// rollback flag to true
	req = req.WithContext(middleware.SetRollback(req.Context(), true))
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.Use(middleware.TransactionMiddleware)
	router.HandleFunc(reqData.Path, reqData.Router)
	router.ServeHTTP(rr, req)

	return rr
}
