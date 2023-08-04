package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Vote struct {
	VoterID     string
	CandidateID string
}

type Status struct {
	Code    int
	Message string
}

func serveRoot(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		candidates := getCandidatesVote()
		totalVotes := 0
		for _, votes := range candidates {
			totalVotes += votes
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"Results": func() []map[string]interface{} {
				results := make([]map[string]interface{}, 0, len(candidates))
				for candidate, votes := range candidates {
					results = append(results, map[string]interface{}{
						"CandidateID": candidate,
						"Votes":       votes,
					})
				}
				return results
			}(),
			"TotalVotes": totalVotes,
		})
	case "POST":
		var vote Vote
		if err := json.NewDecoder(r.Body).Decode(&vote); err != nil {
			writeVoterResponse(w, Status{Code: 400, Message: "Bad Request. Vote can not be saved"})
			return
		}
		if err := saveVote(vote); err != nil {
			writeVoterResponse(w, Status{Code: 500, Message: "Internal Server Error. Vote can not be saved"})
			return
		}
		writeVoterResponse(w, Status{Code: 201, Message: "Vote saved sucessfully"})
	default:
		writeVoterResponse(w, Status{Code: 405, Message: "Method Not Allowed. Vote can not be saved"})
	}
}

// TODO: Replace with actual implementation of these functions
func getCandidatesVote() map[string]int {
	return map[string]int{
		"candidate1": 10,
		"candidate2": 20,
	}
}

func saveVote(vote Vote) error {
	return nil
}

func writeVoterResponse(w http.ResponseWriter, status Status) {
	w.WriteHeader(status.Code)
	json.NewEncoder(w).Encode(status)
}

// TestServeRoot_e6109c0b6f tests the serveRoot function
func TestServeRoot_e6109c0b6f(t *testing.T) {
	// Test GET method
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(serveRoot)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"Results":[{"CandidateID":"candidate1","Votes":10},{"CandidateID":"candidate2","Votes":20}],"TotalVotes":30}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// Test POST method
	vote := Vote{
		VoterID:     "voter1",
		CandidateID: "candidate1",
	}
	voteJSON, _ := json.Marshal(vote)
	req, err = http.NewRequest("POST", "/", bytes.NewBuffer(voteJSON))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(serveRoot)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	expected = `{"Code":201,"Message":"Vote saved sucessfully"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// Test unsupported method
	req, err = http.NewRequest("PUT", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(serveRoot)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}

	expected = `{"Code":405,"Message":"Method Not Allowed. Vote can not be saved"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
