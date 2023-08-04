package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Status struct {
	Code    int
	Message string
}

func writeVoterResponse(w http.ResponseWriter, status Status) {
	w.Header().Set("Content-Type", "application/json")
	resp, err := json.Marshal(status)
	if err != nil {
		log.Println("error marshaling response to vote request. error: ", err)
	}
	w.Write(resp)
}

func TestWriteVoterResponse_11f1d592d2(t *testing.T) {
	t.Run("successful response", func(t *testing.T) {
		status := Status{Code: 200, Message: "Success"}
		req, err := http.NewRequest("GET", "/vote", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			writeVoterResponse(w, status)
		})

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		expected := `{"Code":200,"Message":"Success"}`
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expected)
		}
	})

	t.Run("error in marshaling", func(t *testing.T) {
		status := Status{Code: 200, Message: string([]byte{0x80, 0x81, 0x82})}
		req, err := http.NewRequest("GET", "/vote", nil)
		if err != nil {
			t.Fatal(err)
		}

		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer func() {
			log.SetOutput(nil)
		}()

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			writeVoterResponse(w, status)
		})

		handler.ServeHTTP(rr, req)

		if buf.Len() == 0 {
			t.Error("Expected an error log, but did not get one")
		}
	})
}
