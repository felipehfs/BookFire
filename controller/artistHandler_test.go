// package controller test the controller
package controller

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

const (
	url = "http://localhost:8081/artists/"
)

func TestCreateArtist(t *testing.T) {
	var request = struct {
		Name        string
		Email       string
		Description string
	}{
		Name:  "Belchior",
		Email: "belchior@ig.com.br",
		Description: `Apenas um latino americano, 
			sem conta no banco, sem parentes importantes e vindo do interior`,
	}

	// Convert to bytes the body of request
	jsonData, err := json.Marshal(request)
	if err != nil {
		t.Error(err)
	}

	client := &http.Client{}
	// setup the request to server
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Error(err)
	}
	// make the request to server
	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}
	// If the status code is not ok the error must be send
	checkRequestOK(resp, t)
}

func checkRequestOK(response *http.Response, t *testing.T) {
	if response.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}
		t.Errorf("Status code different received: %s %s", response.Status, string(body))
	}
}

func TestReadArtists(t *testing.T) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Error(err)
	}
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		t.Error(err)
	}

	checkRequestOK(response, t)
}
