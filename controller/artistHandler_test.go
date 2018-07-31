// package controller test the controller
package controller

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
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
	if resp.Status != "200" {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		log.Println(string(body))
	}
}
