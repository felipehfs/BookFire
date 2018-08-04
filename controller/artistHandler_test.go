// package controller test the controller
package controller

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/bookfire/model"
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
	artists, err := getArtists()

	if err != nil {
		t.Error(err)
	}
	// Check the persistence of the data
	last := len(artists) - 1
	equals := artists[last].Name == request.Name && artists[last].Email == request.Email
	if !equals {
		t.Errorf("Values are different: expected %s when %s", request.Name, artists[0].Name)
	}
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
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		t.Error(err)
	}

	checkRequestOK(response, t)
}

// getArtists retrieves all data from the server
// if exists
func getArtists() ([]model.Artist, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	var artists []model.Artist
	body, err := ioutil.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &artists); err != nil {
		return nil, err
	}

	return artists, nil
}

func TestFindArtists(t *testing.T) {
	searched, err := getArtists()
	if err != nil {
		t.Error(err)
	}

	tt := []struct {
		id       string
		name     string
		expected int
	}{
		{searched[0].ID.String(), "busca que funciona", 200},
		{"5b608966767bb623cf13c303", "busca que falha", 404},
	}

	for _, tr := range tt {
		t.Run(tr.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", url+tr.id, nil)
			if err != nil {
				t.Fail()
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Fail()
			}

			if resp.StatusCode != tr.expected {
				t.Errorf("Status code expected %d received %d", tr.expected, resp.StatusCode)
			}
		})
	}
}
