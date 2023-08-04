package main

import (
	"sync"
	"testing"
)

var once sync.Once
var candidateVotesStore map[string]int

func getCandidatesVote() map[string]int {
	once.Do(func() {
		candidateVotesStore = make(map[string]int)
	})
	return candidateVotesStore
}

func TestGetCandidatesVote_152da079ca(t *testing.T) {
	votes := getCandidatesVote()
	if len(votes) != 0 {
		t.Error("Expected an empty map, got ", votes)
	}

	votes["John"] = 5
	votes = getCandidatesVote()
	if len(votes) != 1 {
		t.Error("Expected a map with one entry, got ", votes)
	}

	votes["Jane"] = 3
	votes = getCandidatesVote()
	if votes["Jane"] != 3 {
		t.Error("Expected Jane to have 3 votes, got ", votes["Jane"])
	}
}
