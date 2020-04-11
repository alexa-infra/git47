package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetRef(t *testing.T) {
	r := prepareRepository(t)

	testData := map[string]string{
		"master": "60a58ae38710f264b2c00f77c82ae44419381a3f",
		"60a58ae38710f264b2c00f77c82ae44419381a3f": "60a58ae38710f264b2c00f77c82ae44419381a3f",
		"new-branch": "377229569f4a7ae706ed3a376117dabee4cec8f8",
		"v1":         "377229569f4a7ae706ed3a376117dabee4cec8f8",
	}

	for refName, refHash := range testData {
		req, _ := http.NewRequest("GET", "/r/memory/tree/master", nil)
		req = mux.SetURLVars(req, map[string]string{"ref": refName})

		hash, err := getRef(req, r)
		if assert.Nil(t, err) {
			assert.Equal(t, refHash, hash.String(), fmt.Sprintf("Invalid hash for ref %s", refName))
		}
	}
}
