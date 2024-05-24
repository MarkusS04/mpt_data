package auth

import (
	"bytes"
	"encoding/json"
	"mpt_data/database"
	"mpt_data/test/vars"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMain(m *testing.M) {
	if err := database.Connect(vars.GetDbPAth()); err != nil {
		panic(err)
	}
	m.Run()
}

func TestLoginHandler(t *testing.T) {
	// Create a sample user login payload
	payload, err := json.Marshal(vars.UserAPI)
	if err != nil {
		t.Error("error test prep:", err)
		return
	}

	// Create a mock request with the payload
	req, err := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the login handler with the mock request
	handler := http.HandlerFunc(login)
	handler.ServeHTTP(rr, req)

	// Check the response status code is what you expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response headers for the Authorization header
	authHeader := rr.Header().Get("Authorization")
	if authHeader == "" {
		t.Error("Authorization header not set in response")
	}
}
