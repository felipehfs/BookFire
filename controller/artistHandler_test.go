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

func TestUpdateArtist(t *testing.T) {
	artists, err := getArtists()
	if err != nil {
		t.Error(err)
	}

	last := len(artists) - 1
	if last < 0 {
		t.Errorf("empty list")
	}
	updated := artists[last]
	updated.Name = "Alter"
	updated.Email = "alterado@gmail.com"

	newURL := url + updated.ID.Hex()

	body, err := json.Marshal(updated)
	if err != nil {
		t.Error(err)
	}
	req, err := http.NewRequest("PUT", newURL, bytes.NewBuffer(body))

	if err != nil {
		t.Error(err)
	}
	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Status expected %v got %v", http.StatusOK, res.StatusCode)
	}
}

func TestDeleteArtist(t *testing.T) {
	artists, err := getArtists()
	if err != nil {
		t.Error(err)
	}
	if len(artists) < 0 {
		t.Error("empty list")
	}
	last := artists[len(artists)-1]
	tt := []struct {
		id           string
		expectedCode int
	}{
		{last.ID.Hex(), 200},
		{"5b608966767bb623cf13c303", 500},
	}
	for _, tr := range tt {
		req, err := http.NewRequest("DELETE", url+tr.id, nil)
		if err != nil {
			t.Error(err)
		}
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			t.Error(err)
		}

		if res.StatusCode != tr.expectedCode {
			t.Errorf("Error: expected %v got %v\n", tr.expectedCode, res.StatusCode)
		}
	}
}

func TestUpdateManyCase(t *testing.T) {
	tt := []struct {
		name           string
		id             string
		body           string
		statusExpected int
	}{
		{"Update which exists",
			"5b6057e8767bb623cf13c305",
			`{"name": "Belchior updated", "email": "belchiorupdated@gmail.com"}`,
			200},
		{"Update which id not exists",
			"5b6057e8767bb423cf13c305",
			`{"name": "Belchior", "email": "belchior@gmail.com"}`,
			500},
		{
			"Update which id returns error",
			"5b6057e8767bb623cf13c305",
			`{ "_id": 5b6057e8767bb623cf13c303, 
				"name": "Belchior error", "email": "belchior@gmail.com",
				"description": "fjfjs"}`,
			500},
		{
			"Update which values the same",
			"5b47efb6391bc8b26eb2deac",
			`{"name": "Ken thompson", "email": "kenthompson@gmail.com", "description": "O cara que criou o grande C"}`,
			200},
	}

	for _, tr := range tt {
		req, err := http.NewRequest("PUT", url+tr.id, bytes.NewBuffer([]byte(tr.body)))
		if err != nil {
			t.Error(err)
		}

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != tr.statusExpected {
			t.Errorf("Error: expected %d got %d\n", tr.statusExpected, resp.StatusCode)
		}
	}
}
